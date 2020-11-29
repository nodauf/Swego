package controllers

import (
    "fmt"
    "log"

    clr "github.com/ropnop/go-clr"
)

var TARGET_VERSION = "v2"


func PrintEmbeddedFiles(){
    _, children_files := listEmbeddedFiles()
    for _, value := range children_files{
        fmt.Println(value)
    }
}

func RunEmbeddedBinary(binary string, arguments []string){
    binBytes := readEmbeddedBinary(binary)

    retCode, err := clr.ExecuteByteArray(TARGET_VERSION, binBytes, arguments)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] %s returned exit code: %d\n", binary, retCode)
}
