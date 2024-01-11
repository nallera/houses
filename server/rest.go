package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type RestClient interface {
	Get(url string) ([]byte, error)
}

func NewRestClient() RestClient {
	return &restClient{}
}

type restClient struct{}

func (r *restClient) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, resp.Body)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes, nil
}
