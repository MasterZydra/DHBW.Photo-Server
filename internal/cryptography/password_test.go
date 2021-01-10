/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package cryptography

import (
	"testing"
)

func TestCreatePassword(t *testing.T) {
	hexHash, err := CreatePassword("test123")
	if err != nil || hexHash == "" {
		t.Errorf("Error while creating a salted hash: %v", err)
	}
}

func TestComparePasswordTrue(t *testing.T) {
	hexHash, _ := CreatePassword("test123")
	ok, err := ComparePassword(hexHash, "test123")
	if !ok || err != nil {
		t.Errorf("The passwords do not match; expected the hashes to match")
	}
}

func TestComparePasswordFalse(t *testing.T) {
	hexHash, _ := CreatePassword("test123")
	ok, _ := ComparePassword(hexHash, "123test")
	if ok {
		t.Errorf("The passwords match; expected the hashes to not match")
	}
}

func TestComparePasswordWrongHex(t *testing.T) {
	hexHash, _ := CreatePassword("test123")
	hexHash += "abcdefghijklmnopqrstuvwxyz"
	_, err := ComparePassword(hexHash, "123test")
	if err == nil {
		t.Errorf("The hex hash is valid; expected it not to be valid")
	}
}
