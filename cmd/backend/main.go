package main

import (
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/auth"
	"DHBW.Photo-Server/internal/user"
	"log"
	"net/http"
)

func main() {
	port := "3000"

	// TODO: mustParams-Wrapper einbauen? https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
	// ermöglicht Registrierung
	http.HandleFunc("/register", registerHandler)

	// gibt Thumbnails mit den Infos dazu von index bis length zurück
	http.HandleFunc("/thumbnails", auth.Wrapper(auth.Authenticate(), thumbnailsHander))

	// lädt Image hoch
	http.HandleFunc("/upload", auth.Wrapper(auth.Authenticate(), uploadHandler))

	// TODO: Jones: klären, ob nur Pfad oder das Bild an sich -> andere sollen nicht Link erraten können (gilt auch für Thumbnails)
	// Gibt Bild + Infos zurück
	http.HandleFunc("/image", auth.Wrapper(auth.Authenticate(), imageHandler))

	log.Println("backend listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
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
	um := user.NewUserManager()
	err = um.Register(data.Username, data.Password)
	if err != nil {
		res.Error = err.Error()
		return
	}
}

func thumbnailsHander(w http.ResponseWriter, r *http.Request) {
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
	// TODO: David: implementieren
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: David: implementieren
}
