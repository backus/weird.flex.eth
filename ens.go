package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

type ENSClient struct {
	client *ethclient.Client
	cache  FileSystemCache
}

func NewENSClient(infuraUrl string) ENSClient {
	client, err := ethclient.Dial(infuraUrl)
	check(err)

	cache := NewFileSystemCache("data/ens")

	return ENSClient{client, cache}
}

func (domain ENSDomain) CacheKey() string {
	return string(domain)
}

func (client ENSClient) CachedResolve(domain ENSDomain) (string, error) {
	var address string
	if client.cache.IsCached(domain) {
		address = string(client.cache.ReadCache(domain))
	} else {
		result, err := client.Resolve(domain)

		if err != nil {
			return "", err
		}

		address = result

		client.cache.WriteCache(domain, []byte(address))
	}

	return address, nil
}

func (client ENSClient) Resolve(domain ENSDomain) (string, error) {
	address, err := ens.Resolve(client.client, string(domain))

	if err != nil {
		return "", err
	}

	return address.String(), nil
}
