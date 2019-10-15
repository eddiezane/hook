package hook

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

// Hook represents a single hook configuration.
type Hook struct {
	Method  string      `yaml:"method"`
	Headers http.Header `yaml:"headers"`
	Body    string      `yaml:"body"`
	Params  url.Values  `yaml:"params"`
}

// NewFromRequest creates a new Hook from the given HTTP Request.
func NewFromRequest(r *http.Request) (*Hook, error) {
	h := &Hook{
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

// NewFromPath creates a new Hook from the given path.
func NewFromPath(path string) (*Hook, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return New(bs)
}

// New creates a new Hook from the given bytestring.
func New(bs []byte) (*Hook, error) {
	h := &Hook{}
	if err := yaml.Unmarshal(bs, h); err != nil {
		return nil, err
	}
	return h, nil
}

// Dump TODO(eddiezane): Is this the right method?
func (h *Hook) Dump() ([]byte, error) {
	return yaml.Marshal(h)
}

// Fire sends an HTTP request to the given target.
func (h *Hook) Fire(target string) (*http.Response, error) {
	r, err := h.toRequest(target)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(r)
}

// toRequest converts the hook into a HTTP request.
func (h *Hook) toRequest(target string) (*http.Request, error) {
	r, err := http.NewRequest(h.Method, target, nil)
	if err != nil {
		return nil, err
	}

	r.Header = h.Headers
	r.URL.RawQuery = h.Params.Encode()

	if h.Body != "" {
		reader := strings.NewReader(h.Body)
		r.Body = ioutil.NopCloser(reader)
	}

	return r, nil
}
