package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
)

func newTestServer () *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		username, _, err := req.BasicAuth()
		if !err || username == "" {
			rw.Write([]byte(`{"Error":"Could not get username"}`))
			return
		}
		// Send valid response to be tested
		rw.Write([]byte(`{"Error":""}`))
	}))
}

func prepareTests() (*httptest.Server, UserInput, []*os.File) {
	testPath := "../../test"
	files := ReadJPEGsFromPath(testPath)
	if files == nil || len(files) < 1 {
		log.Fatal("Error while reading test image")
	}

	// Start a local HTTP server
	server := newTestServer()

	input := UserInput{
		Username: "testuser",
		Password: "12345",
		Host:     server.URL,
		Path:     testPath,
	}
	return server, input, files
}

func ExampleUploadPhoto() {
	server, input, files := prepareTests()

	wg := sync.WaitGroup{}
	wg.Add(1)
	UploadPhoto(input, files[0], &wg)
	wg.Wait()
	server.Close()

	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
}

func ExampleUploadPhotoInvalidImg()  {
	_, input, _ := prepareTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send invalid response
		rw.Write([]byte(`Test`))
	}))
	input.Host = server.URL

	wg := sync.WaitGroup{}
	wg.Add(1)
	UploadPhoto(input, nil, &wg)
	wg.Wait()
	server.Close()

	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
	// invalid argument
}

func ExampleUploadPhotoErrorDoingReq()  {
	server, input, files := prepareTests()
	input.Host = ""

	wg := sync.WaitGroup{}
	wg.Add(1)
	UploadPhoto(input, files[0], &wg)
	wg.Wait()
	server.Close()

	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
	// Post "/upload": unsupported protocol scheme ""
}

func ExampleUploadPhotoInvalidJsonRes() {
	_, input, files := prepareTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send invalid response
		rw.Write([]byte(`Test`))
	}))
	input.Host = server.URL

	wg := sync.WaitGroup{}
	wg.Add(1)
	UploadPhoto(input, files[0], &wg)
	wg.Wait()
	server.Close()

	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
	// invalid character 'T' looking for beginning of value
}

func ExampleUploadPhotoWithoutUname() {
	server, input, files := prepareTests()
	input.Username = ""

	wg := sync.WaitGroup{}
	wg.Add(1)
	UploadPhoto(input, files[0], &wg)
	wg.Wait()
	server.Close()

	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
	// Could not get username
}

func ExampleUploadPhotos() {
	server, input, files := prepareTests()

	UploadPhotos(input, files)
	server.Close()

	// Output:
	// --------------
	// Folder successfully read in
	// --------------
	//
	// --------------
	// Successfully read 1 files
	// --------------
}