package main

import (
	"math/rand"
	"time"
)

var fileName string

func main() {
	rand.Seed(time.Now().UnixNano())
	var db = dbConnection()
	//go crawlWrapper(db)
	setupServer(db)
}
