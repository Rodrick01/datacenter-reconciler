package topology

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"text/template"

	"github.com/user/ixp-lab/internal/models"
)

// Generator encapsula la lógica para crear la topología final de red.
type Generator struct {
	logger       *slog.Logger
	templatePath string
	outputPath   string
}

// NewGenerator inicializa el generador de topología de Containerlab con SRE logging.
func NewGenerator(logger *slog.Logger, tmplPath, outPath string) *Generator {
	return &Generator{
		logger:       logger,
		templatePath: tmplPath,
		outputPath:   outPath,
	}
}

// Render procesa la plantilla YAML utilizando el modelo de red en memoria.
// Exige un contexto para manejo de timeout y requiere que el Topology se pase por puntero
// para evitar copias costosas en redes extensas.
func (g *Generator) Render(ctx context.Context, topo *models.Topology) error {
	// Verificar si el contexto ha sido cancelado antes de realizar I/O
	select {
	case <-ctx.Done():
		return fmt.Errorf("renderizado abortado por timeout o cancelación: %w", ctx.Err())
	default:
	}

	g.logger.Info("Iniciando renderizado de topología",
		slog.String("network_name", topo.Name),
		slog.Int("nodes_count", len(topo.Nodes)),
		slog.Int("links_count", len(topo.Links)),
	)

	// 1. Cargar la plantilla base
	tmpl, err := template.ParseFiles(g.templatePath)
	if err != nil {
		return fmt.Errorf("fallo al parsear la plantilla %s: %w", g.templatePath, err)
	}

	// 2. Crear el archivo destino YAML de Containerlab
	f, err := os.Create(g.outputPath)
	if err != nil {
		return fmt.Errorf("fallo al crear archivo de salida %s: %w", g.outputPath, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			g.logger.Error("Fallo al cerrar descriptor de archivo", slog.String("error", cerr.Error()))
		}
	}()

	// 3. Renderizar y aplicar los datos a la plantilla
	err = tmpl.Execute(f, topo) // topo is already a pointer
	if err != nil {
		return fmt.Errorf("fallo durante la ejecución de la plantilla: %w", err)
	}

	g.logger.Info("Topología renderizada exitosamente",
		slog.String("output_file", g.outputPath),
	)

	return nil
}
