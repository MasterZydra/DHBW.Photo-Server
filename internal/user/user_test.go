package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"encoding/hex"
	"net/http"
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

func TestFromCsvWithCookie(t *testing.T) {
	var csvLine = []string{
		DHBW_Photo_Server.User1Name,
		DHBW_Photo_Server.Pw1Hash,
		DHBW_Photo_Server.CookieValue1,
	}
	newUser := FromCsv(csvLine)
	if newUser.Name != DHBW_Photo_Server.User1Name || newUser.Cookie.Value != DHBW_Photo_Server.CookieValue1 {
		t.Error("Some information is wrong from the csv file")
	}
}

func TestFromCsvNoCookie(t *testing.T) {
	var csvLine = []string{
		DHBW_Photo_Server.User1Name,
		DHBW_Photo_Server.Pw1Hash,
	}
	newUser := FromCsv(csvLine)
	if newUser.Name != DHBW_Photo_Server.User1Name || newUser.Cookie.Value != "" {
		t.Error("Some information is wrong from the csv file")
	}
}
func TestFromCsvNoCookieBlank(t *testing.T) {
	var csvLine = []string{
		DHBW_Photo_Server.User1Name,
		DHBW_Photo_Server.Pw1Hash,
		"",
	}
	newUser := FromCsv(csvLine)
	if newUser.Name != DHBW_Photo_Server.User1Name || newUser.Cookie.Value != "" {
		t.Error("Some information is wrong from the csv file")
	}
}

func TestUser_ToCsvNoCookie(t *testing.T) {
	user := User{
		Name:     "test",
		password: "mostcomplexpasswordever",
	}
	userCsvRow := user.ToCsv()
	if len(userCsvRow) != 3 || userCsvRow[0] != user.Name || userCsvRow[1] != user.password || userCsvRow[2] != "" {
		t.Error("Error while converting user object to csv")
	}
}

func TestUser_ToCsvCookie(t *testing.T) {
	cookie := http.Cookie{
		Name:  DHBW_Photo_Server.CookieName,
		Value: "programmerOnlyWriteDumbTestValues",
	}
	user := User{
		Name:     "test",
		password: "mostcomplexpasswordever",
		Cookie:   cookie,
	}
	userCsvRow := user.ToCsv()
	if len(userCsvRow) != 3 || userCsvRow[0] != user.Name || userCsvRow[1] != user.password || userCsvRow[2] != cookie.Value {
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
