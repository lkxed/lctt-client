package feeder

import (
	"lctt-client/internal/configurar"
	"lctt-client/internal/helper"
	"time"
)

type Item struct {
	Title      string
	Link       string
	PubDate    time.Time
	Categories []string
	Summary    string
}

var (
	Links []string
)

func init() {
	for hostname, website := range configurar.Websites {
		feed := website.Feed
		if len(feed) == 0 {
			continue
		}
		baseUrl := "https://" + hostname
		link := helper.ConcatUrl(baseUrl, feed)
		Links = append(Links, link)
	}
}
