package hook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewFromRequest(t *testing.T) {
	body := "test=body"
	b := strings.NewReader(body)
	r, err := http.NewRequest(http.MethodPost, "/", b)
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
		t.Errorf("expected body to be %s go %s", body, h.Body)
	}
}

func TestToRequest(t *testing.T) {
	headers := http.Header{
		"foo":  []string{"bar", "baz"},
		"herp": []string{"derp"},
	}
	body := "test=body"
	h := &Hook{
		Method:  http.MethodPost,
		Headers: headers,
		Body:    body,
	}

	r, err := h.toRequest("localhost")
	if err != nil {
		t.Error(err)
	}

	if r.Method != http.MethodPost {
		t.Errorf("expected method to me POST got %s", r.Method)
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
