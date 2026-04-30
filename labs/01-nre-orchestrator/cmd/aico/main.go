package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"datacenter-reconciler/internal/ai"
	"datacenter-reconciler/internal/fabric"
	"datacenter-reconciler/internal/sensor"
)

func main() {
	// 1. Configuración de Logger SRE
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("Iniciando AICO Gateway (AI Consensus Orchestrator + Telemetry)")

	// 2. Inicialización de Proveedores IA Multi-Cloud
	claudeKey := os.Getenv("ANTHROPIC_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if claudeKey == "" || geminiKey == "" {
		logger.Warn("API Keys no detectadas. El orquestador fallará al contactar las nubes si no están mockeadas.")
	}

	thinker := ai.NewClaudeProvider(claudeKey)
	auditor := ai.NewGeminiProvider(geminiKey)
	consensusEngine := ai.NewConsensusEngine(logger, thinker, auditor)

	// 3. Inicialización del Controlador Fabric (Nokia SR Linux)
	targetRouter := "srl-edge-01.mgmt.local"
	gNMIController := fabric.NewGNMIController(logger, targetRouter)

	// 4. SRE Context con Graceful Shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// 5. Setup de la Tubería de Telemetría (Fan-In Pattern)
	events := make(chan sensor.NetworkEvent, 100) // Canal con buffer
	var wg sync.WaitGroup

	// Instanciamos los sensores
	ebpfSensor := sensor.NewEBPFSensor(logger)
	gnmiSensor := sensor.NewGNMIStreamSensor(logger)

	// Lanzamos los sensores en goroutines independientes
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := ebpfSensor.Start(ctx, events); err != nil && err != context.Canceled {
			logger.Error("eBPF Sensor falló", slog.String("error", err.Error()))
		}
	}()
	go func() {
		defer wg.Done()
		if err := gnmiSensor.Start(ctx, events); err != nil && err != context.Canceled {
			logger.Error("gNMI Sensor falló", slog.String("error", err.Error()))
		}
	}()

	// 6. El Loop Central del Orquestador (AICO Gateway)
	logger.Info("AICO Gateway escuchando eventos de red (eBPF & gNMI)...")

OrchestratorLoop:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Señal de apagado recibida. Cerrando AICO Gateway...")
			break OrchestratorLoop

		case event := <-events:
			logger.WarnContext(ctx, "¡ANOMALÍA DETECTADA!", slog.String("source", event.Source))

			// Convertimos el evento a JSON para dárselo de contexto a la IA
			eventJSON, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				logger.Error("Fallo empaquetando evento para la IA", slog.String("error", err.Error()))
				continue
			}

			// Creamos el Prompt de Contexto para el Motor de Consenso
			networkContext := fmt.Sprintf(`ALERTA DE RED CRÍTICA. 
Se ha recibido el siguiente evento desde la telemetría del router:
%s

Acción requerida: Analiza la anomalía y genera una configuración de mitigación en formato YANG (JSON) nativo para Nokia SR Linux (ej. aplicar Blackhole, Traffic Engineering o FlowSpec).`, string(eventJSON))

			// Las IAs debaten y devuelven el JSON YANG consensuado (Timeboxed)
			aiCtx, aiCancel := context.WithTimeout(ctx, 45*time.Second)
			finalYANG, err := consensusEngine.GenerateAutonomousRemediation(aiCtx, networkContext)
			aiCancel()

			if err != nil {
				logger.ErrorContext(ctx, "Fallo en el consenso de IA", slog.String("error", err.Error()))
				continue // SRE: No apagamos el gateway, seguimos escuchando
			}

			logger.InfoContext(ctx, "Consenso alcanzado con éxito, inyectando a la red...")

			// El controlador aplica el JSON generado autónomamente al router Nokia vía gNMI
			if err := gNMIController.ApplyAutonomousYANG(ctx, []byte(finalYANG)); err != nil {
				logger.ErrorContext(ctx, "Fallo inyectando configuración gNMI", slog.String("error", err.Error()))
			} else {
				logger.InfoContext(ctx, "AICO completó el ciclo autónomo de mitigación. SR Linux protegido.")
			}
		}
	}

	wg.Wait()
	logger.Info("Apagado seguro completado.")
}
