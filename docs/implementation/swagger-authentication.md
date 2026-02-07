# Swagger Authentication Implementation

## Overview

This project implements **Basic Authentication** to protect Swagger documentation across different environments with flexible configuration.

## Architecture

**Middleware**: `internal/shared/web/middleware/swagger_auth.go`
**Configuration**: `configs/config.go`
**Route Registration**: `internal/infra/web/register_routes.go`

## How It Works

### Development Environment
- **Access**: Open (no authentication)
- **URL**: `http://localhost:8080/swagger/index.html`
- **Credentials**: Not required

### Staging Environment
- **Access**: Protected with Basic Auth
- **URL**: `http://staging-server:8080/swagger/index.html`
- **Credentials**: Required (configured in `.env.staging`)

### Production Environment
- **Access**: Disabled by default (or highly protected if enabled)
- **URL**: `http://production-server:8080/swagger/index.html`
- **Credentials**: Optional (for emergency debugging)

## Configuration

### Environment Variables

```bash
# Environment type (determines authentication behavior)
SERVER_APP_ENVIRONMENT=development|staging|production

# Enable/disable Swagger completely
SERVER_APP_SWAGGER_ENABLED=true|false

# Basic Auth credentials (required for staging/production if enabled)
SERVER_APP_SWAGGER_USER=username
SERVER_APP_SWAGGER_PASS=password
```

### Configuration Files

We provide three example configurations:

1. **`.env.development`** - Open access for local development
2. **`.env.staging`** - Protected with Basic Auth
3. **`.env.production`** - Disabled or highly secured

Copy the appropriate file to `.env` based on your environment:

```bash
# Development
cp cmd/server/.env.development cmd/server/.env

# Staging
cp cmd/server/.env.staging cmd/server/.env

# Production
cp cmd/server/.env.production cmd/server/.env
```

## Usage Examples

### Development (No Authentication)

**Configuration** (`.env`):
```bash
SERVER_APP_ENVIRONMENT=development
SERVER_APP_SWAGGER_ENABLED=true
SERVER_APP_SWAGGER_USER=
SERVER_APP_SWAGGER_PASS=
```

**Access**:
- Open browser: `http://localhost:8080/swagger/index.html`
- No credentials required ✅

---

### Staging (Basic Auth)

**Configuration** (`.env`):
```bash
SERVER_APP_ENVIRONMENT=staging
SERVER_APP_SWAGGER_ENABLED=true
SERVER_APP_SWAGGER_USER=swagger_admin
SERVER_APP_SWAGGER_PASS=StagingPass123!
```

**Access**:
1. Open browser: `http://staging-server:8080/swagger/index.html`
2. Browser shows popup asking for credentials:
   ```
   ┌─────────────────────────────────────┐
   │  Authentication Required            │
   │                                     │
   │  The server says:                   │
   │  "Swagger Documentation -           │
   │   Restricted Access"                │
   │                                     │
   │  Username: [swagger_admin        ] │
   │  Password: [•••••••••••••        ] │
   │                                     │
   │  [Cancel]           [Sign In]      │
   └─────────────────────────────────────┘
   ```
3. Enter credentials
4. Access granted ✅

**Alternative - cURL**:
```bash
curl -u swagger_admin:StagingPass123! \
  http://staging-server:8080/swagger/doc.json
```

---

### Production (Disabled)

**Configuration** (`.env`):
```bash
SERVER_APP_ENVIRONMENT=production
SERVER_APP_SWAGGER_ENABLED=false
SERVER_APP_SWAGGER_USER=
SERVER_APP_SWAGGER_PASS=
```

**Access**:
- Open browser: `http://production-server:8080/swagger/index.html`
- Response: `403 Forbidden - "Swagger documentation is disabled"`

---

### Production (Emergency Access)

**Configuration** (`.env`):
```bash
SERVER_APP_ENVIRONMENT=production
SERVER_APP_SWAGGER_ENABLED=true
SERVER_APP_SWAGGER_USER=emergency_admin
SERVER_APP_SWAGGER_PASS=VeryStr0ngProductionP@ssw0rd!
```

**Access**:
1. Open browser: `http://production-server:8080/swagger/index.html`
2. Browser requests authentication
3. Enter strong credentials
4. Access granted (for debugging only) ⚠️

## Security Best Practices

### ✅ DO's

1. **Use strong passwords in production/staging**
   ```bash
   # Bad
   SERVER_APP_SWAGGER_PASS=123456
   
   # Good
   SERVER_APP_SWAGGER_PASS=K$9mP#xL2@nQ5zR!vY8wT
   ```

2. **Keep credentials in secrets management** (AWS Secrets Manager, Vault, etc.)

3. **Disable Swagger in production by default**
   ```bash
   SERVER_APP_SWAGGER_ENABLED=false
   ```

4. **Use different credentials per environment**
   - Dev: No auth
   - Staging: Medium strength
   - Production: Very strong or disabled

5. **Rotate credentials regularly** (especially in staging)

### ❌ DON'Ts

1. **Never commit real credentials to git**
   ```bash
   # Add to .gitignore
   cmd/server/.env
   cmd/server/.env.staging
   cmd/server/.env.production
   ```

2. **Don't use the same password as other systems**

3. **Don't share production Swagger credentials in Slack/Email**

4. **Don't leave Swagger open in production indefinitely**

## Troubleshooting

### "Authentication required" in development

**Problem**: Getting auth popup in dev environment

**Solution**: Check your `.env` file:
```bash
# Must be "development"
SERVER_APP_ENVIRONMENT=development
```

---

### "Swagger authentication not configured"

**Problem**: Getting 503 error in staging/production

**Solution**: Set credentials in `.env`:
```bash
SERVER_APP_SWAGGER_USER=your_username
SERVER_APP_SWAGGER_PASS=your_password
```

---

### Browser keeps asking for credentials

**Problem**: Popup appears repeatedly even with correct password

**Solution**:
1. Clear browser cache/cookies
2. Try incognito mode
3. Verify credentials in `.env` file
4. Check for typos (no extra spaces)

---

### How to bypass authentication temporarily

**For development/testing only**:
```bash
# Change environment temporarily
SERVER_APP_ENVIRONMENT=development
```

Then restart the server.

## Middleware Implementation Details

### Authentication Flow

```
Request to /swagger/*
       ↓
SwaggerBasicAuth Middleware
       ↓
Check SWAGGER_ENABLED
       ├─ false → 403 Forbidden
       └─ true → Continue
              ↓
       Check ENVIRONMENT
              ├─ development → Allow
              └─ staging/production
                     ↓
              Validate Basic Auth
                     ├─ valid → Allow
                     └─ invalid → 401 Unauthorized
```

### Response Codes

- **200 OK** - Authenticated successfully
- **401 Unauthorized** - Wrong credentials or no credentials provided
- **403 Forbidden** - Swagger disabled or access denied
- **503 Service Unavailable** - Authentication not configured properly

## Integration with CI/CD

### Docker Build

```dockerfile
# Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

# Build with environment-specific config
RUN go build -o server cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/server /server
COPY --from=builder /app/cmd/server/.env.production /cmd/server/.env

CMD ["/server"]
```

### Kubernetes Secrets

```yaml
# k8s/swagger-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: swagger-credentials
type: Opaque
stringData:
  SWAGGER_USER: "staging_admin"
  SWAGGER_PASS: "SecurePassword123!"
```

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: api
        env:
        - name: SERVER_APP_SWAGGER_USER
          valueFrom:
            secretKeyRef:
              name: swagger-credentials
              key: SWAGGER_USER
        - name: SERVER_APP_SWAGGER_PASS
          valueFrom:
            secretKeyRef:
              name: swagger-credentials
              key: SWAGGER_PASS
```

## Testing Authentication

### Manual Test

```bash
# Test without credentials (should fail in staging/production)
curl http://localhost:8080/swagger/doc.json

# Test with credentials
curl -u username:password http://localhost:8080/swagger/doc.json
```

### Automated Test

```go
// internal/shared/web/middleware/swagger_auth_test.go
package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/refortunato/go_app_base/internal/shared/web/middleware"
)

func TestSwaggerBasicAuth_Production(t *testing.T) {
    os.Setenv("SERVER_APP_ENVIRONMENT", "production")
    os.Setenv("SERVER_APP_SWAGGER_USER", "testuser")
    os.Setenv("SERVER_APP_SWAGGER_PASS", "testpass")
    
    router := gin.New()
    router.Use(middleware.SwaggerBasicAuth())
    router.GET("/test", func(c *gin.Context) {
        c.String(200, "ok")
    })
    
    // Without auth - should fail
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusUnauthorized {
        t.Errorf("Expected 401, got %d", w.Code)
    }
    
    // With auth - should succeed
    req.SetBasicAuth("testuser", "testpass")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Summary

✅ **Development**: Open access for fast development
✅ **Staging**: Protected with Basic Auth for team testing
✅ **Production**: Disabled by default, or highly secured for emergency debugging
✅ **Flexible**: Easy to configure via environment variables
✅ **Secure**: Multiple layers of protection
