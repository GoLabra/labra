package handler

import (
	"net/http"
	"strings"

	"github.com/GoLabra/labra/config"
)

const defaultCSP = "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'"

// SecurityHeaders sets baseline HTTP security headers on every response.
// HSTS is only set in production and only when the request is HTTPS (direct TLS or X-Forwarded-Proto=https).
func SecurityHeaders(cfg *config.Config) func(http.Handler) http.Handler {
	csp := defaultCSP
	appEnv := ""

	if cfg != nil {
		appEnv = strings.TrimSpace(cfg.Environment)
		if v := strings.TrimSpace(cfg.CSPPolicy); v != "" {
			csp = v
		}
	}

	// Prevent header injection via newline characters
	if strings.ContainsAny(csp, "\r\n") {
		csp = defaultCSP
	}

	isProd := strings.EqualFold(appEnv, "prod") || strings.EqualFold(appEnv, "production")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()

			h.Set("X-Frame-Options", "DENY")
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-XSS-Protection", "0") // disabled in favor of CSP
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Content-Security-Policy", csp)
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

			// HSTS: production only, and only over HTTPS
			if isProd && isHTTPSRequest(r) {
				h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isHTTPSRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	// Common reverse-proxy header
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}
