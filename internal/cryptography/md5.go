package cryptography

import "crypto/md5"

func HashImage(image []byte) [md5.Size]byte {
	return md5.Sum(image)
}