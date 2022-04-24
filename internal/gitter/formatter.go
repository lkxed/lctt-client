package gitter

import (
	"strings"
)

// ReformatBranch according to https://mirrors.edge.kernel.org/pub/software/scm/git/docs/git-check-ref-format.html
func reformatBranch(branch string) string {
	branch = strings.TrimSuffix(branch, ".md")
	// 1. They cannot end with the sequence .lock
	branch = strings.TrimSuffix(branch, ".lock")
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
	branch = strings.TrimPrefix(branch, "/")
	branch = strings.TrimSuffix(branch, "/")
	branch = strings.ReplaceAll(branch, "//", "-")
	// 7. They cannot end with a dot .
	branch = strings.TrimSuffix(branch, ".")
	// 8. They cannot contain a sequence @{
	branch = strings.ReplaceAll(branch, "@{", "-")
	// 10. They cannot contain a \
	branch = strings.ReplaceAll(branch, "\\", "-")

	return branch
}

func formatCommitMessage(action string, category string, filename string) string {
	message := CommitMessage
	message = strings.ReplaceAll(message, "<ACTION>", action)
	message = strings.ReplaceAll(message, "<CATEGORY>", category)
	message = strings.ReplaceAll(message, "<FILENAME>", filename)
	return message
}

func formatRequestTitle(action string, category string, filename string) string {
	title := RequestTitle
	title = strings.ReplaceAll(title, "<ACTION>", action)
	title = strings.ReplaceAll(title, "<CATEGORY>", category)
	title = strings.ReplaceAll(title, "<FILENAME>", filename)
	return title
}

func formatRequestBody(status string) string {
	body := RequestBody
	body = strings.ReplaceAll(body, "<STATUS>", status)
	body = strings.ReplaceAll(body, "<USERNAME>", Username)
	return body
}
