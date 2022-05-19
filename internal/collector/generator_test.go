package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://news.itsfoss.com/inkscape-1-2-release/")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
