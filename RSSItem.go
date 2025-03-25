package main

type RSSItem struct {
	Title        string       `xml:"title"`
	Link         string       `xml:"link"`
	Description  string       `xml:"description"`
	PubDate      string       `xml:"pubDate"`
	MediaContent MediaContent `xml:"media:content"`
}

type MediaContent struct {
	URL    string `xml:"url,attr"`
	Medium string `xml:"medium,attr"`
}
