package hook

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNew_http(t *testing.T) {
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

	if reflect.DeepEqual(url.Values{}, r.URL.Query()) {
		t.Errorf("expected params to be %v got %v", url.Values{}, r.URL.Query())
	}
}

func TestNew_byte(t *testing.T) {
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
	bs := []byte(yml)
	h, err := New(bs)
	if err != nil {
		t.Error(err)
	}
	if h == nil {
		t.Error("expected h to be defined")
	}

	yml = `: bad yaml`
	bs = []byte(yml)
	h, err = New(bs)
	if err == nil {
		t.Error("expected error but got nil")
	}
	if h != nil {
		t.Errorf("expected Hook to be nil but got %v", h)
	}
}
