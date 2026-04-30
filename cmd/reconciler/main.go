package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"datacenter-reconciler/internal/config"
	"datacenter-reconciler/internal/engine"
	"datacenter-reconciler/internal/netbox"
)

func main() {
	// Logger SRE Standard
	// Usamos formato JSON (slog.NewJSONHandler) nativo porque herramientas como
	// Elasticsearch, Splunk o Datadog lo ingieren nativamente sin parseo grok pesado.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Contexto de Interrupción del SO (Graceful Shutdown)
	// Si hacemos un kill -9 (SIGKILL) las conexiones TCP activas de Netbox quedan rotas (RST).
	// Capturamos SIGTERM/SIGINT, cancelamos el context padre, lo que desencadena que
	// todos los request HTTP/gRPC se aborten limpiamente con context.Canceled.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.InfoContext(ctx, "Iniciando Node Reconciler Tier-1", slog.String("version", "v1.0.0.enterprise"))

	// 1. Carga de Variables de Entorno Segura (Fail-Fast)
	// Fallará si falta NETBOX_URL y detendrá el arranque, evitando comportamientos inconsistentes.
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.ErrorContext(ctx, "Fallo fatal de configuracion inicial", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 2. Dependency Injection Invertida
	// Instanciamos abstracciones desde afuera para facilitar TestMocks.
	nbClient := netbox.NewHTTPClient(cfg.NetboxURL, cfg.NetboxToken)
	
	orchEngine := engine.NewReconciler(logger, nbClient, cfg.MaxConcurrentWorkers)

	// 3. Kickstart Motor Concurrente
	if err := orchEngine.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "Motor de orquestacion fallo incomprensiblemente", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.InfoContext(ctx, "Apagado en frio finalizado con exito")
}
