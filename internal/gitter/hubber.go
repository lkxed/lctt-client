package gitter

import (
	"context"
	"github.com/google/go-github/v43/github"
	"lctt-client/internal/helper"
	"log"
	"path"
)

func fork() *github.Repository {
	log.Println("Forking upstream...")

	owner := UpstreamOwner
	repo := path.Base(UpstreamRepository)
	repository, _, err := client.Repositories.CreateFork(context.Background(), owner, repo, nil)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			log.Printf("Forked: %s\n", path.Join(owner, repo))
		} else {
			helper.ExitIfError(err)
		}
	}
	return repository
}

func checkOpenPR(title string) bool {
	openPRs := listOpenPRs()
	for _, pr := range openPRs {
		if title == *pr.Title {
			return true
		}
	}
	return false
}

func listOpenPRs() []*github.PullRequest {
	owner := UpstreamOwner
	repo := path.Base(UpstreamRepository)
	// Because default "State" filter is "open"
	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, nil)
	helper.ExitIfError(err)
	return prs
}

func createPR(branch string, title string, body string) *github.PullRequest {
	log.Println("Creating pull request...")

	owner := UpstreamOwner
	head := Username + ":" + branch
	base := UpstreamBranch
	repo := path.Base(UpstreamRepository)
	maintainerCanModify := true
	draft := false
	newPR := &github.NewPullRequest{
		Title:               &title,
		Head:                &head,
		Base:                &base,
		Body:                &body,
		MaintainerCanModify: &maintainerCanModify,
		Draft:               &draft,
	}
	pr, _, err := client.PullRequests.Create(context.Background(), owner, repo, newPR)
	helper.ExitIfError(err)

	log.Printf("Created: %s\n", title)
	return pr
}

func deleteOriginRepository() {
	log.Println("Deleting origin repository...")

	owner := Username
	repo := path.Base(OriginRepository)
	_, err := client.Repositories.Delete(context.Background(), owner, repo)
	helper.ExitIfError(err)

	log.Printf("Deleted: %s\n", path.Join(owner, repo))
}

func checkOriginRepository() bool {
	log.Println("Getting origin repository...")

	owner := Username
	repo := path.Base(OriginRepository)
	_, _, err := client.Repositories.Get(context.Background(), owner, repo)

	if err != nil {
		log.Println("Repository not exist.")
		return false
	}
	log.Printf("Got: %s,\n", path.Join(owner, repo))
	return true
}

func searchForExistence(query string) bool {
	result, _, err := client.Search.Code(context.Background(), query, &github.SearchOptions{
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	})
	helper.ExitIfError(err)

	return result.GetTotal() > 0
}
