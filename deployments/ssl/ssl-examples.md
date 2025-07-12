# SSL/TLS Configuration Examples for URL Shortener Service

## Overview

This document provides comprehensive examples for configuring SSL/TLS in the URL Shortener service across different deployment scenarios.

## Development Environment

### Self-Signed Certificate Generation

For development environments, you can generate self-signed certificates:

```go
package main

import (
    "log"
    "url-shortener/internal/config"
    "url-shortener/internal/core/services"
    "url-shortener/internal/infrastructure/ssl"
)

func main() {
    cfg := config.New()
    logger := services.NewLoggingService(cfg)
    
    certManager := ssl.NewCertificateManager(cfg, logger)
    
    // Create certificate request
    req := certManager.GetDefaultCertificateRequest()
    req.CommonName = "localhost"
    req.DNSNames = []string{"localhost", "127.0.0.1", "url-shortener.local"}
    req.ValidityDays = 365
    
    // Generate certificate
    certPEM, keyPEM, err := certManager.GenerateSelfSignedCertificate(req)
    if err != nil {
        log.Fatalf("Failed to generate certificate: %v", err)
    }
    
    // Save to files
    err = certManager.SaveCertificateToFile(
        certPEM, keyPEM,
        "./certs/server.crt",
        "./certs/server.key",
    )
    if err != nil {
        log.Fatalf("Failed to save certificate: %v", err)
    }
    
    log.Println("Certificate generated successfully!")
}
```

### Environment Configuration

```bash
# .env.development
GO_ENV=development
TLS_ENABLED=true
TLS_CERT_FILE=./certs/server.crt
TLS_KEY_FILE=./certs/server.key
TLS_MIN_VERSION=1.2
TRUST_PROXY=false
ENABLE_HTTPS=false  # Allow HTTP in development
```

## Production Environment

### Production Certificate Configuration

```bash
# .env.production
GO_ENV=production
TLS_ENABLED=true
TLS_CERT_FILE=/etc/ssl/certs/url-shortener.crt
TLS_KEY_FILE=/etc/ssl/private/url-shortener.key
TLS_MIN_VERSION=1.2
TRUST_PROXY=true
ENABLE_HTTPS=true
ENABLE_SECURITY_HEADERS=true
```

### Let's Encrypt with Certbot

```bash
#!/bin/bash
# certbot-setup.sh

# Install Certbot
sudo apt update
sudo apt install -y certbot

# Generate certificate
sudo certbot certonly --standalone \
  --email admin@yourdomain.com \
  --agree-tos \
  --no-eff-email \
  -d yourdomain.com \
  -d www.yourdomain.com

# Create symbolic links
sudo ln -sf /etc/letsencrypt/live/yourdomain.com/fullchain.pem /etc/ssl/certs/url-shortener.crt
sudo ln -sf /etc/letsencrypt/live/yourdomain.com/privkey.pem /etc/ssl/private/url-shortener.key

# Set permissions
sudo chmod 644 /etc/ssl/certs/url-shortener.crt
sudo chmod 600 /etc/ssl/private/url-shortener.key
sudo chown root:ssl-cert /etc/ssl/private/url-shortener.key

# Setup auto-renewal
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

## Docker Configuration

### Dockerfile with SSL Support

```dockerfile
# Multi-stage build for production
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest

# Install SSL certificates and create ssl-cert group
RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S ssl-cert \
    && adduser -S appuser -G ssl-cert

# Create directories for certificates
RUN mkdir -p /etc/ssl/certs /etc/ssl/private \
    && chown root:ssl-cert /etc/ssl/private \
    && chmod 750 /etc/ssl/private

WORKDIR /root/

COPY --from=builder /app/server .
COPY --chown=root:ssl-cert deployments/ssl/certs/ /etc/ssl/

# Switch to non-root user
USER appuser

EXPOSE 8080 8443

CMD ["./server"]
```

### Docker Compose with SSL

```yaml
version: '3.8'

services:
  url-shortener:
    build: .
    ports:
      - "8080:8080"
      - "8443:8443"
    environment:
      - GO_ENV=production
      - TLS_ENABLED=true
      - TLS_CERT_FILE=/etc/ssl/certs/server.crt
      - TLS_KEY_FILE=/etc/ssl/private/server.key
      - TLS_MIN_VERSION=1.2
      - PORT=8080
      - TLS_PORT=8443
    volumes:
      - ./certs:/etc/ssl/certs:ro
      - ./private:/etc/ssl/private:ro
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/ssl/certs:ro
      - ./private:/etc/ssl/private:ro
    depends_on:
      - url-shortener
    restart: unless-stopped
```

## Kubernetes Deployment

### SSL Certificate Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: url-shortener-tls
  namespace: production
type: kubernetes.io/tls
data:
  tls.crt: LS0tLS1CRUdJTi... # base64 encoded certificate
  tls.key: LS0tLS1CRUdJTi... # base64 encoded private key
```

### Deployment with TLS

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: url-shortener
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: url-shortener
  template:
    metadata:
      labels:
        app: url-shortener
    spec:
      containers:
      - name: url-shortener
        image: url-shortener:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 8443
          name: https
        env:
        - name: GO_ENV
          value: "production"
        - name: TLS_ENABLED
          value: "true"
        - name: TLS_CERT_FILE
          value: "/etc/ssl/certs/tls.crt"
        - name: TLS_KEY_FILE
          value: "/etc/ssl/private/tls.key"
        volumeMounts:
        - name: tls-certs
          mountPath: /etc/ssl/certs
          readOnly: true
        - name: tls-private
          mountPath: /etc/ssl/private
          readOnly: true
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: tls-certs
        secret:
          secretName: url-shortener-tls
          items:
          - key: tls.crt
            path: tls.crt
      - name: tls-private
        secret:
          secretName: url-shortener-tls
          items:
          - key: tls.key
            path: tls.key
```

## Certificate Management

### Automated Certificate Renewal

```bash
#!/bin/bash
# certificate-renewal.sh

SERVICE_NAME="url-shortener"
CERT_PATH="/etc/ssl/certs/url-shortener.crt"
LOG_FILE="/var/log/cert-renewal.log"

log_message() {
    echo "$(date): $1" >> "$LOG_FILE"
}

check_certificate_expiry() {
    local days_until_expiry
    days_until_expiry=$(openssl x509 -in "$CERT_PATH" -noout -dates | grep notAfter | cut -d= -f2 | xargs -I {} date -d {} +%s)
    local current_time=$(date +%s)
    local seconds_until_expiry=$((days_until_expiry - current_time))
    local days_until_expiry=$((seconds_until_expiry / 86400))
    
    echo $days_until_expiry
}

main() {
    log_message "Starting certificate renewal check"
    
    local days_until_expiry
    days_until_expiry=$(check_certificate_expiry)
    
    if [ "$days_until_expiry" -lt 30 ]; then
        log_message "Certificate expires in $days_until_expiry days, renewing..."
        
        # Renew with Certbot
        certbot renew --quiet
        
        if [ $? -eq 0 ]; then
            log_message "Certificate renewed successfully"
            
            # Restart service
            systemctl reload "$SERVICE_NAME"
            log_message "Service reloaded"
        else
            log_message "Certificate renewal failed"
            exit 1
        fi
    else
        log_message "Certificate is valid for $days_until_expiry more days"
    fi
}

main "$@"
```

### Certificate Monitoring Script

```go
package main

import (
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "log"
    "time"
)

func checkCertificateExpiry(certPath string) error {
    certPEM, err := ioutil.ReadFile(certPath)
    if err != nil {
        return fmt.Errorf("failed to read certificate: %w", err)
    }
    
    block, _ := pem.Decode(certPEM)
    if block == nil {
        return fmt.Errorf("failed to decode PEM block")
    }
    
    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        return fmt.Errorf("failed to parse certificate: %w", err)
    }
    
    now := time.Now()
    daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
    
    fmt.Printf("Certificate: %s\n", certPath)
    fmt.Printf("Subject: %s\n", cert.Subject.String())
    fmt.Printf("Issuer: %s\n", cert.Issuer.String())
    fmt.Printf("Not Before: %s\n", cert.NotBefore.Format(time.RFC3339))
    fmt.Printf("Not After: %s\n", cert.NotAfter.Format(time.RFC3339))
    fmt.Printf("Days until expiry: %d\n", daysUntilExpiry)
    
    if daysUntilExpiry < 30 {
        log.Printf("WARNING: Certificate expires in %d days!", daysUntilExpiry)
    }
    
    if daysUntilExpiry < 7 {
        log.Printf("CRITICAL: Certificate expires in %d days!", daysUntilExpiry)
    }
    
    return nil
}

func main() {
    certPaths := []string{
        "/etc/ssl/certs/url-shortener.crt",
        "./certs/server.crt",
    }
    
    for _, certPath := range certPaths {
        if err := checkCertificateExpiry(certPath); err != nil {
            log.Printf("Error checking %s: %v", certPath, err)
        }
        fmt.Println("---")
    }
}
```

## Security Best Practices

### TLS Configuration Hardening

```go
// Production TLS config with security hardening
func createProductionTLSConfig() *tls.Config {
    return &tls.Config{
        MinVersion: tls.VersionTLS12,
        MaxVersion: tls.VersionTLS13,
        CipherSuites: []uint16{
            // TLS 1.3 cipher suites (configured automatically)
            // TLS 1.2 cipher suites
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
            tls.CurveP384,
        },
        PreferServerCipherSuites: false, // Client preference for modern browsers
        SessionTicketsDisabled:   false, // Enable for performance
        Renegotiation:           tls.RenegotiateNever,
        InsecureSkipVerify:      false,
    }
}
```

### HSTS Configuration

```go
// Add security headers including HSTS
func addSecurityHeaders(w http.ResponseWriter, r *http.Request) {
    // HSTS Header (only for HTTPS)
    if r.TLS != nil {
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
    }
    
    // Additional security headers
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")
    w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
    w.Header().Set("Content-Security-Policy", "default-src 'self'")
}
```

## Troubleshooting

### Common SSL/TLS Issues

1. **Certificate chain issues**:
   ```bash
   # Check certificate chain
   openssl s_client -connect yourdomain.com:443 -showcerts
   ```

2. **Cipher suite problems**:
   ```bash
   # Test SSL configuration
   nmap --script ssl-enum-ciphers -p 443 yourdomain.com
   ```

3. **Certificate validation**:
   ```bash
   # Validate certificate
   openssl x509 -in server.crt -text -noout
   ```

4. **Key/certificate mismatch**:
   ```bash
   # Compare certificate and key
   openssl x509 -noout -modulus -in server.crt | openssl md5
   openssl rsa -noout -modulus -in server.key | openssl md5
   ```

This comprehensive SSL/TLS configuration ensures secure communication across all deployment scenarios while maintaining flexibility for different environments.