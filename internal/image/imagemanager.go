/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

import (
	"math"
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
func NewImageManager(userName string) (*ImageManager, error) {
	return ReadContent(userName)
}

// Check if the ImageManager already contains the given uploaded image.
// It compares the file names (case-insensitive) and the hash values of the
// image.
func (im *ImageManager) Contains(image *UploadImage) bool {
	for _, i := range im.images {
		if strings.Compare(i.hash, image.hash) == 0 ||
			strings.ToLower(i.Name) == strings.ToLower(image.Name) {
			return true
		}
	}
	return false
}

// Add an UploadImage to this ImageManager.
// The UploadImage will be stored in the associated user directory.
// The image will also be added to the Image array and stored in the content
// file.
func (im *ImageManager) AddImageUpload(image *UploadImage) error {
	// Set path to user directory
	image.SetUserPath(im.user)
	// Store file to disk
	err := image.SaveImageToDisk()
	if err != nil {
		return err
	}
	// Store thumbnail to disk
	err = image.GenerateAndSaveThumbnailToDisk()
	if err != nil {
		return err
	}
	// Add file to image array
	im.AddImage(&image.Image)
	// Sort and store content file
	im.sort()
	err = WriteContent(im.user, im)
	if err != nil {
		return err
	}
	return nil
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

// Return the pointer to Image object which has the given name.
func (im *ImageManager) GetImage(name string) *Image {
	for _, i := range im.images {
		if strings.Compare(i.Name, name) == 0 {
			return i
		}
	}
	return nil
}

// Return a pointer to an array of Image object pointer which are in the ImageManager.
// The result Images are defined by the given start index and length.
func (im *ImageManager) GetThumbnails(start, length int) []*Image {
	if start >= len(im.images) {
		// return empty array
		return []*Image{}
	}

	end := int64(math.Min(float64(start+length), float64(len(im.images))))
	images := im.images[start:end]
	return images
}

// Returns total number of Images of current ImageManager.
func (im *ImageManager) TotalImages() int {
	return len(im.images)
}
