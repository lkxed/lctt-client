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
		s = TrimLeft(s, "\n")
		s = TrimLeft(s, "\t")
	}
	for strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\t") {
		s = TrimRight(s, "\n")
		s = TrimRight(s, "\t")
	}
	return s
}

func TrimRight(s string, cut string) string {
	lastIndex := strings.LastIndex(s, cut)
	if len(s)-lastIndex == len(cut) {
		return s[:lastIndex]
	}
	return s
}

func TrimLeft(s string, cut string) string {
	index := strings.Index(s, cut)
	if index == 0 {
		return s[len(cut):]
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
