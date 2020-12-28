package image

import (
	"io/ioutil"
	"os"
	"testing"
)

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

func TestImageManager_AddImageUpload(t *testing.T) {
	// Test data
	fileName := "MyImg.jpg"
	// Read bytes of example image
	raw, err := ioutil.ReadFile("../../test/example_imgs/img1.jpg")
	if err != nil {
		t.Errorf("Error reading image: %v", err)
	}
	// Init upload image
	upimg := UploadImage{
		Raw: raw,
		Image: Image{
			Name: fileName,
			Date: "01.01.2020",
			Hash: "d41d8cd98f00b204e9800998ecf8427e"}}

	// Add image to ImageManager
	imgMan := ImageManager{user: "../../test/output"}
	imgMan.AddImageUpload(&upimg)

	// Check if image is stored to directory
	dir, err := os.Open("../../test/output")
	if err != nil {
		t.Errorf("Failed to open output folder: %v", err)
		return
	}
	// List all files and folders in output folder
	fileInfo, err := dir.Readdir(0)
	if err != nil {
		t.Errorf("Failed to read folder content: %v", err)
		return
	}
	// Search in array if image is contained
	found := false
	for _, info := range fileInfo {
		if info.Name() == fileName {
			found = true
			break
		}
	}
	// Check if image was found
	if !found {
		t.Errorf("File %v not found in output folder", fileName)
	}

	// Check if image is in image array in ImageManager
	if len(imgMan.images) != 1 {
		t.Errorf("Too much images in ImageManager")
		return
	}
	if img := imgMan.images[0];
	img.Name != upimg.Name || img.Date != upimg.Date || img.Hash != upimg.Hash {
		t.Errorf("Image in ImageManager does not match with UploadImage")
	}

	// ToDo Load Image and check if raw data matches too
}

// ToDo Test AddImage
// ToDo Test Sort
// ToDo Test GetPhoto
// ToDo Test GetThumbnail