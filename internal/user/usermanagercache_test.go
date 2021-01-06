package user

import (
	"DHBW.Photo-Server"
	"testing"
)

func TestUserManagerCacheNew(t *testing.T) {
	usersFile = DHBW_Photo_Server.TestUserFile
	userManager = nil
	userManagerBefore := userManager
	newUserManager := UserManagerCache()
	if userManagerBefore != nil || newUserManager == nil {
		t.Error("Something went wrong while getting new UserManager")
	}
}

func TestUserManagerCacheExists(t *testing.T) {
	usersFile = "someweirdfile.csv"
	userManager = NewUserManager()
	newUserManager := UserManagerCache()
	if newUserManager.UsersFile != usersFile {
		t.Error("Something went wrong while getting existing UserManager")
	}
}

func TestResetUserManagerCache(t *testing.T) {
	userManager = NewUserManager()
	ResetUserManagerCache()
	if userManager != nil {
		t.Error("Something went wrong while resetting the userManager cache")
	}
}
