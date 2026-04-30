# Site Reliability Engineering (SRE) & Network Automation Portfolio

<div align="center">
  <p><strong>A curated collection of laboratories demonstrating Enterprise-grade Engineering</strong></p>
  <p>Go | Terraform | Docker | Containerlab | CI/CD | AI Integrations</p>
</div>

---

## 📖 Executive Summary

This monorepo serves as a technical showcase of my hands-on experience across multiple disciplines: Software Engineering, Cloud Infrastructure, and Network Reliability Engineering (NRE). 

Each folder in `labs/` represents an isolated, production-ready implementation of a specific technology or pattern. 

---

## 📂 The Laboratories

### 1. `labs/01-nre-orchestrator` (The Masterpiece)
**Technologies:** Go, gNMI, Docker, NetBox
*   An Enterprise-grade orchestrator that reconciles a Datacenter Fabric (Spine-Leaf) state from a Single Source of Truth (NetBox) to physical/virtual routers (Nokia SR Linux).
*   Implements **strict Go concurrency patterns** (Worker Pools, Buffered Channels, Zero-Copy data structures).
*   Includes **AICO (AI Consensus Orchestrator)**: An experimental feature that connects the orchestrator to a dual-LLM (Thinker/Auditor) pipeline to automatically resolve network anomalies using pure YANG models.

### 2. `labs/02-mikrotik-ospf-mesh`
**Technologies:** Containerlab, MikroTik RouterOS v7, OSPF
*   A 4-node virtualized mesh topology simulating an ISP core.
*   Demonstrates fast-convergence and Zero-Touch provisioning via `.rsc` configuration scripts.

### 3. `labs/03-terraform-aws-vpc`
**Technologies:** Terraform (HCL), AWS
*   A clean, modular Infrastructure-as-Code (IaC) deployment.
*   Provisions a highly-available AWS Virtual Private Cloud (VPC), Public Subnet, Internet Gateway, and Security Groups, ready for EKS or EC2 clusters.

### 4. `labs/04-docker-multistage`
**Technologies:** Docker, Go
*   Demonstrates container optimization using **Multi-stage Builds**.
*   Compiles a Go web server in a heavy builder image and extracts the static binary into a `scratch` container, resulting in a production image of ~5MB with near-zero attack surface.

### 5. `labs/05-containerlab-basic`
**Technologies:** Containerlab, Linux Namespaces
*   A foundational network topology connecting two lightweight Alpine Linux containers point-to-point. 
*   Used for rapid, low-overhead network protocol testing.

---

## ⚙️ Continuous Integration (CI/CD)

This repository enforces quality via **GitHub Actions** (`.github/workflows/ci.yml`). Every push to the `main` branch automatically:
1.  Downloads dependencies and runs `go test` and `go build` for all Go projects.
2.  Executes `terraform fmt -check` and `terraform validate` on the AWS infrastructure.
3.  Lints all YAML topology files.

> *"Automating infrastructure with the rigor of software development."*
