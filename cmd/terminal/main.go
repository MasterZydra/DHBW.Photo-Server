package main

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

	// send https request
	// use encoded as body
	// TODO: add restconsumer

	UploadPhotos(input, filePointers)
}