package utils

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/yeka/zip"
)

// Dirlisting to manages directory listings
type Dirlisting struct {
	Name          string
	ChildrenDir   []string
	ChildrenFiles []string
	ServerUA      string
	Embedded      bool
}

// CopyToArray convert a list to a string
func CopyToArray(src *list.List) []string {
	dst := make([]string, src.Len())

	i := 0
	for e := src.Front(); e != nil; e = e.Next() {
		dst[i] = e.Value.(string)
		i = i + 1
	}

	return dst
}

// Min returns the minimum between two integer
func Min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// ParseCSV format. Use to parse the header Encoding
func ParseCSV(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	dataTmp := make([]string, len(splitted))

	for i, val := range splitted {
		dataTmp[i] = strings.TrimSpace(val)
	}

	return dataTmp
}

func ParseRange(data string) int64 {
	stop := (int64)(0)
	part := 0
	for i := 0; i < len(data) && part < 2; i = i + 1 {
		if part == 0 {
			if data[i] == '=' {
				part = 1
			}

			continue
		}

		if part == 1 {
			if data[i] == ',' || data[i] == '-' {
				part = 2
			} else {
				if 48 <= data[i] && data[i] <= 57 {
					// ... convert the char to integer and add it!
					stop = (stop * 10) + (((int64)(data[i])) - 48)
				} else {
					part = 2
				}
			}
		}
	}

	return stop
}

//AddfiletoZip will add a file to a zip file
func AddfiletoZip(path string, f *os.File, zipw *zip.Writer, encrypted bool, password string) {
	filePathName := f.Name()

	body, err := ioutil.ReadFile(filePathName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	//Create the file to the zip zipw
	var w io.Writer
	if encrypted {
		w, err = zipw.Encrypt(path, password, zip.StandardEncryption)
	} else {
		w, err = zipw.Create(path)
	}
	if err != nil {
		log.Fatal(err)
	}
	// Copy the data of the local file f into the zip
	_, err = io.Copy(w, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
	return
}

// AddRicefiletoZip add a embedded file (rice file) to a zip
func AddRicefiletoZip(path string, f *rice.File, filePathName string, zipw *zip.Writer, encrypted bool, password string) {
	//        body, err := ioutil.ReadFile(filePathName)
	statInfo, err := f.Stat()
	if err != nil {
		log.Println("500 Internal Error : stat() failure for the file: " + filePathName)
		return
	}
	buf := make([]byte, statInfo.Size())
	var body []byte
	//		n := 0
	for err == nil {
		_, err = f.Read(buf)
		body = append(body, buf...)
		//output_writer.Write(body[0:n])
	}

	//        if err != nil {
	//            log.Fatalf("unable to read file: %v", err)
	//        }

	//Create the file to the zip zipw
	var w io.Writer
	if encrypted {
		w, err = zipw.Encrypt(path, password, zip.StandardEncryption)
	} else {
		w, err = zipw.Create(path)
	}
	if err != nil {
		log.Fatal(err)
	}
	// Copy the data of the local file f into the zip
	_, err = io.Copy(w, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
	return
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ZipDirectory will create a zip file from a directory
func ZipDirectory(f *os.File, encrypted bool) string {
	// Absolute path to the directory
	directoryPathName := f.Name()
	statinfo, _ := f.Stat()

	// Name of the diretory to download
	directoryName := statinfo.Name()

	// Create the zip file
	zipFile, err := os.Create(directoryPathName + "/" + directoryName + ".zip")
	if err != nil {
		log.Fatalln("os.Create: " + err.Error())
	}
	zipFilePath := zipFile.Name()
	zipw := zip.NewWriter(zipFile)

	// Iterate recursively in the folder folderPathName
	err = filepath.Walk(directoryPathName,
		func(path string, info os.FileInfo, err error) error {
			// Take the relative path from the root directory of the web server
			// Windows trick
			if runtime.GOOS == "windows" {
				directoryPathName = strings.Replace(directoryPathName, "/", "\\", -1)
				path = strings.Replace(path, "/", "\\", -1)
			}
			zipPath := directoryName + strings.SplitAfter(path, directoryPathName)[1]

			if err != nil {
				return err
			}

			// Don't add folder and the zip itself
			if !info.IsDir() && info.Name() != directoryName+".zip" {
				// Open the file to zip
				f, err = os.Open(path)

				if err != nil {
					return err
				}

				AddfiletoZip(zipPath, f, zipw, encrypted, "infected")
			}
			return nil
		})
	zipw.Close()
	if err != nil {
		log.Println("walk directory zip: " + err.Error())
	}

	return zipFilePath
}

// SearchAndReplace function in the byte slice
func SearchAndReplace(SearchAndReplaceMap map[string]string, buf []byte) []byte {
	if len(SearchAndReplaceMap) > 0 {
		for searchAndReplaceOld, searchAndReplaceNew := range SearchAndReplaceMap {
			buf = bytes.ReplaceAll(buf, []byte(searchAndReplaceOld), []byte(searchAndReplaceNew))
		}
	}
	return buf
}

// Check if there is an error
func Check(e error, customMessage string) {
	if e != nil {
		fmt.Println(customMessage)
		panic(e)
	}
}
