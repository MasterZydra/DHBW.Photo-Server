package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"

	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/user"
	"DHBW.Photo-Server/internal/util"
)

func main() {
	// Setup
	err := util.CheckExistAndCreateFolder(dhbwphotoserver.ImageDir())
	if err != nil {
		log.Fatalf("Error creating image folder: %v", err)
	}

	port := "3000"

	// API endpoint to register a new user
	http.HandleFunc("/register", mustParam(registerHandler, http.MethodPost))

	// returns thumbnails from index to length of currently authenticated user
	http.HandleFunc("/thumbnails", user.AuthHandlerWrapper(
		user.AuthHandler(),
		mustParam(thumbnailsHandler, http.MethodGet, "index", "length"),
	))

	// uploads an image to the users image folder and generates a thumbnail
	http.HandleFunc("/upload", user.AuthHandlerWrapper(
		user.AuthHandler(),
		mustParam(uploadHandler, http.MethodPost),
	))

	// returns one image object with it's information
	http.HandleFunc("/image", user.AuthHandlerWrapper(
		user.AuthHandler(),
		mustParam(imageHandler, http.MethodGet, "name"),
	))

	log.Println("backend listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

// The mustParam wrapper function is used to check if the correct HTTP method is used (POST or GET)
// on the current API endpoint.
// It also checks if all necessary parameters params are provided on a GET request.
func mustParam(h http.HandlerFunc, method string, params ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check if method is allowed
		if r.Method != method {
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

// Register a new user with provided username, password and passwordConfirmation parameters.
// These parameters are sent in a POST request as JSON.
// It is decoded from JSON into the struct api.RegisterReqData.
// After that it checks if the parameters password and passwordConfirmation match.
// Once this is done it can execute the um.Register() function
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

	res.Images = GetThumbnails(strings.ToLower(username), index, length)
	res.TotalImages = GetTotalImages(username)
	return
}

// Upload a new image. The image has to be base64 encoded in the JSON struct.
// The name of the image and the creation date will also be sent via the JSON struct.
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
	err = UploadImage(strings.ToLower(username), data.Filename, data.CreationDate, imageBytes)
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
	res.Image = GetImage(strings.ToLower(username), imgName)
	return
}
