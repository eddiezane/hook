package hook

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/viper"
)

const (
	// AppName is the name of the app.
	AppName = "hook"
)

var (
	// DefaultCatalog is the default catalog installed automatically for users.
	DefaultCatalog = &RemoteConfig{
		Name: "@",
		URL:  "https://github.com/eddiezane/hook-catalog",
	}
)

// Initcfg intializes the config environment. This is not done as a part of
// the standard Go init function so that we can override environment specific
// variables (such as $HOME) in tests.
func Initcfg() {
	viper.SetEnvPrefix(AppName)

	// Set locations for where to look for config files.
	if cfg, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(cfg, ".config", AppName))
	}
	if cfg, err := os.UserConfigDir(); err == nil {
		viper.AddConfigPath(filepath.Join(cfg, AppName))
	}

	cache, err := os.UserCacheDir()
	if err != nil {
		// Fallback to system tmp directory if user dir cannot be found.
		cache = os.TempDir()
	}
	viper.SetDefault("cache", filepath.Join(cache, AppName))

	viper.SetDefault("catalog.remote", DefaultCatalog)

	viper.SetConfigName(AppName)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("error reading config file: ", err)
			return
		}
	}
	log.Println("read config", viper.ConfigFileUsed())
}

// Config descirbes configuration for hook settings.
type Config struct {
	Cache   string
	Catalog *CatalogConfig
}

// CatalogConfig describes the hook remote catalog.
type CatalogConfig struct {
	Remote []*RemoteConfig
}

// Sort deterministically sorts the Catalog, useful for consistent ordering for
// tests.
func (c CatalogConfig) Sort() {
	sort.Slice(c.Remote, func(i, j int) bool {
		return c.Remote[i].Name < c.Remote[j].Name
	})
}

// RemoteConfigSet maps remote config names to the complete remote config.
type RemoteConfigSet map[string]*RemoteConfig

// GetRemoteConfigs gets the current remote catalog configuration.
// Returns a map so that values can easily be treated as a set.
func GetRemoteConfigs() (RemoteConfigSet, error) {
	var cfg []*RemoteConfig
	if err := viper.UnmarshalKey("catalog.remote", &cfg); err != nil {
		return nil, err
	}
	m := make(map[string]*RemoteConfig)
	for _, c := range cfg {
		m[c.Name] = c
	}
	return m, nil
}

// Get returns the RemoteConfig if it exists, or an error if it was missing.
func (r RemoteConfigSet) Get(name string) (*RemoteConfig, error) {
	cfg, ok := r[name]
	if !ok {
		return nil, fmt.Errorf("catalog %s not found", name)
	}
	return cfg, nil
}

// GetRemoteConfig resolves a catalog name into the underlying remote config.
// If the catalog doesn't exist, an error is returned.
func GetRemoteConfig(catalog string) (*RemoteConfig, error) {
	cfg, err := GetRemoteConfigs()
	if err != nil {
		return nil, err
	}
	return cfg.Get(catalog)
}
