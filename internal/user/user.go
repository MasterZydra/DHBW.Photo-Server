package user

import (
	"DHBW.Photo-Server/internal/cryptography"
)

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

func FromCsv(csvLine []string) User {
	return User{
		Name:     csvLine[0],
		password: csvLine[1],
	}
}

func (u *User) ToCsv() []string {
	return []string{u.Name, u.password}
}

func (u *User) ComparePassword(clearPassword string) (bool, error) {
	return cryptography.ComparePassword(u.password, clearPassword)
}
