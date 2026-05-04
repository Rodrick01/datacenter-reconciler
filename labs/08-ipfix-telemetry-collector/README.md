# High-Performance IPFIX Telemetry Collector

This laboratory demonstrates a production-grade **IPFIX (NetFlow v10) Collector** built entirely in Go. Designed for Tier-1 ISPs and large-scale datacenters, this telemetry aggregator is capable of ingesting massive UDP flow bursts, decoding binary structs with zero-memory-copy paradigms, and exposing aggregated metrics to Prometheus/Grafana for real-time observability.

---

## 🌟 Enterprise-Grade Technical Features

### 1. ⚡ Zero-Copy Binary Decoding
- **Memory Optimization:** Instead of creating intermediate Go structs for every packet, the decoder uses slice windowing and `encoding/binary` to extract integers and IPs directly from the raw byte buffer. This prevents the Garbage Collector from thrashing under heavy loads (e.g., 10+ Gbps flow sampling).
- **Socket Buffer Tuning:** Modifies the OS-level UDP `ReadBuffer` explicitly up to 10MB to handle network micro-bursts without dropping packets at the kernel level.

### 2. 🛡️ Idempotent Template Caching
- **State Resilience:** IPFIX separates "Templates" (data schemas) from the actual "Data Records". The collector uses `sync.RWMutex` to maintain an idempotent cache of templates per Observation Domain ID.
- **Out-of-Order Handling:** If Data Records arrive before their corresponding templates (or if UDP packets are duplicated/lost), the decoder fails gracefully, dropping the record without corrupting internal states or crashing the daemon.

### 3. 📊 SRE Observability (Prometheus)
- **Thread-Safe Metrics:** Uses `promauto` to maintain thread-safe counters for bytes transferred, packet drops, and parsed flows.
- **Graceful Shutdown:** Implements `context.Context` combined with `sync.WaitGroup` to capture OS interrupts (SIGTERM), ensuring the HTTP exporter and UDP sockets are cleanly flushed and closed before process termination.

## 🛠️ Technology Stack

| Component | Technology |
| :--- | :--- |
| **Telemetry Protocol** | IPFIX (NetFlow v10) over UDP |
| **Collector Daemon** | Go (`net.ListenUDP`, `context`, `sync`) |
| **Metrics Exporter** | Prometheus Client Golang (`/metrics`) |
| **Visualization** | Grafana (Auto-provisioned Dashboards) |
| **Infrastructure** | Docker Compose |

## 🚀 Usage Guide

### 1. Spin up the Telemetry Stack
Start the Go Collector, Prometheus, and Grafana using Docker Compose:
```bash
docker-compose up -d --build
```

### 2. Verify Dashboards
Navigate to Grafana (http://localhost:3000) (admin/admin). The `IPFIX Telemetry & SRE Metrics` dashboard will be automatically provisioned.

### 3. Inject Test Flows
To test the collector, you can inject synthetic IPFIX traffic using tools like `tcpreplay` or a Python script targeting `udp://localhost:4739`. Watch the Prometheus counters (`ipfix_udp_packets_received_total`) increment in real-time.
