package main

import (
	"io/ioutil"
	"log"
	"net/http"

	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/auth"
	"DHBW.Photo-Server/internal/user"
	"DHBW.Photo-Server/internal/util"
)

// TODO: Jones Documentation

func main() {
	// Setup
	err := util.CheckExistAndCreateFolder(DHBW_Photo_Server.ImageDir)
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

func thumbnailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var res api.ThumbnailsRes
	defer jsonUtil.EncodeResponse(w, &res)

	// decode data
	var data api.ThumbnailsReq
	err := jsonUtil.DecodeBody(r, &data)
	if err != nil {
		res.Error = err.Error()
		return
	}

	//username, _, ok := r.BasicAuth()
	//if !ok {
	//	res.Error = "Could not get username"
	//	return
	//}

	// TODO: David: implementieren
}

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

	// TODO: David imagecreationdate optional machen
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

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// TODO: David: implementieren
}
