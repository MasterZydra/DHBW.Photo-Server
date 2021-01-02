package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/auth"
	"DHBW.Photo-Server/internal/user"
	"DHBW.Photo-Server/internal/util"
)

// TODO: Jones Documentation

func main() {
	// Setup
	err := util.CheckExistAndCreateFolder(dhbwphotoserver.ImageDir)
	if err != nil {
		log.Fatalf("Error creating image folder: %v", err)
	}

	port := "3000"
	// TODO: mustParams-Wrapper einbauen? https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
	// ermöglicht Registrierung
	http.HandleFunc("/register", registerHandler)

	// gibt Thumbnails mit den Infos dazu von index bis length zurück
	http.HandleFunc("/thumbnails", auth.HandlerWrapper(auth.AuthenticateHandler(), thumbnailsHandler))

	// lädt Image hoch
	http.HandleFunc("/upload", auth.HandlerWrapper(auth.AuthenticateHandler(), uploadHandler))

	// Gibt Bild + Infos zurück
	http.HandleFunc("/image", auth.HandlerWrapper(auth.AuthenticateHandler(), imageHandler))

	log.Println("backend listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Jones Method-Prüfung in Wrapper bauen
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var res api.RegisterRes
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.RegisterReq
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		return
	}

	// check if confirmed password is correct
	if data.Password != data.PasswordConfirmation {
		res.Error = "The passwords do not match"
		return
	}

	// execute register function
	um := user.NewUserManager()
	err = um.Register(data.Username, data.Password)
	if err != nil {
		res.Error = err.Error()
		return
	}
}

// Request a range of thumbnails. The result is a JSON object.
// The request need to be send via the GET method and contain the parameters
// index and length. Both parameters have to be integers.
func thumbnailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var res api.ThumbnailsRes
	defer jsonUtil.EncodeResponse(w, &res)

	// ToDo David: Return message if not all paramters are given

	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		res.Error = "Invalid index. Index must be an Integer"
	}

	length, err := strconv.Atoi(r.URL.Query().Get("length"))
	if err != nil {
		res.Error = "Invalid index. Index must be an Integer"
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	res.Images = GetThumbnail(username, index, length)
	return
}

// Upload a new image. The image has to be given in the message body.
// In the header the image name must be given with the key "imagename".
// If known the creation date of the image (read from file system) can be given as key "imagecreationdate".
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var res api.ImageUploadRes
	defer jsonUtil.EncodeResponse(w, &res)

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	// ToDo: David - Move default value for date in NewUploadImage
	imgname := r.Header.Get("imagename")
	imgcreation := r.Header.Get("imagecreationdate")

	// Read body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			// ToDo Error treatment
		}
	}

	errorString := UploadImage(username, imgname, imgcreation, body)
	if errorString != "" {
		res.Error = errorString
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Request details to an image. The result is a JSON object.
// The request need to be send via the GET method and contain the parameter name which is the file name.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var res api.ImageRes
	defer jsonUtil.EncodeResponse(w, &res)

	imgname := r.URL.Query().Get("name")
	fmt.Print("Imagename: %v", imgname)
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	res.Image = GetImage(username, imgname)
	return
	// TODO: David: implementieren
}
