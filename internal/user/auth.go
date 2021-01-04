package user

import (
	"net/http"
	"strings"
)

// interface and type declarations for the abstraction used in the authentication wrapper functions
type Authenticator interface {
	Authenticate(user, password string, r *http.Request) bool
}

type AuthenticatorFunc func(user, password string, r *http.Request) bool

func (af AuthenticatorFunc) Authenticate(user, password string, r *http.Request) bool {
	return af(user, password, r)
}

// The authentication wrapper AuthHandlerWrapper can be used as second argument for http.HandleFunc
// It returns a function from type http.HandlerFunc.
// Within this function the username and password is retrieved from the reqest object via BasicAuth().
// These are passed to an Authenticate function from the interface of the first argument of AuthHandlerWrapper.
// If this Authenticate function returns true, the handler (second argument) is executed.
// If not it adds the Header "www-authenticate" to the response and writes the Status 401 Unauthorized.
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

// The AuthFileServerWrapper is basically the same as the AuthHandlerWrapper but returns a http.Handler
// and executes the handler via ServeHTTP().
// If the Autenticate function returned false it writes the status 403 Forbidden.
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

// The AuthHandler just returns the local AuthenticatorFunc authenticate.
func AuthHandler() AuthenticatorFunc {
	return authenticate
}

// Executes authenticate and checks if the currently requested file (image) belongs to the current authenticated user.
func AuthFileServer() AuthenticatorFunc {
	return func(username, password string, r *http.Request) bool {
		urlParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		return authenticate(username, password, r) &&
			len(urlParts) > 1 && strings.ToLower(urlParts[1]) == strings.ToLower(username)
	}
}

// Gets the imageManager from cache and executes Authenticate with the given username and password.
func authenticate(username, password string, r *http.Request) bool {
	um := GetImageManager()
	ok := um.Authenticate(username, password)
	return ok
}
