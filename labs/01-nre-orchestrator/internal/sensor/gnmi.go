package sensor

import (
	"context"
	"log/slog"
	"time"
)

// GNMIStreamSensor simula un cliente de Streaming Telemetry (gNMI Subscribe)
// que monitorea el plano de control (BGP, OSPF) del router.
type GNMIStreamSensor struct {
	logger *slog.Logger
}

func NewGNMIStreamSensor(logger *slog.Logger) *GNMIStreamSensor {
	return &GNMIStreamSensor{logger: logger}
}

func (s *GNMIStreamSensor) Start(ctx context.Context, events chan<- NetworkEvent) error {
	s.logger.InfoContext(ctx, "Iniciando gNMI Streaming Sensor (Subscribe request simulado)")

	// Simulamos un evento de caída BGP posterior al ataque volumétrico
	timer := time.NewTimer(8 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		s.logger.InfoContext(ctx, "gNMI Sensor detenido por contexto")
		return ctx.Err()
	case <-timer.C:
		s.logger.WarnContext(ctx, "gNMI Update: BGP Session flapped")
		
		events <- NetworkEvent{
			Source:      "gNMI_ControlPlane_Sensor",
			Severity:    "HIGH",
			Description: "Sesión BGP con AS 65002 perdida (Hold Timer Expired)",
			Metrics: map[string]string{
				"peer":   "10.0.0.12",
				"status": "IDLE",
			},
		}
	}

	return nil
}
