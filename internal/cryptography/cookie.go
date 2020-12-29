package cryptography

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
)

// TODO: Test schreiben
func GenerateCookie(name string, valuePrefix string, maxAge int) http.Cookie {
	valuePrefixBytes := []byte(valuePrefix)
	hasher := sha1.New()
	hasher.Write(valuePrefixBytes)
	if valuePrefix != "" {
		valuePrefix += DHBW_Photo_Server.CookieValueSeparator
	}
	cookieValue := valuePrefix + base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return http.Cookie{
		Name:   name,
		Value:  cookieValue,
		MaxAge: maxAge,
	}
}
