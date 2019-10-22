package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
			return strings.Compare(iName, jName) < 0
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

	if err := updateConfig(); err != nil {
		t.Fatal(err)
	}

	if err := addConfig("https://example.com/foo", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := updateConfig(); err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"git", "clone", "https://example.com/foo", filepath.Join(cachedir(t), "foo")}}
	if diff := cmp.Diff(want, c.commands()); diff != "" {
		t.Error(diff)
	}

	c.reset()

	if err := addConfig("https://example.com/bar", "tacocat", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := updateConfig(); err != nil {
		t.Fatal(err)
	}
	want = [][]string{
		{
			"git",
			"--git-dir", filepath.Join(cachedir(t), "foo", ".git"),
			"--work-tree", filepath.Join(cachedir(t), "foo"),
			"fetch", "origin", "",
		},
		{
			"git",
			"--git-dir", filepath.Join(cachedir(t), "foo", ".git"),
			"--work-tree", filepath.Join(cachedir(t), "foo"),
			"-c", "advice.detachedHead=false",
			"checkout", "FETCH_HEAD",
		},
		{"git", "clone", "-b", "v1", "https://example.com/bar", filepath.Join(cachedir(t), "tacocat")},
	}
	if diff := cmp.Diff(want, c.commands()); diff != "" {
		t.Error(diff)
	}
}
