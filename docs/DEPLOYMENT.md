# 🚀 Deployment Guide

## Overview

This guide covers various deployment strategies for the Gin API Template, from development to production environments.

## 📋 Prerequisites

- Docker and Docker Compose
- Kubernetes cluster (for K8s deployment)
- PostgreSQL database (for production)
- Load balancer (for production)
- SSL certificates (for HTTPS)

## 🏠 Local Development

### Using Docker Compose

```bash
# Start all services
make up-build

# Services started:
# - API server on port 8080
# - PostgreSQL on port 5432
# - pgAdmin on port 5050 (with --profile admin)

# View logs
make logs-api

# Stop services
make down
```

### Environment Variables for Development

```env
APP_ENV=development
LOG_LEVEL=debug
LOG_FORMAT=text
DB_DRIVER=postgres
DB_DSN=host=postgres user=api_user password=secure_password_123 dbname=my_api_db port=5432 sslmode=disable
```

## 🏗️ Production Deployment

### 1. Environment Setup

#### Production Environment Variables

```env
# Application
APP_NAME=MyAPI
APP_ENV=production
PORT=8080

# Server Configuration
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
MAX_BODY_SIZE=10485760  # 10MB

# Database (PostgreSQL recommended)
DB_DRIVER=postgres
DB_DSN=host=your-db-host user=api_user password=your-secure-password dbname=production_db port=5432 sslmode=require
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m

# JWT (CRITICAL: Use strong secrets!)
JWT_SECRET=your-super-secret-256-bit-key-here
JWT_EXP_MINUTES=15m
JWT_REFRESH_MINUTES=7d
JWT_ISSUER=your-api-production

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Security
RATE_LIMIT_RPS=5.0
RATE_LIMIT_BURST=10
AUTH_RATE_LIMIT=3
CORS_ENABLED=true
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### 2. Docker Production Build

#### Production Dockerfile

The included Dockerfile is already optimized for production:

```dockerfile
# Multi-stage build with optimizations
FROM golang:1.24.4-alpine AS builder
# ... build stage

FROM scratch
# Ultra-minimal final image
COPY --from=builder /app/api /api
EXPOSE 8080
USER 65534:65534
ENTRYPOINT ["/api"]
```

#### Build and Push

```bash
# Build production image
docker build -t your-registry/gin-api:latest .

# Push to registry
docker push your-registry/gin-api:latest
```

### 3. Cloud Deployment Options

#### A. Docker Swarm

```yaml
# docker-stack.yml
version: '3.8'

services:
  api:
    image: your-registry/gin-api:latest
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - JWT_SECRET=${JWT_SECRET}
      - DB_DSN=${DB_DSN}
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
    networks:
      - api-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - api
    networks:
      - api-network

networks:
  api-network:
    driver: overlay
```

Deploy:
```bash
docker stack deploy -c docker-stack.yml gin-api-stack
```

#### B. AWS ECS

```json
{
  "family": "gin-api",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "gin-api",
      "image": "your-registry/gin-api:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "APP_ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "JWT_SECRET",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:jwt-secret"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/gin-api",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

## ☸️ Kubernetes Deployment

### 1. Configuration Maps and Secrets

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gin-api-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
  RATE_LIMIT_RPS: "5.0"
  RATE_LIMIT_BURST: "10"
---
apiVersion: v1
kind: Secret
metadata:
  name: gin-api-secrets
type: Opaque
data:
  JWT_SECRET: <base64-encoded-secret>
  DB_DSN: <base64-encoded-connection-string>
```

### 2. Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gin-api
  labels:
    app: gin-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gin-api
  template:
    metadata:
      labels:
        app: gin-api
    spec:
      containers:
      - name: gin-api
        image: your-registry/gin-api:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: gin-api-config
        - secretRef:
            name: gin-api-secrets
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

### 3. Service and Ingress

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: gin-api-service
spec:
  selector:
    app: gin-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
---
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: gin-api-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    secretName: gin-api-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gin-api-service
            port:
              number: 80
```

### 4. Deploy to Kubernetes

```bash
# Apply configurations
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml

# Check deployment status
kubectl get deployments
kubectl get pods
kubectl logs -l app=gin-api
```

## 🗄️ Database Setup

### PostgreSQL Production Setup

```sql
-- Create database and user
CREATE DATABASE production_api_db;
CREATE USER api_user WITH ENCRYPTED PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE production_api_db TO api_user;

-- Connect to the new database
\c production_api_db;

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO api_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO api_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO api_user;
```

### Database Migration

The application automatically runs GORM migrations on startup. For production, consider:

1. **Backup before deployment**
2. **Run migrations manually** if preferred
3. **Use database migration tools** like `golang-migrate`

## 🔧 Monitoring and Observability

### Health Checks

The API provides comprehensive health checks:

- **Liveness**: `GET /health/live` - Container is running
- **Readiness**: `GET /health/ready` - Application is ready to serve traffic
- **Health**: `GET /health/` - Detailed service status

### Logging

Production logging configuration:

```yaml
# In production use JSON logging
LOG_FORMAT: json
LOG_LEVEL: info

# Log aggregation with tools like:
# - ELK Stack (Elasticsearch, Logstash, Kibana)
# - Fluentd/Fluent Bit
# - AWS CloudWatch
# - Google Cloud Logging
```

### Metrics (Future Enhancement)

Consider adding:
- Prometheus metrics endpoint
- Grafana dashboards
- Application Performance Monitoring (APM)

## 🔒 Security Considerations

### Production Security Checklist

- [ ] **Strong JWT secrets** (256-bit random)
- [ ] **HTTPS only** (SSL/TLS certificates)
- [ ] **Database encryption** at rest and in transit
- [ ] **Network policies** (firewall rules)
- [ ] **Container security** (non-root user, minimal image)
- [ ] **Rate limiting** configured appropriately
- [ ] **CORS** configured for your domains only
- [ ] **Secrets management** (not in environment variables)
- [ ] **Regular security updates**
- [ ] **Monitoring and alerting**

### SSL/TLS Configuration

#### Nginx Configuration Example

```nginx
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    ssl_certificate /etc/ssl/certs/your-cert.pem;
    ssl_certificate_key /etc/ssl/private/your-key.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    ssl_dhparam /etc/ssl/certs/dhparam.pem;

    location / {
        proxy_pass http://gin-api-service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

## 🚦 CI/CD Pipeline

### GitHub Actions Example

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - run: go test ./...

  build-and-deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Docker image
        run: docker build -t ${{ secrets.REGISTRY }}/gin-api:${{ github.sha }} .
        
      - name: Push to registry
        run: |
          echo ${{ secrets.REGISTRY_PASSWORD }} | docker login -u ${{ secrets.REGISTRY_USER }} --password-stdin
          docker push ${{ secrets.REGISTRY }}/gin-api:${{ github.sha }}
          
      - name: Deploy to Kubernetes
        run: |
          echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig
          kubectl set image deployment/gin-api gin-api=${{ secrets.REGISTRY }}/gin-api:${{ github.sha }}
```

## 📊 Performance Tuning

### Application Optimization

```env
# Database connection pooling
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m

# Server timeouts
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s

# Rate limiting
RATE_LIMIT_RPS=5.0
RATE_LIMIT_BURST=10
```

### Resource Limits

```yaml
# Kubernetes resource limits
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

## 🔄 Backup and Recovery

### Database Backups

```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > backup_$DATE.sql
aws s3 cp backup_$DATE.sql s3://your-backup-bucket/
```

### Application State

The application is stateless, so recovery involves:
1. Restore database from backup
2. Deploy latest application version
3. Verify health checks

## 📞 Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check connection string
   - Verify network connectivity
   - Check credentials

2. **JWT Token Issues**
   - Verify JWT_SECRET is consistent
   - Check token expiration
   - Validate token format

3. **Rate Limiting**
   - Check rate limit configuration
   - Monitor rate limit metrics
   - Adjust limits if needed

### Debug Commands

```bash
# Check application logs
kubectl logs -l app=gin-api

# Port forward for debugging
kubectl port-forward deployment/gin-api 8080:8080

# Check health endpoints
curl https://api.yourdomain.com/health/
curl https://api.yourdomain.com/health/live
curl https://api.yourdomain.com/health/ready
```

## 🎯 Next Steps

After successful deployment:

1. **Monitor application performance**
2. **Set up alerting and notifications**
3. **Implement log aggregation**
4. **Add metrics and dashboards**
5. **Plan for scaling**
6. **Security audits and updates**

For detailed API documentation, see [`docs/api.md`](./api.md).
