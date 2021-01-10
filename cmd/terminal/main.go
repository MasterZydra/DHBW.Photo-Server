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
	UploadPhotos(input, filePointers)

	// ask if another folder should be uploaded
	var uploadAnother bool = UploadAnotherFolder()

	for uploadAnother {
		// read new path
		input.Path = EnterNewPath()

		// read files from the path the user has entered
		filePointers = ReadJPEGsFromPath(input.Path)

		// send https request
		UploadPhotos(input, filePointers)

		// ask if another folder should be uploaded
		uploadAnother = UploadAnotherFolder()
	}
}