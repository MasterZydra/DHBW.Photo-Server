package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
}

func TestMustParamWrongMethod(t *testing.T) {
	server := createServer(MustParam(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}), http.MethodPost))

	response, _ := http.Get(server.URL)
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Wrong status code")
	}
}

func TestMustParamMissingGetParams(t *testing.T) {
	server := createServer(MustParam(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}), http.MethodGet, "param1", "param2"))

	response, _ := http.Get(server.URL)
	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestMustParamWithGetParams(t *testing.T) {
	server := createServer(MustParam(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}), http.MethodGet, "param1", "param2"))

	response, _ := http.Get(server.URL + "?param1=some&param2=thing")
	if response.StatusCode != http.StatusOK {
		t.Error("Wrong status code")
	}
}
