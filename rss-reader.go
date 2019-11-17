package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var fileName string

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
}

func getHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	var rss = getRss()
	var jsonEncoder = json.NewEncoder(res)
	err := jsonEncoder.Encode(rss)
	if err != nil {
		log.Fatalf("Was not able to encode %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getRss() Rss {
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

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/list").HandlerFunc(getHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/public")))

	err := http.ListenAndServe(":8800", r)
	if err != nil {
		log.Fatal(err)
	}

	// fileName = "https://www.elle.com/rss/all.xml/"
	// resp, err := http.Get(fileName)
	// if err != nil {
	// 	log.Fatalf("Was not able to get %v", fileName)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("Status error: %v", err)
	// 	return
	// }

	// xmlFile, err := os.Create("response.xml")
	// if err != nil {
	// 	log.Fatalf("Was not able to create file %v", err)
	// }

	// defer xmlFile.Close()

	// if _, err := io.Copy(xmlFile, resp.Body); err != nil {
	// 	log.Fatalf("Was not able to copy %v", err)
	// }
}
