package main

import (
	"DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/user"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// TODO: Jones Documentation

// TODO: Port über go flag konfigurierbar machen können
const port = "4443"

type TemplateVariables struct {
	Global GlobalVariables
	Local  interface{}
}

type GlobalVariables struct {
	Username string
	LoggedIn bool
}

func main() {
	// serve images directory
	fs := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", user.AuthFileServerWrapper(
		user.AuthFileServer(),
		http.StripPrefix("/images", fs),
	))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/home", user.AuthHandlerWrapper(user.AuthHandler(), homeHandler))
	// TODO: jones: uploadHandler implementieren
	http.HandleFunc("/upload", user.AuthHandlerWrapper(user.AuthHandler(), uploadHandler))
	// TODO: Jones: Frontend-Struktur überlegen und ApiCalls fürs Backend vorbereiten

	log.Println("web listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "index"
	layout(w, r, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(100 << 20) // max 100MB
		if err != nil {
			internalServerError(w, err)
			return
		}

		files := r.MultipartForm.File["images"]
		for _, fh := range files {
			file, err := fh.Open()
			if err != nil {
				internalServerError(w, err)
				return
			}

			// read file content into buffer
			buf := make([]byte, fh.Size)
			fReader := bufio.NewReader(file)
			_, err = fReader.Read(buf)
			if err != nil {
				internalServerError(w, err)
				return
			}

			// create request data with base64 encoded image buffer
			data := api.UploadReqData{
				Base64Image:  base64.StdEncoding.EncodeToString(buf),
				Filename:     fh.Filename,
				CreationDate: time.Now().Local(),
			}

			req, err := newPostRequest("upload", data)
			if err != nil {
				internalServerError(w, err)
				return
			}

			var res api.UploadResData
			err = callApi(r, req, &res)
			if err != nil {
				badRequest(w, err)
				return
			}
		}
	}
	layout(w, r, nil)
}

// TODO: jones tests schreiben
func registerHandler(w http.ResponseWriter, r *http.Request) {
	defer layout(w, r, nil)

	if r.Method == http.MethodPost {
		data := api.RegisterReqData{
			Username:             r.FormValue("user"),
			Password:             r.FormValue("password"),
			PasswordConfirmation: r.FormValue("passwordConfirmation"),
		}

		req, err := newPostRequest("register", data)
		if err != nil {
			internalServerError(w, err)
			return
		}

		var res api.RegisterResData
		err = callApi(r, req, &res)
		if err != nil {
			badRequest(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var res api.ThumbnailsResData
	payload := url.Values{
		"index":  {"0"},
		"length": {"25"},
	}
	req, err := newGetRequest("thumbnails?" + payload.Encode())
	if err != nil {
		internalServerError(w, err)
		return
	}

	err = callApi(r, req, &res)
	if err != nil {
		internalServerError(w, err)
		return
	}

	layout(w, r, res)
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

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// encode data to jsonUtil
//jsonBytes, err := json.Marshal(data)
//
//// prepare request
//req, err := http.NewRequest(method, DHBW_Photo_Server.BackendHost+url, bytes.NewBuffer(jsonBytes))
//if err != nil {
//return err
//}

func newPostRequest(url string, data interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return newRequest(http.MethodPost, url, jsonBytes)
}

func newGetRequest(url string) (*http.Request, error) {
	return newRequest(http.MethodGet, url, nil)
}

func newRequest(method string, url string, data []byte) (*http.Request, error) {
	return http.NewRequest(method, DHBW_Photo_Server.BackendHost+url, bytes.NewBuffer(data))
}

func callApi(r *http.Request, req *http.Request, res api.BaseRes) error {
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
	jsonBytes, err := ioutil.ReadAll(apiRes.Body)
	if err != nil {
		return err
	}

	// decode data from jsonUtil into result struct
	err = json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return err
	}

	// check for error from backend in res
	errorString := res.GetError()
	if errorString != "" {
		return errors.New(errorString)
	}

	return nil
}
