package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	viper "github.com/spf13/viper"
)

var cfgFile string
var help bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Swego",
	Short: "Swiss army knife Webserver in Golang",
	Long: `Alternative to SimpleHTTPServer of python. This is a cross-platform webserver with basic functionnalities likes upload, download folder as zip
	https, create folder from the web, private folder, ...
	Default will run subcommand web`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return webCmd.PreRunE(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		webCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Swego.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "Help message")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".Swego" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".Swego")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}

	// Bind the webcmd flags to the config file
	viper.BindPFlags(webCmd.Flags())
	viper.BindPFlags(runCmd.Flags())
}
