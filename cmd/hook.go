package cmd

import (
	"os"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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

// GenMarkdownTree generates markdown documentation for the command.
func GenMarkdownTree(path string) error {
	return doc.GenMarkdownTree(rootCmd, path)
}
