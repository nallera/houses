package house

type Repository interface {
	GetHousesWithPagination(perPage, numberOfPages int) ([]*House, error)
}
