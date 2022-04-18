package gitter

import (
	"errors"
	"fmt"
	"lctt-client/internal/helper"
	"log"
	"path"
	"strings"
	"unicode"
)

// ReplaceUrls replaces original urls with linux.cn urls.
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

// Collect an article.
func Collect(category string, filename string) {
	open()
	pull()

	branch := initBranch(filename)
	checkout(branch)

	previewPath := path.Join(helper.PreviewDir, filename)
	relativePath := path.Join("sources", category, filename)
	filepath := path.Join(LocalRepository, relativePath)
	helper.Move(previewPath, filepath)

	checkWorkTreeStatus()

	add(relativePath)

	message := formatCommitMessage("手动选题", category, filename)
	commit(message)

	push(branch)

	title := formatRequestTitle("手动选题", category, filename)
	body := formatRequestBody("collected")
	createPR(branch, title, body)

	checkout(UpstreamBranch)
}

// Request to translate an article.
func Request(category string, filename string) {
	open()

	branch := initBranch(filename)
	checkout(branch)

	// Update the file, fill in translator's GitHub username.
	relativePath := path.Join("sources", category, filename)
	filepath := path.Join(LocalRepository, relativePath)
	content := string(helper.ReadFile(filepath))
	// For historical reasons, there are two formats:
	// translator: ( )
	// translator: " "
	if strings.Contains(content, `translator: " "`) {
		translator := fmt.Sprintf(`translator: "%s"`, Username)
		content = strings.Replace(content, `translator: " "`, translator, 1)
	} else {
		translator := fmt.Sprintf("translator: (%s)", Username)
		content = strings.Replace(content, `translator: ( )`, translator, 1)
	}
	helper.Write(filepath, []byte(content))

	checkWorkTreeStatus()

	add(relativePath)

	message := formatCommitMessage("申领原文", category, filename)
	commit(message)

	push(branch)

	title := formatRequestTitle("申领原文", category, filename)
	body := formatRequestBody("being translated")
	exists := checkOpenPR(title)
	if !exists {
		createPR(branch, title, body)
	}
	//checkout(UpstreamBranch)
}

func Complete(category string, filename string, force bool) error {
	open()

	branch := initBranch(filename)
	checkout(branch)

	sourcesRelativePath := path.Join("sources", category, filename)
	sourcesPath := path.Join(LocalRepository, sourcesRelativePath)
	content := string(helper.ReadFile(sourcesPath))
	// Decide whether the translation is complete by
	// checking if Chinese characters consist more than 50% of it.
	rest := strings.Split(content, "======")[1]
	translation := strings.Split(rest,
		"--------------------------------------------------------------------------------")[0]
	translation = helper.ClearSpace(translation)
	var count int
	for _, c := range translation {
		if unicode.Is(unicode.Han, c) {
			count++
		}
	}
	zhHansPercentage := float64(count) / float64(len(translation))
	log.Printf("Chinese characters consist %.1f%% of your translation.\n", zhHansPercentage*100)
	if !force && zhHansPercentage < 0.1 {
		return errors.New("translation not completed")
	}

	// In case somebody forgets to change it
	if strings.Contains(content, "译者ID") {
		content = strings.Replace(content, "译者ID", Username, 2)
		helper.Write(sourcesPath, []byte(content))
	}

	translatedRelativePath := path.Join("translated", category, filename)
	translatedPath := path.Join(LocalRepository, translatedRelativePath)
	helper.Move(sourcesPath, translatedPath)

	checkWorkTreeStatus()

	add(sourcesRelativePath)
	add(translatedRelativePath)

	message := formatCommitMessage("提交译文", category, filename)
	commit(message)

	push(branch)

	title := formatRequestTitle("提交译文", category, filename)
	body := formatRequestBody("translated")
	exists := checkOpenPR(title)
	if !exists {
		createPR(branch, title, body)
	}

	checkout(UpstreamBranch)
	return nil
}

func initBranch(filename string) string {
	branch := strings.ReplaceAll(LocalBranch, "<FILENAME>", filename)
	branch = reformatBranch(branch)
	if !hasLocalBranch(branch) {
		createLocalBranch(branch)
	}
	return branch
}

func checkWorkTreeStatus() {
	log.Println("Checking worktree status...")
	if isClean, changes := inspectStatus(); isClean {
		log.Fatalln("No changes since last commit.")
	} else {
		fmt.Println(changes)
	}
}
