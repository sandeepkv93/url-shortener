# Reverse Proxy Configurations for URL Shortener Service

This directory contains production-ready reverse proxy configurations for deploying the URL Shortener service behind various load balancers and reverse proxies.

## Overview

The URL Shortener service is designed to work behind reverse proxies and load balancers. The service includes built-in middleware to handle proxy headers correctly and extract real client IPs for analytics and rate limiting.

## Supported Configurations

### 1. Nginx Reverse Proxy (`nginx.conf`)

**Features:**
- SSL termination with modern TLS settings
- Rate limiting for API and redirect endpoints
- Compression and caching
- Security headers and attack protection
- Health check and metrics endpoint routing
- Load balancing between multiple backend instances

**Setup:**
```bash
# Copy configuration
sudo cp nginx.conf /etc/nginx/sites-available/url-shortener
sudo ln -s /etc/nginx/sites-available/url-shortener /etc/nginx/sites-enabled/

# Update SSL certificate paths and domain names
sudo nano /etc/nginx/sites-available/url-shortener

# Test and reload
sudo nginx -t
sudo systemctl reload nginx
```

**Required Modules:**
- `ssl`
- `http_realip_module`
- `http_limit_req_module`
- `http_limit_conn_module`
- `http_gzip_module`

### 2. Apache Reverse Proxy (`apache.conf`)

**Features:**
- SSL termination and security
- Load balancing with health checks
- Compression and caching
- Security headers and attack protection
- Balancer manager for monitoring

**Setup:**
```bash
# Enable required modules
sudo a2enmod ssl rewrite proxy proxy_http proxy_balancer lbmethod_byrequests headers deflate expires

# Copy configuration
sudo cp apache.conf /etc/apache2/sites-available/url-shortener.conf
sudo a2ensite url-shortener.conf

# Update SSL certificate paths and domain names
sudo nano /etc/apache2/sites-available/url-shortener.conf

# Test and reload
sudo apache2ctl configtest
sudo systemctl reload apache2
```

### 3. AWS Application Load Balancer (`aws-alb.yaml`)

**Features:**
- Auto-scaling target groups
- SSL termination with ACM certificates
- Health checks and monitoring
- WAF integration for security
- CloudWatch alarms

**Setup:**
```bash
# Deploy using CloudFormation
aws cloudformation create-stack \
  --stack-name url-shortener-alb \
  --template-body file://aws-alb.yaml \
  --parameters \
    ParameterKey=VpcId,ParameterValue=vpc-xxxxxxxx \
    ParameterKey=SubnetIds,ParameterValue="subnet-xxxxxxxx,subnet-yyyyyyyy" \
    ParameterKey=CertificateArn,ParameterValue=arn:aws:acm:region:account:certificate/xxxxxxxx \
    ParameterKey=DomainName,ParameterValue=yourdomain.com \
  --capabilities CAPABILITY_IAM
```

### 4. Google Cloud Load Balancer (`gcp-lb.yaml`)

**Features:**
- Global load balancing with Cloud CDN
- Managed SSL certificates
- Cloud Armor security policies
- Auto-scaling instance groups
- Cloud Monitoring integration

**Setup:**
```bash
# Initialize Terraform
terraform init

# Plan deployment
terraform plan \
  -var="project_id=your-project-id" \
  -var="domain_name=yourdomain.com" \
  -var="environment=production"

# Deploy
terraform apply
```

## Backend Service Configuration

### Environment Variables

Configure the following environment variables for your URL Shortener service when running behind a reverse proxy:

```bash
# Enable proxy trust
TRUST_PROXY=true

# Proxy headers (defaults shown)
PROXY_HEADER=X-Forwarded-For
REAL_IP_HEADER=X-Real-IP
FORWARDED_PROTO_HEADER=X-Forwarded-Proto

# Security headers
ENABLE_SECURITY_HEADERS=true
ENABLE_SECURE_COOKIES=true  # For HTTPS

# Health check endpoints
ENABLE_HEALTH_ENDPOINTS=true
```

### Docker Configuration

When running in containers, ensure the service binds to all interfaces:

```bash
# Bind to all interfaces
HOST=0.0.0.0
PORT=8080

# Or in Docker
docker run -p 8080:8080 -e HOST=0.0.0.0 -e PORT=8080 url-shortener
```

## Security Considerations

### 1. Trusted Proxy Configuration

The service automatically trusts common proxy IP ranges:
- Private networks (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
- Loopback (127.0.0.0/8, ::1/128)
- Common cloud provider ranges

For custom proxy configurations, ensure your reverse proxy IPs are in the trusted ranges.

### 2. SSL/TLS Configuration

- Always use modern TLS versions (1.2+)
- Use strong cipher suites
- Enable HSTS headers
- Consider certificate pinning for high-security environments

### 3. Rate Limiting

Configure rate limiting at both the reverse proxy and application levels:
- **API endpoints**: 10-50 requests/second per IP
- **Redirect endpoints**: 50-100 requests/second per IP
- **Admin endpoints**: 5-10 requests/second per IP

### 4. Monitoring and Logging

Enable comprehensive logging:
- Access logs with real client IPs
- Error logs with detailed information
- Security event logs
- Performance metrics

## Health Checks

The service provides multiple health check endpoints:

- `/health` - Basic health check
- `/healthz` - Kubernetes-style health check
- `/ready` - Readiness check
- `/live` - Liveness check

Configure your load balancer to use these endpoints for health monitoring.

## Load Balancing Strategies

### Session Persistence

The URL Shortener service is stateless and doesn't require session persistence. Use round-robin or least-connections load balancing for optimal performance.

### Backend Instance Configuration

- **Minimum instances**: 2 (for high availability)
- **Maximum instances**: Scale based on load
- **Health check interval**: 30 seconds
- **Health check timeout**: 5 seconds
- **Unhealthy threshold**: 3 consecutive failures

## Troubleshooting

### Common Issues

1. **Real IP not detected**
   - Verify `TRUST_PROXY=true` is set
   - Check proxy headers are being forwarded
   - Ensure proxy IP is in trusted ranges

2. **SSL redirects not working**
   - Verify `X-Forwarded-Proto` header is set
   - Check `ENABLE_HTTPS=true` in production

3. **Rate limiting issues**
   - Check if real client IP is being detected
   - Verify rate limiting configuration
   - Monitor rate limiting logs

### Debug Headers

In development mode, the service adds debug headers:
- `X-Real-IP-Detected`: Detected real client IP
- `X-HTTPS-Detected`: Whether HTTPS was detected
- `X-Proxy-Trust-Enabled`: Whether proxy trust is enabled

## Performance Optimization

### Backend Optimization

- Enable compression at reverse proxy level
- Configure appropriate cache headers
- Use connection keep-alive
- Optimize timeout settings

### Database and Cache

- Use connection pooling
- Configure Redis for session storage
- Implement proper database indexes
- Monitor query performance

### Monitoring

Set up monitoring for:
- Response times and throughput
- Error rates (4xx, 5xx)
- Backend health and availability
- SSL certificate expiration

## Example Production Setup

```bash
# Nginx upstream configuration
upstream url_shortener_backend {
    least_conn;
    server 10.0.1.10:8080 max_fails=3 fail_timeout=30s;
    server 10.0.1.11:8080 max_fails=3 fail_timeout=30s;
    server 10.0.1.12:8080 max_fails=3 fail_timeout=30s backup;
}

# Backend service configuration
docker run -d \
  --name url-shortener-1 \
  -p 8080:8080 \
  -e GO_ENV=production \
  -e HOST=0.0.0.0 \
  -e PORT=8080 \
  -e TRUST_PROXY=true \
  -e ENABLE_SECURITY_HEADERS=true \
  -e DATABASE_URL=postgres://... \
  -e REDIS_URL=redis://... \
  url-shortener:latest
```

This setup provides a robust, scalable, and secure deployment of the URL Shortener service behind a reverse proxy.