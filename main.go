package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetchWebsiteContent(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch website content: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse website content: %w", err)
	}

	return doc, nil
}

func parseArticles(baseURL string, doc *goquery.Document, locator string) []RSSItem {
	var items []RSSItem

	doc.Find(locator).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2").Text())
		link, _ := s.Find("a").Attr("href")
		description := strings.TrimSpace(s.Find("p").Text())
		coverImageLink, _ := s.Find("img").Attr("src")
		pubDate := time.Now()

		if timeString, exists := s.Find("time").Attr("content"); exists {
			if parsedTime, err := time.Parse(time.RFC3339Nano, timeString); err == nil {
				pubDate = parsedTime
			}
		}

		if strings.HasPrefix(link, "/") {
			link = baseURL + link
		}

		if strings.HasPrefix(coverImageLink, "/") {
			coverImageLink = baseURL + coverImageLink
		}

		item := RSSItem{
			Title:       title,
			Link:        link,
			Description: description,
			PubDate:     pubDate.Format(time.RFC822),
			Media: MediaContent{
				URL:    coverImageLink,
				Medium: "image",
			},
		}
		items = append(items, item)
	})

	return items
}

func readExistingFeed(filename string) (*RSSDocument, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &RSSDocument{
			Version:    "2.0",
			XMLNSMedia: "http://search.yahoo.com/mrss/",
			Channel: RSSChannel{
				Title:       "Web Scraper Feed",
				Description: "Automatically generated RSS feed",
				Link:        "https://fshchudlo.github.io/blog-to-rss-parser/feed.xml",
			},
		}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read feed file: %w", err)
	}

	var feed RSSDocument
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed file: %w", err)
	}

	if feed.XMLNSMedia == "" {
		feed.XMLNSMedia = "http://search.yahoo.com/mrss/"
	}

	return &feed, nil
}

func saveRSSFeed(filename string, feed *RSSDocument) error {
	xmlData, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RSS feed: %w", err)
	}

	xmlContent := append([]byte(xml.Header), xmlData...)
	if err := os.WriteFile(filename, xmlContent, 0644); err != nil {
		return fmt.Errorf("failed to write RSS feed to file: %w", err)
	}

	return nil
}

func appendWithoutDuplicates(existingItems, newItems []RSSItem) []RSSItem {
	existingLinks := make(map[string]struct{}, len(existingItems))
	for _, item := range existingItems {
		existingLinks[item.Link] = struct{}{}
	}

	for _, item := range newItems {
		if _, exists := existingLinks[item.Link]; !exists {
			existingItems = append(existingItems, item)
			existingLinks[item.Link] = struct{}{}
		}
	}

	return existingItems
}

func main() {
	const feedFileName = "feed.xml"
	websites := map[string]string{
		"https://blog.bitdrift.dev": "article",
	}

	xmlFeed, err := readExistingFeed(feedFileName)
	if err != nil {
		log.Fatalf("Error reading existing feed: %v", err)
	}

	for url, locator := range websites {
		content, err := fetchWebsiteContent(url)
		if err != nil {
			log.Printf("Error fetching content from %s: %v", url, err)
			continue
		}

		parsedItems := parseArticles(url, content, locator)
		xmlFeed.Channel.Items = appendWithoutDuplicates(xmlFeed.Channel.Items, parsedItems)
	}

	if err := saveRSSFeed(feedFileName, xmlFeed); err != nil {
		log.Fatalf("Error saving RSS feed: %v", err)
	}
}
