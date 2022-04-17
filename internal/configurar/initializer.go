package configurar

import (
	"lctt-client/internal/helper"
	"path"
)

type Website struct {
	Title     string `yaml:"title"`
	Summary   string `yaml:"summary"`
	Author    string `yaml:"author"`
	Date      string `yaml:"date"`
	Content   string `yaml:"content"`
	Exclusion string `yaml:"exclusion"`
	Feed      string `yaml:"feed"`
}

type settings struct {
	Git struct {
		User struct {
			Name  string `yaml:"name"`
			Email string `yaml:"email"`
		}
		Local struct {
			Repository string `yaml:"repository"`
			Branch     string `yaml:"branch"`
		}
		Remote struct {
			Upstream struct {
				Repository string `yaml:"repository"`
				Branch     string `yaml:"branch"`
			}
		}
		Commit struct {
			Message string `yaml:"message"`
		}
		Hub struct {
			Username    string `yaml:"username"`
			AccessToken string `yaml:"access-token"`
			PullRequest struct {
				Title string `yaml:"title"`
				Body  string `yaml:"body"`
			} `yaml:"pull-request"`
		}
	}
	Editor  string `yaml:"editor"`
	Browser string `yaml:"browser"`
}

var (
	Websites map[string]Website
	Settings settings
)

func init() {
	loadWebsites(path.Join(helper.ConfigDir, "websites.yml"))
	loadSettings(path.Join(helper.ConfigDir, "settings.yml"))
}
