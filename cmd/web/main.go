package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := "4443"

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", serverTemplate)

	log.Println("web listening on " + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func serverTemplate(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "cmd", "web")
	layout := filepath.Join(dir, "templates", "layout.html")
	publicFile := filepath.Join(dir, "public", filepath.Clean(r.URL.Path))

	siteStat, err := os.Stat(publicFile)
	if err != nil && os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	if siteStat.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(layout, publicFile)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
