package collector

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"lctt-client/internal/configurar"
	"lctt-client/internal/helper"
	"log"
	"strconv"
	"strings"
	"time"
)

// Parse the article of the given link
func Parse(link string) Article {
	log.Println("Loading webpage...")
	// scrape the webpage
	res := helper.Scrape(link)
	defer func(Body io.ReadCloser) {
		helper.ExitIfError(Body.Close())
	}(res.Body)

	// load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	helper.ExitIfError(err)

	log.Println("Parsing article...")

	// Parse the base url
	temp := strings.Split(link, "//")
	protocol := temp[0]
	hostname := strings.Split(temp[1], "/")[0]
	website, exists := configurar.Websites[hostname]
	if !exists {
		log.Fatalf("Website: %s is not supported.\n", hostname)
	}
	baseUrl := protocol + "//" + hostname

	// Parse the title
	title := doc.Find(website.Title).First().Text()
	title = helper.ClearSpace(title)

	// Parse the summary
	// css pseudo-class :first-child is not supported, hence the workaround
	summarySelector := website.Summary
	summarySelector = strings.ReplaceAll(summarySelector, ":first-child", "")
	summary := doc.Find(summarySelector).First().Text()
	// note that summary is treated as plain text (without styles)
	summary = helper.ClearSpace(summary)

	// Parse the author
	author := parseAuthor(doc, website.Author, baseUrl)

	// Parse the date
	date := parseDate(doc, website.Date)

	// Parse the texts
	texts, urls := parseTexts(doc, website.Content, website.Exclusion, baseUrl)

	return Article{link, title, summary, author, date, texts, urls}
}

func parseAuthor(doc *goquery.Document, selector string, baseUrl string) Author {
	var author Author
	authorAnchor := doc.Find(selector).First()
	if authorAnchor.Size() == 0 { // means the author is manually specified
		authorInfo := strings.Split(selector, ",")
		authorName := authorInfo[0]
		authorLink := helper.ClearSpace(authorInfo[1])
		author = Author{Name: authorName, Link: authorLink}
	} else {
		authorLink := authorAnchor.AttrOr("href", "")
		authorLink = helper.ConcatUrl(baseUrl, authorLink)
		authorName := authorAnchor.Text()
		authorName = helper.ClearSpace(authorName)
		author = Author{Name: authorName, Link: authorLink}
	}
	return author
}

func parseDate(doc *goquery.Document, selector string) string {
	var date string
	var datetime string
	dateNode := doc.Find(selector).First()
	if dateNode.Is("time") {
		datetime = dateNode.AttrOr("datetime", "")
	} else if dateNode.Is("span") {
		datetime = dateNode.AttrOr("content", "")
	}
	// if provided with the right datetime format
	if len(datetime) > 0 && strings.Contains(datetime, "-") {
		date = strings.Split(datetime, "T")[0]
		date = strings.ReplaceAll(date, "-", "")
	}
	// final approach: try the date layout below
	if len(date) == 0 && dateNode.Size() > 0 {
		dateText := dateNode.First().Text()
		dateText = helper.ClearSpace(dateText)
		parsedTime, parseError := time.Parse("January 2, 2006", dateText)
		if parseError != nil { // shouldn't leave date empty, should we?
			parsedTime = time.Now()
		}
		date = parsedTime.Format("20060102")
	}
	return date
}

func parseTexts(doc *goquery.Document, selector string, exclusion string, baseUrl string) ([]string, []string) {
	var texts []string
	var urls []string
	var urlNo int
	tags := "h2, h3, h4, p, span, amp-img, img, video, iframe, ul, ol, code, pre"
	doc.Find(selector).
		Find(tags).
		Not(exclusion).
		// TODO process <code> inside <li>, I forgot the problem...will deal with it when it occurs again
		Each(func(elementIndex int, s *goquery.Selection) {
			if s.Is("h2") { // process <h2> tags
				texts = append(texts, "### "+s.Text())
			} else if s.Is("h3") { // process <h3> tags
				texts = append(texts, "#### "+s.Text())
			} else if s.Is("h4") {
				texts = append(texts, "##### "+s.Text()) // process <h4> tags, there shouldn't be <h5> or smaller ones
			} else if s.Is("p") { // process <p> tags
				text := helper.ClearSpace(s.Text())
				if len(text) > 0 {
					// process inner <a> tags
					s.Find("a").Each(func(_ int, as *goquery.Selection) {
						urlNo++
						url := as.AttrOr("href", "")
						if !strings.HasPrefix(url, "#") { // ignore in-page anchors
							urls = append(urls, url)
							a := helper.ClearSpace(as.Text())
							substitute := "[" + a + "][" + strconv.Itoa(urlNo) + "]"
							text = strings.ReplaceAll(text, a, substitute)
						}
					})
					s.ChildrenFiltered("code, .code, .inline-code").Each(func(_ int, cs *goquery.Selection) {
						code := cs.Text()
						text = strings.ReplaceAll(text, code, "`"+code+"`")
					})
					s.ChildrenFiltered("strong").Each(func(_ int, ss *goquery.Selection) {
						strong := ss.Text()
						text = strings.ReplaceAll(text, strong, "**"+strong+"**")
					})
					s.ChildrenFiltered("em").Each(func(_ int, es *goquery.Selection) {
						em := es.Text()
						text = strings.ReplaceAll(text, em, "*"+em+"*")
					})
					hasBlockQuoteParents := s.ParentsFiltered("blockquote").Size() > 0
					if hasBlockQuoteParents {
						text = "> " + text
					}
					texts = append(texts, text)
				}
			} else if s.Is("span") {
				otherTags := strings.ReplaceAll(tags, "span, ", "")
				hasOtherTagsParents := s.ParentsFiltered(otherTags).Size() > 0
				if !hasOtherTagsParents {
					text := helper.ClearSpace(s.Text())
					texts = append(texts, text)
				}
			} else if s.Is("amp-img, img") { // process <amp-img> & <img> tags
				hasAmpImgParents := s.ParentsFiltered("amp-img").Size() > 0
				if s.Is("amp-img") || (s.Is("img") && !hasAmpImgParents) {
					urlNo++
					url := s.AttrOr("src", "")
					// if src empty, try data-lazy-src attribute
					if len(url) == 0 || strings.HasPrefix(url, "data:image") {
						url = s.AttrOr("data-lazy-src", "")
					}
					url = helper.ConcatUrl(baseUrl, url)
					urls = append(urls, url)
					var title string
					// if in <figure>, try and get <figcaption>
					if s.Parent().Is("figure") {
						cs := s.Parent().Find("figcaption").First()
						if cs.Size() > 0 {
							title = helper.TrimSpace(cs.Text())
						}
					}
					// if no <figcaption>, use title instead
					if len(title) == 0 {
						title = s.AttrOr("title", "")
					}
					// if missing title, use alt instead
					if len(title) == 0 {
						title = s.AttrOr("alt", "")
					}
					imgText := "![" + title + "][" + strconv.Itoa(urlNo) + "]"
					texts = append(texts, imgText)
				}
			} else if s.Is("video") { // process HTML5 <video> tags
				urlNo++
				url := s.AttrOr("src", "")
				if len(url) == 0 {
					url = s.Find("source").First().AttrOr("src", "")
				}
				urls = append(urls, url)
				texts = append(texts, "![]["+strconv.Itoa(urlNo)+"]")
			} else if s.Is("iframe") { // process YouTube Embed Videos <iframe> tags
				urlNo++
				url := s.AttrOr("src", "")
				url = strings.Split(url, "?")[0]
				url = strings.ReplaceAll(url, "www.youtube.com/embed", "youtu.be")
				urls = append(urls, url)
				texts = append(texts, "![A Video from YouTube]["+strconv.Itoa(urlNo)+"]")
			} else if s.Is("ul") { // process <ul> tags
				var items []string
				s.Find("li").Each(func(_ int, lis *goquery.Selection) {
					liText := helper.ClearSpace(lis.Text())
					if len(liText) > 0 {
						items = append(items, "* "+liText)
					}
				})
				if len(items) > 0 {
					text := strings.Join(items, "\n")
					texts = append(texts, text)
				}
			} else if s.Is("ol") { // process <ol> tags
				var items []string
				itemNo := 0
				s.Find("li").Each(func(_ int, lis *goquery.Selection) {
					liText := helper.ClearSpace(lis.Text())
					if len(liText) > 0 {
						itemNo++
						items = append(items, strconv.Itoa(itemNo)+". "+liText)
					}
				})
				if len(items) > 0 {
					text := strings.Join(items, "\n")
					texts = append(texts, text)
				}
			} else if s.Is("pre, code") { // process <pre> & <code> tags
				hasCodeDescendants := s.Find("code").Size() > 0
				hasPParents := s.ParentsFiltered("p").Size() > 0
				if (s.Is("pre") && !hasCodeDescendants) || (s.Is("code") && !hasPParents) {
					code := helper.TrimSpace(s.Text())
					if len(code) > 0 {
						text := "```\n" + code + "\n```"
						texts = append(texts, text)
					}
				}
			}
		})
	return texts, urls
}
