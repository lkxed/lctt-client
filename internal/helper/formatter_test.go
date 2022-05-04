package helper

import (
	"fmt"
	"testing"
)

func TestConcatFilename(t *testing.T) {
	date := "20220425"
	title := "ESI Group Collaborates With ENSAM, Open Sources Its “Inspector” Software"
	want := "20220425 ESI Group Collaborates With ENSAM, Open Sources Its -Inspector- Software"
	got := ConcatFilename(date, title)
	if got != want {
		t.Fatalf("Want %s, but got %s", want, got)
	}
	fmt.Println(got)
}
