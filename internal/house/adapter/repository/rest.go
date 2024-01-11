package repository

import (
	"encoding/json"
	"fmt"
	"houses/internal/house"
	"houses/server"
)

var HouseBaseURL = "http://app-homevision-staging.herokuapp.com/api_project/houses"

type RestHouses struct {
	Houses []*RestHouse `json:"houses"`
}

type RestHouse struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

func RestHouseToApp(rh *RestHouse) *house.House {
	return &house.House{
		Id:        rh.Id,
		Address:   rh.Address,
		Homeowner: rh.Homeowner,
		Price:     rh.Price,
		PhotoURL:  rh.PhotoURL,
	}
}

func RestHousesToApp(rh *RestHouses) []*house.House {
	var houses []*house.House
	for _, h := range rh.Houses {
		house := RestHouseToApp(h)
		houses = append(houses, house)
	}

	return houses
}

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

	var houses RestHouses

	jsonErr := json.Unmarshal(housesJSON, &houses)

	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse houses data: %v", jsonErr)
	}

	return RestHousesToApp(&houses), nil
}

