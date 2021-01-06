package user

import (
	"DHBW.Photo-Server"
	"testing"
)

func TestImageManagerCacheNew(t *testing.T) {
	usersFile = DHBW_Photo_Server.TestUserFile
	userManager = nil
	userManagerBefore := userManager
	newUserManager := ImageManagerCache()
	if userManagerBefore != nil || newUserManager == nil {
		t.Error("Something went wrong while getting new ImageManager")
	}
}

func TestImageManagerCacheExists(t *testing.T) {
	usersFile = "someweirdfile.csv"
	userManager = NewUserManager()
	newUserManager := ImageManagerCache()
	if newUserManager.UsersFile != usersFile {
		t.Error("Something went wrong while getting existing ImageManager")
	}
}

func TestResetImageManagerCache(t *testing.T) {
	userManager = NewUserManager()
	ResetImageManagerCache()
	if userManager != nil {
		t.Error("Something went wrong while resetting the userManager cache")
	}
}
