package main

import (
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
}

func TestMustParamWrongMethod(t *testing.T) {
	server := createServer(MustParam(func(w http.ResponseWriter, r *http.Request) {
		return
	}, http.MethodPost))

	response, _ := http.Get(server.URL)
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Wrong status code")
	}
}

func TestMustParamMissingGetParams(t *testing.T) {
	server := createServer(MustParam(func(w http.ResponseWriter, r *http.Request) {
		return
	}, http.MethodGet, "param1", "param2"))

	response, _ := http.Get(server.URL)
	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestMustParamWithGetParams(t *testing.T) {
	server := createServer(MustParam(func(w http.ResponseWriter, r *http.Request) {
		return
	}, http.MethodGet, "param1", "param2"))

	response, _ := http.Get(server.URL + "?param1=some&param2=thing")
	if response.StatusCode != http.StatusOK {
		t.Error("Wrong status code")
	}
}

func newPostRequest(url string, data interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
}

func executeRequest(r *http.Request) (*http.Response, error) {
	c := http.Client{}
	return c.Do(r)
}

func TestRegisterHandler(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	data := api.RegisterReqData{
		Username:             "test",
		Password:             "sec123",
		PasswordConfirmation: "sec123",
	}
	req, _ := newPostRequest(server.URL, data)
	response, err := executeRequest(req)

	if err != nil || response.StatusCode != http.StatusOK {
		t.Error("Status code wrong")
	}
}

func TestRegisterHandlerInvalidJson(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	jsonBytes := []byte{14, 5, 86}
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(jsonBytes))

	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong")
	}
}

func TestRegiststerHandlerPasswordsNotMatch(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	data := api.RegisterReqData{
		Username:             "test",
		Password:             "sec123",
		PasswordConfirmation: "someothervalue",
	}
	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong")
	}
}

func TestRegiststerHandlerInvalidUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	data := api.RegisterReqData{
		Username:             "inv?lid",
		Password:             "sec123",
		PasswordConfirmation: "sec123",
	}
	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong")
	}
}
