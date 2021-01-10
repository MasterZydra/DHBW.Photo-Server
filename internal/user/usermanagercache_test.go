/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package user

import (
	"DHBW.Photo-Server"
	"testing"
)

func TestUserManagerCacheNew(t *testing.T) {
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	userManager = nil
	userManagerBefore := userManager
	newUserManager := UserManagerCache()
	if userManagerBefore != nil || newUserManager == nil {
		t.Error("Something went wrong while getting new UserManager")
	}
}

func TestUserManagerCacheExists(t *testing.T) {
	DHBW_Photo_Server.SetUsersFile("somenonexistingcsv.csv")
	userManager = NewUserManager()
	newUserManager := UserManagerCache()
	if newUserManager.UsersFile != DHBW_Photo_Server.UsersFile() {
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
