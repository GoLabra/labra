package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"
)

const (
	// Body size limits (keep these tight for login/signup).
	MaxBodyLoginBytes  int64 = 4 << 10 // 4 KB
	MaxBodySignupBytes int64 = 8 << 10 // 8 KB

	// Field limits.
	MaxEmailLen    = 254
	MinPasswordLen = 12
	MaxPasswordLen = 128
	MaxNameLen     = 64
)

// decodeJSONLimited reads JSON safely with:
// - http.MaxBytesReader (request body limit)
// - DisallowUnknownFields (reject unexpected fields)
// - EOF check (reject trailing junk)
func decodeJSONLimited(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		// If the body is too large, MaxBytesReader makes reads fail.
		// Commonly shows as "http: request body too large".
		return err
	}

	// Ensure there's no extra JSON (e.g. "{}{}")
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("invalid JSON (multiple values)")
	}

	return nil
}

func sanitizeEmail(s string) string {
	// Trim and lowercase for canonical comparison.
	return strings.ToLower(strings.TrimSpace(s))
}

func sanitizeName(s string) string {
	// Trim outer whitespace; keep inner spaces intact.
	return strings.TrimSpace(s)
}

// validateEmailRFC5322 validates format using net/mail parsing.
// It’s not “deliverability”, just syntax sanity (which is what we want here).
func validateEmailRFC5322(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > MaxEmailLen {
		return fmt.Errorf("email too long (max %d)", MaxEmailLen)
	}
	// ParseAddress allows "Name <a@b.com>" too; we only accept pure addr-spec.
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}
	if addr.Address != email {
		// Reject display-name forms and odd normalizations.
		return errors.New("invalid email format")
	}
	// Basic sanity: must contain @
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func validatePasswordLen(pw string) error {
	if pw == "" {
		return errors.New("password is required")
	}
	if len(pw) < MinPasswordLen {
		return fmt.Errorf("password too short (min %d)", MinPasswordLen)
	}
	if len(pw) > MaxPasswordLen {
		return fmt.Errorf("password too long (max %d)", MaxPasswordLen)
	}
	return nil
}

func validateNameField(label, v string) error {
	if v == "" {
		return fmt.Errorf("%s is required", label)
	}
	if len(v) > MaxNameLen {
		return fmt.Errorf("%s too long (max %d)", label, MaxNameLen)
	}
	return nil
}

func (d *LoginFormData) Sanitize() {
	d.Email = sanitizeEmail(d.Email)
	// Do NOT TrimSpace password; spaces might be intentional.
}

func (d *LoginFormData) Validate() error {
	if err := validateEmailRFC5322(d.Email); err != nil {
		return err
	}
	// For login, we still bound password length to avoid pathological input.
	if err := validatePasswordLen(d.Password); err != nil {
		return err
	}
	return nil
}

func (d *SignupFormData) Sanitize() {
	d.Email = sanitizeEmail(d.Email)
	d.FirstName = sanitizeName(d.FirstName)
	d.LastName = sanitizeName(d.LastName)
}

func (d *SignupFormData) Validate() error {
	if err := validateEmailRFC5322(d.Email); err != nil {
		return err
	}
	if err := validatePasswordLen(d.Password); err != nil {
		return err
	}
	if err := validateNameField("firstName", d.FirstName); err != nil {
		return err
	}
	if err := validateNameField("lastName", d.LastName); err != nil {
		return err
	}
	return nil
}
