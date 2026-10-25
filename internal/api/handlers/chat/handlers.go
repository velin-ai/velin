package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/velin-ai/velin/internal/api/ai"
	"github.com/velin-ai/velin/internal/api/middlewares"
	"github.com/velin-ai/velin/internal/api/models"
	"github.com/velin-ai/velin/internal/store"
)

const chatHistoryLimit = 20

func CreateNewChat(
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

	var req models.CreateNewChatReq

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

	if len(req.Title) > 75 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "title_is_too_long",
		})

		return
	}

	var chatID int64

	err = store.DB.QueryRow(
		r.Context(),
		`
		INSERT INTO chats (user_id,title)
		VALUES ($1, $2)
		RETURNING id;
		`,
		claims.UserID,
		req.Title,
	).Scan(&chatID)

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})
		return
	}

	w.WriteHeader(201)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "chat_created_successfully",
		"id":      chatID,
		"title":   req.Title,
	})

}

func ListMyChats(
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
		SELECT id, title, created_at
		FROM chats
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

	chats := make([]models.Chat, 0)

	for rows.Next() {
		var chat models.Chat

		err = rows.Scan(
			&chat.ID,
			&chat.Title,
			&chat.CreatedAt,
		)

		if err != nil {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "database_error",
			})

			return
		}

		chats = append(chats, chat)
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
		"chats": chats,
	})
}

func ChatMessage(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	info, ok := middlewares.APIKeyFromContext(r.Context())

	if !ok {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_api_key",
		})

		return
	}

	chatID, err := strconv.ParseInt(
		chi.URLParam(r, "chat_id"), 10, 64,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_chat_id",
		})

		return
	}

	var req models.ChatMessageReq

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "invalid_request",
		})

		return
	}

	if strings.TrimSpace(req.Model) == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "model_is_required",
		})

		return
	}

	if strings.TrimSpace(req.Content) == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "content_is_required",
		})

		return
	}

	var chatExists bool

	err = store.DB.QueryRow(
		r.Context(),
		`
		SELECT EXISTS (
			SELECT 1
			FROM chats
			WHERE id = $1 AND user_id = $2
		);
		`,
		chatID,
		info.UserID,
	).Scan(&chatExists)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	if !chatExists {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "chat_not_found",
		})

		return
	}

	rows, err := store.DB.Query(
		r.Context(),
		`
		SELECT role, content
		FROM (
			SELECT id, role, content
			FROM messages
			WHERE chat_id = $1
			ORDER BY id DESC
			LIMIT $2
		) AS recent
		ORDER BY id;
		`,
		chatID,
		chatHistoryLimit,
	)

	if err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	defer rows.Close()

	history := make([]ai.Message, 0, chatHistoryLimit+1)

	for rows.Next() {
		var message ai.Message

		err = rows.Scan(
			&message.Role,
			&message.Content,
		)

		if err != nil {
			w.WriteHeader(500)

			json.NewEncoder(w).Encode(map[string]any{
				"error": "database_error",
			})

			return
		}

		history = append(history, message)
	}

	if err := rows.Err(); err != nil {
		w.WriteHeader(500)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "database_error",
		})

		return
	}

	rows.Close()

	history = append(history, ai.Message{
		Role:    "user",
		Content: req.Content,
	})

	answer, err := ai.ChatWithOllama(
		req.Model,
		history,
		false,
	)

	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusBadGateway)

		json.NewEncoder(w).Encode(map[string]any{
			"error": "ai_error",
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

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO messages (chat_id, role, content)
		VALUES ($1, 'user', $2), ($1, 'assistant', $3);
		`,
		chatID,
		req.Content,
		answer,
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
		"message":  "chat_completed_successfully",
		"response": answer,
	})
}
