package cryptography

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// TODO: Jones Documentation

const (
	saltSize   = 16
	iterations = 1e4
)

// CreatePassword generates a new salt and passes this with the provided clear password to the hashPassword function.
// After that it encodes this new hashed password as a hex value and returns it.
func CreatePassword(password string) (string, error) {
	salt, err := createSalt(saltSize)
	if err != nil {
		return "", err
	}
	hashedPassword := hashPassword([]byte(password), salt)
	return hex.EncodeToString(hashedPassword), nil
}

// ComparePassword takes a hashed password in hex encoding and a clear text password.
// It decodes the hex hashed password into bytes and then hashes the clear password.
// After that it can compare the bytes from the hashed clear password and the decoded hashed password
// and return a boolean.
func ComparePassword(hexHash string, password string) (bool, error) {
	hashBytes, err := hex.DecodeString(hexHash)
	if err != nil {
		return false, err
	}
	compareHashBytes := hashPassword([]byte(password), hashBytes[:saltSize])
	return bytes.Equal(hashBytes, compareHashBytes), nil
}

// The internal function createSalt creates a random salt with the given n and returns these bytes.
func createSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// hashPassword takes the password bytes and the salt bytes and hashes the password with creating a new key
// and appending it to the salt.
func hashPassword(pw []byte, salt []byte) []byte {
	ret := make([]byte, len(salt))
	copy(ret, salt)
	return append(ret, createKey(pw, salt, iterations, sha256.Size, sha256.New)...)
}

// copied from "golang.org/x/crypto/pbkdf2"
func createKey(password []byte, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
