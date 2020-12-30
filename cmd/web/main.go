package main

import (
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/auth"
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

var whitelistPaths = []string{"/", "/index.html", "/register.html"}

type TemplateVariables struct {
	Global GlobalVariables
	Local  interface{}
}

type GlobalVariables struct {
	Username string
	LoggedIn bool
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/home", auth.Wrapper(auth.Authenticate(), homeHandler))

	log.Println("web listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func layout(w http.ResponseWriter, r *http.Request, data interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		internalServerError(w, err)
		return
	}
	dir := filepath.Join(wd, "cmd", "web")
	layout := filepath.Join(dir, "templates", "layout.html")
	publicFile := filepath.Join(dir, "public", filepath.Clean(r.URL.Path)+".html")

	// check if site exists or is a directory
	siteStat, err := os.Stat(publicFile)
	if err != nil && os.IsNotExist(err) || siteStat.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(layout, publicFile)
	if err != nil {
		internalServerError(w, err)
		return
	}

	// set template variables
	username, _, ok := r.BasicAuth()
	loggedIn := true
	if !ok {
		username = ""
		loggedIn = false
	}
	templateVars := TemplateVariables{
		Global: GlobalVariables{Username: username, LoggedIn: loggedIn},
		Local:  data,
	}

	// execute template and send data with it
	err = tmpl.ExecuteTemplate(w, "layout", templateVars)
	if err != nil {
		internalServerError(w, err)
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "index"
	layout(w, r, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data := api.RegisterReq{
			Username:             r.FormValue("user"),
			Password:             r.FormValue("password"),
			PasswordConfirmation: r.FormValue("passwordConfirmation"),
		}
		var res api.RegisterRes
		err := callApi(r, "register", data, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if res.Error != "" {
			http.Error(w, res.Error, http.StatusBadRequest)
			return
		}

		//logout(w, r)
		http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
		return
	}
	layout(w, r, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ImageRes
	_ = callApi(r, "images", nil, &res)

	layout(w, r, res)
}

func internalServerError(w http.ResponseWriter, error error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+error.Error(), http.StatusInternalServerError)
}

func callApi(r *http.Request, url string, data interface{}, res interface{}) error {
	// encode data to json
	jsonBytes, err := json.Marshal(data)

	// prepare request
	req, err := http.NewRequest(http.MethodPost, BackendHost+url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	// set basic auth for backend request if available
	username, pw, ok := r.BasicAuth()
	if ok {
		req.SetBasicAuth(username, pw)
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
