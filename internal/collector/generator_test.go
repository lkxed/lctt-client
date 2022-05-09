package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://news.itsfoss.com/tails-5-0-release/")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
