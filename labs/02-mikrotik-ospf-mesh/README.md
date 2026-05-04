# MikroTik OSPF Mesh Architecture

This repository section demonstrates a fully automated deployment of a highly available OSPF (Open Shortest Path First) mesh topology using MikroTik RouterOS virtual instances. It showcases advanced Network Reliability Engineering (NRE) principles by treating dynamic routing infrastructure as code.

Designed to emulate a Carrier-grade or Enterprise WAN core, this project automatically provisions and establishes an IGP (Interior Gateway Protocol) domain, ensuring immediate reconvergence and seamless path redundancy upon link failures.

---

## 🌟 Enterprise-Grade Technical Features

### 1. 🏗️ Infrastructure as Code (IaC) Topology
- **Declarative Blueprint:** The entire multi-node router topology is defined in a `topology.yaml` configuration, allowing for deterministic deployments and version control of the network layout.
- **Zero-Touch Provisioning (ZTP):** Custom `*.rsc` startup scripts (`R1-startup.rsc`, etc.) are dynamically bound to each router instance. This ensures that interfaces, IP addresses, and OSPF processes are fully configured at boot without any manual console intervention.

### 2. ⚡ High Availability Routing (OSPF)
- **Mesh Resiliency:** The routers are interconnected in a full or partial mesh. OSPF dynamically calculates the shortest path tree (SPF), automatically redirecting traffic over alternate links if a primary connection is disrupted.
- **Microsegmentation & Point-to-Point Adjacencies:** Interfaces are strictly configured with `/30` or `/31` point-to-point subnets to minimize broadcast traffic and accelerate OSPF adjacency establishment.

### 3. 🛡️ Network Reliability Standards
- **Configuration Consistency:** By centralizing the startup configurations, configuration drift across the router fleet is eliminated. Every deployment is identical and reproducible.
- **Scalable Design:** The topology is designed to easily integrate BGP on top of the established OSPF IGP, providing a solid foundation for more complex Tier-1 ISP architectures.

## 🛠️ Technology Stack

| Component | Technology |
| :--- | :--- |
| **Routing Engine** | MikroTik RouterOS |
| **Topology Manager** | Containerlab / YAML |
| **Dynamic Routing Protocol** | OSPFv2 (Interior Gateway Protocol) |
| **Provisioning** | RouterOS Scripting (`.rsc`) |

## 🚀 Usage Guide

1. Deploy the network simulation utilizing the defined topology:
   ```bash
   sudo clab deploy -t topology.yaml
   ```
2. Verify OSPF neighbor adjacencies across the nodes:
   ```bash
   docker exec -it clab-ospf-r1 ssh admin@localhost /routing ospf neighbor print
   ```
3. Test network resiliency by shutting down an interface and observing OSPF reconvergence.
