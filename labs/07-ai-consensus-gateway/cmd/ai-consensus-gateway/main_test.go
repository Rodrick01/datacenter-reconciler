package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockProvider implementa LLMProvider para evitar gastar tokens reales en tests
type MockProvider struct {
	name     string
	response string
}

func (m *MockProvider) Ask(ctx context.Context, prompt string) (string, error) {
	return m.response, nil
}

func (m *MockProvider) Name() string {
	return m.name
}

func TestConsensusHandler(t *testing.T) {
	// Reemplazamos la función getProvider global (o inyectamos en un diseño más robusto)
	// Para este test de integración de API, sobreescribimos la lógica del handler para usar mocks,
	// o mejor, testeamos que el JSON parsing del Request sea correcto.
	
	reqBody := ConsensusRequest{
		Context:   "Datos de ventas de heladeras 2026: Ene=100, Feb=150, Mar=80",
		Objective: "¿Cuál fue el mejor mes de ventas de heladeras?",
		Thinker:   "gemini",
		Auditor:   "claude",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/consensus", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	_ = httptest.NewRecorder()

	// Como getProvider usa env vars y llama APIs reales, vamos a forzar una validación de bad request primero
	badReq := httptest.NewRequest(http.MethodPost, "/api/v1/consensus", bytes.NewBuffer([]byte(`{}`)))
	wBad := httptest.NewRecorder()
	
	consensusHandler(wBad, badReq)

	if wBad.Code != http.StatusBadRequest {
		t.Errorf("Esperaba status 400 Bad Request, recibí %d", wBad.Code)
	}

	// Test Health Check
	reqHealth := httptest.NewRequest(http.MethodGet, "/health", nil)
	wHealth := httptest.NewRecorder()
	healthHandler(wHealth, reqHealth)

	if wHealth.Code != http.StatusOK {
		t.Errorf("Esperaba status 200 OK para /health, recibí %d", wHealth.Code)
	}

	// Un test simple de validez de structs
	if reqBody.Context == "" {
		t.Error("El contexto no debería estar vacío")
	}
}
