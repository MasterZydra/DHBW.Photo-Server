package image

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"encoding/csv"
	"os"
	"path"
)
import "fmt"

var imagedir 	= DHBW_Photo_Server.ImageDir
var thumbdir 	= DHBW_Photo_Server.ThumbDir
var usercontent = DHBW_Photo_Server.UserContent

// Read "content.csv" for given user. The user has to be equal to the folder name.
// It returns an initialized ImageManager struct.
func ReadContent(user string) *ImageManager {
	// Open file
	csvFile, err := os.Open(path.Join(user, usercontent))
	if err != nil {
		// Error stuff
		fmt.Println(err)
	}

	// Initialize reader
	reader := csv.NewReader(csvFile)
	reader.Comma = '|'

	// Read file
	images, err := reader.ReadAll()
	if err != nil {
		// Error stuff
		fmt.Println(err)
	}

	// Store all data in ImageManager struct
	imageManager := &ImageManager{user: user}
	for _, img := range images {
		imageManager.AddImage(NewImage(img[0], img[1], img[2]))
	}

	return imageManager
}

// Write "content.csv" for given user. The user has to be equal to the folder name.
func WriteContent(user string, imgs *ImageManager) error {
	// Create new file
	f, err := os.Create(path.Join(user, usercontent))
	if err != nil {
		return err
	}

	// Initialize writer
	writer := csv.NewWriter(f)
	writer.Comma = '|'

	// Build array that will be stored in CSV format
	var imgArray [][]string
	for _, img := range imgs.images {
		imgArray = append(imgArray, []string{img.Name, img.Date.Format("2006-01-02"), img.Hash})
	}

	// Write file
	err = writer.WriteAll(imgArray)
	if err != nil {
		return err
	}
	writer.Flush()

	// Close file
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
