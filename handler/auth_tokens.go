package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/domain/svc"
	"github.com/GoLabra/labra/entgql/ent"
	"github.com/golang-jwt/jwt"
)

var persistAdminRefreshSession = defaultPersistAdminRefreshSession
var rotateAdminRefreshSession = defaultRotateAdminRefreshSession

func accessTokenTTL(appConfig *config.AppConfig) time.Duration {
	if appConfig == nil || appConfig.AccessTokenTTLMin <= 0 {
		return 15 * time.Minute
	}
	return time.Duration(appConfig.AccessTokenTTLMin) * time.Minute
}

func refreshTokenTTL(appConfig *config.AppConfig) time.Duration {
	if appConfig == nil || appConfig.RefreshTokenTTLHours <= 0 {
		return 7 * 24 * time.Hour
	}
	return time.Duration(appConfig.RefreshTokenTTLHours) * time.Hour
}

func issueAccessToken(appConfig *config.AppConfig, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(accessTokenTTL(appConfig)).Unix(),
		"sub":  email,
		"role": role,
	})

	return token.SignedString([]byte(appConfig.SecretKey))
}

func newOpaqueToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashOpaqueToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func writeTokenResponse(w http.ResponseWriter, accessToken string) {
	response, _ := json.Marshal(map[string]string{
		"token": accessToken,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func issueAdminSession(ctx context.Context, appConfig *config.AppConfig, user *ent.AdminUser, role *ent.Role) (string, string, time.Duration, time.Duration, error) {
	accessToken, err := issueAccessToken(appConfig, user.Email, role.Name)
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshToken, err := newOpaqueToken()
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshTTL := refreshTokenTTL(appConfig)
	if err := persistAdminRefreshSession(ctx, user.ID, hashOpaqueToken(refreshToken), time.Now().Add(refreshTTL)); err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, refreshToken, accessTokenTTL(appConfig), refreshTTL, nil
}

func getAdminService(ctx context.Context) (*svc.Service, error) {
	service, ok := ctx.Value(constants.AdminServiceContextValue).(*svc.Service)
	if !ok || service == nil {
		return nil, errors.New(svc.ErrServiceNotSetInContext)
	}
	return service, nil
}

func defaultPersistAdminRefreshSession(ctx context.Context, adminUserID, tokenHash string, expiresAt time.Time) error {
	service, err := getAdminService(ctx)
	if err != nil {
		return err
	}

	return service.AdminRefreshToken.Persist(ctx, adminUserID, tokenHash, expiresAt)
}

func defaultRotateAdminRefreshSession(ctx context.Context, rawRefreshToken string, appConfig *config.AppConfig) (*ent.AdminUser, *ent.Role, string, string, time.Duration, time.Duration, error) {
	service, err := getAdminService(ctx)
	if err != nil {
		return nil, nil, "", "", 0, 0, err
	}

	adminUser, role, err := service.AdminRefreshToken.Use(ctx, hashOpaqueToken(rawRefreshToken))
	if err != nil {
		return nil, nil, "", "", 0, 0, err
	}

	accessToken, newRefreshToken, accessTTL, refreshTTL, err := issueAdminSession(ctx, appConfig, adminUser, role)
	if err != nil {
		return nil, nil, "", "", 0, 0, err
	}

	return adminUser, role, accessToken, newRefreshToken, accessTTL, refreshTTL, nil
}
