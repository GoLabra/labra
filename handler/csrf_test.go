package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFMiddleware_CookieJWTRequiresCSRF(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/anything", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: "dummy"})
	rr := httptest.NewRecorder()

	CSRFMiddleware(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func TestCSRFMiddleware_AuthorizationSkipsCSRF(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/anything", nil)
	req.Header.Set("Authorization", "Bearer dummy")
	rr := httptest.NewRecorder()

	CSRFMiddleware(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
