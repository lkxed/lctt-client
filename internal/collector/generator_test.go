package collector

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	article := Parse("https://ostechnix.com/php-mysql-where-clause/")
	filepath, content := Generate(article)
	fmt.Println(filepath)
	fmt.Println(string(content))
}
