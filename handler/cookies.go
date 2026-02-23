package handler

import (
	"net/http"
	"strings"
	"time"
)

const (
	jwtCookieName = "jwt"
	jwtMaxAgeSecs = 60 * 60 * 24 // 24h
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
	secure := isSecureRequest(r)

	c := &http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   jwtMaxAgeSecs,
		Expires:  time.Now().Add(time.Duration(jwtMaxAgeSecs) * time.Second),
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
