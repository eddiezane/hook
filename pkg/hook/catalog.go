package hook

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ParsePath takes a hook name of the form <catalog>@<path> and returns the
// individual pieces.
func ParsePath(uri string) (catalog string, path string) {
	s := strings.SplitN(uri, "@", 2)
	if len(s) == 1 {
		return "", s[0]
	}

	if s[0] == "" && strings.Contains(uri, "@") {
		return "@", s[1]
	}

	return s[0], s[1]
}

// Catalog represents a mechanism for fetching hook configurations.
type Catalog interface {
	Open(path string) (*os.File, error)
}

// RemoteConfig describes a single catalog remote.
type RemoteConfig struct {
	Name     string
	URL      string
	Revision string `yaml:",omitempty"`
}

// Path returns the cache path where the remote cache exists.
func (rc RemoteConfig) Path() string {
	return filepath.Join(viper.GetString("cache"), rc.Name)
}

func (rc *RemoteConfig) isCached() bool {
	_, err := os.Stat(rc.Path())
	return !os.IsNotExist(err)
}

// Open returns the file corresponding to the file in the remote config,
// cloning the config locally if it has not occured yet.
func (rc *RemoteConfig) Open(path string) (*os.File, error) {
	if !rc.isCached() {
		if err := rc.Clone(); err != nil {
			return nil, err
		}
	}

	return openFile(filepath.Join(rc.Path(), path))
}

var newCommand func(name, command string, args ...string) runnable = execCommand

type runnable interface {
	Run() error
}

func execCommand(_, command string, args ...string) runnable {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// Clone clones the remote catalog to the hook cache.
func (rc *RemoteConfig) Clone() error {
	if rc.Name == "" {
		return errors.New("RemoteConfig must have a name")
	}

	args := []string{"clone"}
	if rc.Revision != "" {
		args = append(args, "-b", rc.Revision)
	}
	args = append(args, rc.URL, rc.Path())

	cmd := newCommand(rc.Name, "git", args...)
	return cmd.Run()
}

// Update pulls the latest version of the catalog to the hook cache.
// If the catalog does not exist in the cache, it is cloned.
func (rc *RemoteConfig) Update() error {
	if rc.Name == "" {
		return errors.New("RemoteConfig must have a name")
	}

	if !rc.isCached() {
		return rc.Clone()
	}

	gitdir := filepath.Join(rc.Path(), ".git")
	fetch := newCommand(rc.Name, "git",
		"--git-dir", gitdir,
		"--work-tree", rc.Path(),
		"fetch", "origin", rc.Revision)
	if err := fetch.Run(); err != nil {
		return err
	}

	checkout := newCommand(rc.Name, "git",
		"--git-dir", gitdir,
		"--work-tree", rc.Path(),
		"-c", "advice.detachedHead=false",
		"checkout", "FETCH_HEAD")
	return checkout.Run()
}

// LocalCatalog handles reading configuration locally.
type LocalCatalog struct{}

// Open returns the local file. This is a wrapper around os.Open.
func (LocalCatalog) Open(path string) (*os.File, error) {
	return openFile(path)
}

// openFile opens the given file, allowing for fuzzing of the extension.
func openFile(path string) (*os.File, error) {
	if filepath.Ext(path) != "" {
		return os.Open(path)
	}

	var err error
	for _, ext := range []string{"", ".yaml", ".yml"} {
		var f *os.File
		f, err = os.Open(path + ext)
		if err == nil {
			return f, nil
		}
	}
	return nil, err
}
