package middleware

import (
	"context"
	"net/http"
	"strings"

	"todo/internal/token"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := token.Parse(tokenStr)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if claims["type"] != "access" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, uint(sub))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
