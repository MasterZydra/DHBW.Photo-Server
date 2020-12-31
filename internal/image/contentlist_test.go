package image

import (
	"testing"
)

func TestReadContent(t *testing.T) {
	user := "../../test"
	// Test images
	image1 := NewImage("img1", "2020-11-20", "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("img2", "2020-11-21", "d41d8cdb8f0db204a9800498ecf8427e")

	// Overwrite output file name
	imagedir = ""
	usercontent = "contentReadTest.csv"

	// Read content file
	readImages := ReadContent(user)
	if readImages == nil || len(readImages.images) != 2 || readImages.user != user {
		t.Errorf("File was not read correctly")
		return
	}

	// Check if read content is correct
	if img := readImages.images; img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].Hash != image1.Hash ||
		img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].Hash != image2.Hash {
		t.Errorf("Read content is not correct")
	}
}

func TestReadContent_FileNotExist(t *testing.T) {
	user := "../../test"
	// Overwrite output file name
	imagedir = ""
	usercontent = "contentWhichDoesNotExist.csv"

	// Read content file
	readImages := ReadContent(user)
	if readImages == nil || len(readImages.images) != 0 || readImages.user != user {
		t.Errorf("File was not read correctly")
		return
	}
}

func TestWriteContent(t *testing.T) {
	// Test images
	image1 := NewImage("Image 1", "2020-01-02", "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("Image 2", "2020-01-01", "d41d8cd98f10b214e5803998ecf8427e")

	// Overwrite output file name
	imagedir = ""
	usercontent = "contentWriteTest.csv"

	// Creat new image manager, fill it with images and write that data
	imgman := ImageManager{}
	imgman.AddImage(image1)
	imgman.AddImage(image2)
	WriteContent("../../test/output", &imgman)

	// Read file again
	readImages := ReadContent("../../test/output")
	if readImages == nil {
		t.Errorf("Written file could not be read")
		return
	}

	// Check if read content is correct
	if img := readImages.images; img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].Hash != image1.Hash ||
		img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].Hash != image2.Hash {
		t.Errorf("Content is not what was written")
	}
}
