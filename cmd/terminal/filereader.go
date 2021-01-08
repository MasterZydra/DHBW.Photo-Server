package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadJPEGsFromPath(path string) []*os.File {
	// try reading files from path
	// if the path is not valid an error will occur
	files, err := ioutil.ReadDir(path)
	if err != nil {
		WriteMessage("Something went wrong while reading the folder. Try it again")
		log.Fatal(err)
	}

	WriteMessage("Folder successfully read in")

	// fill an array with files to return it
	var result []*os.File
	var fileSeparator string = string(os.PathSeparator)

	// for each file in the directory the user has entered ...
	for _, fileInfo := range files {
		// ... it is checked if it is a JPEG image
		if !fileInfo.IsDir() && isJPG(fileInfo.Name()) {
			// get file-pointer
			file, err := os.Open(path + fileSeparator + fileInfo.Name())
			if err != nil {
				log.Fatal(err)
			}

			// and add the pointer to the array
			result = append(result, file)
		}
	}

	// Information for the user how many images have been read in
	WriteMessage("Successfully read " + strconv.Itoa(len(result)) + " files")

	return result
}

func isJPG(filename string) bool {
	// check if file is jpg
	filenameSplitted := strings.Split(filename, ".")
	filenameExtension := strings.ToLower(filenameSplitted[len(filenameSplitted)-1])
	return strings.EqualFold(filenameExtension, "jpg") || strings.EqualFold(filenameExtension, "jpeg")
}