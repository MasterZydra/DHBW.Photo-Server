package main

import (
	"strconv"
)

type UserInput struct {
	Username string
	Password string
	Host string
	Path string
}

func main() {
	// Welcome user
	Greet()

	// read user input
	input := WaitForUserInput()

	// read files from the path the user has entered
	filePointers := ReadJPEGsFromPath(input.Path)

	// encode files to base 64
	encoded := EncodeFilesToBase64(filePointers)

	// output anything that uses encoded so that the variable is used and no error occurs
	WriteMessage("Encoded " + strconv.Itoa(len(encoded)) + " files")
	// send https request
	// use encoded as body
	// TODO: add restconsumer

	UploadPhotos(input, encoded)
}