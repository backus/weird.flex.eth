package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

type IgnoreList struct {
	path string
	list []ENSDomain
}

func LoadIgnoreList(path string) []ENSDomain {
	result, err := os.ReadFile(path)
	check(err)

	data := []ENSDomain{}
	json.Unmarshal(result, &data)

	return data
}

func NewIgnoreList(path string) IgnoreList {
	pathType, err := checkPathType(path)
	check(err)

	data := []ENSDomain{}

	if pathType == IsFile {
		data = LoadIgnoreList(path)
	} else {
		serialized, err := json.Marshal(data)
		check(err)

		os.WriteFile(path, serialized, 0777)
	}

	return IgnoreList{path, data}
}

func (list IgnoreList) Has(domain ENSDomain) bool {
	for _, item := range list.list {
		if item == domain {
			return true
		}
	}

	return false
}

func (list IgnoreList) Add(domain ENSDomain) {
	currentList := LoadIgnoreList(list.path)
	newList := append(currentList, domain)

	serialized, err := json.Marshal(newList)
	check(err)

	os.WriteFile(list.path, serialized, 0777)
}

type ENSClient struct {
	client     *ethclient.Client
	cache      FileSystemCache
	ignoreList IgnoreList
}

func NewENSClient(infuraUrl string) ENSClient {
	client, err := ethclient.Dial(infuraUrl)
	check(err)

	cache := NewFileSystemCache("data/ens")

	ignoreList := NewIgnoreList("data/ens/ignore.json")

	return ENSClient{client, cache, ignoreList}
}

func (domain ENSDomain) CacheKey() string {
	return string(domain)
}

func (client ENSClient) CachedResolve(domain ENSDomain) (ETHAddress, error) {
	if client.ignoreList.Has(domain) {
		return "", errors.New("domain failed to resolve and has been marked as ignored")
	}

	data, err := client.cache.WithRawCache(domain, func() ([]byte, error) {
		address, err := client.Resolve(domain)
		if err != nil {
			return nil, err
		} else {
			return []byte(address), nil
		}
	})

	if err != nil {
		return ETHAddress(""), err
	}

	return ETHAddress(data), nil
}

func (client ENSClient) Resolve(domain ENSDomain) (ETHAddress, error) {
	address, err := ens.Resolve(client.client, string(domain))

	if err != nil {
		client.ignoreList.Add(domain)
		return "", err
	}

	return ETHAddress(address.String()), nil
}
