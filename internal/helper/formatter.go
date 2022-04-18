package helper

import (
	"strings"
)

func ClearSpace(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.TrimSpace(s)
	return s
}

func TrimSpace(s string) string {
	s = strings.TrimSpace(s)
	for strings.HasPrefix(s, "\n") || strings.HasPrefix(s, "\t") {
		s = strings.TrimLeft(s, "\n")
		s = strings.TrimLeft(s, "\t")
	}
	for strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\t") {
		s = strings.TrimRight(s, "\n")
		s = strings.TrimRight(s, "\t")
	}
	return s
}

func ConcatUrl(baseUrl string, url string) string {
	if len(url) == 0 || strings.HasPrefix(url, "http") {
		return url
	} else if strings.HasPrefix(url, "/") {
		return baseUrl + url
	} else {
		return baseUrl + "/" + url
	}
}

func ConcatFilename(date string, title string) string {
	return date + " " + title + ".md"
}
