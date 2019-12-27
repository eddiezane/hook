package hook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewFromRequest(t *testing.T) {
	body := "test=body"
	b := strings.NewReader(body)
	r, err := http.NewRequest(http.MethodPost, "http://localhost?query=value&other=one&other=two", b)
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("foo", "bar")
	r.Header.Add("foo", "baz")
	r.Header.Add("herp", "derp")

	h, err := NewFromRequest(r)
	if err != nil {
		t.Error(err)
	}

	if h.Method != http.MethodPost {
		t.Error("expected method to be POST")
	}

	for k, values := range r.Header {
		for i, v := range values {
			if j := h.Headers[k][i]; j != v {
				t.Errorf("expected header %s to equal %s but got %s", k, v, j)
			}
		}
	}

	if h.Body != body {
		t.Errorf("expected body to be %s got %s", body, h.Body)
	}

	q := url.Values{
		"query": {"value"},
		"other": {"one", "two"},
	}
	if !reflect.DeepEqual(q, h.Params) {
		t.Errorf("expected params to be %v got %v", q, h.Params)
	}
}

func TestToRequest(t *testing.T) {
	headers := http.Header{
		"foo":  []string{"bar", "baz"},
		"herp": []string{"derp"},
	}
	body := "test=body"
	params := url.Values{
		"foo":  []string{"bar", "baz"},
		"taco": []string{"cat"},
	}

	h := &Hook{
		Method:  http.MethodPost,
		Headers: headers,
		Body:    body,
		Params:  params,
	}

	r, err := h.toRequest("localhost")
	if err != nil {
		t.Error(err)
	}

	if r.Method != http.MethodPost {
		t.Errorf("expected method to be POST got %s", r.Method)
	}

	if !reflect.DeepEqual(params, r.URL.Query()) {
		t.Errorf("expected params to be %v got %v", params, r.URL.Query())
	}

	if !reflect.DeepEqual(headers, r.Header) {
		t.Errorf("expected headers to be %v got %v", headers, r.Header)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
	}
	if body != string(b) {
		t.Errorf("expected body to be %s got %s", body, b)
	}
}

func TestToRequest_empty_params(t *testing.T) {
	h := &Hook{}

	r, err := h.toRequest("localhost")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(url.Values{}, r.URL.Query()) {
		t.Errorf("expected params to be %v got %v", url.Values{}, r.URL.Query())
	}
}

func TestNew(t *testing.T) {
	yml := `
method: POST
headers:
  foo:
  - bar
  - baz
  herp:
  - derp
body: test=body
`
	hook := &Hook{
		Method: http.MethodPost,
		Headers: http.Header{
			"foo":  []string{"bar", "baz"},
			"herp": []string{"derp"},
		},
		Body: "test=body",
	}

	testcases := []struct {
		name string
		yml  string
		want []*Hook
	}{
		{
			name: "singledoc",
			yml:  yml,
			want: []*Hook{hook},
		},
		{
			name: "multidoc",
			yml:  fmt.Sprintf("%s---\n%s", yml, yml),
			want: []*Hook{hook, hook},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := New(strings.NewReader(tc.yml))
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(h, tc.want); diff != "" {
				t.Error(diff)
			}
		})
	}

	yml = `: bad yaml`
	h, err := New(strings.NewReader(yml))
	if err == nil {
		t.Error("expected error but got nil")
	}
	if h != nil {
		t.Errorf("expected Hook to be nil but got %v", h)
	}
}

// mockHTTP is a HTTP server that will capture HTTP requests as text.
type mockHTTP struct {
	req []string
}

func (m *mockHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.req = append(m.req, string(b))
	w.WriteHeader(http.StatusOK)
}

func TestFire(t *testing.T) {
	m := &mockHTTP{}
	srv := httptest.NewServer(m)
	defer srv.Close()
	client := srv.Client()
	// Redirect all requests to the fake server.
	// This allows us to send all traffic to the fake server but use
	// deterministic values in the request (i.e. host).
	u, _ := url.Parse(srv.URL)
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(u),
	}

	// TODO(wlynch): We should probably refactor fire to use a custom client.
	http.DefaultClient = client

	paths, err := filepath.Glob(filepath.Join("testdata", "*.hook"))
	if err != nil {
		t.Fatalf("filepath.Glob: %v", err)
	}

	for _, path := range paths {
		m.req = []string{}
		name := strings.TrimSuffix(filepath.Base(path), ".hook")
		t.Run(name, func(t *testing.T) {
			hooks, err := NewFromPath(path)
			if err != nil {
				t.Fatalf("NewFromPath: %v", err)
			}

			for _, h := range hooks {
				if _, err := h.Fire("http://example.com"); err != nil {
					t.Fatalf("Fire(%v): %v", h, err)
				}
			}

			b, err := ioutil.ReadFile(filepath.Join("testdata", name+".http"))
			if err != nil {
				t.Fatalf("ReadFile(http): %v", err)
			}
			want := []string{string(b)}

			if diff := cmp.Diff(want, m.req); diff != "" {
				t.Errorf("diff: %s", diff)
			}
		})
	}
}
