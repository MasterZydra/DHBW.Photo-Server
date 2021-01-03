package image

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/util"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestNewImageManager(t *testing.T) {
	// Test data
	user := "../../test"
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21")
	image1 := NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("img2", date2, "d41d8cdb8f0db204a9800498ecf8427e")

	// Overwrite output file name
	imagedir = ""
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
	if img := imgMan.images; img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].hash != image1.hash ||
		img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].hash != image2.hash {
		t.Errorf("Read content is not correct")
	}
}

func TestImageManager_Contains_Hash(t *testing.T) {
	// Test images
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21")
	image1 := UploadImage{Image: *NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")}
	image2 := UploadImage{Image: *NewImage("img2", date2, "d41d8cdb8f0db204a9800498ecf8427e")}

	// Init ImageManager
	imgMan := ImageManager{}
	imgMan.AddImage(&(image1.Image))

	// Execute tests
	if !imgMan.Contains(&image1) {
		t.Errorf("Existing image1 not detected")
	}
	if imgMan.Contains(&image2) {
		t.Errorf("Wrongly detects image2 as alread contained")
	}
}

func TestImageManager_Contains_Filename(t *testing.T) {
	// Test images
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21")
	image1 := UploadImage{Image: *NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")}
	image2 := UploadImage{Image: *NewImage("img1", date2, "d41d8cdb8f0db204a9800498ecf8427e")}
	image3 := UploadImage{Image: *NewImage("img2", date2, "d41d8cdb8f0db204a9800498ecf8427e")}

	// Init ImageManager
	imgMan := ImageManager{}
	imgMan.AddImage(&(image1.Image))

	// Execute tests
	if !imgMan.Contains(&image1) {
		t.Errorf("Existing image1 not detected")
	}
	if !imgMan.Contains(&image2) {
		t.Errorf("Existing image1 not detected")
	}
	if imgMan.Contains(&image3) {
		t.Errorf("Wrongly detects image3 as alread contained")
	}
}

func TestImageManager_Contains_WithExampleImages(t *testing.T) {
	// Test images
	// Read bytes of example image 1
	raw1, err := ioutil.ReadFile("../../test/example_imgs/img1.jpg")
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}
	// Read bytes of example image 2
	raw2, err := ioutil.ReadFile("../../test/example_imgs/img2.jpg")
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}

	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21")
	image1 := NewUploadImage("img1", date1, raw1)
	image2 := NewUploadImage("img2", date2, raw2)

	// Init ImageManager
	imgMan := ImageManager{}
	imgMan.AddImage(&(image1.Image))

	// Execute tests
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
		return
	}
	// Init upload image

	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-01-01")
	upimg := NewUploadImage(fileName, date, raw)

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
	if img := imgMan.images[0]; img.Name != upimg.Name || img.Date != upimg.Date || img.hash != upimg.hash {
		t.Errorf("Image in ImageManager does not match with UploadImage")
	}

	// Read saved image to check if raw data is equal
	readRawImage, err := ioutil.ReadFile(path.Join("../../test/output", fileName))
	if err != nil {
		t.Errorf("Error reading saved image: %v", err)
		return
	}
	if bytes.Compare(raw, readRawImage) != 0 {
		t.Errorf("Writen and read image raw data is not equal")
	}

	// ToDo Test if image is in content file
}

func TestImageManager_AddImage(t *testing.T) {
	// Test images

	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21")
	image1 := NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("img2", date2, "d41d8cdb8f0db204a9800498ecf8427e")

	// Init ImageManager
	imgMan := ImageManager{}

	// Add first image
	imgMan.AddImage(image1)
	// Check if image is in image array in ImageManager
	if len(imgMan.images) != 1 {
		t.Errorf("Too much images in ImageManager")
		return
	}
	if img := imgMan.images[0]; img.Name != image1.Name || img.Date != image1.Date || img.hash != image1.hash {
		t.Errorf("Image in ImageManager does not match with UploadImage")
		return
	}

	// Add second image
	imgMan.AddImage(image2)
	// Check if image is in image array in ImageManager
	if len(imgMan.images) != 2 {
		t.Errorf("Too much images in ImageManager")
		return
	}
	if img := imgMan.images[1]; img.Name != image2.Name || img.Date != image2.Date || img.hash != image2.hash {
		t.Errorf("Image in ImageManager does not match with UploadImage")
	}
}

func TestImageManager_Sort(t *testing.T) {
	// Test images
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-01-21")
	date3, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21")

	image1 := NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")
	image0 := NewImage("img0", date2, "d41d8cdb8f0db204a9800498ecf8427e")
	image2 := NewImage("img2", date3, "d41d8cdb8f0db204a9800498ecf8427e")

	// Init ImageManager
	imgMan := ImageManager{}
	imgMan.AddImage(image1)
	imgMan.AddImage(image0)
	imgMan.AddImage(image2)

	// Sort and check order
	imgMan.sort()
	if img := imgMan.images; img[0].Name != "img2" || img[1].Name != "img1" || img[2].Name != "img0" {
		t.Errorf("Images are in the wrong order")
	}
}

func TestImageManager_GetImage(t *testing.T) {
	// Test data
	user := "../../test/output"
	// Read test image
	raw, err := util.ReadRawImage("../../test/example_imgs/img1.jpg")
	if err != nil {
		t.Errorf("Could not read example image: %v", err)
		return
	}
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-12-31")
	upimg := NewUploadImage("img1.jpg", date, raw)

	// Overwrite output file name
	imagedir = ""
	usercontent = "contentGetImageTest.csv"

	// Init ImageManager
	imgMan := NewImageManager(user)
	imgMan.AddImageUpload(&upimg)

	// Get image and compare raw data
	imgManRead := NewImageManager(user)
	readimg := imgManRead.GetImage("img1.jpg")
	if readimg.Name != "img1.jpg" || readimg.hash != upimg.hash {
		t.Errorf("Image from GetImage is not equal to stored image")
	}

	// Test with invalid file
	readInvalid := imgManRead.GetImage("invalidImage.jpg")
	if readInvalid != nil {
		t.Errorf("Wrong result for invalid image name")
	}
}

// ToDo Test GetImage
// ToDo Test GetThumbnail
