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
	twitter := NewTwitterClient(bearerToken)
	ens := NewENSClient(os.Getenv("INFURA_URL"))
	etherscan := NewEtherscanClient(os.Getenv("ETHERSCAN_API_KEY"))

	seed, err := LoadTwitterScrapeSeed()
	check(err)
	seed = seed.Inflate(twitter)
	seed.Persist()

	userPool := seed.LoadFollowing(twitter)

	fmt.Printf("Total users in pool: %d\n", len(userPool))
	fmt.Println()

	extractedDomains := ENSDomains(userPool)

	fmt.Printf("Total ENS domains extracted: %d\n", len(extractedDomains))

	for i, domain := range extractedDomains {
		fmt.Printf("%d. %s\n", i+1, domain)
	}

	fmt.Println("Resolving ENS domains to ETH addresses...")

	ethAddresses := []ETHAddress{}

	for i, domain := range extractedDomains {
		address, err := ens.CachedResolve(domain)

		if err != nil {
			fmt.Printf("%d. %s -> (resolve failed)\n", i+1, domain)
		} else {
			ethAddresses = append(ethAddresses, address)
			fmt.Printf("%d. %s -> %s\n", i+1, domain, address)
		}
	}

	for i, ethAddress := range ethAddresses {
		balance := etherscan.CachedGetBalance(ethAddress)
		wei, err := parseWei(balance.Result)
		eth := weiToEth(wei)
		check(err)
		fmt.Printf("%d. %s balance = %f\n", i+1, ethAddress, eth)
	}

	fmt.Println("Done!")
}
