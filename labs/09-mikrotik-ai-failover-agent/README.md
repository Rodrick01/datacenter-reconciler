# MikroTik AI Failover Agent (RouterOS v7 Containers)

Este proyecto implementa un agente escrito en Go "Vanilla" (sin dependencias externas pesadas) diseñado para ejecutarse nativamente dentro de un contenedor en un router MikroTik con RouterOS v7.

El agente se conecta a la REST API local del router para extraer la tabla de enrutamiento y la envía a un **AI Consensus Gateway** externo. Si el gateway detecta degradación en el ISP principal (por latencia o packet loss), devuelve un JSON estructurado para inyectar un PATCH que modifica la distancia administrativa (`distance`), ejecutando un failover "Self-Healing".

---

## 🏗️ Arquitectura
1. **RouterOS v7 Container:** El MikroTik corre el agente internamente consumiendo menos de 15MB de disco y recursos mínimos.
2. **Vanilla Go:** Código ultra liviano. Sin frameworks externos.
3. **Multi-Stage Build:** `Dockerfile` optimizado que usa `golang:alpine` para compilar un binario estático y `alpine:latest` para minimizar el runtime footprint.

## 🚀 Instrucciones de Despliegue (Plug-and-Play)

Hemos implementado un `Makefile` para automatizar la compilación y el despliegue del entorno virtual.

### 1. Compilación del Agente (Multi-Arquitectura)
El agente se compila en un contenedor estático. Usá el comando correspondiente según tu router:

**Para equipos MikroTik ARM64 (ej. RB5009, CCR2116):**
```bash
make build-arm64
```

**Para equipos MikroTik x86_64 o entornos virtuales (CHR):**
```bash
make build-amd64
```
*Esto generará un archivo `ai-agent-*.tar` listo para subir al router.*

### 2. Levantar el Entorno de Pruebas (Opcional)
Si no tenés un MikroTik físico, podés levantar uno virtualizado junto con dos ISPs simulados usando Containerlab:
```bash
make deploy-lab
```
*(Requiere tener la imagen `vrnetlab/vr-routeros:7.14` pre-cargada en tu Docker).*


### 2. Preparación en MikroTik
Asegurate de que el router tiene el usuario listo y la API expuesta:
```routeros
# Crear usuario exclusivo para el agente
/user add name=ai_agent password=modo_cientifico group=full

# Habilitar REST API
/ip service enable www
```

### 3. Cargar y Configurar el Container en MikroTik
1. Subí el archivo `ai-agent-*.tar` generado al almacenamiento del router (Files).

2. Configurá la interfaz Virtual Ethernet (VETH) y agregala a tu Bridge local (ej. `bridge-local`):
```routeros
# Crear la interfaz VETH para el container (asumiendo que tu red local es 192.168.88.0/24)
/interface veth add name=veth-ai-agent address=192.168.88.253/24 gateway=192.168.88.1

# Agregar el VETH al bridge de la LAN para que tenga salida y conectividad
/interface bridge port add bridge=bridge-local interface=veth-ai-agent
```

3. Configurá las variables de entorno para el agente:
```routeros
/container envs add name=ai_envs key=MIKROTIK_IP value="192.168.88.1"
/container envs add name=ai_envs key=MIKROTIK_USER value="ai_agent"
/container envs add name=ai_envs key=MIKROTIK_PASS value="modo_cientifico"
/container envs add name=ai_envs key=GATEWAY_URL value="http://<TU_IP_DEL_GATEWAY>:8080/api/v1/consensus"
```

4. Instanciá e iniciá el contenedor:
```routeros
/container add file=ai-agent-arm64.tar interface=veth-ai-agent envlist=ai_envs logging=yes
/container start 0
```

Revisá los logs del router para ver cómo el agente audita la tabla de ruteo y solicita decisiones al AI Gateway.
