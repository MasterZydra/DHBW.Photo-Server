/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package user

import (
	"encoding/hex"
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "Helmut"
	user := NewUser(name, "12345678")
	if user.Name != name || user.password == "" || len(user.OrderList.Entries) > 0 {
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
