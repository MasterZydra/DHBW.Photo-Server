/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import (
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func UploadPhoto(input UserInput, img *os.File, wg *sync.WaitGroup) {
	// read data from image-file
	rawImg, err := ioutil.ReadAll(img)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	// create data for request body as defined in apistructs.go
	// image data will be base64-encoded
	// the filename is taken from the *os.file
	data := api.UploadReqData{
		Base64Image:  base64.StdEncoding.EncodeToString(rawImg),
		Filename:     filepath.Base(img.Name()),
		CreationDate: time.Now().Local(),
	}

	// json encoding of data
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	// create a new request with method POST
	// the url is the host entered by the user linked to the endpoint /upload
	// as body the json encoded data is set
	req, err := http.NewRequest("POST", input.Host+"/upload", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	// set the request's authorization header, use the username and password entered by the user
	req.SetBasicAuth(input.Username, input.Password)

	// skip certificate verification to avoid error: "x509: certificate signed by unknown authority"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// create a new http client and send the request
	client := http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	// read the response body (it should be a json with the fields as defined in apistructs.go)
	jsonResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	// parses JSON-encoded response
	var backendRes api.UploadResData
	err = json.Unmarshal(jsonResponse, &backendRes)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	// if an error was returned in the response, it will be displayed to the user
	resError := backendRes.GetError()
	if resError != "" {
		fmt.Println(errors.New(resError))
	}
	// the counter is decreased by one
	wg.Done()
}

func UploadPhotos(input UserInput, imgs []*os.File) {
	// a WaitGroup waits for a collection of goroutines to finish
	wg := sync.WaitGroup{}
	for _, img := range imgs {
		// for each goroutine the counter is increased by one
		wg.Add(1)
		go UploadPhoto(input, img, &wg)
	}
	// wait until the counter is zero again
	wg.Wait()
}
