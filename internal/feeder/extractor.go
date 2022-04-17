package feeder

func ExtractLinks(items []Item) []string {
	var links []string
	for _, item := range items {
		links = append(links, item.Link)
	}
	return links
}
