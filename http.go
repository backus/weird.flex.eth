package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Finish constructing and submit a GET request, returning any error encountered
// as well as returning an error if the response status is not 200 OK.
func StrictGetRequest(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected a 200 OK status code, but received %s while requesting %s", response.Status, url)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
