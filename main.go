package main

import (
	"fmt"
	"houses/internal/house/adapter/repository"
	"houses/internal/house/app"
	"houses/server"
)

func main() {
	// instantiate service
	restClient := server.NewRestClient()
	houseRepo := repository.NewRestRepositoryClient(restClient)
	houseService := app.NewHouseService(houseRepo)

	// request pages concurrently until success

	houses, err := houseService.GetHouses(13, 2)
	if err != nil {
		println(fmt.Sprintf("error: %+v", err))
	}

	for _, h := range houses {
		println(fmt.Sprintf("%+v", h))
	}

	// concurrently download the photos
}
