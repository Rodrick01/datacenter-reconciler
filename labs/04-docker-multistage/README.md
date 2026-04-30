# Multi-stage Docker Builds (Go)

Este laboratorio demuestra cómo crear contenedores **extremadamente ligeros y seguros** utilizando el patrón *Multi-stage Build* de Docker.

## 📝 Concepto
Cuando compilamos aplicaciones en Go, Node.js o Java, necesitamos muchas herramientas (compiladores, SDKs, librerías) que resultan en imágenes de cientos de megabytes (ej: `golang:1.22` pesa ~800MB). 

Sin embargo, el binario resultante no necesita el compilador para ejecutarse en producción. El patrón *Multi-stage* resuelve esto:
1.  **Stage 1 (Builder):** Usa la imagen pesada para compilar el código.
2.  **Stage 2 (Runtime):** Usa una imagen vacía (`scratch`) o muy ligera (`alpine`) y **sólo copia el binario final** desde el Stage 1.

**Resultado:** Una imagen final de ~5 MB que es mucho más rápida de desplegar y tiene una superficie de ataque (seguridad) casi nula.

## 🚀 Uso

### 1. Construir la imagen
```bash
docker build -t go-minimal-api .
```

### 2. Verificar el tamaño
```bash
docker images | grep go-minimal-api
```
*(Verás que la imagen pesa menos de 10MB)*

### 3. Ejecutar el contenedor
```bash
docker run -d -p 8080:8080 go-minimal-api
```

### 4. Probar
```bash
curl http://localhost:8080
```
