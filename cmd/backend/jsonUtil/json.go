package jsonUtil

import (
	"encoding/json"
	"net/http"
)

// TODO: Jones Documentation

func DecodeBody(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	return nil
}

func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
