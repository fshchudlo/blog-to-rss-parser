package main

import (
	"log"

	"blog-to-rss-parser/utils"
)

func main() {
	const feedFileName = "feed.xml"
	websites := map[string]string{
		"https://www.anthropic.com/engineering": "main article article",
		"https://blog.bitdrift.dev":             "article",
	}

	for url, locator := range websites {
		log.Printf("Processing %s", url)

		xmlFeed, err := utils.ReadExistingFeedFile(feedFileName)
		if err != nil {
			log.Fatalf("Error reading existing feed: %v", err)
		}

		content, err := utils.FetchWebsiteContent(url)
		if err != nil {
			log.Printf("Error fetching content from %s: %v", url, err)
			continue
		}

		parsedItems := utils.ParseArticles(url, content, locator)
		xmlFeed.Channel.Items = utils.MergeRSSItems(xmlFeed.Channel.Items, parsedItems)

		if err := utils.SaveRSSFeedFile(feedFileName, xmlFeed); err != nil {
			log.Printf("Error saving RSS feed: %v", err)
		}
	}
}
