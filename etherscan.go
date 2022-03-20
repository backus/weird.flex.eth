package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
)

type EtherscanClient struct {
	apiKey string
	cache  FileSystemCache
}

func NewEtherscanClient(apiKey string) EtherscanClient {
	return EtherscanClient{apiKey, NewFileSystemCache("data/eth")}
}

const ApiUrl = "https://api.etherscan.io/api"

type ETHAddress string

func (address ETHAddress) CacheKey() string {
	return string(address)
}

type GetBalanceResponse struct {
	Status  string
	Message string
	Result  string
}

type BalanceCheck ETHAddress

func (subject BalanceCheck) CacheKey() string {
	return fmt.Sprintf("%s.balance", string(subject))
}

func (client EtherscanClient) CachedGetBalance(address ETHAddress) GetBalanceResponse {
	logger.Debug("Looking up balance for %s", address)
	var result GetBalanceResponse
	balanceCheck := BalanceCheck(address)

	data, err := client.cache.WithRawCache(balanceCheck, func() ([]byte, error) {
		balanceResult := client.GetBalance(address)

		return balanceResult, nil
	})

	check(err)
	json.Unmarshal(data, &result)

	// if client.cache.IsCached(balanceCheck) {
	// 	rawData := string(client.cache.ReadCache(balanceCheck))
	// 	json.Unmarshal([]byte(rawData), &result)
	// } else {
	// 	result = client.GetBalance(ETHAddress(balanceCheck))
	// 	serialized, err := json.MarshalIndent(result, "", "  ")
	// 	check(err)

	// 	// Etherscan rate limit = 5 requests per second
	// 	time.Sleep(200 * time.Millisecond)

	// 	client.cache.WriteCache(balanceCheck, []byte(serialized))
	// }

	return result
}

func (client EtherscanClient) GetBalance(address ETHAddress) []byte {
	logger.Debug("Looking up balance for %s", address)

	url := apiUrl(map[string]string{
		"module":  "account",
		"action":  "balance",
		"address": string(address),
		"tag":     "latest",
		"apikey":  client.apiKey,
	})

	body, err := StrictGetRequest(url, nil)
	check(err)

	return body

	// var result GetBalanceResponse
	// json.Unmarshal(body, &result)

	// return result
}

type GetPriceResponse struct {
	Status  string
	Message string
	Result  struct {
		Ethusd string `json:"ethusd"`
	}
}

func (client EtherscanClient) GetETHUSDPrice() *big.Float {
	logger.Debug("Fetching ETH/USD price")

	url := apiUrl(map[string]string{
		"module": "stats",
		"action": "ethprice",
		"apikey": client.apiKey,
	})

	responseBody, err := StrictGetRequest(url, nil)
	check(err)

	var result GetPriceResponse

	json.Unmarshal(responseBody, &result)

	price, err := parseBigFloat(result.Result.Ethusd)
	check(err)
	return price
}

func apiUrl(params map[string]string) string {
	baseUrl, err := url.Parse(ApiUrl)
	check(err)

	query := baseUrl.Query()

	for key, value := range params {
		query.Add(key, value)
	}

	baseUrl.RawQuery = query.Encode()

	return baseUrl.String()
}
