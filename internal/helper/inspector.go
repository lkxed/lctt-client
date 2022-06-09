package helper

import (
	"regexp"
	"strings"
	"unicode"
)

func StringSliceContains(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}

func ArticleZhHansPercentage(filepath string) float64 {
	// Process file: Decide whether the translation is complete by
	// checking if Chinese bytes consist more than 75% of it.
	// This is a rough estimation, better algorithms needed.
	content := string(ReadFile(filepath))
	rest := strings.Split(content, "======")[1]
	// extract main content
	translation := strings.Split(rest,
		"--------------------------------------------------------------------------------")[0]
	// exclude code blocks
	re := regexp.MustCompile("```[\\w|\\W]*```")
	translation = string(re.ReplaceAll([]byte(translation), []byte{}))
	// exclude spaces
	translation = strings.ReplaceAll(translation, " ", "")
	translation = strings.ReplaceAll(translation, "\n", "")
	translation = strings.ReplaceAll(translation, "\t", "")
	var count int
	for _, c := range translation {
		if unicode.Is(unicode.Han, c) {
			count++
		}
	}
	zhHansPercentage := float64(count) * 3 / float64(len(translation))
	return zhHansPercentage
}

func IsBeingTranslated(filepath string) bool {
	content := string(ReadFile(filepath))
	return !strings.Contains(content, `translator: " "`) &&
		!strings.Contains(content, `translator: ( )`) &&
		strings.Contains(content, "译者ID")
}
