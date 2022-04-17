package feeder

import (
	"strings"
	"time"
)

func FilterPubDate(items []Item, since time.Time) []Item {
	var recentItems []Item
	for _, item := range items {
		if item.PubDate.After(since) {
			recentItems = append(recentItems, item)
		}
	}
	return recentItems
}

func FilterCategories(items []Item, prefer string) []Item {
	var interestingItems []Item
	prefer = strings.ToLower(prefer)
	for _, item := range items {
		for _, category := range item.Categories {
			category = strings.ToLower(category)
			if strings.Contains(category, prefer) {
				interestingItems = append(interestingItems, item)
				break
			}
		}
	}
	return interestingItems
}
