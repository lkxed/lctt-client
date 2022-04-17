package collector

import "testing"

func TestParse(t *testing.T) {
	article := Parse("https://www.freecodecamp.org/news/how-to-fork-a-github-repository/")
	for i, url := range article.Urls {
		t.Logf("index: %d, url: %s\n", i, url)
	}

	for _, text := range article.Texts {
		t.Log(text)
	}
}
