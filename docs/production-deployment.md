# Production Deployment and Operations Guide

This guide provides comprehensive instructions for deploying, monitoring, and maintaining the URL Shortener service in production environments. It covers everything from initial server setup to ongoing operational procedures.

## Table of Contents

1. [Prerequisites and Planning](#prerequisites-and-planning)
2. [Infrastructure Setup](#infrastructure-setup)
3. [Initial Deployment](#initial-deployment)
4. [Configuration Management](#configuration-management)
5. [Security Configuration](#security-configuration)
6. [Monitoring and Observability](#monitoring-and-observability)
7. [Backup and Disaster Recovery](#backup-and-disaster-recovery)
8. [Operational Procedures](#operational-procedures)
9. [Scaling and Performance](#scaling-and-performance)
10. [Troubleshooting](#troubleshooting)
11. [Maintenance and Updates](#maintenance-and-updates)

## Prerequisites and Planning

### System Requirements

#### Minimum Production Environment
- **CPU**: 2 vCPU cores
- **RAM**: 4GB
- **Storage**: 20GB SSD
- **Network**: 1Gbps connection
- **OS**: Ubuntu 22.04 LTS or CentOS 8+

#### Recommended Production Environment
- **CPU**: 4 vCPU cores
- **RAM**: 8GB
- **Storage**: 50GB SSD
- **Network**: 1Gbps connection with redundancy
- **OS**: Ubuntu 22.04 LTS

#### High-Availability Setup
- **Load Balancer**: 2 instances
- **Application Servers**: 3+ instances
- **Database**: Primary + 2 replicas
- **Cache**: Redis cluster (3+ nodes)

### Required Services

#### External Services
- **Domain Name** - Registered domain with DNS control
- **SSL Certificate** - Let's Encrypt or commercial certificate
- **Email Service** - For notifications (optional)
- **Monitoring Service** - External monitoring (optional)

#### Cloud Provider Requirements
- **Compute Instances** - Virtual machines or containers
- **Load Balancer** - Application load balancer
- **Database** - Managed PostgreSQL or self-hosted
- **Cache** - Managed Redis or self-hosted
- **Storage** - Block storage for data persistence
- **Networking** - VPC with security groups

### Pre-Deployment Checklist

```bash
# Server access and credentials
□ SSH access to production servers
□ Root or sudo privileges
□ Database credentials and access
□ SSL certificates obtained
□ Domain DNS configured

# Security preparations
□ Firewall rules planned
□ Security groups configured
□ SSL/TLS certificates ready
□ Backup encryption keys generated

# Monitoring setup
□ Log aggregation service configured
□ Monitoring endpoints accessible
□ Alert notification channels set up
□ Performance baseline established

# Application configuration
□ Environment variables documented
□ Database migrations tested
□ External service integrations tested
□ Load testing completed
```

## Infrastructure Setup

### Single-Server Deployment

#### Server Preparation

```bash
# 1. Update system packages
sudo apt update && sudo apt upgrade -y

# 2. Install required packages
sudo apt install -y \
    curl \
    wget \
    unzip \
    git \
    nginx \
    certbot \
    python3-certbot-nginx \
    htop \
    iotop \
    netstat-nat

# 3. Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 4. Create application user
sudo useradd -m -s /bin/bash urlshortener
sudo mkdir -p /opt/urlshortener
sudo chown urlshortener:urlshortener /opt/urlshortener
```

#### Automated Setup Script

```bash
# Use the provided setup script
./scripts/setup.sh prod

# Or manual setup
sudo -u urlshortener bash << 'EOF'
cd /opt/urlshortener
git clone https://github.com/your-org/url-shortener.git .
cp .env.production .env
# Edit .env with production values
EOF
```

### Multi-Server Deployment

#### Load Balancer Setup (Nginx)

```nginx
# /etc/nginx/sites-available/urlshortener
upstream backend {
    least_conn;
    server app1.internal:8080 max_fails=3 fail_timeout=30s;
    server app2.internal:8080 max_fails=3 fail_timeout=30s;
    server app3.internal:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css application/json application/javascript text/javascript;

    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # API endpoints
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }

    # Short URL redirects
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static files (if serving frontend)
    location /static/ {
        alias /opt/urlshortener/frontend/dist/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

#### Database Setup (PostgreSQL)

```bash
# For managed database, configure connection settings
# For self-hosted PostgreSQL:

# Install PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Configure PostgreSQL
sudo -u postgres psql << 'EOF'
CREATE USER urlshortener WITH PASSWORD 'secure_password_here';
CREATE DATABASE urlshortener OWNER urlshortener;
GRANT ALL PRIVILEGES ON DATABASE urlshortener TO urlshortener;
\q
EOF

# Configure connection limits and performance
sudo tee -a /etc/postgresql/*/main/postgresql.conf << 'EOF'
# Performance tuning
shared_buffers = 1GB
effective_cache_size = 3GB
maintenance_work_mem = 256MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200

# Connection settings
max_connections = 200
listen_addresses = '*'
EOF

# Configure authentication
sudo tee -a /etc/postgresql/*/main/pg_hba.conf << 'EOF'
host    urlshortener    urlshortener    10.0.0.0/8       md5
EOF

sudo systemctl restart postgresql
```

#### Cache Setup (Redis)

```bash
# Install Redis
sudo apt install -y redis-server

# Configure Redis for production
sudo tee /etc/redis/redis.conf << 'EOF'
# Network configuration
bind 127.0.0.1 10.0.0.10
port 6379
protected-mode yes
requirepass your_redis_password_here

# Memory configuration
maxmemory 2gb
maxmemory-policy allkeys-lru

# Persistence
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dir /var/lib/redis

# Logging
loglevel notice
logfile /var/log/redis/redis-server.log

# Security
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command DEBUG ""
EOF

sudo systemctl restart redis-server
```

## Initial Deployment

### Using Deployment Scripts

The automated deployment script handles the entire deployment process:

```bash
# Deploy to production server
./scripts/deploy.sh --host your-server.com deploy

# Deploy with custom configuration
./scripts/deploy.sh \
    --host your-server.com \
    --user urlshortener \
    --path /opt/urlshortener \
    deploy

# Deploy with backup
./scripts/deploy.sh --host your-server.com --backup deploy
```

### Manual Deployment Steps

#### 1. Application Deployment

```bash
# Connect to production server
ssh urlshortener@your-server.com

# Navigate to application directory
cd /opt/urlshortener

# Pull latest code
git pull origin main

# Build and deploy with Docker Compose
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# Run database migrations
docker-compose -f docker-compose.prod.yml exec backend go run cmd/migrate/main.go

# Verify deployment
curl http://localhost:8080/health/checks
```

#### 2. SSL Certificate Setup

```bash
# Obtain SSL certificate with Let's Encrypt
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Test certificate renewal
sudo certbot renew --dry-run

# Set up automatic renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

#### 3. Firewall Configuration

```bash
# Configure UFW firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# For multi-server setup, allow internal communication
sudo ufw allow from 10.0.0.0/8 to any port 8080
sudo ufw allow from 10.0.0.0/8 to any port 5432
sudo ufw allow from 10.0.0.0/8 to any port 6379
```

#### 4. Service Configuration

```bash
# Create systemd service for automatic startup
sudo tee /etc/systemd/system/urlshortener.service << 'EOF'
[Unit]
Description=URL Shortener Service
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/urlshortener
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml down
TimeoutStartSec=0
User=urlshortener
Group=urlshortener

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl enable urlshortener.service
sudo systemctl start urlshortener.service
```

## Configuration Management

### Environment Variables

#### Production Environment File

```bash
# /opt/urlshortener/.env
# Application Configuration
APP_ENV=production
APP_DEBUG=false
APP_URL=https://yourdomain.com
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=urlshortener
DB_PASSWORD=secure_database_password
DB_NAME=urlshortener
DB_SSLMODE=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secure_redis_password
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10

# JWT Configuration
JWT_SECRET=very_secure_jwt_secret_key_here
JWT_EXPIRES_IN=24h
JWT_REFRESH_EXPIRES_IN=7d

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_GLOBAL=100
RATE_LIMIT_AUTH=5
RATE_LIMIT_URL_CREATION=10

# Security Configuration
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_CREDENTIALS=true

# External Services
GEOLOCATION_SERVICE_URL=https://ipapi.co
GEOLOCATION_API_KEY=your_geolocation_api_key

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout

# Monitoring Configuration
PROMETHEUS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s

# Backup Configuration
BACKUP_ENABLED=true
BACKUP_SCHEDULE=0 2 * * *
BACKUP_RETENTION_DAYS=30
S3_BACKUP_BUCKET=your-backup-bucket
S3_ACCESS_KEY=your_s3_access_key
S3_SECRET_KEY=your_s3_secret_key
```

### Secrets Management

#### Using Docker Secrets

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  backend:
    image: urlshortener/backend:latest
    secrets:
      - db_password
      - jwt_secret
      - redis_password
    environment:
      DB_PASSWORD_FILE: /run/secrets/db_password
      JWT_SECRET_FILE: /run/secrets/jwt_secret
      REDIS_PASSWORD_FILE: /run/secrets/redis_password

secrets:
  db_password:
    file: ./secrets/db_password.txt
  jwt_secret:
    file: ./secrets/jwt_secret.txt
  redis_password:
    file: ./secrets/redis_password.txt
```

#### Using External Secret Management

```bash
# Example with HashiCorp Vault
vault kv put secret/urlshortener \
    db_password="secure_password" \
    jwt_secret="secure_jwt_secret" \
    redis_password="secure_redis_password"

# Example with AWS Secrets Manager
aws secretsmanager create-secret \
    --name "urlshortener/production" \
    --description "URL Shortener production secrets" \
    --secret-string '{"db_password":"secure_password","jwt_secret":"secure_jwt_secret"}'
```

## Security Configuration

### Application Security

#### Security Headers

```go
// Already implemented in middleware, verify configuration
func SecurityMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Security headers
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            w.Header().Set("X-XSS-Protection", "1; mode=block")
            w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            w.Header().Set("Content-Security-Policy", "default-src 'self'")
            
            next.ServeHTTP(w, r)
        })
    }
}
```

#### Input Validation

```bash
# Ensure all input validation is enabled
grep -r "validate:" backend/internal/core/domain/

# Check for SQL injection prevention
grep -r "parameterized\|prepared" backend/internal/infrastructure/database/
```

### Network Security

#### Firewall Rules

```bash
# Production firewall configuration
sudo ufw reset
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH (restrict to specific IPs in production)
sudo ufw allow from YOUR_ADMIN_IP to any port 22

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow internal communication (adjust IP ranges)
sudo ufw allow from 10.0.0.0/8 to any port 8080
sudo ufw allow from 10.0.0.0/8 to any port 5432
sudo ufw allow from 10.0.0.0/8 to any port 6379

sudo ufw enable
```

#### SSL/TLS Configuration

```bash
# Test SSL configuration
curl -I https://yourdomain.com
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com

# Check SSL rating
curl -s "https://api.ssllabs.com/api/v3/analyze?host=yourdomain.com" | jq .
```

### Database Security

#### PostgreSQL Security

```sql
-- Create read-only user for analytics
CREATE USER urlshortener_readonly WITH PASSWORD 'readonly_password';
GRANT CONNECT ON DATABASE urlshortener TO urlshortener_readonly;
GRANT USAGE ON SCHEMA public TO urlshortener_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO urlshortener_readonly;

-- Ensure proper permissions
REVOKE ALL ON DATABASE urlshortener FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
```

#### Redis Security

```bash
# Verify Redis security configuration
redis-cli config get requirepass
redis-cli config get protected-mode
redis-cli config get bind
```

## Monitoring and Observability

### Health Monitoring

#### Application Health Checks

```bash
# Set up health check monitoring
cat > /etc/systemd/system/urlshortener-health.service << 'EOF'
[Unit]
Description=URL Shortener Health Check
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/curl -f http://localhost:8080/health/checks || /usr/bin/systemctl restart urlshortener
User=urlshortener

[Install]
WantedBy=multi-user.target
EOF

cat > /etc/systemd/system/urlshortener-health.timer << 'EOF'
[Unit]
Description=Run URL Shortener Health Check
Requires=urlshortener-health.service

[Timer]
OnCalendar=*:0/5
Persistent=true

[Install]
WantedBy=timers.target
EOF

sudo systemctl enable urlshortener-health.timer
sudo systemctl start urlshortener-health.timer
```

#### External Monitoring

```bash
# Example monitoring script
#!/bin/bash
# /opt/urlshortener/scripts/external-health-check.sh

ENDPOINT="https://yourdomain.com/health/checks"
TIMEOUT=10
RETRIES=3

for i in $(seq 1 $RETRIES); do
    if curl -f -m $TIMEOUT "$ENDPOINT" > /dev/null 2>&1; then
        echo "Health check passed"
        exit 0
    fi
    sleep 5
done

echo "Health check failed after $RETRIES attempts"
# Send alert notification here
exit 1
```

### Logging

#### Centralized Logging Setup

```yaml
# docker-compose.prod.yml - Add logging configuration
version: '3.8'
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    # Or use external logging driver
    # logging:
    #   driver: syslog
    #   options:
    #     syslog-address: "tcp://logserver:514"
```

#### Log Rotation

```bash
# Configure logrotate
sudo tee /etc/logrotate.d/urlshortener << 'EOF'
/var/log/urlshortener/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 urlshortener urlshortener
    postrotate
        docker-compose -f /opt/urlshortener/docker-compose.prod.yml restart backend
    endscript
}
EOF
```

### Metrics Collection

#### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'urlshortener'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/health/metrics/system'
    scrape_interval: 30s
    
  - job_name: 'urlshortener-health'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/health'
    scrape_interval: 10s
```

#### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "URL Shortener Metrics",
    "panels": [
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, http_request_duration_seconds_bucket{job=\"urlshortener\"})",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{job=\"urlshortener\"}[5m])",
            "legendFormat": "Requests/sec"
          }
        ]
      }
    ]
  }
}
```

### Alerting

#### Alert Rules

```yaml
# alert-rules.yml
groups:
  - name: urlshortener-alerts
    rules:
      - alert: ServiceDown
        expr: up{job="urlshortener"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "URL Shortener service is down"
          description: "The URL Shortener service has been down for more than 1 minute."

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, http_request_duration_seconds_bucket{job="urlshortener"}) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is above 1 second."

      - alert: HighErrorRate
        expr: rate(http_requests_total{job="urlshortener",status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is above 10% for the last 5 minutes."

      - alert: DatabaseConnectionFailure
        expr: urlshortener_database_health != 1
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connection failure"
          description: "Unable to connect to the database."
```

## Backup and Disaster Recovery

### Automated Backup Strategy

#### Database Backups

```bash
# Use the provided backup script
./scripts/backup.sh backup --db-only

# Or set up automated backups
crontab -e
# Add: 0 2 * * * /opt/urlshortener/scripts/backup.sh backup --quiet
```

#### Complete System Backup

```bash
# Full system backup including files and database
./scripts/backup.sh backup

# Schedule regular backups
cat > /etc/cron.d/urlshortener-backup << 'EOF'
# Daily backup at 2 AM
0 2 * * * urlshortener /opt/urlshortener/scripts/backup.sh backup --quiet

# Weekly cleanup on Sunday at 3 AM
0 3 * * 0 urlshortener /opt/urlshortener/scripts/backup.sh cleanup --quiet

# Monthly full maintenance on first Sunday at 4 AM
0 4 1-7 * 0 urlshortener /opt/urlshortener/scripts/backup.sh maintenance --quiet
EOF
```

#### Offsite Backup

```bash
# Configure S3 backup
export S3_BACKUP_BUCKET="your-backup-bucket"
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"

# Test S3 backup
./scripts/backup.sh backup --s3-upload
```

### Disaster Recovery

#### Recovery Procedures

```bash
# 1. Restore from backup
./scripts/backup.sh restore --backup-date 2024-01-15

# 2. Verify data integrity
docker-compose -f docker-compose.prod.yml exec backend go run cmd/verify/main.go

# 3. Test application functionality
curl http://localhost:8080/health/checks

# 4. Update DNS if needed (for disaster recovery site)
# Update DNS records to point to recovery site
```

#### Recovery Testing

```bash
# Regular recovery testing schedule
# Monthly: Test backup restore on staging environment
# Quarterly: Full disaster recovery simulation
# Annually: Complete infrastructure rebuild

# Example recovery test script
#!/bin/bash
BACKUP_DATE=$(date -d "yesterday" +%Y-%m-%d)
./scripts/backup.sh restore --backup-date $BACKUP_DATE --test-mode
./scripts/test.sh integration
```

## Operational Procedures

### Deployment Process

#### Rolling Deployment

```bash
# 1. Pre-deployment checks
./scripts/deploy.sh --host server1.com health-check
./scripts/test.sh all

# 2. Deploy to first server
./scripts/deploy.sh --host server1.com deploy

# 3. Verify deployment
./scripts/deploy.sh --host server1.com health-check

# 4. Deploy to remaining servers
for server in server2.com server3.com; do
    ./scripts/deploy.sh --host $server deploy
    ./scripts/deploy.sh --host $server health-check
done
```

#### Blue-Green Deployment

```bash
# 1. Deploy to green environment
./scripts/deploy.sh --host green.internal deploy

# 2. Run smoke tests
./scripts/test.sh e2e --host green.internal

# 3. Switch traffic
# Update load balancer configuration to point to green

# 4. Monitor for issues
./scripts/deploy.sh --host green.internal status

# 5. Decommission blue environment (after verification)
```

### Routine Maintenance

#### Weekly Maintenance

```bash
#!/bin/bash
# weekly-maintenance.sh

# 1. Update system packages
sudo apt update && sudo apt upgrade -y

# 2. Clean Docker resources
docker system prune -f

# 3. Analyze database performance
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U urlshortener -c "SELECT * FROM pg_stat_activity;"

# 4. Check log sizes
du -sh /var/log/urlshortener/*

# 5. Verify backups
./scripts/backup.sh list | tail -7

# 6. Run health checks
./scripts/deploy.sh health-check
```

#### Monthly Maintenance

```bash
#!/bin/bash
# monthly-maintenance.sh

# 1. Database maintenance
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U urlshortener -c "VACUUM ANALYZE;"

# 2. SSL certificate check
sudo certbot certificates

# 3. Security updates
sudo apt update && sudo apt list --upgradable
sudo unattended-upgrades

# 4. Performance review
./scripts/backup.sh monitor --report

# 5. Backup cleanup
./scripts/backup.sh cleanup --retention 30
```

### Performance Monitoring

#### Daily Performance Checks

```bash
# CPU and memory usage
htop
free -h
df -h

# Database performance
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U urlshortener -c "
    SELECT query, calls, total_time, mean_time 
    FROM pg_stat_statements 
    ORDER BY total_time DESC 
    LIMIT 10;"

# Cache performance
docker-compose -f docker-compose.prod.yml exec redis redis-cli info stats

# Application metrics
curl -s http://localhost:8080/health/metrics/application | jq
```

#### Performance Alerts

```bash
# Set up performance monitoring alerts
cat > /opt/urlshortener/scripts/performance-monitor.sh << 'EOF'
#!/bin/bash

# CPU usage alert
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
    echo "High CPU usage: $CPU_USAGE%" | mail -s "Performance Alert" admin@example.com
fi

# Memory usage alert
MEMORY_USAGE=$(free | grep Mem | awk '{printf "%.2f", $3/$2 * 100.0}')
if (( $(echo "$MEMORY_USAGE > 85" | bc -l) )); then
    echo "High memory usage: $MEMORY_USAGE%" | mail -s "Performance Alert" admin@example.com
fi

# Disk usage alert
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | cut -d'%' -f1)
if [ $DISK_USAGE -gt 85 ]; then
    echo "High disk usage: $DISK_USAGE%" | mail -s "Performance Alert" admin@example.com
fi
EOF

# Schedule performance monitoring
echo "*/5 * * * * /opt/urlshortener/scripts/performance-monitor.sh" | crontab -
```

## Scaling and Performance

### Horizontal Scaling

#### Adding Application Servers

```bash
# 1. Provision new server
# 2. Install application using setup script
./scripts/setup.sh prod

# 3. Deploy application
./scripts/deploy.sh --host new-server.com deploy

# 4. Add to load balancer configuration
# Update nginx upstream configuration

# 5. Test load distribution
for i in {1..10}; do
    curl -H "Host: yourdomain.com" http://new-server.com/health
done
```

#### Database Scaling

```sql
-- Set up read replicas
-- Primary database configuration
ALTER SYSTEM SET wal_level = replica;
ALTER SYSTEM SET max_wal_senders = 3;
ALTER SYSTEM SET wal_keep_segments = 64;
SELECT pg_reload_conf();

-- Create replication user
CREATE USER replicator REPLICATION LOGIN CONNECTION LIMIT 1 ENCRYPTED PASSWORD 'replicator_password';
```

#### Cache Scaling

```bash
# Redis cluster setup
redis-cli --cluster create \
    cache1.internal:6379 \
    cache2.internal:6379 \
    cache3.internal:6379 \
    --cluster-replicas 1
```

### Performance Optimization

#### Database Optimization

```sql
-- Add indexes for common queries
CREATE INDEX CONCURRENTLY idx_short_urls_user_id ON short_urls(user_id);
CREATE INDEX CONCURRENTLY idx_clicks_short_url_id ON clicks(short_url_id);
CREATE INDEX CONCURRENTLY idx_clicks_created_at ON clicks(created_at);

-- Analyze query performance
EXPLAIN ANALYZE SELECT * FROM short_urls WHERE user_id = 1 ORDER BY created_at DESC LIMIT 10;
```

#### Application Optimization

```bash
# Profile application performance
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Check for goroutine leaks
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

### Load Testing

#### Performance Testing

```bash
# Install load testing tools
sudo apt install -y apache2-utils wrk

# Basic load test
ab -n 10000 -c 100 http://localhost:8080/health

# Advanced load test with wrk
wrk -t12 -c400 -d30s --script=load-test.lua http://localhost:8080/

# Load test specific endpoints
cat > load-test.lua << 'EOF'
request = function()
    local path = "/api/v1/urls"
    local body = '{"original_url": "https://example.com/test-' .. math.random(1000000) .. '"}'
    local headers = {
        ["Content-Type"] = "application/json",
        ["Authorization"] = "Bearer YOUR_TEST_TOKEN"
    }
    return wrk.format("POST", path, headers, body)
end
EOF
```

## Troubleshooting

### Common Issues

#### Application Not Starting

```bash
# Check service status
sudo systemctl status urlshortener

# Check logs
docker-compose -f docker-compose.prod.yml logs backend

# Check configuration
docker-compose -f docker-compose.prod.yml config

# Verify environment variables
docker-compose -f docker-compose.prod.yml exec backend env | grep -E "(DB_|REDIS_|JWT_)"
```

#### Database Connection Issues

```bash
# Test database connectivity
docker-compose -f docker-compose.prod.yml exec postgres psql -U urlshortener -c "SELECT 1;"

# Check database logs
docker-compose -f docker-compose.prod.yml logs postgres

# Verify database configuration
psql -h localhost -p 5432 -U urlshortener -c "\l"
```

#### High CPU/Memory Usage

```bash
# Identify resource-intensive processes
htop
ps aux --sort=-%cpu | head -10
ps aux --sort=-%mem | head -10

# Check application metrics
curl http://localhost:8080/health/metrics/system

# Profile application
go tool pprof http://localhost:8080/debug/pprof/profile
```

#### SSL Certificate Issues

```bash
# Check certificate status
sudo certbot certificates

# Test SSL configuration
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com

# Renew certificate manually
sudo certbot renew --force-renewal
```

### Emergency Procedures

#### Service Recovery

```bash
# 1. Immediate service restart
sudo systemctl restart urlshortener

# 2. If restart fails, rollback to previous version
./scripts/deploy.sh rollback

# 3. Check health after rollback
curl http://localhost:8080/health/checks

# 4. If still failing, restore from backup
./scripts/backup.sh restore --latest
```

#### Database Recovery

```bash
# 1. Stop application to prevent further writes
sudo systemctl stop urlshortener

# 2. Restore database from backup
./scripts/backup.sh restore --db-only --backup-date $(date -d "yesterday" +%Y-%m-%d)

# 3. Verify data integrity
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U urlshortener -c "SELECT COUNT(*) FROM short_urls;"

# 4. Restart application
sudo systemctl start urlshortener
```

### Log Analysis

#### Application Logs

```bash
# View real-time logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Search for errors
docker-compose -f docker-compose.prod.yml logs backend | grep -i error

# Analyze access patterns
grep "POST /api/v1/urls" /var/log/nginx/access.log | wc -l

# Monitor specific error codes
tail -f /var/log/nginx/access.log | grep " 5[0-9][0-9] "
```

#### Performance Analysis

```bash
# Database query analysis
docker-compose -f docker-compose.prod.yml exec postgres \
    psql -U urlshortener -c "
    SELECT query, calls, total_time/calls as avg_time 
    FROM pg_stat_statements 
    WHERE calls > 100 
    ORDER BY avg_time DESC 
    LIMIT 10;"

# Response time analysis
awk '{print $10}' /var/log/nginx/access.log | sort -n | tail -100
```

## Maintenance and Updates

### Regular Updates

#### Application Updates

```bash
# 1. Test updates in staging environment
git checkout staging
./scripts/test.sh all

# 2. Deploy to production with backup
./scripts/deploy.sh --backup deploy

# 3. Verify deployment
./scripts/deploy.sh health-check

# 4. Monitor for issues
tail -f /var/log/nginx/access.log
```

#### System Updates

```bash
# Weekly system updates
sudo apt update && sudo apt list --upgradable
sudo apt upgrade -y

# Docker updates
sudo apt update docker-ce docker-ce-cli containerd.io
sudo systemctl restart docker

# Restart services if needed
sudo systemctl restart urlshortener
```

#### Security Updates

```bash
# Enable automatic security updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure unattended-upgrades

# Configure automatic updates
sudo tee /etc/apt/apt.conf.d/50unattended-upgrades << 'EOF'
Unattended-Upgrade::Allowed-Origins {
    "${distro_id}:${distro_codename}-security";
    "${distro_id}ESMApps:${distro_codename}-apps-security";
    "${distro_id}ESM:${distro_codename}-infra-security";
};
Unattended-Upgrade::AutoFixInterruptedDpkg "true";
Unattended-Upgrade::MinimalSteps "true";
Unattended-Upgrade::Remove-Unused-Dependencies "true";
Unattended-Upgrade::Automatic-Reboot "false";
EOF
```

### Database Maintenance

#### Regular Maintenance Tasks

```sql
-- Weekly maintenance
VACUUM ANALYZE;

-- Monthly maintenance
REINDEX DATABASE urlshortener;

-- Check for long-running queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query 
FROM pg_stat_activity 
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';
```

#### Performance Monitoring

```bash
# Monitor database performance
cat > /opt/urlshortener/scripts/db-monitor.sh << 'EOF'
#!/bin/bash
docker-compose -f /opt/urlshortener/docker-compose.prod.yml exec postgres \
    psql -U urlshortener -c "
    SELECT 
        schemaname,
        tablename,
        attname,
        n_distinct,
        correlation
    FROM pg_stats
    WHERE schemaname = 'public'
    ORDER BY tablename, attname;"
EOF
```

### Capacity Planning

#### Resource Monitoring

```bash
# Set up capacity monitoring
cat > /opt/urlshortener/scripts/capacity-monitor.sh << 'EOF'
#!/bin/bash

# Database size monitoring
DB_SIZE=$(docker-compose -f /opt/urlshortener/docker-compose.prod.yml exec postgres \
    psql -U urlshortener -t -c "SELECT pg_size_pretty(pg_database_size('urlshortener'));")

# URL count monitoring
URL_COUNT=$(docker-compose -f /opt/urlshortener/docker-compose.prod.yml exec postgres \
    psql -U urlshortener -t -c "SELECT COUNT(*) FROM short_urls;")

# Click count monitoring
CLICK_COUNT=$(docker-compose -f /opt/urlshortener/docker-compose.prod.yml exec postgres \
    psql -U urlshortener -t -c "SELECT COUNT(*) FROM clicks;")

# Log metrics
echo "$(date): DB_SIZE=$DB_SIZE, URLs=$URL_COUNT, CLICKS=$CLICK_COUNT" >> /var/log/urlshortener/capacity.log

# Alert if growth rate is high
# Add alerting logic here
EOF

# Schedule capacity monitoring
echo "0 */6 * * * /opt/urlshortener/scripts/capacity-monitor.sh" | crontab -
```

### Documentation Updates

#### Keeping Documentation Current

```bash
# Update configuration documentation
./scripts/deploy.sh --host server.com --dry-run deploy > deployment-log.txt

# Update performance baselines
./scripts/test.sh performance --baseline > performance-baseline.txt

# Update operational procedures
git log --oneline --since="1 month ago" -- docs/ > recent-doc-changes.txt
```

---

This production deployment and operations guide provides comprehensive coverage of deploying, monitoring, and maintaining the URL Shortener service in production environments. For additional support, refer to the [troubleshooting documentation](./troubleshooting.md) and [developer onboarding guide](./developer-onboarding.md).