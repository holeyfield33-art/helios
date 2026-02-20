# Multi-stage build: Go binary + Python conformance harness
# Stage 1: Build Go binary
FROM golang:1.23-alpine AS go-builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /helios ./cmd/helios/
RUN CGO_ENABLED=0 go test ./... -count=1

# Stage 2: Runtime with Go binary + Python
FROM python:3.12-slim

# Copy Go binary
COPY --from=go-builder /helios /usr/local/bin/helios

# Copy test vectors and scripts
COPY test_vectors/ /app/test_vectors/
COPY scripts/ /app/scripts/
COPY implementations/python/ /app/implementations/python/

WORKDIR /app

# Make scripts executable
RUN chmod +x /app/scripts/*.sh

# Default: run cross-language check
CMD ["/app/scripts/cross_check.sh"]
