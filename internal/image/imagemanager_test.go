/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

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

const testDir = "../../test"
const testOutputDir = "../../test/output"

func TestNewImageManager(t *testing.T) {
	// Test data
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 14:13:43")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 23:12:24")
	image1 := NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("img2", date2, "d41d8cdb8f0db204a9800498ecf8427e")

	// Overwrite output file name
	DHBW_Photo_Server.SetImageDir("")
	usercontent = "contentNewImageManagerTest.csv"

	// Init ImageManager for given user path
	imgMan, err := NewImageManager(testDir)
	if err != nil {
		t.Errorf("Something went wrong creating a new ImageManager: %v", err)
		return
	}

	if imgMan == nil {
		t.Errorf("Something went wrong creating a ImageManager from user path")
		return
	}

	// Check if given parameter is stored in object
	if imgMan.user != testDir {
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
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 14:13:43")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 23:12:24")
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
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 14:13:43")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 23:12:24")
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
	raw1, err := ioutil.ReadFile(path.Join(testDir, "example_imgs/img1.jpg"))
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}
	// Read bytes of example image 2
	raw2, err := ioutil.ReadFile(path.Join(testDir, "example_imgs/img2.jpg"))
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}

	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 23:12:24")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 14:13:43")
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
	raw, err := ioutil.ReadFile(path.Join(testDir, "example_imgs/img1.jpg"))
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}

	// Clean up before running test logic
	_ = os.Remove(path.Join(testOutputDir, "MyImg.jpg"))
	_ = os.Remove(path.Join(testOutputDir, "content.csv"))
	_ = os.RemoveAll(path.Join(testOutputDir, thumbdir))

	// Overwrite default settings
	DHBW_Photo_Server.SetImageDir("")

	// Init upload image
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-01-01 23:12:24")
	upimg := NewUploadImage(fileName, date, raw)

	// Add image to ImageManager
	imgMan := ImageManager{user: testOutputDir}
	_ = imgMan.AddImageUpload(&upimg)

	// Check if image is stored to directory
	dir, err := os.Open(testOutputDir)
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
	readRawImage, err := ioutil.ReadFile(path.Join(testOutputDir, fileName))
	if err != nil {
		t.Errorf("Error reading saved image: %v", err)
		return
	}
	if bytes.Compare(raw, readRawImage) != 0 {
		t.Errorf("Writen and read image raw data is not equal")
	}
}

func TestImageManager_AddImage(t *testing.T) {
	// Test images
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 23:12:24")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 15:41:20")
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
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 08:20:13")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-01-21 06:30:02")
	date3, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 08:20:20")

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
	// Read test image
	raw, err := util.ReadRawImage(path.Join(testDir, "example_imgs/img1.jpg"))
	if err != nil {
		t.Errorf("Could not read example image: %v", err)
		return
	}
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-12-31 03:45:12")
	upimg := NewUploadImage("img1.jpg", date, raw)

	// Overwrite output file name
	DHBW_Photo_Server.SetImageDir("")
	usercontent = "contentGetImageTest.csv"

	// Init ImageManager
	imgMan, err := NewImageManager(testOutputDir)
	if err != nil {
		t.Errorf("Something went wrong creating a new ImageManager: %v", err)
		return
	}
	_ = imgMan.AddImageUpload(&upimg)

	// Get image and compare raw data
	imgManRead, err := NewImageManager(testOutputDir)
	if err != nil {
		t.Errorf("Something went wrong creating a new ImageManager: %v", err)
		return
	}
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

func TestImageManager_TotalImages(t *testing.T) {
	// Overwrite output file name
	DHBW_Photo_Server.SetImageDir("")
	usercontent = "contentNewImageManagerTest.csv"

	// Init ImageManager for given user path
	imgMan, err := NewImageManager(testDir)
	if err != nil {
		t.Errorf("Something went wrong creating a new ImageManager: %v", err)
		return
	}
	totalImages := imgMan.TotalImages()

	if totalImages != 2 {
		t.Error("Wrong count of total images in ImageManager")
	}
}

func TestImageManager_GetThumbnails(t *testing.T) {
	// Overwrite output file name
	DHBW_Photo_Server.SetImageDir("")
	usercontent = "contentLongList.csv"

	// Init ImageManager
	imgMan, err := NewImageManager(testDir)
	if err != nil {
		t.Errorf("Something went wrong creating a new ImageManager: %v", err)
		return
	}

	// Test 0 to 4
	thumbs := imgMan.GetThumbnails(0, 5)
	if thumbs[0] != imgMan.images[0] || thumbs[1] != imgMan.images[1] || thumbs[2] != imgMan.images[2] || thumbs[3] != imgMan.images[3] || thumbs[4] != imgMan.images[4] ||
		len(thumbs) != 5 {
		t.Errorf("False thumbnails returned (1)")
	}

	// Test latest ones
	thumbs = imgMan.GetThumbnails(6, 5)
	if thumbs[0] != imgMan.images[6] || thumbs[1] != imgMan.images[7] || thumbs[2] != imgMan.images[8] ||
		len(thumbs) != 3 {
		t.Errorf("False thumbnails returned (2)")
	}

	// Test if start index is to big
	thumbs = imgMan.GetThumbnails(100, 5)
	if len(thumbs) != 0 {
		t.Errorf("Returned thumbnails even if start index is greater then the highest index")
	}
}
