package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/jwtrefresh"
)

type RefreshRequest struct {
	RefreshToken    string `json:"refresh_token"`
	RefreshTokenAlt string `json:"refreshToken"`
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	appConfig, ok := r.Context().Value("config").(*config.AppConfig)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("config not found")
		return
	}

	adminEntClient, ok := r.Context().Value(constants.AdminEntClientContextValue).(*ent.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("admin ent client not found")
		return
	}

	refreshToken, err := extractRefreshTokenFromRequest(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid refresh request: %v", err)
		return
	}

	claims, err := jwtrefresh.ParseRefreshToken(refreshToken, appConfig.SecretKey, jwtrefresh.SubjectTypeAdmin)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("invalid refresh token: %v", err)
		return
	}

	tx, err := adminEntClient.Tx(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to start refresh transaction: %v", err)
		return
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	stored, err := jwtrefresh.LoadRefreshToken(r.Context(), tx, appConfig.DBDialect, refreshToken)
	if err != nil {
		if errors.Is(err, jwtrefresh.ErrRefreshTokenNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed loading refresh token: %v", err)
		return
	}

	if stored.SubjectType != claims.SubjectType || stored.SubjectEmail != claims.Subject || stored.RoleName != claims.Role {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if stored.RevokedAtUnix.Valid || stored.ExpiresAtUnix <= time.Now().UTC().Unix() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	revoked, err := jwtrefresh.RevokeRefreshToken(r.Context(), tx, appConfig.DBDialect, refreshToken, time.Now().UTC())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed revoking refresh token: %v", err)
		return
	}
	if !revoked {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pair, err := jwtrefresh.IssueTokenPair(
		appConfig.SecretKey,
		claims.Subject,
		claims.Role,
		claims.SubjectType,
		appConfig.AccessTokenTTL,
		appConfig.RefreshTokenTTL,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed issuing new token pair: %v", err)
		return
	}

	err = jwtrefresh.SaveRefreshToken(
		r.Context(),
		tx,
		appConfig.DBDialect,
		pair.RefreshToken,
		claims.Subject,
		claims.SubjectType,
		claims.Role,
		pair.RefreshExpiresAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed storing rotated refresh token: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed committing refresh transaction: %v", err)
		return
	}
	committed = true

	SetJWTCookieWithTTL(w, r, pair.AccessToken, "", jwtrefresh.EffectiveAccessTTL(appConfig.AccessTokenTTL))

	response, _ := json.Marshal(map[string]string{
		"token":         pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func extractRefreshTokenFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	var payload RefreshRequest
	err := decodeRefreshJSON(r, &payload)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	token := strings.TrimSpace(payload.RefreshToken)
	if token == "" {
		token = strings.TrimSpace(payload.RefreshTokenAlt)
	}
	if token == "" {
		token = extractBearerToken(r.Header.Get("Authorization"))
	}
	if token == "" {
		if cookie, err := r.Cookie("refresh_jwt"); err == nil {
			token = strings.TrimSpace(cookie.Value)
		}
	}
	if token == "" {
		return "", errors.New("missing refresh token")
	}

	return token, nil
}

func decodeRefreshJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return io.EOF
	}

	return json.Unmarshal(body, dst)
}

func extractBearerToken(authorizationHeader string) string {
	auth := strings.TrimSpace(authorizationHeader)
	if auth == "" {
		return ""
	}
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return ""
	}
	return strings.TrimSpace(auth[len("Bearer "):])
}
