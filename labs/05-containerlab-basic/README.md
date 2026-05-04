# Network Virtualization with Containerlab

This laboratory introduces foundational concepts of modern network orchestration using **Containerlab**. It demonstrates how to rapidly deploy and connect virtualized network elements using lightweight Docker containers, bridging the gap between traditional networking hardware and Cloud-Native DevOps workflows.

Designed for Network Reliability Engineering (NRE) automation, this environment replaces heavy virtual machine hypervisors with high-density containerized routing topologies, facilitating rapid prototyping and CI/CD network testing.

---

## 🌟 Enterprise-Grade Technical Features

### 1. 🏗️ High-Density Network Emulation
- **Resource Efficiency:** By leveraging the Docker runtime, Containerlab simulates complete network nodes (routers, switches, and hosts) using a fraction of the RAM and CPU required by traditional emulators like GNS3 or EVE-NG.
- **Microsecond Boot Times:** Network nodes boot almost instantaneously, drastically accelerating the validation loop for configuration changes and network automation scripts.

### 2. ⚡ Declarative YAML Topologies
- **Infrastructure as Code (IaC):** The entire lab topology, including node definitions, virtual image mappings, and point-to-point links (veth pairs), is declared in a single `basic.yaml` file.
- **Reproducible Environments:** The declarative nature of the topology guarantees that any engineer can spin up the exact same network state in seconds, mitigating "it works on my machine" discrepancies.

### 3. 🛡️ Seamless CLI Integration
- **Native Management:** Containerlab interacts natively with standard Linux networking tools (like `iproute2` and `bridge-utils`) and the Docker CLI.
- **Automated Lifecycle Management:** Commands such as `deploy`, `destroy`, and `graph` abstract away the underlying complexity of manually wiring Linux namespaces and virtual ethernet cables.

## 🛠️ Technology Stack

| Component | Technology |
| :--- | :--- |
| **Virtualization Runtime** | Docker Engine |
| **Topology Manager** | Containerlab |
| **Configuration Standard** | YAML |

## 🚀 Usage Guide

### 1. Topology Deployment
Provision the virtual network by pointing Containerlab to the declarative topology file:
```bash
sudo clab deploy -t basic.yaml
```

### 2. Verify Connectivity
Access the virtual nodes and test data-plane connectivity using standard tools (like Ping or Traceroute):
```bash
docker exec -it clab-basic-node1 ping 192.168.1.2
```

### 3. Teardown
Clean up the namespaces and containers quickly to free up host resources:
```bash
sudo clab destroy -t basic.yaml
```
