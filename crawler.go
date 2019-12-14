package main

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func crawlWrapper(db *sql.DB) {
	for true {
		crawl(db)
		time.Sleep(1 * time.Minute)
	}
}

func crawlSingleSource(source *Source, db *sql.DB, waitGroup *sync.WaitGroup) {
	var items = getXML(source.Link, source.LastPubDate)
	insertItems(db, source.ID, items)
	if waitGroup != nil {
		waitGroup.Done()
	}
}

func crawl(db *sql.DB) {
	startTime := time.Now()
	log.Println("Start crawl")
	var sources = getAllSources(db)

	var waitGroup sync.WaitGroup
	for _, source := range sources {
		log.Printf("Fetch data for source %v:\n", source.Link)
		waitGroup.Add(1)
		go crawlSingleSource(source, db, &waitGroup)
	}

	waitGroup.Wait()
	endTime := time.Now()
	log.Printf("End crawl. Elapsed %v:\n", endTime.Sub(startTime))
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
