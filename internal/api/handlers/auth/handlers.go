package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/velin-ai/velin/internal/admin/config"
	"github.com/velin-ai/velin/internal/api/models"
	"github.com/velin-ai/velin/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_request",
		})

		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "username_or_password_is_required",
		})

		return
	}

	if len(req.Username) > 20 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "username_too_long",
		})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		12,
	)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "bcrypt_hash_error",
		})

		return
	}

	_, err = store.DB.Exec(
		r.Context(),
		`
		INSERT INTO users (username, password)
		VALUES ($1,$2);
		`,
		req.Username,
		string(hashedPassword),
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				w.WriteHeader(409)

				json.NewEncoder(w).Encode(map[string]any{
					"error": "user_already_exists",
				})

				return
			}
		}

		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	w.WriteHeader(201)

	json.NewEncoder(w).Encode(map[string]any{
		"error": "user_created",
	})
}

func LoginHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_req",
		})

		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "username_or_password_is_required",
		})

		return
	}

	if len(req.Username) > 20 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "username_too_long",
		})

		return
	}

	var (
		userID     int64
		passwordDB string
	)

	err = store.DB.QueryRow(
		r.Context(),
		`
		SELECT id, password
		FROM users
		WHERE username = $1;
		`,
		req.Username,
	).Scan(
		&userID,
		&passwordDB,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(401)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid_credentials",
			})

			return
		}

		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordDB),
		[]byte(req.Password),
	)

	if err != nil {
		w.WriteHeader(401)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_credentials",
		})

		return
	}

	claims := models.Claims{
		UserID:   userID,
		Username: req.Username,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(30 * time.Minute),
			),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := config.GetJWTSecret()

	if err != nil {
		log.Println(err)

		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "secret_error",
		})

		return
	}

	accessToken, err := access.SignedString(secret)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "jwt_error",
		})

		return
	}

	randomToken, err := CreateRandomToken()

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "generate_refresh_token_error",
		})

		return
	}

	_, err = store.DB.Exec(
		r.Context(),
		`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1,$2,$3);
		`,
		userID,
		HashToken(randomToken),
		time.Now().Add(refreshTokenTTL),
	)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(map[string]any{
		"message":      "login_successfully",
		"accessToken":  accessToken,
		"refreshToken": randomToken,
		"user_id":      userID,
	})
}

func RefreshHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RefreshReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_request",
		})

		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "refresh_token_is_required",
		})

		return
	}

	tokenHash := HashToken(req.RefreshToken)

	var (
		userID    int64
		username  string
		expiresAt time.Time
	)

	err = store.DB.QueryRow(
		r.Context(),
		`
		SELECT u.id, u.username, rt.expires_at
		FROM refresh_tokens rt
		JOIN users u ON u.id = rt.user_id
		WHERE rt.token = $1;
		`,
		tokenHash,
	).Scan(
		&userID,
		&username,
		&expiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(401)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid_refresh_token",
			})

			return
		}

		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	if time.Now().After(expiresAt) {
		_, _ = store.DB.Exec(
			r.Context(),
			`DELETE FROM refresh_tokens WHERE token = $1;`,
			tokenHash,
		)

		w.WriteHeader(401)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "refresh_token_expired",
		})

		return
	}

	claims := models.Claims{
		UserID:   userID,
		Username: username,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(30 * time.Minute),
			),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := config.GetJWTSecret()

	if err != nil {
		log.Println(err)

		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "secret_error",
		})

		return
	}

	accessToken, err := access.SignedString(secret)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "jwt_error",
		})

		return
	}

	randomToken, err := CreateRandomToken()

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "generate_refresh_token_error",
		})

		return
	}

	tx, err := store.DB.Begin(r.Context())

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	defer tx.Rollback(r.Context())

	tag, err := tx.Exec(
		r.Context(),
		`DELETE FROM refresh_tokens WHERE token = $1;`,
		tokenHash,
	)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	if tag.RowsAffected() == 0 {
		w.WriteHeader(401)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_refresh_token",
		})

		return
	}

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1,$2,$3);
		`,
		userID,
		HashToken(randomToken),
		time.Now().Add(refreshTokenTTL),
	)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	err = tx.Commit(r.Context())

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(map[string]any{
		"message":      "refresh_successfully",
		"accessToken":  accessToken,
		"refreshToken": randomToken,
		"user_id":      userID,
	})
}
