# Terraform AWS VPC Base

Este laboratorio demuestra los fundamentos de **Infraestructura como Código (IaC)** utilizando Terraform. 

Provisiona una arquitectura de red base y segura en AWS, lista para desplegar instancias EC2 o clusters de EKS.

## 🏗️ Arquitectura Desplegada
*   **1 Virtual Private Cloud (VPC)**
*   **1 Internet Gateway (IGW)**
*   **1 Subred Pública (Public Subnet)**
*   **1 Tabla de Rutas (Route Table)** asociada a la subred pública.
*   **1 Security Group (SG)** configurado para permitir tráfico HTTP (80) y SSH (22).

## 🚀 Uso

### 1. Inicializar Terraform
Descarga los proveedores necesarios (AWS):
```bash
terraform init
```

### 2. Validar el plan
Verifica qué recursos serán creados:
```bash
terraform plan
```

### 3. Aplicar los cambios
Despliega la infraestructura en AWS (requiere AWS CLI configurado con tus credenciales):
```bash
terraform apply
```

### 4. Limpieza
Destruye todos los recursos para evitar cobros de AWS:
```bash
terraform destroy
```
