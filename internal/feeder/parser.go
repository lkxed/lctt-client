package feeder

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/SlyMarbo/rss"
	"html"
	"lctt-client/internal/helper"
	"log"
	"strings"
)

func ParseAll() []Item {
	log.Println("Parsing feed...this may take a while...")

	var items []Item
	for _, link := range Links {
		items = append(items, parse(link)...)
	}

	log.Println("Completed.")
	return items
}

func parse(link string) []Item {
	feed, err := rss.Fetch(link)
	helper.ExitIfError(err)
	var items []Item
	for _, item := range feed.Items {
		title := item.Title
		link := item.Link
		categories := item.Categories
		// Parse the summary
		summary := item.Summary
		if len(summary) == 0 {
			summary = item.Content
		}
		summary = html.UnescapeString(summary)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(summary))
		helper.ExitIfError(err)
		// set summary as the first non-empty <p>'s text
		doc.Find("p").EachWithBreak(func(i int, s *goquery.Selection) bool {
			text := helper.ClearSpace(s.Text())
			if len(text) != 0 {
				summary = text
				return false
			}
			return true
		})
		pubDate := item.Date
		items = append(items, Item{
			Title:      title,
			Link:       link,
			PubDate:    pubDate,
			Categories: categories,
			Summary:    summary,
		})
	}
	log.Println(link)
	return items
}
