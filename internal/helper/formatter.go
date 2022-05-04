package helper

import (
	"regexp"
	"strings"
)

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
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	title = string(re.ReplaceAll([]byte(title), []byte("-")))
	return date + " " + title + ".md"
}
