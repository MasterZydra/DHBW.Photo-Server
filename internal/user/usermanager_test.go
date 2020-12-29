package user

import (
	"DHBW.Photo-Server"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	csvFile, err := os.Create(DHBW_Photo_Server.TestUserFile)
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(csvFile)
	var data = [][]string{
		{DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Hash, DHBW_Photo_Server.CookieValue1},
		{DHBW_Photo_Server.User2Name, DHBW_Photo_Server.Pw2Hash, DHBW_Photo_Server.CookieValue2},
	}
	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewUsersManager(t *testing.T) {
	um1 := NewUsersManager()
	um2 := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	um3 := NewUsersManager("")
	if um1.UsersFile != DHBW_Photo_Server.ProdUserFile || um2.UsersFile != DHBW_Photo_Server.TestUserFile || um3.UsersFile != DHBW_Photo_Server.ProdUserFile {
		t.Error("At least one users file is not correct in the usermanager")
	}
}

func TestAddUserCount(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	usersCountBefore := len(um.Users)
	newUser := User{
		Name:     "testuser",
		password: DHBW_Photo_Server.Pw1Hash,
	}
	um.AddUser(&newUser)
	if usersCountBefore == len(um.Users) {
		t.Errorf("Users should have more than %v users after adding one to it", usersCountBefore)
	}
}
func TestAddUserContent(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	newUser := NewUser("manuela", "testPW")
	um.AddUser(&newUser)
	lastUser := um.Users[len(um.Users)-1]
	if lastUser.Name != newUser.Name || lastUser.password != newUser.password {
		t.Error("Last Username isn't the one added before")
	}
}

func TestLoadUsersCount(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	usersCountBefore := len(um.Users)
	_ = um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Errorf("Users should be more than %v after loading it from usersFile.csv", usersCountBefore)
	}
}

func TestLoadUsersContent(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	max := um.Users[0]
	ana := um.Users[1]
	if max.Name != DHBW_Photo_Server.User1Name || max.password != DHBW_Photo_Server.Pw1Hash || ana.Name != DHBW_Photo_Server.User2Name || ana.password != DHBW_Photo_Server.Pw2Hash {
		t.Error("At least one Username wasn't loaded correctly from the usersfile")
	}
}

func TestLoadUsersWrongFile(t *testing.T) {
	um := NewUsersManager("wrongfile.csv")
	err := um.LoadUsers()
	if err == nil {
		t.Error("There should be an error thrown because wrongfile.csv doesn't exist")
	}
}

func TestLoadUsersMultiple(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	usersCountBetween := len(um.Users)
	_ = um.LoadUsers()
	if usersCountBetween != len(um.Users) {
		t.Error("If loadUsers is executed twice it should still have the same amount of users")
	}
}

func TestUsersManager_GetUserSuccess(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	max := um.GetUser(DHBW_Photo_Server.User1Name)
	if max.Name != DHBW_Photo_Server.User1Name || max.password != DHBW_Photo_Server.Pw1Hash {
		t.Errorf("Something went wrong while getting user %v", DHBW_Photo_Server.User1Name)
	}
}

func TestUsersManager_GetUserFail(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	username := "unknownUser"
	unknown := um.GetUser(username)
	if unknown != nil {
		t.Errorf("Shouldn't get user %v", username)
	}
}

func TestUsersManager_GetUserByCookieSuccess(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	cookie := http.Cookie{
		Value: DHBW_Photo_Server.User1Name + DHBW_Photo_Server.CookieValueSeparator + "TEa8eQ_-ZVoeZaz6z5XaUM1NOQI=",
	}
	userObj := um.GetUserByCookie(&cookie)
	if userObj == nil || userObj.Name != DHBW_Photo_Server.User1Name {
		t.Error("User object shouldn't be nil. It should get the user by the provided cookie")
	}
}

func TestUsersManager_GetUserByCookieUnknownUser(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	cookie := http.Cookie{
		Value: "unknownUser" + DHBW_Photo_Server.CookieValueSeparator + "TEa8eQ_-ZVoeZaz6z5XaUM1NOQI=",
	}
	userObj := um.GetUserByCookie(&cookie)
	if userObj != nil {
		t.Error("User object shouldn't be nil. It should get the user by the provided cookie")
	}
}

func TestUsersManager_GetUserByCookieNoUser(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	cookie := http.Cookie{
		Value: "TEa8eQ_-ZVoeZaz6z5XaUM1NOQI=",
	}
	userObj := um.GetUserByCookie(&cookie)
	if userObj != nil {
		t.Error("User object shouldn't be nil. It should get the user by the provided cookie")
	}
}

func TestStoreUsers(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	_ = um.LoadUsers()
	usersCountBefore := len(um.Users)
	newUser := NewUser("manuela", "1234")
	um.AddUser(&newUser)
	_ = um.StoreUsers()
	_ = um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Error("Storing the new Username has not worked")
	}
}

func TestUsersManager_UserExistsTrue(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	exists, _ := um.UserExists(DHBW_Photo_Server.User1Name)
	if !exists {
		t.Errorf("user '%v' should exist", DHBW_Photo_Server.User1Name)
	}
}

func TestUsersManager_UserExistsFalse(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	username := "userThatDoesntExist"
	exists, _ := um.UserExists(username)
	if exists {
		t.Errorf("user '%v' should exist", username)
	}
}

func TestRegister(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	err := um.Register("robert", "1234")
	if err != nil {
		t.Errorf("There was an error during registration: %v", err)
	}
}

func TestUsersManager_RegisterWrongUsername(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	err := um.Register("robert*", "1234")
	if err == nil {
		t.Errorf("There should be an error stating that the username is invalid")
	}
}

func TestUsersManager_RegisterExistingUser(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	err := um.Register(DHBW_Photo_Server.User1Name, "0987")
	if err == nil {
		t.Errorf("You shouldn't be able to add the Username %v, since it already exists", DHBW_Photo_Server.User1Name)
	}
}

func TestUsersManager_AuthenticateCorrect(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	ok, _ := um.Authenticate(DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)
	if !ok {
		t.Errorf("Authentication should be valid, but it isn't")
	}
}

func TestUsersManager_AuthenticateWrongUser(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	ok, _ := um.Authenticate("wrongUserName", DHBW_Photo_Server.Pw1Clear)
	if ok {
		t.Errorf("username should be wrong, but it seems to be correct")
	}
}

func TestUsersManager_AuthenticateWrongPW(t *testing.T) {
	um := NewUsersManager(DHBW_Photo_Server.TestUserFile)
	ok, _ := um.Authenticate(DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw2Clear)
	if ok {
		t.Errorf("password should be wrong, but it seems to be correct")
	}
}
