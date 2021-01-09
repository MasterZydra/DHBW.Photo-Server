package image

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/cryptography"
	"DHBW.Photo-Server/internal/util"
	"bytes"
	"image/jpeg"
	"path"
	"time"
)

// The UploadImage struct is used to represent an uploaded image.
// It uses an Image struct to store the name, creation date and hash values.
// Additionally it contains the raw data of the image and the path of the user
// directory where it will be stored.
type UploadImage struct {
	Image
	Raw      []byte
	userPath string
}

// Create a new UploadImage object from the given image name, creation date and
// the raw data of the image itself. It calculates the additional data like the
// hash value.
func NewUploadImage(name string, creationDate time.Time, raw []byte) UploadImage {
	// Generate hash value
	hash := cryptography.HashToString(cryptography.HashImage(raw))

	// Extract creation date from EXIF data
	img := NewImage(name, creationDate, hash)
	exifData, err := parseRawExifDataFromFile(bytes.NewReader(raw))
	if err == nil {
		exifDate, err := time.Parse(dhbwphotoserver.TimeLayout, string(getDateFromData(exifData)))
		if err == nil {
			img = NewImage(name, exifDate, hash)
		}
	}

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
	// Check if folder exists and if necessary create folder
	err := util.CheckExistAndCreateFolder(path.Join(dhbwphotoserver.ImageDir(), i.userPath))
	if err != nil {
		return err
	}

	// Write image to disk
	return util.WriteRawImage(path.Join(dhbwphotoserver.ImageDir(), i.userPath, i.Name), i.Raw)
}

func (i *UploadImage) GenerateAndSaveThumbnailToDisk() error {
	// Check if folder exists and if necessary create folder
	err := util.CheckExistAndCreateFolder(path.Join(dhbwphotoserver.ImageDir(), i.userPath, thumbdir))
	if err != nil {
		return err
	}

	// Create image.Image object
	original, err := jpeg.Decode(bytes.NewReader(i.Raw))
	if err != nil {
		return err
	}

	thumbnail := GenerateThumbnail(original, 200)

	// Encode thumbnail to get byte slice and store it
	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf,thumbnail, nil)
	if err != nil {
		return err
	}

	// Write thumbnail to disk
	return util.WriteRawImage(path.Join(dhbwphotoserver.ImageDir(), i.userPath, thumbdir, i.Name), imageBuf.Bytes())
}
