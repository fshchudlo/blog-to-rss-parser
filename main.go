package main

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func fetchWebsiteContent(url string) (*goquery.Document, error) {
	// Send HTTP request to the website
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func parseArticles(baseUrl string, doc *goquery.Document, locator string) []RSSItem {
	var newItems []RSSItem

	doc.Find(locator).Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2").Text()
		link, _ := s.Find("a").Attr("href")
		description := s.Find("p").Text()
		pubDate := time.Now()

		if timeString, exists := s.Find("time").Attr("content"); exists {
			pubDate, _ = time.Parse(time.RFC3339Nano, timeString)
		}

		newItem := RSSItem{
			Title:       strings.TrimSpace(title),
			Link:        baseUrl + link,
			Description: strings.TrimSpace(description),
			PubDate:     pubDate.Format(time.RFC822),
		}
		newItems = append(newItems, newItem)
	})
	return newItems
}

func readExistingFeed(filename string) (*RSSDocument, error) {
	// If feed file doesn't exist, return empty RSSDocument structure
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &RSSDocument{
			Version: "2.0",
			Channel: RSSChannel{
				Title:       "Web Scraper Feed",
				Description: "Automatically generated RSSDocument feed",
			},
		}, nil
	}

	// Read existing feed file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var feed RSSDocument
	err = xml.Unmarshal(data, &feed)
	return &feed, err
}

func saveRSSFeed(filename string, feed *RSSDocument) error {
	// Marshal RSS feed to XML
	xmlData, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	// Prepend XML header
	xmlContent := []byte(xml.Header)
	xmlContent = append(xmlContent, xmlData...)

	// Write to file
	return os.WriteFile(filename, xmlContent, 0644)
}

func main() {
	feedFileName := "feed.xml"
	websites := map[string]string{
		"https://blog.bitdrift.dev": "article",
	}
	xmlFeed, err := readExistingFeed(feedFileName)

	if err != nil {
		log.Fatal(fmt.Errorf("error reading existing feed: %v", err))
	}

	for url, locator := range websites {
		content, err := fetchWebsiteContent(url)
		if err != nil {
			log.Fatal(err)
		}
		parsedItems := parseArticles(url, content, locator)
		xmlFeed.Channel.Items = append(xmlFeed.Channel.Items, parsedItems...)
	}
	err = saveRSSFeed(feedFileName, xmlFeed)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading existing feed: %v", err))
	}
}
