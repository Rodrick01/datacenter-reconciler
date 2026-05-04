# AI Consensus Orchestrator (AICO) & Autonomous Nokia SR Linux Fabric

This laboratory represents the pinnacle of autonomous Network Reliability Engineering (NRE). It implements a closed-loop, self-healing **Leaf-Spine** datacenter fabric utilizing white-box routers running **Nokia SR Linux**.

The core of the project is the **AICO Gateway (AI Consensus Orchestrator)**, a Go-based engine that ingest ultra-low latency telemetry via **eBPF** and **gNMI**, detects anomalies, and leverages a multi-LLM consensus loop (Claude + Gemini) to debate, verify, and automatically inject remediation configurations in native YANG format.

---

## 🌟 Enterprise-Grade Technical Features

### 1. 🏗️ Autonomous Leaf-Spine Fabric
- **Nokia SR Linux White-Boxes:** The datacenter topology is fully containerized using Containerlab, running native `ghcr.io/nokia/srlinux` images arranged in a strict Clos (Leaf-Spine) architecture.
- **eBPF & gNMI Sensor Pipelines:** Instead of legacy SNMP, the orchestrator utilizes concurrent Go routines to establish gNMI telemetry streams and eBPF kernel-level sensors, detecting network anomalies at microsecond precision.

### 2. 🧠 Multi-LLM Consensus Engine (Thinker/Auditor Loop)
- **AI-Driven Remediation:** When a critical network anomaly is detected (e.g., a massive route leak or DDoS vector), the event context is passed to the AI gateway.
- **Debate & Verification:** A primary LLM (Thinker) proposes a remediation strategy in native Nokia YANG format. A secondary LLM (Auditor) verifies the proposal against safety constraints to prevent catastrophic network loops or isolation.
- **Zero-Touch Injection:** Once consensus is reached, the Go controller safely applies the autonomous YANG payload directly to the Nokia routers via gNMI, neutralizing the threat dynamically.

### 3. ⚡ SRE Paradigms & Go Architecture
- **Fail-Safe Goroutines (Fan-In Pattern):** Telemetry sensors push events into a buffered pipeline handled by the central orchestrator, isolated with `sync.WaitGroup`.
- **Context Timeboxing:** LLM API calls are strictly bound by `context.WithTimeout` to guarantee that network auto-remediation decisions are either taken rapidly or bypassed securely.
- **Structured SRE Logging:** `log/slog` is utilized natively for JSON-structured, traceable event lifecycle management.

## 🛠️ Technology Stack

| Component | Technology |
| :--- | :--- |
| **Network Fabric** | Nokia SR Linux (Containerlab) |
| **Telemetry & Observability** | eBPF Kernel Probes + gNMI Streams |
| **Orchestrator** | Go (Concurrency, Context Propagation) |
| **Decision Engine** | Anthropic Claude & Google Gemini |
| **Configuration Standard** | Native YANG (JSON) |

## 🚀 Usage Guide (Zero-Touch Deployment)

We have implemented a `Makefile` for a fully automated, plug-and-play experience.

### 1. Clone and Deploy
```bash
cd labs/01-nre-orchestrator
export GEMINI_API_KEY="your-api-key"
export OPENAI_API_KEY="your-api-key"
make deploy
```
This single command will:
- Build the minimalist Alpine Docker images for the Go agents (`aico` and `reconciler`).
- Spin up NetBox and its dependencies via Docker Compose.
- Deploy the Nokia SR Linux Spine-Leaf topology via Containerlab.
- Launch the Go agents alongside the network, connected via `network_mode: host` to access both Containerlab nodes and NetBox.

### 2. View Logs
To monitor the AI consensus and reconciliation process in real-time:
```bash
make logs
```

### 3. Teardown
To cleanly destroy the topology and containers:
```bash
make destroy
```
