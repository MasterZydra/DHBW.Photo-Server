package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const port = "4443"

type templateData struct {
	BackendHost string
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", serverTemplate)

	log.Println("web listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func serverTemplate(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	layout := filepath.Join(wd, "templates", "layout.html")
	publicFile := filepath.Join(wd, "public", filepath.Clean(r.URL.Path))

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

	err = tmpl.ExecuteTemplate(w, "layout", getTemplateData())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func getTemplateData() templateData {
	return templateData{
		BackendHost: "https://localhost:3000/",
	}
}
