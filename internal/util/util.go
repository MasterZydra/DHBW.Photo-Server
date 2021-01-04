package util

import (
	"io/ioutil"
	"os"
)

// Returns read bytes for given file in param "path"
func ReadRawImage(path string) ([]byte, error) {
	readRawImage, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return readRawImage, nil
}

// Write raw bytes into file with given name
func WriteRawImage(name string, image []byte) error {
	// Open new file
	imgFile, err := os.Create(name)
	if err != nil {
		return err
	}

	// Write data
	imgFile.Write(image)

	// Close file
	imgFile.Close()
	if err != nil {
		return err
	}
	return nil
}

// Check if given path exists.
// If false, then not create the necessary folders.
func CheckExistAndCreateFolder(path string) error {
	file, err := os.Stat(path)
	if os.IsNotExist(err) || !file.Mode().IsDir() {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
