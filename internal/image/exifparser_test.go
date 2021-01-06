package image

import (
	"os"
	"testing"
)

func TestReadDate(t *testing.T) {
	filename := "../../test/Testbild.jpg"
	file, err := os.Open(filename)
	if err != nil {
		t.Errorf("Error reading image: %v", err)
	}

	exifData, err := parseRawExifDataFromFile(file)
	if err != nil {
		t.Error(err)
	}

	date := getDateFromData(exifData)

	wantedDate := "2020:08:20 14:23:45"

	if string(date) != wantedDate {
		t.Errorf("Error parsing date from Exif-data. Wanted: %v, Got: %v", wantedDate, string(date))
	}
}