package user

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type UsersManager struct {
	Users     []*User
	UsersFile string
}

func NewUsersManager(args ...string) UsersManager {
	usersFile := "usersFile.csv"
	if args != nil && args[0] != "" {
		usersFile = args[0]
	}
	return UsersManager{
		UsersFile: usersFile,
	}
}

func (um *UsersManager) AddUser(user *User) {
	um.Users = append(um.Users, user)
}

func (um *UsersManager) LoadUsers() error {
	// load Users from Users file into Users array

	csvFile, err := os.Open(um.UsersFile)
	if err != nil {
		return err
	}

	csvReader := csv.NewReader(csvFile)
	um.Users = []*User{}
	for {
		res, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		newUser := User{
			Name:     res[0],
			password: res[1],
		}
		um.AddUser(&newUser)
	}

	err = csvFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (um *UsersManager) StoreUsers() error {
	// save Users array into usersfile

	csvFile, err := os.Create(um.UsersFile)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(csvFile)
	var users [][]string
	for _, user := range um.Users {
		users = append(users, user.ToCsv())
	}

	err = csvWriter.WriteAll(users)
	if err != nil {
		return err
	}

	csvWriter.Flush()
	err = csvFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (um *UsersManager) Register(name string, password string) error {
	// check if Username already exists, and yes: error; if not add it to usersfile
	// TODO: create Username folder?
	err := um.LoadUsers()
	if err != nil {
		return err
	}

	for _, user := range um.Users {
		if strings.ToLower(name) == strings.ToLower(user.Name) {
			return errors.New("Username '" + name + "' already exists")
		}
	}

	newUser := NewUser(name, password)
	um.AddUser(&newUser)

	err = um.StoreUsers()
	if err != nil {
		return err
	}
	err = um.LoadUsers()
	if err != nil {
		return err
	}

	return nil
}

func (um *UsersManager) Authenticate(user string, pw string) bool {
	_ = um.LoadUsers()
	for _, userObj := range um.Users {
		if userObj.Name == user {
			ok, _ := userObj.ComparePassword(pw)
			if ok {
				return true
			}
		}
	}
	return false
}
