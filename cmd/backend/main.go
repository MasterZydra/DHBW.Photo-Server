/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import (
	"DHBW.Photo-Server/internal/auth"
	"DHBW.Photo-Server/internal/util"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	dhbwphotoserver "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/user"
)

func main() {
	// Read flags
	port := flag.Int("port", dhbwphotoserver.BackendDefaultPort, "Port the server listens on")
	imagePath := flag.String("imagepath", dhbwphotoserver.DefaultImageDir, "Path of image directory")
	flag.Parse()

	// Convert flags into right format
	portStr := strconv.Itoa(*port)
	dhbwphotoserver.SetImageDir(*imagePath)

	// Setup
	err := util.CheckExistAndCreateFolder(dhbwphotoserver.ImageDir())
	if err != nil {
		log.Fatalf("Error creating image folder: %v", err)
	}

	// API endpoint to register a new user
	http.HandleFunc("/register", MustParam(RegisterHandler, http.MethodPost))

	// returns thumbnails from index to length of currently authenticated user
	http.HandleFunc("/thumbnails", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(ThumbnailsHandler, http.MethodGet, "index", "length"),
	))

	// uploads an image to the users image folder and generates a thumbnail
	http.HandleFunc("/upload", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(UploadHandler, http.MethodPost),
	))

	// returns one image object with it's information
	http.HandleFunc("/image", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(imageHandler, http.MethodGet, "name"),
	))

	// GET request that returns all order list entries
	http.HandleFunc("/orderList", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(OrderListEntryHandler, http.MethodGet),
	))

	// adds a new entry to the order list of the current user with the passed ImageName
	http.HandleFunc("/addOrderListEntry", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(AddOrderListEntryHandler, http.MethodPost),
	))

	// change parameter of a order list entry (e.g. Format or NumberOfPrints
	http.HandleFunc("/changeOrderListEntry", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(ChangeOrderListEntryHandler, http.MethodPost),
	))

	// remove a entry from the order list of the current user
	http.HandleFunc("/removeOrderListEntry", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(RemoveOrderListEntryHandler, http.MethodPost),
	))

	// delete the complete order list
	http.HandleFunc("/deleteOrderList", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(DeleteOrderListHandler, http.MethodPost),
	))

	// download the current order list as a zip file with a content.json containing it's metadata
	http.HandleFunc("/downloadOrderList", auth.AuthHandlerWrapper(
		auth.AuthHandler(),
		MustParam(DownloadOrderList, http.MethodGet),
	))

	log.Println("backend listening on https://localhost:" + portStr)
	log.Fatalln(http.ListenAndServeTLS(":"+portStr, "cert.pem", "key.pem", nil))
}

// The MustParam wrapper function is used to check if the correct HTTP method is used (POST or GET)
// on the current API endpoint.
// It also checks if all necessary parameters params are provided on a GET request.
func MustParam(h http.HandlerFunc, method string, params ...string) http.HandlerFunc {
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
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var res api.RegisterResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.RegisterReqData
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if confirmed password is correct
	if data.Password != data.PasswordConfirmation {
		res.Error = "The passwords do not match"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// execute register function
	um := user.UserManagerCache()
	err = um.Register(data.Username, data.Password)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Request a range of thumbnails. The result is a JSON object.
// The request need to be send via the GET method and contain the parameters
// index and length. Both parameters have to be integers.
func ThumbnailsHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ThumbnailsResData
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if parameter "index" is an integer
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		res.Error = "Invalid index. Index must be an Integer"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if parameter "length" is an integer
	length, err := strconv.Atoi(r.URL.Query().Get("length"))
	if err != nil {
		res.Error = "Invalid length. Length must be an Integer"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Load thumbnails from associated ImageManager
	res.Images = GetThumbnails(strings.ToLower(username), index, length)
	res.TotalImages = GetTotalImages(username)
	return
}

// Upload a new image. The image has to be base64 encoded in the JSON struct.
// The name of the image and the creation date will also be sent via the JSON struct.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var res api.UploadResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.UploadReqData
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(data.Base64Image)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Save image for associated user
	err = UploadImage(strings.ToLower(username), data.Filename, data.CreationDate, imageBytes)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// TODO: Jones: tests
// ATTENTION!
// There are NO TESTS for this function because this API endpoint would be used for the comment feature which is not implemented.
// Request details to an image. The result is a JSON object.
// The request need to be send via the GET method and contain the parameter name which is the file name.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ImageResData
	defer jsonUtil.EncodeResponse(w, &res)

	// Get username from basic authentication
	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get parameter "name" from url
	imgName := r.URL.Query().Get("name")

	// Load image from associated ImageManager
	res.Image = GetImage(strings.ToLower(username), imgName)
	return
}

// Gets all order list entries and writes it in OrderListResData.OrderList
func OrderListEntryHandler(w http.ResponseWriter, r *http.Request) {
	var res api.OrderListResData
	defer jsonUtil.EncodeResponse(w, &res)

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	um := user.UserManagerCache()
	usr := um.GetUser(username)

	res.OrderList = usr.OrderList.Entries
}

// API endpoint for adding a new order list entry by passing a name of an image
func AddOrderListEntryHandler(w http.ResponseWriter, r *http.Request) {
	var res api.AddOrderListEntryResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.AddOrderListEntryReqData
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img := GetImage(username, data.ImageName)
	if img == nil {
		res.Error = "Could not get image '" + data.ImageName + "'"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	um := user.UserManagerCache()
	usr := um.GetUser(username)

	err = usr.OrderList.AddOrderEntry(img)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Removes a specific order list entry with an imageName.
func RemoveOrderListEntryHandler(w http.ResponseWriter, r *http.Request) {
	var res api.RemoveOrderListEntryResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.RemoveOrderListEntryReqData
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	um := user.UserManagerCache()
	usr := um.GetUser(username)

	exists := usr.OrderList.RemoveOrderEntry(data.ImageName)
	if !exists {
		res.Error = "Image '" + data.ImageName + "' does not exist for user '" + usr.Name + "'"
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Modifies the Format and NumberOfPrints of an order list entry.
func ChangeOrderListEntryHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ChangeOrderListEntryResData
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.ChangeOrderListEntryReqData
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	um := user.UserManagerCache()
	usr := um.GetUser(username)

	_, entry := usr.OrderList.GetOrderEntry(data.ImageName)
	if entry == nil {
		res.Error = "Could not find order entry with image '" + data.ImageName + "'"
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	entry.Format = data.Format
	entry.NumberOfPrints = data.NumberOfPrints
}

// Empties the complete order list array of the current user
func DeleteOrderListHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ChangeOrderListEntryResData
	defer jsonUtil.EncodeResponse(w, &res)

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	um := user.UserManagerCache()
	usr := um.GetUser(username)

	usr.OrderList = &user.OrderList{Entries: []*user.OrderListEntry{}}
}

// Calls user.CreateOrderListZipFile and returns the zip files content to base64 encoded string.
func DownloadOrderList(w http.ResponseWriter, r *http.Request) {
	var res api.DownloadOrderListResData
	defer jsonUtil.EncodeResponse(w, &res)

	username, _, ok := r.BasicAuth()
	if !ok {
		res.Error = "Could not get username"
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	um := user.UserManagerCache()
	usr := um.GetUser(username)

	zipFileName := filepath.Join(dhbwphotoserver.ImageDir(), usr.Name, "order-list-download.zip")
	err := user.CreateOrderListZipFile(zipFileName, usr.Name, usr.OrderList)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	zipFileBytes, err := ioutil.ReadFile(zipFileName)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Base64ZipFile = base64.StdEncoding.EncodeToString(zipFileBytes)

	err = os.Remove(zipFileName)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
