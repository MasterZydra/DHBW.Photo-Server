package cryptography

import "crypto/md5"

func HashImage(image string) [md5.Size]byte {
	data := []byte(image)
	return md5.Sum(data)
}