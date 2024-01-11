package house

type Repository interface {
	GetHouses(perPage, pageNumber int) ([]*House, error)
}
