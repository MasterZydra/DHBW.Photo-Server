package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/cryptography"
	"DHBW.Photo-Server/internal/image"
	"DHBW.Photo-Server/internal/order"
	"errors"
)

// User is used to represent one user entry in the usersFile
// It holds the users name, the users hashed password and the orderList from the user
type User struct {
	Name      string
	password  string
	OrderList []*order.ListEntry
}

// Returns a new User with a new hashed password and the passed name
func NewUser(name string, password string) User {
	pw, _ := cryptography.CreatePassword(password)
	return User{
		Name:      name,
		password:  pw,
		OrderList: []*order.ListEntry{},
	}
}

// Returns a new user from a csv line (e.g. usersFile.csv)
func FromCsv(csvLine []string) User {
	return User{
		Name:      csvLine[0],
		password:  csvLine[1],
		OrderList: []*order.ListEntry{},
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

// Checks first if this image is already in the orderList
// After that it adds a new order entry with the passed image object pointer
func (u *User) AddOrderEntry(image *image.Image) error {
	_, entry := u.GetOrderEntry(image.Name)
	if entry != nil {
		return errors.New("Image '" + image.Name + "' already in order list")
	}
	newEntry := order.ListEntry{
		Image:          image,
		Format:         DHBW_Photo_Server.OrderListFormats[1], // Letter (8.5 x 11)
		NumberOfPrints: 1,
	}
	u.OrderList = append(u.OrderList, &newEntry)
	return nil
}

// Gets the order entry with the passed imageName string and removes this entry from the orderList preserving order
func (u *User) RemoveOrderEntry(imageName string) bool {
	i, entry := u.GetOrderEntry(imageName)
	if entry != nil {
		u.OrderList = append(u.OrderList[:i], u.OrderList[i+1:]...)
		return true
	}
	return false
}

// Gets the order entry with the image that has the name of the passed imageName
// Returns it's index in the orderList and the a pointer to the order.ListEntry
func (u *User) GetOrderEntry(imageName string) (index int, entry *order.ListEntry) {
	for i, entry := range u.OrderList {
		if entry.Image.Name == imageName {
			return i, entry
		}
	}
	return -1, nil
}
