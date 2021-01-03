package main

import (
	"DHBW.Photo-Server/internal/image"
	"errors"
	"time"
)

var imageManagers = make(map[string]*image.ImageManager)

func getImageManager(user string) *image.ImageManager {
	imgman, exists := imageManagers[user]
	if !exists {
		imgman = image.NewImageManager(user)
		imageManagers[user] = imgman
	}
	return imgman
}

func UploadImage(user string, name string, creationDate time.Time, raw []byte) error {
	upimg := image.NewUploadImage(name, creationDate, raw)

	imgman := getImageManager(user)

	if imgman.Contains(&upimg) {
		return errors.New("Image is already stored")
	}
	imgman.AddImageUpload(&upimg)
	return nil
}

func GetImage(username, imagename string) *image.Image {
	imgman := getImageManager(username)
	return imgman.GetImage(imagename)
}

func GetThumbnail(username string, start, length int) []*image.Image {
	imgman := getImageManager(username)
	return *imgman.GetThumbnails(start, length)
}
