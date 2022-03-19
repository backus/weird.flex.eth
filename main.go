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
	client := NewTwitterClient(bearerToken)

	// debugCache()
	seed, err := LoadTwitterScrapeSeed()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[initial] Twitter Scrape Seed: %+v\n", seed)

	seed = seed.Inflate(client)

	fmt.Printf("[inflated] Twitter Scrape Seed: %+v\n", seed)

	// twitterAuth := TwitterAuth{bearerToken}

	// fmt.Println("Hello world!")
	// fmt.Printf("Twitter auth debug = %s\n", client)
	// // client.writeToRequestBin("go client attempt 2")

	// result, err := client.LookupUsers([]string{"backus"})
	// check(err)
	// fmt.Println("Users List Result =", result)
	// sampleJson1 := `{ "name": "John", "age": 28 }`
	// sampleJson2 := `{ "name": "John", "age": 28, "dob": "1993-10-07" }`
	// sampleJson3 := `{ "name": "John" }`
	// // println(apiRoute("/foo/bar", map[string]string{"q1": "hello", "q2": "goodbye"}))
	// ParseUserJSON(sampleJson1)
	// ParseUserJSON(sampleJson2)
	// ParseUserJSON(sampleJson3)
	// println("Done")
}
