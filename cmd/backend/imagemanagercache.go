/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import (
	"DHBW.Photo-Server/internal/image"
	"errors"
	"strings"
	"time"
)

// This map is used to store all ImageManager objects in RAM to reduce the time for disk access.
var imageManagers = make(map[string]*image.ImageManager)

// This functions checks the map of ImageManagers for the given user.
// If it exists, it returns the pointer to the object, otherwise it loads the data from the disk.
func getImageManager(username string) (*image.ImageManager, error) {
	username = strings.ToLower(username)
	imgman, exists := imageManagers[username]
	if !exists {
		imgman, err := image.NewImageManager(username)
		if err != nil {
			return nil, err
		}
		imageManagers[username] = imgman
	}
	return imgman, nil
}

// Encapsulated logic to add an uploaded image
func UploadImage(username string, name string, creationDate time.Time, raw []byte) error {
	upimg := image.NewUploadImage(name, creationDate, raw)

	imgman, err := getImageManager(username)
	if err != nil {
		return err
	}

	if imgman.Contains(&upimg) {
		return errors.New("Image is already stored")
	}
	return imgman.AddImageUpload(&upimg)
}

// Encapsulated logic to get an image by its name for the given user
func GetImage(username, imagename string) (*image.Image, error) {
	imgman, err := getImageManager(username)
	if err != nil {
		return &image.Image{}, err
	}
	return imgman.GetImage(imagename), nil
}

// Encapsulated logic to get the thumbnails for the given range and user
func GetThumbnails(username string, start, length int) ([]*image.Image, error) {
	imgman, err := getImageManager(username)
	if err != nil {
		return []*image.Image{}, err
	}
	return imgman.GetThumbnails(start, length), nil
}

// Encapsulated logic to get the total amount of images for the given user
func GetTotalImages(username string) (int, error) {
	imgman, err := getImageManager(username)
	if err != nil {
		return 0, err
	}
	return imgman.TotalImages(), nil
}
