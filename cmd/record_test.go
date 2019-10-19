package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/eddiezane/hook/pkg/hook"
	"github.com/kylelemons/godebug/diff"
)

func testfile(t *testing.T, path string) *os.File {
	t.Helper()

	f, err := ioutil.TempFile("", path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func deletefile(t *testing.T, f *os.File) {
	t.Helper()
	if err := os.Remove(f.Name()); err != nil {
		t.Fatal(err)
	}
}

func readfile(t *testing.T, f *os.File) string {
	t.Helper()
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func reqString(t *testing.T, req *http.Request) string {
	t.Helper()

	h, err := hook.NewFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	s, err := h.Dump()
	if err != nil {
		t.Fatal(err)
	}
	return string(s)
}

func TestRecord(t *testing.T) {
	f := testfile(t, "hook.yml")
	defer deletefile(t, f)

	r, err := newRecorder(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer r.close()

	srv := httptest.NewServer(r)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	s := `method: GET
headers:
  Accept-Encoding:
  - gzip
  User-Agent:
  - Go-http-client/1.1
body: ""
`

	client := http.DefaultClient

	if _, err := client.Do(req); err != nil {
		t.Fatal(err)
	}
	got := readfile(t, f)
	if d := diff.Diff(s, got); d != "" {
		t.Error(d)
	}

	// Make request again to test appends.
	if _, err := client.Do(req); err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("%s---\n%s", s, s)
	got = readfile(t, f)
	if d := diff.Diff(want, got); d != "" {
		t.Error(d)
	}
}
