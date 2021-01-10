/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/image"
	"testing"
	"time"
)

func TestUser_AddOrderEntryAlreadyInList(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.OrderList.AddOrderEntry(img)
	err := usr.OrderList.AddOrderEntry(img)
	if err == nil || err.Error() != "Image 'img1' already in order list" {
		t.Error("No or wrong error while adding order entry that should already exist")
	}
}

func TestUser_AddOrderEntry(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	err := usr.OrderList.AddOrderEntry(img)
	if err != nil || usr.OrderList.Entries[0].Image.Name != "img1" {
		t.Error("Error or wrong image while adding image as a order entry to user")
	}
}

func TestUser_RemoveOrderEntrySuccessful(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.OrderList.AddOrderEntry(img)
	ok := usr.OrderList.RemoveOrderEntry(img.Name)
	if !ok {
		t.Error("Image img1 should be removable from orderList")
	}
}

func TestUser_RemoveOrderEntryFail(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	ok := usr.OrderList.RemoveOrderEntry(img.Name)
	if ok {
		t.Error("Image img1 should not be removable from orderList")
	}
}

func TestUser_GetOrderEntryFound(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.OrderList.AddOrderEntry(img)
	index, entry := usr.OrderList.GetOrderEntry(img.Name)
	if index != 0 || entry.Image.Name != img.Name {
		t.Error("orderList entry index wrong or wrong image")
	}
}

func TestUser_GetOrderEntryNotFound(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	index, entry := usr.OrderList.GetOrderEntry(img.Name)
	if index != -1 || entry != nil {
		t.Error("orderList entry index wrong or entry actually found")
	}
}
