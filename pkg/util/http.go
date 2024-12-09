package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func HTTPPost(url string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func HTTPDelete(url string, body io.Reader) ([]byte, error) {
	fmt.Println(url)
	req, _ := http.NewRequest(http.MethodDelete, url, body)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		return nil, errors.New(resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func HTTPPostWithMultipart(req *http.Request) ([]byte, error) {
	fmt.Println(req.Header, req.Body, req.MultipartForm)
	client := new(http.Client)

	resp, err := client.Do(req)
	fmt.Println(resp.Status)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)

}
