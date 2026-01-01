package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"blog-to-rss-parser/models"
)

func ReadExistingFeedFile(filename string) (*models.RSSDocument, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &models.RSSDocument{
			Version:    "2.0",
			XMLNSMedia: "http://search.yahoo.com/mrss/",
			Channel: models.RSSChannel{
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

	var feed models.RSSDocument
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed file: %w", err)
	}

	if feed.XMLNSMedia == "" {
		feed.XMLNSMedia = "http://search.yahoo.com/mrss/"
	}

	return &feed, nil
}

func SaveRSSFeedFile(filename string, feed *models.RSSDocument) error {
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

func MergeRSSItems(existingItems, newItems []models.RSSItem) []models.RSSItem {
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
