package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "0.0.1"

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints out the version of hook",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Hook version %s", version))
	},
}
