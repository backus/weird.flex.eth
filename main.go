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

	seed, err := LoadTwitterScrapeSeed()
	check(err)
	seed = seed.Inflate(client)
	seed.Persist()

	userPool := seed.LoadFollowing(client)

	fmt.Printf("Total users in pool: %d\n", len(userPool))
	fmt.Println()

	extractedDomains := ENSDomains(userPool)

	fmt.Printf("Total ENS domains extracted: %d\n", len(extractedDomains))

	for i, domain := range extractedDomains {
		fmt.Printf("%d. %s\n", i, domain)
	}
}
