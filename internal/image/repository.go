package image

type Repository interface {
	GetImages(imageMetadata []*Metadata) error
}
