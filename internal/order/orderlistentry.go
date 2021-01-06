package order

import "DHBW.Photo-Server/internal/image"

// TODO: jones: Documentation

type ListEntry struct {
	Image          *image.Image
	Format         string
	NumberOfPrints int
}
