/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package jsonUtil

import (
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(f)
}

func TestDecodeBodyValid(t *testing.T) {
	jsonString := `{"SomeString": "test", "SomeInt": 1234}`
	body := bytes.NewReader([]byte(jsonString))
	req, _ := http.NewRequest(http.MethodGet, "/", body)

	var data api.TestReqData
	err := DecodeBody(req, &data)

	if err != nil || data.SomeString != "test" || data.SomeInt != 1234 {
		t.Error("json wasn't successfully decoded, when it should have been")
	}
}

func TestDecodeBodyInvalid(t *testing.T) {
	jsonString := `{someInvalid: "json"}`
	body := bytes.NewReader([]byte(jsonString))
	req := httptest.NewRequest(http.MethodGet, "/", body)

	var data api.TestReqData
	err := DecodeBody(req, &data)

	if err == nil {
		t.Error("json was successfully decoded, when it shouldn't have been")
	}
}

func TestEncodeResponseValid(t *testing.T) {
	server := createServer(func(w http.ResponseWriter, r *http.Request) {
		res := api.TestResData{
			Error:            "",
			SomeResultString: "fromserver",
			SomeResultInt:    12,
		}
		_ = EncodeResponse(w, &res)
	})

	response, err := http.Get(server.URL)

	var res api.TestResData
	jsonBytes, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(jsonBytes, &res)

	if err != nil || res.Error != "" || res.SomeResultString != "fromserver" || res.SomeResultInt != 12 {
		t.Error("Could not get ")
	}
}
