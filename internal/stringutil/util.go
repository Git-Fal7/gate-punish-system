package stringutil

import "strings"

// utils
func ReplaceAll(str string, replaceMap map[string]string) string {
	for key, value := range replaceMap {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}
