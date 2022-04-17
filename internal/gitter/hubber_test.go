package gitter

import (
	"testing"
)

func TestFork(t *testing.T) {
	fork()
}

func TestCreatePullRequest(t *testing.T) {
	createPR("branch", "title", "body")
}

func TestDeleteRepository(t *testing.T) {
	deleteOriginRepository()
}
