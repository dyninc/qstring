package qstring

import (
	"strings"
)

// parseTag splits a struct field's qstring tag into its name and, if an
// optional omitempty option was provided, a boolean indicating this is
// returned
func parseTag(tag string) (string, bool) {
	if idx := strings.Index(tag, ","); idx != -1 {
		if tag[idx+1:] == "omitempty" {
			return tag[:idx], true
		}
		return tag[:idx], false
	}
	return tag, false
}
