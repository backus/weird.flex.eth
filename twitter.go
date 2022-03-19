package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TwitterAuth struct {
	bearer string
}

type TwitterClient struct {
	auth TwitterAuth
}

const Hostname string = "https://eoukkwxovtn2fsw.m.pipedream.net"

func NewTwitterClient(bearerToken string) TwitterClient {
	auth := TwitterAuth{bearerToken}
	client := TwitterClient{auth}

	return client
}

func (auth TwitterAuth) DebugText() string {
	return auth.bearer
}

// Kind of equivalent to defining TwitterClient#inspect in Ruby
func (tw TwitterClient) String() string {
	return fmt.Sprintf("TwitterClient(bearer=%s)", tw.auth.bearer)
}

func (tw TwitterClient) writeToRequestBin(value string) {
	// url :=
	// curl -d '{
	// 	"type": "cURL"
	// }'   -H "Content-Type: application/json"   https://eoukkwxovtn2fsw.m.pipedream.net
	payload := map[string]string{"test": value}

	rawBody := postReq(tw, apiRoute("/route", map[string]string{"foo": "bar"}), payload)
	fmt.Printf("Response body (raw): %s\n", rawBody)
	// serialized, err := json.Marshal(payload)

	// check(err)

	// bodyBuffer := bytes.NewBuffer(serialized)
	// response, err := http.Post(url, "application/json", bodyBuffer)

	// check(err)

	// defer response.Body.Close()
	// body, err := ioutil.ReadAll(response.Body)
	// check(err)
	// strBody := string(body)
	// fmt.Printf("Response body: %s\n", strBody)
}

func apiRoute(path string, query map[string]string) string {
	baseUrl, err := url.Parse(Hostname)
	check(err)
	baseUrl.Path = path

	givenQuery := baseUrl.Query()

	for key, value := range query {
		givenQuery.Add(key, value)
	}

	baseUrl.RawQuery = givenQuery.Encode()

	return baseUrl.String()
}

func postReq(tw TwitterClient, url string, payload interface{}) string {
	serialized, err := json.Marshal(payload)

	check(err)

	bodyBuffer := bytes.NewBuffer(serialized)
	response, err := http.Post(url, "application/json", bodyBuffer)

	check(err)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	check(err)
	strBody := string(body)
	return strBody
}

type User struct {
	Name string `json:"name`
	Age  *int   `json:"age",omitempty`
}

func (u User) String() string {
	if u.Age != nil {
		return fmt.Sprintf("User(name=%s, age=%d)", u.Name, *u.Age)
	} else {
		return fmt.Sprintf("User(name=%s, age=Nil)", u.Name)
	}
}

func ParseUserJSON(raw string) {
	var user User
	json.Unmarshal([]byte(raw), &user)
	fmt.Printf("Parsed user: %s\n", user)
}
