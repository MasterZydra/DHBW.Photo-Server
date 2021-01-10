/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

import (
	"DHBW.Photo-Server/internal/util"
	"bytes"
	"image/jpeg"
	"path"
	"testing"
)

func TestGenerateThumbnail(t *testing.T) {
	// Read example image
	raw, err := util.ReadRawImage(path.Join(testDir, "example_imgs/img1.jpg"))
	if err != nil {
		t.Errorf("Error reading example image: %v", err)
		return
	}

	// Create image.Image object
	original, err := jpeg.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Errorf("Unexpected error while decoding: %v", err)
	}

	// Test that if image is smaller then thumbnail width, return original
	thumbnail := GenerateThumbnail(original, 2000)
	if thumbnail != original {
		t.Errorf("Did not return correct image")
	}

	// Test that if image is bigger then thumbnail width, return thumbnail
	thumbnail = GenerateThumbnail(original, 200)
	if thumbnail == original {
		t.Errorf("Did not return correct image")
	}
}
