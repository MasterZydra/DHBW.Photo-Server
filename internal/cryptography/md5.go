/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package cryptography

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash image with md5 algorithm. The image needs to be a pointer to a byte array.
// It returns a new byte array which contains the hashed result.
func HashImage(image []byte) []byte {
	hash := md5.Sum(image)
	return hash[:]
}

// Convert a given byte array to a string.
// This can be used to convert the hashed image to a string and store it.
func HashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}
