package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var Run bool

var Args string
var Binary string
var List bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
}
