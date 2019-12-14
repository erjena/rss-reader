package main

var fileName string

func main() {
	var db = dbConnection()
	go crawlWrapper(db)
	setupServer(db)
}
