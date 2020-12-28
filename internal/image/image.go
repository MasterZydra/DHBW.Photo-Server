package image

type Image struct {
	Name string
	Date string
	Hash string
}

func NewImage(name string, date string) *Image {
	// ToDo Implement
	return &Image{} //Name: name, Date: date} //, Hash: cryptography.HashImage()}
}

func (i *Image) GetThumb() *[]byte {
	// ToDo Implement
	return &[]byte{}
}

func (i *Image) GetOriginal() *[]byte {
	// ToDo Implement
	return &[]byte{}
}
