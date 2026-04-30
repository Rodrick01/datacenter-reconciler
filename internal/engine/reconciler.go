package engine

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"datacenter-reconciler/internal/fabric"
	"datacenter-reconciler/internal/ipam"
	"datacenter-reconciler/internal/netbox"
)

// Reconciler es el orquestador maestro (Control Plane).
type Reconciler struct {
	logger        *slog.Logger
	netboxClient  netbox.Client
	maxWorkers    int
}

// NewReconciler inyecta las dependencias necesarias.
func NewReconciler(logger *slog.Logger, nb netbox.Client, maxConcurrent int) *Reconciler {
	return &Reconciler{
		logger:       logger,
		netboxClient: nb,
		maxWorkers:   maxConcurrent,
	}
}

// Run comienza el flujo completo del orquestador adoptando un patron Fan-Out Fan-In.
//
// Patrón de Concurrencia "Fan-Out / Fan-In" (Worker Pool Dinámico):
// ¿Por qué lanzamos goroutines concurrentes? 
// Al interactuar con el I/O de red de 200 switches L3 por gRPC, procesarlos secuencialmente
// tomaría demasiados minutos. Goroutines proveen Thread-Safety y escalabilidad I/O.
// Limitamos el nivel máximo de workers pasándolos de a lotes o usando un semaforo limitante (chan struct{}).
func (r *Reconciler) Run(ctx context.Context) error {
	devices, err := r.netboxClient.FetchDesiredState(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo cargar la base de inventario desde NetBox: %w", err)
	}

	if len(devices) == 0 {
		r.logger.InfoContext(ctx, "No hay switches para reconciliar. Fin.")
		return nil
	}

	r.logger.InfoContext(ctx, "Inicializando provisionamiento concurrente", slog.Int("device_count", len(devices)))

	// Utilizamos un WaitGroup para prevenir terminar la ejecución principal (main())
	// antes de que todos los workers (hilos/goroutines) acaben sus transacciones gNMI.
	var wg sync.WaitGroup
	// Semaphore pattern limitador: Un canal bufferizado al máximo de capacidad (ej: 20 slots).
	concurrencySem := make(chan struct{}, r.maxWorkers)

	for _, dev := range devices {
		wg.Add(1) // Registramos el trabajo
		
		// Inyección de goroutine anónima atada a la iteración dinámica.
		// ¡IMPORTANTE!: En arquitecturas concurrentes, siempre aislar las variables de loop pasándolas como valor explícito `d`
		// (Nota: Go 1.22 soluciona la semántica del cierre de loops de forma nativa, pero es cultura SRE dejarlo explícito).
		go func(d *netbox.DeviceState) {
			defer wg.Done()

			// Adquirir un slot del semáforo. Bloquea si ya hay demasiados en ejecución paralela.
			concurrencySem <- struct{}{}
			// Liberar el slot al finalizar (usando defer).
			defer func() { <-concurrencySem }()

			// Timebox por switch: Para evitar que un switch colgado se apodere de un worker perpetuamente.
			nodeCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
			defer cancel()

			r.processSingleNode(nodeCtx, d)

		}(dev)
	}

	// Fan-In: Esperamos aquí congelados hasta que todos los workers reporten Done().
	wg.Wait()
	r.logger.InfoContext(ctx, "Reconciliación de todo el datacenter completada.")
	
	return nil
}

// processSingleNode engloba la lógica IPAM + SSoT + Switch I/O.
func (r *Reconciler) processSingleNode(ctx context.Context, state *netbox.DeviceState) {
	log := r.logger.With(slog.String("hostname", state.Hostname))

	// 1. IPAM Matemático: Calculamos IPs offline para este hostname particular.
	attrs, err := ipam.DeterministicAllocate(state.Hostname)
	if err != nil {
		log.ErrorContext(ctx, "Imposible deducir matemática IPAM", slog.String("error", err.Error()))
		return
	}

	// Enriquecemos el state pointer con la matemática
	state.ASN = attrs.ASN
	state.Loopback = attrs.Loopback

	// 2. Registramos IPAM en NetBox asegurando Consistencia.
	if err := r.netboxClient.EnsureLoopbackExists(ctx, state.Hostname, state.Loopback); err != nil {
		log.WarnContext(ctx, "No pudimos asentar la IPAM en NetBox, pero procedemos a intentar impactar switch", slog.String("error", err.Error()))
	} else {
		log.DebugContext(ctx, "IPAM registrado satisfactoriamente en NetBox SSoT")
	}

	// 3. Impacto Físico vía el protocolo gNMI
	ctrl := fabric.NewGNMIController(r.logger, state.Hostname)
	if err := ctrl.ReconcileNode(ctx, state); err != nil {
		log.ErrorContext(ctx, "El router falla en convergir", slog.String("error", err.Error()))
	} else {
		log.InfoContext(ctx, "Router aprovisionado exitosamente en estado final P0")
	}
}
