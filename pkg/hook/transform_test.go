package hook

import "testing"

func TestTransformBase64(t *testing.T) {
	decoded := `{"foo":"bar"}`
	path := "foo"
	encoded := `{"foo":"YmFy"}`

	b64t := Base64Transformer{}

	if b64t.Type() != TransformBase64 {
		t.Errorf("Type: want %s, got %s", TransformBase64, b64t.Type())
	}

	encodeOut, err := b64t.Encode(decoded, path)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if encoded != encodeOut {
		t.Errorf("Encode: want %s, got %s", encoded, encodeOut)
	}

	decodeOut, err := b64t.Decode(encoded, path)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if decoded != decodeOut {
		t.Errorf("Decode: want %s, got %s", decoded, decodeOut)
	}
}
