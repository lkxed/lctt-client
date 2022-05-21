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
	title = strings.TrimSpace(title)

	// Parse the summary
	// css pseudo-class :first-child is not supported, hence the workaround
	summarySelector := website.Summary
	summarySelector = strings.ReplaceAll(summarySelector, ":first-child", "")
	summary := doc.Find(summarySelector).First().Text()
	// note that summary is treated as plain text (without styles)
	summary = strings.TrimSpace(summary)

	// Parse the author
	author := parseAuthor(doc, website.Author, baseUrl)

	// Parse the date
	date := parseDate(doc, website.Date)

	// Parse the texts
	texts, urls := parseTexts(doc, website.Content, website.Exclusion, baseUrl)

	return Article{link, title, summary, author, date, texts, urls}
}

// parseAuthor parses author's name and link.
func parseAuthor(doc *goquery.Document, selector string, baseUrl string) Author {
	var author Author
	authorAnchor := doc.Find(selector).First()
	if authorAnchor.Size() == 0 { // means the author is manually specified
		authorInfo := strings.Split(selector, ",")
		authorName := authorInfo[0]
		authorLink := strings.TrimSpace(authorInfo[1])
		author = Author{Name: authorName, Link: authorLink}
	} else {
		authorLink := authorAnchor.AttrOr("href", "")
		authorLink = helper.ConcatUrl(baseUrl, authorLink)
		authorName := authorAnchor.Text()
		authorName = strings.TrimSpace(authorName)
		author = Author{Name: authorName, Link: authorLink}
	}
	return author
}

// parseDate parses the publishing date.
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
		dateText = strings.TrimSpace(dateText)
		parsedTime, parseError := time.Parse("January 2, 2006", dateText)
		if parseError != nil { // shouldn't leave date empty, should we?
			parsedTime = time.Now()
		}
		date = parsedTime.Format("20060102")
	}
	return date
}

// parseRichContent parses rich contents inside a container, such as <p>.
// TODO parse the container recursively with all of its descendents.
func parseRichContent(s *goquery.Selection, baseUrl string, urlNo int, urls []string) (string, int, []string) {
	var text string
	s.Contents().Each(func(_ int, ps *goquery.Selection) {
		if ps.Is("a") && ps.Find("img").Size() == 0 {
			url := ps.AttrOr("href", "")
			if url != "" && !strings.HasPrefix(url, "#") { // ignore in-page anchors
				urlNo++
				url = helper.ConcatUrl(baseUrl, url)
				urls = append(urls, url)
				a := strings.TrimSpace(ps.Text())
				a = "[" + a + "][" + strconv.Itoa(urlNo) + "]"
				text += a
			}
		} else if ps.Is("code, .code, .inline-code") {
			code := strings.TrimSpace(ps.Text())
			if len(code) > 0 {
				code = "`" + code + "`"
				text += code
			}
		} else if ps.Is("strong") {
			strong := strings.TrimSpace(ps.Text())
			if len(strong) > 0 {
				aChildren := ps.ChildrenFiltered("a")
				if aChildren.Size() > 0 {
					url := aChildren.First().AttrOr("href", "")
					if url != "" {
						urlNo++
						url = helper.ConcatUrl(baseUrl, url)
						urls = append(urls, url)
						strong = "[" + strong + "][" + strconv.Itoa(urlNo) + "]"
					}
				}
				strong = "**" + strong + "**"
				text += strong
			}
		} else if ps.Is("em") {
			em := strings.TrimSpace(ps.Text())
			if len(em) > 0 {
				aChildren := ps.ChildrenFiltered("a")
				if aChildren.Size() > 0 {
					url := aChildren.First().AttrOr("href", "")
					if url != "" {
						urlNo++
						url = helper.ConcatUrl(baseUrl, url)
						urls = append(urls, url)
						em = "[" + em + "][" + strconv.Itoa(urlNo) + "]"
					}
				}
				em = "*" + em + "*"
				text += em
			}
		} else {
			p := ps.Text()
			m, n := len(text), len(p)
			if m > 0 && n > 0 && text[m-1] == '`' && p[0] != ' ' && p[0] != '.' && p[0] != ',' {
				text += " "
			}
			text += p
		}
	})
	hasBlockQuoteParents := s.ParentsFiltered("blockquote").Size() > 0
	if hasBlockQuoteParents {
		text = "> " + text
	}
	text = strings.TrimSpace(text)
	return text, urlNo, urls
}

// parseTexts parses the main content.
func parseTexts(doc *goquery.Document, selector string, exclusion string, baseUrl string) ([]string, []string) {
	var texts []string
	var urls []string
	var urlNo int
	tags := "h2, h3, h4, p, a, span, amp-img, img, amp-video, video, iframe, ul, ol, code, pre, td"
	selection := doc.Find(selector).Find(tags).Not(exclusion)
	// TODO process <code> inside <li>, I forgot the problem...will deal with it when it occurs again
	selection.Each(func(elementIndex int, s *goquery.Selection) {
		if s.Is("h2") { // process <h2> tags
			texts = append(texts, "### "+s.Text())
		} else if s.Is("h3") { // process <h3> tags
			texts = append(texts, "#### "+s.Text())
		} else if s.Is("h4") {
			texts = append(texts, "##### "+s.Text()) // process <h4> tags, there shouldn't be <h5> or smaller ones
		} else if s.Is("p") { // process <p> tags
			hasLiParents := s.ParentsFiltered("li").Size() > 0
			if !hasLiParents {
				var text string
				text, urlNo, urls = parseRichContent(s, baseUrl, urlNo, urls)
				texts = append(texts, text)
			}
		} else if s.Is("a") {
			otherTags := strings.Replace(tags, "a, ", "", 1)
			hasOtherTagsParents := s.ParentsFiltered(otherTags).Size() > 0
			url := s.AttrOr("href", "")
			url = strings.TrimSpace(url)
			isInPageAnchor := strings.HasPrefix(url, "#")
			hasImgChildren := s.Find("img").Size() > 0
			if url != "" && !hasOtherTagsParents && !isInPageAnchor && !hasImgChildren {
				urlNo++
				url = helper.ConcatUrl(baseUrl, url)
				urls = append(urls, url)
				aText := strings.TrimSpace(s.Text())
				text := "[" + aText + "][" + strconv.Itoa(urlNo) + "]"
				texts = append(texts, text)
			}
		} else if s.Is("span") {
			otherTags := strings.Replace(tags, "span, ", "", 1)
			hasOtherTagsParents := s.ParentsFiltered(otherTags).Size() > 0
			if !hasOtherTagsParents {
				text := strings.TrimSpace(s.Text())
				texts = append(texts, text)
			}
		} else if s.Is("amp-img, img") { // process <amp-img> & <img> tags
			hasAmpImgParents := s.ParentsFiltered("amp-img").Size() > 0
			hidden := s.AttrOr("aria-hidden", "false") == "true"
			if s.Is("amp-img") || (s.Is("img") && !hidden && !hasAmpImgParents) {
				urlNo++
				var url string
				var title string
				// check if it has any <a> parents
				aParents := s.ParentsFiltered("a")
				if aParents.Size() > 0 {
					aParent := aParents.First()
					url = aParent.AttrOr("href", "")
					figCaptionSiblings := aParent.SiblingsFiltered("figcaption")
					if figCaptionSiblings.Size() > 0 {
						title = figCaptionSiblings.First().Text()
					}
				}
				// parse image url
				if url == "" {
					url = s.AttrOr("src", "")
				}
				if strings.HasPrefix(url, "data:image") {
					url = ""
				}
				// if src empty, try data-src attribute
				if url == "" {
					url = s.AttrOr("data-src", "")
				}
				// if data-src empty, try data-lazy-src attribute
				if url == "" {
					url = s.AttrOr("data-lazy-src", "")
				}

				// parse image title
				// if in <figure>, try and get <figcaption>
				if title == "" && s.Parent().Is("figure") {
					cs := s.Parent().Find("c").First()
					if cs.Size() > 0 {
						title = strings.TrimSpace(cs.Text())
					}
				}
				// if no <figcaption>, use title instead
				if title == "" {
					title = s.AttrOr("title", "")
				}
				// if missing title, use alt instead
				if title == "" {
					title = s.AttrOr("alt", "")
				}
				url = helper.ConcatUrl(baseUrl, url)
				urls = append(urls, url)
				title = strings.TrimSpace(title)
				imgText := "![" + title + "][" + strconv.Itoa(urlNo) + "]"
				texts = append(texts, imgText)
			}
		} else if s.Is("amp-video, video") { // process <amp-video> & HTML5 <video> tags
			hasAmpVideoParents := s.ParentsFiltered("amp-video").Size() > 0
			if s.Is("amp-video") || (s.Is("video") && !hasAmpVideoParents) {
				urlNo++
				url := s.AttrOr("src", "")
				source := s.Find("source")
				if len(url) == 0 && source.Size() > 0 {
					url = source.First().AttrOr("src", "")
				}
				url = helper.ConcatUrl(baseUrl, url)
				urls = append(urls, url)
				texts = append(texts, "![]["+strconv.Itoa(urlNo)+"]")
			}
		} else if s.Is("iframe") { // process YouTube Embed Videos <iframe> tags
			url := s.AttrOr("src", "")
			if url == "" {
				url = s.AttrOr("data-src", "")
			}
			if url != "" {
				urlNo++
				url = strings.Split(url, "?")[0]
				url = strings.ReplaceAll(url, "www.youtube.com/embed", "youtu.be")
				urls = append(urls, url)
				title := s.AttrOr("title", "")
				if title == "" {
					title = "A Video from YouTube"
				}
				texts = append(texts, "!["+title+"]["+strconv.Itoa(urlNo)+"]")
			}
		} else if s.Is("ul") || s.Is("ol") { // process <ul> & <ol> tags
			var items []string
			itemNo := 0
			s.Find("li").Each(func(_ int, lis *goquery.Selection) {
				liText := strings.TrimSpace(lis.Text())
				// process <a> tags inside each <li> tag
				lis.Find("a").Each(func(_ int, as *goquery.Selection) {
					url := as.AttrOr("href", "")
					if !strings.HasPrefix(url, "#") { // ignore in-page anchors
						urlNo++
						url = helper.ConcatUrl(baseUrl, url)
						urls = append(urls, url)
						aOld := strings.TrimSpace(as.Text())
						aNew := "[" + aOld + "][" + strconv.Itoa(urlNo) + "]"
						liText = strings.Replace(liText, aOld, aNew, 1)
					}
				})
				lis.Find("code").Each(func(_ int, codes *goquery.Selection) {
					code := strings.TrimSpace(codes.Text())
					if len(code) > 0 {
						codeOld := " " + code
						codeNew := " `" + code + "`"
						liText = strings.Replace(liText, codeOld, codeNew, 1)
					}
				})
				if len(liText) > 0 {
					if s.Is("ol") {
						itemNo++
						liText = strconv.Itoa(itemNo) + ". " + liText
					} else {
						liText = "* " + liText
					}
				}
				items = append(items, liText)
			})
			itemNo = 0
			if len(items) > 0 {
				text := strings.Join(items, "\n")
				texts = append(texts, text)
			}
		} else if s.Is("pre, code") { // process <pre> & <code> tags
			hasCodeDescendants := s.Find("code").Size() > 0
			hasPOrLiParents := s.ParentsFiltered("p, li").Size() > 0
			if (s.Is("pre") && !hasCodeDescendants) || (s.Is("code") && !hasPOrLiParents) {
				code := strings.TrimSpace(s.Text())
				if len(code) > 0 {
					text := "```\n" + code + "\n```"
					texts = append(texts, text)
				}
			}
		} else if s.Is("td") { // TODO transform HTML tables into Markdown tables
			var text string
			s.Contents().Each(func(_ int, tds *goquery.Selection) {
				if goquery.NodeName(tds) == "#text" {
					text += tds.Text()
				} else if tds.Is("strong") {
					text = text + "**" + tds.Text() + "**"
				} else if tds.Is("em") {
					text = text + "*" + tds.Text() + "*"
				}
			})
			check := strings.ReplaceAll(text, "*", "")
			if len(strings.TrimSpace(check)) > 0 {
				text = strings.TrimSpace(text)
				texts = append(texts, text)
			}
		}
	})
	return texts, urls
}
