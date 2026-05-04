package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mikrotik-ai-failover/internal/ai"
	"mikrotik-ai-failover/internal/mikrotik"
)

func main() {
	// 1. SRE Observability: Logging Estructurado
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("Iniciando MikroTik AI Failover Agent en RouterOS v7 Container")

	// 2. Configuración (Normalmente inyectada por variables de entorno en el container de RouterOS)
	routerIP := os.Getenv("MIKROTIK_IP")
	if routerIP == "" {
		routerIP = "192.168.88.1" // Fallback local para pruebas
	}
	routerUser := os.Getenv("MIKROTIK_USER")
	routerPass := os.Getenv("MIKROTIK_PASS")
	if routerUser == "" {
		routerUser = "ai_agent"
		routerPass = "modo_cientifico"
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		logger.Error("Falta la variable de entorno GEMINI_API_KEY")
		os.Exit(1)
	}

	// 3. Inicialización de clientes
	mkClient := mikrotik.NewClient(routerIP, routerUser, routerPass)
	
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	aiClient, err := ai.NewGeminiClient(ctx, geminiKey)
	if err != nil {
		logger.Error("Fallo inicializando cliente Gemini", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer aiClient.Close()

	// 4. Bucle principal del Daemon
	ticker := time.NewTicker(30 * time.Second) // Poll interval
	defer ticker.Stop()

	logger.Info("Agente listo. Iniciando polling de tabla de ruteo...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Señal de apagado recibida. Cerrando agente de forma limpia...")
			return
		case <-ticker.C:
			pollNetworkState(ctx, logger, mkClient, aiClient)
		}
	}
}

func pollNetworkState(ctx context.Context, logger *slog.Logger, mkClient *mikrotik.Client, aiClient *ai.GeminiClient) {
	// PASO A: Obtener estado del router
	pollCtx, cancelPoll := context.WithTimeout(ctx, 10*time.Second)
	defer cancelPoll()

	routes, rawJSON, err := mkClient.GetRoutes(pollCtx)
	if err != nil {
		logger.Error("Error leyendo tabla de ruteo", slog.String("error", err.Error()))
		return
	}

	logger.Info("Tabla de ruteo leída exitosamente", slog.Int("routes_count", len(routes)))

	// Simulamos la detección de un evento de latencia grave en el ISP 1
	// En un entorno SRE real, esta métrica vendría de Prometheus o eBPF.
	latencySpikeDetected := true 
	
	if latencySpikeDetected {
		logger.Warn("Latencia crítica detectada en el ISP primario. Consultando a la IA para remediación...")

		// PASO B y C: Consultar al orquestador AI (Gemini) pasándole la topología
		aiCtx, cancelAI := context.WithTimeout(ctx, 30*time.Second)
		patchJSON, err := aiClient.EvaluateRouting(aiCtx, string(rawJSON))
		cancelAI()

		if err != nil {
			logger.Error("La IA falló al calcular el consenso", slog.String("error", err.Error()))
			return
		}

		logger.Info("Consenso de IA alcanzado", slog.String("remediation_payload", patchJSON))

		// PASO D: Inyectar la remediación en el Router
		// Para el PoC asumimos que la ruta principal es la que tiene ID "*1" o la buscamos.
		// Extraer el ID real de la ruta default (0.0.0.0/0) del ISP1
		var targetRouteID string
		for _, r := range routes {
			if r.Dst == "0.0.0.0/0" && r.Distance < 20 { // ISP Primario típico
				targetRouteID = r.ID
				break
			}
		}

		if targetRouteID == "" {
			// Fallback si no encontramos la ruta dinámicamente
			targetRouteID = "*1" 
		}

		logger.Info("Inyectando parche dinámico", slog.String("route_id", targetRouteID))
		
		patchCtx, cancelPatch := context.WithTimeout(ctx, 5*time.Second)
		err = mkClient.PatchRoute(patchCtx, targetRouteID, []byte(patchJSON))
		cancelPatch()

		if err != nil {
			logger.Error("Fallo inyectando configuración en el router", slog.String("error", err.Error()))
		} else {
			logger.Info("✅ Self-Healing ejecutado con éxito. Router reconfigurado.")
		}
	}
}
