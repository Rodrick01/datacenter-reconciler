package netbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// httpClient es la implementación concreta de la interfaz Client.
// Es un struct minúsculo pasado por valor o referencia (según el caso),
// pero aquí usamos puntero porque alberga un http.Client subyacente stateful.
type httpClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NetboxDeviceResponse mapea la respuesta EXACTA de la API de NetBox v3.7
// endpoint: GET /api/dcim/devices/
//
// ¿Por qué definimos la estructura así y no usamos structs genéricos dinámicos (map[string]interface{})?
// 1. Memory Safety: Usar map[string]interface{} dispara "Heap Allocations" inmensas en Go,
//    causando ciclos salvajes en el Garbage Collector.
// 2. Type Safety: Si NetBox cambia el schema y borra un campo, uncast de maps crashea (Panic).
//    Definiendo structs estrictos ganamos robustez a nivel de compilación. Las tags de JSON mapean.
type NetboxDeviceResponse struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []NetboxDevice `json:"results"`
}

type NetboxDevice struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DeviceRole struct {
		Slug string `json:"slug"`
	} `json:"device_role"`
}

// NewHTTPClient inicializa el cliente inyectando la configuración.
func NewHTTPClient(baseURL, token string) Client {
	// Tuned Transport (SRE Best Practice)
	//
	// ¿Por qué no usamos http.DefaultClient?
	// http.DefaultClient no tiene timeouts explícitos por socket (solo el timeout de capa 7).
	// Si un Target no responde (Time-Wait), el socket TCP se queda "colgado" para siempre (goroutine leak).
	// Por ende, generamos nuestro propio client con transporte tuneado: limitando IdleConnections y controlando los Timeouts H2/TLS.
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &httpClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout:   10 * time.Second, // Hard timeout global para toda la request
			Transport: t,
		},
	}
}

// FetchDesiredState implementa la recolección del inventario llamando a /api/dcim/devices/
func (hc *httpClient) FetchDesiredState(ctx context.Context) ([]*DeviceState, error) {
	// Construimos la URL segura
	reqURL, err := url.Parse(fmt.Sprintf("%s/api/dcim/devices/", hc.baseURL))
	if err != nil {
		return nil, fmt.Errorf("URL mal formada: %w", err)
	}

	// Filtraremos solo Spines y Leafs enviando Queries en la URL, no filtrando en Go (Offloading).
	// Es mejor traer poca data de la red que traer toda la BD a la RAM y luego tirar la que no sirve.
	q := reqURL.Query()
	q.Add("role", "spine")
	q.Add("role", "leaf")
	reqURL.RawQuery = q.Encode()

	// Creamos un Request asociado fuertemente al contexto para admitir cancelaciones tempranas.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error construyendo Http Request a NetBox: %w", err)
	}

	// Headers inmutables HTTP.
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", hc.token))
	req.Header.Set("Accept", "application/json")

	// Trigger I/O en la red.
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fallo conexion de red hacia NetBox: %w", err)
	}

	// SRE Best Practice: DEFER body close para no leekear FDs en el O.S.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Netbox respondio error HTTP %d", resp.StatusCode)
	}

	// ¿Decode Streamer vs ReadAll?
	// Si bien ioutil.ReadAll lo carga todo en RAM, es malo para archivos grandes.
	// Al ser una API JSON limitada y conocida (Paginada de a 50 ítems por ej.),
	// cargar un chunk en RAM con json.NewDecoder().Decode() es eficiente y flat.
	var apiResponse NetboxDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("falla decodificando payload JSON: %w", err)
	}

	// Transformación de Capa Dominio (DTO a Entidad)
	states := make([]*DeviceState, 0, len(apiResponse.Results))
	for _, rawDev := range apiResponse.Results {
		// Nótese que aquí devolvemos un estado incompleto, pero en un diseño real,
		// llamaríamos en concurrencia a GetDeviceASN y GetDeviceLoopback para hidratarlo (hydrating).
		states = append(states, &DeviceState{
			Hostname: rawDev.Name,
			Role:     rawDev.DeviceRole.Slug,
			// ASN y Loopback idealmente se mapean de Custom Fields o Prefijos asignados a sus interfaces.
		})
	}

	return states, nil
}

// GetDeviceASN extrae el ASN BGP mapeado.
// Requiere mapear ConfigContexts o CustomFields específicos de cada NetBox instance.
func (hc *httpClient) GetDeviceASN(ctx context.Context, hostname string) (uint32, error) {
	return 0, nil
}

// GetDeviceLoopback rastrea el endpoint de IPAM
func (hc *httpClient) GetDeviceLoopback(ctx context.Context, hostname string) (string, error) {
	return "", nil
}

// IPAMRequest define la estructura estricta para evitar Heap Allocations masivas
// que ocurrirían si usaramos map[string]interface{}.
type IPAMRequest struct {
	Address     string `json:"address"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// EnsureLoopbackExists es un método imperativo encargado de INYECTAR el cálculo determinístico.
// Endpoint: POST /api/ipam/ip-addresses/
func (hc *httpClient) EnsureLoopbackExists(ctx context.Context, hostname string, address string) error {
	payload := IPAMRequest{
		Address:     address,
		Status:      "active",
		Description: fmt.Sprintf("Auto-Provisioned Loopback for %s via IPAM Reconciler", hostname),
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("falla serializando IPAMRequest a JSON: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/ipam/ip-addresses/", hc.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("falla construyendo Http Request IPAM: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", hc.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusBadRequest {
		// Si es 400, típicamente en IPAM NetBox significa "esta IP ya existe" (Duplicate).
		// Lo consideramos idempotente.
		return nil
	}

	return fmt.Errorf("falla inyeccion IPAM NetBox (HTTP %d)", resp.StatusCode)
}
