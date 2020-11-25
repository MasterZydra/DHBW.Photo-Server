package image

// import "../cryptography"

type Image struct {
	Name string
	Date string
	Hash string
}

func NewImage(name string, date string) *Image {
	// ToDo - Implement
	return &Image{}//Name: name, Date: date} //, Hash: cryptography.HashImage()}
}

func (i *Image) GetThumb() *[]byte {
	return &[]byte{}
}

func (i *Image) GetOriginal() *[]byte {
	return &[]byte{}
}