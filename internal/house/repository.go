package house

type Repository interface {
	GetHousesWithPagination(numberOfHouses, numberOfPages int) ([]*House, error)
}
