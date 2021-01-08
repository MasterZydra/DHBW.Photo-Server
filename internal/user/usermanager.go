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

// UserManager is used to manage a list of Users stored in the UsersFile
type UserManager struct {
	Users     []*User
	UsersFile string
}

// create a new UserManager with the default usersFile and return a pointer to it
func NewUserManager() *UserManager {
	return &UserManager{
		UsersFile: DHBW_Photo_Server.UsersFile(),
	}
}

// adds a User pointer to the Users slice
func (um *UserManager) AddUser(user *User) {
	um.Users = append(um.Users, user)
}

// Load all users from the configured usersFile into the Users slice
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

// Get a specific User pointer via the name that is passed as an argument
// username is not case sensitive.
func (um *UserManager) GetUser(username string) *User {
	for _, userObj := range um.Users {
		if strings.ToLower(userObj.Name) == strings.ToLower(username) {
			return userObj
		}
	}
	return nil
}

// Convert each user in the Users slice to csv string and write it to the configured usersFile
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

// Register first checks if the passed username contains a not allowed character.
// After that it checks if there is already a user with the passed username.
// If it passes those two checks a new user is created, added, stored (into the usersfile).
// Then the Users are loaded again to be up to date.
func (um *UserManager) Register(name string, password string) error {
	// check if user has not allowed characters in it (allowed are: a-z,A-Z,0-9,-,_ and .
	matched, err := regexp.MatchString(DHBW_Photo_Server.UsernameRegexBlacklist, name)
	if matched {
		return errors.New("Username can only contain a-z,A-Z,0-9,-,_ and .")
	}

	// check if Username already exists, and yes: error; if not add it to usersfile
	exists := um.UserExists(name)
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

// Gets the user via the passed username and returns if this user exists or not
func (um *UserManager) UserExists(name string) bool {
	user := um.GetUser(name)
	return user != nil
}

// Gets the user via the passed username and executes the ComparePassword function
// to check if the provided password is correct.
// returns boolean value if username and password is correct or not
func (um *UserManager) Authenticate(user string, pw string) bool {
	userObj := um.GetUser(user)
	if userObj != nil {
		ok, _ := userObj.ComparePassword(pw)
		return ok
	}
	return false
}
