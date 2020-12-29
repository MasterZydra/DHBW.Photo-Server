package main

import (
	"DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/user"
	"DHBW.Photo-Server/internal/util"
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

var whitelistPaths = []string{"/", "/login.html", "/register.html"}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", serverTemplate)

	// backend call paths
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	log.Println("web listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

// TODO: auslagern + tests?
func isLoggedIn(cookie *http.Cookie) (bool, error) {
	um := user.NewUsersManager()
	err := um.LoadUsers()
	if err != nil {
		return false, err
	}
	userObj := um.GetUserByCookie(cookie)
	if userObj != nil && cookie.Value == userObj.Cookie.Value {
		return true, nil
	}
	return false, nil
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func serverTemplate(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	// redirect to login if (no cookie available OR not logged in) AND not a whitelist path
	if !util.ContainsString(whitelistPaths, urlPath) {
		sentCookie, err := r.Cookie(DHBW_Photo_Server.CookieName)
		if sentCookie == nil || err != nil {
			redirectToLogin(w, r)
		} else {
			// redirect to login if cookie is not valid
			loggedIn, err := isLoggedIn(sentCookie)
			if err != nil {
				internalServerError(w, err)
				return
			}
			if !loggedIn {
				redirectToLogin(w, r)
			}
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		internalServerError(w, err)
	}
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
		return
	}

	//err = tmpl.ExecuteTemplate(w, "layout", getTemplateData())
	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err.Error())
		internalServerError(w, err)
		return
	}
}

func internalServerError(w http.ResponseWriter, error error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+error.Error(), http.StatusInternalServerError)
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
	// TODO: wrapper? Sehr viel doppelt in den Handler-Funktionen
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := api.LoginReq{
		Username: r.FormValue("user"),
		Password: r.FormValue("password"),
	}
	var res api.LoginRes
	err := callApi("login", data, &res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if res.Error != "" {
		http.Error(w, res.Error, http.StatusBadRequest)
		return
	}
	if res.Cookie.Name != "" && res.Cookie.Value != "" {
		// set Cookie in browser
		http.SetCookie(w, &res.Cookie)
	}

	http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
}
