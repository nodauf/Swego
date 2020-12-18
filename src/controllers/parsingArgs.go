package controllers

import "flag"
import "os"
import "fmt"
import "strings"
import "SimpleHTTPServer-golang/src/utils"
import "runtime"
import "path"

// Subcommand
var WebCommand *flag.FlagSet
var RunCommand *flag.FlagSet

// Web subcommand
var Bind    *string
var Root_folder *string
var Username *string
var Password *string
var Uses_gzip *bool
var Private *string
var Tls *bool
var Key *string
var Certificate *string

// Run subcommand
//var List *bool
//var Binary *string
//var Args *string
var Run bool
var Binary string
var Args []string

func ParseArgs(){

        //help := flag.Bool("help",false,"Print usage")
        //flag.Parse()
        // Subcommands
        WebCommand = flag.NewFlagSet("web", flag.ExitOnError)
        //RunCommand = flag.NewFlagSet("run", flag.ExitOnError)

        cwd, err := os.Getwd()
        if err != nil {
                fmt.Printf("Error while getting current directory.")
                return
        }

        // Command line parsing for subcommand web
        Bind = WebCommand.String("bind", "8080", "Bind Port")
        Root_folder = WebCommand.String("root", cwd, "Root folder")
        Private = WebCommand.String("private", "private", "Private folder with basic auth, default " + cwd + "/private")
        Username = WebCommand.String("username", "admin", "Username for basic auth, default: admin")
        Password = WebCommand.String("password", "notsecure", "Password for basic auth, default: notsecure")
        Uses_gzip = WebCommand.Bool("gzip", true, "Enables gzip/zlib compression")
        Tls = WebCommand.Bool("tls", false, "Enables HTTPS")
        Key = WebCommand.String("key", "", "HTTPS Key : openssl genrsa -out server.key 2048")
        Certificate = WebCommand.String("certificate", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365")
        helpWeb := WebCommand.Bool("help",false,"Print usage")

        // Command line parsing for subcommand run
        //List = RunCommand.Bool("list",false,"List the embedded files")
        //Binary = RunCommand.String("binary","","Binary to execute")
        //Args = RunCommand.String("args","","Arguments for the binary")
        //helpRun := RunCommand.Bool("help",false,"Print usage")

        // If nothing is specified run the web server
        if len(os.Args) == 1 {
			WebCommand.Parse(os.Args[1:])
			return
		}

        //if len(os.Args) < 2 {
        //    fmt.Println("web or run subcommand is required, default is web")
        //    //WebCommand.Parse(os.Args[1:])
        //    os.Exit(1)
        //}

		switch os.Args[1] {
			case "web":
				WebCommand.Parse(os.Args[2:])
			case "run":
                if runtime.GOOS == "windows"{
					// binary missing
					if len(os.Args) == 2 {
						showUsageRun()
						os.Exit(1)
					}
					Run = true
  //                  RunCommand.Parse(os.Args[2:])
                }else{
                    fmt.Println("run subcommand only available on Windows not on: "+ runtime.GOOS)
                }
			default:
                fmt.Println("web subcommand")
                WebCommand.PrintDefaults()
                fmt.Println("\nrun subcommand")
                showUsageRun()
                //RunCommand.PrintDefaults()
				os.Exit(1)
		}
        // Show usage if help in any subcommand
        if *helpWeb {
            fmt.Println("web subcommand")
			WebCommand.PrintDefaults()
            fmt.Println("\nrun subcommand")
            showUsageRun()
			//RunCommand.PrintDefaults()
			os.Exit(1)
        }
		if WebCommand.Parsed(){
			if *Private != "" {
				// Remove if the last character is /
				if strings.HasSuffix(*Private,"/"){
					*Private = utils.TrimSuffix(*Private, "/")
                }
                // If relative path
                if !strings.HasPrefix(*Private, "/"){
                    *Private = path.Join((*Root_folder), *Private)
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
		}else if Run { //RunCommand.Parsed(){
			Binary = os.Args[2]
			Args = os.Args[3:]
            fmt.Println(Args)
//        // If not listing and no binary select
//        if !*List && *Binary == "" {
//            fmt.Println("You must specify a binary to run")
//            RunCommand.PrintDefaults()
//            os.Exit(1)
//        }
//        // If list and binary
//        if *List && *Binary != "" {
//            fmt.Println("You must specify either binary or list")
//            RunCommand.PrintDefaults()
//            os.Exit(1)
//        }
	 }
}

func showUsageRun(){
    fmt.Printf("Usage:\n%s run <binary> <args>\n", os.Args[0])
	fmt.Println("\nPackaged Binaries:")
	PrintEmbeddedFiles()
}
