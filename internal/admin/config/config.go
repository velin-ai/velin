package config

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	Service   = "velin"
	Account   = "admin-password"
	JWTSecret = "jwt_secret"
)

func GetJWTSecret() (string, error) {
	jwtSecret, err := keyring.Get(Service, JWTSecret)

	if errors.Is(err, keyring.ErrNotFound) {
		return "", fmt.Errorf("jwt secret is not configured: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("accessing jwt secret: %w", err)
	}

	return jwtSecret, nil
}
