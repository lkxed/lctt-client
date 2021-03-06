package collector

import (
	"bytes"
	"lctt-client/internal/helper"
	"log"
	"path"
	"text/template"
)

func Generate(article Article) (string, []byte) {
	log.Println("Generating markdown...")

	t, err := template.
		New("article.tmpl").
		Funcs(template.FuncMap{
			"add": func(op1 int, op2 int) int {
				return op1 + op2
			},
		}).
		ParseFiles(path.Join(helper.ConfigDir, "article.tmpl"))
	helper.ExitIfError(err)

	var buffer bytes.Buffer
	helper.ExitIfError(t.Execute(&buffer, article))

	filename := helper.ConcatFilename(article.Date, article.Title)
	tmpPath := path.Join(helper.TmpDir, filename)
	helper.Write(tmpPath, buffer.Bytes())

	log.Printf("Generated: %s\n", tmpPath)

	return filename, buffer.Bytes()
}
