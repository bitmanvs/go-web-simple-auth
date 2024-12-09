package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash interface {
	MatchesPlaintext(string) bool
	Bytes() []byte
}

// Hash a plaintext password into a secure password hash.
func HashPassword(plaintext string) (PasswordHash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return bcryptPasswordHash(bytes), nil
}

// Converts raw bytes into a password hash. Note that this doesn't perform a
// hash on the bytes. The bytes represent an already-hashed password.
func PasswordHashFromBytes(bytes []byte) PasswordHash {
	return bcryptPasswordHash(bytes)
}
