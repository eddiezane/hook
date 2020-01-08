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
		opts    []hook.Option
	}{
		{
			name: "empty body and params",
			want: `method: GET
headers:
  Accept-Encoding:
  - gzip
  User-Agent:
  - Go-http-client/1.1
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
		{
			name: "json body",
			want: `method: POST
headers:
  Accept-Encoding:
  - gzip
  Content-Length:
  - "14"
  Content-Type:
  - application/json
  User-Agent:
  - Go-http-client/1.1
body: |-
  {
    "foo": "bar"
  }
`,
			method:  http.MethodPost,
			headers: http.Header{"Content-Type": []string{"application/json"}},
			body:    `{"foo": "bar"}`,
		},
		{
			name: "base64 decode field",
			want: `method: POST
headers:
  Accept-Encoding:
  - gzip
  Content-Length:
  - "15"
  Content-Type:
  - application/json
  User-Agent:
  - Go-http-client/1.1
body: |-
  {
    "foo": "bar"
  }
transform:
  base64:
  - foo
`,
			method:  http.MethodPost,
			headers: http.Header{"Content-Type": []string{"application/json"}},
			body:    `{"foo": "YmFy"}`,
			opts:    []hook.Option{hook.DecodeOption(hook.Base64Transformer{}, "foo")},
		},
		{
			name: "base64 decode struct",
			want: `method: POST
headers:
  Accept-Encoding:
  - gzip
  Content-Length:
  - "31"
  Content-Type:
  - application/json
  User-Agent:
  - Go-http-client/1.1
body: |-
  {
    "foo": {
      "bar": "baz"
    }
  }
transform:
  base64:
  - foo
`,
			method:  http.MethodPost,
			headers: http.Header{"Content-Type": []string{"application/json"}},
			body:    `{"foo": "eyJiYXIiOiAiYmF6In0="}`,
			opts:    []hook.Option{hook.DecodeOption(hook.Base64Transformer{}, "foo")},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f := testfile(t, "hook.yml")
			defer deletefile(t, f)

			r, err := newRecorder(f.Name(), tc.opts...)
			if err != nil {
				t.Fatal(err)
			}
			defer r.close()

			srv := httptest.NewServer(r)
			defer srv.Close()

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
		})
	}
}
