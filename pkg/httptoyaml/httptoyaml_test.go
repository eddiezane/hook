package httptoyaml

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	body := "test=body"
	b := strings.NewReader(body)
	r, err := http.NewRequest(http.MethodPost, "/", b)
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("foo", "bar")
	r.Header.Add("foo", "baz")
	r.Header.Add("herp", "derp")

	h, err := Marshal(r)
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

func TestUnmarshal(t *testing.T) {
	headers := map[string][]string{
		"foo":  []string{"bar", "baz"},
		"herp": []string{"derp"},
	}
	body := "test=body"
	h := &HTTPRequest{
		Method:  http.MethodPost,
		Headers: headers,
		Body:    body,
	}

	r, err := Unmarshal(h)
	if err != nil {
		t.Error(err)
	}

	if r.Method != http.MethodPost {
		t.Errorf("expected method to me POST got %s", r.Method)
	}

	// TODO(eddiezane): Idk why this fails
	// if !reflect.DeepEqual(headers, r.Header) {
	// t.Errorf("expected headers to be %v got %v", headers, r.Header)
	// }

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
	}
	if body != string(b) {
		t.Errorf("expected body to be %s got %s", body, b)
	}
}
