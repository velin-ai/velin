package models

import "time"

type CreateNewChatReq struct {
	Title string `json:"title"`
}

type Chat struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessageReq struct {
	Model   string `json:"model"`
	Content string `json:"content"`
}
