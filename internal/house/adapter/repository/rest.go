package repository

import (
	"encoding/json"
	"fmt"
	"houses/internal/house"
	"houses/server"
	"math"
	"sync"
)

var HouseBaseURL = "http://app-homevision-staging.herokuapp.com/api_project/houses"
var MaxRetries = 5
var NumberOfWorkers = 5

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

type RestRepository struct {
	restClient server.RestClient
}

func NewRestRepository(restClient server.RestClient) *RestRepository {
	return &RestRepository{
		restClient: restClient,
	}
}

type HousePageResult struct {
	houses     []*house.House
	pageNumber int
}

func (r *RestRepository) GetHousesWithPagination(numberOfHouses, numberOfPages int) ([]*house.House, error) {
	perPage := int(math.Ceil(float64(numberOfHouses) / float64(numberOfPages)))

	var houses []*house.House

	var wg sync.WaitGroup
	houseChan := make(chan HousePageResult, numberOfPages)
	jobs := make(chan struct{}, NumberOfWorkers)

	for pn := 1; pn <= numberOfPages; pn++ {
		pageNumber := pn
		wg.Add(1)
		// block the calls until the jobs queue has a free slot
		jobs <- struct{}{}

		go func(houseChan chan HousePageResult) {
			defer wg.Done()

			var (
				auxHouses []*house.House
				err       error
			)

			for retry := 1; retry <= MaxRetries; retry++ {
				auxHouses, err = r.getSingleHousesPage(perPage, pageNumber)

				// if there's an error, retry as long as possible
				if err != nil {
					println(fmt.Sprintf("%+v", err))
					if retry == MaxRetries {
						// return a code to understand the process failed
						houseChan <- HousePageResult{
							houses:     nil,
							pageNumber: pageNumber,
						}
						<-jobs
						return
					}
					continue
				}

				break
			}

			houseChan <- HousePageResult{
				houses:     auxHouses,
				pageNumber: pageNumber,
			}
			println(fmt.Sprintf("Successfully got houses page %d", pageNumber))
			<-jobs
		}(houseChan)
	}

	wg.Wait()

	results := make([][]*house.House, numberOfPages)
	for pn := 1; pn <= numberOfPages; pn++ {
		houseResult := <-houseChan

		if houseResult.houses == nil {
			return nil, fmt.Errorf("maximum number of retries exceeded while trying to get page %d", houseResult.pageNumber)
		}

		// the goal of this is to be able to sort the houses by id
		results[houseResult.pageNumber-1] = houseResult.houses
	}

	close(houseChan)

	for _, result := range results {
		houses = append(houses, result...)
	}

	return houses[:numberOfHouses], nil
}

func (r *RestRepository) getSingleHousesPage(perPage, pageNumber int) ([]*house.House, error) {
	url := fmt.Sprintf("%s?page=%d&per_page=%d", HouseBaseURL, pageNumber, perPage)

	println(fmt.Sprintf("request to %s", url))

	housesJSON, err := r.restClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get houses for page %d: %v", pageNumber, err)
	}

	var houses RestHouses

	jsonErr := json.Unmarshal(housesJSON, &houses)

	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse houses data for page %d: %v", pageNumber, jsonErr)
	}

	return RestHousesToApp(&houses), nil
}
