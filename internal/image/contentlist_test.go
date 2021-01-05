package image

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"os"
	"path"
	"testing"
	"time"
)

func TestReadContent(t *testing.T) {
	user := "../../test"
	// Test images
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-20 06:30:02")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	image1 := NewImage("img1", date1, "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("img2", date2, "d41d8cdb8f0db204a9800498ecf8427e")

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
	if img := readImages.images; img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].hash != image1.hash ||
		img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].hash != image2.hash {
		t.Errorf("Read content is not correct")
	}
}

func TestReadContent_UserFolderNotExist(t *testing.T) {
	user := "../../test/someNoneExistingUser"
	// Overwrite output file name
	imagedir = ""

	// Read content file
	readImages := ReadContent(user)
	if readImages == nil || len(readImages.images) != 0 || readImages.user != user {
		t.Errorf("File was not read correctly")
		return
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
	user := "../../test/output"
	// Test images
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-01-02 13:14:15")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-01-01 13:30:16")
	image1 := NewImage("Image 1", date1, "d41d8cd98f00b204e9800998ecf8427e")
	image2 := NewImage("Image 2", date2, "d41d8cd98f10b214e5803998ecf8427e")

	// Overwrite output file name
	imagedir = ""
	usercontent = "contentWriteTest.csv"

	// Clean up
	os.Remove(path.Join(user, usercontent))

	// Creat new image manager, fill it with images and write that data
	imgman := ImageManager{}
	imgman.AddImage(image1)
	imgman.AddImage(image2)
	WriteContent(user, &imgman)

	// Read file again
	readImages := ReadContent(user)
	if readImages == nil {
		t.Errorf("Written file could not be read")
		return
	}

	// Check if read content is correct
	if img := readImages.images; img[0].Name != image1.Name || img[0].Date != image1.Date || img[0].hash != image1.hash ||
		img[1].Name != image2.Name || img[1].Date != image2.Date || img[1].hash != image2.hash {
		t.Errorf("Content is not what was written")
	}
}
