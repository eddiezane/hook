package cmd

import (
	"fmt"
	"os"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fireCommand)
}

var fireCommand = &cobra.Command{
	Use:     "fire",
	Short:   "Fires the selected webhook at a given url",
	Long:    "Fire executes the selected webhook at the given url",
	Example: "hook fire webhooks/twilio/sms.yml http://localhost:3000",
	RunE:    fire,
}

func fire(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("incorrect number of arguments provided. expected %d", 2)
	}

	path := args[0]
	hooks, err := hook.NewFromPath(path)
	if err != nil {
		return err
	}

	target := args[1]
	for _, h := range hooks {
		res, err := h.Fire(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			continue
		}
		fmt.Println(res)
	}
	return nil
}
