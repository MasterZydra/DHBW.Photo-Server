package image

import (
	"fmt"
	"strings"
)

type ImageManager struct {
	images []*Image
	user string
}

func NewImageManager(userName string) *ImageManager {
	return ReadContent(userName)
}

func (im *ImageManager) Contains(image *UploadImage) bool {
	for _, i := range im.images {
		if strings.Compare(i.Hash, image.Hash) == 0 {
			return true
		}
	}
	return false
}

func (im *ImageManager) AddImageUpload(image *UploadImage) {
	// Set path to user directory
	image.SetUserPath(im.user)
	// Store file to disk
	err := image.SaveImageToDisk()
	if err != nil {
		fmt.Errorf("error saving image to disk: %v", err)
	}
	// Add file to image array
	im.AddImage(&image.Image)
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