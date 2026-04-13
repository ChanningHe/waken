package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func BearerAuth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			if subtle.ConstantTimeCompare([]byte(parts[1]), []byte(token)) != 1 {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
