package user

import (
	"DHBW.Photo-Server"
	"testing"
)

func TestGetImageManagerNew(t *testing.T) {
	usersFile = DHBW_Photo_Server.TestUserFile
	userManager = nil
	userManagerBefore := userManager
	newUserManager := GetImageManager()
	if userManagerBefore != nil || newUserManager == nil {
		t.Error("Something went wrong while getting new ImageManager")
	}
}

func TestGetImageManagerExists(t *testing.T) {
	usersFile = "someweirdfile.csv"
	userManager = NewUserManager()
	newUserManager := GetImageManager()
	if newUserManager.UsersFile != usersFile {
		t.Error("Something went wrong while getting existing ImageManager")
	}
}
