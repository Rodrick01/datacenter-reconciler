package ai

import "context"

// LLMProvider define el contrato estricto Multi-Cloud para interactuar con IAs.
type LLMProvider interface {
	Ask(ctx context.Context, prompt string) (string, error)
	Name() string
}
