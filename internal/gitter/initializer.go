package gitter

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
	"lctt-client/internal/configurar"
	"path"
	"strings"
)

var (
	auth               *http.BasicAuth
	repository         *git.Repository
	worktree           *git.Worktree
	client             *github.Client
	UserName           string
	UserEmail          string
	LocalRepository    string
	LocalBranch        string
	UpstreamRepository string
	UpstreamOwner      string
	UpstreamBranch     string
	OriginRepository   string
	CommitMessage      string
	Username           string
	AccessToken        string
	RequestTitle       string
	RequestBody        string
)

func init() {
	Git := configurar.Settings.Git
	UserName = Git.User.Name
	UserEmail = Git.User.Email
	LocalRepository = Git.Local.Repository
	LocalBranch = Git.Local.Branch
	UpstreamRepository = Git.Remote.Upstream.Repository
	UpstreamOwner = path.Base(path.Dir(UpstreamRepository))
	UpstreamBranch = Git.Remote.Upstream.Branch
	CommitMessage = Git.Commit.Message
	Username = Git.Hub.Username
	OriginRepository = strings.Replace(UpstreamRepository, UpstreamOwner, Username, 1)
	AccessToken = Git.Hub.AccessToken
	RequestTitle = Git.Hub.PullRequest.Title
	RequestBody = Git.Hub.PullRequest.Body

	auth = &http.BasicAuth{
		Username: UserName,
		Password: AccessToken,
	}

	client = github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: AccessToken,
	})))
}
