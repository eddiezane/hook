package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func TestDefaultName(t *testing.T) {
	testcases := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "",
			url:  "https://github.com/eddiezane/hook-catalog",
			want: "hook-catalog",
		},
		{
			name: ".git",
			url:  "https://github.com/eddiezane/hook-catalog.git",
			want: "hook-catalog",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := defaultName(tc.url)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func readpath(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	return readfile(t, f)
}

func cachedir(t *testing.T) string {
	t.Helper()
	d, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(d, "hook")
}

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
	hook.Initcfg()

	return d
}

func TestAddConfig(t *testing.T) {
	d := testdirInit(t)
	defer os.RemoveAll(d)

	if err := addConfig("https://example.com/foo", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := addConfig("https://example.com/bar.git", "tacocat", "master"); err != nil {
		t.Fatal(err)
	}

	want := &hook.Config{
		Cache: cachedir(t),
		Catalog: &hook.CatalogConfig{
			Remote: []*hook.RemoteConfig{
				hook.DefaultCatalog,
				{
					Name: "foo",
					URL:  "https://example.com/foo",
				},
				{
					Name:     "tacocat",
					URL:      "https://example.com/bar.git",
					Revision: "master",
				},
			},
		},
	}

	s := readpath(t, filepath.Join(d, ".config", "hook", "hook.yaml"))
	got := new(hook.Config)
	if err := yaml.Unmarshal([]byte(s), got); err != nil {
		t.Fatal(err)
	}
	got.Catalog.Sort()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
