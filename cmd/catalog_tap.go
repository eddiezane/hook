package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

// Get the current configuration.
// Returns a map so that values can easily be treated as a set.
func getRemoteConfig() (map[string]*remoteConfig, error) {
	var cfg []*remoteConfig
	if err := viper.UnmarshalKey("catalog.remote", &cfg); err != nil {
		return nil, err
	}
	m := make(map[string]*remoteConfig)
	for _, c := range cfg {
		m[c.Name] = c
	}
	return m, nil
}

func writeRemoteConfig(m map[string]*remoteConfig) error {
	cfg := make([]*remoteConfig, 0, len(m))
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

	dir := filepath.Join(home, ".config", appName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.yaml", appName))
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

	cfg, err := getRemoteConfig()
	if err != nil {
		return err
	}

	cfg[name] = &remoteConfig{
		Name:     name,
		URL:      url,
		Revision: revision,
	}

	return writeRemoteConfig(cfg)
}

func tap(cmd *cobra.Command, args []string) error {
	initcfg()

	if len(args) < 1 {
		return errors.New("no arguments given")
	}
	url := args[0]

	return addConfig(url, tapName, tapRevision)
}
