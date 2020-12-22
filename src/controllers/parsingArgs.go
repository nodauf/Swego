package controllers

import (
	"SimpleHTTPServer-golang/src/utils"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

// Subcommand
var WebCommand *flag.FlagSet
var RunCommand *flag.FlagSet

// Web subcommand
var Bind *string
var Root_folder *string
var Username *string
var Password *string
var Uses_gzip *bool
var Private *string
var Tls *bool
var Key *string
var Certificate *string
var SearchAndReplace = make(map[string]string)

// Run subcommand
var List *bool
var Binary *string
var Args *string

func ParseArgs() {

	//help := flag.Bool("help",false,"Print usage")
	//flag.Parse()
	// Subcommands
	WebCommand = flag.NewFlagSet("web", flag.ExitOnError)
	RunCommand = flag.NewFlagSet("run", flag.ExitOnError)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return
	}

	// Command line parsing for subcommand web
	Bind = WebCommand.String("bind", "8080", "Bind Port")
	Root_folder = WebCommand.String("root", cwd, "Root folder")
	Private = WebCommand.String("private", "private", "Private folder with basic auth, default "+cwd+"/private")
	Username = WebCommand.String("username", "admin", "Username for basic auth, default: admin")
	Password = WebCommand.String("password", "notsecure", "Password for basic auth, default: notsecure")
	Uses_gzip = WebCommand.Bool("gzip", true, "Enables gzip/zlib compression")
	Tls = WebCommand.Bool("tls", false, "Enables HTTPS")
	Key = WebCommand.String("key", "", "HTTPS Key : openssl genrsa -out server.key 2048")
	Certificate = WebCommand.String("certificate", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365")
	searchAndReplace := WebCommand.String("s", "", "Search and replace string in embedded text files")
	helpWeb := WebCommand.Bool("help", false, "Print usage")

	// Command line parsing for subcommand run
	List = RunCommand.Bool("list", false, "List the embedded files")
	Binary = RunCommand.String("binary", "", "Binary to execute")
	Args = RunCommand.String("args", "", "Arguments for the binary")
	helpRun := RunCommand.Bool("help", false, "Print usage")

	// If this is not run subcommand
	if os.Args[1] != "run" {
		WebCommand.Parse(os.Args[1:])
		//return
	}

	// If the second argument is a subcommand
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		switch os.Args[1] {
		case "web":
			WebCommand.Parse(os.Args[2:])
		case "run":
			if runtime.GOOS == "windows" {
				RunCommand.Parse(os.Args[2:])
			} else {
				fmt.Println("run subcommand only available on Windows not on " + runtime.GOOS)
				showUsage()
				os.Exit(1)
			}
		default:
			showUsage()
			os.Exit(1)
		}
	}
	if *helpWeb {
		showUsage()
		os.Exit(1)
	}
	if WebCommand.Parsed() {
		if *Private != "" {
			// Remove if the last character is /
			if strings.HasSuffix(*Private, "/") {
				*Private = utils.TrimSuffix(*Private, "/")
			}
			// If relative path
			if !strings.HasPrefix(*Private, "/") {
				*Private = path.Join((*Root_folder), *Private)
			}
		}
		if (*Tls || *Key != "" || *Certificate != "") && (!*Tls || *Key == "" || *Certificate == "") {
			fmt.Print("Tls, certificate and/or key arguments missing\n")
			WebCommand.PrintDefaults()
			os.Exit(1)
		} else if *Tls && (!utils.FileExists(*Certificate) || !utils.FileExists(*Key)) { //if TLS enable check if the certificate and key files not exist
			fmt.Printf("Certificate file %s or key file %s not found\n", *Certificate, *Key)
			WebCommand.PrintDefaults()
			os.Exit(1)
		}

		if *searchAndReplace != "" {
			for _, item := range strings.Split(*searchAndReplace, " ") {
				key := strings.Split(item, "=")[0]
				value := strings.Split(item, "=")[1]
				SearchAndReplace[key] = value
			}
		}

	} else if RunCommand.Parsed() {
		if !*List && *Binary == "" {
			fmt.Println("You must specify a binary to run")
			RunCommand.PrintDefaults()
			fmt.Println("\nPackaged Binaries:")
			PrintEmbeddedFiles()
			os.Exit(1)
		}
		//        // If list and binary
		if *List && *Binary != "" {
			fmt.Println("You must specify either binary or list")
			RunCommand.PrintDefaults()
			fmt.Println("\nPackaged Binaries:")
			PrintEmbeddedFiles()
			os.Exit(1)
		}
		if *helpRun {
			RunCommand.PrintDefaults()
			fmt.Println("\nPackaged Binaries:")
			PrintEmbeddedFiles()
			os.Exit(1)
		}
	}
}

func showUsageRun() {
	fmt.Printf("Usage:\n%s run <binary> \"<args>\"\n", os.Args[0])
	fmt.Println("\nPackaged Binaries:")
	PrintEmbeddedFiles()
}

func showUsage() {
	fmt.Println("web subcommand")
	WebCommand.PrintDefaults()
	fmt.Println("\nrun subcommand")
	RunCommand.PrintDefaults()
	fmt.Println("\nPackaged Binaries:")
	PrintEmbeddedFiles()
}
