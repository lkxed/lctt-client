package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"lctt-client/internal/collector"
	"lctt-client/internal/configurar"
	"lctt-client/internal/feeder"
	"lctt-client/internal/gitter"
	"lctt-client/internal/helper"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"
)

func initialize(_ *cli.Context) error {
	gitter.Initialize()
	log.Println("Mission completed. Adios!")
	return nil
}

func feed(c *cli.Context) error {
	items := feeder.ParseAll()
	datePtr := c.Timestamp("since")
	var date time.Time
	if datePtr == nil {
		today := time.Now().Format(layout)
		date, _ = time.Parse(layout, today)
	} else {
		date = *datePtr
	}
	items = feeder.FilterPubDate(items, date)
	category := c.String("prefer")
	if len(category) > 0 {
		items = feeder.FilterCategories(items, category)
	}
	if len(items) == 0 {
		log.Fatalf("No article has been published since %s.\n", date.Format(layout))
	}
	if c.Bool("verbose") {
		feeder.ListVerbose(items)
	} else {
		feeder.List(items)
	}
	if c.Bool("open") {
		browser := configurar.Settings.Browser
		if len(browser) == 0 {
			log.Fatalln("No browser specified in `settings.yml`")
		}
		browserCmd := strings.Split(browser, " ")
		browserCmd = append(browserCmd, feeder.ExtractLinks(items)...)
		cmd := exec.Command(browserCmd[0], browserCmd[1:]...)
		log.Println("Opening articles in browser")
		helper.ExitIfError(cmd.Run())
	}
	log.Println("Have you made up your mind? If so, choose an article to `collect`.")
	log.Println("Anyway, mission completed. Adios!")
	return nil
}

func collect(c *cli.Context) error {
	if c.NArg() != 2 {
		log.Fatalln(c.Command.Usage)
	}
	category := c.Args().Get(0)
	categories := []string{"news", "talk", "tech"}
	var contains bool
	for _, v := range categories {
		if category == v {
			contains = true
		}
	}
	if !contains {
		log.Fatalln()
	}
	link := c.Args().Get(1)
	article := collector.Parse(link)
	filename, _ := collector.Generate(article)
	previewPath := path.Join(helper.PreviewDir, filename)
	if !c.Bool("no-preview") {
		editor := configurar.Settings.Editor
		if len(configurar.Settings.Editor) == 0 {
			log.Fatalln("No editor specified in `settings.yml`")
		}
		editCmd := strings.Split(editor, " ")
		editCmd = append(editCmd, previewPath)
		log.Printf("Previewing article in %s...\n", editCmd[0])
		cmd := exec.Command(editCmd[0], editCmd[1:]...)
		log.Println("Modifications have been saved.")
		helper.ExitIfError(cmd.Run())
	}
	if !c.Bool("no-upload") {
		fmt.Print("Are you ready to upload the article? (yes/no) (default: yes): ")
		var confirmation string
		_, _ = fmt.Scanln(&confirmation)
		confirmation = strings.TrimSpace(confirmation)
		confirmation = strings.ToLower(confirmation)
		if len(confirmation) == 0 || confirmation == "yes" {
			gitter.Collect(category, filename)
			log.Println("Article uploaded. Bravo!")
		} else {
			for true {
				if confirmation == "no" {
					break
				}
				fmt.Print("Are you ready to upload the article? (yes/no) (default: yes): ")
				_, _ = fmt.Scanln(&confirmation)
			}
		}
	}
	log.Println("Mission completed. Adios!")
	return nil
}
