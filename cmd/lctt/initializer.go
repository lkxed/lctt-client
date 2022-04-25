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
		Description: `To get helps with supported commands, type "lctt <command> --help/-h".
Configuration files are in "configs" folder, and temporary files are stored in "tmp" folder.
Note: you should never `,
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
				Name:      "init",
				Usage:     "Initializes the client.",
				UsageText: "lctt init",
				Description: `The "init" command does the following things:
1. Forks the upstream if not forked.
2. Clones the origin if not cloned.
3. Pulls the origin if local exists.`,
				Action: initialize,
			},
			{
				Name:  "feed",
				Usage: "Feeds you a list of articles published recently.",
				Description: `The "feed" command provides you with articles recently published on all subscribed websites.
You can subscribe a website by adding it to "configs/websites.yml".`,
				Flags: []cli.Flag{
					&cli.TimestampFlag{
						Name:        "since",
						Aliases:     []string{"s"},
						Layout:      layout,
						Usage:       "Specifies the `<DATE>` after which articles are published.",
						DefaultText: time.Now().Format(layout),
					},
					//&cli.StringFlag{
					//	Name:        "prefer",
					//	Aliases:     []string{"p"},
					//	Usage:       "Specifies a single `<CATEGORY>` you prefer.",
					//	DefaultText: "ALL",
					//},
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Prints full details of articles.",
					},
					&cli.BoolFlag{
						Name:    "open",
						Aliases: []string{"o"},
						Usage:   `Opens articles in your browser (specified in "configs/settings.yaml").`,
					},
				},
				Action: feed,
			},
			{
				Name:                   "collect",
				UseShortOptionHandling: true,
				Usage:                  "Collects an article with its <CATEGORY> and <LINK>",
				Description: `The "collect" command needs you to specify the <CATEGORY> and <LINK> of an article,
where <CATEGORY> can be "news", "talk" and "tech" and <LINK> belongs to a website defined in "configs/websites.yml".`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Usage:   "Specifies the `<CATEGORY>` of the article.",
					},
					&cli.BoolFlag{
						Name:    "preview",
						Aliases: []string{"p"},
						Usage:   `Auto-preview the article in your editor (specified in "configs/settings.yaml").`,
					},
					&cli.BoolFlag{
						Name:    "upload",
						Aliases: []string{"u"},
						Usage:   "Auto-upload the article.",
					},
				},
				Action: collect,
			},
			{
				Name:  "list",
				Usage: "Lists all articles belongs to the given <CATEGORY>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Usage:   "Specifies the `<CATEGORY>` of articles you are interested in.",
					},
				},
				Action: list,
			},
			{
				Name:                   "request",
				UseShortOptionHandling: true,
				Usage:                  "Requests to translate an article with its <CATEGORY> and <FILENAME>",
				Description: `The "request" command needs you to specify the <CATEGORY> and <FILENAME> of an article,
where <CATEGORY> can be "news", "talk" and "tech".`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "category",
						Aliases:  []string{"c"},
						Required: true,
						Usage:    "Specifies the `<CATEGORY>` of the article.",
					},
					&cli.BoolFlag{
						Name:    "open",
						Aliases: []string{"o"},
						Usage:   "Specifies the `<CATEGORY>` of the article.",
					},
				},
				Action: request,
			},
			{
				Name:  "complete",
				Usage: "Completes the translating process of an article with its <CATEGORY> and <FILENAME>",
				Description: `The "complete" command checks translation, and uploads your translation if the checks passed.
You need to specify the <CATEGORY> and <FILENAME> of an article,
where <CATEGORY> can be "news", "talk" and "tech".`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Usage:   "Specifies the `<CATEGORY>` of the article.",
					},
					&cli.BoolFlag{
						Name:    "modify",
						Aliases: []string{"m"},
						Usage:   "Modifies the translation.",
					},
					&cli.StringFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Ignores the checks, upload anyway.",
					},
				},
				Action: complete,
			},
			{
				Name:      "clean",
				Usage:     "Cleans the client.",
				UsageText: "lctt clean",
				Description: `The "clean" command does the following things (if their relating PRs are merged):
1. Deletes local branches.
2. Deletes origin branches.
3. Deletes temporary files.`,
				Action: clean,
			},
		},
	}
}
