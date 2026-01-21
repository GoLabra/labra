package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GenerateActivationKey generates a unique activation key and returns both
// the plain key (for sending to the user) and the hashed key (for storing in the database).
// The hashed key is generated using bcrypt with cost 14, same as passwords.
func GenerateActivationKey() (plainKey string, hashedKey string, err error) {
	// Generate a secure random key (32 bytes = 256 bits)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("error generating random activation key: %w", err)
	}

	// Encode to base64 URL-safe string for the plain key
	plainKey = base64.URLEncoding.EncodeToString(b)

	// Hash the key using bcrypt (same algorithm and cost as passwords)
	hashedKeyBytes, err := bcrypt.GenerateFromPassword([]byte(plainKey), 14)
	if err != nil {
		return "", "", fmt.Errorf("error hashing activation key: %w", err)
	}

	hashedKey = string(hashedKeyBytes)

	return plainKey, hashedKey, nil
}

// HashActivationKey hashes an existing activation key using bcrypt.
// This is useful if you already have a plain key and need to hash it.
func HashActivationKey(plainKey string) (string, error) {
	hashedKeyBytes, err := bcrypt.GenerateFromPassword([]byte(plainKey), 14)
	if err != nil {
		return "", fmt.Errorf("error hashing activation key: %w", err)
	}
	return string(hashedKeyBytes), nil
}

// VerifyActivationKey verifies if a plain activation key matches the hashed key.
func VerifyActivationKey(plainKey, hashedKey string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(plainKey))
	return err == nil
}



