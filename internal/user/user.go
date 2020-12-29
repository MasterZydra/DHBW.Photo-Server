package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/cryptography"
	"net/http"
)

type User struct {
	Name     string
	password string
	Cookie   http.Cookie
}

func NewUser(name string, password string) User {
	pw, _ := cryptography.CreatePassword(password)
	return User{
		Name:     name,
		password: pw,
	}
}

func FromCsv(csvLine []string) User {
	cookie := http.Cookie{}
	if len(csvLine) > 2 && csvLine[2] != "" {
		cookie = http.Cookie{
			Name:  DHBW_Photo_Server.CookieName,
			Value: csvLine[2],
		}
	}
	return User{
		Name:     csvLine[0],
		password: csvLine[1],
		Cookie:   cookie,
	}
}

func (u *User) ToCsv() []string {
	cookieValue := ""
	if u.Cookie.Value != "" {
		cookieValue = u.Cookie.Value
	}
	return []string{u.Name, u.password, cookieValue}
}

func (u *User) ComparePassword(clearPassword string) (bool, error) {
	return cryptography.ComparePassword(u.password, clearPassword)
}

// TODO: Test schreiben
func (u *User) AddCookie() {
	u.Cookie = cryptography.GenerateCookie(DHBW_Photo_Server.CookieName, u.Name, DHBW_Photo_Server.CookieMaxAge)
}
