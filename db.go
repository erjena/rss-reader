package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"os"
	"strings"
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

	db, err := sql.Open("mysql", configuration.User+":"+configuration.Password+"@/rss_reader?parseTime=true")
	if err != nil {
		log.Fatalf("Was not able to connect to db %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertUserInfo(db *sql.DB, email string, salt string, hash string) int {
	insertStatment, err := db.Prepare("INSERT INTO users (email, password_salt, password_hash, create_time, update_time) VALUES (?, ?, ?, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP())")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStatment.Close()
	insertResult, err := insertStatment.Exec(email, salt, hash)
	if err != nil {
		log.Fatal(err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return int(id)
}

func getUserInfo(db *sql.DB, username string) (int, string, string, error) {
	rows, err := db.Query("SELECT id, password_salt, password_hash FROM users WHERE email=?", username)
	if err != nil {
		log.Print(err)
		return 0, "", "", err
	}
	defer rows.Close()
	if rows.Next() != true {
		log.Printf("User was not found: %v", username)
		return 0, "", "", errors.New("User was not found")
	}
	var id int
	var salt string
	var hash string
	rows.Scan(&id, &salt, &hash)
	return id, salt, hash, nil
}

func insertSource(db *sql.DB, userID int, source string) int {
	insertStatment, err := db.Prepare("INSERT INTO sources (user_id, link, create_time) VALUES (?, ?, CURRENT_TIMESTAMP())")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStatment.Close()

	insertResult, err := insertStatment.Exec(userID, source)
	if err != nil {
		log.Fatal(err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return int(id)
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
		var description = item.Description
		if len(description) > 0 {
			var min = int(math.Min(float64(len(description)), 398))
			description = description[:min]
		}
		_, err = insertStatment.Exec(sourceID, item.Title, item.Link, description, pubdate)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getUserSources(db *sql.DB, userID int) []ResponseToClient {
	rows, err := db.Query("SELECT id, link FROM sources WHERE user_id=?", userID)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	type SourceInfo struct {
		ID   int
		Link string
	}
	var sourcesInfos []*SourceInfo
	for rows.Next() {
		var sourceInfo SourceInfo
		err := rows.Scan(&sourceInfo.ID, &sourceInfo.Link)
		if err != nil {
			log.Print(err)
		}
		sourcesInfos = append(sourcesInfos, &sourceInfo)
	}
	// fmt.Printf("sources %v %v", sourcesInfos[0].ID, sourcesInfos[0].Link)

	var response []ResponseToClient

	for _, source := range sourcesInfos {
		info, err := db.Query("SELECT title, link, description, pubDate FROM items WHERE source_id=?", source.ID)
		if err != nil {
			log.Println(err)
		}
		defer info.Close()

		items := make([]Item, 0)
		for info.Next() {
			var item Item
			u, err := url.Parse(source.Link)
			if err != nil {
				log.Println(err)
			}
			var hostname = u.Hostname()
			var first = strings.Index(hostname, ".")
			var second = strings.LastIndex(hostname, ".")
			item.SourceName = hostname[first:second]
			var pubDate time.Time
			err = info.Scan(&item.Title, &item.Link, &item.Description, &pubDate)
			if err != nil {
				log.Println(err)
			}
			item.PubDate = pubDate.Format(time.RFC1123Z)
			items = append(items, item)
		}
		var sourceResponse = ResponseToClient{SourceID: source.Link, Items: items}
		response = append(response, sourceResponse)
	}

	return response
}

func getSession(db *sql.DB, token string) (int, error) {
	info, err := db.Query("SELECT user_id FROM sessions WHERE token=?", token)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer info.Close()
	if info.Next() == false {
		return 0, errors.New("User ID was not found")
	}
	var userID int
	err = info.Scan(&userID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return userID, nil
}

func getOrCreateToken(db *sql.DB, userID int) string {
	info, err := db.Query("SELECT token FROM sessions WHERE user_id=?", userID)
	if err != nil {
		log.Println(err)
	}
	defer info.Close()

	if info.Next() == true {
		var userToken string
		err := info.Scan(&userToken)
		if err != nil {
			log.Fatal(err)
		}
		return userToken
	}
	newToken := randStringBytes(100)
	insertStatment, err := db.Prepare("INSERT INTO sessions (user_id, token, created_at) VALUES (?, ?, CURRENT_TIMESTAMP())")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStatment.Close()
	_, err = insertStatment.Exec(userID, newToken)
	if err != nil {
		log.Fatal(err)
	}
	return newToken
}

func deleteUserSession(db *sql.DB, token string) error {
	result, err := db.Query("DELETE FROM sessions WHERE token=?", token)
	if err != nil {
		log.Print(err)
		return err
	}
	defer result.Close()
	return nil
}
