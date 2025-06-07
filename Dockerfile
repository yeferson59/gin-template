FROM golang:1.24.4-alpine AS builder

# Instalar dependencias del sistema necesarias
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copiar archivos de dependencias para aprovechar cache de Docker
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar código fuente
COPY cmd/ cmd/
COPY internal/ internal/

# Compilar el binario con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o api ./cmd/api/main.go

# Etapa 2: Imagen final ultra-minimalista
FROM scratch

# Copiar certificados CA y zona horaria desde builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar el binario
COPY --from=builder /app/api /api

# Exponer puerto
EXPOSE 8080

# Usuario no-root para seguridad
USER 65534:65534

# Punto de entrada
ENTRYPOINT ["/api"]
