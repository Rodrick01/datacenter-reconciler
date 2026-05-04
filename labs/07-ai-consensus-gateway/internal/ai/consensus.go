package ai

import (
	"context"
	"fmt"
	"log/slog"
)

// ConsensusEngine maneja el debate entre dos LLMs distintos para tomar decisiones seguras.
type ConsensusEngine struct {
	logger  *slog.Logger
	thinker LLMProvider // La IA primaria que diseña la solución
	auditor LLMProvider // La IA estricta que audita el resultado
}

func NewConsensusEngine(logger *slog.Logger, thinker, auditor LLMProvider) *ConsensusEngine {
	return &ConsensusEngine{
		logger:  logger,
		thinker: thinker,
		auditor: auditor,
	}
}

// GenerateConsensus ejecuta la arquitectura "Dual LLM Verification" de manera genérica.
func (ce *ConsensusEngine) GenerateConsensus(ctx context.Context, contextData, objective string) (string, error) {
	ce.logger.InfoContext(ctx, "Iniciando consenso AI Multi-Cloud", 
		slog.String("thinker", ce.thinker.Name()), 
		slog.String("auditor", ce.auditor.Name()),
	)

	// FASE 1: Creación (Thinker)
	promptThinker := fmt.Sprintf(`Analiza el siguiente contexto y resuelve el objetivo solicitado. Devuelve ÚNICAMENTE la respuesta o el código final, sin explicaciones ni formato markdown a menos que se solicite en el objetivo.
Contexto:
%s

Objetivo:
%s`, contextData, objective)

	proposedSolution, err := ce.thinker.Ask(ctx, promptThinker)
	if err != nil {
		return "", fmt.Errorf("la IA primaria [%s] falló al pensar: %w", ce.thinker.Name(), err)
	}

	ce.logger.DebugContext(ctx, "Propuesta generada por IA primaria", slog.Int("len", len(proposedSolution)))

	// FASE 2: Auditoría (Auditor)
	promptAuditor := fmt.Sprintf(`Actúa como un auditor experto y estricto. Revisa la siguiente propuesta para cumplir el objetivo dado en base al contexto.
Si la propuesta tiene errores, alucinaciones o rompe buenas prácticas, corrígela. 
Devuelve ÚNICAMENTE la respuesta o el código corregido, sin explicaciones ni markdown, a menos que el objetivo lo exija.

Contexto original:
%s

Objetivo original:
%s

Propuesta de la IA anterior a auditar:
%s`, contextData, objective, proposedSolution)

	verifiedSolution, err := ce.auditor.Ask(ctx, promptAuditor)
	if err != nil {
		return "", fmt.Errorf("la IA auditora [%s] falló al validar: %w", ce.auditor.Name(), err)
	}

	ce.logger.InfoContext(ctx, "Consenso alcanzado exitosamente")

	return verifiedSolution, nil
}
