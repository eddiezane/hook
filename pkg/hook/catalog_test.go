package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		in, catalog, path string
	}{
		{
			in:      "@foo",
			catalog: "@",
			path:    "foo",
		},
		{
			in:      "foo@bar",
			catalog: "foo",
			path:    "bar",
		},
		{
			in:      "foo",
			catalog: "",
			path:    "foo",
		},
		{
			in:      "foo@",
			catalog: "foo",
			path:    "",
		},
	}
	for _, tc := range tests {
		if catalog, path := ParsePath(tc.in); catalog != tc.catalog || path != tc.path {
			t.Errorf("ParsePath(%s) = (%s, %s), want (%s, %s)", tc.in, catalog, path, tc.catalog, tc.path)
		}
	}
}

// commandSink captures and stores commands created during execution.
// It guarantees that the returned commands will be consistently ordered by repo
// and by command order.
// It does not guarantee that the repo ordering will be identical to the config
// ordering.
type commandSink struct {
	cmd []*mockCommand
}

func (c *commandSink) reset() {
	c.cmd = []*mockCommand{}
}

// sort ensure consistent return order for comparisons by sorting by
// (catalog name, order event occured).
func (c *commandSink) sort() {
	sort.Slice(c.cmd, func(i, j int) bool {
		iName := c.cmd[i].name
		jName := c.cmd[j].name
		// Sort by catalog name first, then by command ordering.
		if iName != jName {
			return iName < jName
		}
		return c.cmd[i].idx < c.cmd[j].idx
	})
}

func (c *commandSink) commands() [][]string {
	c.sort()
	out := make([][]string, 0, len(c.cmd))
	for _, v := range c.cmd {
		out = append(out, v.args)
	}
	return out
}

func (c *commandSink) record(name, command string, args ...string) runnable {
	cmd := &mockCommand{
		name: name,
		idx:  len(args),
		args: append([]string{command}, args...),
	}
	c.cmd = append(c.cmd, cmd)
	return cmd
}

// mockCommand simulates a exec-ed command without invoking an external shell
// or the network. If mockCommand recognizes the request as a clone, it will
// create the specified directory to simulate a clone of the repo.
type mockCommand struct {
	// Catalog name command is acting on behalf of.
	// Used to guarantee consistent command return order.
	name string
	// Overall command index number. This is not necessarily tied to the specific
	// catalog name, but gives overall ordering of the commands.
	idx int
	// Command arguments that were called.
	args []string
}

func (c mockCommand) Run() error {
	if c.args[1] == "clone" {
		dir := c.args[3]
		return os.MkdirAll(dir, 0700)
	}
	return nil
}

func (c mockCommand) String() string {
	return strings.Join(c.args, " ")
}

func TestCatalogUpdate(t *testing.T) {
	d := testdirInit(t)
	defer os.RemoveAll(d)

	c := &commandSink{}
	newCommand = c.record

	tc := []*RemoteConfig{
		{
			URL:  "http://example.com/foo",
			Name: "foo",
		},
		&RemoteConfig{
			URL:      "https://example.com/bar",
			Name:     "tacocat",
			Revision: "v1",
		},
	}
	for _, rc := range tc {
		t.Run(rc.Name, func(t *testing.T) {
			c.reset()

			if err := rc.Update(); err != nil {
				t.Fatal(err)
			}
			workdir := filepath.Join(cachedir(t), rc.Name)
			cloneCmd := []string{"git", "clone", rc.URL, workdir}
			if rc.Revision != "" {
				cloneCmd = []string{"git", "clone", "-b", rc.Revision, rc.URL, workdir}
			}
			if diff := cmp.Diff([][]string{cloneCmd}, c.commands()); diff != "" {
				t.Error(diff)
			}

			// Simulate directory actually being cloned, run update again.
			os.Mkdir(workdir, os.ModePerm)
			if err := rc.Update(); err != nil {
				t.Fatal(err)
			}

			gitdir := filepath.Join(workdir, ".git")
			want := [][]string{
				cloneCmd,
				{
					"git",
					"--git-dir", gitdir,
					"--work-tree", workdir,
					"fetch", "origin", rc.Revision,
				},
				{
					"git",
					"--git-dir", gitdir,
					"--work-tree", workdir,
					"-c", "advice.detachedHead=false",
					"checkout", "FETCH_HEAD",
				},
			}
			if diff := cmp.Diff(want, c.commands()); diff != "" {
				t.Error(diff)
			}

			fmt.Println(c.commands())
		})
	}
}
