package util

import (
	"bytes"
	"io/ioutil"
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

func TestReadRawImage_FileNotExist(t *testing.T) {
	content, err := ReadRawImage("../../test/notExistingFile.abc")

	// Check result
	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist but get other error: %v", err)
	}
	if content == nil || bytes.Compare(content, []byte{}) != 0 {
		t.Errorf("Did not get empty byte slice: %v", err)
	}
}

func TestWriteRawImage(t *testing.T) {
	// Test data
	filename := "../../test/output/utilWriteRawTest.txt"
	rawdata := []byte{84, 69, 83, 84}

	// Delete file
	os.Remove(filename)

	// Write data
	err := WriteRawImage(filename, rawdata)
	if err != nil {
		t.Errorf("Error while writing raw data: %v", err)
		return
	}

	// Read saved file again
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("writen file could not be opened: %v", err)
		return
	}

	// Check if data matches
	if bytes.Compare(content, rawdata) != 0 {
		t.Errorf("Read data is not equal to the written data")
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
