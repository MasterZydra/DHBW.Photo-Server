package util

import (
	"bytes"
	"os"
	"testing"
)

func TestReadRawImage(t *testing.T) {
	readcontent, err := ReadRawImage("../../test/readRawImageTest.txt")
	if err != nil || bytes.Compare(readcontent, []byte{65,66,67,68,69,70,71,72}) != 0 {
		t.Errorf("Read content is not correct")
	}

	// Test with invalid file
	readInvalid, err := ReadRawImage("invalidImage.jpg")
	if !os.IsNotExist(err) || bytes.Compare(readInvalid, []byte{}) != 0 {
		t.Errorf("Wrong result for invalid image name")
	}
}

func TestCheckExistAndCreateFolder(t *testing.T) {
	path := "../../test/output/newFolder"
	// Clean up before test
	os.Remove(path)
	// Create folder
	err := CheckExistAndCreateFolder(path)
	if err != nil {
		t.Errorf("Something went wrong creating folder: %v", err)
	}
	// Check result
	file, err := os.Stat(path)
	if err != nil || !file.Mode().IsDir() {
		t.Errorf("Folder not created correctly: %v", err)
	}
}
