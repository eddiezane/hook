package cmd

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

const (
	appName = "hook"
)

// initcfg intializes the config environment. This is not done as a part of
// the standard Go init function so that we can override environment specific
// variables (such as $HOME) in tests.
func initcfg() {
	viper.SetEnvPrefix(appName)

	// Set locations for where to look for config files.
	if cfg, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(cfg, ".config", appName))
	}
	if cfg, err := os.UserConfigDir(); err == nil {
		viper.AddConfigPath(filepath.Join(cfg, appName))
	}

	cache, err := os.UserCacheDir()
	if err != nil {
		// Fallback to system tmp directory if user dir cannot be found.
		cache = os.TempDir()
	}
	viper.SetDefault("cache", filepath.Join(cache, appName))

	viper.SetConfigName(appName)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("error reading config file: ", err)
			return
		}
	}
}

type config struct {
	Cache   string
	Catalog *catalogConfig
}

type catalogConfig struct {
	Remote []*remoteConfig
}

func (c catalogConfig) sort() {
	sort.Slice(c.Remote, func(i, j int) bool {
		return strings.Compare(c.Remote[i].Name, c.Remote[j].Name) < 0
	})

}

type remoteConfig struct {
	Name     string
	URL      string
	Revision string `yaml:"omitempty"`
}
