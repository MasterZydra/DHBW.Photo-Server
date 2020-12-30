package image

import (
	"fmt"
	"sort"
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
	// Sort and store content file
	im.sort()
	// ToDo In own func
	err = WriteContent(im.user, im)
	if err != nil {
		fmt.Errorf("error saving content file: %v", err)
	}
}

func (im *ImageManager) AddImage(image *Image) {
	// Add image
	im.images = append(im.images, image)
}

func (im *ImageManager) sort() {
	sort.Slice(im.images, func(i, j int) bool {
		return im.images[i].Date.After(im.images[j].Date)
	})
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