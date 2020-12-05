package main

import (
	user2 "DHBW.Photo-Server/internal/user"
	"net/http"
)

func AuthWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, pw, ok := r.BasicAuth()
		if ok && authenticate(usr, pw) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate",
				"Basic realm=\"Best Go-Server evaaa!\"")
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		}
	}
}

func authenticate(user string, pw string) bool {
	um := user2.NewUsersManager()
	_ = um.LoadUsers()
	for _, userObj := range um.Users {
		if userObj.Name == user {
			ok, err := userObj.ComparePassword(pw)
			if ok && err != nil {
				return true
			}
		}
	}
	return false
}
