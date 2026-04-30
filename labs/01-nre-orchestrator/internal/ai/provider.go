package ai

import "context"

// LLMProvider define el contrato estricto Multi-Cloud para interactuar con IAs.
// Siguiendo las reglas de Separation of Concerns, la lógica de negocio del consenso
// nunca debe saber si está hablando con GPT, Claude o Gemini.
type LLMProvider interface {
	// Ask envía un prompt estructurado y devuelve la respuesta en crudo de la IA.
	Ask(ctx context.Context, prompt string) (string, error)

	// Name devuelve el identificador amigable del proveedor para observabilidad (Logs SRE).
	Name() string
}
