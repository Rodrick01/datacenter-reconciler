# AI Consensus Gateway

This laboratory introduces an **AI Consensus Gateway**, an enterprise-grade Go orchestrator designed to aggregate, evaluate, and mediate responses from multiple Large Language Models (LLMs) such as Claude, GPT, and Gemini. 

Engineered for high-reliability decision-making environments, this gateway abstracts the specific integrations with different AI providers, ensuring that downstream applications consume highly audited, deterministic, and synthesized intelligence rather than raw probabilistic outputs.

---

## 🌟 Enterprise-Grade Technical Features

### 1. 🏗️ Provider Agnosticism & Pluggable Architecture
- **Factory Pattern Implementation:** The application seamlessly instantiates LLM providers dynamically based on payload requests (`Thinker` vs `Auditor` roles). This abstraction guarantees zero vendor lock-in and seamless extendability.
- **Role-Based Consensus Engine:** Implements a multi-agent architectural design where a primary model generates an initial thesis, and a secondary auditor model critically evaluates it against constraints, ultimately converging on a highly refined response.

### 2. ⚡ High-Performance Concurrency & Resiliency
- **Standard Go Layout:** The binary execution `main.go` resides cleanly under `/cmd`, separating business logic (`/internal/ai`) from the HTTP presentation layer.
- **Context Propagation & I/O Safeguards:** Given the inherent latency of external LLM APIs, rigorous `context.Context` bounds (e.g., 120-second hard timeouts) prevent resource exhaustion and hanging HTTP sockets.
- **Fail-Safe SRE Patterns:** Comprehensive JSON validation, immediate descriptor closures (`defer r.Body.Close()`), and graceful fallback mechanisms prevent memory leaks under massive asynchronous loads.

### 3. 🛡️ Telemetry & Auditability
- **Structured JSON Logging:** Native integration with Go's `log/slog` ensures that all critical API transitions, model decisions, and errors are uniformly exported as structured key-value pairs, ready for ingestion by Elasticsearch, Datadog, or Splunk.

## 🛠️ Technology Stack

| Component | Technology |
| :--- | :--- |
| **Core API Gateway** | Go (Standard Library `net/http`) |
| **Model Integrations** | Native API SDKs (Gemini, OpenAI, Anthropic) |
| **Observability** | Structured Logging (`log/slog`) |
| **Concurrency Control** | `context.Context` & Goroutines |

## 🚀 Usage Guide

### 1. Start the Gateway
Compile and execute the gateway ensuring the corresponding API keys are present in the environment:
```bash
GEMINI_API_KEY="..." OPENAI_API_KEY="..." go run cmd/ai-consensus-gateway/main.go
```

### 2. Invoke Consensus Request
Trigger a multi-model consensus workflow via POST request:
```bash
curl -X POST http://localhost:8080/api/v1/consensus \
  -H "Content-Type: application/json" \
  -d '{
    "context": "Given an infrastructure topology...",
    "objective": "Determine the optimal BGP path selection",
    "thinker": "gemini",
    "auditor": "gpt"
  }'
```
