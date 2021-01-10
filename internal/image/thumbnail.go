/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

import (
	"image"
	"math"
)

// Generate a thumbnail for a given Image with a given width. The function
// returns an Image that has an approximated width as the given one,
// based on the chosen implementation of down scaling.
// This function can only down scale. If the given Images width is smaller or
// equal than the given thumbnail width, the original image will be returned.
func GenerateThumbnail(original image.Image, width int) image.Image {
	// Get width and height of original image
	imgWidth := original.Bounds().Max.X
	imgHeight := original.Bounds().Max.Y

	// Check if image size is already smaller or equal than the thumbnail width
	if imgWidth <= width {
		return original
	}

	// Calculate the step size for the selected pixels
	// Using round to get better results than just using an integer division
	stepSize := int(math.Round(float64(imgWidth) / float64(width)))

	// Create a new empty RGBA image with reduced size
	thumbnail := image.NewRGBA(image.Rect(0, 0, imgWidth/stepSize, imgHeight/stepSize))
	// Iterate through every pixel of the new image and assign the responding pixel of the original image
	for y := 0; ; y++ {
		// Calculate yPos once for all x in that row
		yPos := stepSize * y
		if yPos >= imgHeight {
			break
		}
		for x := 0; stepSize*x < imgWidth; x++ {
			// Use Pixel struct to convert the result of the RGBA() func into a color.RGBA object
			thumbnail.SetRGBA(x, y, NewPixel(original.At(stepSize*x, yPos).RGBA()).GetColorRGBA())
		}
	}
	return thumbnail
}
