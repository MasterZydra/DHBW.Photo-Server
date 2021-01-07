package order

import "DHBW.Photo-Server/internal/image"

// ListEntry is used in the user.User object as the orderList.
// It is used to represent one entry in that list.
// It connects the entry with the associated image, it's desired print format and the number of prints.

type ListEntry struct {
	Image          *image.Image
	Format         string
	NumberOfPrints int
}
