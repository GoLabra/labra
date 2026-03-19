package handler

import (
	"app/domain/svc"
	"errors"
	"log"
	"net/http"

	"github.com/GoLabra/labra/config"
	sharedHandler "github.com/GoLabra/labra/handler"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	appConfig, ok := r.Context().Value("config").(*config.AppConfig)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshCookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, _, accessToken, newRefreshToken, accessTTL, refreshTTL, err := rotateUserRefreshSession(r.Context(), refreshCookie.Value, appConfig)
	if err != nil {
		if errors.Is(err, svc.ErrInvalidRefreshToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("error rotating refresh token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sharedHandler.SetJWTCookieWithTTL(w, r, accessToken, "", accessTTL)
	sharedHandler.SetRefreshTokenCookie(w, r, newRefreshToken, "", refreshTTL)
	writeTokenResponse(w, accessToken)
}
