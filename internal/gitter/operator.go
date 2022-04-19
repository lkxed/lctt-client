package gitter

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"lctt-client/internal/helper"
	"log"
	"os"
	"time"
)

func open() {
	r, err := git.PlainOpen(LocalRepository)
	helper.ExitIfError(err)
	repository = r
	worktree, err = r.Worktree()
	helper.ExitIfError(err)
}

func checkout(branch string) {
	inspectOpened()
	if isCurrentBranch(branch) {
		return
	}
	log.Println("Checking out... ")

	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		// since git is transparent to users, this is allowed
		Force: true,
	})
	helper.ExitIfError(err)

	log.Printf("Checked out: %s\n", branch)
}

func createLocalBranch(branch string) {
	inspectOpened()
	log.Println("Creating branch... ")

	ref, err := repository.Head()
	helper.ExitIfError(err)
	ref = plumbing.NewHashReference(plumbing.NewBranchReferenceName(branch), ref.Hash())
	helper.ExitIfError(repository.Storer.SetReference(ref))
	log.Printf("Created: %s\n", branch)
}

func deleteLocalBranch(branch string) error {
	inspectOpened()

	log.Println("Deleting local branch...")

	ref := plumbing.NewBranchReferenceName(branch)
	err := repository.Storer.RemoveReference(ref)

	if err != nil {
		return err
	}
	log.Printf("Deleted: %s\n", branch)
	return nil
}

func deleteOriginBranch(branch string) error {
	inspectOpened()
	log.Println("Deleting origin branch...")

	remote, err := repository.Remote("origin")
	if err != nil {
		return err
	}
	err = remote.Push(&git.PushOptions{
		Auth:     auth,
		RefSpecs: []config.RefSpec{config.RefSpec(":refs/heads/" + branch)},
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	log.Printf("Deleted: %s\n", branch)
	return nil
}

func clone() {
	log.Println("Cloning repository...")

	var err error
	repository, err = git.PlainClone(LocalRepository, false, &git.CloneOptions{
		URL:               OriginRepository,
		SingleBranch:      true,
		ReferenceName:     plumbing.HEAD,
		Auth:              auth,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          os.Stdout,
	})
	helper.ExitIfError(err)

	log.Printf("Cloned: into %s\n", LocalRepository)

	worktree, err = repository.Worktree()
	helper.ExitIfError(err)
}

func add(relativePath string) {
	inspectOpened()

	// git add $filepath
	_, err := worktree.Add(relativePath)
	helper.ExitIfError(err)
}

func commit(message string) {
	inspectOpened()
	log.Println("Committing changes...")

	// git status
	if isClean, changes := inspectStatus(); isClean {
		log.Fatalln("No changes to commit.")
	} else {
		log.Println(changes)
	}

	// git commit -m $message
	_, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  UserName,
			Email: UserEmail,
			When:  time.Now(),
		},
	})
	helper.ExitIfError(err)

	log.Println("Committed.")
}

func pull() {
	inspectOpened()
	log.Println("Pulling changes...this may take a while...")

	err := worktree.Pull(&git.PullOptions{
		RemoteName:        "origin",
		SingleBranch:      true,
		Auth:              auth,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          os.Stdout,
	})
	if err == nil {
		log.Printf("Pulled: from %s:%s\n", UpstreamOwner, UpstreamBranch)
	} else if err == git.NoErrAlreadyUpToDate {
		log.Println("Already up-to-date.")
	} else {
		helper.ExitIfError(err)
	}
}

func push(branch string) {
	inspectOpened()
	log.Println("Pushing changes...this may take a while...")
	ref := plumbing.NewBranchReferenceName(branch)
	refSpec := config.RefSpec(fmt.Sprintf("+%s:%s", ref, ref))
	err := repository.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       auth,
		Progress:   os.Stdout,
		//Prune:    true,
		//Force: 	true,
	})
	if err == nil {
		log.Printf("Pushed: to %s:%s\n.", UpstreamOwner, UpstreamBranch)
	} else if err == git.NoErrAlreadyUpToDate {
		log.Println("Already up-to-date.")
	} else {
		helper.ExitIfError(err)
	}
}
