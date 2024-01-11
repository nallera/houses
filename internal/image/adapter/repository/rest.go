package repository

import (
	"fmt"
	"houses/internal/image"
	"houses/server"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var MaxRetries = 5
var downloadFolder = "images"

type RestRepositoryClient struct {
	restClient server.RestClient
}

func NewRestRepositoryClient(restClient server.RestClient) *RestRepositoryClient {
	return &RestRepositoryClient{
		restClient: restClient,
	}
}

type Result struct {
	err error
	Id  int
}

func (r *RestRepositoryClient) GetImages(imagesMetadata []*image.Metadata) error {
	numberOfImages := len(imagesMetadata)
	errChan := make(chan Result, numberOfImages)

	var wg sync.WaitGroup

	for _, im := range imagesMetadata {
		wg.Add(1)
		imageMetadata := im

		go func(errChan chan Result) {
			defer wg.Done()

			for retry := 1; retry <= MaxRetries; retry++ {
				err := downloadImage(imageMetadata)

				// if there's an error, retry as long as possible
				if err != nil {
					println(fmt.Sprintf("%+v", err))
					if retry == MaxRetries {
						// return the error to communicate the process failed
						errChan <- Result{
							err: err,
							Id:  imageMetadata.Id,
						}
					}
					continue
				}

				errChan <- Result{
					err: nil,
					Id:  imageMetadata.Id,
				}
				println(fmt.Sprintf("Successfully saved image for house %d", imageMetadata.Id))
				break
			}
		}(errChan)
	}

	for in := 0; in < numberOfImages; in++ {
		result := <-errChan

		if result.err != nil {
			return fmt.Errorf("error trying to download the image %d: %v", result.Id, result.err)
		}
	}

	close(errChan)

	return nil
}

func downloadImage(metadata *image.Metadata) error {
	urlSplit := strings.Split(metadata.Url, "/")
	fileName := urlSplit[len(urlSplit)-1]

	downloadDir := "./" + downloadFolder

	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		createErr := os.Mkdir(downloadFolder, os.ModePerm)
		if createErr != nil {
			return fmt.Errorf("error creating the download directory %s: %v", downloadFolder, createErr)
		}
	}

	out, err := os.Create(downloadFolder + "/" + fileName)
	if err != nil {
		return fmt.Errorf("error creating the file %s: %v", fileName, err)
	}
	defer out.Close()

	resp, err := http.Get(metadata.Url)
	if err != nil {
		return fmt.Errorf("error getting the image %s: %v", metadata.Url, err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing the file %s: %v", fileName, err)
	}

	return nil
}
