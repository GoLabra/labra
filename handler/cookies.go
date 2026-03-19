package handler

import (
	"net/http"
	"strings"
	"time"
)

const (
	jwtCookieName          = "jwt"
	refreshTokenCookieName = "refresh_token"
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
	SetJWTCookieWithTTL(w, r, token, domain, 24*time.Hour)
}

func SetJWTCookieWithTTL(w http.ResponseWriter, r *http.Request, token string, domain string, ttl time.Duration) {
	secure := isSecureRequest(r)
	maxAge := int(ttl.Seconds())

	c := &http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
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

func SetRefreshTokenCookie(w http.ResponseWriter, r *http.Request, token string, domain string, ttl time.Duration) {
	secure := isSecureRequest(r)
	maxAge := int(ttl.Seconds())

	c := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  time.Now().Add(ttl),
		Secure:   secure,
		HttpOnly: true,
		SameSite: jwtSameSiteMode(),
	}

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

func ClearRefreshTokenCookie(w http.ResponseWriter, r *http.Request, domain string) {
	secure := isSecureRequest(r)

	c := &http.Cookie{
		Name:     refreshTokenCookieName,
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
