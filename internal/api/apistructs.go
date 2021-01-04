package api

import (
	"DHBW.Photo-Server/internal/image"
	"time"
)

// TODO: Jones Documentation

type BaseRes interface {
	GetError() string
}

type TestReqData struct {
	SomeString  string
	SomeInteger int
}
type TestResData struct {
	Error      string
	SomeResult string
}

func (a TestResData) GetError() string {
	return a.Error
}

type RegisterReqData struct {
	Username             string
	Password             string
	PasswordConfirmation string
}
type RegisterResData struct {
	Error string
}

func (a RegisterResData) GetError() string {
	return a.Error
}

type UploadReqData struct {
	Base64Image  string
	Filename     string
	CreationDate time.Time
}
type UploadResData struct {
	Error string
}

func (a UploadResData) GetError() string {
	return a.Error
}

type ImageResData struct {
	Error string
	Image *image.Image
}

func (a ImageResData) GetError() string {
	return a.Error
}

type ThumbnailsReqData struct {
	Index  int
	Length int
}
type ThumbnailsResData struct {
	Error       string
	Images      []*image.Image
	TotalImages int
}

func (a ThumbnailsResData) GetError() string {
	return a.Error
}
