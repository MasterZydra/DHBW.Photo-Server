package image

import (
	"testing"
)

func TestNewImage(t *testing.T) {
	img := NewImage("img1", "2020-11-20", "d41d8cd98f00b204e9800998ecf8427e")

	if img.Name != "img1" ||
		img.hash != "d41d8cd98f00b204e9800998ecf8427e" ||
		img.Date.Format("2006-01-02") != "2020-11-20" {
		t.Errorf("Image object not filled correctly")
	}
}

func TestNewImage_Invalid(t *testing.T) {
	// ToDo Implement
}