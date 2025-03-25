package main

import "encoding/xml"

type RSSDocument struct {
	XMLName    xml.Name   `xml:"rss"`
	Version    string     `xml:"version,attr"`
	Channel    RSSChannel `xml:"channel"`
	XMLNSMedia string     `xml:"xmlns:media,attr,omitempty"`
}
