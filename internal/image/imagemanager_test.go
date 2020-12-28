package image

import "testing"

func TestNewImageManager(t *testing.T) {
	// Test data
	user := "../../test"
	image1 := Image{Name: "img1", Date: "20.11.2020", Hash: "d41d8cd98f00b204e9800998ecf8427e"}
	image2 := Image{Name: "img2", Date: "21.11.2020", Hash: "d41d8cdb8f0db204a9800498ecf8427e"}

	// Overwrite output file name
	usercontent = "contentNewImageManagerTest.csv"

	// Init ImageManager for given user path
	imgMan := NewImageManager(user)

	if imgMan == nil {
		t.Errorf("Something went wrong creating a ImageManager from user path")
		return
	}

	// Check if given parameter is stored in object
	if imgMan.user != user {
		t.Errorf("Property user is not filled correctly")
	}

	// Check if read content is correct
	if img := imgMan.images;
		img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].Hash != image1.Hash ||
			img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].Hash != image2.Hash {
		t.Errorf("Read content is not correct")
	}
}

func TestImageManager_Contains(t *testing.T) {
	// Test images
	image1 := UploadImage{Image: Image{Name: "img1", Date: "20.11.2020", Hash: "d41d8cd98f00b204e9800998ecf8427e"}}
	image2 := UploadImage{Image: Image{Name: "img2", Date: "21.11.2020", Hash: "d41d8cdb8f0db204a9800498ecf8427e"}}

	imgMan := ImageManager{}
	imgMan.AddImage(&(image1.Image))

	if !imgMan.Contains(&image1) {
		t.Errorf("Existing image1 not detected")
	}
	if imgMan.Contains(&image2) {
		t.Errorf("Wrongly detects image2 as alread contained")
	}
}
