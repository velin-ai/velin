package apikey

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/velin-ai/velin/internal/api/middlewares"
	"github.com/velin-ai/velin/internal/api/models"
	"github.com/velin-ai/velin/internal/store"
)

func CreateNewAPIKeyHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := middlewares.ClaimsFromContext(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "unauthorized",
		})

		return
	}

	var req models.CreateAPIKeyReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_request",
		})

		return
	}

	if strings.TrimSpace(req.Title) == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "title_is_required",
		})

		return
	}

	if len(req.Title) > 50 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "title_is_too_long",
		})

		return
	}

	if req.ExpiresAt != nil && !req.ExpiresAt.After(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "expires_at_must_be_in_the_future",
		})
		return
	}

	// create api api key
	api_key, err := GenerateAPIKey("live")

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "generate_api_key_error",
		})

		return
	}

	hash := sha256.Sum256([]byte(api_key))
	keyHash := hex.EncodeToString(hash[:])

	// save to database
	_, err = store.DB.Exec(
		r.Context(),
		`
	    INSERT INTO api_keys (title, user_id, key_hash, expires_at)
	    VALUES ($1, $2, $3, $4);
		`,
		req.Title,
		claims.UserID,
		keyHash,
		req.ExpiresAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			w.WriteHeader(http.StatusConflict)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "api_key_already_exists",
			})

			return
		}

		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	w.WriteHeader(201)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "api_key_created",
		"api_key": api_key,
		"title":   req.Title,
	})
}

func ListMyAPIKeys(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := middlewares.ClaimsFromContext(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "unauthorized",
		})

		return
	}

	rows, err := store.DB.Query(
		r.Context(),
		`
		SELECT id, title, expires_at, created_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC;
		`,
		claims.UserID,
	)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	defer rows.Close()

	apiKeys := make([]models.APIKey, 0)

	for rows.Next() {
		var apiKey models.APIKey

		err = rows.Scan(
			&apiKey.ID,
			&apiKey.Title,
			&apiKey.ExpiresAt,
			&apiKey.CreatedAt,
		)

		if err != nil {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "database_error",
			})

			return
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(map[string]any{
		"api_keys": apiKeys,
	})
}

func DeleteAPIKey(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := middlewares.ClaimsFromContext(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "unauthorized",
		})

		return
	}

	apiKeyID, err := strconv.ParseInt(
		chi.URLParam(r, "id"), 10, 64,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_api_key_id",
		})

		return
	}

	tag, err := store.DB.Exec(
		r.Context(),
		`
		DELETE FROM api_keys
		WHERE id = $1 AND user_id = $2;
		`,
		apiKeyID,
		claims.UserID,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	if tag.RowsAffected() == 0 {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "api_key_not_found",
		})

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "api_key_deleted",
	})
}
