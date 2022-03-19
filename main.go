package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	// twitterAuth := TwitterAuth{bearerToken}
	client := NewTwitterClient(bearerToken)

	fmt.Println("Hello world!")
	fmt.Printf("Twitter auth debug = %s\n", client)
	// client.writeToRequestBin("go client attempt 1")

	sampleJson1 := `{ "name": "John", "age": 28 }`
	sampleJson2 := `{ "name": "John", "age": 28, "dob": "1993-10-07" }`
	sampleJson3 := `{ "name": "John" }`
	ParseUserJSON(sampleJson1)
	ParseUserJSON(sampleJson2)
	ParseUserJSON(sampleJson3)
	println("Done")
}
