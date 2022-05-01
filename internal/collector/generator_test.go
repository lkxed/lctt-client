package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://opensource.com/article/22/4/go-build-options")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
