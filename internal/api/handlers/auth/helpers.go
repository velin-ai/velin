package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

const (
	refreshTokenByteLength = 32
	refreshTokenTTL        = 30 * 24 * time.Hour
)

// CreateRandomToken returns a cryptographically secure, URL-safe random token
// suitable for use as a refresh token.
func CreateRandomToken() (string, error) {
	token := make([]byte, refreshTokenByteLength)

	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}

// HashToken returns the hex encoded SHA-256 hash of token. Refresh tokens are
// stored hashed so that a database leak does not expose usable tokens.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}
