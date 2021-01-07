package order

import (
	"DHBW.Photo-Server"
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: Jones: tests
// TODO: Jones: Documentation

type Metadata struct {
	Path           string
	Format         string
	NumberOfPrints int
}

func CreateOrderListZipFile(zipFileName string, username string, orderList []*ListEntry) error {
	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zw := zip.NewWriter(newZipFile)
	defer zw.Close()

	var metadata []Metadata

	usersImageRoot := filepath.Join(DHBW_Photo_Server.ImageDir(), username)
	// for each order list entry add the corresponding image file to zip
	for _, entry := range orderList {
		imagePath := filepath.Join(usersImageRoot, entry.Image.Name)
		err = addFileToZip(zw, imagePath)
		if err != nil {
			return err
		}

		meta := Metadata{
			Path:           imagePath,
			Format:         entry.Format,
			NumberOfPrints: entry.NumberOfPrints,
		}
		metadata = append(metadata, meta)
	}

	// convert metadata to indented JSON
	jsonBytes, err := json.MarshalIndent(metadata, "", "	")
	if err != nil {
		return err
	}

	jsonFileName := filepath.Join(usersImageRoot, "content.json")
	err = ioutil.WriteFile(jsonFileName, jsonBytes, 0755)
	if err != nil {
		return err
	}

	// add content JSON file to zip
	err = addFileToZip(zw, jsonFileName)
	if err != nil {
		return err
	}

	err = os.Remove(jsonFileName)
	if err != nil {
		return err
	}

	return nil
}

func addFileToZip(zw *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(stats)
	if err != nil {
		return err
	}

	header.Name = filename

	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}
