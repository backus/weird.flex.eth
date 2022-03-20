package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/dustin/go-humanize"
)

type ENSResolution struct {
	address ETHAddress
	balance big.Float
}

type ENSReportOld struct {
	ens        ENSDomain
	resolution ENSResolution
}

type UserReport struct {
	user    TwitterUser
	domains []ENSReportOld
}

var logger = NewLogger()

func main() {
	logger.SetLevel(LogLevelInfo)

	app := BootstrapApp()

	logger.Debug("Total users in pool: %d\n\n", len(app.users))

	ethPrice := app.etherscan.GetETHUSDPrice()
	logger.Info("ETH/USD price: %f\n", ethPrice)

	userReport := BuildReport(app, app.users)

	sortedResults := userReport.SortedReportList(app.users)

	heading := fmt.Sprintf("| %-15s | %-50s | %11s | %12s |\n", "Twitter handle", "ENS Domain", "ETH Balance", "USD Balance")
	fmt.Printf(heading)
	fmt.Printf("%s\n", strings.Repeat("-", len(heading)))

	for _, userReport := range sortedResults {
		fmt.Printf(
			"| %-15s | %-50s | %11.2f | $%11s | \n",
			userReport.user.Username,
			strings.Join(userReport.ensReportList.domains(), ", "),
			userReport.ensReportList.totalBalance(),
			humanize.Commaf(userReport.ensReportList.totalBalanceUSD(ethPrice)),
		)
	}
}
