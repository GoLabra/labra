package svc

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 14

// hashPasswordValue hashes a required password string (used for Create/Upsert inputs).
func hashPasswordValue(password *string) error {
	if password == nil {
		return fmt.Errorf("password is nil")
	}

	if *password == "" {
		return fmt.Errorf("password must not be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	*password = string(hashed)
	return nil
}

// hashPasswordOptional hashes an optional password (pointer field on Update inputs).
// If password is not present (nil), it does nothing.
func hashPasswordOptional(password **string) error {
	if password == nil || *password == nil {
		return nil
	}

	if **password == "" {
		return fmt.Errorf("password must not be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(**password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	**password = string(hashed)
	return nil
}
