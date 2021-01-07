package image

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/util"
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNewUploadImage(t *testing.T) {
	// Read bytes of example image
	raw, err := ioutil.ReadFile("../../test/example_imgs/img2.jpg")
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}

	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-12-12 13:12:15")

	// Create image
	img := NewUploadImage("myImage.jpg", date, raw)

	// Check struct content
	if img.Name != "myImage.jpg" ||
		img.hash != "fe430f8f373ec2abd6880bfdaebcd9ae" ||
		// ToDo Add check for correct date + Example with image without exif creation date
		//img.Date.Format("2006-01-02") != "2020-12-12" ||
		bytes.Compare(img.Raw, raw) != 0 {
		t.Errorf("UploadImage content does not match with given data")
	}
}

func TestUploadImage_SaveImageToDisk(t *testing.T) {
	// Test data
	rawdata := []byte{84, 69, 83, 84}

	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-01-01 13:12:25")
	// Initialize and save new image
	upimg := UploadImage{
		Raw:      rawdata,
		Image:    *NewImage("imageWriteTest.txt", date, "d41d8cd98f00b204e9800998ecf8427e"),
		userPath: "../../test/output"}
	upimg.SaveImageToDisk()

	// Read saved file again
	content, err := ioutil.ReadFile("../../test/output/imageWriteTest.txt")
	if err != nil {
		t.Errorf("writen file could not be opened: %v", err)
		return
	}

	// Check if data matches
	if bytes.Compare(content, rawdata) != 0 {
		t.Errorf("Read data is not equal to the written data")
	}
}

func TestUploadImage_SetUserPath(t *testing.T) {
	newPath := "NewPath"
	img := UploadImage{userPath: ""}
	img.SetUserPath(newPath)
	if img.userPath != newPath {
		t.Errorf("Path %v expected but received %v", newPath, img.userPath)
	}
}

func TestUploadImage_GenerateAndSaveThumbnailToDisk(t *testing.T) {
	// Prepare test UploadImage
	// Read example image
	raw, err := util.ReadRawImage("../../test/example_imgs/img1.jpg")
	if err != nil {
		t.Errorf("Error reading example image: %v", err)
		return
	}
	// Create new UploadImage
	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-01-01 01:23:02")
	upimg := NewUploadImage("img1.jpg", date, raw)
	upimg.userPath = "test"

	// Overwrite default pathes
	dhbwphotoserver.SetImageDir("../..")
	thumbdir = "output"

	// Clean up
	os.Remove("../../test/output/img1.jpg")

	// Generate and save thumbnail
	err = upimg.GenerateAndSaveThumbnailToDisk()
	if err != nil {
		t.Errorf("Something went wrong generating and saving the thumbnail to disk: %v", err)
	}

	// Check if image is in folder
	_, err = os.Stat("../../test/output/img1.jpg")
	if err != nil {
		t.Errorf("Error getting written image state: %v", err)
	}
}