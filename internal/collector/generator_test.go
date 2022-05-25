package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://www.opensourceforu.com/2022/05/the-basic-concepts-of-shell-scripting/")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
