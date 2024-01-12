package test

import (
	"github.com/stretchr/testify/mock"
	"houses/internal/house"
	"houses/internal/image"
)

type RestClientMock struct {
	mock.Mock
}

func (m *RestClientMock) Get(url string) ([]byte, error) {
	args := m.Called(url)
	d, _ := args.Get(0).([]byte)
	err, _ := args.Get(1).(error)

	return d, err
}

type HouseRestRepositoryMock struct {
	mock.Mock
}

func (m *HouseRestRepositoryMock) GetHousesWithPagination(numberOfHouses, numberOfPages int) ([]*house.House, error) {
	args := m.Called(numberOfHouses, numberOfPages)
	d, _ := args.Get(0).([]*house.House)
	e, _ := args.Get(1).(error)

	return d, e
}

type ImageRestRepositoryMock struct {
	mock.Mock
}

func (m *ImageRestRepositoryMock) GetImages(imagesMetadata []*image.Metadata) error {
	args := m.Called(imagesMetadata)
	e, _ := args.Get(0).(error)

	return e
}
