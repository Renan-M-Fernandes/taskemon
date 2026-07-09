package api

import (
	"net/http"
)

func EnableCORS(next http.Handler, origins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestOrigin := r.Header.Get("Origin")

		for _, allowed := range origins {
			if allowed == "*" || allowed == requestOrigin {
				w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
				w.Header().Set("Vary", "Origin")
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
