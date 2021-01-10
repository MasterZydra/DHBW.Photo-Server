/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"time"
)

// Define pathes and filenames
var thumbdir = dhbwphotoserver.ThumbDir
var usercontent = dhbwphotoserver.UserContent

// The Image struct is used to represent one entry in the content file.
// It only stores the name, creation date and hash value of the raw data.
type Image struct {
	Name string
	Date time.Time
	hash string
}

// Create a new Image object.
func NewImage(name string, date time.Time, hash string) *Image {
	return &Image{Name: name, Date: date, hash: hash}
}

// returns the Date field in a specific format
func (image *Image) FormattedDate() string {
	return image.Date.Format(time.RFC822)
}
