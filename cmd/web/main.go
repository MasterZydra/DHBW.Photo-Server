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
	"strconv"
	"time"
)

// TODO: Jones Documentation

// TODO: Port über go flag konfigurierbar machen können
const port = "4443"

type TemplateVariables struct {
	Global GlobalVariables
	Result interface{}
	Local  interface{}
}

type GlobalVariables struct {
	Username string
	LoggedIn bool
}

var templateFuncMap = template.FuncMap{
	"sub": sub,
	"add": add,
}

func main() {
	// serve images directory
	fs := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", user.AuthFileServerWrapper(
		user.AuthFileServer(),
		http.StripPrefix("/images", fs),
	))

	// serve other static files (css, js etc.)
	staticServer := http.FileServer(http.Dir("cmd/web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", staticServer))

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
	layout(w, r, nil, nil)
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
	layout(w, r, nil, nil)
}

// TODO: jones tests schreiben
func registerHandler(w http.ResponseWriter, r *http.Request) {
	defer layout(w, r, nil, nil)

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
	query := r.URL.Query()
	index := query.Get("index")
	length := query.Get("length")

	indexInt, _ := strconv.Atoi(index)
	lengthInt, _ := strconv.Atoi(length)

	if index == "" || indexInt < 0 {
		index = "0"
		indexInt = 0
	}
	if length == "" || lengthInt < 1 {
		length = "25"
		lengthInt = 25
	}
	//query := r.URL.Query()
	//index, exists := query["index"]
	//if !exists {
	//}
	//length := r.URL.Query().Get("length")

	var res api.ThumbnailsResData
	payload := url.Values{
		"index":  {index},
		"length": {length},
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

	local := struct {
		Index  int
		Length int
	}{indexInt, lengthInt}

	layout(w, r, res, local)
}

func layout(w http.ResponseWriter, r *http.Request, result interface{}, local interface{}) {
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

	tmpl := template.New("photoserver").Funcs(templateFuncMap)
	tmpl, err = tmpl.ParseFiles(layout, publicFile)
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
		Result: result,
		Local:  local,
	}

	// execute template and send data with it
	err = tmpl.Funcs(templateFuncMap).ExecuteTemplate(w, "layout", templateVars)
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
