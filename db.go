package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getRss(userName string) Rss {
	xmlFile, err := os.Open("testData.xml")
	if err != nil {
		log.Fatalf("Was not able to open file %v", xmlFile)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var rss Rss
	xml.Unmarshal(byteValue, &rss)
	// print as json
	// err = json.NewEncoder(os.Stdout).Encode(rss)
	// if err != nil {
	// 	log.Fatalf("Was not able to convert to json %v", err)
	// }
	return rss
}

func dbConnection() *sql.DB {
	configFile, _ := os.Open("dbconfig.json")
	defer configFile.Close()
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Was nor able to read config file %v", err)
	}

	type Configuration struct {
		User     string
		Password string
	}

	var configuration Configuration
	err = json.Unmarshal(byteValue, &configuration)
	if err != nil {
		log.Fatalf("Was not able to parse %v", err)
	}

	db, err := sql.Open("mysql", configuration.User+":"+configuration.Password+"@/rss_reader")
	if err != nil {
		log.Fatalf("Was not able to connect to db %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertSource(db *sql.DB, username string, source string) {
	rows, err := db.Query("SELECT id FROM users WHERE email=?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() != true {
		log.Fatal("User name was not found")
	}
	var userID int
	rows.Scan(&userID)

	insertStatment, err := db.Prepare("INSERT INTO sources (user_id, link, create_time) VALUES (?, ?, CURRENT_TIMESTAMP())")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStatment.Close()

	_, err = insertStatment.Exec(userID, source)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllSources(db *sql.DB) []*Source {
	rows, err := db.Query("SELECT id, link FROM sources")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var sources []*Source
	for rows.Next() {
		var source Source
		err := rows.Scan(&source.ID, &source.Link)
		if err != nil {
			log.Fatal(err)
		}
		sources = append(sources, &source)
	}

	for _, item := range sources {
		time, err := db.Query("SELECT max(pubDate) FROM items WHERE source_id=?", item.ID)
		if err != nil {
			log.Print(err)
		}
		if time.Next() == false {
			log.Println("Publication date was not found")
		}
		err = time.Scan(&item.LastPubDate)
		if err != nil {
			log.Print(err)
		}
	}

	return sources
}

func insertItems(db *sql.DB, sourceID int, items []Item) {
	insertStatment, err := db.Prepare("INSERT INTO items (source_id, title, link, description, pubDate) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStatment.Close()

	for _, item := range items {
		pubdate, _ := time.Parse(time.RFC1123Z, item.PubDate)
		_, err = insertStatment.Exec(sourceID, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			log.Fatal(err)
		}
	}
}