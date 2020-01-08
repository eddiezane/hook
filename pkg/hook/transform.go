package hook

import (
	"encoding/base64"
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	// Transformers are the default set of transformers.
	Transformers = map[TransformStrategy]Transformer{
		TransformBase64: Base64Transformer{},
	}
)

// TransformStrategy denotes an operation to take on a field.
type TransformStrategy string

const (
	// TransformBase64 denotes that the field should be base64 encoded/decoded.
	TransformBase64 TransformStrategy = "base64"
)

// Transformer defines the encoding and decoding methods for message
// transformation.
type Transformer interface {
	Encode(json string, path string) (string, error)
	Decode(json string, path string) (string, error)
	Type() TransformStrategy
}

// Base64Transformer handles base64 transformations.
type Base64Transformer struct{}

// Type returns the strategy type of the transformer.
func (Base64Transformer) Type() TransformStrategy {
	return TransformBase64
}

// Encode takes the given payload + path and base64 encodes the value.
func (Base64Transformer) Encode(json string, path string) (string, error) {
	in := []byte(gjson.Get(json, path).String())
	out := base64.StdEncoding.EncodeToString(in)
	return sjson.Set(json, path, out)
}

// Decode takes the given payload + path and base64 decodes the value.
func (Base64Transformer) Decode(raw string, path string) (string, error) {
	in := gjson.Get(raw, path).String()
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}

	var js map[string]interface{}
	if json.Unmarshal(out, &js) == nil {
		return sjson.SetRaw(raw, path, string(out))
	}
	return sjson.Set(raw, path, out)
}

type decodeOption struct {
	transformer Transformer
	paths       []string
}

// DecodeOption modifies newly hooks by applying the transformer.Decode
// for the specified paths.
func DecodeOption(t Transformer, paths ...string) Option {
	return &decodeOption{
		transformer: t,
		paths:       paths,
	}
}

func (t *decodeOption) Apply(h *Hook) error {
	// Set transform metadata in hook.
	if h.Transform == nil {
		h.Transform = make(map[TransformStrategy][]string)
	}
	tt := t.transformer.Type()
	if h.Transform[tt] == nil {
		h.Transform[tt] = t.paths
	} else {
		h.Transform[tt] = append(h.Transform[tt], t.paths...)
	}

	// Apply transformation.
	for _, p := range t.paths {
		body, err := t.transformer.Decode(h.Body, p)
		if err != nil {
			return err
		}
		h.Body = body
	}
	return nil
}
