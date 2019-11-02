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
	testcases := []struct {
		name    string
		want    string
		method  string
		headers http.Header
		body    string
		query   string
	}{
		{
			name: "empty body and params",
			want: `method: GET
headers:
  Accept-Encoding:
  - gzip
  User-Agent:
  - Go-http-client/1.1
body: ""
params: {}
`,
			method: http.MethodGet,
			body:   "",
			query:  "",
		},
		{
			name: "body and params with headers",
			want: `method: POST
headers:
  Accept-Encoding:
  - gzip
  Captain:
  - Hook
  Content-Length:
  - "5"
  User-Agent:
  - Go-http-client/1.1
body: tacos
params:
  key:
  - value
  other:
  - one
  - two
`,
			method:  http.MethodPost,
			headers: http.Header{"Captain": {"Hook"}},
			body:    "tacos",
			query:   "?key=value&other=one&other=two",
		},
	}

	for _, tc := range testcases {
		f := testfile(t, "hook.yml")

		r, err := newRecorder(f.Name())
		if err != nil {
			t.Fatal(err)
		}

		srv := httptest.NewServer(r)

		req, err := http.NewRequest(tc.method, srv.URL+tc.query, strings.NewReader(tc.body))
		if err != nil {
			t.Fatal(err)
		}

		req.Header = tc.headers

		client := http.DefaultClient

		if _, err := client.Do(req); err != nil {
			t.Fatal(err)
		}
		got := readfile(t, f)
		if d := diff.Diff(tc.want, got); d != "" {
			t.Error(d)
		}

		// Make request again to test appends.
		if _, err := client.Do(req); err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("%s---\n%s", tc.want, tc.want)
		got = readfile(t, f)
		if d := diff.Diff(want, got); d != "" {
			t.Error(d)
		}

		deletefile(t, f)
		r.close()
		srv.Close()
	}
}
