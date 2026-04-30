# --- Stage 1: Builder ---
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# In this simple example we don't have a go.mod, but usually you'd copy it here.
# COPY go.mod go.sum ./
# RUN go mod download

# Copy source code
COPY main.go .

# Build the statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myserver main.go

# --- Stage 2: Minimal Runtime ---
FROM scratch

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/myserver /myserver

# Expose the application port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/myserver"]
