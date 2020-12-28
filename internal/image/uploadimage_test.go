package image

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestUploadImage_SaveImageToDisk(t *testing.T) {
	// Test data
	rawdata := []byte{84,69,83,84}

	// Initialize and save new image
	upimg := UploadImage{
		Raw: rawdata,
		Image: Image{
			Name: "imageWriteTest.txt",
			Date: "01.01.2020",
			Hash: "d41d8cd98f00b204e9800998ecf8427e"},
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
