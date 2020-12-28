package image

import (
	"testing"
)

func TestReadContent(t *testing.T) {
	// Test images
	image1 := Image{Name: "img1", Date: "20.11.2020", Hash: "d41d8cd98f00b204e9800998ecf8427e"}
	image2 := Image{Name: "img2", Date: "21.11.2020", Hash: "d41d8cdb8f0db204a9800498ecf8427e"}

	// Overwrite output file name
	usercontent = "contentReadTest.csv"

	// Read content file
	readImages := ReadContent("../../test")
	if readImages == nil {
		t.Errorf("File could not be read")
		return
	}

	// Check if read content is correct
	if img := readImages.images;
		img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].Hash != image1.Hash ||
			img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].Hash != image2.Hash {
		t.Errorf("Read content is not correct")
	}
}

func TestWriteContent(t *testing.T) {
	// Test images
	image1 := Image{Name: "Image 1", Date: "02.01.2020", Hash: "d41d8cd98f00b204e9800998ecf8427e"}
	image2 := Image{Name: "Image 2", Date: "01.01.2020", Hash: "d41d8cd98f10b214e5803998ecf8427e"}

	// Overwrite output file name
	usercontent = "contentWriteTest.csv"

	// Creat new image manager, fill it with images and write that data
	imgman := ImageManager{}
	imgman.AddImage(&image1)
	imgman.AddImage(&image2)
	WriteContent("../../test/output", &imgman)

	// Read file again
	readImages := ReadContent("../../test/output")
	if readImages == nil {
		t.Errorf("Written file could not be read")
		return
	}

	// Check if read content is correct
	if img := readImages.images;
		img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].Hash != image1.Hash ||
			img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].Hash != image2.Hash {
		t.Errorf("Content is not what was written")
	}
}