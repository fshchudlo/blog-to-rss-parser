package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"blog-to-rss-parser/models"
)

func FetchWebsiteContent(url string) (*goquery.Document, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch website content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse website content: %w", err)
	}
	return doc, nil
}

func ParseArticles(baseURL string, doc *goquery.Document, locator string) []models.RSSItem {
	var items []models.RSSItem

	doc.Find(locator).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h1,h2,h3").Text())
		articleLink, _ := s.Find("a").Attr("href")
		description := strings.TrimSpace(s.Find("p").Text())
		coverImageLink, _ := s.Find("img").Attr("src")
		pubDate := time.Now()

		if timeString, exists := s.Find("time").Attr("content"); exists {
			if parsedTime, err := time.Parse(time.RFC3339Nano, timeString); err == nil {
				pubDate = parsedTime
			}
		}

		articleLink, err := ResolveRelativeUrl(baseURL, articleLink)
		if err != nil {
			log.Printf("Warning: failed to resolve URL %s: %v", articleLink, err)
			return
		}

		item := models.RSSItem{
			Title:       title,
			Link:        articleLink,
			Description: description,
			PubDate:     pubDate.Format(time.RFC822),
		}

		coverImageLink, _ = ResolveRelativeUrl(baseURL, coverImageLink)
		if coverImageLink != "" {
			item.Media = models.MediaContent{
				URL:    coverImageLink,
				Medium: "image",
			}

		}
		items = append(items, item)
	})

	return items
}
