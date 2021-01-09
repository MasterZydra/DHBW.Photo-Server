package main

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"DHBW.Photo-Server/internal/api"
	"DHBW.Photo-Server/internal/image"
	"DHBW.Photo-Server/internal/user"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func resetUsersFile() {
	csvFile, err := os.Create(DHBW_Photo_Server.UsersFile())
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(csvFile)
	var data = [][]string{
		{DHBW_Photo_Server.User1Name, DHBW_Photo_Server.Pw1Hash},
		{DHBW_Photo_Server.User2Name, DHBW_Photo_Server.Pw2Hash},
	}
	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}
}

func resetContentCsv(username string) {
	content := `img1.jpg|2021-01-06 18:27:20|d0d243adf8ab1746e8d904dbd5c4dbd1
img2.jpg|2021-01-06 18:27:19|54d47d125e436487e2d8859b9d691843`
	_ = ioutil.WriteFile(
		filepath.Join(DHBW_Photo_Server.ImageDir(), username, "content.csv"),
		[]byte(content),
		0755,
	)
}

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
}

func newPostRequest(url string, data interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
}

func executeRequest(r *http.Request) (*http.Response, error) {
	c := http.Client{}
	return c.Do(r)
}

func decodeJson(response *http.Response, result interface{}) error {
	// get jsonString from api response
	jsonBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// decode data from jsonUtil into result struct
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return err
	}
	return nil
}

func resetOrderList(username string) {
	um := user.UserManagerCache()
	usr := um.GetUser(username)

	usr.OrderList = &user.OrderList{Entries: []*user.OrderListEntry{}}
}

func TestMustParamWrongMethod(t *testing.T) {
	server := createServer(MustParam(func(w http.ResponseWriter, r *http.Request) {
		return
	}, http.MethodPost))

	response, _ := http.Get(server.URL)
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Wrong status code")
	}
}

func TestMustParamMissingGetParams(t *testing.T) {
	server := createServer(MustParam(func(w http.ResponseWriter, r *http.Request) {
		return
	}, http.MethodGet, "param1", "param2"))

	response, _ := http.Get(server.URL)
	if response.StatusCode != http.StatusBadRequest {
		t.Error("Wrong status code")
	}
}

func TestMustParamWithGetParams(t *testing.T) {
	server := createServer(MustParam(func(w http.ResponseWriter, r *http.Request) {
		return
	}, http.MethodGet, "param1", "param2"))

	response, _ := http.Get(server.URL + "?param1=some&param2=thing")
	if response.StatusCode != http.StatusOK {
		t.Error("Wrong status code")
	}
}

func TestRegisterHandler(t *testing.T) {
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetUsersFile()
	defer resetUsersFile()
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	data := api.RegisterReqData{
		Username:             "test",
		Password:             "sec123",
		PasswordConfirmation: "sec123",
	}
	req, _ := newPostRequest(server.URL, data)
	response, err := executeRequest(req)

	if err != nil || response.StatusCode != http.StatusOK {
		t.Error("Status code wrong")
	}
}

func TestRegisterHandlerInvalidJson(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	jsonBytes := []byte{14, 5, 86}
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(jsonBytes))

	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong")
	}
}

func TestRegisterHandlerPasswordsNotMatch(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	data := api.RegisterReqData{
		Username:             "test",
		Password:             "sec123",
		PasswordConfirmation: "someothervalue",
	}
	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	var res api.RegisterResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "The passwords do not match" {
		t.Error("Status code wrong")
	}
}

func TestRegisterHandlerInvalidUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r)
	})
	data := api.RegisterReqData{
		Username:             "inv?lid",
		Password:             "sec123",
		PasswordConfirmation: "sec123",
	}
	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong")
	}
}

func TestThumbnailsHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=0&length=25", nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.TotalImages != 2 || res.Images[0].Name != "img1.jpg" || res.Images[1].Name != "img2.jpg" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestThumbnailsHandlerNoUsername(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=0&length=25", nil)
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong error message")
	}
}

func TestThumbnailsHandlerInvalidIndex(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=wrong&length=25", nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Invalid index. Index must be an Integer" {
		t.Error("Status code wrong or wrong error message")
	}
}

func TestThumbnailsHandlerInvalidLength(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ThumbnailsHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL+"?index=0&length=wrong", nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.ThumbnailsResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Invalid length. Length must be an Integer" {
		t.Error("Status code wrong or wrong error message")
	}
}

func TestUploadHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r)
	})

	data := api.UploadReqData{
		Base64Image:  "/9j/4AAQSkZJRgABAQEAYABgAAD/4QBoRXhpZgAATU0AKgAAAAgABAEaAAUAAAABAAAAPgEbAAUAAAABAAAARgEoAAMAAAABAAIAAAExAAIAAAARAAAATgAAAAAAAABgAAAAAQAAAGAAAAABcGFpbnQubmV0IDQuMi4xNAAA/9sAQwABAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB/9sAQwEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB/8AAEQgABAAEAwEhAAIRAQMRAf/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC//EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC//EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29/j5+v/aAAwDAQACEQMRAD8A/Tj9hv8AZ08WeJ7T9q8y/tfftl6UfD/7bPx78LodI+MWmRnVI9Ifwsq6zqxv/BeoGbW78zl9RntDZWMrojQafbfPvK/1O8UvFnJ8Dx5n+Fj4G+A2IVKeX2rYngfMJVp8+U4Cp77o8R0afu8/LHlpxtCMU7u8n+T5LkNarleEm+I+Joc0anuwzGlyq1aotObCSfS7vJ6n/9k=",
		Filename:     "img3.jpg",
		CreationDate: time.Now().Local(),
	}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	defer os.Remove(filepath.Join(DHBW_Photo_Server.ImageDir(), "max", "img3.jpg"))
	defer os.Remove(filepath.Join(DHBW_Photo_Server.ImageDir(), "max", DHBW_Photo_Server.ThumbDir, "img3.jpg"))
	resetContentCsv("max")

	var res api.UploadResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.Error != "" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestUploadHandlerNoUsername(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r)
	})

	data := api.UploadReqData{
		Base64Image:  "/9j/4AAQSkZJRgABAQEAYABgAAD/4QBoRXhpZgAATU0AKgAAAAgABAEaAAUAAAABAAAAPgEbAAUAAAABAAAARgEoAAMAAAABAAIAAAExAAIAAAARAAAATgAAAAAAAABgAAAAAQAAAGAAAAABcGFpbnQubmV0IDQuMi4xNAAA/9sAQwABAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB/9sAQwEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB/8AAEQgABAAEAwEhAAIRAQMRAf/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC//EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC//EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29/j5+v/aAAwDAQACEQMRAD8A/Tj9hv8AZ08WeJ7T9q8y/tfftl6UfD/7bPx78LodI+MWmRnVI9Ifwsq6zqxv/BeoGbW78zl9RntDZWMrojQafbfPvK/1O8UvFnJ8Dx5n+Fj4G+A2IVKeX2rYngfMJVp8+U4Cp77o8R0afu8/LHlpxtCMU7u8n+T5LkNarleEm+I+Joc0anuwzGlyq1aotObCSfS7vJ6n/9k=",
		Filename:     "img3.jpg",
		CreationDate: time.Now().Local(),
	}

	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	var res api.UploadResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestUploadHandlerInvalidBase64(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r)
	})

	data := api.UploadReqData{
		Base64Image:  "invalidData!",
		Filename:     "img3.jpg",
		CreationDate: time.Now().Local(),
	}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.UploadResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong or wrong results")
	}
}

func TestUploadHandlerImageAlreadyExists(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r)
	})

	data := api.UploadReqData{
		Base64Image:  "/9j/4AAQSkZJRgABAQEAYABgAAD/4QBoRXhpZgAATU0AKgAAAAgABAEaAAUAAAABAAAAPgEbAAUAAAABAAAARgEoAAMAAAABAAIAAAExAAIAAAARAAAATgAAAAAAAABgAAAAAQAAAGAAAAABcGFpbnQubmV0IDQuMi4xNAAA/9sAQwABAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB/9sAQwEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB/8AAEQgABAAEAwEhAAIRAQMRAf/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC//EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC//EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29/j5+v/aAAwDAQACEQMRAD8A/Tj9hv8AZ08WeJ7T9q8y/tfftl6UfD/7bPx78LodI+MWmRnVI9Ifwsq6zqxv/BeoGbW78zl9RntDZWMrojQafbfPvK/1O8UvFnJ8Dx5n+Fj4G+A2IVKeX2rYngfMJVp8+U4Cp77o8R0afu8/LHlpxtCMU7u8n+T5LkNarleEm+I+Joc0anuwzGlyq1aotObCSfS7vJ6n/9k=",
		Filename:     "img2.jpg",
		CreationDate: time.Now().Local(),
	}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.UploadResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Image is already stored" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestOrderListEntryHandler(t *testing.T) {
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		OrderListEntryHandler(w, r)
	})

	um := user.UserManagerCache()
	usr := um.GetUser("max")

	date := time.Now().Local()
	img := image.NewImage("img1.jpg", date, "54d47d125e436487e2d8859b9d691843")
	_ = usr.OrderList.AddOrderEntry(img)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.OrderListResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.OrderList[0].Image.Name != "img1.jpg" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestOrderListEntryHandlerNoUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		OrderListEntryHandler(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	response, _ := executeRequest(req)

	var res api.OrderListResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestAddOrderListEntryHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		AddOrderListEntryHandler(w, r)
	})
	data := api.AddOrderListEntryReqData{ImageName: "img1.jpg"}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.AddOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.Error != "" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestAddOrderListEntryHandlerNoUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		AddOrderListEntryHandler(w, r)
	})
	data := api.AddOrderListEntryReqData{ImageName: "img1.jpg"}

	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	var res api.AddOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestAddOrderListEntryHandlerInvalidJson(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		AddOrderListEntryHandler(w, r)
	})

	jsonBytes := []byte{14, 5, 86}
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(jsonBytes))
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.AddOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong or wrong results")
	}
}

func TestAddOrderListEntryHandlerWrongImage(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		AddOrderListEntryHandler(w, r)
	})
	data := api.AddOrderListEntryReqData{ImageName: "nonexisting.jpg"}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.AddOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get image '"+data.ImageName+"'" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestRemoveOrderListEntryHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RemoveOrderListEntryHandler(w, r)
	})

	um := user.UserManagerCache()
	usr := um.GetUser("max")

	date := time.Now().Local()
	img := image.NewImage("img1.jpg", date, "54d47d125e436487e2d8859b9d691843")
	_ = usr.OrderList.AddOrderEntry(img)

	data := api.RemoveOrderListEntryReqData{ImageName: "img1.jpg"}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.RemoveOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.Error != "" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestRemoveOrderListEntryHandlerNoUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RemoveOrderListEntryHandler(w, r)
	})

	data := api.RemoveOrderListEntryReqData{ImageName: "img1.jpg"}

	req, _ := newPostRequest(server.URL, data)
	response, _ := executeRequest(req)

	var res api.RemoveOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestRemoveOrderListEntryHandlerInvalidJson(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RemoveOrderListEntryHandler(w, r)
	})

	jsonBytes := []byte{14, 5, 86}
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(jsonBytes))
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong or wrong results")
	}
}

func TestRemoveOrderListEntryHandlerImageNotExists(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		RemoveOrderListEntryHandler(w, r)
	})

	data := api.RemoveOrderListEntryReqData{ImageName: "img1.jpg"}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.RemoveOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Image 'img1.jpg' does not exist for user 'max'" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestChangeOrderListEntryHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ChangeOrderListEntryHandler(w, r)
	})

	um := user.UserManagerCache()
	usr := um.GetUser("max")

	date := time.Now().Local()
	img := image.NewImage("img1.jpg", date, "54d47d125e436487e2d8859b9d691843")
	_ = usr.OrderList.AddOrderEntry(img)

	data := api.ChangeOrderListEntryReqData{
		ImageName:      "img1.jpg",
		Format:         "letter",
		NumberOfPrints: 3,
	}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.ChangeOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || res.Error != "" || usr.OrderList.Entries[0].Format != "letter" || usr.OrderList.Entries[0].NumberOfPrints != 3 {
		t.Error("Status code wrong or wrong results")
	}
}

func TestChangeOrderListEntryHandlerNoUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ChangeOrderListEntryHandler(w, r)
	})

	req, _ := newPostRequest(server.URL, nil)
	response, _ := executeRequest(req)

	var res api.ChangeOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestChangeOrderListEntryHandlerInvalidJson(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ChangeOrderListEntryHandler(w, r)
	})

	jsonBytes := []byte{14, 5, 86}
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(jsonBytes))
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	if response.StatusCode != http.StatusBadRequest {
		t.Error("Status code wrong or wrong results")
	}
}

func TestChangeOrderListEntryHandlerEntryNotExists(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		ChangeOrderListEntryHandler(w, r)
	})

	data := api.ChangeOrderListEntryReqData{
		ImageName:      "img1.jpg",
		Format:         "letter",
		NumberOfPrints: 3,
	}

	req, _ := newPostRequest(server.URL, data)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.ChangeOrderListEntryResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not find order entry with image 'img1.jpg'" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestDeleteOrderListHandler(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		DeleteOrderListHandler(w, r)
	})

	um := user.UserManagerCache()
	usr := um.GetUser("max")

	date := time.Now().Local()
	img := image.NewImage("img1.jpg", date, "54d47d125e436487e2d8859b9d691843")
	_ = usr.OrderList.AddOrderEntry(img)

	req, _ := newPostRequest(server.URL, nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.DeleteOrderListResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusOK || len(usr.OrderList.Entries) > 0 {
		t.Error("Status code wrong or wrong results")
	}
}

func TestDeleteOrderListHandlerNoUsername(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		DeleteOrderListHandler(w, r)
	})

	req, _ := newPostRequest(server.URL, nil)
	response, _ := executeRequest(req)

	var res api.DeleteOrderListResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestDownloadOrderList(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		DownloadOrderList(w, r)
	})

	um := user.UserManagerCache()
	usr := um.GetUser("max")

	date := time.Now().Local()
	img := image.NewImage("img1.jpg", date, "54d47d125e436487e2d8859b9d691843")
	_ = usr.OrderList.AddOrderEntry(img)

	req, _ := newPostRequest(server.URL, nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	var res api.DownloadOrderListResData
	_ = decodeJson(response, &res)

	_, err := base64.StdEncoding.DecodeString(res.Base64ZipFile)

	if response.StatusCode != http.StatusOK || err != nil {
		t.Error("Status code wrong or wrong results")
	}
}

func TestDownloadOrderListNoUsername(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		DownloadOrderList(w, r)
	})

	req, _ := newPostRequest(server.URL, nil)
	response, _ := executeRequest(req)

	var res api.DownloadOrderListResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error != "Could not get username" {
		t.Error("Status code wrong or wrong results")
	}
}

func TestDownloadOrderListWrongImage(t *testing.T) {
	DHBW_Photo_Server.SetImageDir("../../test/example_imgs")
	DHBW_Photo_Server.SetUsersFile(DHBW_Photo_Server.TestUsersFile)
	resetOrderList("max")

	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		DownloadOrderList(w, r)
	})

	um := user.UserManagerCache()
	usr := um.GetUser("max")

	date := time.Now().Local()
	img := image.NewImage("invalid.jpg", date, "54d47d125e436487e2d8859b9d691843")
	_ = usr.OrderList.AddOrderEntry(img)

	req, _ := newPostRequest(server.URL, nil)
	req.SetBasicAuth("max", "pw")
	response, _ := executeRequest(req)

	defer os.Remove(filepath.Join(DHBW_Photo_Server.ImageDir(), "max", "order-list-download.zip"))

	var res api.DownloadOrderListResData
	_ = decodeJson(response, &res)

	if response.StatusCode != http.StatusBadRequest || res.Error == "" {
		t.Error("Status code wrong or wrong results")
	}
}
