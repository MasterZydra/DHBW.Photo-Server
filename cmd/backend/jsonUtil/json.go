package jsonUtil

import (
	"encoding/json"
	"net/http"
)

// decodes json body from http.Request into given data variable
func DecodeBody(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	return nil
}

// encodes given response data into http.ResponseWriter
func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
