package main

import (
	"context"
	"datacenter-reconciler/internal/ai"
	"datacenter-reconciler/internal/fabric"
	"log/slog"
	"os"
	"time"
)

func main() {
	// 1. Configuración de Logger SRE
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("Iniciando AI Consensus Orchestrator (AICO) Multi-Cloud")

	// 2. Inicialización de Proveedores IA Multi-Cloud
	// Obtenemos API Keys de variables de entorno (Security best practice)
	claudeKey := os.Getenv("ANTHROPIC_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if claudeKey == "" || geminiKey == "" {
		logger.Warn("API Keys no detectadas. El orquestador operará en modo Dry-Run o fallará si se intenta red de producción.")
	}

	// Instanciamos los modelos. En esta configuración:
	// - Claude es el 'Thinker' (diseña la configuración).
	// - Gemini es el 'Auditor' (verifica sintaxis YANG para SR Linux).
	thinker := ai.NewClaudeProvider(claudeKey)
	auditor := ai.NewGeminiProvider(geminiKey)

	// 3. Inicialización del Motor de Consenso
	consensusEngine := ai.NewConsensusEngine(logger, thinker, auditor)

	// 4. Inicialización del Controlador Fabric (Nokia SR Linux)
	targetRouter := "srl-edge-01.mgmt.local"
	gNMIController := fabric.NewGNMIController(logger, targetRouter)

	// Contexto SRE con timeout fuerte
	ctx := context.Background()
	// Simulación: Timeout global de 2 minutos para que las IAs debatan
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// --- SIMULACIÓN DEL EVENTO ---
	networkContext := `ALERTA CRÍTICA: Interface ethernet-1/1 en srl-edge-01 reportando 80% de dropped packets.
Patrón BGP Anómalo detectado hacia AS 65002.
Acción requerida: Generar política de Traffic Engineering (Local Preference = 50) para desviar tráfico vía ethernet-1/2.`

	logger.InfoContext(ctx, "Iniciando debate de IAs por evento de red", slog.String("evento", "Congestión BGP"))

	// Las IAs debaten y devuelven el JSON YANG consensuado
	finalYANG, err := consensusEngine.GenerateAutonomousRemediation(ctx, networkContext)
	if err != nil {
		logger.ErrorContext(ctx, "Fallo en el consenso de IA", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.InfoContext(ctx, "Consenso alcanzado con éxito, inyectando a la red...")

	// El controlador aplica el JSON generado autónomamente al router Nokia vía gNMI
	if err := gNMIController.ApplyAutonomousYANG(ctx, []byte(finalYANG)); err != nil {
		logger.ErrorContext(ctx, "Fallo inyectando configuración gNMI", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.InfoContext(ctx, "AICO completó el ciclo autónomo de mitigación. SR Linux actualizado.")
}
