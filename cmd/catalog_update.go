package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

	// newCommand wraps creating of commands to exec in order to allow
	// mocking for testing.
	newCommand func(name, command string, args ...string) runnable = execCommand
)

func init() {
	catalogCmd.AddCommand(updateCmd)
}

type runnable interface {
	Run() error
}

func execCommand(_, command string, args ...string) runnable {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func cloneRepo(name, url, dir, revision string) error {
	args := []string{"clone"}
	if revision != "" {
		args = append(args, "-b", revision)
	}
	args = append(args, url, dir)

	cmd := newCommand(name, "git", args...)
	return cmd.Run()
}

func updateRepo(name, dir, revision string) error {
	gitdir := filepath.Join(dir, ".git")
	fetch := newCommand(name, "git",
		"--git-dir", gitdir,
		"--work-tree", dir,
		"fetch", "origin", revision)
	if err := fetch.Run(); err != nil {
		return err
	}

	checkout := newCommand(name, "git",
		"--git-dir", gitdir,
		"--work-tree", dir,
		"-c", "advice.detachedHead=false",
		"checkout", "FETCH_HEAD")
	return checkout.Run()
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

		dir := cfg.Path()
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := cloneRepo(cfg.Name, cfg.URL, dir, cfg.Revision); err != nil {
				return err
			}
		} else {
			if err := updateRepo(cfg.Name, dir, cfg.Revision); err != nil {
				return err
			}
		}
	}

	return nil
}

func update(cmd *cobra.Command, args []string) error {
	return updateConfig(args...)
}
