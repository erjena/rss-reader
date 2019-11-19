package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"

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

// func insertItems(db *sql.DB) {
// 	var rss = getRss("abc@gmail.com")
// 	insertStatment, err := db.Prepare("INSERT INTO items (source_id, title, link, description, pubDate) VALUES (?, ?, ?, ?, ?)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer insertStatment.Close()

// 	_, err = insertStatment.Exec(2, rss.Channel.Items[0].Title, rss.Channel.Items[0].Link, rss.Channel.Items[0].Description, rss.Channel.Items[0].PubDate)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
