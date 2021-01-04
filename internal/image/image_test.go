package image

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"testing"
	"time"
)

func TestNewImage(t *testing.T) {
	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-11-20")
	img := NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")

	if img.Name != "img1" ||
		img.hash != "d41d8cd98f00b204e9800998ecf8427e" ||
		img.Date != date {
		t.Errorf("Image object not filled correctly")
	}
}

func TestNewImage_Invalid(t *testing.T) {
	// ToDo Implement
}

func TestImage_FormattedDate(t *testing.T) {
	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-11-20")
	img := NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	formattedDate := img.FormattedDate()

	if formattedDate != "20 Nov 20 00:00 UTC" {
		t.Error("Something went wrong during FormattedDate")
	}
}
