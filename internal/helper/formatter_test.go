package helper

import (
	"fmt"
	"testing"
)

func TestTrimLeft(t *testing.T) {
	s := "toTrim, Hello World!"
	s = TrimLeft(s, "toTrim")
	fmt.Println(s)
}

func TestTrimRight(t *testing.T) {
	s := "Hello world! toTrim"
	s = TrimRight(s, "toTrim")
	fmt.Println(s)
}
