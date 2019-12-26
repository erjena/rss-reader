package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const userIDContextKey = "userID"
const sessionCookie = "Session-Token"

func setupServer(db *sql.DB) {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return sessionTokenMiddleware(db, next)
	})
	r.Methods("GET").Path("/api/checkLoggedIn").HandlerFunc(handleCheckLoggedIn)
	r.Methods("GET").Path("/api/list").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		getHandler(db, res, req)
	})
	r.Methods("POST").Path("/api/login").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handleLogin(db, res, req)
	})
	r.Methods("POST").Path("/api/register").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handleRegister(db, res, req)
	})
	r.Methods("POST").Path("/api/sources").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		addSources(db, res, req)
	})
	r.Methods("POST").Path("/api/logout").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		hadleDeleteUserSession(db, res, req)
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/public")))

	err := http.ListenAndServe(":8800", r)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCheckLoggedIn(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	return
}

func handleLogin(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	type RequestBody struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}
	var requestBody *RequestBody
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestBody.Username == nil {
		log.Println("Missing username")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestBody.Password == nil {
		log.Println("Missing password")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// check username and password
	id, salt, hash, err := getUserInfo(db, *requestBody.Username)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	newHash := createHash(*requestBody.Password + salt)
	if newHash != hash {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userToken := getOrCreateToken(db, id)
	setCookie(res, sessionCookie, userToken)
}

func handleRegister(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	type RequestBody struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}
	var requestBody *RequestBody
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestBody.Username == nil {
		log.Println("Missing username")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestBody.Password == nil {
		log.Println("Missing password")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(*requestBody.Username) == 0 {
		log.Println("Missing username")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(*requestBody.Password) == 0 {
		log.Println("Missing password")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	salt := randStringBytes(40)
	passSalt := *requestBody.Password + salt
	hash := createHash(passSalt)
	id := insertUserInfo(db, *requestBody.Username, salt, hash)

	userToken := getOrCreateToken(db, id)
	setCookie(res, sessionCookie, userToken)
}

func getHandler(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(userIDContextKey)
	if userID == nil {
		log.Print("User `id wad not found")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var responseData = getUserSources(db, userID.(int))
	data, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("Was not able to stringify %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func setCookie(res http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:  name,
		Value: value,
	}
	http.SetCookie(res, &cookie)
}

func sessionTokenMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.RequestURI, "/api/") {
			next.ServeHTTP(res, req)
			return
		}
		if strings.HasPrefix(req.RequestURI, "/api/login") {
			next.ServeHTTP(res, req)
			return
		}
		if strings.HasPrefix(req.RequestURI, "/api/register") {
			next.ServeHTTP(res, req)
			return
		}
		token, err := req.Cookie(sessionCookie)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := getSession(db, token.Value)
		if err != nil {
			log.Print(err)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), userIDContextKey, userID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func addSources(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	type SourceInfo struct {
		Link *string `json:"link"`
	}
	var info *SourceInfo
	err := json.NewDecoder(req.Body).Decode(&info)
	if err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Link == nil {
		log.Printf("Missing link")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := req.Context().Value(userIDContextKey)
	if userID == nil {
		log.Print("User `id wad not found")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := insertSource(db, userID.(int), *info.Link)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		return
	}
	source := Source{id, *info.Link, nil}
	crawlSingleSource(&source, db, nil)
}

func hadleDeleteUserSession(db *sql.DB, res http.ResponseWriter, req *http.Request) {
	token, err := req.Cookie(sessionCookie)
	if err != nil {
		log.Print(err)
	}
	err = deleteUserSession(db, token.Value)
	if err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	setCookie(res, sessionCookie, "")
	res.WriteHeader(http.StatusOK)
}
