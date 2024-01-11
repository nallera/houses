package app

import (
	"fmt"
	"houses/internal/house"
	"math"
)

type HouseService interface {
	GetHouses(numberOfHouses, numberOfPages int) ([]*house.House, error)
}

func NewHouseService(houseRestRepository house.Repository) HouseService {
	return &houseService{
		HouseRestRepository: houseRestRepository,
	}
}

type houseService struct {
	HouseRestRepository house.Repository
}

func (hs *houseService) GetHouses(numberOfHouses, numberOfPages int) ([]*house.House, error) {
	perPage := int(math.Ceil(float64(numberOfHouses) / float64(numberOfPages)))
	println(fmt.Sprintf("per page %d", perPage))

	houses, err := hs.HouseRestRepository.GetHousesWithPagination(perPage, numberOfPages)

	if err != nil {
		return nil, fmt.Errorf("error getting houses: %v", err)
	}

	return houses, nil
}
