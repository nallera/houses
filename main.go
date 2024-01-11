package main

import (
	"fmt"
	"houses/internal/house/adapter/repository"
	"houses/internal/house/app"
	"houses/server"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// instantiate service
	restClient := server.NewRestClient()
	houseRepo := repository.NewRestRepositoryClient(restClient)
	houseService := app.NewHouseService(houseRepo)

	// request pages concurrently until success

	houses, err := houseService.GetHouses(10, 2)
	if err != nil {
		println(fmt.Sprintf("error: %+v", err))
	}

	for _, h := range houses {
		println(fmt.Sprintf("%+v", h))
	}

	// concurrently download the photos

	fileUrl := houses[0].PhotoURL
	urlSplit := strings.Split(fileUrl, "/")
	fileName := urlSplit[len(urlSplit)-1]
	out, err := os.Create(fileName)
	defer out.Close()
	resp, err := http.Get(fileUrl)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
}
