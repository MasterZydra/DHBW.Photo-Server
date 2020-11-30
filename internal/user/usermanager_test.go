package user

import (
	"encoding/csv"
	"log"
	"os"
	"testing"
)

const testUserFile = "usersFile_test.csv"
const prodUserFile = "usersFile.csv"

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
		{"Max", "1234"},
		{"Ana", "5678"},
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
		Password: "1234",
	}
	um.AddUser(&newUser)
	if usersCountBefore == len(um.Users) {
		t.Errorf("Users should have more than %v users after adding one to it", usersCountBefore)
	}
}
func TestAddUserContent(t *testing.T) {
	um := NewUsersManager(testUserFile)
	newUser := User{
		Name:     "testuser",
		Password: "1234",
	}
	um.AddUser(&newUser)
	lastUser := um.Users[len(um.Users)-1]
	if lastUser.Name != newUser.Name || lastUser.Password != newUser.Password {
		t.Error("Last user isn't the one added before")
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
	if max.Name != "Max" || max.Password != "1234" || ana.Name != "Ana" || ana.Password != "5678" {
		t.Error("At least one user wasn't loaded correctly from the usersfile")
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
	newUser := User{
		Name:     "manuela",
		Password: "1234",
	}
	um.AddUser(&newUser)
	_ = um.StoreUsers()
	_ = um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Error("Storing the new user has not worked")
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
		t.Errorf("You shouldn't be able to add the user max, since it already exists")
	}
}
