package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GPTProvider struct {
	apiKey     string
	httpClient *http.Client
}

func NewGPTProvider(apiKey string) *GPTProvider {
	return &GPTProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // SRE Rule: Nunca llamadas de red sin timeout
		},
	}
}

func (p *GPTProvider) Name() string {
	return "OpenAI-GPT"
}

func (p *GPTProvider) Ask(ctx context.Context, prompt string) (string, error) {
	// Estructura estricta para evitar Heap Allocations masivas (Regla SRE #2)
	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type payload struct {
		Model    string    `json:"model"`
		Messages []message `json:"messages"`
	}

	body := payload{
		Model: "gpt-4o",
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("fallo serializando request a GPT: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("fallo construyendo request GPT: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error de red contactando OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI respondio error HTTP %d", resp.StatusCode)
	}

	// Parseo simple para PoC (Idealmente con structs estrictos)
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("fallo decodificando respuesta GPT: %w", err)
	}

	// Extraemos el texto de forma defensiva
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("formato de respuesta inesperado de GPT")
	}
	firstChoice := choices[0].(map[string]interface{})
	messageContent := firstChoice["message"].(map[string]interface{})
	
	return fmt.Sprintf("%v", messageContent["content"]), nil
}
