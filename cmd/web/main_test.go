package main

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
}

func newPostReq(url string, data interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return newReq(http.MethodPost, url, jsonBytes)
}

func newGetReq(url string) (*http.Request, error) {
	return newReq(http.MethodGet, url, nil)
}

func newReq(method string, url string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("username", "password")
	return req, nil
}

func TestCallApiPost(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		var data api.TestReqData
		_ = jsonUtil.DecodeBody(r, &data)

		res := api.TestResData{
			Error:            "",
			SomeResultString: data.SomeString + "server",
			SomeResultInt:    data.SomeInt + 1,
		}
		_ = jsonUtil.EncodeResponse(w, &res)
	})
	defer server.Close()

	reqData := api.TestReqData{
		SomeString: "request",
		SomeInt:    12,
	}
	req, _ := newPostReq(server.URL, reqData)

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err != nil || res.SomeResultString != reqData.SomeString+"server" || res.SomeResultInt != reqData.SomeInt+1 {
		t.Error("Error occurred or result data of POST request wrong.")
	}
}

func TestCallApiGet(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		someString := r.URL.Query().Get("someString")
		res := api.TestResData{
			Error:            "",
			SomeResultString: someString + "server",
			SomeResultInt:    123,
		}
		_ = jsonUtil.EncodeResponse(w, &res)
	})

	payload := url.Values{
		"someString": {"test"},
	}
	req, _ := newGetReq(server.URL + "?" + payload.Encode())

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err != nil || res.SomeResultString != "testserver" || res.SomeResultInt != 123 {
		t.Error("Error occurred or result data of GET request wrong")
	}
}

func TestCallApiFailDo(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
	})
	defer server.Close()

	req, _ := newGetReq("invalid://url.com")

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err == nil {
		t.Error("No error occurred")
	}
}

func TestCallApiFailUnmarshalJson(t *testing.T) {
	// nicht testbar, da es sehr schwierig ist ungültige Daten in json.Unmarshal über den Server reinzuschicken
}

func TestCallApiCustomError(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		res := api.TestResData{
			Error: "customServerError",
		}
		_ = jsonUtil.EncodeResponse(w, &res)
	})

	req, _ := newGetReq(server.URL)

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err == nil || err.Error() != "customServerError" {
		t.Error("Error occurred or result data of GET request wrong")
	}
}

func TestNewRequest(t *testing.T) {
	req, err := NewRequest("POST", "some/path", nil)
	if err != nil || req.Method != "POST" || req.URL.Path != "/some/path" || !strings.Contains(DHBW_Photo_Server.BackendHost, req.URL.Host) {
		t.Error("Error while creating new request")
	}
}

func TestNewGetRequest(t *testing.T) {
	req, err := NewGetRequest("some/path")
	if err != nil || req.Method != "GET" {
		t.Error("Error while creating GET request")
	}
}

func TestNewPostRequestValidJson(t *testing.T) {
	data := api.TestReqData{
		SomeString: "test",
		SomeInt:    1,
	}
	req, err := NewPostRequest("some/path", data)
	if err != nil || req.Method != "POST" {
		t.Error("Error while creating POST request")
	}
}

func TestNewPostRequestInvalidJson(t *testing.T) {
	// nicht testbar, da es sehr schwierig ist ungültige Daten in json.Marshal reinzuschicken
}
