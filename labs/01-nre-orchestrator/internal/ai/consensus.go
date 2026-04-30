package ai

import (
	"context"
	"fmt"
	"log/slog"
)

// ConsensusEngine maneja el debate entre dos LLMs distintos para tomar decisiones de red seguras.
type ConsensusEngine struct {
	logger  *slog.Logger
	thinker LLMProvider // La IA primaria que diseña la solución (ej. Claude)
	auditor LLMProvider // La IA estricta que audita el código generado (ej. Gemini/GPT)
}

func NewConsensusEngine(logger *slog.Logger, thinker, auditor LLMProvider) *ConsensusEngine {
	return &ConsensusEngine{
		logger:  logger,
		thinker: thinker,
		auditor: auditor,
	}
}

// GenerateAutonomousRemediation ejecuta la arquitectura "Dual LLM Verification".
// Pide una solución, la verifica con la otra IA y devuelve el JSON YANG final.
func (ce *ConsensusEngine) GenerateAutonomousRemediation(ctx context.Context, networkContext string) (string, error) {
	ce.logger.InfoContext(ctx, "Iniciando consenso AI Multi-Cloud", 
		slog.String("thinker", ce.thinker.Name()), 
		slog.String("auditor", ce.auditor.Name()),
	)

	// FASE 1: Creación (Thinker)
	promptThinker := fmt.Sprintf(`Eres un Senior Network Architect. Analiza este estado de la red y genera una configuración de mitigación en formato JSON YANG nativo para Nokia SR Linux.
Contexto: %s
Solo devuelve el JSON puro, sin markdown ni explicaciones.`, networkContext)

	proposedYANG, err := ce.thinker.Ask(ctx, promptThinker)
	if err != nil {
		return "", fmt.Errorf("la IA primaria [%s] falló al pensar: %w", ce.thinker.Name(), err)
	}
	
	ce.logger.DebugContext(ctx, "Propuesta generada por IA primaria", slog.String("payload_len", fmt.Sprintf("%d bytes", len(proposedYANG))))

	// FASE 2: Auditoría (Auditor)
	promptAuditor := fmt.Sprintf(`Eres un SRE estricto. La siguiente es una propuesta YANG (JSON) para Nokia SR Linux.
Verifica que sea sintácticamente válida, que no destruya el BGP actual y que sea segura de aplicar vía gNMI.
Si encuentras un error, corrígelo. Devuelve ÚNICAMENTE el JSON final válido, nada de texto.

Propuesta Original:
%s`, proposedYANG)

	verifiedYANG, err := ce.auditor.Ask(ctx, promptAuditor)
	if err != nil {
		return "", fmt.Errorf("la IA auditora [%s] falló al validar: %w", ce.auditor.Name(), err)
	}

	ce.logger.InfoContext(ctx, "Consenso alcanzado exitosamente")

	return verifiedYANG, nil
}
