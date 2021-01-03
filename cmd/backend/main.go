package main

import (
	"DHBW.Photo-Server/internal/user"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
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
	// ermöglicht Registrierung
	http.HandleFunc("/register", allowMethod(http.MethodPost, registerHandler))

	// gibt Thumbnails mit den Infos dazu von index bis length zurück
	http.HandleFunc("/thumbnails", user.HandlerWrapper(
		user.AuthenticateHandler(),
		allowMethod(http.MethodGet, thumbnailsHandler),
	))

	// lädt Image hoch
	http.HandleFunc("/upload", user.HandlerWrapper(
		user.AuthenticateHandler(),
		allowMethod(http.MethodPost, uploadHandler),
	))

	// Gibt Bild + Infos zurück
	http.HandleFunc("/image", user.HandlerWrapper(
		user.AuthenticateHandler(),
		allowMethod(http.MethodGet, imageHandler),
	))

	log.Println("backend listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func allowMethod(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
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
	um := user.GetImageManager()
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
	var res api.ThumbnailsRes
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	// Get parameter "index" from url
	strIndex := r.URL.Query().Get("index")
	// Check if parameter "index" is given
	if strIndex == "" {
		res.Error = "Parameter index is missing"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Check if parameter "index" is an integer
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		res.Error = "Invalid index. Index must be an Integer"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Get parameter "length" from url
	strlength := r.URL.Query().Get("length")
	// Check if parameter "length" is given
	if strlength == "" {
		res.Error = "Parameter length is missing"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Check if parameter "length" is an integer
	length, err := strconv.Atoi(strlength)
	if err != nil {
		res.Error = "Invalid length. Length must be an Integer"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Load thumbnails from associated ImageManager
	res.Images = GetThumbnail(strings.ToLower(username), index, length)
	return
}

// Upload a new image. The image has to be given in the message body.
// In the header the image name must be given with the key "imagename".
// If known the creation date of the image (read from file system) can be given as key "imagecreationdate".
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ImageUploadRes
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	// Get key "imagename" from header
	imgname := r.Header.Get("imagename")
	// Check if key "imagename" is given
	if imgname == "" {
		res.Error = "Key imagename is missing"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// ToDo: David - Move default value for date in NewUploadImage
	imgcreation := r.Header.Get("imagecreationdate")
	if imgcreation == "" {
		imgcreation = time.Now().Format("2006-01-02")
	}

	// Read image raw data from body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			// ToDo Error treatment
		}
	}

	// Save image for associated user
	errorString := UploadImage(strings.ToLower(username), imgname, imgcreation, body)
	if errorString != "" {
		res.Error = errorString
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Request details to an image. The result is a JSON object.
// The request need to be send via the GET method and contain the parameter name which is the file name.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ImageRes
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	// Get parameter "name" from url
	imgname := r.URL.Query().Get("name")
	// Check if parameter "name" is given
	if imgname == "" {
		res.Error = "Parameter name is missing"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Load image from associated ImageManager
	res.Image = GetImage(strings.ToLower(username), imgname)
	return
}
