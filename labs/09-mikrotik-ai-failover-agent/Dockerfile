# Stage 1: Build the Go binary
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod ./
# COPY go.sum ./ (Si lo hubiera)
RUN go mod download

COPY . .

# Compilar un binario 100% estático
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ai-agent main.go

# Stage 2: Minimal runtime image
FROM alpine:latest
WORKDIR /root/

# Copiar certificados TLS desde el builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar el binario compilado
COPY --from=builder /app/ai-agent .

# Ejecutar
CMD ["./ai-agent"]
