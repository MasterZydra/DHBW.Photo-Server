package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/image"
	"encoding/hex"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	name := "Helmut"
	user := NewUser(name, "12345678")
	if user.Name != name || user.password == "" || len(user.OrderList) > 0 {
		t.Errorf("Creating user %v didn't work", name)
	}
}

func TestNewUserPwError(t *testing.T) {
	user := NewUser("test", "12345678")
	_, hexConversionError := hex.DecodeString(user.password)
	if hexConversionError != nil {
		t.Error("Creating user didn't work cuz the password is not a valid hex value")
	}
}

func TestUser_ComparePassword(t *testing.T) {
	pw := "harryPotter4eva"
	user := NewUser("ron", pw)
	ok, err := user.ComparePassword(pw)
	if !ok || err != nil {
		t.Error("Error while comparing correct password")
	}
}

func TestUser_AddOrderEntryAlreadyInList(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.AddOrderEntry(img)
	err := usr.AddOrderEntry(img)
	if err == nil || err.Error() != "Image 'img1' already in order list" {
		t.Error("No or wrong error while adding order entry that should already exist")
	}
}

func TestUser_AddOrderEntry(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	err := usr.AddOrderEntry(img)
	if err != nil || usr.OrderList[0].Image.Name != "img1" {
		t.Error("Error or wrong image while adding image as a order entry to user")
	}
}

func TestUser_RemoveOrderEntrySuccessful(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.AddOrderEntry(img)
	ok := usr.RemoveOrderEntry(img.Name)
	if !ok {
		t.Error("Image img1 should be removable from orderList")
	}
}

func TestUser_RemoveOrderEntryFail(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	ok := usr.RemoveOrderEntry(img.Name)
	if ok {
		t.Error("Image img1 should not be removable from orderList")
	}
}

func TestUser_GetOrderEntryFound(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.AddOrderEntry(img)
	index, entry := usr.GetOrderEntry(img.Name)
	if index != 0 || entry.Image.Name != img.Name {
		t.Error("orderList entry index wrong or wrong image")
	}
}

func TestUser_GetOrderEntryNotFound(t *testing.T) {
	usr := NewUser("test", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("img1", date, "d41d8cd98f00b204e9800998ecf8427e")
	index, entry := usr.GetOrderEntry(img.Name)
	if index != -1 || entry != nil {
		t.Error("orderList entry index wrong or entry actually found")
	}
}
