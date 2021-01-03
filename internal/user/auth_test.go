package user

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(auth AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		AuthHandlerWrapper(auth,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "string from server")
			}))
}

func createFileServer(auth AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		AuthFileServerWrapper(auth,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "string from server")
			})))
}

func doRequest(url string, username string, password string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := getRequest(url, username, password)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func getRequest(url string, username string, password string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if req == nil || err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	return req, nil
}

func TestHandlerWrapperWithoutPassword(t *testing.T) {
	server := createServer(func(name, pwd string, r *http.Request) bool {
		return true
	})
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil || res == nil || res.StatusCode != http.StatusUnauthorized {
		t.Error("wrong status code or no error")
	}
}

func TestHandlerWrapperWithCorrectPassword(t *testing.T) {
	var receivedName, receivedPwd string
	server := createServer(func(name, pwd string, r *http.Request) bool {
		receivedName = name
		receivedPwd = pwd
		return true
	})
	defer server.Close()
	res, _ := doRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)

	if res.StatusCode != http.StatusOK || receivedName != DHBW_Photo_Server.User1Name || receivedPwd != DHBW_Photo_Server.Pw1Clear {
		t.Error("Authentication with password didn't work, but it should")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != "string from server\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestHandlerWrapperWithWrongPassword(t *testing.T) {
	var receivedName, receivedPwd string
	server := createServer(func(name, pwd string, r *http.Request) bool {
		receivedName = name
		receivedPwd = pwd
		return false
	})
	defer server.Close()
	res, _ := doRequest(server.URL, "wrong", "credentials")

	if res.StatusCode != http.StatusUnauthorized || receivedName != "wrong" || receivedPwd != "credentials" {
		t.Error("Authentication with wrong password did not fail or status code is wrong")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != http.StatusText(http.StatusUnauthorized)+"\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestFileServerWrapperWithoutPassword(t *testing.T) {
	server := createFileServer(func(name, pwd string, r *http.Request) bool {
		return true
	})
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil || res == nil || res.StatusCode != http.StatusForbidden {
		t.Error("wrong status code or no error")
	}
}

func TestFileServerWrapperWithCorrectPassword(t *testing.T) {
	var receivedName, receivedPwd string
	server := createFileServer(func(name, pwd string, r *http.Request) bool {
		receivedName = name
		receivedPwd = pwd
		return true
	})
	defer server.Close()
	res, _ := doRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)

	if res.StatusCode != http.StatusOK || receivedName != DHBW_Photo_Server.User1Name || receivedPwd != DHBW_Photo_Server.Pw1Clear {
		t.Error("Authentication with password didn't work, but it should")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != "string from server\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestFileServerWrapperWithWrongPassword(t *testing.T) {
	var receivedName, receivedPwd string
	server := createFileServer(func(name, pwd string, r *http.Request) bool {
		receivedName = name
		receivedPwd = pwd
		return false
	})
	defer server.Close()
	res, _ := doRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)

	if res.StatusCode != http.StatusForbidden || receivedName != DHBW_Photo_Server.User1Name || receivedPwd != DHBW_Photo_Server.Pw1Clear {
		t.Error("Authentication with wrong password did not fail or status code is wrong")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != http.StatusText(http.StatusForbidden)+"\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestAuthenticateHandlerWithoutPassword(t *testing.T) {
	usersFile = DHBW_Photo_Server.TestUserFile
	server := createFileServer(AuthHandler())
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil || res == nil || res.StatusCode != http.StatusForbidden {
		t.Error("wrong status code or no error")
	}
}

func TestAuthenticateHandlerWithCorrectPassword(t *testing.T) {
	server := createFileServer(AuthHandler())
	defer server.Close()
	res, _ := doRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)

	if res.StatusCode != http.StatusOK {
		t.Error("Authentication with password didn't work, but it should")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != "string from server\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestAuthenticateHandlerWithWrongPassword(t *testing.T) {
	server := createFileServer(AuthHandler())
	defer server.Close()
	res, _ := doRequest(server.URL, "wrong", "credentials")

	if res.StatusCode != http.StatusForbidden {
		t.Error("Authentication with wrong password did not fail or status code is wrong")
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != http.StatusText(http.StatusForbidden)+"\n" {
			t.Error("wrong message in body from server")
		}
	}
}

func TestAuthenticateFileServerWrongPath(t *testing.T) {
	server := createFileServer(AuthFileServer())
	defer server.Close()

	client := &http.Client{}
	req, _ := getRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)
	req.URL.Path = "/images/wronguser/otherimage.jpg"
	res, _ := client.Do(req)

	if res.StatusCode != http.StatusForbidden {
		t.Error("file server authentication with wrong path has wrong status code")
	}
}

func TestAuthenticateFileServerCorrectPath(t *testing.T) {
	server := createFileServer(AuthFileServer())
	defer server.Close()

	client := &http.Client{}
	req, _ := getRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)
	req.URL.Path = "/images/" + DHBW_Photo_Server.User1Name + "/someimage.jpg"
	res, _ := client.Do(req)

	if res.StatusCode != http.StatusOK {
		t.Error("file server authentication with wrong path has wrong status code")
	}
}

func TestAuthenticateFileServerPathTooShort(t *testing.T) {
	server := createFileServer(AuthFileServer())
	defer server.Close()

	client := &http.Client{}
	req, _ := getRequest(server.URL, DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Clear)
	req.URL.Path = "/someimage.jpg"
	res, _ := client.Do(req)

	if res.StatusCode != http.StatusForbidden {
		t.Error("file server authentication with wrong path has wrong status code")
	}
}
