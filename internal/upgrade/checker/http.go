package checker

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url %s: %v", url, err)
	}

	defer response.Body.Close()

	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for %s: %v", url, err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get url %s: status code: %d, body: %s", url, response.StatusCode, bodyContent)
	}

	return bodyContent, nil
}
