package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"ipfix-telemetry-collector/internal/ipfix"
	"ipfix-telemetry-collector/internal/metrics"
)

func main() {
	// 1. SRE Observability: Structured JSON Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("Iniciando IPFIX Telemetry Collector")

	// 2. Context Propagation & Graceful Shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	// 3. Levantar Endpoint de Prometheus
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    ":2112",
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Prometheus metrics expuestas en :2112/metrics")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error en servidor HTTP Prometheus", slog.String("error", err.Error()))
		}
	}()

	// 4. Iniciar UDP Listener de Alta Performance (Zero-Copy focus)
	addr, err := net.ResolveUDPAddr("udp", ":4739") // Puerto estándar IPFIX
	if err != nil {
		logger.Error("Fallo resolviendo dirección UDP", slog.String("error", err.Error()))
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Error("Fallo abriendo puerto UDP", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Optimizaciones de SRE: Aumentar buffer del socket
	if err := conn.SetReadBuffer(10485760); err != nil { // 10MB Buffer
		logger.Warn("No se pudo configurar ReadBuffer UDP", slog.String("error", err.Error()))
	}

	decoder := ipfix.NewDecoder(logger)
	logger.Info("Escuchando flujos IPFIX en UDP :4739")

	// Pre-alocar buffer para evitar re-alocaciones en el loop
	// Tamaño estándar MTU ~1500, usamos 2048
	buf := make([]byte, 2048)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer conn.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Usamos un timeout corto en ReadFromUDP para no bloquear infinitamente y poder chequear <-ctx.Done()
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				n, remoteAddr, err := conn.ReadFromUDP(buf)

				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue // Timeout esperado
					}
					logger.Error("Error leyendo UDP", slog.String("error", err.Error()))
					continue
				}

				metrics.PacketsReceived.Inc()

				// Procesamiento Zero-Copy: pasamos el slice crudo
				if err := decoder.ParseUDP(buf[:n], remoteAddr); err != nil {
					logger.Debug("Fallo parseando IPFIX", slog.String("error", err.Error()))
				}
			}
		}
	}()

	// 5. Esperar Señal de Cierre
	<-ctx.Done()
	logger.Info("Recibida señal de apagado. Ejecutando graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Fallo al apagar servidor HTTP", slog.String("error", err.Error()))
	}

	wg.Wait()
	logger.Info("IPFIX Telemetry Collector apagado correctamente.")
}
