package app

import (
	"fmt"
	"houses/internal/house"
	"houses/internal/image"
)

type ImageService interface {
	DownloadImages([]*house.House) error
}

func NewImageService(imageRestRepository image.Repository) ImageService {
	return &imageService{
		ImageRestRepository: imageRestRepository,
	}
}

type imageService struct {
	ImageRestRepository image.Repository
}

func housesToImagesMetadata(houses []*house.House) []*image.Metadata {
	var imagesMetadata []*image.Metadata

	for _, h := range houses {
		imagesMetadata = append(imagesMetadata, &image.Metadata{
			Url:     h.PhotoURL,
			Id:      h.Id,
			Address: h.Address,
		})
	}

	return imagesMetadata
}

func (is *imageService) DownloadImages(houses []*house.House) error {
	err := is.ImageRestRepository.GetImages(housesToImagesMetadata(houses))

	if err != nil {
		return fmt.Errorf("error saving images: %v", err)
	}

	return nil
}
