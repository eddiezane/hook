package hook

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// testdirInit creates a test directory, sets relevant environment variables,
// and initializes hook configuration to read from the directory.
func testdirInit(t *testing.T) string {
	t.Helper()
	d, err := ioutil.TempDir("", "hook")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("HOME", d)

	viper.Reset()
	Initcfg()

	return d
}

// testfile is a shortcut to create a new test file.
func testfile(t *testing.T, path string) *os.File {
	t.Helper()

	f, err := ioutil.TempFile("", path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// deletefile is a shortcut to delete a file.
func deletefile(t *testing.T, f *os.File) {
	t.Helper()
	if err := os.Remove(f.Name()); err != nil {
		t.Fatal(err)
	}
}

// readfile is a shortcut to read the contents of a file.
func readfile(t *testing.T, f *os.File) string {
	t.Helper()
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// readpath is a shortcut to read a given path.
func readpath(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	return readfile(t, f)
}

// cachedir is a shortcut to get the cache directory.
func cachedir(t *testing.T) string {
	t.Helper()
	d, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(d, "hook")
}
