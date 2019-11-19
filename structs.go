package main

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
}
