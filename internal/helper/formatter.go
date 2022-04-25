package helper

import (
	"regexp"
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
		s = strings.TrimPrefix(s, "\n")
		s = strings.TrimPrefix(s, "\t")
	}
	for strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\t") {
		s = strings.TrimSuffix(s, "\n")
		s = strings.TrimSuffix(s, "\t")
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
	//  Windows filename can't contain one of these characters: \ / : * ? " < > |
	re := regexp.MustCompile(`[/\\:*?<>|]`)
	title = string(re.ReplaceAll([]byte(title), []byte("-")))
	return date + " " + title + ".md"
}
