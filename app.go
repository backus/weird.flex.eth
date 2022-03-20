package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type App struct {
	twitter   TwitterClient
	ens       ENSClient
	etherscan EtherscanClient
	seedUsers TwitterScrapeSeedInstructions
	users     map[string]TwitterUser
}

func BootstrapApp() App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	twitter := NewTwitterClient(bearerToken)
	ens := NewENSClient(os.Getenv("INFURA_URL"))
	etherscan := NewEtherscanClient(os.Getenv("ETHERSCAN_API_KEY"))

	seed, err := LoadTwitterScrapeSeed()
	check(err)

	seed = seed.Inflate(twitter)
	seed.Persist()

	userPool := seed.LoadFollowing(twitter)

	userMap := make(map[string]TwitterUser)

	for _, user := range userPool {
		userMap[user.Id] = user
	}

	return App{twitter, ens, etherscan, seed, userMap}
}
