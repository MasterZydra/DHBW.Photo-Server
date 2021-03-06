/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"testing"
	"time"
)

func TestNewImage(t *testing.T) {
	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-11-20 15:38:12")
	img := NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")

	if img.Name != "img1" ||
		img.hash != "d41d8cd98f00b204e9800998ecf8427e" ||
		img.Date != date {
		t.Errorf("Image object not filled correctly")
	}
}

func TestImage_FormattedDate(t *testing.T) {
	expected := "21 Nov 20 13:14 UTC"
	date, _ := time.Parse(dhbwphotoserver.TimeLayout, "2020-11-21 13:14:15")
	img := NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	formattedDate := img.FormattedDate()

	if formattedDate != expected {
		t.Errorf("Expected '%v' but received '%v'", expected, formattedDate)
	}
}

func TestImage_FormattedDate_Invalid(t *testing.T) {
	expected := "01 Jan 01 00:00 UTC"
	date, err := time.Parse(dhbwphotoserver.TimeLayout, "2020-11-20")
	if err == nil {
		t.Errorf("Parsing should fail, but did not")
	}
	img := NewImage("", date, "")
	if img.FormattedDate() != expected {
		t.Errorf("Expected '%v' but received '%v'", expected, img.FormattedDate())
	}
}
