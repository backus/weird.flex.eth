package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type TwitterScrapeSeedUser struct {
	Username string  `json:"username"`
	Id       *string `json:"id,omitempty"`
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

const SeedFile = "data/seed.json"

func LoadTwitterScrapeSeed() (TwitterScrapeSeedInstructions, error) {
	seedFileContents, err := ioutil.ReadFile(SeedFile)
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

	if len(usernames) == 0 {
		log.Println("All seed users already have IDs. No inflation necessary")
		return seed
	}

	result, err := client.LookupUsers(usernames)
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

func (seed TwitterScrapeSeedInstructions) Persist() {
	serializedSeed, err := json.Marshal(seed)
	check(err)

	err = ioutil.WriteFile(SeedFile, serializedSeed, 0644)
	check(err)
}
