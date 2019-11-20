package main

var fileName string

func main() {
	var db = dbConnection()
	crawl(db)
	setupServer(db)
}
