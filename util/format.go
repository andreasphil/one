package util

import "strings"

// FormatNestedMap flattens m into a space separated list of key=value pairs,
// with one pair per value.
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
