package cmd

import (
	"os"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hook",
	Short: "Hook is a tool for firing a known collection of webhooks",
}

// Execute runs the root command
func Execute() {
	hook.Initcfg()
	if err := rootCmd.Execute(); err != nil {
		// Cobra already outputs errors returned from RunE
		os.Exit(1)
	}
}
