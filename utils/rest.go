package utils

import (
	"net/http"
	"fmt"
	"io"
)

/**
 * Http get request to a given URL.
 * @param url The URL to send the GET request to.
 * @return The response body as a string, or an error if the request fails.
 */
func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: %s, status code: %d", url, resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}	