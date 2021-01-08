package main

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func resetUsersFile() {
	csvFile, err := os.Create(DHBW_Photo_Server.UsersFile())
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(csvFile)
	var data = [][]string{
		{DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Hash},
		{DHBW_Photo_Server.User2Name, DHBW_Photo_Server.Pw2Hash},
	}
	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}
}

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
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

func decodeJson(response *http.Response, result interface{}) error {
	// get jsonString from api response
	jsonBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// decode data from jsonUtil into result struct
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return err
	}
	return nil
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

func TestRegisterHandler(t *testing.T) {
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetUsersFile()
	defer resetUsersFile()
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

func TestRegisterHandlerPasswordsNotMatch(t *testing.T) {
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

	var res api.RegisterResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "The passwords do not match" {
		t.Error("Status code wrong")
	}
}

func TestRegisterHandlerInvalidUsername(t *testing.T) {
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

func TestThumbnailsHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=0&length=25", nil)
	req.SetBasicAuth("Max", "pw")
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.TotalImages != 2 || res.Images[0].Name != "img1.jpg" || res.Images[1].Name != "img2.jpg" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestThumbnailsHandlerNoUsername(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=0&length=25", nil)
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong error message")
	}
}

func TestThumbnailsHandlerInvalidIndex(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=wrong&length=25", nil)
	req.SetBasicAuth("Max", "pw")
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Invalid index. Index must be an Integer" {
		t.Error("Status code wrong or wrong error message")
	}
}

func TestThumbnailsHandlerInvalidLength(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=0&length=wrong", nil)
	req.SetBasicAuth("Max", "pw")
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Invalid length. Length must be an Integer" {
		t.Error("Status code wrong or wrong error message")
	}
}
