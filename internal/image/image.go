package image

type Image struct {
	Name string
	Date string
	Hash string
}

func (i *Image) GetThumb() *[]byte {
	return &[]byte{}
}

func (i *Image) GetOriginal() *[]byte {
	return &[]byte{}
}