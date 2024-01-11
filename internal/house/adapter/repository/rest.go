package repository

import (
	"encoding/json"
	"fmt"
	"houses/internal/house"
	"houses/server"
)

var HouseBaseURL = "http://app-homevision-staging.herokuapp.com/api_project/houses"

type RestRepositoryClient struct {
	restClient server.RestClient
}

func NewRestRepositoryClient(restClient server.RestClient) *RestRepositoryClient {
	return &RestRepositoryClient{
		restClient: restClient,
	}
}

func (r *RestRepositoryClient) GetHouses(perPage, pageNumber int) ([]*house.House, error) {
	url := fmt.Sprintf("%s?page=%d&per_page=%d", HouseBaseURL, pageNumber, perPage)

	housesJSON, err := r.restClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get houses: %v", err)
	}

	var houses []*house.House

	jsonErr := json.Unmarshal(housesJSON, &houses)

	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse houses data: %v", err)
	}

	return houses, nil
}
