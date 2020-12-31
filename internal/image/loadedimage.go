package image

import "time"

// The LoadedImage struct is used to represent an image or thumbnail which has
// been load from the disk. It will be used to send as response to a request.
type LoadedImage struct {
	Raw		[]byte
	Name	string
	Date	time.Time
}

// Create a new LoadedImage object. With the combination of an Image instance
// and the raw image data all necessary data for the initialization is given.
func NewLoadedImage(image *Image, raw *[]byte) LoadedImage{
	return LoadedImage{Name: image.Name, Date: image.Date, Raw: *raw}
}