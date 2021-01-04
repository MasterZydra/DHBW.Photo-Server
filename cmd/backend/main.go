package main

import (
	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/user"
	"DHBW.Photo-Server/internal/util"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/register", mustParam(registerHandler, http.MethodPost))

	// gibt Thumbnails mit den Infos dazu von index bis length zurück
	http.HandleFunc("/thumbnails", user.AuthHandlerWrapper(
		user.AuthHandler(),
		mustParam(thumbnailsHandler, http.MethodGet, "index", "length"),
	))

	// lädt Image hoch
	http.HandleFunc("/upload", user.AuthHandlerWrapper(
		user.AuthHandler(),
		mustParam(uploadHandler, http.MethodPost),
	))

	// Gibt Bild + Infos zurück
	http.HandleFunc("/image", user.AuthHandlerWrapper(
		user.AuthHandler(),
		mustParam(imageHandler, http.MethodGet, "name"),
	))

	log.Println("backend listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func mustParam(h http.HandlerFunc, method string, params ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			//w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "Method "+r.Method+" not allowed here", http.StatusMethodNotAllowed)
			return
		}

		// check for missing GET parameters
		if method == http.MethodGet && len(params) > 0 {
			data := r.URL.Query()
			for _, param := range params {
				if len(data.Get(param)) == 0 {
					http.Error(w, "missing GET parameter: "+param, http.StatusBadRequest)
					return
				}
			}
		}

		h(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var res api.RegisterResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.RegisterReqData
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
	var res api.ThumbnailsResData
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	// Check if parameter "index" is an integer
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		res.Error = "Invalid index. Index must be an Integer"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Check if parameter "length" is an integer
	length, err := strconv.Atoi(r.URL.Query().Get("length"))
	if err != nil {
		res.Error = "Invalid length. Length must be an Integer"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Load thumbnails from associated ImageManager
	res.Images = GetThumbnails(username, index, length)
	res.TotalImages = GetTotalImages(username)
	return
}

// Upload a new image. The image has to be base64 encoded in the json struct.
// The name of the image and the creation date will also be sent via the json struct.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var res api.UploadResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.UploadReqData
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		return
	}

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(data.Base64Image)
	if err != nil {
		res.Error = err.Error()
		return
	}

	// Save image for associated user
	err = UploadImage(username, data.Filename, data.CreationDate, imageBytes)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Request details to an image. The result is a JSON object.
// The request need to be send via the GET method and contain the parameter name which is the file name.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ImageResData
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		return
	}

	// Get parameter "name" from url
	imgName := r.URL.Query().Get("name")

	// Load image from associated ImageManager
	res.Image = GetImage(username, imgName)
	return
}
