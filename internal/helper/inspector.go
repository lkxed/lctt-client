package helper

import "strings"

func StringSliceContains(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}
