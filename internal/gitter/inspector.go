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
	helper.ExitIfError(err)
	return name == head.Name()
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
		first := s[0]
		switch first {
		case '?':
			s = strings.Replace(s, "?  ", "[?] `", 1)
			s = strings.Replace(s, "\n", "`", 1)
		case 'A':
			s = strings.Replace(s, "A  ", "[+] `", 1)
			s = strings.Replace(s, "\n", "`", 1)
		case 'D':
			s = strings.Replace(s, "D  ", "[-] `", 1)
			s = strings.Replace(s, "\n", "`", 1)
		case 'M':
			s = strings.Replace(s, "M  ", "[M] `", 1)
			s = strings.Replace(s, "\n", "`", 1)
		}
		s = strings.Replace(s, "\n", "", 1)
	}
	return status.IsClean(), s
}
