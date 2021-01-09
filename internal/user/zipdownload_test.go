package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/image"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateOrderListZipFile(t *testing.T) {
	usr := NewUser("max", "12345678")
	date1, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	date2, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img1 := image.NewImage("img1.jpg", date1, "d41d8cd98f00b204e9800998ecf8427e")
	img2 := image.NewImage("img2.jpg", date2, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.OrderList.AddOrderEntry(img1)
	_ = usr.OrderList.AddOrderEntry(img2)

	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	zipFileName := filepath.Join(DHBW_Photo_Server.ImageDir(), usr.Name, "download.zip")
	err := CreateOrderListZipFile(zipFileName, usr.Name, usr.OrderList)
	defer os.Remove(zipFileName)

	_, statErr := os.Stat(zipFileName)

	if err != nil || os.IsNotExist(statErr) {
		t.Error("Something went wrong while creating zip file")
	}
}

func TestCreateOrderListZipFileInvalidImage(t *testing.T) {
	usr := NewUser("max", "12345678")
	date, _ := time.Parse(DHBW_Photo_Server.TimeLayout, "2020-11-21 08:35:59")
	img := image.NewImage("invalidfile.jpg", date, "d41d8cd98f00b204e9800998ecf8427e")
	_ = usr.OrderList.AddOrderEntry(img)

	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	zipFileName := filepath.Join(DHBW_Photo_Server.ImageDir(), usr.Name, "download.zip")
	err := CreateOrderListZipFile(zipFileName, usr.Name, usr.OrderList)
	defer os.Remove(zipFileName)

	if err == nil || !os.IsNotExist(err) {
		t.Error("File invalidfile.jpg should not exists")
	}
}
