package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

// This function uses bcrypt to hash the data passed as a parameter.
// Most of the time, this is used for passwords to ensure they are
// never stored in plain text in the database.
func HashPassword(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), 14)
	return string(bytes), err
}

// It checks the provided data against the stored hash.
// Example usage:
// CheckPasswordHash("hello", "asd445..")
func CheckPasswordHash(data, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}
