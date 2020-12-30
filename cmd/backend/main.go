package main

import (
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/auth"
	"DHBW.Photo-Server/internal/user"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	port := "3000"

	// TODO: mustParams-Wrapper einbauen? https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/images", auth.Wrapper(auth.Authenticate(), imagesHandler, []string{}))

	log.Println("backend listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	responseString := "<html><body>Hallo</body></html>"
	w.Write([]byte(responseString))
}

func decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func encode(w http.ResponseWriter, v interface{}) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var res api.RegisterRes

	// decode data
	var data api.RegisterReq
	err := decode(r, &data)
	if err != nil {
		res.Error = err.Error()
	}

	// check if confirmed password is correct
	if data.Password != data.PasswordConfirmation {
		res.Error = "The passwords do not match"
	}

	// execute register function
	um := user.NewUsersManager()
	err = um.Register(data.Username, data.Password)
	if err != nil {
		res.Error = err.Error()
	}

	if err == nil {
		res.Username = data.Username
	}

	// encode data
	// TODO: in Wrapper packen?
	_ = encode(w, &res)
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	responseString := "<html><body>Erfolreich Basic Auth bezwungen</body></html>"
	w.Write([]byte(responseString))
}
