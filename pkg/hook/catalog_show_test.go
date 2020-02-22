package hook

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestShowHook(t *testing.T) {
	b := new(bytes.Buffer)
	if err := ShowHook(b, filepath.Join("testdata", "a")); err != nil {
		t.Fatalf("ShowHook: %v", err)
	}

	want := `---
method: POST
headers:
  foo:
  - bar
  - baz
  herp:
  - derp
body: test=body
`
	if diff := cmp.Diff(want, b.String()); diff != "" {
		t.Error(diff)
	}
}
