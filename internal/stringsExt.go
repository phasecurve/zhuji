// Package internal
package internal

import "strings"

func SplitRemoveEmpty(value, sep string) []string {
	parts := strings.Split(value, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}

func TrimSuffix(srcString string, toRemove rune) string {
	rs := []rune(srcString)
	var replaced string
	if len(rs) > 0 && rs[len(rs)-1] == toRemove {
		replaced = string(rs[:len(rs)-1])
	}
	return replaced
}
