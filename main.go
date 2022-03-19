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

	twitterAuth := TwitterAuth{bearerToken}

	fmt.Println("Hello world!")
	fmt.Printf("Twitter auth debug = %s\n", twitterAuth.DebugText())
}
