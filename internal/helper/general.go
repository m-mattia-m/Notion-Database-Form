package helper

import "strings"

func IfArrayElementContainsString(needle string, haystack []string) bool {
	for _, key := range haystack {
		if strings.Contains(needle, key) {
			return true
		}
	}
	return false
}
