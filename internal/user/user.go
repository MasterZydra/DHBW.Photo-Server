package user

import (
	"DHBW.Photo-Server/internal/cryptography"
)

// User is used to represent one user entry in the usersFile
// It holds the users name and the users hashed password
type User struct {
	Name     string
	password string
}

// Returns a new User with a new hashed password and the passed name
func NewUser(name string, password string) User {
	pw, _ := cryptography.CreatePassword(password)
	return User{
		Name:     name,
		password: pw,
	}
}

// Returns a new user from a csv line (e.g. usersFile.csv)
func FromCsv(csvLine []string) User {
	return User{
		Name:     csvLine[0],
		password: csvLine[1],
	}
}

// Converts the current User to an array of strings, so it can be written into a csv file with a csvWriter
func (u *User) ToCsv() []string {
	return []string{u.Name, u.password}
}

// Compares the given clear password with the password of the current user and returns a boolean (or error)
func (u *User) ComparePassword(clearPassword string) (bool, error) {
	return cryptography.ComparePassword(u.password, clearPassword)
}
