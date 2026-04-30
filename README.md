# Datacenter Reconciler: Automated Network Fabric Provisioning

<div align="center">
  <h3>A Tier-1 ISP Grade Orchestrator for Spine-Leaf Fabrics</h3>
  <p><strong>Go | Site Reliability Engineering (SRE) | Network Reliability Engineering (NRE)</strong></p>
</div>

---

## 📖 Executive Summary

This repository serves as a technical showcase of an Enterprise-grade Network Automation Orchestrator. It is designed to reconcile the desired state of a Datacenter Fabric (stored in a Single Source of Truth like NetBox) with the actual physical/virtual routers (Nokia SR Linux via gNMI).

The project rigorously applies **Cloud-Native patterns, Site Reliability Engineering (SRE) principles, and advanced Go concurrency models** to ensure a reliable, idempotent, and highly performant network provisioning lifecycle.

---

## 🛠️ 1. Software Engineering Architecture (Senior Go)

The codebase strictly adheres to the [Standard Go Project Layout](https://github.com/golang-standards/project-layout), keeping executables in `/cmd` and private business logic safely encapsulated in `/internal`.

### Key Design Patterns
*   **Dependency Injection (DI) & Interface Segregation:** The `netbox.Client` is defined as a strict interface (`internal/netbox/client.go`). The orchestrator never interacts directly with HTTP requests, ensuring clear *Separation of Concerns* and facilitating 100% testable code via mocks.
*   **Fan-Out / Fan-In Concurrency:** Network I/O is slow. Instead of sequential processing, the `engine.Reconciler` implements a dynamic Worker Pool using goroutines and WaitGroups.
*   **Semaphore Bottlenecking:** To prevent DDoS-ing the NetBox API or overwhelming local sockets, concurrency is capped using a buffered channel semaphore (`concurrencySem <- struct{}{}`).
*   **Memory Efficiency (Zero-Copy & Pointers):** Large data structures (like HTTP clients and router states) are passed by reference (`*netbox.DeviceState`). JSON payloads use strictly defined structs (e.g., `IPAMRequest`) instead of dynamic `map[string]interface{}` to prevent excessive Garbage Collector (Heap) allocations.
*   **Slice Capacity Pre-Allocation:** Where possible, slice capacities are known ahead of time and allocated efficiently (`make([]*DeviceState, 0, len(results))`) to prevent expensive re-allocations in memory.

---

## 🛡️ 2. Site Reliability Engineering (SRE) Practices

This tool is built to run unattended in production environments (like Kubernetes or Nomad clusters). It implements strict defensive programming:

*   **Fail-Fast Configuration:** Configuration relies on the 12-Factor App methodology. If critical environment variables (`NETBOX_URL`, `NETBOX_TOKEN`) are missing, the application crashes immediately at boot rather than failing silently hours later.
*   **Graceful Shutdown & Context Propagation:** Uses `signal.NotifyContext` to catch `SIGTERM` / `SIGINT`. All HTTP and gNMI requests receive this `context.Context`. During a pod eviction, ongoing network I/O is gracefully canceled to prevent TCP RSTs and file descriptor leaks.
*   **Explicit Error Handling:** Errors are never ignored (`_`). Every failure is wrapped with `fmt.Errorf("context: %w", err)` to preserve the stack trace and context for operators.
*   **Defensive I/O (Timeouts & FDs):** Custom `http.Transport` tuning prevents goroutine leaks from hanging sockets. Every response body is strictly closed via `defer resp.Body.Close()`.
*   **Idempotency First:** `POST` and `PUT` operations (like IPAM injection or BGP convergence) first evaluate the current state (e.g., HTTP 400 Bad Request on duplicates or gNMI GET checks) before attempting state mutation.
*   **Structured Logging:** Implements Go's native `log/slog` with `slog.NewJSONHandler`. Logs emit machine-readable JSON containing telemetry (hostname, latencies), ready for immediate ingestion by Datadog, Splunk, or Elasticsearch without expensive Grok parsing.

---

## 🌐 3. Network Reliability Engineering (NRE)

The networking logic pushes toward a mathematically deterministic infrastructure, minimizing human input and "drift".

*   **Topology:** Orchestrates a strict Clos Network (Spine-Leaf) defined in `topology.yaml` using containerized Nokia SR Linux instances (`ghcr.io/nokia/srlinux`).
*   **BGP Unnumbered (RFC 5549):** The controller relies on IPv6 Link-Local addressing for BGP peering, removing the need for complex /30 or /31 IPv4 point-to-point subnets.
*   **Deterministic IPAM:** Instead of manually querying a database for "the next available IP", `ipam.allocator` uses mathematics. Spines share `ASN 65000` (to avoid route reflection loops) and loopbacks are dynamically calculated (`10.0.0.11 + SpineID`). Leafs use unique ASNs for eBGP Multi-Path (ECMP).
*   **Zero-Touch Provisioning (ZTP) Readiness:** The Go script registers the mathematically deduced IPAM back into the NetBox SSoT automatically, enforcing that the Source of Truth reflects the mathematical truth.
*   **YANG / gNMI Modeling:** Configuration is abstracted into state trees manipulated via gRPC, removing fragile CLI scraping (Screen-scraping) and moving toward declarative networking.

---

## 🤖 4. AI Consensus Orchestrator (AICO)

A cutting-edge experimental feature that introduces **Autonomous Multi-Cloud AI Remediation**. AICO bridges the gap between Large Language Models (LLMs) and physical network infrastructure using an intelligent "Dual-LLM Verification" architecture.

### Multi-Cloud AI Architecture
*   **Provider Agnostic:** Designed with strict `LLMProvider` interfaces, allowing seamless swapping between major AI clouds (Google Gemini, Anthropic Claude, OpenAI GPT) simply by injecting the correct API keys.
*   **The Consensus Engine:**
    *   **Phase 1 (The Thinker):** One AI model (e.g., Claude) acts as the Senior Architect. It ingests network anomalies (e.g., BGP route drops, interface congestion) and generates a mitigation strategy formatted as native Nokia SR Linux **YANG (JSON)**.
    *   **Phase 2 (The Auditor):** A second, distinct AI model (e.g., Gemini or GPT-4) acts as the strict Site Reliability Engineer. It takes the output from the Thinker and verifies it against YANG schemas, ensuring the payload is syntactically valid and non-destructive before allowing execution.
*   **Autonomous Injection:** Once consensus is reached, the verified YANG JSON is injected directly into the Nokia SR Linux router's backplane via gNMI, fixing the network anomaly in real-time with zero human intervention.

This demonstrates the peak of "Network-as-a-Platform", proving that network elements can be autonomously managed by intelligent, multi-cloud heuristic engines.

---

## 🚀 Getting Started

### Prerequisites
*   Docker & Docker Compose
*   [Containerlab](https://containerlab.dev/) (For spinning up network topologies)
*   Go 1.22+

### 📂 Labs Repository (`/labs`)

This project acts as a centralized repository for various network laboratory topologies:
1.  **`topology.yaml` (Root)**: The baseline Nokia SR Linux Spine-Leaf Fabric.
2.  **`labs/mikrotik-ospf-mesh/`**: A 4-node MikroTik RouterOS v7 mesh topology running OSPF. Designed to demonstrate fast-convergence, OSPF interface templates, and Zero-Touch provisioning via `.rsc` configuration scripts.

To deploy any lab, navigate to its directory and run:
```bash
sudo containerlab deploy -t topology.yaml
```

### Running the Orchestrator (Go)

1.  **Spin up the NetBox SSoT:**

    ```bash
    docker-compose up -d
    ```
2.  **Export Credentials:**
    ```bash
    export NETBOX_URL="http://localhost:8080"
    export NETBOX_TOKEN="your_token_here"
    ```
3.  **Run the Reconciler:**
    ```bash
    go run cmd/reconciler/main.go
    ```

> *"Automatizando redes con el rigor del desarrollo de software."*
