package feeder

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	items := parse("https://www.linuxuprising.com/feeds/posts/default")
	log.Println(items)
}
