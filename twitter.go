package main

import (
	"encoding/json"
	"fmt"
	"log"
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

const Hostname string = "https://api.twitter.com"

func NewTwitterClient(bearerToken string) TwitterClient {
	auth := TwitterAuth{bearerToken}
	client := TwitterClient{auth}

	return client
}

// API Response type for /2/users/by/:username
type TwitterAPIUsersList struct {
	Data []struct {
		Id       string
		Name     string
		Username string
	}
}

type TwitterAPIListFollowingRequestOptions struct {
	PaginationToken *string
}

// This type is used to create a Cacheable request
type TwitterAPIListFollowingRequestInput struct {
	path   string
	params map[string]string
}

func (req TwitterAPIListFollowingRequestInput) CacheKey() string {
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

// This is a struct a single user in the list we get back from /2/users/:id/following
// More importantly though, this is also the Twitter user struct we pass around for analysis
type TwitterUser struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Url         *string `json:"url"`
	Username    string  `json:"username"`
	Description string  `json:"description"`
}

const MaxDescriptionLength = 50

// When we print out a user as part of the final report, we want to show some more context
// on the user via their bio, but we don't want to print the entire thing
func (user TwitterUser) ShortDescription() string {
	desc := strings.Split(user.Description, "\n")[0]

	if len(desc) > MaxDescriptionLength {
		desc = fmt.Sprintf("%s...", desc[:MaxDescriptionLength])
	}

	return desc
}

// Struct for the API Response for a single page from /2/users/:id/following
type PaginatedUserList struct {
	Data []TwitterUser
	Meta struct {
		ResultCount int     `json:"result_count"`
		NextToken   *string `json:"next_token"`
	}
}

// Given a UserID, get a single page of 1000 users they are following via /2/users/:id/following
func (tw TwitterClient) CachedListFollowing(userId string, cache FileSystemCache, options TwitterAPIListFollowingRequestOptions) PaginatedUserList {
	// TODO: This method could be slimmed down a lot if I change up how I compute the cache key
	path := fmt.Sprintf("/2/users/%s/following", userId)
	params := make(map[string]string)

	if options.PaginationToken != nil {
		params["pagination_token"] = *options.PaginationToken
	}

	requestInput := TwitterAPIListFollowingRequestInput{path, params}
	paginatedUserList, err := WithJSONCache(cache, requestInput, func() (PaginatedUserList, error) {
		return tw.ListFollowing(userId, cache, options), nil
	})
	check(err)

	return paginatedUserList
}

func (tw TwitterClient) ListFollowing(userId string, cache FileSystemCache, options TwitterAPIListFollowingRequestOptions) PaginatedUserList {
	path := fmt.Sprintf("/2/users/%s/following", userId)
	params := make(map[string]string)
	params["max_results"] = "1000"
	params["user.fields"] = strings.Join(UserFields, ",")

	if options.PaginationToken != nil {
		params["pagination_token"] = *options.PaginationToken
	}

	url := apiRoute(path, params)
	var paginatedUserList PaginatedUserList
	var rawResponse []byte

	logger.Debug("Performing live request for users %s is following\n", userId)

	rawResponse, err := tw.get(url)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(rawResponse, &paginatedUserList)

	return paginatedUserList
}

// Expand a list of usernames into user IDs.
func (tw TwitterClient) LookupUsers(usernames []string) (TwitterAPIUsersList, error) {
	serializedQuery := strings.Join(usernames, ",")
	uri := apiRoute("/2/users/by", map[string]string{"usernames": serializedQuery})

	responseBody, err := tw.get(uri)
	if err != nil {
		return TwitterAPIUsersList{}, err
	}

	var userList TwitterAPIUsersList
	json.Unmarshal(responseBody, &userList)

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

func (tw TwitterClient) get(url string) ([]byte, error) {
	return StrictGetRequest(url, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", tw.auth.bearer),
	})
}

// Facade that reads each page from `ListFollowing`.
// NOTE: This doesn't play well with Twitter rate limits. If you hit a rate limit, just run the program again.
func (tw TwitterClient) ListAllFollowing(userId string, cache FileSystemCache) []TwitterUser {
	var following []TwitterUser

	options := TwitterAPIListFollowingRequestOptions{}

	for {
		followingPage := tw.CachedListFollowing(userId, cache, options)
		following = append(following, followingPage.Data...)
		options.PaginationToken = followingPage.Meta.NextToken

		if options.PaginationToken == nil {
			break
		}
	}

	return following
}
