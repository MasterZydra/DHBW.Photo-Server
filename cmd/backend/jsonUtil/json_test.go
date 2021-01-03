package jsonUtil

import (
	"DHBW.Photo-Server/internal/api"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeBodyValid(t *testing.T) {
	jsonString := `{"SomeString": "test", "SomeInteger": 1234}`
	body := bytes.NewReader([]byte(jsonString))
	req, _ := http.NewRequest(http.MethodGet, "/", body)

	var data api.TestReqData
	err := DecodeBody(req, &data)

	if err != nil || data.SomeString != "test" || data.SomeInteger != 1234 {
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

//func TestEncodeResponseValidNoError(t *testing.T) {
//	res := api.TestResData{
//		"",
//		"test",
//	}
//	w := // TODO: jones wo krieg ich einen http.ResponseWriter her?
//	err := EncodeResponse(w, &res)
//}

//func TestEncodeResponseValidError(t *testing.T) {
//
//}

//func TestEncodeResponseInvalid(t *testing.T) {
//
//}
