package hook

import (
	"strings"
)

// ParsePath takes a hook name of the form <catalog>@<path> and returns the
// individual pieces.
func ParsePath(name string) (catalog string, path string) {
	s := strings.SplitN(name, "@", 2)
	if len(s) == 1 {
		return "", s[0]
	}

	if s[0] == "" && strings.Contains(name, "@") {
		return "@", s[1]
	}

	return s[0], s[1]
}
