package main

import (
	"testing"
)

func ExampleReadJPEGsFromPath() {
	ReadJPEGsFromPath("./")
	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	// 
	// --------------
	// Successfully read 1 files
	// --------------
}

func TestIsJpgCheck (t *testing.T) {
	jpgFilenames := [4]string{"Test.jpg", "Image.jpeg", "Bild.jpg", "IMG_20201027_131455.jpg"}
	otherFilenames := [4]string{"Image.png", "Test.gif", "Bild.svg", "IMG_20201027_131455.png"}

	for _, file := range jpgFilenames {
		result := isJPG(file)

		if !result {
			t.Error("Error while checking filename '" + file + "'")
		}
	}

	for _, file := range otherFilenames {
		result := isJPG(file)

		if result {
			t.Error("Error while checking filename '" + file + "'")
		}
	}
}