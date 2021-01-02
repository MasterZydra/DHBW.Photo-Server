package util

import (
	"io/ioutil"
	"os"
)

// TODO: Jones Documentation

// TODO: Jones Tests schreiben, wenns überhaupt noch benötigt wird
func ContainsString(haystack []string, needle string) bool {
	for _, currentNeedle := range haystack {
		if currentNeedle == needle {
			return true
		}
	}
	return false
}

// Returns read bytes for given file in param "path"
func ReadRawImage(path string) ([]byte, error) {
	readRawImage, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return readRawImage, nil
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
