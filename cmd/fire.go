package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fireCommand)

	// TODO(eddiezane): Slurp up all hooks here?
}

var fireCommand = &cobra.Command{
	Use:   "fire",
	Short: "Fires the selected webhook at a given url",
	Run:   fire,
}

func fire(cmd *cobra.Command, args []string) {
	fmt.Println("fire command")
}
