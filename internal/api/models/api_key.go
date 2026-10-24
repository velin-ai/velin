package models

import "time"

type CreateAPIKeyReq struct {
	Title     string     `json:"title"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type APIKey struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}
