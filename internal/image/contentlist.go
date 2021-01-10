/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

/*
 * This file contains functions to read and write the content file.
 * The content file contains the list of all images in the user folder.
 * It contains the name, creation date and hash.
 */

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/util"
	"encoding/csv"
	"os"
	"path"
	"time"
)

// Read content file for given user. The user has to be equal to the folder name.
// It returns an initialized ImageManager struct.
func ReadContent(user string) (*ImageManager, error) {
	// Open file
	csvFile, err := os.Open(path.Join(dhbwphotoserver.ImageDir(), user, usercontent))
	if os.IsNotExist(err) {
		return &ImageManager{user: user}, nil
	}
	if err != nil {
		return nil, err
	}

	// Initialize reader
	reader := csv.NewReader(csvFile)
	reader.Comma = '|'

	// Read file
	images, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Store all data in ImageManager struct
	imageManager := &ImageManager{user: user}
	for _, img := range images {
		date, err := time.Parse(dhbwphotoserver.TimeLayout, img[1])
		if err != nil {
			return nil, err
		}
		imageManager.AddImage(NewImage(img[0], date, img[2]))
	}

	return imageManager, nil
}

// Write content file for given user. The user has to be equal to the folder name.
func WriteContent(user string, imgs *ImageManager) error {
	err := util.CheckExistAndCreateFolder(path.Join(dhbwphotoserver.ImageDir(), user))
	if err != nil {
		return err
	}

	// Create new file
	f, err := os.Create(path.Join(dhbwphotoserver.ImageDir(), user, usercontent))
	if err != nil {
		return err
	}

	// Initialize writer
	writer := csv.NewWriter(f)
	writer.Comma = '|'

	// Build array that will be stored in CSV format
	var imgArray [][]string
	for _, img := range imgs.images {
		imgArray = append(imgArray, []string{img.Name, img.Date.Format(dhbwphotoserver.TimeLayout), img.hash})
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
