package main

import (
	"DHBW.Photo-Server/internal/image"
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

func UploadImage(user, name, creationDate string, raw []byte) string {
	upimg := image.NewUploadImage(name, creationDate, raw)

	imgman := getImageManager(user)

	if imgman.Contains(&upimg) {
		return "Image is already stored"
	}
	imgman.AddImageUpload(&upimg)
	return ""
}

func GetImage(username, imagename string) *image.Image {
	imgman := getImageManager(username)
	return imgman.GetImage(imagename)
}

func GetThumbnail(username string, start, length int) []*image.Image {
	imgman := getImageManager(username)
	return *imgman.GetThumbnails(start, length)
}
