package main

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func crawlSingleSource(source *Source, db *sql.DB) {
	var items = getXML(source.Link, source.LastPubDate)
	insertItems(db, source.ID, items)
}

func crawl(db *sql.DB) {
	var sources = getAllSources(db)
	for _, source := range sources {
		crawlSingleSource(source, db)
	}
}

func getXML(link string, lastPubDate *time.Time) []Item {
	resp, err := http.Get(link)
	if err != nil {
		log.Printf("Was not able to get %v", link)
		return make([]Item, 0)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status error: %v", err)
		return make([]Item, 0)
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return make([]Item, 0)
	}
	var rss Rss
	err = xml.Unmarshal(byteValue, &rss)
	if err != nil {
		log.Print(err)
		return make([]Item, 0)
	}

	var newData []Item
	for _, item := range rss.Channel.Items {
		itemPubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Print(err)
			continue
		}

		if lastPubDate == nil || itemPubDate.After(*lastPubDate) {
			newData = append(newData, item)
		}
	}

	return newData
}
