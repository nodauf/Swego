package utils

import "strings"
import "container/list"


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

