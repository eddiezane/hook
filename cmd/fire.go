package cmd

import (
	"fmt"
	"os"

	"github.com/eddiezane/captain-hook/pkg/hook"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fireCommand)

	// TODO(eddiezane): Slurp up all hooks here?
}

var fireCommand = &cobra.Command{
	Use:   "fire",
	Short: "Fires the selected webhook at a given url",
	RunE:  fire,
}

func fire(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		cmd.Usage()
		os.Exit(1)
	}
	path := args[0]

	h, err := hook.NewFromPath(path)
	if err != nil {
		return err
	}

	target := args[1]
	res, err := h.Fire(target)
	fmt.Println(res)
	return err
}
