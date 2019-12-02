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
	r.Methods("GET").Path("/list").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		getHandler(db, res, req)
	})
	r.Methods("POST").Path("/sources").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		addSources(db, res, req)
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/public")))

	err := http.ListenAndServe(":8800", r)
	if err != nil {
		log.Fatal(err)
	}
}

func getHandler(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	// fmt.Printf("request body %v", req.Body)
	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	res.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// log.Println(string(body))

	// type UserName struct {
	// 	User *string `json:"user"`
	// }
	// var userName UserName
	// err := json.NewDecoder(req.Body).Decode(&userName)
	// if err != nil {
	// 	log.Print(err)
	// 	res.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// err = json.Unmarshal(body, &userName)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Printf("user name %v", userName.User)

	var responseData = getUserSources(db, 1)
	data, err := json.Marshal(responseData)
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

	id := insertSource(db, *info.User, *info.Link)
	source := Source{id, *info.Link, nil}
	crawlSingleSource(&source, db)
}
