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

// TODO: Port über go flag konfigurierbar machen können
const port = "4443"

// structs to send different variables to template
type TemplateVariables struct {
	Global GlobalVariables
	Result interface{}
	Local  interface{}
}
type GlobalVariables struct {
	Username string
	LoggedIn bool
	ThumbDir string
}

// templateFuncMap initializes template functions to use in html templates.
var templateFuncMap = template.FuncMap{
	"sub": Sub,
	"add": Add,
}

var webRoot = filepath.Join("cmd", "web")
var layoutName = "layout"

var backendServerRoot = DHBW_Photo_Server.BackendHost

func main() {
	// serve images directory
	fs := http.FileServer(http.Dir(DHBW_Photo_Server.ImageDir()))
	http.Handle("/images/", user.AuthFileServerWrapper(
		user.AuthFileServer(),
		http.StripPrefix("/images", fs),
	))

	// serve other static files (css, js etc.)
	//staticServer := http.FileServer(http.Dir("cmd/web/assets"))
	//http.Handle("/assets/", http.StripPrefix("/assets", staticServer))

	// Handlers with auth wrappers if needed
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/home", user.AuthHandlerWrapper(user.AuthHandler(), HomeHandler))
	http.HandleFunc("/upload", user.AuthHandlerWrapper(user.AuthHandler(), UploadHandler))

	// listen and start server
	log.Println("web listening on https://localhost:" + port)
	log.Fatalln(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil))
}

// If navigating to the root the RootHandler sets index as the current path to load index.html in Layout.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "index"
	Layout(w, r, nil, nil)
}

// The UploadHandler parses the files from the HTML upload form and triggers a new API call for each file.
// If upload fails all following uploads are aborted.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

			req, err := NewPostRequest("upload", data)
			if err != nil {
				internalServerError(w, err)
				return
			}

			var res api.UploadResData
			err = CallApi(r, req, &res)
			if err != nil {
				badRequest(w, err)
				return
			}
		}
	}
	Layout(w, r, nil, nil)
}

// The RegisterHandler takes user, password and passwordConfirmation from the HTML form
// and triggers a API call to /register where it sends these values.
// If successful the site simply reloads.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data := api.RegisterReqData{
			Username:             r.FormValue("user"),
			Password:             r.FormValue("password"),
			PasswordConfirmation: r.FormValue("passwordConfirmation"),
		}

		req, err := NewPostRequest("register", data)
		if err != nil {
			internalServerError(w, err)
			return
		}

		var res api.RegisterResData
		err = CallApi(r, req, &res)
		if err != nil {
			badRequest(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	Layout(w, r, nil, nil)
}

// If a user is logged and visits /home in the HomeHandler is triggered.
// It sets a default index and length GET values for the API call (/thumbnails)
// to get the thumbnails from index to length.
// Afterwards these variables are made available to the Layout via TemplateVariables.Local
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	index := query.Get("index")
	length := query.Get("length")

	indexInt, _ := strconv.Atoi(index)
	lengthInt, _ := strconv.Atoi(length)

	// set default values
	if index == "" || indexInt < 0 {
		index = "0"
		indexInt = 0
	}
	if length == "" || lengthInt < 1 {
		length = "25"
		lengthInt = 25
	}

	// create GET request for API call
	payload := url.Values{
		"index":  {index},
		"length": {length},
	}
	req, err := NewGetRequest("thumbnails?" + payload.Encode())
	if err != nil {
		internalServerError(w, err)
		return
	}

	var res api.ThumbnailsResData
	err = CallApi(r, req, &res)
	if err != nil {
		badRequest(w, err)
		return
	}

	local := struct {
		Index  int
		Length int
	}{indexInt, lengthInt}

	Layout(w, r, res, local)
}

// Layout creates a new template, assigns its template functions from templateFuncMap,
// loads the Layout file + the current requested HTML file and adds the TemplateVariables
// so they are accessible in the HTML templates.
// After that the template is executed
func Layout(w http.ResponseWriter, r *http.Request, result interface{}, local interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		internalServerError(w, err)
		return
	}

	dir := filepath.Join(wd, webRoot)
	layout := filepath.Join(dir, "templates", layoutName+".html")
	publicFile := filepath.Join(dir, "public", filepath.Clean(r.URL.Path)+".html")

	// check if site exists or is a directory
	siteStat, err := os.Stat(publicFile)
	if err != nil && os.IsNotExist(err) || siteStat.IsDir() {
		http.NotFound(w, r)
		return
	}

	// create template and add templateFuncMap
	tmpl := template.New("photoserver").Funcs(templateFuncMap)
	tmpl, err = tmpl.ParseFiles(layout, publicFile)
	if err != nil {
		internalServerError(w, err)
		return
	}

	// set TemplateVariables
	username, _, ok := r.BasicAuth()
	loggedIn := true
	if !ok {
		username = ""
		loggedIn = false
	}
	templateVars := TemplateVariables{
		Global: GlobalVariables{
			Username: username,
			LoggedIn: loggedIn,
			ThumbDir: DHBW_Photo_Server.ThumbDir,
		},
		Result: result,
		Local:  local,
	}

	// execute template and send data with it
	err = tmpl.Funcs(templateFuncMap).ExecuteTemplate(w, layoutName, templateVars)
	if err != nil {
		internalServerError(w, err)
		return
	}
}

// simple helper for internal server error output
func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
}

// simple helper for a bad request response
func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// converts data to json bytes and passes it to NewRequest to return a new POST request
func NewPostRequest(url string, data interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return NewRequest(http.MethodPost, url, jsonBytes)
}

// simple wrapper for a new GET request
func NewGetRequest(url string) (*http.Request, error) {
	return NewRequest(http.MethodGet, url, nil)
}

// returns a new request with the configured BackendHost as a prefix to the url
func NewRequest(method string, url string, data []byte) (*http.Request, error) {
	return http.NewRequest(method, backendServerRoot+url, bytes.NewBuffer(data))
}

// Wrapper for making an API call.
// CallApi takes two http.Request objects and a result object api.BaseRes.
// The first request object is used to pass the basic auth credentials to the other request.
// After that the second request is executed.
// After getting a response, it is parsed into the given result object
func CallApi(r *http.Request, req *http.Request, res api.BaseRes) error {
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
