package image

import (
	"time"
)

type Image struct {
	Name string
	Date time.Time
	Hash string
}

func NewImage(name string, date string, hash string) *Image {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		// ToDo Implement
	}
	// ToDo Implement
	return &Image{
		Name: name,
		Date: d,
		Hash: hash}
}

func (i *Image) GetThumb() *[]byte {
	// ToDo Implement
	return &[]byte{}
}

func (i *Image) GetOriginal() *[]byte {
	// ToDo Implement
	return &[]byte{}
}
