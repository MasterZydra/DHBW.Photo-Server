package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func UploadPhoto(input UserInput, img *os.File) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", input.Host + "/img", img)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(input.Username, input.Password)

	resp, err := client.Do(req)

	fmt.Println(resp.Status)
}

func UploadPhotos(input UserInput, imgs []*os.File) {
	// TODO: implement
}
