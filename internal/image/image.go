package image

import "DHBW.Photo-Server/internal/cryptography"

type Image struct {
	Name string
	Date string
	Hash string
}

type UploadImage struct {
	Image
	Raw []byte
}

func NewUploadImage(name string, raw []byte) UploadImage {
	hash := cryptography.HashToString(cryptography.HashImage(&raw))
	return UploadImage{Image: Image{Name: name, Hash: hash}, Raw: raw}
}

func NewImage(name string, date string) *Image {
	// ToDo - Implement
	return &Image{} //Name: name, Date: date} //, Hash: cryptography.HashImage()}
}

func (i *Image) GetThumb() *[]byte {
	return &[]byte{}
}

func (i *Image) GetOriginal() *[]byte {
	return &[]byte{}
}
