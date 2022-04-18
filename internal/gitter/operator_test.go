package gitter

import (
	"log"
	"testing"
)

func TestOpen(t *testing.T) {
	open()
	inspectOpened()
}

func TestCheckout(t *testing.T) {
	open()
	checkout("master")
	log.Println(inspectStatus())
}

func TestDeleteLocalBranch(t *testing.T) {
	open()
	deleteLocalBranch("test-1")
}

func TestDeleteOriginBranch(t *testing.T) {
	open()
	deleteOriginBranch("test-1")
}

func TestPush(t *testing.T) {
	open()
	push("20220412-Meet-Lite-XL--A-Lightweight,-Open-Source-Text-Editor-for-Linux-Users")
}
