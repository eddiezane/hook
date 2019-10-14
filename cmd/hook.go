package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hook",
	Short: "Hook is a tool for firing a known collection of webhooks",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra already outputs errors returned from RunE
		os.Exit(1)
	}
}
