package hook

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Hook represents a single hook configuration.
type Hook struct {
	Method  string      `yaml:"method"`
	Headers http.Header `yaml:"headers"`
	Body    string      `yaml:"body"`
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
func NewFromPath(path string) ([]*Hook, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return New(f)
}

// New creates a new Hook from the given bytestring.
func New(r io.Reader) ([]*Hook, error) {
	res := []*Hook{}
	d := yaml.NewDecoder(r)

	var err error
	for err == nil {
		h := new(Hook)
		err = d.Decode(h)
		if err == nil {
			res = append(res, h)
		}
	}
	if err != io.EOF {
		return nil, err
	}
	return res, nil
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

	if h.Body != "" {
		reader := strings.NewReader(h.Body)
		r.Body = ioutil.NopCloser(reader)
	}

	return r, nil
}
