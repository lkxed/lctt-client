package gitter

import (
	"log"
	"testing"
)

func TestInspectStatus(t *testing.T) {
	open()
	isClean, status := inspectStatus()
	log.Println(isClean, status)
}
