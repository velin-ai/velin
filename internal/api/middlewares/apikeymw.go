package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/velin-ai/velin/internal/api/models"
	"github.com/velin-ai/velin/internal/store"
)

const apiKeyContextKey contextKey = "api_key"

// APIKeyInfo holds the authenticated API key and the id of its owner.
type APIKeyInfo struct {
	APIKey models.APIKey
	UserID int64
}

// APIKeyAuth validates the API key sent in the "Authorization: Bearer <api-key>"
// header against the api_keys table. On success the matched API key and its
// user id are stored in the request context and can be read with
// APIKeyFromContext.
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing_api_key")

				return
			}

			apiKey := strings.TrimSpace(
				strings.TrimPrefix(authHeader, "Bearer "),
			)

			if apiKey == "" {
				writeError(w, http.StatusUnauthorized, "missing_api_key")

				return
			}

			hash := sha256.Sum256([]byte(apiKey))
			keyHash := hex.EncodeToString(hash[:])

			var info APIKeyInfo

			err := store.DB.QueryRow(
				r.Context(),
				`
				SELECT id, title, user_id, expires_at, created_at
				FROM api_keys
				WHERE key_hash = $1
				  AND (expires_at IS NULL OR expires_at > NOW());
				`,
				keyHash,
			).Scan(
				&info.APIKey.ID,
				&info.APIKey.Title,
				&info.UserID,
				&info.APIKey.ExpiresAt,
				&info.APIKey.CreatedAt,
			)

			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					writeError(w, http.StatusUnauthorized, "invalid_api_key")

					return
				}

				writeError(w, http.StatusInternalServerError, "database_error")

				return
			}

			ctx := context.WithValue(
				r.Context(),
				apiKeyContextKey,
				info,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// APIKeyFromContext returns the API key and user id stored by the APIKeyAuth middleware.
func APIKeyFromContext(ctx context.Context) (APIKeyInfo, bool) {
	info, ok := ctx.Value(apiKeyContextKey).(APIKeyInfo)

	return info, ok
}
