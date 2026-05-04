FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Compilamos ambos agentes de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o aico-agent ./cmd/aico
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o reconciler-agent ./cmd/reconciler

# Imagen final minimalista
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/aico-agent .
COPY --from=builder /app/reconciler-agent .

# Default command (puede ser sobrescrito en el docker-compose)
CMD ["./reconciler-agent"]
