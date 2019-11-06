package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tapCmd = &cobra.Command{
		Use:   "tap <url>",
		Short: "Adds the given URL to the catalog config.",
		// TODO(wlynch): Update this once the command is working,
		Long:    "TAP COMMAND IS WIP AND IS NOT FUNCTIONAL YET.",
		Example: "hook catalog tap https://github.com/eddiezane/hook-catalog",
		RunE:    tap,
	}

	tapName     string
	tapRevision string
)

func init() {
	tapCmd.Flags().StringVarP(&tapName, "name", "n", "", "Custom name of the remote. If not specified, a name will be inferred from the URL.")
	tapCmd.Flags().StringVarP(&tapName, "revision", "r", "", "Name of the revision to use when cloning the remote. If not specified, the default branch cloned will be used.")
	catalogCmd.AddCommand(tapCmd)
}

func defaultName(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(filepath.Base(u.Path), ".git"), nil
}

func writeRemoteConfig(m hook.RemoteConfigSet) error {
	cfg := make([]*hook.RemoteConfig, 0, len(m))
	for _, v := range m {
		cfg = append(cfg, v)
	}

	viper.Set("catalog.remote", cfg)

	err := viper.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		path, err := touchConfig()
		if err != nil {
			return err
		}
		log.Println("writing config file:", path)
		return viper.WriteConfigAs(path)
	}
	return err
}

// Create new default config file if one does not exist.
// Currently, viper does not create new config files as part of its write.
// See https://github.com/spf13/viper/issues/390 for more details.
func touchConfig() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".config", hook.AppName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.yaml", hook.AppName))
	f, err := os.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	return path, f.Close()
}

func addConfig(url, name, revision string) error {
	if name == "" {
		var err error
		name, err = defaultName(url)
		if err != nil {
			return err
		}
	}

	cfg, err := hook.GetRemoteConfigs()
	if err != nil {
		return err
	}
	cfg[name] = &hook.RemoteConfig{
		Name:     name,
		URL:      url,
		Revision: revision,
	}

	return writeRemoteConfig(cfg)
}

func tap(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("no arguments given")
	}
	url := args[0]

	return addConfig(url, tapName, tapRevision)
}
