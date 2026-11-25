// Package normalize removes separators from user input.
package normalize

import "strings"

// Digits removes spaces, dots and dashes from s. It returns the result and
// true if the result contains only ASCII digits. An empty result returns true.
func Digits(s string) (string, bool) {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == ' ' || r == '.' || r == '-':
			continue
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			return "", false
		}
	}
	return b.String(), true
}
