package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(auth AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		AuthWrapper(auth,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello client")
			}))
}

func TestWithoutPassword(t *testing.T) {
	server := createServer(func(name, pwd string) bool {
		return true
	})
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil || res == nil || res.StatusCode != http.StatusUnauthorized {
		t.Error("wrong status code or no error")
	}
}

func TestWithCorrectPassword(t *testing.T) {
}
