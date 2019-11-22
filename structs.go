package main

import "time"

// Rss - main struct
type Rss struct {
	Channel Channel `xml:"channel" json:"channel"`
}

// Channel struct for RSS
type Channel struct {
	Title string `xml:"title" json:"title"`
	Link  string `xml:"link" json:"link"`
	Items []Item `xml:"item" json:"items"`
}

// Item struct for each item in Channel
type Item struct {
	Title       string `xml:"title" json:"title"`
	Link        string `xml:"link" json:"link"`
	Description string `xml:"description" json:"description"`
	PubDate     string `xml:"pubDate" json:"pubDate"`
	SourceName  string `json:"sourceName"`
}

// Source stores source info
type Source struct {
	ID          int
	Link        string
	LastPubDate *time.Time
}

// ResponseToClient to send response with source name and array of Items
type ResponseToClient struct {
	SourceID string `json:"sourceID"`
	Items    []Item `json:"items"`
}
