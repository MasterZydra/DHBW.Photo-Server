package cryptography

import (
	"bytes"
	"testing"
)

func TestHashImage(t *testing.T) {
	raw := []byte{'D', 'a', 's', 'i', 's', 't', 'e', 'i', 'n', 'T', 'e', 's', 't'}
	out := []byte{105, 4, 82, 118, 13, 169, 49, 43, 95, 62, 46, 250, 73, 161, 224, 172}
	hash := HashImage(&raw)
	if bytes.Compare(hash, out) != 0 {
		t.Errorf("HashImage(...) = %v, want %v", hash, out)
	}
}

func TestHashToString(t *testing.T) {
	raw := []byte{105, 4, 82, 118, 13, 169, 49, 43, 95, 62, 46, 250, 73, 161, 224, 172}
	out := "690452760da9312b5f3e2efa49a1e0ac"
	encoded :=  HashToString(raw)
	if encoded != out {
		t.Errorf("HashToString(...) = %v, want %v", encoded, out)
	}
}
