package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// MikrotikRoute define los campos relevantes para el patching
type MikrotikRoute struct {
	Id       string `json:".id"`
	DstCheck string `json:"dst-address"`
	Gateway  string `json:"gateway"`
	Distance int    `json:"distance"`
	Active   bool   `json:"active"`
}

// ConsensusRequest es el payload que espera el AI Gateway
type ConsensusRequest struct {
	Context   string `json:"context"`
	Objective string `json:"objective"`
	Thinker   string `json:"thinker"`
	Auditor   string `json:"auditor"`
}

// ConsensusResponse es el output esperado del AI Gateway
type ConsensusResponse struct {
	ConsensusResult string `json:"consensus_result"`
}

func main() {
	// Configuraciones desde variables de entorno con fallbacks
	routerIP := os.Getenv("MIKROTIK_IP")
	if routerIP == "" {
		routerIP = "192.168.88.1"
	}
	mikrotikUser := os.Getenv("MIKROTIK_USER")
	if mikrotikUser == "" {
		mikrotikUser = "ai_agent"
	}
	mikrotikPass := os.Getenv("MIKROTIK_PASS")
	if mikrotikPass == "" {
		mikrotikPass = "modo_cientifico"
	}
	gatewayURL := os.Getenv("GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://192.168.88.254:8080/api/v1/consensus"
	}

	// Cliente HTTP reutilizable con timeout global
	client := &http.Client{Timeout: 10 * time.Second}

	log.Println("Iniciando MikroTik AI Failover Agent...")

	// PASO A: Obtener la tabla de ruteo
	routesURL := fmt.Sprintf("http://%s/rest/ip/route", routerIP)
	req, err := http.NewRequest(http.MethodGet, routesURL, nil)
	if err != nil {
		log.Fatalf("Error creando request para MikroTik: %v", err)
	}
	req.SetBasicAuth(mikrotikUser, mikrotikPass)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error leyendo rutas del MikroTik: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Código de error HTTP de MikroTik: %d", resp.StatusCode)
	}

	bodyRoutes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error leyendo body de rutas: %v", err)
	}

	log.Printf("Tabla de ruteo obtenida (%d bytes)", len(bodyRoutes))

	// PASO B: Armar el prompt para el AI Gateway
	aiReq := ConsensusRequest{
		Context: string(bodyRoutes),
		Objective: `Analizar la tabla de ruteo. El ISP1 tiene latencia. Devolver UNICAMENTE un JSON válido para la REST API de MikroTik (método PATCH) que cambie la 'distance' a 20.`,
		Thinker: "gemini",
		Auditor: "claude",
	}

	aiPayload, err := json.Marshal(aiReq)
	if err != nil {
		log.Fatalf("Error armando payload de IA: %v", err)
	}

	// PASO C: POST al AI Gateway con timeout extendido de 120s
	log.Printf("Consultando al AI Gateway (%s)...", gatewayURL)
	aiClient := &http.Client{Timeout: 120 * time.Second}
	gwReq, err := http.NewRequest(http.MethodPost, gatewayURL, bytes.NewBuffer(aiPayload))
	if err != nil {
		log.Fatalf("Error creando request para AI Gateway: %v", err)
	}
	gwReq.Header.Set("Content-Type", "application/json")

	gwResp, err := aiClient.Do(gwReq)
	if err != nil {
		log.Fatalf("Error contactando AI Gateway: %v", err)
	}
	defer gwResp.Body.Close()

	if gwResp.StatusCode != http.StatusOK {
		log.Fatalf("AI Gateway retornó error HTTP: %d", gwResp.StatusCode)
	}

	var aiResult ConsensusResponse
	if err := json.NewDecoder(gwResp.Body).Decode(&aiResult); err != nil {
		log.Fatalf("Error decodificando respuesta de la IA: %v", err)
	}

	log.Printf("Consenso de IA alcanzado: %s", aiResult.ConsensusResult)

	// PASO C.5: Sanity Check estricto antes de tocar el router
	var remediation map[string]interface{}
	if err := json.Unmarshal([]byte(aiResult.ConsensusResult), &remediation); err != nil {
		log.Fatalf("❌ CRÍTICO: La IA no devolvió un JSON válido. Abortando remediación: %v", err)
	}

	// MikroTik espera que 'distance' sea un número entre 1 y 255
	if distanceFloat, ok := remediation["distance"].(float64); !ok || distanceFloat < 1 || distanceFloat > 255 {
		log.Fatalf("❌ CRÍTICO: La IA devolvió una métrica fuera de rango o de tipo inválido. Abortando.")
	}

	log.Println("✅ Sanity check superado. Payload seguro para inyección.")

	// PASO D: Inyectar parche dinámico (simulando que la IA detectó ID "*1")
	// En un caso real completo, la IA o el agente debe inferir este ID dinámicamente.
	routeID := "*1"
	patchURL := fmt.Sprintf("http://%s/rest/ip/route/%s", routerIP, routeID)
	
	patchReq, err := http.NewRequest(http.MethodPatch, patchURL, bytes.NewBufferString(aiResult.ConsensusResult))
	if err != nil {
		log.Fatalf("Error creando request de parche: %v", err)
	}
	patchReq.SetBasicAuth(mikrotikUser, mikrotikPass)
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := client.Do(patchReq)
	if err != nil {
		log.Fatalf("Error ejecutando PATCH en MikroTik: %v", err)
	}
	defer patchResp.Body.Close()

	if patchResp.StatusCode >= 200 && patchResp.StatusCode < 300 {
		log.Println("✅ Self-Healing ejecutado con éxito. Router reconfigurado.")
	} else {
		log.Fatalf("❌ Fallo la ejecución en el router. HTTP Status: %d", patchResp.StatusCode)
	}
}
