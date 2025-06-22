# API Documentation

## Overview

This is a professional Gin API template with JWT authentication, rate limiting, structured logging, and comprehensive middleware support.

## Base URL

```
http://localhost:8080
```

## Authentication

This API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Rate Limiting

- General endpoints: 10 requests per second per IP
- Authentication endpoints: 5 requests per minute per IP

## Health Check Endpoints

### GET /health/

Returns the overall health status of the application.

**Response:**
```json
{
  "success": true,
  "message": "Health check completed",
  "data": {
    "status": "ok",
    "timestamp": "2023-12-01T12:00:00Z",
    "version": "1.0.0",
    "services": {
      "database": "ok"
    }
  }
}
```

### GET /health/live

Liveness probe for Kubernetes.

### GET /health/ready

Readiness probe for Kubernetes.

## Authentication Endpoints

### POST /api/auth/register

Register a new user.

**Request Body:**
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "Password123!"
}
```

**Validation Rules:**
- Username: 3-30 characters, alphanumeric with underscores and hyphens
- Email: Valid email format
- Password: Minimum 8 characters with uppercase, lowercase, number, and special character

**Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

### POST /api/auth/login

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "username": "testuser",
  "password": "Password123!"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

## Protected Endpoints

All endpoints below require authentication via JWT token.

### GET /api/protected/

Example protected endpoint.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200):**
```json
{
  "success": true,
  "message": "Access granted to protected resource",
  "data": {
    "user_id": 1,
    "email": "test@example.com",
    "username": "testuser",
    "message": "You have successfully accessed a protected resource"
  }
}
```

### GET /api/users/me

Get current user profile.

**Response (200):**
```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": "Additional details about the error"
  }
}
```

### Common Error Codes

- `BAD_REQUEST` - Invalid request data
- `UNAUTHORIZED` - Authentication required or invalid
- `FORBIDDEN` - Access denied
- `NOT_FOUND` - Resource not found
- `CONFLICT` - Resource already exists
- `VALIDATION_ERROR` - Input validation failed
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `INTERNAL_SERVER_ERROR` - Server error

## Status Codes

- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `429` - Too Many Requests
- `500` - Internal Server Error

## Request/Response Headers

### Security Headers (Applied to all responses)

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000; includeSubDomains`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'`

### Request ID

Each request receives a unique ID in the `X-Request-ID` header for tracking purposes.

## Environment Variables

See `.env.example` for all available configuration options.

## Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

### Access protected resource

```bash
curl -X GET http://localhost:8080/api/protected/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
