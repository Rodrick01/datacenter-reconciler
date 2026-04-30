package config

import (
	"errors"
	"os"
	"time"
)

// AppConfig agrupa todas las configuraciones globales inyectadas al orquestador.
type AppConfig struct {
	// Entorno NetBox
	NetboxURL   string
	NetboxToken string

	// Ajustes de SRE & Concurrencia
	MaxConcurrentWorkers int
	ExecutionTimeout     time.Duration
}

// LoadConfig lee las variables de entorno basándose en la metodología 12-Factor App.
//
// ¿Por qué no usamos flags o archivos JSON estáticos?
// En entornos Cloud-Native (Nomad/Kubernetes), montar archivos de configuración
// introduce latencia y complejiza los secretos. Las variables de entorno son inyectadas
// de forma segura por el orquestador principal (Nomad Vault integration).
// De este modo, la aplicación es inmutable y 100% portable.
func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		// Defaulting: Valores resilientes por defecto si no son alterados
		MaxConcurrentWorkers: 20, // Paralelismo tope para no degradar NetBox
		ExecutionTimeout:     30 * time.Second,
	}

	url, ok := os.LookupEnv("NETBOX_URL")
	if !ok || url == "" {
		// Fail Fast pattern: Si falta la config crítica, la app muere inmediatamente al inicio
		// en lugar de fallar silenciosamente 10 horas después en tiempo de ejecución.
		return nil, errors.New("missing vital environment variable: NETBOX_URL")
	}
	cfg.NetboxURL = url

	token, ok := os.LookupEnv("NETBOX_TOKEN")
	if !ok || token == "" {
		return nil, errors.New("missing vital environment variable: NETBOX_TOKEN")
	}
	cfg.NetboxToken = token

	return cfg, nil
}
