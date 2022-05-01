package gitter

import (
	"github.com/go-git/go-git/v5/plumbing"
	"lctt-client/internal/helper"
	"log"
	"strings"
)

func notOpened() bool {
	return repository == nil || worktree == nil
}

func inspectOpened() {
	if notOpened() {
		log.Fatalln("Repository not opened.")
	}
}

func isCurrentBranch(branch string) bool {
	name := plumbing.NewBranchReferenceName(branch)
	head, err := repository.Head()
	if err != nil || name != head.Name() {
		return false
	}
	return true
}

func hasBranch(name plumbing.ReferenceName) bool {
	var exists bool
	if iter, err := repository.Branches(); err == nil {
		iterFunc := func(reference *plumbing.Reference) error {
			if name == reference.Name() {
				exists = true
				return nil
			}
			return nil
		}
		_ = iter.ForEach(iterFunc)
	}
	return exists
}

func hasLocalBranch(branch string) bool {
	inspectOpened()
	name := plumbing.NewBranchReferenceName(branch)
	return hasBranch(name)
}

func inspectStatus() (bool, string) {
	inspectOpened()
	status, err := worktree.Status()
	helper.ExitIfError(err)
	var s string
	if !status.IsClean() {
		s = status.String()
		s = strings.ReplaceAll(s, "?  ", "[?] ")
		s = strings.ReplaceAll(s, "A  ", "[+] ")
		s = strings.ReplaceAll(s, "D  ", "[-] ")
		s = strings.ReplaceAll(s, "M  ", "[M] ")
	}
	s = strings.TrimSpace(s)
	return status.IsClean(), s
}
