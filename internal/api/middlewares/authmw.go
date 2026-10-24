package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/velin-ai/velin/internal/admin/config"
	"github.com/velin-ai/velin/internal/api/models"
)

type contextKey string

const claimsKey contextKey = "claims"

// Auth validates the JWT access token sent in the "Authorization: Bearer <token>"
// header. On success the parsed claims are stored in the request context and can
// be read with ClaimsFromContext, UserIDFromContext or UsernameFromContext.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing_token")

				return
			}

			tokenString := strings.TrimSpace(
				strings.TrimPrefix(authHeader, "Bearer "),
			)

			if tokenString == "" {
				writeError(w, http.StatusUnauthorized, "missing_token")

				return
			}

			secret, err := config.GetJWTSecret()

			if err != nil {
				log.Println(err)

				writeError(w, http.StatusInternalServerError, "secret_error")

				return
			}

			claims := &models.Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(_ *jwt.Token) (any, error) {
					return []byte(secret), nil
				},
				jwt.WithValidMethods([]string{
					jwt.SigningMethodHS256.Alg(),
				}),
			)

			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid_token")

				return
			}

			ctx := context.WithValue(
				r.Context(),
				claimsKey,
				claims,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// ClaimsFromContext returns the claims stored by the Auth middleware.
func ClaimsFromContext(ctx context.Context) (*models.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*models.Claims)

	return claims, ok
}

// UserIDFromContext returns the user id stored by the Auth middleware.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	claims, ok := ClaimsFromContext(ctx)

	if !ok {
		return 0, false
	}

	return claims.UserID, true
}

// UsernameFromContext returns the username stored by the Auth middleware.
func UsernameFromContext(ctx context.Context) (string, bool) {
	claims, ok := ClaimsFromContext(ctx)

	if !ok {
		return "", false
	}

	return claims.Username, true
}

func writeError(w http.ResponseWriter, status int, code string) {
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]any{
		"error": code,
	})
}
