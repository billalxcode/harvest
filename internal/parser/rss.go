package parser

import (
	"encoding/xml"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}
type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func Decode(body []byte) (RSS, error) {
	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return RSS{}, err
	}

	return rss, nil
}
