package user

import (
	"encoding/hex"
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "Helmut"
	user := NewUser(name, "12345678")
	if user.Name != name || user.password == "" {
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

func TestUser_ToCsv(t *testing.T) {
	user := User{
		Name:     "test",
		password: "mostcomplexpasswordever",
	}
	userCsvRow := user.ToCsv()
	if len(userCsvRow) != 2 || userCsvRow[0] != user.Name || userCsvRow[1] != user.password {
		t.Error("Error while converting user object to csv")
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
