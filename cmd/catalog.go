package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(catalogCmd)
}

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Subcommands for catalog operations",
	Long:  "CATALOG COMMANDS ARE WIP AND ARE NOT COMPLETE.",
}
