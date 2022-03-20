package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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
	var result GetBalanceResponse

	balanceCheck := BalanceCheck(address)

	if client.cache.IsCached(balanceCheck) {
		rawData := string(client.cache.ReadCache(balanceCheck))
		json.Unmarshal([]byte(rawData), &result)
	} else {
		result = client.GetBalance(ETHAddress(balanceCheck))
		serialized, err := json.MarshalIndent(result, "", "  ")
		check(err)

		// Etherscan rate limit = 5 requests per second
		time.Sleep(200 * time.Millisecond)

		client.cache.WriteCache(balanceCheck, []byte(serialized))
	}

	return result
}

func (client EtherscanClient) GetBalance(address ETHAddress) GetBalanceResponse {
	url := apiUrl(map[string]string{
		"module":  "account",
		"action":  "balance",
		"address": string(address),
		"tag":     "latest",
		"apikey":  client.apiKey,
	})

	response, err := http.Get(url)
	check(err)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	check(err)

	var result GetBalanceResponse

	json.Unmarshal(body, &result)

	return result
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
