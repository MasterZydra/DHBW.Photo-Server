package main

import (
	"DHBW.Photo-Server/internal/api"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// outsourced function from web/main.go to update an order list entry
func UpdateOrderListEntry(w http.ResponseWriter, r *http.Request) {
	numberOfPrints, err := strconv.Atoi(r.FormValue("numberOfPrints"))
	if err != nil {
		numberOfPrints = 1
	}
	data := api.ChangeOrderListEntryReqData{
		ImageName:      r.FormValue("imageName"),
		Format:         r.FormValue("format"),
		NumberOfPrints: numberOfPrints,
	}

	req, err := NewPostRequest("changeOrderListEntry", data)
	if err != nil {
		internalServerError(w, err)
		return
	}

	var res api.ChangeOrderListEntryResData
	err = CallApi(r, req, &res)
	if err != nil {
		badRequest(w, err)
		return
	}
}

// outsourced function from web/main.go to remove an order list entry
// redirects with a GET request to /order-list afterwards
func RemoveOrderListEntry(w http.ResponseWriter, r *http.Request, imageToRemove string) {
	data := api.RemoveOrderListEntryReqData{ImageName: imageToRemove}
	req, err := NewPostRequest("removeOrderListEntry", data)
	if err != nil {
		internalServerError(w, err)
		return
	}

	var res api.RemoveOrderListEntryResData
	err = CallApi(r, req, &res)
	if err != nil {
		badRequest(w, err)
		return
	}
	http.Redirect(w, r, "/order-list", http.StatusFound)
}

// outsourced function from web/main.go to delete a order list
// redirects with a GET request to /order-list afterwards
func DeleteOrderList(w http.ResponseWriter, r *http.Request) {
	req, err := NewPostRequest("deleteOrderList", nil)
	if err != nil {
		internalServerError(w, err)
		return
	}

	var res api.DeleteOrderListResData
	err = CallApi(r, req, &res)
	if err != nil {
		badRequest(w, err)
		return
	}
	http.Redirect(w, r, "/order-list", http.StatusFound)
}

// outsourced function from web/main.go to download the order list of the current user as a .zip file
// The file will be retrieved as base64 encoded string and converted to a real file.
// After that the file can be served to the browser and deleted afterwards.
func DownloadOrderList(w http.ResponseWriter, r *http.Request) {
	req, err := NewGetRequest("downloadOrderList")
	if err != nil {
		internalServerError(w, err)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		internalServerError(w, errors.New("Username couldn't be retrieved from Basic auth"))
		return
	}
	zipFileName := "order-list-" + username + "-download.zip"

	var res api.DownloadOrderListResData
	err = CallApi(r, req, &res)
	if err != nil {
		badRequest(w, err)
		return
	}

	// convert backend result back to bytes to write the zipfile as a file
	zipFileBytes, err := base64.StdEncoding.DecodeString(res.Base64ZipFile)
	if err != nil {
		internalServerError(w, err)
		return
	}

	err = ioutil.WriteFile(zipFileName, zipFileBytes, 0755)
	if err != nil {
		internalServerError(w, err)
		return
	}

	// download it to browser
	http.ServeFile(w, r, zipFileName)

	err = os.Remove(zipFileName)
	if err != nil {
		internalServerError(w, err)
		return
	}
	return
}

// Outsourced function from web/main.go to add all chosen images as an order list entry.
// Redirects with a GET request to /order-list afterwards.
func AddOrderListEntries(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		internalServerError(w, err)
		return
	}
	if len(r.PostForm) > 0 {
		for _, value := range r.PostForm {
			data := api.AddOrderListEntryReqData{ImageName: value[0]}

			req, err := NewPostRequest("addOrderListEntry", data)
			if err != nil {
				internalServerError(w, err)
				return
			}

			var res api.AddOrderListEntryResData
			err = CallApi(r, req, &res)
			if err != nil {
				badRequest(w, err)
				return
			}
		}
		http.Redirect(w, r, "/order-list", http.StatusFound)
		return
	}
}
