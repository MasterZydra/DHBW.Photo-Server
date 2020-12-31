package image

import (
	"DHBW.Photo-Server/internal/util"
	"fmt"
	"path"
	"sort"
	"strings"
)

type ImageManager struct {
	images []*Image
	user   string
}

func NewImageManager(userName string) *ImageManager {
	return ReadContent(userName)
}

func (im *ImageManager) Contains(image *UploadImage) bool {
	for _, i := range im.images {
		if strings.Compare(i.Hash, image.Hash) == 0 ||
			strings.ToLower(i.Name) == strings.ToLower(image.Name){
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

// Get raw data of given image name for current user.
// Returns pointer to the byte slice.
// Returns nil if file was not found.
func (im *ImageManager) GetImage(name string) *LoadedImage {
	return im.getRawImage(path.Join(imagedir,im.user), name)
}

// Get raw data of thumbnail of given image name for current user.
// Returns pointer to the byte slice.
// Returns nil if file was not found.
func (im *ImageManager) GetThumbnail(name string) *LoadedImage {
	return im.getRawImage(path.Join(imagedir,im.user,thumbdir), name)
}

// Get raw data of given image name in given directory.
// Returns pointer to the byte slice.
// Returns nil if file was not found.
func (im *ImageManager) getRawImage(imgPath, name string) *LoadedImage {
	for _, i := range im.images {
		if strings.Compare(i.Name, name) == 0 {
			raw, err := util.ReadRawImage(path.Join(imgPath,name))
			if err != nil {
				fmt.Printf("error reading file: %v\n",err)
				return nil
			}
			loadedImage := NewLoadedImage(i, &raw)
			return &loadedImage
		}
	}
	return nil
}

