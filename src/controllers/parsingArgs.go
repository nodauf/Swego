package controllers

import "flag"
import "os"
import "fmt"
import "strings"
import "SimpleHTTPServer-golang/src/utils"


// Subcommand
var WebCommand *flag.FlagSet
var RunCommand *flag.FlagSet

// Web subcommand
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

// Run subcommand
var List *bool
var Binary *string
var Args *string

func ParseArgs(){

        help := flag.Bool("help",false,"Print usage")
        flag.Parse()
        // Subcommands
        WebCommand = flag.NewFlagSet("web", flag.ContinueOnError)
        RunCommand = flag.NewFlagSet("run", flag.ContinueOnError)

        cwd, err := os.Getwd()
        if err != nil {
                fmt.Printf("Error while getting current directory.")
                return
        }

        // Command line parsing for subcommand web
        Bind = WebCommand.String("bind", "8080", "Bind Port")
        Root_folder = WebCommand.String("root", cwd, "Root folder")
        Public = WebCommand.String("public", "", "Default " + cwd + " public folder")
        Private = WebCommand.String("private", "private", "Private folder with basic auth, default " + cwd + "/private")
        Username = WebCommand.String("username", "admin", "Username for basic auth, default: admin")
        Password = WebCommand.String("password", "notsecure", "Password for basic auth, default: notsecure")
        Uses_gzip = WebCommand.Bool("gzip", true, "Enables gzip/zlib compression")
        Tls = WebCommand.Bool("tls", false, "Enables HTTPS")
        Key = WebCommand.String("key", "", "HTTPS Key : openssl genrsa -out server.key 2048")
        Certificate = WebCommand.String("certificate", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365")

        // Command line parsing for subcommand run
        List = RunCommand.Bool("list",false,"List the embedded files")
        Binary = RunCommand.String("binary","","Binary to execute")
        Args = RunCommand.String("args","","Arguments for the binary")

        if *help {
            fmt.Println("web subcommand")
			WebCommand.PrintDefaults()
            fmt.Println("\nrun subcommand")
			RunCommand.PrintDefaults()
			os.Exit(1)
        }

        if len(os.Args) < 2 {
            fmt.Println("web or run subcommand is required")
            os.Exit(1)
        }

		switch os.Args[1] {
			case "web":
				WebCommand.Parse(os.Args[2:])
			case "run":
				RunCommand.Parse(os.Args[2:])
			default:
                fmt.Println("web subcommand")
                WebCommand.PrintDefaults()
                fmt.Println("\nrun subcommand")
                RunCommand.PrintDefaults()
				os.Exit(1)
		}
		if WebCommand.Parsed(){
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
				WebCommand.PrintDefaults()
				os.Exit(1)
			}else if *Tls && (!utils.FileExists(*Certificate) || !utils.FileExists(*Key)) { //if TLS enable check if the certificate and key files not exist
				fmt.Printf("Certificate file %s or key file %s not found\n", *Certificate, *Key)
				WebCommand.PrintDefaults()
				os.Exit(1)
			}
	}else if RunCommand.Parsed(){
        if !*List && *Binary == "" {
            fmt.Println("You must specify a binary to run")
            RunCommand.PrintDefaults()
            os.Exit(1)
        }
    }
}
