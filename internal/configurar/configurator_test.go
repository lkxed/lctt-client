package configurar

import (
	"lctt-client/internal/helper"
	"path"
	"testing"
)

func TestLoadWebsites(t *testing.T) {
	loadWebsites(path.Join(helper.ConfigDir, "websites.yml"))
	for k, v := range Websites {
		t.Logf("%s: {\n", k)
		t.Logf("  title: %s,\n", v.Title)
		t.Logf("  summary: %s,\n", v.Summary)
		t.Logf("  author: %s,\n", v.Author)
		t.Logf("  date: %s,\n", v.Date)
		t.Logf("  content: %s\n", v.Content)
		t.Logf("}\n\n")
	}
}

func TestLoadSettings(t *testing.T) {
	loadSettings(path.Join(helper.ConfigDir, "settings.yml"))
	t.Log(Settings)
}
