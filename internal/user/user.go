/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package user

import (
	"DHBW.Photo-Server/internal/cryptography"
)

// User is used to represent one user entry in the usersFile
// It holds the users name, the users hashed password and the orderList from the user
type User struct {
	Name      string
	password  string
	OrderList *OrderList
}

// Returns a new User with a new hashed password and the passed name
func NewUser(name string, password string) User {
	pw, _ := cryptography.CreatePassword(password)
	return User{
		Name:      name,
		password:  pw,
		OrderList: &OrderList{Entries: nil},
	}
}

// Returns a new user from a csv line (e.g. usersFile.csv)
func FromCsv(csvLine []string) User {
	return User{
		Name:      csvLine[0],
		password:  csvLine[1],
		OrderList: &OrderList{Entries: nil},
	}
}

// Converts the current User to an array of strings, so it can be written into a csv file with a csvWriter
// Note that the OrderList will not be stored in the csv
func (u *User) ToCsv() []string {
	return []string{u.Name, u.password}
}

// Compares the given clear password with the password of the current user and returns a boolean (or error)
func (u *User) ComparePassword(clearPassword string) (bool, error) {
	return cryptography.ComparePassword(u.password, clearPassword)
}
