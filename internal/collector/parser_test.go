package collector

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	article := Parse("https://www.freecodecamp.org/news/how-to-fork-a-github-repository/")
	for i, url := range article.Urls {
		t.Logf("index: %d, url: %s\n", i, url)
	}

	for _, text := range article.Texts {
		t.Log(text)
	}
}

func TestContents(t *testing.T) {
	html :=
		`<html>
			<p>This is a very <em>IMPORTANT</em> message</p>
		</html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	doc.Find("p").Contents().Each(func(i int, s *goquery.Selection) {
		fmt.Println(i, s.Text())
	})
}
