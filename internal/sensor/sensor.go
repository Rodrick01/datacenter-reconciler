package sensor

import (
	"context"
)

// NetworkEvent representa una anomalía estructurada detectada en el plano de datos o control.
type NetworkEvent struct {
	Source      string            `json:"source"`
	Severity    string            `json:"severity"`
	Description string            `json:"description"`
	Metrics     map[string]string `json:"metrics"`
}

// Sensor define el contrato estricto para cualquier fuente de telemetría.
// Siguiendo las reglas de concurrencia en Go, el sensor no bloquea;
// inyecta los eventos en el canal que se le pasa por parámetro.
type Sensor interface {
	Start(ctx context.Context, events chan<- NetworkEvent) error
}
