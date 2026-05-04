package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/user/ixp-lab/internal/models"
	"github.com/user/ixp-lab/internal/topology"
)

func main() {
	// 1. Configurar Logger Estructurado (SRE standard para ISP Tier 1)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Propagación de Contexto (Context Propagation) con Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Definición de la Topología Modernizada (pasando pesados structs por puntero)
	topo := &models.Topology{
		Name:       "ixp-lab",
		MgmtSubnet: "10.254.0.0/24",
		Nodes: []models.Node{
			{Name: "r1", Kind: "vr-vmx", Image: "vrnetlab/vr-vmx:latest", MgmtIPv4: "10.254.0.11"},
			{Name: "r3", Kind: "vr-vmx", Image: "vrnetlab/vr-vmx:latest", MgmtIPv4: "10.254.0.13"},
			{Name: "r4", Kind: "vr-vmx", Image: "vrnetlab/vr-vmx:latest", MgmtIPv4: "10.254.0.14"},
			// Nodos GoBGP usando contenedores minimalistas (Zero-copy network configs)
			{
				Name:     "g2",
				Kind:     "linux",
				Image:    "osrg/gobgp:latest",
				MgmtIPv4: "10.254.0.102",
				Binds:    []string{"./configs/g2_gobgpd.toml:/etc/gobgp/gobgpd.toml:ro"},
				Exec: []string{
					"ip addr add 10.1.12.2/24 dev eth1",
					"ip addr add 10.173.176.2/24 dev eth2",
					"gobgpd -f /etc/gobgp/gobgpd.toml --pprof-disable --api-hosts=':50051' -p 2112 -d",
				},
			},
			{
				Name:     "rs",
				Kind:     "linux",
				Image:    "osrg/gobgp:latest",
				MgmtIPv4: "10.254.0.150",
				Binds:    []string{"./configs/rs_gobgpd.toml:/etc/gobgp/gobgpd.toml:ro"},
				Exec: []string{
					"ip addr add 10.173.176.254/24 dev eth1",
					"gobgpd -f /etc/gobgp/gobgpd.toml --pprof-disable --api-hosts=':50051' -p 2112 -d",
				},
			},
			// Observability Stack
			{Name: "util", Kind: "linux", Image: "alpine:3.18", MgmtIPv4: "10.254.0.250"},
			{
				Name:     "prometheus",
				Kind:     "linux",
				Image:    "prom/prometheus:latest",
				MgmtIPv4: "10.254.0.201",
				Ports:    []string{"9090:9090"},
				Binds:    []string{"./configs/prometheus.yml:/etc/prometheus/prometheus.yml:ro"},
			},
			{
				Name:     "grafana",
				Kind:     "linux",
				Image:    "grafana/grafana:latest",
				MgmtIPv4: "10.254.0.202",
				Ports:    []string{"3000:3000"},
				Binds: []string{
					"./configs/grafana/provisioning:/etc/grafana/provisioning",
					"./configs/grafana/dashboards:/var/lib/grafana/dashboards",
				},
			},
			// Switch bridge L2
			{Name: "ix-switch", Kind: "bridge", Image: "", MgmtIPv4: ""},
		},
		Links: []models.Link{
			{Endpoints: []string{"r1:eth1", "r4:eth1"}},
			{Endpoints: []string{"r1:eth2", "g2:eth1"}},
			{Endpoints: []string{"g2:eth2", "ix-switch:eth1"}},
			{Endpoints: []string{"rs:eth1", "ix-switch:eth2"}},
			{Endpoints: []string{"r3:eth1", "ix-switch:eth3"}},
			{Endpoints: []string{"r1:eth3", "ix-switch:eth4"}},
		},
	}

	// 4. Inicializar y ejecutar Generador
	gen := topology.NewGenerator(logger, "topology.tmpl", "ixp-lab.clab.yml")

	if err := gen.Render(ctx, topo); err != nil {
		// Manejo de Error Explícito (Fail-Safe), terminación controlada
		logger.Error("Fallo crítico durante el renderizado de la topología", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Despliegue preparado para orquestación via Containerlab. Ejecute: `clab deploy -t ixp-lab.clab.yml`")
}
