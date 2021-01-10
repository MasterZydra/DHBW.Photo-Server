/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package jsonUtil

import (
	"encoding/json"
	"net/http"
)

// decodes json body from http.Request into given data variable
func DecodeBody(r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// encodes given response data into http.ResponseWriter
func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}
	return nil
}
