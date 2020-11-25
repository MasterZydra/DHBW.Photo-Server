package cryptography

import (
	"crypto/md5"
	"encoding/hex"
)

func HashImage(image *[]byte) []byte {
	hash := md5.Sum(*image)
	return hash[:]
}

func HashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}