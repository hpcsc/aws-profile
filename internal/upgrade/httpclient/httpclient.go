package httpclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func GetUrlWithAuthorization(url string, authorization string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request for url %s: %v", url, err)
	}

	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}

	c := http.Client{}
	response, err := c.Do(request)

	if err != nil {
		return nil, fmt.Errorf("failed to get url %s: %v", url, err)
	}

	return readBodyContent(url, response)
}

func GetUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url %s: %v", url, err)
	}

	return readBodyContent(url, response)
}

func readBodyContent(url string, response *http.Response) ([]byte, error) {
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

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %v", filepath, err)
	}

	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file from %s: %v", url, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get url %s: status code: %d", url, response.StatusCode)
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write file to %s: %v", filepath, err)
	}

	return nil
}
