/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import (
	"path/filepath"
	"testing"
)

func ExampleReadJPEGsFromPath() {
	ReadJPEGsFromPath("../../test/")
	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
}

func TestReadJPEGsFromPath(t *testing.T) {
	path := "../../test"
	jpegs := ReadJPEGsFromPath(path)
	// there should be one JPEG (Testbild.jpg)
	if len(jpegs) < 1 || filepath.Base(jpegs[0].Name()) != "Testbild.jpg" {
		t.Errorf("Error: Testbild not found")
	}
}

func TestIsJpgCheck(t *testing.T) {
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
