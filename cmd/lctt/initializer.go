package main

import (
	"github.com/urfave/cli/v2"
	"time"
)

var (
	app    *cli.App
	layout string
)

func init() {
	layout = "2006-01-02"
	app = &cli.App{
		Name:      "LCTT Client",
		HelpName:  "lctt",
		Usage:     "Aims to be THE All-In-One client for LCTT (https://linux.cn/lctt/)",
		UsageText: "lctt <command> [<options>] [<arguments>]",
		Version:   "0.0.1",
		Authors: []*cli.Author{
			{
				Name:  "lkxed",
				Email: "lkxed@outlook.com",
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initializes the client.",
				Description: "The `init` command does the following things:\n" +
					"1. Forks the upstream if not forked.\n" +
					"2. Clones the origin if not cloned.\n" +
					"3. Pulls the origin if local exists.\n",
				Action: initialize,
			},
			{
				Name:  "feed",
				Usage: "Feeds you a list of articles published recently.",
				Description: "The `feed` command provides you with articles recently published on all subscribed websites.\n" +
					"You can subscribe a website by adding it to `configs/websites.yml`.",
				Flags: []cli.Flag{
					&cli.TimestampFlag{
						Name:        "since",
						Layout:      layout,
						Usage:       "Specifies the `<DATE>` after which articles are published.",
						DefaultText: time.Now().Format(layout),
					},
					&cli.StringFlag{
						Name:        "prefer",
						Usage:       "Specifies a single `<CATEGORY>` you prefer.",
						DefaultText: "ALL",
					},
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Prints full details of articles.",
					},
					&cli.BoolFlag{
						Name:  "open",
						Usage: "Opens articles in your browser (specified in `configs/settings.yaml`).",
					},
				},
				Action: feed,
			},
			{
				Name:      "collect",
				Usage:     "Collects an article with its <CATEGORY> and <LINK>",
				UsageText: "lctt collect [<options>] <CATEGORY> <LINK>",
				Description: "The `collect` command needs you to specify the <CATEGORY> and <LINK> of an article,\n" +
					"where <CATEGORY> can be `news`, `talk` and `tech` and <LINK> belongs to a website defined in `configs/websites.yml`.",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-preview",
						Usage: "Disables auto-preview the article in your editor (specified in `configs/settings.yaml`).",
					},
					&cli.BoolFlag{
						Name:  "no-upload",
						Usage: "Disables auto-upload the article.",
					},
				},
				Action: collect,
			},
		},
	}

}
