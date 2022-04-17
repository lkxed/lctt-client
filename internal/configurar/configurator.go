package configurar

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"lctt-client/internal/helper"
)

func loadWebsites(filename string) {
	file, err := ioutil.ReadFile(filename)
	helper.ExitIfError(err)

	err = yaml.Unmarshal(file, &Websites)
	helper.ExitIfError(err)
}

func loadSettings(filename string) {
	file, err := ioutil.ReadFile(filename)
	helper.ExitIfError(err)

	err = yaml.Unmarshal(file, &Settings)
	helper.ExitIfError(err)
}
