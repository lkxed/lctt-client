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

// ReformatBranch according to https://mirrors.edge.kernel.org/pub/software/scm/git/docs/git-check-ref-format.html
func ReformatBranch(branch string) string {
	branch = strings.TrimRight(branch, ".md")
	// 1. They cannot end with the sequence .lock
	branch = strings.TrimRight(branch, ".lock")
	// 3. They cannot have two consecutive dots .. anywhere.
	branch = strings.ReplaceAll(branch, "..", "-")
	// 4. They cannot have ASCII control characters, space, tilde ~, caret ^, or colon : anywhere
	branch = strings.ReplaceAll(branch, " ", "-")
	branch = strings.ReplaceAll(branch, "~", "-")
	branch = strings.ReplaceAll(branch, "^", "-")
	branch = strings.ReplaceAll(branch, ":", "-")
	// 5. They cannot have question-mark ?, asterisk *, or open bracket [ anywhere
	branch = strings.ReplaceAll(branch, "?", "-")
	branch = strings.ReplaceAll(branch, "*", "-")
	branch = strings.ReplaceAll(branch, "[", "-")
	// 6. They cannot begin or end with a slash / or contain multiple consecutive slashes
	branch = strings.TrimRight(branch, "/")
	branch = strings.ReplaceAll(branch, "//", "-")
	// 7. They cannot end with a dot .
	branch = strings.TrimRight(branch, ".")
	// 8. They cannot contain a sequence @{
	branch = strings.ReplaceAll(branch, "@{", "-")
	// 10. They cannot contain a \
	branch = strings.ReplaceAll(branch, "\\", "-")

	return branch
}
