package cmd

import (
	"os"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:     "show <url>...",
		Short:   "Show the catalog config(s).",
		Example: "hook catalog show @github/push",
		RunE:    show,
	}
)

func init() {
	catalogCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) error {
	return hook.ShowHook(os.Stdout, args...)
}
