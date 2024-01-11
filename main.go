package main

import (
	"fmt"
	houserepo "houses/internal/house/adapter/repository"
	houseapp "houses/internal/house/app"
	imagerepo "houses/internal/image/adapter/repository"
	"houses/internal/image/app"
	"houses/server"
)

func main() {
	// instantiate services
	houseRestClient := server.NewRestClient()
	houseRepo := houserepo.NewRestRepositoryClient(houseRestClient)
	houseService := houseapp.NewHouseService(houseRepo)

	imagesRestClient := server.NewRestClient()
	imageRepo := imagerepo.NewRestRepositoryClient(imagesRestClient)
	imagesService := app.NewImageService(imageRepo)

	// request pages concurrently until success

	houses, err := houseService.GetHouses(20, 4)
	if err != nil {
		println(fmt.Sprintf("error: %+v", err))
	}

	for _, h := range houses {
		println(fmt.Sprintf("%+v", h))
	}

	// concurrently download the photos
	err = imagesService.DownloadImages(houses)
	if err != nil {
		println(fmt.Sprintf("error: %+v", err))
	}
}
