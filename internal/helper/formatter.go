package helper

import (
	"regexp"
	"strings"
)

func ConcatUrl(baseUrl string, url string) string {
	url = strings.ReplaceAll(url, " ", "%20")
	if len(url) == 0 || strings.HasPrefix(url, "http") {
		return url
	} else if strings.HasPrefix(url, "/") {
		return baseUrl + url
	} else {
		return baseUrl + "/" + url
	}
}

func ConcatFilename(date string, title string) string {
	// Windows filename can't contain one of these characters: \ / : * ? " < > |.
	// This will replace other related punctuations as well just in case.
	re := regexp.MustCompile(`[\\/:*?"'’‘“”()<>|]`)
	title = string(re.ReplaceAll([]byte(title), []byte("-")))
	return date + " " + title + ".md"
}
