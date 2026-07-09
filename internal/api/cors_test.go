package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnableCORSAddsHeaders(t *testing.T) {
	cfg := loadConfig(t)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:8123")
	rr := httptest.NewRecorder()

	EnableCORS(next, cfg.Server.CORSOrigins).ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatal("Access-Control-Allow-Origin should be set")
	}
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("Access-Control-Allow-Methods should be set")
	}
	if rr.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatal("Access-Control-Allow-Headers should be set")
	}
}

func TestEnableCORSHandlesOptions(t *testing.T) {
	cfg := loadConfig(t)
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/tasks/ash", nil)
	rr := httptest.NewRecorder()

	EnableCORS(next, cfg.Server.CORSOrigins).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status mismatch: got %d, expect %d", rr.Code, http.StatusNoContent)
	}
	if called {
		t.Fatal("next handler should not be called for OPTIONS")
	}
}
