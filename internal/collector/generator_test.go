package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://opensource.com/article/22/5/gnu-c-compiler")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
