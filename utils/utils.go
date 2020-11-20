package utils

import "strings"
import "container/list"


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

