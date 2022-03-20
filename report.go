package main

import (
	"math"
	"math/big"
	"sort"
)

type ENSReport struct {
	domain  ENSDomain
	valid   bool
	address *ETHAddress
	balance *big.Float // Denominated in ETH, not Wei
}

type UserENSReportMap map[string][]ENSReport

type UserENSReport struct {
	user          TwitterUser
	ensReportList ENSReportList
}

type ENSReportList struct {
	reports []ENSReport
}

func (reportList ENSReportList) totalBalance() float64 {
	total := float64(0)

	for _, report := range reportList.reports {
		if !report.valid {
			continue
		}

		balance64, _ := report.balance.Float64()
		total += balance64
	}

	return total
}

func (reportList ENSReportList) totalBalanceUSD(ethPrice *big.Float) float64 {
	price64, _ := ethPrice.Float64()

	return math.Round(reportList.totalBalance() * price64)
}

func (reportList ENSReportList) domains() []string {
	domains := []string{}

	for _, report := range reportList.reports {
		domains = append(domains, string(report.domain))
	}

	return domains
}

func BuildReport(app App, users map[string]TwitterUser) UserENSReportMap {
	userReport := UserENSReportMap{}

	for _, user := range users {
		domains := user.ENSDomains()

		if len(domains) == 0 {
			continue
		}

		reports := []ENSReport{}

		for _, domain := range domains {
			var report ENSReport

			address, err := app.ens.CachedResolve(domain)

			if err != nil {
				report = ENSReport{domain, false, nil, nil}
			} else {
				balance := app.etherscan.CachedGetBalance(address)
				wei, err := parseBigFloat(balance.Result)
				check(err)
				eth := weiToEth(wei)

				report = ENSReport{domain, true, &address, eth}
			}

			reports = append(reports, report)
		}

		userReport[user.Id] = reports
	}

	return userReport
}

func (reportMap UserENSReportMap) SortedReportList(userMap map[string]TwitterUser) []UserENSReport {
	reportList := []UserENSReport{}

	for userId, reports := range reportMap {
		user := userMap[userId]
		userEnsReport := UserENSReport{user, ENSReportList{reports}}

		reportList = append(reportList, userEnsReport)
	}

	sort.SliceStable(reportList, func(i, j int) bool {
		return reportList[i].ensReportList.totalBalance() > reportList[j].ensReportList.totalBalance()
	})

	return reportList
}
