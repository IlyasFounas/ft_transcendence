package code_generation

import (
	"crypto/rand"
	"encoding/base32"
)

func RandGenerate(nb_octets int) (string, error) {
	key := make([]byte, nb_octets)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	encodedKey := base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(key)
	return encodedKey, nil
}
