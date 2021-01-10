/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package image

import "image/color"

// The Pixel struct is used to represent one Pixel of an image.
// It stores the red, green, blue and alpha value of that position.
type Pixel struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

// Create a new Pixel object.
// The constructor can be used to directly create a pixel from the result of
// jpegImage.At(x,y).RGBA() by using it as input parameters for the function.
func NewPixel(r, g, b, a uint32) Pixel {
	return Pixel{R: r, G: g, B: b, A: a}
}

// Return a color.RGBA object.
// Convert the unsigned 32 bit integers into unsigned 8 bit integers.
func (p Pixel) GetColorRGBA() color.RGBA {
	return color.RGBA{R: uint8(p.R / 256), G: uint8(p.G / 256), B: uint8(p.B / 256), A: uint8(p.A / 256)}
}
