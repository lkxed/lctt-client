package helper

import (
	"os"
	"path"
	"strings"
)

var (
	ClientDir string
	ConfigDir string
	TmpDir    string
)

func init() {
	wd, _ := os.Getwd()
	for len(wd) > 1 && !strings.HasSuffix(wd, "/lctt-client") {
		wd = path.Dir(wd)
	}
	ClientDir = wd

	ConfigDir = path.Join(ClientDir, "configs")
	TmpDir = path.Join(ClientDir, "tmp")
}
