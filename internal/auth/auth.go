package auth

import (
	"DHBW.Photo-Server/internal/user"
	"DHBW.Photo-Server/internal/util"
	"net/http"
)

type Authenticator interface {
	Authenticate(user, password string) bool
}

type AuthenticatorFunc func(user, password string) bool

func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

func Wrapper(authenticator Authenticator, handler http.HandlerFunc, whitelist []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, pw, ok := r.BasicAuth()
		if util.ContainsString(whitelist, r.URL.Path) || ok && authenticator.Authenticate(usr, pw) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate",
				"Basic realm=\"Please Enter Credentials\"")
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		}
	}
}

// TODO: tests schreiben?
func Authenticate() AuthenticatorFunc {
	return func(username, password string) bool {
		um := user.NewUsersManager()
		ok, _ := um.Authenticate(username, password)
		return ok
	}
}
