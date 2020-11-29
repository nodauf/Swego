package controllers

import (
    "fmt"

)


func PrintEmbeddedFiles(){
    _, children_files := listEmbeddedFiles()
    for _, value := range children_files{
        fmt.Println(value)
    }
}

// To avoid generating error on Linux
func RunEmbeddedBinary(binary string, arguments []string){

}

