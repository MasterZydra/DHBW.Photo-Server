package image

import "time"

type LoadedImage struct {
	Raw		[]byte
	Name	string
	Date	time.Time
}

func NewLoadedImage(image *Image, raw *[]byte) LoadedImage{
	return LoadedImage{Name: image.Name, Date: image.Date, Raw: *raw}
}