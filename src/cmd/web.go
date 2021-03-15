package cmd

import (
	"Swego/src/utils"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var Web bool

var Bind int
var IP string
var TLS bool
var TLSCertificate string
var TLSKey string
var DisableListing bool
var Gzip bool
var Oneliners bool
var Username string
var Password string
var PrivateFolder string
var RootFolder string
var SearchAndReplaceMap = make(map[string]string)

var promptPassword bool
var cwd string
var searchAndReplace string

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the webserver (default subcommand)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (TLS || TLSKey != "" || TLSCertificate != "") && (!TLS || TLSKey == "" || TLSCertificate == "") {
			return errors.New("Tls, certificate and/or key arguments missing")

		} else if TLS && (!utils.FileExists(TLSCertificate) || !utils.FileExists(TLSKey)) { //if TLS enable check if the certificate and key files not exist
			return errors.New("Certificate file " + TLSCertificate + " or key file " + TLSKey + " not found")

		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if promptPassword {
			//reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Print("\n")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//text, _ := reader.ReadString('\n')
			Password = strings.TrimSpace(string(bytePassword))
		}

		// For RootFolder
		// Remove if the last character is /
		if strings.HasSuffix(RootFolder, "/") {
			RootFolder = utils.TrimSuffix(RootFolder, "/")
		}
		// If relative path
		if !strings.HasPrefix(RootFolder, "/") {
			RootFolder = path.Join((cwd), RootFolder)
		}

		// For PrivateFolder
		// Remove if the last character is /
		if strings.HasSuffix(PrivateFolder, "/") {
			PrivateFolder = utils.TrimSuffix(PrivateFolder, "/")
		}
		// If relative path
		if !strings.HasPrefix(PrivateFolder, "/") {
			PrivateFolder = path.Join((RootFolder), PrivateFolder)
		}

		if searchAndReplace != "" {
			for _, item := range strings.Split(searchAndReplace, " ") {
				key := strings.Split(item, "=")[0]
				value := strings.Split(item, "=")[1]
				SearchAndReplaceMap[key] = value
			}
		}

		Web = true

	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return
	}

	webCmd.Flags().IntVarP(&Bind, "bind", "b", 8080, "Bind Port")
	webCmd.Flags().StringVar(&IP, "ip", "0.0.0.0", "Binding IP")
	webCmd.Flags().BoolVarP(&Gzip, "gzip", "g", true, "Enables gzip/zlib compression")
	webCmd.Flags().BoolVarP(&Oneliners, "oneliners", "o", false, "Generate oneliners to download files")
	webCmd.Flags().StringVarP(&RootFolder, "root", "r", cwd, "Root folder")
	webCmd.Flags().StringVarP(&searchAndReplace, "searchAndReplace", "s", "", "Search and replace string in embedded text files")

	webCmd.Flags().StringVarP(&Username, "username", "u", "admin", "Username for basic auth")
	webCmd.Flags().StringVarP(&Password, "password", "p", "notsecure", "Password for basic auth")
	webCmd.Flags().StringVar(&PrivateFolder, "private", cwd+"/private", "Private folder with basic auth")
	webCmd.Flags().BoolVar(&promptPassword, "promptPassword", false, "Prompt for for basic auth's password")

	webCmd.Flags().BoolVar(&TLS, "tls", false, "Enables HTTPS")
	webCmd.Flags().StringVarP(&TLSCertificate, "certificate", "c", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365")
	webCmd.Flags().StringVarP(&TLSKey, "key", "k", "", "HTTPS Key : openssl genrsa -out server.key 2048")

	webCmd.Flags().BoolVarP(&DisableListing, "disableListing", "d", false, "Disable directory listing")

}
