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
var DownloadFolder = "images"
var NumberOfWorkers = 5

type RestRepository struct {
	restClient server.RestClient
}

func NewRestRepository(restClient server.RestClient) *RestRepository {
	return &RestRepository{
		restClient: restClient,
	}
}

type Result struct {
	Id  int
	err error
}

func (r *RestRepository) GetImages(imagesMetadata []*image.Metadata) error {
	numberOfImages := len(imagesMetadata)

	downloadDir := "./" + DownloadFolder

	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		createErr := os.Mkdir(DownloadFolder, os.ModePerm)
		if createErr != nil {
			return fmt.Errorf("error creating the download directory %s: %v", DownloadFolder, createErr)
		}
	}

	var wg sync.WaitGroup
	jobs := make(chan struct{}, NumberOfWorkers)
	errChan := make(chan Result, numberOfImages)

	for _, im := range imagesMetadata {
		imageMetadata := im
		wg.Add(1)
		// block the calls until the jobs queue has a free slot
		jobs <- struct{}{}

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
							Id:  imageMetadata.Id,
							err: err,
						}
						<-jobs
						return
					}
					continue
				}
			}

			errChan <- Result{
				Id:  imageMetadata.Id,
				err: nil,
			}
			<-jobs
		}(errChan)
	}

	wg.Wait()

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
	urlSplit := strings.Split(metadata.Url, ".")
	extension := urlSplit[len(urlSplit)-1]
	fileName := fmt.Sprintf("%d-%s.%s", metadata.Id, metadata.Address, extension)
	fileNameWithPath := DownloadFolder + "/" + fileName

	if _, err := os.Stat(fileNameWithPath); os.IsNotExist(err) {
		resp, err := http.Get(metadata.Url)
		if err != nil {
			return fmt.Errorf("error getting the image %s: %v", metadata.Url, err)
		}
		defer resp.Body.Close()

		out, err := os.Create(fileNameWithPath)
		if err != nil {
			return fmt.Errorf("error creating the file %s: %v", fileName, err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("error writing the file %s: %v", fileName, err)
		}

		println(fmt.Sprintf("Successfully saved image for house %d", metadata.Id))

		return nil
	}

	println(fmt.Sprintf("Image for house %d already existed in the download folder", metadata.Id))

	return nil
}
