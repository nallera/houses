package app

import (
	"fmt"
	"houses/internal/house"
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
	houses, err := hs.HouseRestRepository.GetHousesWithPagination(numberOfHouses, numberOfPages)

	if err != nil {
		return nil, fmt.Errorf("error getting houses: %v", err)
	}

	return houses, nil
}
