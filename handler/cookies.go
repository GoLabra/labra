package handler

import (
	"net/http"
	"strings"
	"time"
)

const (
	jwtCookieName = "jwt"
	defaultJWTTTL = 24 * time.Hour
)

// Decide Secure based on TLS (direct https) OR reverse proxy header.
func isSecureRequest(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

// SameSite: Strict is safest (best CSRF posture), but can break some flows.
// Lax is acceptable and common for session cookies.

// func jwtSameSiteMode() http.SameSite {
// 	return http.SameSiteLaxMode
// }

func jwtSameSiteMode() http.SameSite {
	return http.SameSiteStrictMode
}

func SetJWTCookie(w http.ResponseWriter, r *http.Request, token string, domain string) {
	SetJWTCookieWithTTL(w, r, token, domain, defaultJWTTTL)
}

func SetJWTCookieWithTTL(w http.ResponseWriter, r *http.Request, token string, domain string, ttl time.Duration) {
	secure := isSecureRequest(r)
	if ttl <= 0 {
		ttl = defaultJWTTTL
	}
	maxAgeSeconds := int(ttl.Seconds())
	if maxAgeSeconds <= 0 {
		maxAgeSeconds = 1
	}

	c := &http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAgeSeconds,
		Expires:  time.Now().Add(ttl),
		Secure:   secure,
		HttpOnly: true,
		SameSite: jwtSameSiteMode(),
	}

	// Only set Domain if provided (do NOT send an empty Domain attribute)
	if d := strings.TrimSpace(domain); d != "" {
		c.Domain = d
	}

	http.SetCookie(w, c)
}

func ClearJWTCookie(w http.ResponseWriter, r *http.Request, domain string) {
	secure := isSecureRequest(r)

	c := &http.Cookie{
		Name:     jwtCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   secure,
		HttpOnly: true,
		SameSite: jwtSameSiteMode(),
	}

	if d := strings.TrimSpace(domain); d != "" {
		c.Domain = d
	}

	http.SetCookie(w, c)
}
