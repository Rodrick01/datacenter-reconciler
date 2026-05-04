package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Rodrick01/ai-consensus-gateway/internal/ai"
)

type ConsensusRequest struct {
	Context   string `json:"context"`
	Objective string `json:"objective"`
	Thinker   string `json:"thinker,omitempty"`   // ej: "gemini", "claude", "gpt"
	Auditor   string `json:"auditor,omitempty"`   // ej: "gemini", "claude", "gpt"
}

type ConsensusResponse struct {
	ConsensusResult string `json:"consensus_result"`
	Error           string `json:"error,omitempty"`
}

var logger *slog.Logger

func getProvider(name string) ai.LLMProvider {
	switch name {
	case "claude":
		return ai.NewClaudeProvider(os.Getenv("ANTHROPIC_API_KEY"))
	case "gpt":
		return ai.NewGPTProvider(os.Getenv("OPENAI_API_KEY"))
	case "gemini":
		fallthrough
	default:
		// Default to Gemini si no se especifica
		return ai.NewGeminiProvider(os.Getenv("GEMINI_API_KEY"))
	}
}

func consensusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido. Usa POST.", http.StatusMethodNotAllowed)
		return
	}

	var req ConsensusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Context == "" || req.Objective == "" {
		http.Error(w, "context y objective son obligatorios", http.StatusBadRequest)
		return
	}

	thinker := getProvider(req.Thinker)
	auditor := getProvider(req.Auditor)

	engine := ai.NewConsensusEngine(logger, thinker, auditor)

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	result, err := engine.GenerateConsensus(ctx, req.Context, req.Objective)
	
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ConsensusResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(ConsensusResponse{ConsensusResult: result})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed. Use GET.", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "UP", "service": "ai-consensus-gateway"})
}

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/consensus", consensusHandler)
	mux.HandleFunc("/health", healthHandler)

	logger.Info("Iniciando AI Consensus Gateway", slog.String("puerto", port))
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 130 * time.Second, // Alto porque LLMs son lentos
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Fallo al iniciar el servidor", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
