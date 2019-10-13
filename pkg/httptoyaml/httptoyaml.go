package httptoyaml

import (
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

type HTTPRequest struct {
	Method  string              `yaml:"method"`
	Headers map[string][]string `yaml:"headers"`
	Body    string              `yaml:"body"`
}

// Marshal marshals an http request into a struct
func Marshal(r *http.Request) (*HTTPRequest, error) {
	h := &HTTPRequest{
		Method:  r.Method,
		Headers: r.Header,
	}

	if r.Body != http.NoBody {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		h.Body = string(b)
	}
	return h, nil
}

// Unmarshal unmarshals a strcut into an http request
func Unmarshal(h *HTTPRequest) (*http.Request, error) {
	r, err := http.NewRequest(h.Method, "", nil)
	if err != nil {
		return nil, err
	}

	r.Header = h.Headers

	if h.Body != "" {
		reader := strings.NewReader(h.Body)
		rc := ioutil.NopCloser(reader)
		r.Body = rc
	}

	return r, nil
}

// Dump TODO(eddiezane): Is this the right method?
func (h *HTTPRequest) Dump() ([]byte, error) {
	out, err := yaml.Marshal(h)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func Slurp(bs []byte) (*HTTPRequest, error) {
	h := &HTTPRequest{}
	err := yaml.Unmarshal(bs, h)
	if err != nil {
		return nil, err
	}
	return h, nil
}
