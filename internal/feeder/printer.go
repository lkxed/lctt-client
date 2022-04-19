package feeder

import (
	"fmt"
	"strings"
)

func List(items []Item) {
	fmt.Println()
	fmt.Println("[ARTICLE LIST]")
	fmt.Println()
	for _, item := range items {
		fmt.Printf("- %s \"%s\"\n", item.PubDate.Format("2006-01-02"), item.Title)
		fmt.Printf("  %s\n", item.Link)
		fmt.Println()
	}
	fmt.Println("[END]")
	fmt.Println()
}

func ListVerbose(items []Item) {
	fmt.Println()
	fmt.Println("[ARTICLE LIST]")
	fmt.Println()
	for _, item := range items {
		fmt.Printf(`- "%s""\n`, item.Title)
		fmt.Printf("  %s\n", item.Link)
		fmt.Printf("  Published at %s\n", item.PubDate.Format("2006-01-02 15:04"))
		fmt.Printf("  [%s]\n", strings.Join(item.Categories, ", "))
		fmt.Printf("  > %s\n", item.Summary)
		fmt.Println()
	}
	fmt.Println("[END]")
	fmt.Println()
}
