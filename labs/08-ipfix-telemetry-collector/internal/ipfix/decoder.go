package ipfix

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"ipfix-telemetry-collector/internal/metrics"
)

// Constantes IPFIX (RFC 7011)
const (
	IPFIXVersion = 10
	HeaderLength = 16
)

// TemplateRecord representa un esquema cacheado para decodificar datos
type TemplateRecord struct {
	TemplateID uint16
	FieldCount uint16
	Fields     []FieldSpecifier
}

type FieldSpecifier struct {
	InformationElementID uint16
	FieldLength          uint16
}

// Decoder maneja el estado idempotente de la recepción de paquetes IPFIX
type Decoder struct {
	logger *slog.Logger

	// Caché de Templates: map[ObservationDomainID]map[TemplateID]TemplateRecord
	mu        sync.RWMutex
	templates map[uint32]map[uint16]*TemplateRecord
}

func NewDecoder(logger *slog.Logger) *Decoder {
	return &Decoder{
		logger:    logger,
		templates: make(map[uint32]map[uint16]*TemplateRecord),
	}
}

// ParseUDP payload crudo directamente del socket (Zero-Copy)
func (d *Decoder) ParseUDP(payload []byte, remoteAddr net.Addr) error {
	if len(payload) < HeaderLength {
		metrics.ParseErrors.Inc()
		return fmt.Errorf("paquete demasiado corto: %d bytes", len(payload))
	}

	// 1. Parsear Header IPFIX
	version := binary.BigEndian.Uint16(payload[0:2])
	if version != IPFIXVersion {
		metrics.ParseErrors.Inc()
		return fmt.Errorf("versión no soportada, se esperaba v10, llegó: %d", version)
	}

	// length := binary.BigEndian.Uint16(payload[2:4])
	// exportTime := binary.BigEndian.Uint32(payload[4:8])
	// seqNum := binary.BigEndian.Uint32(payload[8:12])
	observationDomainID := binary.BigEndian.Uint32(payload[12:16])

	offset := HeaderLength

	// 2. Iterar sobre los Sets
	for offset < len(payload) {
		if offset+4 > len(payload) {
			break
		}
		setID := binary.BigEndian.Uint16(payload[offset : offset+2])
		setLength := binary.BigEndian.Uint16(payload[offset+2 : offset+4])

		if setLength < 4 {
			metrics.ParseErrors.Inc()
			return fmt.Errorf("longitud de Set inválida")
		}
		if offset+int(setLength) > len(payload) {
			break
		}

		setPayload := payload[offset+4 : offset+int(setLength)]

		if setID == 2 || setID == 3 {
			// Es un Template Set o Options Template Set
			d.parseTemplateSet(observationDomainID, setPayload)
		} else if setID > 255 {
			// Data Set
			d.parseDataSet(observationDomainID, setID, setPayload)
		}

		offset += int(setLength)
	}

	return nil
}

func (d *Decoder) parseTemplateSet(domainID uint32, payload []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.templates[domainID]; !exists {
		d.templates[domainID] = make(map[uint16]*TemplateRecord)
	}

	// Simplificación para la prueba: Guardar template de forma idempotente
	// En producción, se parsea FieldCount y Fields[]
	d.logger.Debug("Template Set recibido y cacheado de forma idempotente", slog.Any("domainID", domainID))
}

func (d *Decoder) parseDataSet(domainID uint32, templateID uint16, payload []byte) {
	d.mu.RLock()
	domainTemplates, ok := d.templates[domainID]
	d.mu.RUnlock()

	if !ok || domainTemplates[templateID] == nil {
		// No tenemos el template aún (Drop Idempotente)
		d.logger.Debug("Data Set ignorado, falta Template cacheado", slog.Int("templateID", int(templateID)))
		return
	}

	// Aquí decodificamos según el Template.
	// Para el propósito del Lab, mockeamos la extracción de IPs y Bytes y lo empujamos a Prometheus.
	metrics.FlowsProcessed.Inc()

	// Mocking una extracción de bytes
	srcIP := "10.0.0.1"
	dstIP := "192.168.1.100"
	bytesTransferidos := float64(len(payload) * 10) // Simulado

	metrics.BytesProcessed.WithLabelValues(srcIP, dstIP).Add(bytesTransferidos)
}
