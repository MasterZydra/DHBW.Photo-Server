package user

import "DHBW.Photo-Server/internal/cryptography"

type User struct {
	Name     string
	password string
}

func NewUser(name string, password string) User {
	pw, _ := cryptography.CreatePassword(password)
	return User{
		Name:     name,
		password: pw,
	}
}

func (u *User) ToCsv() []string {
	return []string{u.Name, u.password}
}

//func (u *User) ComparePassword(hexHash string) {
//	cryptography.ComparePassword(hexHash, u.password)
//	// TODO
//
//}
