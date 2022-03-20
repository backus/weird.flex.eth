package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gosimple/slug"
)

type TwitterAuth struct {
	bearer string
}

type TwitterClient struct {
	auth TwitterAuth
}

type CachingTwitterClient struct {
	client TwitterClient
	cache  FileSystemCache
}

const Hostname string = "https://api.twitter.com"

func NewTwitterClient(bearerToken string) TwitterClient {
	auth := TwitterAuth{bearerToken}
	client := TwitterClient{auth}

	return client
}

func NewCachedTwitterClient(bearerToken string) CachingTwitterClient {
	client := NewTwitterClient(bearerToken)
	cache := NewFileSystemCache("data")

	return CachingTwitterClient{client, cache}
}

func (auth TwitterAuth) DebugText() string {
	return auth.bearer
}

func (tw TwitterClient) String() string {
	return fmt.Sprintf("TwitterClient(bearer=%s)", tw.auth.bearer)
}

type TwitterUser struct {
	Id       string
	Name     string
	Username string
}

type TwitterUserListData struct {
	Data []TwitterUser
}

type ListFollowingOptions struct {
	PaginationToken *string
}

type ListFollowingRequestInput struct {
	path   string
	params map[string]string
}

func (req ListFollowingRequestInput) CacheKey() string {
	token, isPresent := req.params["pagination_token"]
	var tokenDisplay string

	if isPresent {
		tokenDisplay = token
	} else {
		tokenDisplay = "Nil"
	}

	return slug.Make(req.path + "?pagination_token=" + tokenDisplay)
}

var UserFields = []string{
	"created_at",
	"description",
	"entities",
	"id",
	"location",
	"name",
	"pinned_tweet_id",
	"profile_image_url",
	"protected",
	"public_metrics",
	"url",
	"username",
	"verified",
	"withheld",
}

type TwitterUserFollowing struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Url         *string `json:"url"`
	Username    string  `json:"username"`
	Description string  `json:"description"`
}

const MaxDescriptionLength = 50

func (user TwitterUserFollowing) ShortDescription() string {
	desc := strings.Split(user.Description, "\n")[0]

	if len(desc) > MaxDescriptionLength {
		desc = fmt.Sprintf("%s...", desc[:MaxDescriptionLength])
	}

	return desc
}

type PaginatedUserList struct {
	Data []TwitterUserFollowing
	Meta struct {
		ResultCount int     `json:"result_count"`
		NextToken   *string `json:"next_token"`
	}
}

func (tw TwitterClient) ListAllFollowing(userId string, cache FileSystemCache) []TwitterUserFollowing {
	var following []TwitterUserFollowing

	options := ListFollowingOptions{}

	for {
		followingPage := tw.ListFollowing(userId, cache, options)
		following = append(following, followingPage.Data...)
		options.PaginationToken = followingPage.Meta.NextToken

		if options.PaginationToken == nil {
			break
		}
	}

	return following
}

func (tw TwitterClient) ListFollowing(userId string, cache FileSystemCache, options ListFollowingOptions) PaginatedUserList {
	path := fmt.Sprintf("/2/users/%s/following", userId)
	params := make(map[string]string)
	params["max_results"] = "1000"
	params["user.fields"] = strings.Join(UserFields, ",")

	if options.PaginationToken != nil {
		params["pagination_token"] = *options.PaginationToken
	}

	requestInput := ListFollowingRequestInput{path, params}
	url := apiRoute(path, params)
	var paginatedUserList PaginatedUserList
	var rawResponse string

	if cache.IsCached(requestInput) {
		logger.Debug("Cache hit for %s\n", requestInput.CacheKey())
		rawResponse = string(cache.ReadCache(requestInput))
	} else {
		logger.Debug("Performing live request for %s\n", requestInput.CacheKey())

		rawResponse = submitGetRequest(tw, url)
		cache.WriteCache(requestInput, []byte(rawResponse))
	}

	json.Unmarshal([]byte(rawResponse), &paginatedUserList)

	return paginatedUserList
}

func (tw TwitterClient) LookupUsers(usernames []string) (TwitterUserListData, error) {
	client := &http.Client{}

	serializedQuery := strings.Join(usernames, ",")
	uri := apiRoute("/2/users/by", map[string]string{"usernames": serializedQuery})

	req, err := http.NewRequest("GET", uri, nil)
	check(err)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tw.auth.bearer))

	response, err := client.Do(req)
	check(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	check(err)

	var userList TwitterUserListData
	json.Unmarshal([]byte(body), &userList)

	return userList, nil
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

func submitGetRequest(tw TwitterClient, url string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	check(err)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tw.auth.bearer))

	response, err := client.Do(req)
	check(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)
	strBody := string(body)

	if response.StatusCode != 200 {
		log.Fatalf("Error! Received status code %s while requesting %s\nResponse body = %s\n", response.Status, url, strBody)
	}

	return strBody
}
