package main

import (
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const port = "4443"
const BackendHost = "https://localhost:3000/"

type templateData struct {
	BackendHost string
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", serverTemplate)

	// backend call paths
	http.HandleFunc("/register", registerHandler)

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
		BackendHost,
	}
}

func callApi(url string, data interface{}, res interface{}) error {
	// encode data to json
	jsonBytes, err := json.Marshal(data)

	// prepare request
	req, err := http.NewRequest(http.MethodPost, BackendHost+url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	// skip certificate verification to avoid error: "x509: certificate signed by unknown authority"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := http.Client{Transport: tr}

	// execute API request
	apiRes, err := c.Do(req)
	if err != nil {
		return err
	}

	// get jsonString from api response
	jsonBytes, err = ioutil.ReadAll(apiRes.Body)
	if err != nil {
		return err
	}

	// decode data from json into result struct
	err = json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return err
	}
	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := api.RegisterReq{
		Username:             r.FormValue("user"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("passwordConfirmation"),
	}
	var res api.RegisterRes
	err := callApi("register", data, &res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if res.Error != "" {
		http.Error(w, res.Error, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
}
