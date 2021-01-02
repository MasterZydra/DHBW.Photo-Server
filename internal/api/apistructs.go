package api

import (
	"DHBW.Photo-Server/internal/image"
)

// TODO: Jones Documentation

type BaseRes interface {
	GetError() string
}

type TestReq struct {
	SomeString  string
	SomeInteger int
}
type TestRes struct {
	Error      string
	SomeResult string
}

func (a TestRes) GetError() string {
	return a.Error
}

type RegisterReq struct {
	Username             string
	Password             string
	PasswordConfirmation string
}
type RegisterRes struct {
	Error string
}

func (a RegisterRes) GetError() string {
	return a.Error
}

type ImageUploadRes struct {
	Error string
}

func (a ImageUploadRes) GetError() string {
	return a.Error
}

type ImageReq struct {
}
type ImageRes struct {
	Error string
	Image *image.Image
}

func (a ImageRes) GetError() string {
	return a.Error
}

type ThumbnailsReq struct {
	Index  int
	Length int
}
type ThumbnailsRes struct {
	Error  string
	Images []*image.Image
}

func (a ThumbnailsRes) GetError() string {
	return a.Error
}
