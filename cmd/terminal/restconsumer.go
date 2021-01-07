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
	rawImg, err := ioutil.ReadAll(img)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	data := api.UploadReqData{
		Base64Image:  base64.StdEncoding.EncodeToString(rawImg),
		Filename:     filepath.Base(img.Name()),
		CreationDate: time.Now().Local(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	req, err := http.NewRequest("POST", input.Host + "/upload", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	req.SetBasicAuth(input.Username, input.Password)

	// skip certificate verification to avoid error: "x509: certificate signed by unknown authority"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	jsonResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	var backendRes api.UploadResData
	err = json.Unmarshal(jsonResponse, &backendRes)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	resError := backendRes.GetError()
	if resError != "" {
		fmt.Println(errors.New(resError))
	}
	wg.Done()
}

func UploadPhotos(input UserInput, imgs []*os.File) {
	// TODO: test
	wg := sync.WaitGroup{}
	for _, img := range imgs {
		wg.Add(1)
		go UploadPhoto(input, img, &wg)
	}
	wg.Wait()
}
