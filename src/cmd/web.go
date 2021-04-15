package cmd

import (
	"Swego/src/utils"
	"crypto/tls"
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
var tlsCertificate string

// TLSKey is the tls key
var tlsKey string

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

// TLSConfig contains the configuration for the webserver
var TLSConfig tls.Config

// cert contains certificate and private key
var cert *tls.Certificate

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
		tlsCertificate = viper.GetString("Certificate")
		tlsKey = viper.GetString("Key")
		DisableListing = viper.GetBool("DisableListing")
		Gzip = viper.GetBool("Gzip")
		Oneliners = viper.GetBool("Oneliners")
		Username = viper.GetString("Username")
		Password = viper.GetString("Password")
		RootFolder = viper.GetString("Root")
		PrivateFolder = viper.GetString("Private")

		promptPassword = viper.GetBool("promptPassword")
		searchAndReplace = viper.GetString("searchAndReplace")

		if TLS && (tlsKey == "" || tlsCertificate == "") {
			var err error
			cert, err = utils.GenerateTLSCertificate("")
			if err != nil {
				return errors.New("Error while generating certificate: " + err.Error())
			}
			TLSConfig.Certificates = append(TLSConfig.Certificates, *cert)
		} else if (TLS || tlsKey != "" || tlsCertificate != "") && (!TLS || tlsKey == "" || tlsCertificate == "") {
			return errors.New("Tls, certificate and/or key arguments missing")

		} else if TLS && (!utils.FileExists(tlsCertificate) || !utils.FileExists(tlsKey)) { //if TLS enable check if the certificate and key files not exist
			return errors.New("Certificate file " + tlsCertificate + " or key file " + tlsKey + " not found")

		} else if TLS && utils.FileExists(tlsCertificate) && utils.FileExists(tlsKey){
			cer, err := tls.LoadX509KeyPair(tlsCertificate, tlsKey)
			if err != nil {
				return errors.New(err.Error())
			}
			cert = &cer
			TLSConfig.Certificates = append(TLSConfig.Certificates, *cert)
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
	webCmd.Flags().StringVarP(&RootFolder, "root", "r", cwd, "Root folder")
	webCmd.Flags().StringVarP(&searchAndReplace, "searchAndReplace", "s", "", "Search and replace string in embedded text files")

	webCmd.Flags().StringVarP(&Username, "username", "u", "admin", "Username for basic auth")
	webCmd.Flags().StringVarP(&Password, "password", "p", "notsecure", "Password for basic auth")
	webCmd.Flags().StringVar(&PrivateFolder, "private", cwd+"/private", "Private folder with basic auth")
	webCmd.Flags().BoolVar(&promptPassword, "promptPassword", false, "Prompt for for basic auth's password")

	webCmd.Flags().BoolVar(&TLS, "tls", false, "Enables HTTPS (for web and webdav)")
	webCmd.Flags().StringVarP(&tlsCertificate, "certificate", "c", "", "HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 (for web and webdav)")
	webCmd.Flags().StringVarP(&tlsKey, "key", "k", "", "HTTPS Key : openssl genrsa -out server.key 2048 (for web and webdav)")

	webCmd.Flags().BoolVarP(&DisableListing, "disableListing", "d", false, "Disable directory listing")

}
