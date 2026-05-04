# Terraform AWS VPC Architecture

This module provides a production-ready Infrastructure as Code (IaC) implementation for provisioning a secure, foundational Virtual Private Cloud (VPC) environment within Amazon Web Services (AWS) using HashiCorp Terraform.

Designed to reflect Site Reliability Engineering (SRE) best practices, this repository ensures deterministic, repeatable, and version-controlled infrastructure deployment, serving as the networking backbone for EC2 workloads, EKS clusters, or serverless architectures.

---

## 🌟 Enterprise-Grade Technical Features

### 1. 🏗️ Declarative Cloud Infrastructure
- **Idempotency:** State management is handled entirely by Terraform, guaranteeing that identical code always yields an identical infrastructure state without configuration drift.
- **Modularity:** The deployment logic is cleanly separated into `main.tf`, `variables.tf`, and `outputs.tf`, adhering to Terraform modular design patterns for maximum reusability and parameterization.

### 2. ⚡ Foundational Networking Setup
- **VPC & Subnet Isolation:** Provisions a fully customized Virtual Private Cloud alongside designated public subnets, strictly controlling the IP addressing scheme.
- **Internet Gateway & Routing:** Automatically attaches an Internet Gateway (IGW) and binds the corresponding Route Tables to ensure outbound internet connectivity for public-facing assets while maintaining granular routing control.

### 3. 🛡️ Security & Access Control
- **Security Group Hardening:** Enforces strict inbound traffic policies (e.g., restricted HTTP/SSH access) using precisely scoped Security Groups (SG) attached at the instance level.

## 🛠️ Technology Stack

| Component | Technology |
| :--- | :--- |
| **Cloud Provider** | Amazon Web Services (AWS) |
| **Infrastructure as Code** | HashiCorp Terraform |
| **State Management** | Local State (Extensible to S3/DynamoDB) |

## 🚀 Usage Guide

### 1. Initialization
Initialize the working directory containing Terraform configuration files and download necessary provider plugins:
```bash
terraform init
```

### 2. Plan Generation
Create an execution plan to preview the changes before provisioning resources:
```bash
terraform plan
```

### 3. Deployment
Apply the configuration to provision the infrastructure in the target AWS account (ensure AWS CLI credentials are valid):
```bash
terraform apply
```

### 4. Teardown
Destroy the Terraform-managed infrastructure to prevent unnecessary AWS billing:
```bash
terraform destroy
```
