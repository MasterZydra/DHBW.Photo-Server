/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

type UserInput struct {
	Username string
	Password string
	Host     string
	Path     string
}

func main() {
	// Welcome user
	Greet()

	// read user input
	input := WaitForUserInput()

	// read files from the path the user has entered
	filePointers := ReadJPEGsFromPath(input.Path)

	// send https request
	UploadPhotos(input, filePointers)
}
