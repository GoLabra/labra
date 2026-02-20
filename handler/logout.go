package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/ent"
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

	adminEntClient, ok := r.Context().Value(constants.AdminEntClientContextValue).(*ent.Client)
	if !ok || adminEntClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("admin ent client not found")
		return
	}

	startTokenRevocationCleanupJob(adminEntClient)

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
		jwt.WithKey(jwt_hs.HS256, []byte(appConfig.Secrets.SecretKey)),
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
		subjectType = jwtrefresh.SubjectTypeAdmin
	}

	expiresAt, err := claimUnixTime(claims, "exp")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("missing/invalid exp claim: %v", err)
		return
	}

	err = tokenrevocation.RevokeToken(
		r.Context(),
		adminEntClient,
		adminEntClient.DialectName(),
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

	ClearJWTCookie(w, r, "")
	ClearCSRFCookie(w, isSecureRequest(r))

	response, _ := json.Marshal(map[string]string{
		"status": "logged_out",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
