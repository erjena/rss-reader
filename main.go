package main

var fileName string

func main() {
	var db = dbConnection()
	setupServer(db)

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
