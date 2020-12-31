package image

import (
	"DHBW.Photo-Server/internal/cryptography"
	"DHBW.Photo-Server/internal/util"
	"os"
	"path"
)

// The UploadImage struct is used to represent an uploaded image.
// It uses an Image struct to store the name, creation date and hash values.
// Additionally it contains the raw data of the image and the path of the user
// directory where it will be stored.
type UploadImage struct {
	Image
	Raw []byte
	userPath string
}

// Create a new UploadImage object from the given image name, creation date and
// the raw data of the image itself. It calculates the additional data like the
// hash value.
func NewUploadImage(name string, creationDate string, raw []byte) UploadImage {
	hash := cryptography.HashToString(cryptography.HashImage(&raw))
	// ToDo Read Exif -> if no date use given creationdate
	img := NewImage(name, creationDate, hash)
	return UploadImage{Image: *img, Raw: raw}
}

// Set or change user path. The path defines where the uploaded image will be
// saved on the server.
func (i *UploadImage) SetUserPath(path string) {
	i.userPath = path
}

// Save the raw data of the image as a file in the user directory with the
// given file name.
func (i *UploadImage) SaveImageToDisk() error {
	err := util.CheckExistAndCreateFolder(path.Join(imagedir, i.userPath))
	if err != nil {
		return err
	}
	// Open new file
	imgFile, err := os.Create(path.Join(imagedir, i.userPath, i.Name))
	if err != nil {
		return err
	}

	// Write data
	imgFile.Write(i.Raw)

	// Close file
	imgFile.Close()
	if err != nil {
		return err
	}
	return nil
}