package image

import (
	"time"
)

// The Image struct is used to represent one entry in the content file.
// It only stores the name, creation date and hash value of the raw data.
type Image struct {
	Name string
	Date time.Time
	Hash string
}

// Create a new Image object.
// The date string will be converted into Time.
func NewImage(name string, date string, hash string) *Image {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		// ToDo Implement Error Handling
	}
	return &Image{Name: name, Date: d, Hash: hash}
}

