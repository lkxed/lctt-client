package gitter

import (
	"errors"
	"fmt"
	"lctt-client/internal/helper"
	"log"
	"path"
	"strings"
	"time"
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
	open()
	checkout(UpstreamBranch)
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
	_ = add(relativePath)

	message := formatCommitMessage("手动选题", category, filename)
	commit(message)

	push(branch)

	title := formatRequestTitle("手动选题", category, filename)
	body := formatRequestBody("collected")
	exists := checkOpenPRContains(filename)
	if !exists {
		createPR(branch, title, body)
	}

	checkout(UpstreamBranch)
	log.Printf("Collected: %s.\n", relativePath)
}

// List sources (of a certain category).
func List(category string, since *time.Time, all bool) []string {
	contentPath := path.Join("sources", category)

	var nameFilter func(s string) bool
	if since == nil {
		nameFilter = nil
	} else {
		nameFilter = func(s string) bool {
			date, err := time.Parse("20060102", s[:8])
			if err != nil {
				return false
			}
			return date.After(*since)
		}
	}

	var contentFilter func(s string) bool
	if all {
		contentFilter = nil
	} else {
		contentFilter = func(s string) bool {
			return strings.Contains(s, `translator: " "`) || strings.Contains(s, `translator: ( )`)
		}
	}
	
	filenames, err := getDirFilenames(contentPath, nameFilter, contentFilter)
	if err != nil {
		return nil
	}
	return filenames
}

// Request to translate an article.
func Request(category string, filename string) {
	relativePath := path.Join("sources", category, filename)
	filepath := path.Join(LocalRepository, relativePath)
	exists := helper.CheckPath(filepath)
	if !exists {
		log.Fatalf("Invalid path: %s.\n", filepath)
	}

	open()
	checkout(UpstreamBranch)

	branch := initBranch(filename)
	checkout(branch)

	// Update the file, fill in translator's GitHub username.
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
	_ = add(relativePath)

	message := formatCommitMessage("申领原文", category, filename)
	commit(message)

	push(branch)

	title := formatRequestTitle("申领原文", category, filename)
	body := formatRequestBody("being translated")
	exists = checkOpenPRContains(filename)
	if !exists {
		createPR(branch, title, body)
	}
	//checkout(UpstreamBranch)
	log.Printf("Requested %s.\n", relativePath)
}

func Complete(category string, filename string, force bool) error {
	tmpPath := path.Join(helper.TmpDir, filename)
	exists := helper.CheckPath(tmpPath)
	if !exists {
		log.Fatalf("Invalid path: %s.\n", tmpPath)
	}

	open()
	checkout(UpstreamBranch)

	branch := initBranch(filename)
	checkout(branch)

	zhHansPercentage := helper.ArticleZhHansPercentage(tmpPath)
	log.Printf("Chinese characters consist %.1f%% of your translation.\n", zhHansPercentage*100)
	if !force && zhHansPercentage <= 0.70 {
		return errors.New("translation not completed")
	}

	content := string(helper.ReadFile(tmpPath))

	// In case somebody forgets to change it
	if strings.Contains(content, "译者ID") {
		content = strings.Replace(content, "译者ID", Username, 2)
		helper.Write(tmpPath, []byte(content))
	}

	// Copy it to "translated" folder for git operations.
	translatedRelativePath := path.Join("translated", category, filename)
	translatedPath := path.Join(LocalRepository, translatedRelativePath)
	helper.Copy(tmpPath, translatedPath)

	title := formatRequestTitle("提交译文", category, filename)
	exists = checkOpenPRContains(filename)
	existsComplete := checkOpenPRContains(title)

	sourcesRelativePath := path.Join("sources", category, filename)
	sourcesPath := path.Join(LocalRepository, sourcesRelativePath)

	if !existsComplete {
		// Remove the source for git operations.
		helper.Remove(sourcesPath)
	}

	// Check worktree status to make sure changes have been made.
	checkWorkTreeStatus()

	if !existsComplete {
		// Add file deletion & creation changes.
		err := add(sourcesRelativePath)
		helper.ExitIfError(err)
	}

	err := add(translatedRelativePath)
	helper.ExitIfError(err)

	action := "提交译文"
	if existsComplete {
		action = "修改译文"
	}
	message := formatCommitMessage(action, category, filename)
	commit(message)

	push(branch)

	if !exists {
		body := formatRequestBody("translated")
		createPR(branch, title, body)
	}

	log.Printf("Completed: %s.\n", translatedRelativePath)
	checkout(UpstreamBranch)
	return nil
}

func Clean() error {
	open()
	checkout(UpstreamBranch)

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
		filepath := path.Join(helper.TmpDir, filename)
		isBeingTranslated := helper.IsBeingTranslated(filepath)
		if !isPROpen && !isBeingTranslated {
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
