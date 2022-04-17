package gitter

import (
	"lctt-client/internal/helper"
	"log"
	"path"
	"strings"
)

// ReplaceUrls replaces original urls with linux.cn urls
// TODO need support with linux.cn search api, ask @wxy
func ReplaceUrls(urls []string) {
}

func Initialize() {
	if !checkOriginRepository() {
		fork()
	}
	dotGit := path.Join(LocalRepository, ".git")
	if !helper.CheckPath(dotGit) {
		helper.MkdirAll(LocalRepository)
		clone()
	} else {
		open()
		pull()
	}
}

func Collect(category string, filename string) {
	previewPath := path.Join(helper.PreviewDir, filename)
	relativePath := path.Join("sources", category, filename)
	filepath := path.Join(LocalRepository, relativePath)
	branch := strings.ReplaceAll(LocalBranch, "<FILENAME>", filename)
	branch = helper.ReformatBranch(branch)
	open()
	pull()
	if !hasLocalBranch(branch) {
		createLocalBranch(branch)
	}
	checkout(branch)
	helper.Move(previewPath, filepath)
	log.Println("Checking worktree status...")
	if isClean, changes := inspectStatus(); isClean {
		log.Fatalln("No changes since last commit.")
	} else {
		log.Println(changes)
	}
	add(relativePath)
	message := strings.ReplaceAll(CommitMessage, "<ACTION>", "自动选题")
	message = strings.ReplaceAll(message, "<CATEGORY>", category)
	message = strings.ReplaceAll(message, "<FILENAME>", filename)
	commit(message)
	title := strings.ReplaceAll(RequestTitle, "<ACTION>", "自动选题")
	title = strings.ReplaceAll(title, "<CATEGORY>", category)
	title = strings.ReplaceAll(title, "<FILENAME>", filename)
	body := strings.ReplaceAll(RequestBody, "<STATUS>", "collected")
	body = strings.ReplaceAll(RequestBody, "<USERNAME>", Username)
	push(branch)
	createPR(branch, title, body)
}
