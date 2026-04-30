# Containerlab: Topología Linux Básica

Este laboratorio (o "sonzera") demuestra la capacidad de crear topologías de red virtuales utilizando contenedores ligeros mediante **Containerlab**.

Es el bloque fundacional para probar conectividad, protocolos y scripts de automatización sin necesidad de levantar pesadas máquinas virtuales o hardware real.

## 🏗️ Topología
*   **2 Nodos (pc1, pc2):** Basados en la imagen `alpine:latest`.
*   **1 Enlace Punto a Punto:** Conectando `eth1` de `pc1` con `eth1` de `pc2`.
*   **Direccionamiento:**
    *   pc1: `10.0.0.1/24`
    *   pc2: `10.0.0.2/24`

## 🚀 Uso

### 1. Desplegar el Laboratorio
Requiere tener Docker y Containerlab instalados.
```bash
sudo containerlab deploy -t basic.yaml
```

### 2. Verificar Conectividad
Prueba que `pc1` puede alcanzar a `pc2` enviando paquetes ICMP (Ping):
```bash
docker exec -it clab-basic-linux-topo-pc1 ping 10.0.0.2 -c 4
```

### 3. Destruir el Laboratorio
Para limpiar el entorno y los namespaces de red:
```bash
sudo containerlab destroy -t basic.yaml
```
