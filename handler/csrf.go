package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
)

const (
	csrfCookieName = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
)

// Require CSRF only when:
// - method is state-changing
// - and auth is cookie-based (jwt cookie exists)
// - and there is NO Authorization header (API clients)
func csrfRequired(r *http.Request) bool {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		// continue
	default:
		return false
	}

	// Don't require CSRF on login/signup endpoints
	switch r.URL.Path {
	case "/admin/login", "/admin/signup", "/login":
		return false
	}

	if strings.TrimSpace(r.Header.Get("Authorization")) != "" {
		return false
	}

	_, err := r.Cookie("jwt")
	return err == nil
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !csrfRequired(r) {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie(csrfCookieName)
		if err != nil || c.Value == "" {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		h := r.Header.Get(csrfHeaderName)
		if h == "" {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		if subtle.ConstantTimeCompare([]byte(c.Value), []byte(h)) != 1 {
			http.Error(w, "CSRF token invalid", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetCSRFCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		Secure:   secure,               // true in prod (https)
		HttpOnly: false,                // frontend must be able to read it
		SameSite: http.SameSiteLaxMode, // typical for session cookies
	})
}

func ClearCSRFCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   secure,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
}
