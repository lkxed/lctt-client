package helper

import (
	"io/ioutil"
	"os"
	"path"
)

func CheckPath(filepath string) bool {
	_, err := os.Open(filepath)
	return err == nil
}

func MkdirAll(dir string) {
	err := os.MkdirAll(dir, 0755)
	ExitIfError(err)
}

func ListDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	return filenames, nil
}

func Write(filepath string, content []byte) {
	MkdirAll(path.Dir(filepath))
	err := ioutil.WriteFile(filepath, content, 0644)
	ExitIfError(err)
}

func Remove(filepath string) {
	err := os.Remove(filepath)
	ExitIfError(err)
}

func ReadFile(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	ExitIfError(err)
	return content
}

func Copy(oldPath string, newPath string) {
	content := ReadFile(oldPath)
	Write(newPath, content)
}

func Move(oldPath string, newPath string) {
	Copy(oldPath, newPath)
	Remove(oldPath)
}
