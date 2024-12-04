package util

import (
	"io"
	"net/http"
)

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func HttpPost(url string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
