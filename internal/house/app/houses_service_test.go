package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"houses/internal/house"
	"houses/test"
	"testing"
)

func TestHouseService_GetHouses(t *testing.T) {
	type depFields struct {
		houseRepoMock *test.HouseRestRepositoryMock
	}
	type input struct {
		numberOfHouses int
		numberOfPages  int
	}
	type output struct {
		houses []*house.House
		err    error
	}

	numberOfPages := 2
	numberOfHouses := 4

	housesResponse := test.MakeHousesResponse()

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "return houses successfully",
			in: input{
				numberOfHouses: numberOfHouses,
				numberOfPages:  numberOfPages,
			},
			on: func(df *depFields) {
				df.houseRepoMock.On("GetHousesWithPagination", numberOfHouses, numberOfPages).Return(housesResponse, nil).Once()
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, housesResponse, out.houses)
			},
		},
		{
			name: "error getting houses",
			in: input{
				numberOfHouses: numberOfHouses,
				numberOfPages:  numberOfPages,
			},
			on: func(df *depFields) {
				df.houseRepoMock.On("GetHousesWithPagination", numberOfHouses, numberOfPages).Return(nil, errors.New("test error")).Once()
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error getting houses: test error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			housesRepoMock := new(test.HouseRestRepositoryMock)
			s := NewHouseService(housesRepoMock)

			f := &depFields{houseRepoMock: housesRepoMock}

			tt.on(f)

			// When
			houses, err := s.GetHouses(tt.in.numberOfHouses, tt.in.numberOfPages)

			o := output{houses: houses, err: err}
			// Then
			tt.assert(t, &o)
			housesRepoMock.AssertExpectations(t)
		})
	}
}
