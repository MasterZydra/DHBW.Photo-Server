package main

import (
	"log"
	"net/http"
)

func main() {
	port := "3000"

	http.HandleFunc("/", mainHandler)
	log.Println("backend listening on " + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	responseString := "<html><body>Hallo</body></html>"
	w.Write([]byte(responseString))
}
