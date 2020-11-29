package utils

import(
	"strings"
	"container/list"
	"io/ioutil"
    "path/filepath"
	"os"
	"log"
	"bytes"
	"io"
    "runtime"

    "github.com/GeertJohan/go.rice"
	"github.com/yeka/zip"
)

const fs_maxbufsize = 4096

type Params struct {
	Name string
	Children_dir []string
	Children_files []string
	Messages	string
}

// Manages directory listings
type Dirlisting struct {
        Name           string
        Children_dir   []string
        Children_files []string
        ServerUA       string
        Embedded       bool
}

func TrimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}

func CopyToArray(src *list.List) []string {
        dst := make([]string, src.Len())

        i := 0
        for e := src.Front(); e != nil; e = e.Next() {
                dst[i] = e.Value.(string)
                i = i + 1
        }

        return dst
}

func Min(x int64, y int64) int64 {
        if x < y {
                return x
        }
        return y
}

func ParseCSV(data string) []string {
        splitted := strings.SplitN(data, ",", -1)

        data_tmp := make([]string, len(splitted))

        for i, val := range splitted {
                data_tmp[i] = strings.TrimSpace(val)
        }

        return data_tmp
}

func ParseRange(data string) int64 {
        stop := (int64)(0)
        part := 0
        for i := 0; i < len(data) && part < 2; i = i + 1 {
                if part == 0 { // part = 0 <=> equal isn't met.
                        if data[i] == '=' {
                                part = 1
                        }

                        continue
                }

                if part == 1 { // part = 1 <=> we've met the equal, parse beginning
                        if data[i] == ',' || data[i] == '-' {
                                part = 2 // part = 2 <=> OK DUDE.
                        } else {
                                if 48 <= data[i] && data[i] <= 57 { // If it's a digit ...
                                        // ... convert the char to integer and add it!
                                        stop = (stop * 10) + (((int64)(data[i])) - 48)
                                } else {
                                        part = 2 // Parsing error! No error needed : 0 = from start.
                                }
                        }
                }
        }

        return stop
}

//func AddfiletoEncryptedZip(path string, f *os.File, zipw *zip.Writer, password string) {
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
        }else{
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
				body = append(body, buf ...)
                //output_writer.Write(body[0:n])
        }

//        if err != nil {
//            log.Fatalf("unable to read file: %v", err)
//        }

        //Create the file to the zip zipw
        var w io.Writer
        if encrypted {
            w, err = zipw.Encrypt(path, password, zip.StandardEncryption)
        }else{
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

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func ZipDirectory(f *os.File, encrypted bool) string{
	// Absolute path to the directory
	directoryPathName := f.Name()
    statinfo, _ := f.Stat()

	// Name of the diretory to download
	directoryName := statinfo.Name()

	// Create the zip file
	zipFile, err := os.Create(directoryPathName+"/"+directoryName+".zip")
	if err != nil {
		log.Fatalln("os.Create: "+err.Error())
	}
	zipFilePath := zipFile.Name()
	zipw := zip.NewWriter(zipFile)

	// Iterate recursively in the folder folderPathName
	err = filepath.Walk(directoryPathName,
		func(path string, info os.FileInfo, err error) error {
			// Take the relative path from the root directory of the web server
            // Windows trick
            if runtime.GOOS == "windows" {
                directoryPathName = strings.Replace(directoryPathName,"/","\\",-1)
                path = strings.Replace(path,"/","\\",-1)
            }
			zipPath := directoryName + strings.SplitAfter(path,directoryPathName)[1]

			if err != nil {
				return err
			}

			// Don't add folder and the zip itself
			if(!info.IsDir() && info.Name() != directoryName+".zip"){
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
		log.Println("walk directory zip: "+err.Error())
	}

    return zipFilePath
}
