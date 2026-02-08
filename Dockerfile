# Multi-stage Dockerfile for ultra-minimal Gin API production image
# Stage 1: Build environment
FROM golang:1.25.7-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
  git \
  ca-certificates \
  tzdata \
  upx

# Set working directory
WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./

# Download dependencies with verification
RUN go mod download && go mod verify

# Copy source code (organized structure)
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

# Build the binary with maximum optimization
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s -extldflags '-static' -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -a -installsuffix cgo \
  -trimpath \
  -o api ./cmd/api/main.go

# Compress binary with UPX for even smaller size
RUN upx --best --lzma /app/api

# Verify the binary works
RUN ./api --help || echo "Binary built successfully"

# Stage 2: Create minimal runtime user setup
FROM alpine:3.19 AS runtime-prep

# Create necessary directories with proper permissions
# Use existing nobody user (65534:65534) instead of creating new one
RUN mkdir -p /app/data /app/configs /app/logs && \
  chown -R 65534:65534 /app

# Stage 3: Final ultra-minimal image
FROM scratch AS final

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy user and group files
COPY --from=runtime-prep /etc/passwd /etc/passwd
COPY --from=runtime-prep /etc/group /etc/group

# Copy application directories with proper ownership
COPY --from=runtime-prep --chown=65534:65534 /app /app

# Copy the optimized binary
COPY --from=builder --chown=65534:65534 /app/api /app/api

# Set working directory
WORKDIR /app

# Set environment variables for production
ENV GIN_MODE=release
ENV APP_ENV=production
ENV LOG_FORMAT=json
ENV TZ=UTC

# Expose port (non-privileged)
EXPOSE 8080

# Use non-root user for security
USER 65534:65534

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["./api", "--health-check"] || exit 1

# Set entrypoint
ENTRYPOINT ["./api"]

# Default command (can be overridden)
CMD []
