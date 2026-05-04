package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey     string
	httpClient *http.Client
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *GeminiProvider) Name() string {
	return "Google-Gemini"
}

func (p *GeminiProvider) Ask(ctx context.Context, prompt string) (string, error) {
	type part struct {
		Text string `json:"text"`
	}
	type content struct {
		Parts []part `json:"parts"`
	}
	type payload struct {
		Contents []content `json:"contents"`
	}

	body := payload{
		Contents: []content{
			{Parts: []part{{Text: prompt}}},
		},
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("fallo serializando request a Gemini: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent?key=%s", p.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("fallo construyendo request Gemini: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error de red contactando Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Google respondio error HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("fallo decodificando respuesta Gemini: %w", err)
	}

	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("formato de respuesta inesperado de Gemini")
	}
	firstCandidate := candidates[0].(map[string]interface{})
	contentMap := firstCandidate["content"].(map[string]interface{})
	parts := contentMap["parts"].([]interface{})
	firstPart := parts[0].(map[string]interface{})

	return fmt.Sprintf("%v", firstPart["text"]), nil
}
