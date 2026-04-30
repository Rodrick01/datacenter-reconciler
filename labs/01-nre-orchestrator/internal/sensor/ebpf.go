package sensor

import (
	"context"
	"log/slog"
	"time"
)

// EBPFSensor simula un programa enganchado al XDP (eXpress Data Path) del kernel Linux
// en el router Nokia SR Linux, leyendo paquetes a velocidad de línea.
type EBPFSensor struct {
	logger *slog.Logger
}

func NewEBPFSensor(logger *slog.Logger) *EBPFSensor {
	return &EBPFSensor{logger: logger}
}

func (s *EBPFSensor) Start(ctx context.Context, events chan<- NetworkEvent) error {
	s.logger.InfoContext(ctx, "Iniciando eBPF Sensor (XDP hook attach simulado)")

	// Simulamos que el sensor de eBPF lee mapas de memoria compartida durante un tiempo
	// y de repente detecta un pico de tráfico SYN.
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		s.logger.InfoContext(ctx, "eBPF Sensor detenido por contexto")
		return ctx.Err()
	case <-timer.C:
		s.logger.WarnContext(ctx, "eBPF Map Alert: Volumetric spike detected")
		
		// Inyectamos el evento al pipeline de telemetría (Fan-In)
		events <- NetworkEvent{
			Source:      "eBPF_XDP_Sensor",
			Severity:    "CRITICAL",
			Description: "Ataque DDoS Volumétrico detectado en la interface ethernet-1/1 (TCP SYN Flood)",
			Metrics: map[string]string{
				"pps_current": "15000000",
				"pps_normal":  "120000",
				"target_ip":   "10.0.0.50",
			},
		}
	}

	return nil
}
