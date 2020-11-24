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
var Tls *bool
var Key *string
var Certificate *string

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
        Tls = flag.Bool("tls", false, "Enables HTTPS")
        Key = flag.String("key", "", "HTTPS Key : openssl genrsa -out server.key 2048")
        Certificate = flag.String("certificate", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365")

        flag.Parse()
        if *Private != "" {
            // Remove if the last character is /
            if strings.HasSuffix(*Private,"/"){
                *Private = utils.TrimSuffix(*Private, "/")
            }
        }
        if *Public != "" {
            // Remove if the last character is /
            if strings.HasSuffix(*Public,"/"){
                *Public = utils.TrimSuffix(*Public, "/")
            }
        }
        if (*Tls || *Key != "" || *Certificate != "") && (!*Tls || *Key == "" || *Certificate == ""){
            fmt.Print("Tls, certificate and/or key arguments missing\n")
            flag.PrintDefaults()
            os.Exit(1)
        }else if *Tls && (!utils.FileExists(*Certificate) || !utils.FileExists(*Key)) { //if TLS enable check if the certificate and key files not exist
            fmt.Printf("Certificate file %s or key file %s not found\n", *Certificate, *Key)
            flag.PrintDefaults()
            os.Exit(1)
        }
}
