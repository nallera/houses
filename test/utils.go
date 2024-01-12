package test

import (
	"houses/internal/house"
)

func MakeHousesResponse() []*house.House {
	return []*house.House{
		{
			Id:        0,
			Address:   "test address",
			Homeowner: "test homeowner",
			Price:     100,
			PhotoURL:  "http://testurl.com",
		},
		{
			Id:        1,
			Address:   "test address",
			Homeowner: "test homeowner",
			Price:     200,
			PhotoURL:  "http://testurl.com",
		},
		{
			Id:        2,
			Address:   "test address",
			Homeowner: "test homeowner",
			Price:     2500,
			PhotoURL:  "http://testurl.com",
		},
		{
			Id:        3,
			Address:   "test address",
			Homeowner: "test homeowner",
			Price:     4500,
			PhotoURL:  "http://testurl.com",
		},
	}
}
