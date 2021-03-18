package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

// Arguments

// Run : bool to know if it's the run subcommand
var Run bool

// Args for the binary to run
var Args string

// Binary to run
var Binary string

// List : Action to list the embedded binaries
var List bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an embedded binary",
	Long:  `Run an embedded binary`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		Args = viper.GetString("Args")
		Binary = viper.GetString("Binary")
		List = viper.GetBool("List")
		if !List && Binary == "" {
			return errors.New("You must specify a binary to run")

		} else if List && Binary != "" {
			return errors.New("You must specify either binary or list")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		Run = true
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&Args, "args", "a", "", "Arguments for the binary")
	runCmd.Flags().StringVarP(&Binary, "binary", "b", "", "Binary to execute")
	runCmd.Flags().BoolVarP(&List, "list", "l", false, "List embedded binaries")

	viper.BindPFlags(runCmd.Flags())

}
