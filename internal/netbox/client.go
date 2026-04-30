package netbox

import (
	"context"
)

// DeviceState resume de manera estructurada el estado deseado de un nodo en la red.
// Se recomienda manejar instancias pequeñas de esta estructura por valor o referencias
// en batch, utilizando punteros para evitar overhead de copia.
type DeviceState struct {
	Hostname string
	Role     string // spine o leaf
	ASN      uint32
	Loopback string // Formato CIDR, ej. "10.0.0.11/32"
}

// Client define el contrato estricto de la interfaz con la API de NetBox.
// Garantiza Separation of Concerns: nuestro reconciliador de red no conoce de REST ni HTTP,
// solo exige este contrato de lectura para operar y comparar estados.
type Client interface {
	// FetchDesiredState extrae eficientemente desde NetBox el estado objetivo
	// de todos los dispositivos relevantes para nuestra infraestructura automatizada.
	// Recibe context.Context requerido para timeouts de red tempranos y gracefully cancelations.
	FetchDesiredState(ctx context.Context) ([]*DeviceState, error)

	// GetDeviceASN extrae específicamente el BGP ASN (ej. 65000 para spines, 651xx para leafs)
	// de un nodo en particular.
	GetDeviceASN(ctx context.Context, hostname string) (uint32, error)

	// GetDeviceLoopback provee la extracción exacta y predecible de la IP Loopback IPv4
	// para el router id y orquestación BGP Unnumbered de RFC 5549.
	GetDeviceLoopback(ctx context.Context, hostname string) (string, error)

	// EnsureLoopbackExists es un método imperativo encargado de registrar la
	// IPAM matemática en el SSoT de manera idempotente.
	EnsureLoopbackExists(ctx context.Context, hostname string, address string) error
}
