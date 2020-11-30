package user

import "DHBW.Photo-Server/internal/cryptography"

// TODO: besprechen: Password private machen? Damit mans nicht editieren kann
type User struct {
	Name     string
	Password string
}

func NewUser(name string, password string) User {
	pw, _ := cryptography.CreatePassword(password)
	return User{
		Name:     name,
		Password: pw,
	}
}

func (u *User) ToCsv() []string {
	return []string{u.Name, u.Password}
}
