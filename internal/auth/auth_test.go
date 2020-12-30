package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(auth AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		Wrapper(auth,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "string from server")
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

func doRequestWithPassword(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if req == nil || err != nil {
		return nil, err
	}
	req.SetBasicAuth("<username>", "<password>")
	res, err := client.Do(req)
	return res, nil
}

func TestWithCorrectPassword(t *testing.T) {
	var receivedName, receivedPwd string
	server := createServer(func(name, pwd string) bool {
		receivedName = name
		receivedPwd = pwd
		return true
	})
	defer server.Close()
	res, _ := doRequestWithPassword(server.URL)

	if res.StatusCode != http.StatusOK || receivedName != "<username>" || receivedPwd != "<password>" {
		t.Error("Authentication with password didn't work, but it should")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != "string from server\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestWithWrongPassword(t *testing.T) {
	var receivedName, receivedPwd string
	server := createServer(func(name, pwd string) bool {
		receivedName = name
		receivedPwd = pwd
		return false
	})
	defer server.Close()
	res, _ := doRequestWithPassword(server.URL)

	if res.StatusCode != http.StatusUnauthorized || receivedName != "<username>" || receivedPwd != "<password>" {
		t.Error("Authentication with wrong password did not fail")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != http.StatusText(http.StatusUnauthorized)+"\n" {
			t.Error("wrong message in body from server")
		}
	}
}
