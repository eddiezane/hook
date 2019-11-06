package hook

import "testing"

func TestParsePath(t *testing.T) {
	tests := []struct {
		in, catalog, path string
	}{
		{
			in:      "@foo",
			catalog: "@",
			path:    "foo",
		},
		{
			in:      "foo@bar",
			catalog: "foo",
			path:    "bar",
		},
		{
			in:      "foo",
			catalog: "",
			path:    "foo",
		},
		{
			in:      "foo@",
			catalog: "foo",
			path:    "",
		},
	}
	for _, tc := range tests {
		if catalog, path := ParsePath(tc.in); catalog != tc.catalog || path != tc.path {
			t.Errorf("ParsePath(%s) = (%s, %s), want (%s, %s)", tc.in, catalog, path, tc.catalog, tc.path)
		}
	}
}
