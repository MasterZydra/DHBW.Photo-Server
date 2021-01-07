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

// TODO: Jones: Documentation
// TODO: Jones: Test

func UpdateOrderListEntry(w http.ResponseWriter, r *http.Request) {
	// update order list entry
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

func DownloadOrderList(w http.ResponseWriter, r *http.Request) {
	req, err := NewGetRequest("downloadOrderList")
	if err != nil {
		internalServerError(w, err)
		return
	}

	var res api.DownloadOrderListResData
	err = CallApi(r, req, &res)
	if err != nil {
		badRequest(w, err)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		internalServerError(w, errors.New("Username couldn't be retrieved from Basic auth"))
		return
	}
	zipFileName := "order-list-" + username + "-download.zip"

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
