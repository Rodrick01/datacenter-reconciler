package fabric

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	
	"datacenter-reconciler/internal/netbox"
)

var (
	ErrStateUnchanged = errors.New("estado gNMI ya covergido (idempotency ok)")
)

// GNMIController abstrae la capa de hardware (Nokia SR Linux).
type GNMIController struct {
	logger *slog.Logger
	target string
	// gnmiClient *gNMI.Client // Aqui iría el puntero a la conexión gRPC real
}

func NewGNMIController(logger *slog.Logger, target string) *GNMIController {
	return &GNMIController{
		logger: logger,
		target: target,
	}
}

// ReconcileNode opera una mutación atómica en el State Tree del router.
func (c *GNMIController) ReconcileNode(ctx context.Context, state *netbox.DeviceState) error {
	// Logger estructurado: facilita agregación en ElasticSearch/Datadog.
	logEntry := c.logger.With(
		slog.String("hostname", state.Hostname),
		slog.String("role", state.Role),
	)

	logEntry.InfoContext(ctx, "Comienza reconciliacion BGP vía gNMI")

	// FASE 1: Extract (Idempotency check preventivo)
	// No disparamos un gNMI Set ciego. Cuesta ciclos de CPU del router recompilar rutas.
	currentASN, err := c.fetchCurrentASN(ctx)
	if err != nil {
		return fmt.Errorf("gNMI Get fallido [fase de extraccion]: %w", err)
	}

	if currentASN == state.ASN {
		logEntry.InfoContext(ctx, "Ignorando SetRequest: El nodo ya posee la configuracion BGP Unnumbered correcta")
		return nil
	}

	// FASE 2: Load (Mutación)
	logEntry.InfoContext(ctx, "Inyectando modelo YANG bgp-unnumbered", slog.Uint64("asn", uint64(state.ASN)))
	if err := c.applyBGP(ctx, state.ASN, state.Loopback); err != nil {
		return fmt.Errorf("gNMI Set fallo [fase mutacion]: %w", err)
	}

	return nil
}

// fetchCurrentASN recupera la información via path: "urn:srl_nokia/bgp/bgp-instance[name=default]/asn"
func (c *GNMIController) fetchCurrentASN(ctx context.Context) (uint32, error) {
	// gNMI Get implementation map
	return 0, nil
}

// applyBGP construye un SetRequest Update/Replace con un payload JSON IETF (YANG)
func (c *GNMIController) applyBGP(ctx context.Context, asn uint32, loopback string) error {
	// Estructura del Payload:
	// {"bgp": {"autonomous-system": 65000, "router-id": "10.0.0.11", ...}}
	// jsonYANG, _ := json.Marshal(configStruct)
	// client.Set(ctx, &gnmi.SetRequest{Update: []gnmi.Update{...}})
	return nil
}
