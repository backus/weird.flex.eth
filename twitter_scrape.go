package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TwitterScrapeSeedUser struct {
	Username string  `json:"username"`
	Id       *string `json:"id"`
}

func (user TwitterScrapeSeedUser) String() string {
	var idDisplay string
	if user.Id != nil {
		idDisplay = *user.Id
	} else {
		idDisplay = "Nil"
	}

	return fmt.Sprintf("TwitterScrapeSeedUser(username=%s, id=%s)", user.Username, idDisplay)
}

type TwitterScrapeSeedInstructions struct {
	Users []TwitterScrapeSeedUser `json:"users"`
}

func LoadTwitterScrapeSeed() (TwitterScrapeSeedInstructions, error) {
	seedFileContents, err := ioutil.ReadFile("data/seed.json")
	if err != nil {
		return TwitterScrapeSeedInstructions{}, err
	}

	var seed TwitterScrapeSeedInstructions
	err = json.Unmarshal(seedFileContents, &seed)
	if err != nil {
		return TwitterScrapeSeedInstructions{}, err
	}

	return seed, nil
}

func (seed TwitterScrapeSeedInstructions) Inflate(client TwitterClient) TwitterScrapeSeedInstructions {
	var usernames []string

	for _, user := range seed.Users {
		if user.Id == nil {
			usernames = append(usernames, user.Username)
		}
	}

	result, err := client.FakeLookupUsers(usernames)
	check(err)

	idMap := make(map[string]string)

	for _, user := range result.Data {
		idMap[user.Username] = user.Id
	}

	var inflatedSeedUsers []TwitterScrapeSeedUser

	for _, user := range seed.Users {
		newUser := user
		if newUser.Id == nil {
			value, isPresent := idMap[user.Username]

			if isPresent {
				newUser.Id = &value
			}
		}

		inflatedSeedUsers = append(inflatedSeedUsers, newUser)
	}

	return TwitterScrapeSeedInstructions{inflatedSeedUsers}
}
