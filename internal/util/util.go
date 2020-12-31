package util

import (
	"io/ioutil"
)

// TODO: Tests schreiben
func ContainsString(haystack []string, needle string) bool {
	for _, currentNeedle := range haystack {
		if currentNeedle == needle {
			return true
		}
	}
	return false
}

// Returns read bytes for given file in param "path"
func ReadRawImage(path string) ([]byte, error){
	readRawImage, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return readRawImage, nil
}