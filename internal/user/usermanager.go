package user

import (
	"encoding/csv"
	"io"
	"log"
	"os"
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

func (um *UsersManager) LoadUsers() {
	// load Users from Users file into Users array

	csvFile, err := os.Open(um.UsersFile)
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(csvFile)
	um.Users = []*User{}
	for {
		res, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		newUser := User{
			Name:     res[0],
			Password: res[1],
		}
		um.AddUser(&newUser)
	}

	err = csvFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (um *UsersManager) StoreUsers() {
	// save Users array into usersfile

	csvFile, err := os.Create(um.UsersFile)
	if err != nil {
		log.Fatal(err)
	}

	csvWriter := csv.NewWriter(csvFile)
	var users [][]string
	for _, user := range um.Users {
		users = append(users, user.ToCsv())
	}

	err = csvWriter.WriteAll(users)
	if err != nil {
		log.Fatal(err)
	}

	csvWriter.Flush()
	err = csvFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (um *UsersManager) Register(name string, password string) {
	// ???
}
