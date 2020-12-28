package image

import (
	"DHBW.Photo-Server/internal/cryptography"
	"os"
	"path"
)

type UploadImage struct {
	Image
	Raw []byte
	userPath string
}

func (i *UploadImage) SaveImageToDisk() error {
	// Open new file
	imgFile, err := os.Create(path.Join(i.userPath, i.Name))
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