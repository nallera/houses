package main

import (
	"flag"
	"fmt"
	houserepo "houses/internal/house/adapter/repository"
	houseapp "houses/internal/house/app"
	imagerepo "houses/internal/image/adapter/repository"
	"houses/internal/image/app"
	"houses/server"
)

func main() {
	// parse the args to the program
	numberOfHouses := flag.Int("houses", 10, "the total number of houses to retrieve")
	numberOfPages := flag.Int("pages", 2, "the number of pages to retrieve")
	flag.Parse()

	if *numberOfPages > *numberOfHouses {
		panic("the number of houses must be greater than the number of pages")
	}

	println(fmt.Sprintf("nh: %d, np: %d", *numberOfHouses, *numberOfPages))

	// instantiate services
	houseRestClient := server.NewRestClient()
	houseRepo := houserepo.NewRestRepository(houseRestClient)
	houseService := houseapp.NewHouseService(houseRepo)

	imagesRestClient := server.NewRestClient()
	imageRepo := imagerepo.NewRestRepository(imagesRestClient)
	imagesService := app.NewImageService(imageRepo)

	// request pages concurrently until success

	houses, err := houseService.GetHouses(*numberOfHouses, *numberOfPages)
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
