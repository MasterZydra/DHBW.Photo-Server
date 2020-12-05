package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", AuthWrapper(mainHandler))
	log.Fatalln(http.ListenAndServeTLS(":4443", "cert.pem", "key.pem", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	responseString := "<html><body>Hallo</body></html>"
	w.Write([]byte(responseString))
}
