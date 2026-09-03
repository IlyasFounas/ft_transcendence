package join_code

import (
	"crypto/rand"
	"encoding/base32"
)

func RandGenerate() (string, error) {
	key := make([]byte, 4)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	encodedKey := base32.StdEncoding.
	WithPadding(base32.NoPadding).
	EncodeToString(key)
	return encodedKey, nil
}