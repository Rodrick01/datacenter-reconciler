FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ipfix-collector ./cmd/ipfix-collector

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ipfix-collector .

EXPOSE 4739/udp
EXPOSE 2112

CMD ["./ipfix-collector"]
