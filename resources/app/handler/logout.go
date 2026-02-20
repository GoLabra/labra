package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/ent"
	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/jwtrefresh"
	"github.com/GoLabra/labra/tokenrevocation"

	jwt_hs "github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	appConfig, ok := r.Context().Value("config").(*config.AppConfig)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("config not found")
		return
	}

	entClient, ok := r.Context().Value(constants.EntClientContextValue).(*ent.Client)
	if !ok || entClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ent client not found")
		return
	}

	tokenString := extractJWTFromRequest(r)
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("token not found")
		return
	}

	appJWT, err := jwt.ParseString(
		tokenString,
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithKey(jwt_hs.HS256, []byte(appConfig.SecretKey)),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("error parsing token: %v", err)
		return
	}

	claims, err := appJWT.AsMap(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("error parsing claims: %v", err)
		return
	}

	subjectEmail := claimStringValue(claims, jwtrefresh.ClaimSubject)
	if subjectEmail == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("missing sub claim")
		return
	}

	subjectType := claimStringValue(claims, jwtrefresh.ClaimSubjectType)
	if subjectType == "" {
		subjectType = jwtrefresh.SubjectTypeUser
	}

	expiresAt, err := claimUnixTime(claims, "exp")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("missing/invalid exp claim: %v", err)
		return
	}

	err = tokenrevocation.RevokeToken(
		r.Context(),
		entClient,
		appConfig.DBDialect,
		tokenString,
		subjectEmail,
		subjectType,
		expiresAt,
		time.Now().UTC(),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to revoke token during logout: %v", err)
		return
	}

	response, _ := json.Marshal(map[string]string{
		"status": "logged_out",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func claimStringValue(claims map[string]interface{}, key string) string {
	raw, ok := claims[key]
	if !ok || raw == nil {
		return ""
	}

	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func claimUnixTime(claims map[string]interface{}, key string) (time.Time, error) {
	raw, ok := claims[key]
	if !ok || raw == nil {
		return time.Time{}, fmt.Errorf("missing %s claim", key)
	}

	switch v := raw.(type) {
	case time.Time:
		return v.UTC(), nil
	case *time.Time:
		if v == nil {
			return time.Time{}, fmt.Errorf("invalid %s claim", key)
		}
		return v.UTC(), nil
	case int64:
		return time.Unix(v, 0).UTC(), nil
	case int:
		return time.Unix(int64(v), 0).UTC(), nil
	case int32:
		return time.Unix(int64(v), 0).UTC(), nil
	case float64:
		return time.Unix(int64(v), 0).UTC(), nil
	case float32:
		return time.Unix(int64(v), 0).UTC(), nil
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid %s claim: %w", key, err)
		}
		return time.Unix(n, 0).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("invalid %s claim type %T", key, raw)
	}
}
