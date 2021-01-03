package image

import (
	"time"
)

// The Image struct is used to represent one entry in the content file.
// It only stores the name, creation date and hash value of the raw data.
type Image struct {
	Name string
	Date time.Time
	hash string
}

// Create a new Image object.
// The date string will be converted into Time.
func NewImage(name string, date time.Time, hash string) *Image {
	return &Image{Name: name, Date: date, hash: hash}
}
