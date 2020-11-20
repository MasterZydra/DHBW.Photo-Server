package main

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadJPEGsFromPath(path string) []*os.File {
	// TODO: validate Path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		WriteMessage("Something went wrong while reading the folder. Try it again")
		log.Fatal(err)
	}

	WriteMessage("Folder successfully read in.")

	var result []*os.File
	var fileSeparator string = string(os.PathSeparator)

	for _, fileInfo := range files {
		if !fileInfo.IsDir() && isJPG(fileInfo.Name()) {
			// get file-pointer
			file, err := os.Open(path + fileSeparator + fileInfo.Name())
			if err != nil {
				log.Fatal(err)
			}

			result = append(result, file)
		}
	}

	WriteMessage("Successfully read " + strconv.Itoa(len(result)) + " files")

	return result
}

func EncodeFilesToBase64(files []*os.File) []string {
	var encodedStrings []string

	for _, file := range files {
		// encode each file and add it to slice
		encodedStrings = append(encodedStrings, encodeSingleFile(file))
	}

	return encodedStrings
}

func encodeSingleFile(file *os.File) string {
	//first read content of the file:
	content, _ := ioutil.ReadAll(bufio.NewReader(file))

	// then encode
	encoded := base64.StdEncoding.EncodeToString(content)

	// return encoded string
	return encoded
}

func isJPG(filename string) bool {
	// check if file is jpg
	filenameSplitted := strings.Split(filename, ".")
	filenameExtension := strings.ToLower(filenameSplitted[len(filenameSplitted)-1])
	return strings.EqualFold(filenameExtension, "jpg") || strings.EqualFold(filenameExtension, "jpeg")
}