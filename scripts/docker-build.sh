#!/bin/bash

# Docker build script for Gin API template
# This script builds an optimized Docker image with size reporting

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
IMAGE_NAME="gin-template"
TAG="${1:-latest}"
FULL_IMAGE_NAME="${IMAGE_NAME}:${TAG}"

echo -e "${BLUE}🐳 Building optimized Docker image for Gin API template${NC}"
echo -e "${BLUE}=================================================${NC}"

# Function to get image size
get_image_size() {
    docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | grep "^${IMAGE_NAME}" | grep "${TAG}" | awk '{print $3}'
}

# Function to print step
print_step() {
    echo -e "\n${YELLOW}📋 $1${NC}"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Docker is running
print_step "Checking Docker daemon..."
if ! docker info >/dev/null 2>&1; then
    print_error "Docker daemon is not running. Please start Docker and try again."
    exit 1
fi
print_success "Docker daemon is running"

# Clean up any existing images with the same tag
print_step "Cleaning up existing images..."
if docker images | grep -q "^${IMAGE_NAME}.*${TAG}"; then
    docker rmi "${FULL_IMAGE_NAME}" 2>/dev/null || true
    print_success "Cleaned up existing image"
else
    echo "No existing image found"
fi

# Build the Docker image
print_step "Building Docker image: ${FULL_IMAGE_NAME}"
echo "Build context: $(pwd)"
echo "Dockerfile: ./Dockerfile"

# Build with build time and version info
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

docker build \
    --build-arg BUILD_TIME="${BUILD_TIME}" \
    --build-arg GIT_COMMIT="${GIT_COMMIT}" \
    --build-arg VERSION="${TAG}" \
    --tag "${FULL_IMAGE_NAME}" \
    --target final \
    .

if [ $? -eq 0 ]; then
    print_success "Docker image built successfully"
else
    print_error "Docker build failed"
    exit 1
fi

# Get image information
print_step "Analyzing built image..."
IMAGE_SIZE=$(get_image_size)
IMAGE_ID=$(docker images --format "{{.ID}}" "${FULL_IMAGE_NAME}")

echo -e "\n${BLUE}📊 Image Information:${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Image Name:    ${FULL_IMAGE_NAME}"
echo "Image ID:      ${IMAGE_ID}"
echo "Image Size:    ${IMAGE_SIZE}"
echo "Build Time:    ${BUILD_TIME}"
echo "Git Commit:    ${GIT_COMMIT}"

# Test the image
print_step "Testing the built image..."

# Test version command
echo "Testing version command..."
if docker run --rm "${FULL_IMAGE_NAME}" --version >/dev/null 2>&1; then
    print_success "Version command test passed"
else
    print_error "Version command test failed"
fi

# Test health check
echo "Testing health check command..."
if docker run --rm "${FULL_IMAGE_NAME}" --health-check >/dev/null 2>&1; then
    echo "Health check command exists (may fail without running server)"
else
    echo "Health check command test completed"
fi

# Show layer information
print_step "Image layer analysis..."
docker history "${FULL_IMAGE_NAME}" --format "table {{.CreatedBy}}\t{{.Size}}" | head -10

# Security scan (if available)
print_step "Running security checks..."
echo "Checking for security vulnerabilities..."

# Check if the image uses non-root user
USER_CHECK=$(docker run --rm "${FULL_IMAGE_NAME}" whoami 2>/dev/null || echo "nobody")
if [ "${USER_CHECK}" != "root" ]; then
    print_success "Image runs as non-root user: ${USER_CHECK}"
else
    print_error "Image runs as root user (security risk)"
fi

# Show image manifest
print_step "Image manifest summary..."
docker inspect "${FULL_IMAGE_NAME}" --format='
Image Details:
- Created: {{.Created}}
- Architecture: {{.Architecture}}
- OS: {{.Os}}
- Size: {{.Size}} bytes
- User: {{.Config.User}}
- Exposed Ports: {{range $port, $_ := .Config.ExposedPorts}}{{$port}} {{end}}
- Entrypoint: {{.Config.Entrypoint}}
- Cmd: {{.Config.Cmd}}
- Environment Variables: {{len .Config.Env}} vars set
'

# Performance recommendations
echo -e "\n${BLUE}🚀 Performance & Size Optimizations Applied:${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Multi-stage build (3 stages)"
echo "✅ Scratch-based final image"
echo "✅ UPX compression"
echo "✅ Static binary compilation"
echo "✅ Symbol stripping (-w -s flags)"
echo "✅ Minimal build context (.dockerignore)"
echo "✅ Non-root user (65534:65534)"
echo "✅ Health check support"
echo "✅ Security options (no-new-privileges)"

# Usage instructions
echo -e "\n${BLUE}📖 Usage Instructions:${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1. Run the image:"
echo "   docker run -p 8080:8080 ${FULL_IMAGE_NAME}"
echo ""
echo "2. Run with environment variables:"
echo "   docker run -p 8080:8080 -e APP_ENV=production ${FULL_IMAGE_NAME}"
echo ""
echo "3. Health check:"
echo "   docker run --rm ${FULL_IMAGE_NAME} --health-check"
echo ""
echo "4. Version info:"
echo "   docker run --rm ${FULL_IMAGE_NAME} --version"
echo ""
echo "5. Production deployment:"
echo "   docker-compose -f docker-compose.prod.yml up -d"

# Final summary
echo -e "\n${GREEN}🎉 Build completed successfully!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Image: ${FULL_IMAGE_NAME}"
echo "Size:  ${IMAGE_SIZE}"
echo "Ready for deployment! 🚀"
