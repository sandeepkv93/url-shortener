# Docker Deployment Guide

This guide covers Docker-based deployment for the URL Shortener service, including development and production configurations.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
- [Production Deployment](#production-deployment)
- [Docker Configurations](#docker-configurations)
- [Environment Variables](#environment-variables)
- [Data Persistence](#data-persistence)
- [Monitoring](#monitoring)
- [Security](#security)
- [Troubleshooting](#troubleshooting)

## Overview

The URL Shortener service uses Docker for containerization with the following components:

- **Backend**: Go application with multi-stage Dockerfile
- **Frontend**: React/Vite application with Nginx in production
- **Database**: PostgreSQL 15 with optimized configuration
- **Cache**: Redis 7 with production tuning
- **Reverse Proxy**: Nginx with SSL termination (production)
- **Monitoring**: Prometheus and Grafana (optional)

## Prerequisites

### System Requirements

**Development:**
- Docker 24.0+
- Docker Compose 2.x
- 4GB RAM minimum
- 10GB disk space

**Production:**
- Docker 24.0+
- Docker Compose 2.x
- 8GB RAM minimum
- 50GB SSD storage
- SSL certificates

### Installation

```bash
# Install Docker (Ubuntu/Debian)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

## Development Setup

### Quick Start

```bash
# 1. Clone repository
git clone <repository-url>
cd url-shortener

# 2. Copy environment files
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# 3. Start development environment
docker-compose up -d

# 4. View logs
docker-compose logs -f

# 5. Access services
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Database Admin: http://localhost:8081
# Redis Admin: http://localhost:8082
```

### Development Services

The development setup includes additional tools:

```bash
# Start with development tools
docker-compose --profile tools up -d

# Start with mail testing
docker-compose --profile mail-testing up -d

# Start everything
docker-compose --profile tools --profile mail-testing up -d
```

### Development Commands

```bash
# View service status
docker-compose ps

# Follow logs for specific service
docker-compose logs -f backend

# Restart a service
docker-compose restart backend

# Rebuild and restart
docker-compose up -d --build backend

# Execute commands in container
docker-compose exec backend sh
docker-compose exec postgres psql -U urlshortener

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Live Reload

Development containers support live reload:

- **Backend**: Source code is mounted to `/app`, restart required for Go changes
- **Frontend**: Hot reload via Vite dev server on port 5173

## Production Deployment

### Pre-deployment Setup

1. **Prepare server**:
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Create application user
sudo useradd -r -m -s /bin/bash urlshortener
sudo usermod -aG docker urlshortener

# Create application directory
sudo mkdir -p /opt/urlshortener
sudo chown urlshortener:urlshortener /opt/urlshortener
```

2. **Deploy application**:
```bash
# Switch to application user
sudo su - urlshortener
cd /opt/urlshortener

# Clone repository
git clone <repository-url> .

# Create data directories
mkdir -p data/{postgres,redis,prometheus,grafana}
mkdir -p ssl nginx-sites monitoring logs

# Set permissions
chmod 700 data/
chmod 600 .env.production
```

3. **Configure environment**:
```bash
# Copy and edit production environment
cp .env.production .env
nano .env

# Update all CHANGE_THIS_* values with strong passwords
# Set correct domain names and URLs
# Configure SSL certificate paths
```

### SSL Certificate Setup

#### Using Let's Encrypt

```bash
# Install Certbot
sudo apt install certbot

# Generate certificates
sudo certbot certonly --standalone -d yourdomain.com -d api.yourdomain.com

# Copy certificates to application directory
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem /opt/urlshortener/ssl/
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem /opt/urlshortener/ssl/
sudo chown urlshortener:urlshortener /opt/urlshortener/ssl/*

# Set up automatic renewal
echo "0 12 * * * /usr/bin/certbot renew --quiet && docker-compose -f /opt/urlshortener/docker-compose.prod.yml restart nginx" | sudo crontab -
```

#### Using Custom Certificates

```bash
# Copy your certificates
cp your-certificate.pem /opt/urlshortener/ssl/fullchain.pem
cp your-private-key.pem /opt/urlshortener/ssl/privkey.pem
chmod 600 /opt/urlshortener/ssl/*
```

### Nginx Configuration

Create production Nginx configuration:

```bash
# Create main nginx configuration
cat > nginx-prod.conf << 'EOF'
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
                    
    access_log /var/log/nginx/access.log main;
    
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 20m;
    
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types
        application/javascript
        application/json
        application/xml
        text/css
        text/javascript
        text/plain
        text/xml
        image/svg+xml;
    
    include /etc/nginx/conf.d/*.conf;
}
EOF

# Create site configuration
mkdir -p nginx-sites
cat > nginx-sites/default.conf << 'EOF'
# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name yourdomain.com api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

# Main frontend
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    
    add_header Strict-Transport-Security "max-age=63072000" always;
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    
    location / {
        proxy_pass http://frontend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Short URL redirection
    location ~ ^/[a-zA-Z0-9_-]+$ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# API backend
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    
    location / {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
EOF
```

### Production Deployment

```bash
# Deploy with production configuration
docker-compose -f docker-compose.prod.yml up -d

# Check service status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Test services
curl https://yourdomain.com/health
curl https://api.yourdomain.com/health
```

### Monitoring Setup

```bash
# Start with monitoring
docker-compose -f docker-compose.prod.yml --profile monitoring up -d

# Access monitoring
# Prometheus: http://server:9090
# Grafana: http://server:3001 (admin/your-password)
```

## Docker Configurations

### Backend Dockerfile

The backend uses multi-stage builds:

1. **Builder stage**: Compiles Go application
2. **Development stage**: Includes development tools and live reload
3. **Production stage**: Minimal scratch-based image

Build targets:
```bash
# Development build
docker build --target development -t urlshortener-backend:dev backend/

# Production build  
docker build --target production -t urlshortener-backend:prod backend/
```

### Frontend Dockerfile

The frontend also uses multi-stage builds:

1. **Builder stage**: Installs dependencies and builds assets
2. **Development stage**: Runs Vite dev server
3. **Production stage**: Serves built assets with Nginx

Build targets:
```bash
# Development build
docker build --target development -t urlshortener-frontend:dev frontend/

# Production build
docker build --target production -t urlshortener-frontend:prod frontend/
```

### Database Configuration

PostgreSQL is optimized for production:

- **Memory**: Shared buffers, work memory, maintenance work memory
- **WAL**: Write-ahead logging optimization
- **Connections**: Connection limits and pooling
- **Performance**: Statistics, query planner settings

### Redis Configuration

Redis production configuration includes:

- **Memory management**: LRU eviction policy
- **Persistence**: RDB snapshots for durability
- **Security**: Password protection, command renaming
- **Performance**: Pipeline optimization, connection limits

## Environment Variables

### Development

Development uses simplified configuration:

```bash
# Database
DATABASE_URL=postgres://urlshortener:password123@postgres:5432/urlshortener?sslmode=disable

# Redis
REDIS_URL=redis://:redis123@redis:6379

# Security (weak for development)
JWT_SECRET=development-jwt-secret-change-in-production

# Features
RATE_LIMIT_ENABLED=false
LOG_LEVEL=debug
```

### Production

Production requires secure configuration:

```bash
# Strong passwords
POSTGRES_PASSWORD=random-64-character-password
REDIS_PASSWORD=random-64-character-password
JWT_SECRET=random-64-character-secret

# Security settings
BCRYPT_COST=12
RATE_LIMIT_ENABLED=true
LOG_LEVEL=info
```

### Environment Validation

Validate configuration before deployment:

```bash
# Backend validation
docker-compose exec backend go run scripts/validate-env.go

# Check required variables
docker-compose config
```

## Data Persistence

### Volume Management

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect urlshortener_postgres_data

# Backup volume
docker run --rm -v urlshortener_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres-backup.tar.gz -C /data .

# Restore volume
docker run --rm -v urlshortener_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres-backup.tar.gz -C /data
```

### Database Backups

Automated backup script:

```bash
#!/bin/bash
# backup.sh
BACKUP_DIR="/opt/urlshortener/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/postgres_backup_$DATE.sql"

mkdir -p "$BACKUP_DIR"

# Create backup
docker-compose exec -T postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB > "$BACKUP_FILE"

# Compress backup
gzip "$BACKUP_FILE"

# Remove old backups (keep last 30 days)
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

Schedule with cron:
```bash
# Add to crontab
0 2 * * * /opt/urlshortener/backup.sh >> /var/log/urlshortener-backup.log 2>&1
```

## Monitoring

### Health Checks

All services include health checks:

```bash
# Check service health
docker-compose ps

# Manual health check
curl http://localhost:8080/health
```

### Resource Monitoring

```bash
# Monitor resource usage
docker stats

# View container logs
docker-compose logs -f --tail=100

# Monitor specific service
docker-compose logs -f backend
```

### Prometheus Metrics

Configure monitoring:

```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'url-shortener-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'
    
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:5432']
    
  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']
```

## Security

### Container Security

Production containers implement security best practices:

- **Non-root users**: All services run as non-root
- **Read-only filesystems**: Production containers are read-only
- **No new privileges**: Security option prevents privilege escalation
- **Resource limits**: CPU and memory limits prevent resource exhaustion
- **Network isolation**: Backend services in isolated network

### Network Security

```bash
# View networks
docker network ls

# Inspect network
docker network inspect urlshortener_backend-network
```

### Secrets Management

For enhanced security in production:

```bash
# Use Docker secrets
echo "strong-password" | docker secret create postgres_password -
echo "jwt-secret-key" | docker secret create jwt_secret -

# Reference in compose file
services:
  backend:
    secrets:
      - postgres_password
      - jwt_secret
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/postgres_password
      JWT_SECRET_FILE: /run/secrets/jwt_secret

secrets:
  postgres_password:
    external: true
  jwt_secret:
    external: true
```

## Troubleshooting

### Common Issues

#### Service Won't Start

```bash
# Check logs
docker-compose logs service-name

# Check configuration
docker-compose config

# Validate environment
docker-compose exec backend env | grep -E "DATABASE|REDIS|JWT"
```

#### Database Connection Issues

```bash
# Test database connection
docker-compose exec postgres psql -U urlshortener -c "SELECT 1;"

# Check database logs
docker-compose logs postgres

# Verify credentials
docker-compose exec backend sh -c 'echo $DATABASE_URL'
```

#### Redis Connection Issues

```bash
# Test Redis connection
docker-compose exec redis redis-cli ping

# Check Redis logs
docker-compose logs redis

# Test authentication
docker-compose exec redis redis-cli -a redis123 ping
```

#### Performance Issues

```bash
# Monitor resource usage
docker stats

# Check service health
docker-compose ps

# View service metrics
curl http://localhost:8080/metrics
```

### Debugging

#### Enable Debug Mode

```bash
# Set debug environment
export LOG_LEVEL=debug
docker-compose up -d backend

# View debug logs
docker-compose logs -f backend
```

#### Container Shell Access

```bash
# Access backend container
docker-compose exec backend sh

# Access database
docker-compose exec postgres psql -U urlshortener

# Access Redis
docker-compose exec redis redis-cli -a redis123
```

### Log Analysis

```bash
# View all logs
docker-compose logs

# Follow logs with timestamps
docker-compose logs -f -t

# Filter logs by service
docker-compose logs backend | grep ERROR

# Export logs
docker-compose logs > application.log
```

### Recovery Procedures

#### Database Recovery

```bash
# Stop services
docker-compose down

# Restore from backup
docker run --rm -v urlshortener_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres-backup.tar.gz -C /data

# Start services
docker-compose up -d
```

#### Full System Recovery

```bash
# Stop all services
docker-compose down -v

# Pull latest images
docker-compose pull

# Rebuild and start
docker-compose up -d --build
```

This comprehensive Docker guide provides everything needed to deploy and maintain the URL Shortener service using Docker in both development and production environments.