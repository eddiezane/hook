package hook

import (
	"fmt"
	"io"
)

// ShowHook writes the given hooks specified by the URI to a Writer.
func ShowHook(w io.Writer, uri ...string) error {
	for _, u := range uri {
		hooks, err := NewFromPath(u)
		if err != nil {
			return err
		}
		for _, h := range hooks {
			b, err := h.Dump()
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "---\n%s", b)
		}
	}
	return nil
}
