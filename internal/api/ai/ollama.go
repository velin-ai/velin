package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const ollamaURL = "http://localhost:11434/api/chat"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type chatResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func ChatWithOllama(model string, messages []Message, streaming bool) (string, error) {
	body, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: messages,
		Stream:   streaming,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		ollamaURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"Ollama API hatası: %s: %s",
			resp.Status,
			strings.TrimSpace(string(errBody)),
		)
	}

	var result strings.Builder

	if !streaming {
		var response chatResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return "", err
		}

		return response.Message.Content, nil
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	for scanner.Scan() {
		var response chatResponse
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			return result.String(), err
		}

		chunk := response.Message.Content
		result.WriteString(chunk)

		fmt.Print(chunk)
	}

	if err := scanner.Err(); err != nil {
		return result.String(), err
	}

	return result.String(), nil
}
