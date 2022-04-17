package helper

import (
	"log"
	"net/http"
)

// Scrape a webpage of the given link
func Scrape(link string) *http.Response {
	// request the HTML page.
	res, err := http.Get(link)
	ExitIfError(err)

	if res.StatusCode != 200 {
		log.Fatalf("%d: %s\n", res.StatusCode, res.Status)
	}
	return res
}
