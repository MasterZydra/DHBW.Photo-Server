package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/image"
	"errors"
)

// The OrderList holds the ListEntries of the order list and provides functions to manage them
type OrderList struct {
	Entries []*OrderListEntry
}

// Checks first if this image is already in the orderList
// After that it adds a new order entry with the passed image object pointer
func (ol *OrderList) AddOrderEntry(image *image.Image) error {
	_, entry := ol.GetOrderEntry(image.Name)
	if entry != nil {
		return errors.New("Image '" + image.Name + "' already in order list")
	}
	newEntry := OrderListEntry{
		Image:          image,
		Format:         DHBW_Photo_Server.OrderListFormats[1], // Letter (8.5 x 11)
		NumberOfPrints: 1,
	}
	ol.Entries = append(ol.Entries, &newEntry)
	return nil
}

// Gets the order entry with the passed imageName string and removes this entry from the orderList preserving order
func (ol *OrderList) RemoveOrderEntry(imageName string) bool {
	i, entry := ol.GetOrderEntry(imageName)
	if entry != nil {
		ol.Entries = append(ol.Entries[:i], ol.Entries[i+1:]...)
		return true
	}
	return false
}

// Gets the order entry with the image that has the name of the passed imageName
// Returns it's index in the orderList and the a pointer to the order.OrderListEntry
func (ol *OrderList) GetOrderEntry(imageName string) (index int, entry *OrderListEntry) {
	for i, entry := range ol.Entries {
		if entry.Image.Name == imageName {
			return i, entry
		}
	}
	return -1, nil
}
