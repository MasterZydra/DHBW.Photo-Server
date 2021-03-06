/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/cmd/backend/jsonUtil"
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
}

func createFileServer(f http.Handler) *httptest.Server {
	return httptest.NewServer(f)
}

func newPostReq(url string, data interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return newReq(http.MethodPost, url, jsonBytes)
}

func newGetReq(url string) (*http.Request, error) {
	return newReq(http.MethodGet, url, nil)
}

func newReq(method string, url string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("username", "password")
	return req, nil
}

func TestCallApiPost(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		var data api.TestReqData
		_ = jsonUtil.DecodeBody(r, &data)

		res := api.TestResData{
			Error:            "",
			SomeResultString: data.SomeString + "server",
			SomeResultInt:    data.SomeInt + 1,
		}
		_ = jsonUtil.EncodeResponse(w, &res)
	})
	defer server.Close()

	reqData := api.TestReqData{
		SomeString: "request",
		SomeInt:    12,
	}
	req, _ := newPostReq(server.URL, reqData)

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err != nil || res.SomeResultString != reqData.SomeString+"server" || res.SomeResultInt != reqData.SomeInt+1 {
		t.Error("Error occurred or result data of POST request wrong.")
	}
}

func TestCallApiGet(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		someString := r.URL.Query().Get("someString")
		res := api.TestResData{
			Error:            "",
			SomeResultString: someString + "server",
			SomeResultInt:    123,
		}
		_ = jsonUtil.EncodeResponse(w, &res)
	})
	defer server.Close()

	payload := url.Values{
		"someString": {"test"},
	}
	req, _ := newGetReq(server.URL + "?" + payload.Encode())

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err != nil || res.SomeResultString != "testserver" || res.SomeResultInt != 123 {
		t.Error("Error occurred or result data of GET request wrong")
	}
}

func TestCallApiFailDo(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
	})
	defer server.Close()

	req, _ := newGetReq("invalid://url.com")

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err == nil {
		t.Error("No error occurred")
	}
}

func TestCallApiFailUnmarshalJson(t *testing.T) {
	// nicht testbar, da es sehr schwierig ist ungültige Daten in json.Unmarshal über den Server reinzuschicken
}

func TestCallApiCustomError(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		res := api.TestResData{
			Error: "customServerError",
		}
		_ = jsonUtil.EncodeResponse(w, &res)
	})
	defer server.Close()

	req, _ := newGetReq(server.URL)

	var res api.TestResData
	err := CallApi(req, req, &res)
	if err == nil || err.Error() != "customServerError" {
		t.Error("Error occurred or result data of GET request wrong")
	}
}

func TestNewRequest(t *testing.T) {
	backendServerRoot = DHBW_Photo_Server.BackendHost
	req, err := NewRequest("POST", "some/path", nil)
	if err != nil || req.Method != "POST" || req.URL.Path != "/some/path" || !strings.Contains(DHBW_Photo_Server.BackendHost, req.URL.Host) {
		t.Error("Error while creating new request")
	}
}

func TestNewGetRequest(t *testing.T) {
	req, err := NewGetRequest("some/path")
	if err != nil || req.Method != "GET" {
		t.Error("Error while creating GET request")
	}
}

func TestNewPostRequestValidJson(t *testing.T) {
	data := api.TestReqData{
		SomeString: "test",
		SomeInt:    1,
	}
	req, err := NewPostRequest("some/path", data)
	if err != nil || req.Method != "POST" {
		t.Error("Error while creating POST request")
	}
}

func TestNewPostRequestInvalidJson(t *testing.T) {
	// nicht testbar, da es sehr schwierig ist über ein interface{} ungültige Daten in json.Marshal reinzuschicken
}

func TestLayoutFileNotFound(t *testing.T) {
	webRoot = "someWrongWebRoot"
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		Layout(w, r, nil, nil)
	})
	defer server.Close()

	response, err := http.Get(server.URL)

	if err != nil || response.StatusCode != http.StatusNotFound {
		t.Error("Response has not the status 404")
	}
}

func TestLayoutWrongTemplate(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "wrong-template"

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		Layout(w, r, nil, nil)
	})
	defer server.Close()

	server.URL = server.URL + "/test"
	response, _ := http.Get(server.URL)

	if response.StatusCode != http.StatusInternalServerError {
		t.Error("Should throw a error 500 since template has wrong syntax")
	}
}

func TestLayoutValidTemplate(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		backendResult := api.TestResData{
			Error:            "",
			SomeResultString: "serverstring",
			SomeResultInt:    54,
		}
		local := struct{ LocalString string }{"whatever"}
		r.SetBasicAuth("name", "pw")

		Layout(w, r, backendResult, local)
	})
	defer server.Close()

	server.URL = server.URL + "/test"
	response, _ := http.Get(server.URL)

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	html := string(bodyBytes)

	expectedHtml := "nametrueserverstring540108whatever"

	if response.StatusCode != http.StatusOK || html != expectedHtml {
		t.Error("valid Layout did not load correctly")
	}
}

func TestCacheWrapper(t *testing.T) {
	var cacheControlHeader string
	server := createFileServer(CacheWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cacheControlHeader = w.Header().Get("Cache-Control")
	})))

	_, _ = http.Get(server.URL)

	if cacheControlHeader == "" {
		t.Error("Cache-Control header is empty!")
	}
}

func TestRootHandler(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	var indexPath string
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RootHandler(w, r)
		indexPath = r.URL.Path
	})
	defer server.Close()

	_, _ = http.Get(server.URL)

	if indexPath != "index" {
		t.Error("RootHandler couldn't redirect from / to /index.html")
	}
}

func postBinaryData(server *httptest.Server, postData string) (resp *http.Response, err error) {
	return http.Post(
		server.URL,
		`multipart/form-data; boundary=xxx`,
		ioutil.NopCloser(strings.NewReader(postData)),
	)
}

func postFormData(server *httptest.Server, postData string) (resp *http.Response, err error) {
	return http.Post(
		server.URL,
		`application/x-www-form-urlencoded`,
		ioutil.NopCloser(strings.NewReader(postData)),
	)
}

func TestUploadHandlerMultipartFormError(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r)
	})
	defer server.Close()

	response, _ := postBinaryData(server, "invaliddata")

	if response.StatusCode != http.StatusInternalServerError {
		t.Error("Parsing the multipart form should fail here")
	}
}

func TestUploadHandlerMultipleImages(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test"
		UploadHandler(w, r)
	})
	defer webServer.Close()

	var expectedImages []string
	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.UploadResData
		defer jsonUtil.EncodeResponse(w, &res)

		var data api.UploadReqData
		_ = jsonUtil.DecodeBody(r, &data)
		imageBytes, _ := base64.StdEncoding.DecodeString(data.Base64Image)
		expectedImages = append(expectedImages, string(imageBytes))
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "upload"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	// from: https://golang.org/src/net/http/request_test.go?h=Request
	postData := `--xxx
Content-Disposition: form-data; name="images"; filename="file1.jpg"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data1
--xxx
Content-Disposition: form-data; name="images"; filename="file2.jpg"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data2
--xxx--
`
	response, _ := postBinaryData(webServer, postData)

	if response.StatusCode != http.StatusOK || expectedImages[0] != "binary data1" || expectedImages[1] != "binary data2" {
		t.Error("Uploading the two images should work here")
	}
}

func TestUploadHandlerBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test"
		UploadHandler(w, r)
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return // do nothing: will result in a "unexpected end of JSON" error
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "upload"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	// from: https://golang.org/src/net/http/request_test.go?h=Request
	postData := `--xxx
Content-Disposition: form-data; name="images"; filename="file1.jpg"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data1
--xxx--
`
	response, _ := postBinaryData(webServer, postData)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Uploading the two images should work here")
	}
}

func TestRegisterHandlerLayout(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test"
		RegisterHandler(w, r)
	})
	defer webServer.Close()

	response, _ := http.Get(webServer.URL)

	if response.StatusCode != http.StatusOK {
		t.Error("Something went wrong while testing registrationHandler")
	}
}

func TestRegisterHandlerSuccessful(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test"
		RegisterHandler(w, r)
	})
	defer webServer.Close()

	var expectedUsername string
	var expectedPw string
	var expectedPwConfirmation string
	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.RegisterResData
		defer jsonUtil.EncodeResponse(w, &res)

		var data api.RegisterReqData
		_ = jsonUtil.DecodeBody(r, &data)
		expectedUsername = data.Username
		expectedPw = data.Password
		expectedPwConfirmation = data.PasswordConfirmation
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "register"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	postData := `--xxx
Content-Disposition: form-data; name="user"

benutzer1
--xxx
Content-Disposition: form-data; name="password"

securepw123
--xxx
Content-Disposition: form-data; name="passwordConfirmation"

securepw123
--xxx--
`

	response, _ := postBinaryData(webServer, postData)

	if response.StatusCode != http.StatusOK || expectedUsername != "benutzer1" || expectedPw != "securepw123" || expectedPwConfirmation != "securepw123" {
		t.Error("Something went wrong while testing registrationHandler")
	}
}

func TestRegisterHandlerBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test"
		RegisterHandler(w, r)
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return // do nothing: will result in a "unexpected end of JSON" error
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "register"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	postData := `--xxx
Content-Disposition: form-data; name="user"

benutzer1
--xxx--
`
	response, _ := postBinaryData(webServer, postData)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Something went wrong while testing registrationHandler")
	}
}

func TestHomeHandler(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		HomeHandler(w, r)
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.ThumbnailsResData
		defer jsonUtil.EncodeResponse(w, &res)

		var data api.ThumbnailsReqData
		_ = jsonUtil.DecodeBody(r, &data)
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "thumbnails"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	response, _ := http.Get(webServer.URL)

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	html := string(bodyBytes)

	if response.StatusCode != http.StatusOK || html != "nameindex:0length:25" {
		t.Error("Response body of /home is incorrect or wrong status code")
	}
}

func TestHomeHandlerBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		HomeHandler(w, r)
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return // do nothing: will result in a "unexpected end of JSON" error
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "thumbnails"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	response, _ := http.Get(webServer.URL)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Something went wrong while testing HomeHandler")
	}
}

func TestHomeHandlerAddOrderListEntries(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		HomeHandler(w, r)
	})
	defer webServer.Close()

	var imageNames []string
	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.AddOrderListEntryResData
		defer jsonUtil.EncodeResponse(w, &res)

		var data api.AddOrderListEntryReqData
		_ = jsonUtil.DecodeBody(r, &data)
		imageNames = append(imageNames, data.ImageName)
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "addOrderListEntry"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	postData := `toOrderListimg1.jpg=img1.jpg&toOrderListimg2.jpg=img2.jpg`
	response, _ := postFormData(webServer, postData)

	if response.StatusCode != http.StatusOK || imageNames[0] != "img1.jpg" || imageNames[1] != "img2.jpg" {
		t.Error("wrong status code or wrong imageNames received at the server")
	}
}

func TestHomeHandlerAddOrderListEntriesInvalidForm(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		HomeHandler(w, r)
	})
	defer webServer.Close()

	postData := `invaliddata:;_!"§$%&/()"`
	response, _ := postFormData(webServer, postData)

	if response.StatusCode != http.StatusInternalServerError {
		t.Error("wrong status code or wrong imageNames received at the server")
	}
}

func TestHomeHandlerAddOrderListEntriesBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		HomeHandler(w, r)
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "addOrderListEntry"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	postData := `invalid=data`
	response, _ := postFormData(webServer, postData)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("wrong status code")
	}
}

func TestUpdateOrderListEntry(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		if r.Method == http.MethodPost {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	var imageName string
	var format string
	var numberOfPrints int
	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.ChangeOrderListEntryResData
		defer jsonUtil.EncodeResponse(w, &res)

		if r.Method == http.MethodPost {
			var data api.ChangeOrderListEntryReqData
			_ = jsonUtil.DecodeBody(r, &data)

			imageName = data.ImageName
			format = data.Format
			numberOfPrints = data.NumberOfPrints
		}
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "changeOrderListEntry"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	postData := "imageName=img1.jpg&format=Letter&numberOfPrints='3'"
	response, _ := postFormData(webServer, postData)

	if response.StatusCode != http.StatusOK || imageName != "img1.jpg" || format != "Letter" || numberOfPrints != 1 {
		t.Error("Wrong status code or wrong parameters received on the server")
	}
}

func TestUpdateOrderListEntryBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		if r.Method == http.MethodPost {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "changeOrderListEntry"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	postData := "invaliddata=blabla"
	response, _ := postFormData(webServer, postData)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestRemoveOrderListEntry(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		imageToRemove := r.URL.Query().Get("ImageToRemove")
		if imageToRemove != "" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	var imageName string
	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.RemoveOrderListEntryResData
		defer jsonUtil.EncodeResponse(w, &res)
		if r.Method == http.MethodPost {
			var data api.RemoveOrderListEntryReqData
			_ = jsonUtil.DecodeBody(r, &data)

			imageName = data.ImageName
		}
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "removeOrderListEntry"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"ImageToRemove": {"img1.jpg"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusOK || imageName != "img1.jpg" {
		t.Error("Wrong status code or wrong parameters received on the server")
	}
}

func TestRemoveOrderListEntryBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		imageToRemove := r.URL.Query().Get("ImageToRemove")
		if imageToRemove != "" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "removeOrderListEntry"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"ImageToRemove": {"img1.jpg"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestDeleteOrderList(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		deleteOrderList := r.URL.Query().Get("deleteOrderList")
		if deleteOrderList == "1" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.DeleteOrderListResData
		defer jsonUtil.EncodeResponse(w, &res)
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "deleteOrderList"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"deleteOrderList": {"1"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusOK {
		t.Error("Wrong status code")
	}
}

func TestDeleteOrderListBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		deleteOrderList := r.URL.Query().Get("deleteOrderList")
		if deleteOrderList == "1" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "deleteOrderList"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"deleteOrderList": {"1"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestDownloadOrderList(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		downloadOrderList := r.URL.Query().Get("downloadOrderList")
		if downloadOrderList == "1" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	var base64String string
	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.DownloadOrderListResData
		defer jsonUtil.EncodeResponse(w, &res)

		base64String = "SGVsbXV0IE5lZW1hbm4="
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "downloadOrderList"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"downloadOrderList": {"1"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusOK || base64String != "SGVsbXV0IE5lZW1hbm4=" {
		t.Error("Wrong status code")
	}
}

func TestDownloadOrderListBadApiRequest(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		downloadOrderList := r.URL.Query().Get("downloadOrderList")
		if downloadOrderList == "1" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		return
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "downloadOrderList"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"downloadOrderList": {"1"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestDownloadOrderListNoBasicAuth(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		downloadOrderList := r.URL.Query().Get("downloadOrderList")
		if downloadOrderList == "1" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	payload := url.Values{
		"downloadOrderList": {"1"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusInternalServerError {
		t.Error("Wrong status code")
	}
}

func TestDownloadOrderListInvalidBase64(t *testing.T) {
	webRoot = "../../test/html"
	layoutName = "valid-template"

	webServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/test-home"
		r.SetBasicAuth("name", "pw")
		downloadOrderList := r.URL.Query().Get("downloadOrderList")
		if downloadOrderList == "1" {
			OrderListHandler(w, r)
		}
	})
	defer webServer.Close()

	backendServer := createServer(func(w http.ResponseWriter, r *http.Request) {
		var res api.DownloadOrderListResData
		defer jsonUtil.EncodeResponse(w, &res)

		res.Base64ZipFile = "INVALIDBASE64StRing"
	})
	BackendUrlRoot := backendServer.URL + "/"
	backendServer.URL = BackendUrlRoot + "downloadOrderList"
	backendServerRoot = BackendUrlRoot

	defer backendServer.Close()

	payload := url.Values{
		"downloadOrderList": {"1"},
	}
	response, _ := http.Get(webServer.URL + "/?" + payload.Encode())

	if response.StatusCode != http.StatusInternalServerError {
		t.Error("Wrong status code")
	}
}
