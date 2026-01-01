package models

import "encoding/xml"

type RSSDocument struct {
	XMLName    xml.Name   `xml:"rss"`
	Version    string     `xml:"version,attr"`
	Channel    RSSChannel `xml:"channel"`
	XMLNSMedia string     `xml:"xmlns:media,attr"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description string       `xml:"description"`
	PubDate     string       `xml:"pubDate"`
	Media       MediaContent `xml:"http://search.yahoo.com/mrss/ content"`
}

type MediaContent struct {
	XMLName xml.Name `xml:"http://search.yahoo.com/mrss/ content"`
	URL     string   `xml:"url,attr"`
	Medium  string   `xml:"medium,attr"`
}
