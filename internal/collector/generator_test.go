package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://www.opensourceforu.com/2022/05/esi-group-collaborates-with-ensam-open-sources-its-inspector-software/")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
