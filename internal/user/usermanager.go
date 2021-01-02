package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

// TODO: Jones Documentation

type UserManager struct {
	Users     []*User
	UsersFile string
}

func NewUserManager(args ...string) UserManager {
	// TODO: LoadUsers direkt beim NewUserManager ausführen -> refactoring
	usersFile := "usersFile.csv"
	if args != nil && args[0] != "" {
		usersFile = args[0]
	}
	return UserManager{
		UsersFile: usersFile,
	}
}

func (um *UserManager) AddUser(user *User) {
	um.Users = append(um.Users, user)
}

func (um *UserManager) LoadUsers() error {
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

		newUser := FromCsv(res)
		um.AddUser(&newUser)
	}

	err = csvFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (um *UserManager) GetUser(username string) *User {
	for _, userObj := range um.Users {
		if strings.ToLower(userObj.Name) == strings.ToLower(username) {
			return userObj
		}
	}
	return nil
}

func (um *UserManager) StoreUsers() error {
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

func (um *UserManager) Register(name string, password string) error {
	// check if user has not allowed characters in it (allowed are: a-z,A-Z,0-9,-,_ and .
	matched, err := regexp.MatchString(DHBW_Photo_Server.UsernameRegexBlacklist, name)
	if matched {
		return errors.New("Username can only contain a-z,A-Z,0-9,-,_ and .")
	}

	// check if Username already exists, and yes: error; if not add it to usersfile
	exists, err := um.UserExists(name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Username '" + name + "' already exists")
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

func (um *UserManager) UserExists(name string) (bool, error) {
	err := um.LoadUsers()
	if err != nil {
		return false, err
	}
	user := um.GetUser(name)
	if user != nil {
		return true, nil
	}
	return false, nil
}

func (um *UserManager) Authenticate(user string, pw string) (bool, error) {
	err := um.LoadUsers()
	if err != nil {
		return false, err
	}
	userObj := um.GetUser(user)
	if userObj != nil {
		ok, _ := userObj.ComparePassword(pw)
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// TODO: test
func (um *UserManager) AuthenticateHashedPassword(username string, hashedPw string) bool {
	userObj := um.GetUser(username)
	if userObj != nil {
		return userObj.password == hashedPw
	}
	return false
}
