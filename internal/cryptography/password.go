package cryptography

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

const (
	saltSize   = 16
	iterations = 1e4
)

func CreatePassword(password string) (string, error) {
	salt, err := createSalt(saltSize)
	if err != nil {
		return "", err
	}
	hashedPassword := hashPassword([]byte(password), salt)
	return hex.EncodeToString(hashedPassword), nil
}

func ComparePassword(hexHash string, password string) (bool, error) {
	hashBytes, err := hex.DecodeString(hexHash)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hashBytes, hashPassword([]byte(password), hashBytes[:saltSize])), nil
}

func createSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

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
