package api

import (
	"net/http"
	"os"
)

func APIKeyMiddleware(next http.Handler) http.Handler {
	requiredKey := os.Getenv("API_KEY")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" || apiKey != requiredKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
