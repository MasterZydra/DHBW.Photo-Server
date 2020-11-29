package user

import "testing"

func TestNewUsersManager(t *testing.T) {
	um1 := NewUsersManager()
	um2 := NewUsersManager("usersFile_test.csv")
	um3 := NewUsersManager("")
	if um1.UsersFile != "usersFile.csv" || um2.UsersFile != "usersFile_test.csv" || um3.UsersFile != "usersFile.csv" {
		t.Error("At least one users file is not correct in the usermanager")
	}
}

func TestAddUserCount(t *testing.T) {
	um := NewUsersManager("usersFile_test.csv")
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
	um := NewUsersManager("usersFile_test.csv")
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
	um := NewUsersManager("usersFile_test.csv")
	usersCountBefore := len(um.Users)
	um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Errorf("Users should be more than %v after loading it from usersFile.csv", usersCountBefore)
	}
}

func TestLoadUsersContent(t *testing.T) {
	um := NewUsersManager("usersFile_test.csv")
	um.LoadUsers()
	max := um.Users[0]
	ana := um.Users[1]
	if max.Name != "Max" || max.Password != "1234" || ana.Name != "Ana" || ana.Password != "5678" {
		t.Error("At least one user wasn't loaded correctly from the usersfile")
	}
}

func TestLoadUsersMultiple(t *testing.T) {
	um := NewUsersManager("usersFile_test.csv")
	um.LoadUsers()
	usersCountBetween := len(um.Users)
	um.LoadUsers()
	if usersCountBetween != len(um.Users) {
		t.Errorf("If loadUsers is executed twice it should still have the same amount of users")
	}
}

func TestStoreUsers(t *testing.T) {
	um := NewUsersManager("usersFile_test.csv")
	um.LoadUsers()
	usersCountBefore := len(um.Users)
	newUser := User{
		Name:     "testuser",
		Password: "1234",
	}
	um.AddUser(&newUser)
	um.StoreUsers()
	um.LoadUsers()
	if usersCountBefore == len(um.Users) {
		t.Error("Storing the new user has not worked")
	}
}
