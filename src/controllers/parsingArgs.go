package controllers

import "flag"
import "os"
import "fmt"
import "strings"
import "SimpleHTTPServer-golang/src/utils"


var Bind    *string
var Root_folder *string
var Username *string
var Password *string
var Uses_gzip *bool
var Public *string
var Private *string

func ParseArgs(){
        cwd, err := os.Getwd()
        if err != nil {
                fmt.Printf("Error while getting current directory.")
                return
        }

        // Command line parsing
        Bind = flag.String("bind", "8080", "Bind Port")
        Root_folder = flag.String("root", cwd, "Root folder")
        Public = flag.String("public", "", "Default " + cwd + " public folder")
        Private = flag.String("private", "private", "Private folder with basic auth, default " + cwd + "/private")
        Username = flag.String("username", "admin", "Username for basic auth, default: admin")
        Password = flag.String("password", "notsecure", "Password for basic auth, default: notsecure")
        Uses_gzip = flag.Bool("gzip", true, "Enables gzip/zlib compression")

        flag.Parse()
        if *Private != "" {
            // Remove if the last character is /
            if strings.HasSuffix(*Private,"/"){
                *Private = utils.TrimSuffix(*Private, "/")
            }
        }
}
