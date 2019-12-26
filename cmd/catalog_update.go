package cmd

import (
	"log"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:     "update <url>",
		Short:   "Adds the given URL to the catalog config.",
		Long:    "update updates a given catalog to the configured revision.",
		Example: "hook catalog update <name>",
		RunE:    update,
	}
)

func init() {
	catalogCmd.AddCommand(updateCmd)
}

func updateConfig(names ...string) error {
	rc, err := hook.GetRemoteConfigs()
	if err != nil {
		return err
	}

	if len(names) == 0 {
		names = make([]string, 0, len(rc))
		for k := range rc {
			names = append(names, k)
		}
	}

	for _, n := range names {
		cfg, err := rc.Get(n)
		if err != nil {
			log.Print(err)
			continue
		}

		log.Printf("updating %s@%s", cfg.Name, cfg.Revision)

		if err := cfg.Update(); err != nil {
			return err
		}
	}

	return nil
}

func update(cmd *cobra.Command, args []string) error {
	return updateConfig(args...)
}
