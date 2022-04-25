package helper

import (
	"fmt"
	"testing"
)

func TestConcatFilename(t *testing.T) {
	date := "20220425"
	title := "LKXED: Who is he\\she?"
	want := "20220425 LKXED- Who is he-she-.md"
	got := ConcatFilename(date, title)
	if got != want {
		t.Fatalf("Want %s, but got %s", want, got)
	}
	fmt.Println(got)
}
