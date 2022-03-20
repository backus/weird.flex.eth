package main

import (
	"log"
	"math/big"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/joho/godotenv"
)

type ENSResolution struct {
	address ETHAddress
	balance big.Float
}

type ENSReport struct {
	ens        ENSDomain
	resolution ENSResolution
}

type UserReport struct {
	user    TwitterUser
	domains []ENSReport
}

var logger = NewLogger()

func main() {
	logger.SetLevel(LogLevelInfo)

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

	logger.Debug("Total users in pool: %d\n\n", len(userPool))

	ethPrice := etherscan.GetETHUSDPrice()
	logger.Info("ETH/USD price: %f\n", ethPrice)

	userReports := []UserReport{}

	for _, user := range userPool {
		domains := user.ENSDomains()
		if len(domains) == 0 {
			continue
		}
		ensReports := []ENSReport{}
		for _, domain := range domains {
			address, err := ens.CachedResolve(domain)
			if err != nil {
				continue
			}
			balance := etherscan.CachedGetBalance(address)
			wei, err := parseBigFloat(balance.Result)
			check(err)
			eth := weiToEth(wei)

			ensReport := ENSReport{domain, ENSResolution{address, *eth}}
			ensReports = append(ensReports, ensReport)
		}
		if len(ensReports) == 0 {
			continue
		}
		userReport := UserReport{user, ensReports}
		userReports = append(userReports, userReport)
	}

	logger.Info("Total users with ENS domains: %d\n", len(userReports))

	logger.Debug("Resolving ENS domains to ETH addresses...\n\n")

	sort.SliceStable(userReports, func(i, j int) bool {
		return userReports[i].domains[0].resolution.balance.Cmp(&userReports[j].domains[0].resolution.balance) != -1
	})

	for userIndex, userReport := range userReports {
		logger.Info("\n%d. %s - %s\n", userIndex+1, userReport.user.Username, userReport.user.ShortDescription())

		for _, ensReport := range userReport.domains {
			dollarBalanceFloat := big.NewFloat(1e18)
			dollarBalanceFloat.Mul(&ensReport.resolution.balance, ethPrice)
			dollarBalance, _ := dollarBalanceFloat.Float32()

			logger.Info("   %-25s = %-5.2f ETH = $%s\n", ensReport.ens, &ensReport.resolution.balance, humanize.Comma(int64(dollarBalance)))
		}
	}
}
