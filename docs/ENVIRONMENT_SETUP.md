# Environment Setup Guide

This guide provides comprehensive instructions for setting up the URL Shortener service in different environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Environment](#development-environment)
- [Production Environment](#production-environment)
- [Environment Variables](#environment-variables)
- [Database Setup](#database-setup)
- [Redis Setup](#redis-setup)
- [SSL/TLS Configuration](#ssltls-configuration)
- [Health Checks](#health-checks)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher
- **PostgreSQL**: Version 15 or higher
- **Redis**: Version 7 or higher
- **Docker**: Version 24 or higher (optional)
- **Docker Compose**: Version 2.x (optional)

### Recommended System Specifications

#### Development
- **CPU**: 2+ cores
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 10GB free space

#### Production
- **CPU**: 4+ cores
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 50GB+ SSD storage
- **Network**: 1Gbps connection

## Development Environment

### Quick Start with Docker

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd url-shortener
   ```

2. **Copy environment files**:
   ```bash
   cp .env.example .env
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```

3. **Start services with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

4. **Verify installation**:
   ```bash
   # Check backend health
   curl http://localhost:8080/health

   # Check frontend
   curl http://localhost:3000
   ```

### Manual Development Setup

#### 1. Database Setup

**Install PostgreSQL**:
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# Windows
# Download from https://www.postgresql.org/download/windows/
```

**Create database and user**:
```sql
-- Connect as postgres user
sudo -u postgres psql

-- Create database and user
CREATE DATABASE urlshortener;
CREATE USER urlshortener WITH ENCRYPTED PASSWORD 'password123';
GRANT ALL PRIVILEGES ON DATABASE urlshortener TO urlshortener;

-- Exit psql
\q
```

#### 2. Redis Setup

**Install Redis**:
```bash
# Ubuntu/Debian
sudo apt install redis-server

# macOS
brew install redis
brew services start redis

# Windows
# Download from https://github.com/MicrosoftArchive/redis/releases
```

**Configure Redis**:
```bash
# Edit Redis configuration
sudo nano /etc/redis/redis.conf

# Add password authentication
requirepass redis123

# Restart Redis
sudo systemctl restart redis
```

#### 3. Backend Setup

**Install Go dependencies**:
```bash
cd backend
go mod download
go mod verify
```

**Configure environment**:
```bash
# Copy and edit .env file
cp .env.example .env
nano .env
```

**Run database migrations**:
```bash
go run cmd/migrate/main.go
```

**Start backend server**:
```bash
go run cmd/server/main.go
```

#### 4. Frontend Setup

**Install Node.js dependencies**:
```bash
cd frontend
npm install
```

**Configure environment**:
```bash
# Copy and edit .env file
cp .env.example .env
nano .env
```

**Start development server**:
```bash
npm run dev
```

## Production Environment

### System Preparation

#### 1. Security Hardening

**Update system packages**:
```bash
sudo apt update && sudo apt upgrade -y
```

**Configure firewall**:
```bash
# Install UFW
sudo apt install ufw

# Allow SSH
sudo ufw allow ssh

# Allow HTTP/HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Allow backend port (if not behind reverse proxy)
sudo ufw allow 8080

# Enable firewall
sudo ufw enable
```

**Create application user**:
```bash
sudo useradd -r -m -s /bin/bash urlshortener
sudo mkdir -p /opt/urlshortener
sudo chown urlshortener:urlshortener /opt/urlshortener
```

#### 2. Database Production Setup

**Install and configure PostgreSQL**:
```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Configure PostgreSQL
sudo -u postgres psql

-- Create production database
CREATE DATABASE urlshortener_prod;
CREATE USER urlshortener_prod WITH ENCRYPTED PASSWORD 'STRONG_PRODUCTION_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE urlshortener_prod TO urlshortener_prod;

-- Configure connection limits
ALTER USER urlshortener_prod CONNECTION LIMIT 20;

\q
```

**Optimize PostgreSQL configuration**:
```bash
# Edit PostgreSQL configuration
sudo nano /etc/postgresql/15/main/postgresql.conf

# Key production settings
shared_buffers = 256MB                    # 25% of RAM
effective_cache_size = 1GB               # 75% of RAM
work_mem = 4MB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1                   # For SSD storage
max_connections = 100

# Restart PostgreSQL
sudo systemctl restart postgresql
```

#### 3. Redis Production Setup

**Install and configure Redis**:
```bash
# Install Redis
sudo apt install redis-server

# Configure Redis for production
sudo nano /etc/redis/redis.conf

# Key production settings
bind 127.0.0.1
port 6379
requirepass STRONG_REDIS_PASSWORD
maxmemory 512mb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000

# Enable Redis service
sudo systemctl enable redis-server
sudo systemctl restart redis-server
```

#### 4. Application Deployment

**Build backend application**:
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -o bin/urlshortener cmd/server/main.go
```

**Build frontend application**:
```bash
cd frontend
npm run build
```

**Deploy application**:
```bash
# Copy backend binary
sudo cp backend/bin/urlshortener /opt/urlshortener/
sudo chown urlshortener:urlshortener /opt/urlshortener/urlshortener
sudo chmod +x /opt/urlshortener/urlshortener

# Copy frontend build
sudo cp -r frontend/dist /opt/urlshortener/public
sudo chown -R urlshortener:urlshortener /opt/urlshortener/public
```

#### 5. Systemd Service Configuration

**Create backend service**:
```bash
sudo nano /etc/systemd/system/urlshortener-backend.service
```

```ini
[Unit]
Description=URL Shortener Backend
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=urlshortener
Group=urlshortener
WorkingDirectory=/opt/urlshortener
ExecStart=/opt/urlshortener/urlshortener
Restart=always
RestartSec=5
EnvironmentFile=/opt/urlshortener/.env

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/urlshortener

[Install]
WantedBy=multi-user.target
```

**Enable and start service**:
```bash
sudo systemctl daemon-reload
sudo systemctl enable urlshortener-backend
sudo systemctl start urlshortener-backend
sudo systemctl status urlshortener-backend
```

## Environment Variables

### Backend Environment Variables

```bash
# Server Configuration
GO_ENV=production                              # Environment: development, production
PORT=8080                                      # Server port
HOST=0.0.0.0                                  # Server host

# Database Configuration
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
DB_HOST=localhost                              # Database host
DB_PORT=5432                                   # Database port
DB_NAME=urlshortener                          # Database name
DB_USER=urlshortener                          # Database user
DB_PASSWORD=password123                        # Database password
DB_SSLMODE=require                            # SSL mode: disable, require, verify-ca, verify-full
DB_MAX_OPEN_CONNS=25                          # Maximum open connections
DB_MAX_IDLE_CONNS=10                          # Maximum idle connections
DB_CONN_MAX_LIFETIME=300s                     # Connection maximum lifetime

# Redis Configuration
REDIS_URL=redis://:password@host:6379/0       # Redis connection URL
REDIS_HOST=localhost                          # Redis host
REDIS_PORT=6379                               # Redis port
REDIS_PASSWORD=redis123                       # Redis password
REDIS_DB=0                                    # Redis database number
REDIS_MAX_RETRIES=3                           # Maximum retry attempts
REDIS_POOL_SIZE=10                            # Connection pool size

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key          # JWT signing key (256-bit)
JWT_ACCESS_EXPIRY=15m                         # Access token expiry
JWT_REFRESH_EXPIRY=168h                       # Refresh token expiry (7 days)

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://yourdomain.com   # Allowed origins (comma-separated)
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE      # Allowed methods
CORS_ALLOWED_HEADERS=*                        # Allowed headers
CORS_ALLOW_CREDENTIALS=true                   # Allow credentials

# Rate Limiting
RATE_LIMIT_ENABLED=true                       # Enable rate limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=100           # Global rate limit
AUTH_RATE_LIMIT_REQUESTS_PER_MINUTE=5        # Auth endpoints rate limit
URL_CREATION_RATE_LIMIT_REQUESTS_PER_MINUTE=10  # URL creation rate limit

# Security Configuration
BCRYPT_COST=12                                # Bcrypt hashing cost
SESSION_SECRET=your-session-secret            # Session secret key
ENCRYPTION_KEY=your-32-byte-encryption-key    # AES encryption key

# External Services
GEOLOCATION_API_KEY=your-geolocation-api-key  # IP geolocation service
GEOLOCATION_SERVICE_URL=https://api.service.com  # Geolocation service URL

# Monitoring and Logging
LOG_LEVEL=info                                # Log level: debug, info, warn, error
LOG_FORMAT=json                               # Log format: json, text
METRICS_ENABLED=true                          # Enable metrics collection
HEALTH_CHECK_ENABLED=true                     # Enable health checks

# File Storage
UPLOAD_MAX_SIZE=10MB                          # Maximum upload size
STATIC_FILES_PATH=/opt/urlshortener/public   # Static files directory

# Application Configuration
BASE_URL=https://yourdomain.com              # Base URL for short links
FRONTEND_URL=https://yourdomain.com          # Frontend URL
SHORT_CODE_LENGTH=6                          # Default short code length
DEFAULT_EXPIRY_HOURS=8760                    # Default URL expiry (1 year)
MAX_CUSTOM_ALIAS_LENGTH=50                   # Maximum custom alias length
```

### Frontend Environment Variables

```bash
# API Configuration
VITE_API_BASE_URL=https://api.yourdomain.com/api/v1  # Backend API URL
VITE_BASE_URL=https://yourdomain.com                 # Frontend base URL

# Authentication
VITE_JWT_STORAGE_KEY=url_shortener_token            # JWT storage key
VITE_REFRESH_TOKEN_STORAGE_KEY=url_shortener_refresh # Refresh token storage key

# Features
VITE_ENABLE_ANALYTICS=true                          # Enable analytics features
VITE_ENABLE_QR_CODES=true                          # Enable QR code features
VITE_ENABLE_BULK_OPERATIONS=true                   # Enable bulk operations

# UI Configuration
VITE_APP_TITLE=URL Shortener                       # Application title
VITE_APP_DESCRIPTION=Professional URL shortening   # Application description
VITE_THEME_PRIMARY_COLOR=#667eea                   # Primary theme color
VITE_PAGINATION_SIZE=10                            # Default pagination size

# External Services
VITE_GOOGLE_ANALYTICS_ID=GA-XXXXXXXXX              # Google Analytics ID
VITE_SENTRY_DSN=https://sentry.io/dsn              # Sentry error tracking

# Development
VITE_DEBUG_MODE=false                              # Enable debug mode
VITE_MOCK_API=false                               # Use mock API responses
```

### Environment Validation

#### Backend Validation Script

Create `backend/scripts/validate-env.go`:

```go
package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    required := []string{
        "GO_ENV", "PORT", "DATABASE_URL", "REDIS_URL", 
        "JWT_SECRET", "CORS_ALLOWED_ORIGINS",
    }
    
    missing := []string{}
    
    for _, env := range required {
        if os.Getenv(env) == "" {
            missing = append(missing, env)
        }
    }
    
    if len(missing) > 0 {
        fmt.Printf("Missing required environment variables: %s\n", 
            strings.Join(missing, ", "))
        os.Exit(1)
    }
    
    // Validate JWT secret length
    if len(os.Getenv("JWT_SECRET")) < 32 {
        fmt.Println("JWT_SECRET must be at least 32 characters long")
        os.Exit(1)
    }
    
    fmt.Println("Environment validation passed")
}
```

#### Frontend Validation Script

Create `frontend/scripts/validate-env.js`:

```javascript
const requiredEnvVars = [
    'VITE_API_BASE_URL',
    'VITE_BASE_URL'
];

const missingVars = requiredEnvVars.filter(envVar => !process.env[envVar]);

if (missingVars.length > 0) {
    console.error('Missing required environment variables:', missingVars.join(', '));
    process.exit(1);
}

// Validate URLs
const urlVars = ['VITE_API_BASE_URL', 'VITE_BASE_URL'];
for (const urlVar of urlVars) {
    const url = process.env[urlVar];
    if (url && !url.match(/^https?:\/\/.+/)) {
        console.error(`Invalid URL format for ${urlVar}: ${url}`);
        process.exit(1);
    }
}

console.log('Environment validation passed');
```

## Database Setup

### Production Database Configuration

#### Connection Pool Optimization

```go
// In backend/internal/infrastructure/database/postgres.go
func configureConnectionPool(db *gorm.DB) {
    sqlDB, _ := db.DB()
    
    // Set maximum number of open connections
    sqlDB.SetMaxOpenConns(25)
    
    // Set maximum number of idle connections
    sqlDB.SetMaxIdleConns(10)
    
    // Set maximum lifetime of connections
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

#### Database Monitoring

```sql
-- Monitor connection usage
SELECT count(*) as connections, state 
FROM pg_stat_activity 
WHERE datname = 'urlshortener_prod' 
GROUP BY state;

-- Monitor slow queries
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements 
WHERE mean_time > 1000 
ORDER BY mean_time DESC 
LIMIT 10;
```

### Migration Management

#### Production Migration Strategy

```bash
# Create migration
migrate create -ext sql -dir migrations add_indexes

# Test migration in staging
migrate -path migrations -database "$STAGING_DATABASE_URL" up

# Backup production database before migration
pg_dump $PRODUCTION_DATABASE_URL > backup_$(date +%Y%m%d_%H%M%S).sql

# Run migration in production
migrate -path migrations -database "$PRODUCTION_DATABASE_URL" up

# Verify migration
migrate -path migrations -database "$PRODUCTION_DATABASE_URL" version
```

## Redis Setup

### Production Redis Configuration

#### Memory Optimization

```bash
# /etc/redis/redis.conf

# Memory management
maxmemory 1gb
maxmemory-policy allkeys-lru

# Persistence
save 900 1      # Save if at least 1 key changed in 900 seconds
save 300 10     # Save if at least 10 keys changed in 300 seconds
save 60 10000   # Save if at least 10000 keys changed in 60 seconds

# Network
tcp-keepalive 300
timeout 0

# Security
requirepass your_strong_redis_password
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command DEBUG ""
```

#### Redis Monitoring

```bash
# Monitor Redis performance
redis-cli info stats
redis-cli info memory
redis-cli info clients

# Monitor slow commands
redis-cli slowlog get 10
```

## SSL/TLS Configuration

### Let's Encrypt with Certbot

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d yourdomain.com -d api.yourdomain.com

# Test renewal
sudo certbot renew --dry-run

# Add automatic renewal to crontab
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

### Nginx Configuration

```nginx
# /etc/nginx/sites-available/urlshortener
server {
    listen 80;
    server_name yourdomain.com api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    ssl_dhparam /etc/nginx/dhparam.pem;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=63072000" always;
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Referrer-Policy "strict-origin-when-cross-origin";
    
    # Frontend
    location / {
        root /opt/urlshortener/public;
        try_files $uri $uri/ /index.html;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # Backend API
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

## Health Checks

### Application Health Checks

#### Backend Health Check

The backend includes comprehensive health checks at `/health`:

```json
{
  "status": "healthy",
  "service": "url-shortener",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "components": {
    "database": {
      "status": "up",
      "last_checked": "2024-01-15T10:30:00Z",
      "duration": "2ms"
    },
    "redis": {
      "status": "up",
      "last_checked": "2024-01-15T10:30:00Z",
      "duration": "1ms"
    }
  }
}
```

#### System-Level Health Checks

```bash
# Create health check script
sudo nano /opt/urlshortener/health-check.sh
```

```bash
#!/bin/bash

# Health check script for URL Shortener

check_service() {
    local service=$1
    local url=$2
    local expected_status=$3
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" || echo "000")
    
    if [ "$response" -eq "$expected_status" ]; then
        echo "✓ $service is healthy"
        return 0
    else
        echo "✗ $service is unhealthy (HTTP $response)"
        return 1
    fi
}

echo "URL Shortener Health Check - $(date)"
echo "=================================="

# Check backend
check_service "Backend API" "http://localhost:8080/health" 200
backend_status=$?

# Check database connectivity
check_service "Database Health" "http://localhost:8080/health" 200
db_status=$?

# Check Redis connectivity
redis_status=0
if ! redis-cli -a "$REDIS_PASSWORD" ping > /dev/null 2>&1; then
    echo "✗ Redis is unhealthy"
    redis_status=1
else
    echo "✓ Redis is healthy"
fi

# Overall status
if [ $backend_status -eq 0 ] && [ $db_status -eq 0 ] && [ $redis_status -eq 0 ]; then
    echo "✓ All services are healthy"
    exit 0
else
    echo "✗ Some services are unhealthy"
    exit 1
fi
```

```bash
# Make executable
sudo chmod +x /opt/urlshortener/health-check.sh

# Add to crontab for monitoring
echo "*/5 * * * * /opt/urlshortener/health-check.sh >> /var/log/urlshortener-health.log 2>&1" | sudo crontab -
```

## Troubleshooting

### Common Issues

#### Database Connection Issues

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log

# Test database connection
psql -h localhost -U urlshortener -d urlshortener_prod -c "SELECT 1;"
```

#### Redis Connection Issues

```bash
# Check Redis status
sudo systemctl status redis-server

# Check Redis logs
sudo tail -f /var/log/redis/redis-server.log

# Test Redis connection
redis-cli -a "$REDIS_PASSWORD" ping
```

#### Backend Service Issues

```bash
# Check service status
sudo systemctl status urlshortener-backend

# View service logs
sudo journalctl -u urlshortener-backend -f

# Check application logs
sudo tail -f /opt/urlshortener/app.log
```

#### Performance Issues

```bash
# Monitor system resources
htop
iotop
nethogs

# Monitor database performance
SELECT * FROM pg_stat_activity WHERE state = 'active';

# Monitor Redis performance
redis-cli info stats
```

### Log Analysis

#### Centralized Logging Setup

```bash
# Install and configure rsyslog for centralized logging
sudo apt install rsyslog

# Configure application logging
echo "local0.*    /var/log/urlshortener/app.log" | sudo tee -a /etc/rsyslog.conf
sudo systemctl restart rsyslog
```

#### Log Rotation

```bash
# Configure logrotate
sudo nano /etc/logrotate.d/urlshortener
```

```
/var/log/urlshortener/*.log {
    daily
    missingok
    rotate 30
    compress
    notifempty
    create 644 urlshortener urlshortener
    postrotate
        systemctl reload urlshortener-backend
    endscript
}
```

### Monitoring and Alerting

#### Basic Monitoring Script

```bash
#!/bin/bash
# /opt/urlshortener/monitor.sh

ALERT_EMAIL="admin@yourdomain.com"
LOG_FILE="/var/log/urlshortener-monitor.log"

log_message() {
    echo "$(date): $1" >> "$LOG_FILE"
}

send_alert() {
    local subject="$1"
    local message="$2"
    echo "$message" | mail -s "$subject" "$ALERT_EMAIL"
    log_message "ALERT: $subject"
}

# Check disk space
disk_usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$disk_usage" -gt 85 ]; then
    send_alert "High Disk Usage" "Disk usage is at ${disk_usage}%"
fi

# Check memory usage
memory_usage=$(free | awk 'NR==2{printf "%.2f", $3*100/$2}')
if (( $(echo "$memory_usage > 85" | bc -l) )); then
    send_alert "High Memory Usage" "Memory usage is at ${memory_usage}%"
fi

# Check service health
if ! /opt/urlshortener/health-check.sh > /dev/null 2>&1; then
    send_alert "Service Health Check Failed" "URL Shortener health check failed"
fi

log_message "Monitor check completed"
```

This comprehensive environment setup guide provides everything needed to deploy and maintain the URL Shortener service in both development and production environments.