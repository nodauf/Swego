package utils

import(
	"strings"
	"container/list"
	"io/ioutil"
	"os"
	"log"
	"bytes"
	"io"

	"github.com/yeka/zip"
)

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

//func AddfiletoEncryptedZip(f *os.File, zipw *zip.Writer, password string) {
func AddfiletoEncryptedZip(path string, f *os.File, zipw *zip.Writer, password string) {
        filePathName := f.Name()

        body, err := ioutil.ReadFile(filePathName)
        if err != nil {
            log.Fatalf("unable to read file: %v", err)
        }

        if err != nil {
            log.Fatalln(err)
        }

        //Create the file to the zip zipw
        w, err := zipw.Encrypt(path, password, zip.StandardEncryption)
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
