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
	viper "github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// Arguments

// Web : bool to know if it's the web subcommand
var Web bool

// Bind Port
var Bind int

// IP where the webserver will listen
var IP string

// TLS : bool if TLS is enabled or not
var TLS bool

// TLSCertificate is the tls certificate
var TLSCertificate string

// TLSKey is the tls key
var TLSKey string

// DisableListing : option to disable directory listing
var DisableListing bool

// Gzip enable gzip compression
var Gzip bool

// Oneliners : boolean to enable generation of oneliners to download and execute files
var Oneliners bool

// Username for the private folder
var Username string

// Password for the private folder
var Password string

// PrivateFolder is the path for the private folder
var PrivateFolder string

// RootFolder is the path for the root folder
var RootFolder string

// SearchAndReplaceMap is the map which contains the information to search and replace string by another
var SearchAndReplaceMap = make(map[string]string)

// Webdav to enable the webdav server
var Webdav bool

// Webdav to enable the webdav server
var WebdavPort int

var promptPassword bool
var cwd string
var searchAndReplace string

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the webserver (default subcommand)",
	Long:  `Start the webserver (default subcommand)`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Map the variable in the config file
		Bind = viper.GetInt("bind")
		IP = viper.GetString("IP")
		TLS = viper.GetBool("TLS")
		TLSCertificate = viper.GetString("Certificate")
		TLSKey = viper.GetString("Key")
		DisableListing = viper.GetBool("DisableListing")
		Gzip = viper.GetBool("Gzip")
		Oneliners = viper.GetBool("Oneliners")
		Username = viper.GetString("Username")
		Password = viper.GetString("Password")
		RootFolder = viper.GetString("Root")
		PrivateFolder = viper.GetString("Private")
		Webdav = viper.GetBool("webdav")
		WebdavPort = viper.GetInt("webdavPort")

		promptPassword = viper.GetBool("promptPassword")
		searchAndReplace = viper.GetString("searchAndReplace")
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
			utils.Check(err, "Error when reading password")

			//text, _ := reader.ReadString('\n')
			Password = strings.TrimSpace(string(bytePassword))
		}

		// For RootFolder
		// Remove if the last character is /
		if strings.HasSuffix(RootFolder, "/") {
			RootFolder = strings.TrimSuffix(RootFolder, "/")
		}
		// If relative path
		if !strings.HasPrefix(RootFolder, "/") {
			RootFolder = path.Join((cwd), RootFolder)
		}

		// For PrivateFolder
		// Remove if the last character is /
		if strings.HasSuffix(PrivateFolder, "/") {
			PrivateFolder = strings.TrimSuffix(PrivateFolder, "/")
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
	webCmd.Flags().StringVarP(&RootFolder, "root", "r", cwd, "Root folder (for web and webdav)")
	webCmd.Flags().StringVarP(&searchAndReplace, "searchAndReplace", "s", "", "Search and replace string in embedded text files")

	webCmd.Flags().StringVarP(&Username, "username", "u", "admin", "Username for basic auth")
	webCmd.Flags().StringVarP(&Password, "password", "p", "notsecure", "Password for basic auth")
	webCmd.Flags().StringVar(&PrivateFolder, "private", cwd+"/private", "Private folder with basic auth")
	webCmd.Flags().BoolVar(&promptPassword, "promptPassword", false, "Prompt for for basic auth's password")

	webCmd.Flags().BoolVar(&TLS, "tls", false, "Enables HTTPS (for web and webdav)")
	webCmd.Flags().StringVarP(&TLSCertificate, "certificate", "c", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 (for web and webdav)")
	webCmd.Flags().StringVarP(&TLSKey, "key", "k", "", "HTTPS Key : openssl genrsa -out server.key 2048 (for web and webdav)")

	webCmd.Flags().BoolVarP(&DisableListing, "disableListing", "d", false, "Disable directory listing")

	webCmd.Flags().BoolVarP(&Webdav, "webdav", "w", false, "Enable webdav (easier for copy with windows and for capture Net-NTLM hashes")
	webCmd.Flags().IntVar(&WebdavPort, "webdavPort", 8081, "Port for webdav")

}
