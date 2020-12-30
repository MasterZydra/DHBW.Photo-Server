package api

import (
	"DHBW.Photo-Server/internal/image"
)

type BaseRes interface {
	GetError() string
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
	Error string
}

func (a ThumbnailsRes) GetError() string {
	return a.Error
}
