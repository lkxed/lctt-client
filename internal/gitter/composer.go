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
	hasForked := checkOriginRepository()
	if !hasForked {
		fork()
	}
	dotGit := path.Join(LocalRepository, ".git")
	hasCloned := helper.CheckPath(dotGit)
	if !hasCloned {
		helper.MkdirAll(LocalRepository)
		clone()
	} else {
		open()
		pull()
	}
}

// Collect an article.
func Collect(category string, filename string) {
	log.Println("Collecting...")
	open()
	pull()

	branch := initBranch(filename)
	checkout(branch)

	tmpPath := path.Join(helper.TmpDir, filename)
	relativePath := path.Join("sources", category, filename)
	filepath := path.Join(LocalRepository, relativePath)
	helper.Copy(tmpPath, filepath)

	// Check worktree status to make sure changes have been made.
	checkWorkTreeStatus()

	// Add file creation & modification changes.
	add(relativePath)

	message := formatCommitMessage("手动选题", category, filename)
	commit(message)

	push(branch)

	title := formatRequestTitle("手动选题", category, filename)
	body := formatRequestBody("collected")
	createPR(branch, title, body)

	checkout(UpstreamBranch)
	log.Printf("Collected: %s.\n", relativePath)
}

// Request to translate an article.
func Request(category string, filename string) {
	log.Println("Requesting...")
	open()

	branch := initBranch(filename)
	checkout(branch)

	// Update the file, fill in translator's GitHub username.
	relativePath := path.Join("sources", category, filename)
	filepath := path.Join(LocalRepository, relativePath)
	// Copy it to "tmp/" folder and process there
	tmpPath := path.Join(helper.TmpDir, filename)
	helper.Copy(filepath, tmpPath)

	// Process file: Add translator's GitHub ID.
	// For historical reasons, there are two formats:
	// translator: ( )
	// translator: " "
	content := string(helper.ReadFile(tmpPath))
	if strings.Contains(content, `translator: " "`) {
		translator := fmt.Sprintf(`translator: "%s"`, Username)
		content = strings.Replace(content, `translator: " "`, translator, 1)
	} else {
		translator := fmt.Sprintf("translator: (%s)", Username)
		content = strings.Replace(content, `translator: ( )`, translator, 1)
	}
	helper.Write(tmpPath, []byte(content))

	// Copy it back to "sources" for git operations.
	helper.Copy(tmpPath, filepath)

	// Check worktree status to make sure changes have been made.
	checkWorkTreeStatus()

	// Add file modification changes.
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
	log.Printf("Requested %s.\n", relativePath)
}

func Complete(category string, filename string, force bool) error {
	log.Println("Completing...")
	open()

	branch := initBranch(filename)
	checkout(branch)

	sourcesRelativePath := path.Join("sources", category, filename)
	sourcesPath := path.Join(LocalRepository, sourcesRelativePath)
	// Copy it to "tmp/" folder and process there
	tmpPath := path.Join(helper.TmpDir, filename)
	helper.Copy(sourcesPath, tmpPath)

	// Process file: Decide whether the translation is complete by
	// checking if Chinese characters consist more than 15% of it.
	// This is a rough estimation, better algorithms needed.
	content := string(helper.ReadFile(tmpPath))
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
	if !force && zhHansPercentage < 0.15 {
		return errors.New("translation not completed")
	}

	// In case somebody forgets to change it
	if strings.Contains(content, "译者ID") {
		content = strings.Replace(content, "译者ID", Username, 2)
		helper.Write(tmpPath, []byte(content))
	}

	// Copy it to "translated" folder for git operations.
	translatedRelativePath := path.Join("translated", category, filename)
	translatedPath := path.Join(LocalRepository, translatedRelativePath)
	helper.Copy(tmpPath, translatedPath)

	// Remove the source for git operations.
	helper.Remove(sourcesPath)

	// Check worktree status to make sure changes have been made.
	checkWorkTreeStatus()

	// Add file deletion & creation changes.
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

	log.Printf("Completed: %s.\n", translatedRelativePath)
	checkout(UpstreamBranch)
	return nil
}

func Clean() error {
	open()

	filenames, err := helper.ListDir(helper.TmpDir)
	if err != nil || len(filenames) == 0 {
		return errors.New("nothing to clean")
	}
	prs := listOpenPRs()
	var titles []string
	for _, pr := range prs {
		titles = append(titles, *pr.Title)
	}
	for _, filename := range filenames {
		isPROpen := helper.StringSliceContains(titles, filename)
		if !isPROpen {
			cleanBranch(filename)
			// Remove temporary files
			tmpPath := path.Join(helper.TmpDir, filename)
			helper.Remove(tmpPath)
		}
	}

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

func cleanBranch(filename string) {
	branch := reformatBranch(filename)
	_ = deleteLocalBranch(branch)
	_ = deleteOriginBranch(branch)
}

func checkWorkTreeStatus() {
	log.Println("Checking worktree status...")
	if isClean, changes := inspectStatus(); isClean {
		log.Fatalln("No changes since last commit.")
	} else {
		fmt.Println(changes)
	}
}
