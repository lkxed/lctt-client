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
	log.Println("Initializing...")

	gitter.Initialize()

	log.Println("Initialized.")
	log.Println("Mission Complete. Adios!")
	return nil
}

func feed(c *cli.Context) error {
	log.Println("Feeding...")

	items := feeder.ParseAll()
	datePtr := c.Timestamp("since")
	var date time.Time
	yesterday := time.Now().AddDate(0, 0, -1).Format(layout)
	if datePtr == nil {
		date, _ = time.Parse(layout, yesterday)
	} else {
		date = *datePtr
	}
	items = feeder.FilterPubDate(items, date)
	//category := c.String("prefer")
	//if len(category) > 0 {
	//	items = feeder.FilterCategories(items, category)
	//}
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
	log.Println("Anyway, Mission Complete. Adios!")
	return nil
}

func collect(c *cli.Context) error {
	if c.NArg() != 1 {
		log.Fatalln(c.Command.Usage)
	}
	log.Println("Collecting...")

	link := c.Args().Get(0)
	article := collector.Parse(link)
	filename, _ := collector.Generate(article)
	tmpPath := path.Join(helper.TmpDir, filename)
	if c.Bool("preview") {
		editor := configurar.Settings.Editor
		if len(configurar.Settings.Editor) == 0 {
			log.Fatalln("No editor specified in `settings.yml`")
		}
		editCmd := strings.Split(editor, " ")
		editCmd = append(editCmd, tmpPath)
		log.Printf("Previewing article in %s...\n", editCmd[0])
		cmd := exec.Command(editCmd[0], editCmd[1:]...)
		helper.ExitIfError(cmd.Run())
		log.Println("Modifications have been saved.")
	}
	var category string
	if c.Bool("upload") {
		category = c.String("category")
		categories := []string{"news", "talk", "tech"}
		if !helper.StringSliceContains(categories, category) {
			log.Fatalln("To upload, you must specify the <CATEGORY>.")
		}
		confirmation := "yes"
		for true {
			fmt.Print("Are you ready to upload the article? (yes/no) (default: yes): ")
			_, _ = fmt.Scanln(&confirmation)
			confirmation = strings.TrimSpace(confirmation)
			confirmation = strings.ToLower(confirmation)
			if confirmation == "no" || confirmation == "yes" {
				break
			}
		}
		if confirmation == "yes" {
			gitter.Collect(category, filename)
		}
	}

	log.Println("Mission Complete. Adios!")
	return nil
}

func request(c *cli.Context) error {
	if c.NArg() != 1 {
		log.Fatalln(c.Command.Usage)
	}
	log.Println("Requesting...")

	filename := c.Args().Get(0)
	category := c.String("category")
	categories := []string{"news", "talk", "tech"}
	if !helper.StringSliceContains(categories, category) {
		log.Fatalln("To upload, you must specify the <CATEGORY>.")
	}

	gitter.Request(category, filename)

	if c.Bool("open") {
		editor := configurar.Settings.Editor
		if len(configurar.Settings.Editor) == 0 {
			log.Fatalln("No editor specified in `settings.yml`")
		}
		editCmd := strings.Split(editor, " ")
		tmpPath := path.Join(helper.TmpDir, filename)
		editCmd = append(editCmd, tmpPath)
		log.Printf("Opening article in %s...\n", editCmd[0])
		cmd := exec.Command(editCmd[0], editCmd[1:]...)
		helper.ExitIfError(cmd.Run())
	}

	log.Println("Mission Complete. Adios!")
	return nil
}

func list(c *cli.Context) error {
	log.Println("Listing...")

	category := c.String("category")
	categories := []string{"news", "talk", "tech"}
	if category != "" && !helper.StringSliceContains(categories, category) {
		log.Fatalln("<CATEGORY> must be `news`, `talk` or `tech`.")
	}
	filenames := gitter.List(category)
	if len(filenames) == 0 {
		return nil
	}
	fmt.Println()
	fmt.Println("[ARTICLE LIST]")
	fmt.Println()
	for _, filename := range filenames {
		if filename == "README.md" {
			continue
		}
		fmt.Printf("- [%s]: %s\n", category, filename)
	}
	fmt.Println()
	fmt.Println("[END]")
	fmt.Println()
	log.Println("Have you made up your mind? If so, choose an article to `collect`.")
	log.Println("Anyway, Mission Complete. Adios!")
	return nil
}

func complete(c *cli.Context) error {
	if c.NArg() != 1 {
		log.Fatalln(c.Command.Usage)
	}
	log.Println("Completing...")

	filename := c.Args().Get(0)

	if c.Bool("modify") {
		editor := configurar.Settings.Editor
		if len(configurar.Settings.Editor) == 0 {
			log.Fatalln("No editor specified in `settings.yml`")
		}
		editCmd := strings.Split(editor, " ")
		tmpPath := path.Join(helper.TmpDir, filename)
		editCmd = append(editCmd, tmpPath)
		log.Printf("Opening article in %s...\n", editCmd[0])
		cmd := exec.Command(editCmd[0], editCmd[1:]...)
		helper.ExitIfError(cmd.Run())

		confirmation := "yes"
		for true {
			fmt.Print("Are you ready to upload the article? (yes/no) (default: yes): ")
			_, _ = fmt.Scanln(&confirmation)
			confirmation = strings.TrimSpace(confirmation)
			confirmation = strings.ToLower(confirmation)
			if confirmation == "no" || confirmation == "yes" {
				break
			}
		}
	}

	category := c.String("category")
	categories := []string{"news", "talk", "tech"}
	if !helper.StringSliceContains(categories, category) {
		log.Fatalln("To upload, you must specify the <CATEGORY>.")
	}

	force := c.Bool("force")
	err := gitter.Complete(category, filename, force)
	if err != nil {
		log.Fatalln("Your translation is not complete. Please complete it and try again.")
	}

	log.Println("Mission Complete. Adios!")
	return nil
}

func clean(_ *cli.Context) error {
	log.Println("Cleaning...")

	err := gitter.Clean()
	if err != nil {
		log.Println("Such a tidy workspace! Nothing to clean.")
	} else {
		log.Println("Cleaned.")
	}

	log.Println("Mission Complete. Adios!")
	return nil
}
