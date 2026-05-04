package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClaudeProvider struct {
	apiKey     string
	httpClient *http.Client
}

func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *ClaudeProvider) Name() string {
	return "Anthropic-Claude"
}

func (p *ClaudeProvider) Ask(ctx context.Context, prompt string) (string, error) {
	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type payload struct {
		Model     string    `json:"model"`
		MaxTokens int       `json:"max_tokens"`
		Messages  []message `json:"messages"`
	}

	body := payload{
		Model:     "claude-3-opus-20240229",
		MaxTokens: 4096,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("fallo serializando request a Claude: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("fallo construyendo request Claude: %w", err)
	}

	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error de red contactando Anthropic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Anthropic respondio error HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("fallo decodificando respuesta Claude: %w", err)
	}

	contentArr, ok := result["content"].([]interface{})
	if !ok || len(contentArr) == 0 {
		return "", fmt.Errorf("formato de respuesta inesperado de Claude")
	}
	firstBlock := contentArr[0].(map[string]interface{})
	
	return fmt.Sprintf("%v", firstBlock["text"]), nil
}
