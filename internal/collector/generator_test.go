package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://www.linuxuprising.com/2022/01/extension-manager-search-and-install.html")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
