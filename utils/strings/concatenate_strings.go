package string

import "strings"

// Concatenate any number of strings provided as slice of strings in the same order
func ConcatenateStrings(strs ...string) string {
	return strings.Join(strs, "")
}
