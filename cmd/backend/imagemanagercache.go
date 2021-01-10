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

var imageManagers = make(map[string]*image.ImageManager)

func getImageManager(username string) *image.ImageManager {
	username = strings.ToLower(username)
	imgman, exists := imageManagers[username]
	if !exists {
		imgman = image.NewImageManager(username)
		imageManagers[username] = imgman
	}
	return imgman
}

func UploadImage(username string, name string, creationDate time.Time, raw []byte) error {
	upimg := image.NewUploadImage(name, creationDate, raw)

	imgman := getImageManager(username)

	if imgman.Contains(&upimg) {
		return errors.New("Image is already stored")
	}
	return imgman.AddImageUpload(&upimg)
}

func GetImage(username, imagename string) *image.Image {
	imgman := getImageManager(username)
	return imgman.GetImage(imagename)
}

func GetThumbnails(username string, start, length int) []*image.Image {
	imgman := getImageManager(username)
	return imgman.GetThumbnails(start, length)
}

func GetTotalImages(username string) int {
	imgman := getImageManager(username)
	return imgman.TotalImages()
}
