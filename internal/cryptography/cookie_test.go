package cryptography

import (
	DHBW_Photo_Server "DHBW.Photo-Server"
	"strings"
	"testing"
)

func TestGenerateCookiePrefix(t *testing.T) {
	name := "customCookieName"
	valuePrefix := "custom_"
	maxAge := 1234
	cookie := GenerateCookie(name, valuePrefix, maxAge)

	if cookie.Name != name || !strings.HasPrefix(cookie.Value, valuePrefix+DHBW_Photo_Server.CookieValueSeparator) || cookie.MaxAge != maxAge {
		t.Error("Something went wrong during cookie generation")
	}
}

func TestGenerateCookieNoPrefix(t *testing.T) {
	name := "customCookieName"
	valuePrefix := ""
	maxAge := 1234
	cookie := GenerateCookie(name, valuePrefix, maxAge)
	if cookie.Name != name || strings.HasPrefix(cookie.Value, DHBW_Photo_Server.CookieValueSeparator) || cookie.MaxAge != maxAge {
		t.Error("Something went wrong during cookie generation")
	}
}
