package apikey

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func GenerateAPIKey(env string) (string, error) {
	if env != "live" && env != "test" {
		return "", errors.New("invalid API key environment")
	}

	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return "vln_" + env + "_" + hex.EncodeToString(b), nil
}
