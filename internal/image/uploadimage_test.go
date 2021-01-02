package image

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestNewUploadImage(t *testing.T) {
	// Read bytes of example image
	raw, err := ioutil.ReadFile("../../test/example_imgs/img2.jpg")
	if err != nil {
		t.Errorf("Error reading image: %v", err)
		return
	}

	// Create image
	img := NewUploadImage("myImage.jpg", "2020-12-12", raw)

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
	rawdata := []byte{84,69,83,84}

	// Initialize and save new image
	upimg := UploadImage{
		Raw: rawdata,
		Image: *NewImage("imageWriteTest.txt", "2020-01-01", "d41d8cd98f00b204e9800998ecf8427e"),
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