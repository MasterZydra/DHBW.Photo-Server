package user

import (
	"net/http"
	"strings"
)

// TODO: Jones Documentation

type Authenticator interface {
	Authenticate(user, password string, r *http.Request) bool
}

type AuthenticatorFunc func(user, password string, r *http.Request) bool

func (af AuthenticatorFunc) Authenticate(user, password string, r *http.Request) bool {
	return af(user, password, r)
}

func AuthHandlerWrapper(authenticator Authenticator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, pw, ok := r.BasicAuth()
		if ok && authenticator.Authenticate(usr, pw, r) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-authenticate", "Basic realm=\"Please Enter Credentials\"")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func AuthFileServerWrapper(authenticator Authenticator, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usr, pw, ok := r.BasicAuth()
		if ok && authenticator.Authenticate(usr, pw, r) {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
	})
}

func AuthHandler() AuthenticatorFunc {
	return authenticate
}

// führt authenticate aus und überprüft, ob die aktuell angefragte Datei zum Benutzer gehört
func AuthFileServer() AuthenticatorFunc {
	return func(username, password string, r *http.Request) bool {
		urlParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		return authenticate(username, password, r) &&
			len(urlParts) > 1 && strings.ToLower(urlParts[1]) == strings.ToLower(username)
	}
}

func authenticate(username, password string, r *http.Request) bool {
	um := GetImageManager()
	ok, _ := um.Authenticate(username, password)
	return ok
}
