/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package DHBW_Photo_Server

import (
	"os"
	"strings"
)

const (
	// Image folder configuration
	DefaultImageDir = "images" // absolute or relative to working directory
	ThumbDir        = "previews"
	UserContent     = "content.csv"

	// Default values for server configurations
	BackendDefaultPort = 3000
	BackendHost        = "https://localhost:3000/"
	WebDefaultUrl      = "https://localhost"
	WebDefaultPort     = 4443

	// userfiles
	TestUsersFile = "../../test/usersFile_test.csv"
	ProdUsersFile = "usersFile.csv"

	// user1
	User1Name = "max"
	Pw1Clear  = "test123"
	Pw1Hash   = "73a64b63aeb9e71d4e10df824ab4a9d32ce1911bf343d085fb67dec7aba0fb711bda08780efc5d9291df3e8e1c7a66b2"

	// user2
	User2Name = "ana"
	Pw2Clear  = "123test"
	Pw2Hash   = "e0f3ae3d616e121df29464a191b1e5cec18c84190490550230b3a8f93930b71e46a90e876b6896839996675259096fd4"

	// other
	UsernameRegexBlacklist   = `[^a-z^A-Z^0-9\^_\^.\^-]`
	TimeLayout               = "2006-01-02 15:04:05"
	ImagesCacheControlMaxAge = "2592000"
)

var OrderListFormats = []string{"Junior Legal (8 x 5)", "Letter (8.5 x 11)", "Legal (8.5 x 14)", "Tabloid (11 x 17)"}

var imageDir = DefaultImageDir

func SetImageDir(image string) {
	imageDir = image
}

func ImageDir() string {
	return strings.Trim(strings.Trim(imageDir, "/"), string(os.PathSeparator))
}

var usersFile = ProdUsersFile

func SetUsersFile(file string) {
	usersFile = file
}

func UsersFile() string {
	return usersFile
}
