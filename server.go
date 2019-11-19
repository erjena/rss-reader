package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupServer(db *sql.DB) {
	r := mux.NewRouter()
	r.Methods("GET").Path("/list").HandlerFunc(getHandler)
	r.Methods("POST").Path("/sources").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		addSources(db, res, req)
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/public")))

	err := http.ListenAndServe(":8800", r)
	if err != nil {
		log.Fatal(err)
	}
}

func getHandler(res http.ResponseWriter, req *http.Request) {
	var rss = getRss("username")
	data, err := json.Marshal(rss)
	if err != nil {
		log.Fatalf("Was not able to stringify %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func addSources(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	type SourceInfo struct {
		User *string `json:"user"`
		Link *string `json:"link"`
	}
	var info *SourceInfo
	err := json.NewDecoder(req.Body).Decode(&info)
	if err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.User == nil {
		log.Printf("Missing user name")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Link == nil {
		log.Printf("Missing link")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	insertSource(db, *info.User, *info.Link)
}

// func addFeedElements(db *sql.DB, res http.ResponseWriter, req *http.Request) {
// 	insertItems()
// }
