package image

import "strings"

type ImageManager struct {
	images []*Image
}

func NewImageManager(userName string) *ImageManager {
	// ToDo Implement
	return &ImageManager{}
}

func (im *ImageManager) Contains(image *UploadImage) bool {
	for _, i := range im.images {
		if strings.Compare(i.Hash, image.Hash) == 0 {
			return true
		}
	}
	return false
}

func (im *ImageManager) AddImage(image *Image) {
	im.images = append(im.images, image)
}

func (im *ImageManager) Sort() {
	// Do something
}

func (im *ImageManager) GetPhoto(name string) *[]byte {
	for _, i := range im.images {
		if strings.Compare(i.Name, name) == 0 {
			return &[]byte{}
		}
	}
	return nil
}

func (im *ImageManager) GetThumbnail(name string) *[]byte {
	for _, i := range im.images {
		if strings.Compare(i.Name, name) == 0 {
			return &[]byte{}
		}
	}
	return nil
}