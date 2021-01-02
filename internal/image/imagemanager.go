package image

import (
	"DHBW.Photo-Server/internal/util"
	"fmt"
	"path"
	"sort"
	"strings"
)

// The ImageManager struct is used to represent all uploaded images of one user.
// It contains the user as a string and an array of all images represented as
// Image objects.
type ImageManager struct {
	images []*Image
	user   string
}

// Create a new ImageManager object.
// The user name is needed to load an existing content file for the given user.
// If no user directory exists or no content file exists for this user it
// returns an empty ImageManager only filled with the user name.
func NewImageManager(userName string) *ImageManager {
	return ReadContent(userName)
}

// Check if the ImageManager already contains the given uploaded image.
// It compares the file names (case-insensitive) and the hash values of the
// image.
func (im *ImageManager) Contains(image *UploadImage) bool {
	for _, i := range im.images {
		if strings.Compare(i.hash, image.hash) == 0 ||
			strings.ToLower(i.Name) == strings.ToLower(image.Name){
			return true
		}
	}
	return false
}

// Add an UploadImage to this ImageManager.
// The UploadImage will be stored in the associated user directory.
// The image will also be added to the Image array and stored in the content
// file.
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

// Add an Image object to the Image array.
// Before calling this function, the function Contains should be called to make
// sure that the image is not already in the list.
func (im *ImageManager) AddImage(image *Image) {
	// Add image
	im.images = append(im.images, image)
}

// Sort the Image array by date descending
func (im *ImageManager) sort() {
	sort.Slice(im.images, func(i, j int) bool {
		return im.images[i].Date.After(im.images[j].Date)
	})
}

// Get raw data of given image name for current user.
// Returns pointer to an LoadedImage object.
// Returns nil if file was not found.
func (im *ImageManager) GetImage(name string) *LoadedImage {
	return im.getRawImage(path.Join(imagedir,im.user), name)
}

// Get raw data of thumbnail of given image name for current user.
// Returns pointer to an LoadedImage object.
// Returns nil if file was not found.
func (im *ImageManager) GetThumbnail(name string) *LoadedImage {
	return im.getRawImage(path.Join(imagedir,im.user,thumbdir), name)
}

// Get raw data of given image name in given directory.
// Returns pointer to an LoadedImage object.
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

