package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/entgql/domain/svc"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	appConfig, ok := r.Context().Value("config").(*config.AppConfig)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil || refreshCookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, _, accessToken, newRefreshToken, accessTTL, refreshTTL, err := rotateAdminRefreshSession(r.Context(), refreshCookie.Value, appConfig)
	if err != nil {
		if errors.Is(err, svc.ErrInvalidRefreshToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("error rotating refresh token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookieDomain := ""
	SetJWTCookieWithTTL(w, r, accessToken, cookieDomain, accessTTL)
	SetRefreshTokenCookie(w, r, newRefreshToken, cookieDomain, refreshTTL)
	writeTokenResponse(w, accessToken)
}
