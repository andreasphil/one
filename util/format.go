package util

import "strings"

// FormatNestedMap flattens m into a space separated list of key=value pairs,
// with one pair per value. The order of the pairs is not stable, because it
// follows the iteration order of m.
func FormatNestedMap(m map[string][]string) string {
	result := strings.Builder{}

	for k, vv := range m {
		for _, v := range vv {
			result.WriteString(k)
			result.WriteString("=")
			result.WriteString(v)
			result.WriteString(" ")
		}
	}

	return strings.TrimSpace(result.String())
}
