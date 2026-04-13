# ── Build Stage ──
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies first to optimize Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api_binary ./cmd/api/main.go

# ── Final Stage (Small Image) ──
FROM alpine:latest

WORKDIR /app

# Install CA certificates for external HTTPS requests (e.g., Keycloak JWKS, S3/MinIO APIs)
RUN apk --no-cache add ca-certificates tzdata

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/api_binary .

# Expose the API port
EXPOSE 8080

# Command to run the executable
CMD ["./api_binary"]
