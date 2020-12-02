package user

import (
	"encoding/csv"
	"log"
	"os"
	"testing"
)

const testUserFile = "usersFile_test.csv"
const prodUserFile = "usersFile.csv"

//pw: test123
const pw1Hash = "6dfbf8730f569dba965ead781f536f7b5ccc2f6b9824f0e49e6878b349a94bc9186c7d7145df80e841def14f3dd70791"

//pw: 123test
const pw2Hash = "e9fa8567977ba0db64bc5d5f18118d377032a4820c38ed404400b52bdb6751b9e27c0beb37e35f2bf75608c634a28990"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	csvFile, err := os.Create(testUserFile)
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(csvFile)
	var data = [][]string{
		{"Max", pw1Hash},
		{"Ana", pw2Hash},
	}
	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewUsersManager(t *testing.T) {
	um1 := NewUsersManager()
	um2 := NewUsersManager(testUserFile)
	um3 := NewUsersManager("")
	if um1.UsersFile != prodUserFile || um2.UsersFile != testUserFile || um3.UsersFile != prodUserFile {
		t.Error("At least one users file is not correct in the usermanager")
	}
}

func TestAddUserCount(t *testing.T) {
	um := NewUsersManager(testUserFile)
	usersCountBefore := len(um.Users)
	newUser := User{
		Name:     "testuser",
		password: pw1Hash,
	}
	um.AddUser(&newUser)
	if usersCountBefore == len(um.Users) {
		t.Errorf("Users should have more than %v users after adding one to it", usersCountBefore)
	}
}
func TestAddUserContent(t *testing.T) {
	um := NewUsersManager(testUserFile)
	newUser := NewUser("manuela", "testPW")
	um.AddUser(&newUser)
	lastUser := um.Users[len(um.Users)-1]
	if lastUser.Name != newUser.Name || lastUser.password != newUser.password {
		t.Error("Last User isn't the one added before")
	}
}

func TestLoadUsersCount(t *testing.T) {
	um := NewUsersManager(testUserFile)
	usersCountBefore := len(um.Users)
	_ = um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Errorf("Users should be more than %v after loading it from usersFile.csv", usersCountBefore)
	}
}

func TestLoadUsersContent(t *testing.T) {
	um := NewUsersManager(testUserFile)
	_ = um.LoadUsers()
	max := um.Users[0]
	ana := um.Users[1]
	if max.Name != "Max" || max.password != pw1Hash || ana.Name != "Ana" || ana.password != pw2Hash {
		t.Error("At least one User wasn't loaded correctly from the usersfile")
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
	um := NewUsersManager(testUserFile)
	_ = um.LoadUsers()
	usersCountBetween := len(um.Users)
	_ = um.LoadUsers()
	if usersCountBetween != len(um.Users) {
		t.Error("If loadUsers is executed twice it should still have the same amount of users")
	}
}

func TestStoreUsers(t *testing.T) {
	um := NewUsersManager(testUserFile)
	_ = um.LoadUsers()
	usersCountBefore := len(um.Users)
	newUser := NewUser("manuela", "1234")
	um.AddUser(&newUser)
	_ = um.StoreUsers()
	_ = um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Error("Storing the new User has not worked")
	}
}

func TestRegister(t *testing.T) {
	um := NewUsersManager(testUserFile)
	err := um.Register("robert", "1234")
	if err != nil {
		t.Errorf("There was an error during registration: %v", err)
	}
}

func TestRegisterExistingUser(t *testing.T) {
	um := NewUsersManager(testUserFile)
	err := um.Register("max", "0987")
	if err == nil {
		t.Errorf("You shouldn't be able to add the User max, since it already exists")
	}
}
