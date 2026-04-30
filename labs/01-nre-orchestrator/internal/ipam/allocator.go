package ipam

import (
	"fmt"
	"regexp"
	"strconv"
)

// Constantes maestras para la red
const (
	SpineBaseASN     uint32 = 65000
	LeafBaseASN      uint32 = 65100
	LoopbackBaseIPv4 string = "10.0.0."
	SpineOffset      int    = 10  // Spines empiezan en 10.0.0.11
	LeafOffset       int    = 100 // Leafs empiezan en 10.0.0.101
)

var nameRegex = regexp.MustCompile(`^(spine|leaf)(\d+)$`)

// NetworkAttributes contiene los parámetros auto-calculados de un nodo.
type NetworkAttributes struct {
	Role     string
	ID       int
	ASN      uint32
	Loopback string
}

// DeterministicAllocate calcula el ASN y Loopback matemáticamente basado en el hostname.
// Esto elimina la necesidad de asignar IPs manualmente; el sistema se "calcula a sí mismo".
func DeterministicAllocate(hostname string) (*NetworkAttributes, error) {
	matches := nameRegex.FindStringSubmatch(hostname)
	if len(matches) != 3 {
		return nil, fmt.Errorf("hostname %s no cumple el estándar (ej. spine1, leaf2)", hostname)
	}

	role := matches[1]
	idStr := matches[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("error parseando ID de nodo: %w", err)
	}

	attr := &NetworkAttributes{
		Role: role,
		ID:   id,
	}

	switch role {
	case "spine":
		// Todos los Spines BGP en esta arquitectura particular comparten el mismo ASN para no reflejar rutas entre ellos
		attr.ASN = SpineBaseASN
		attr.Loopback = fmt.Sprintf("%s%d/32", LoopbackBaseIPv4, SpineOffset+id)
	case "leaf":
		// Los Leafs necesitan un ASN único para eBGP multi-path (ECMP)
		attr.ASN = LeafBaseASN + uint32(id)
		attr.Loopback = fmt.Sprintf("%s%d/32", LoopbackBaseIPv4, LeafOffset+id)
	}

	return attr, nil
}
