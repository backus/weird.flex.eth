package main

import (
	"regexp"
)

type ENSDomain string

func ENSDomains(users []TwitterUserFollowing) []ENSDomain {
	extractedDomains := []ENSDomain{}

	for _, user := range users {
		extractedDomains = append(extractedDomains, user.ENSDomains()...)
	}

	return extractedDomains
}

func (user TwitterUserFollowing) ENSDomains() []ENSDomain {
	var domains []ENSDomain

	domains = append(domains, findENSDomain(user.Username)...)
	domains = append(domains, findENSDomain(user.Name)...)
	domains = append(domains, findENSDomain(user.Description)...)

	return domains
}

/*
 * Based on first Google result for "ENS Domain Regular Expression"
 * @see https://www.regextester.com/111178
 *
 * NOTE: This was more helpful than scouring multiple ENS packages for an actual regex 🤷‍♂️
 */
var ensPattern = regexp.MustCompile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.eth`)

func findENSDomain(input string) []ENSDomain {
	matches := ensPattern.FindAll([]byte(input), -1)

	var strings []ENSDomain

	for _, match := range matches {
		strings = append(strings, ENSDomain(string(match)))
	}

	return strings
}
