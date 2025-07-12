# Health Check System Documentation

The URL Shortener service includes a comprehensive health check system that provides multiple levels of health monitoring and diagnostics. This system is designed for production monitoring, Kubernetes deployments, and operational visibility.

## Health Check Endpoints

All health check endpoints are available under the `/health` path and return JSON responses with appropriate HTTP status codes.

### Basic Health Checks

#### `GET /health`
**Primary health endpoint** - Returns overall system health status.

**Response Codes:**
- `200 OK` - System is healthy
- `206 Partial Content` - System is degraded but functional
- `503 Service Unavailable` - System is unhealthy

**Response Format:**
```json
{
  "status": "healthy|degraded|unhealthy",
  "version": "1.0.0",
  "uptime": "72h30m15s",
  "timestamp": "2024-01-15T10:30:00Z",
  "components": {
    "database": {
      "status": "up|degraded|down",
      "message": "Database is healthy",
      "last_checked": "2024-01-15T10:30:00Z",
      "duration": "15ms",
      "metadata": {
        "open_connections": 5,
        "in_use": 2,
        "idle": 3,
        "max_open_connections": 25
      }
    },
    "cache": {
      "status": "up|degraded|down",
      "message": "Cache is healthy",
      "last_checked": "2024-01-15T10:30:00Z",
      "duration": "8ms",
      "metadata": {}
    },
    "system": {
      "status": "up|degraded|down",
      "message": "System resources are healthy",
      "last_checked": "2024-01-15T10:30:00Z",
      "duration": "2ms",
      "metadata": {
        "memory_usage_percent": 45.2,
        "memory_allocated": 134217728,
        "memory_total": 268435456,
        "goroutines": 25,
        "gc_cycles": 142
      }
    }
  },
  "checks": {}
}
```

#### `GET /health/health` (Alias)
Alternative endpoint that returns the same response as `/health`.

#### `GET /health/healthz` (Kubernetes Style)
Kubernetes-compatible health check endpoint that returns the same response as `/health`.

### Kubernetes Health Checks

#### `GET /health/livez`
**Liveness probe** - Indicates if the application is running and can receive traffic.

**Response Codes:**
- `200 OK` - Application is alive

**Response Format:**
```json
{
  "status": "alive",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Usage in Kubernetes:**
```yaml
livenessProbe:
  httpGet:
    path: /health/livez
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

#### `GET /health/readyz`
**Readiness probe** - Indicates if the application is ready to serve traffic.

**Response Codes:**
- `200 OK` - Application is ready
- `503 Service Unavailable` - Application is not ready

**Response Format:**
```json
{
  "status": "ready|not ready",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Usage in Kubernetes:**
```yaml
readinessProbe:
  httpGet:
    path: /health/readyz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3
```

### Detailed Health Information

#### `GET /health/checks`
**Detailed health checks** - Returns comprehensive health check results.

**Response Codes:**
- `200 OK` - All critical checks passed
- `206 Partial Content` - Some non-critical checks failed or degraded
- `503 Service Unavailable` - Critical checks failed

**Response Format:**
```json
{
  "checks": {
    "database_connectivity": {
      "name": "Database Connectivity",
      "status": "pass|warn|fail",
      "message": "Database is accessible",
      "critical": true,
      "duration": "15ms",
      "last_run": "2024-01-15T10:30:00Z",
      "next_run": "2024-01-15T10:30:30Z"
    },
    "database_performance": {
      "name": "Database Performance",
      "status": "pass|warn|fail",
      "message": "Database queries are performing well",
      "critical": false,
      "duration": "45ms",
      "last_run": "2024-01-15T10:30:00Z",
      "next_run": "2024-01-15T10:31:00Z"
    },
    "cache_connectivity": {
      "name": "Cache Connectivity",
      "status": "pass|warn|fail",
      "message": "Cache is accessible",
      "critical": true,
      "duration": "8ms",
      "last_run": "2024-01-15T10:30:00Z",
      "next_run": "2024-01-15T10:30:30Z"
    },
    "cache_performance": {
      "name": "Cache Performance",
      "status": "pass|warn|fail",
      "message": "Cache operations are performing well",
      "critical": false,
      "duration": "12ms",
      "last_run": "2024-01-15T10:30:00Z",
      "next_run": "2024-01-15T10:31:00Z"
    },
    "disk_space": {
      "name": "Disk Space",
      "status": "pass|warn|fail",
      "message": "Disk space is adequate",
      "critical": false,
      "duration": "1ms",
      "last_run": "2024-01-15T10:30:00Z",
      "next_run": "2024-01-15T10:35:00Z"
    },
    "memory_usage": {
      "name": "Memory Usage",
      "status": "pass|warn|fail",
      "message": "Memory usage is normal",
      "critical": false,
      "duration": "2ms",
      "last_run": "2024-01-15T10:30:00Z",
      "next_run": "2024-01-15T10:31:00Z"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Component-Specific Health Checks

#### `GET /health/database`
**Database health** - Returns database-specific health information.

**Response Codes:**
- `200 OK` - Database is healthy
- `206 Partial Content` - Database is degraded
- `503 Service Unavailable` - Database is down

**Response Format:**
```json
{
  "status": "up|degraded|down",
  "message": "Database is healthy",
  "last_checked": "2024-01-15T10:30:00Z",
  "duration": "15ms",
  "metadata": {
    "open_connections": 5,
    "in_use": 2,
    "idle": 3,
    "max_open_connections": 25
  }
}
```

#### `GET /health/cache`
**Cache health** - Returns cache-specific health information.

**Response Codes:**
- `200 OK` - Cache is healthy
- `206 Partial Content` - Cache is degraded
- `503 Service Unavailable` - Cache is down

**Response Format:**
```json
{
  "status": "up|degraded|down",
  "message": "Cache is healthy",
  "last_checked": "2024-01-15T10:30:00Z",
  "duration": "8ms",
  "metadata": {}
}
```

#### `GET /health/external`
**External services health** - Returns health status of external service dependencies.

**Response Codes:**
- `200 OK` - All external services are healthy
- `206 Partial Content` - Some external services are degraded
- `503 Service Unavailable` - Critical external services are down

**Response Format:**
```json
{
  "status": "healthy|degraded|unhealthy",
  "services": {},
  "count": 0
}
```

### System Metrics

#### `GET /health/metrics/system`
**System-level metrics** - Returns detailed system resource metrics.

**Response Codes:**
- `200 OK` - Always returns success

**Response Format:**
```json
{
  "cpu": {
    "usage": 0.0,
    "cores": 4
  },
  "memory": {
    "total": 268435456,
    "used": 134217728,
    "free": 134217728,
    "available": 134217728,
    "usage_percent": 50.0
  },
  "processes": {
    "total": 25,
    "running": 25
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### `GET /health/metrics/application`
**Application-level metrics** - Returns application-specific metrics.

**Response Codes:**
- `200 OK` - Always returns success

**Response Format:**
```json
{
  "goroutine_count": 25,
  "heap_size": 134217728,
  "database_connections": {
    "open_connections": 5,
    "in_use_connections": 2,
    "idle_connections": 3,
    "max_open_connections": 25
  },
  "cache_metrics": {},
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Version Information

#### `GET /health/version`
**Version and build information** - Returns service version and build details.

**Response Codes:**
- `200 OK` - Always returns success

**Response Format:**
```json
{
  "service": "url-shortener",
  "version": "1.0.0",
  "build_time": "2024-01-15T10:00:00Z",
  "commit": "latest",
  "go_version": "1.21"
}
```

#### `GET /health/info` (Alias)
Alternative endpoint that returns the same response as `/health/version`.

## Health Check Configuration

### Service Integration

To integrate the health check system into your application:

```go
// Create health service
healthService := services.NewHealthService(db, cacheService, "1.0.0")

// Create health handler
healthHandler := handlers.NewHealthHandler(healthService)

// Configure router with health handler
router := routes.NewRouterBuilder().
    WithHealthHandler(healthHandler).
    // ... other handlers
    Build()
```

### Health Check Intervals

The system implements different check intervals based on criticality:

- **Database Connectivity**: Every 30 seconds
- **Cache Connectivity**: Every 30 seconds
- **Database Performance**: Every 60 seconds
- **Cache Performance**: Every 60 seconds
- **Memory Usage**: Every 60 seconds
- **Disk Space**: Every 5 minutes

### Thresholds and Alerts

#### Memory Usage
- **Warning**: > 90% usage
- **Critical**: > 95% usage

#### Database Connections
- **Warning**: > 90% of max connections used
- **Critical**: Unable to connect

#### Cache Performance
- **Warning**: Operations taking > 1 second
- **Critical**: Operations failing

#### Query Performance
- **Warning**: Queries taking > 5 seconds
- **Critical**: Query failures

## Monitoring Integration

### Prometheus Metrics

The health check system is designed to integrate with Prometheus monitoring:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'url-shortener'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/health/metrics/system'
    scrape_interval: 30s
```

### Grafana Dashboard

Example Grafana queries for health monitoring:

```promql
# Application uptime
up{job="url-shortener"}

# Memory usage percentage
memory_usage_percent{job="url-shortener"}

# Database connection pool usage
database_connections_used{job="url-shortener"} / database_connections_max{job="url-shortener"}

# Response time percentiles
histogram_quantile(0.95, http_request_duration_seconds_bucket{job="url-shortener"})
```

### Alerting Rules

Example Prometheus alerting rules:

```yaml
groups:
- name: url-shortener-health
  rules:
  - alert: ServiceDown
    expr: up{job="url-shortener"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "URL Shortener service is down"

  - alert: HighMemoryUsage
    expr: memory_usage_percent{job="url-shortener"} > 90
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage detected"

  - alert: DatabaseConnectionPoolExhausted
    expr: (database_connections_used{job="url-shortener"} / database_connections_max{job="url-shortener"}) > 0.9
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "Database connection pool nearly exhausted"
```

## Production Deployment

### Docker Health Checks

```dockerfile
# Dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD curl -f http://localhost:8080/health/livez || exit 1
```

### Docker Compose

```yaml
# docker-compose.yml
services:
  app:
    image: url-shortener:latest
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/readyz"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s
```

### Load Balancer Configuration

#### Nginx

```nginx
upstream backend {
    server app1:8080 max_fails=3 fail_timeout=30s;
    server app2:8080 max_fails=3 fail_timeout=30s;
}

server {
    location /health {
        access_log off;
        proxy_pass http://backend;
        proxy_set_header Host $host;
    }
    
    location / {
        proxy_pass http://backend;
        # Health check for upstream
        proxy_next_upstream error timeout http_503;
    }
}
```

#### HAProxy

```haproxy
backend url-shortener-backend
    balance roundrobin
    option httpchk GET /health/readyz
    http-check expect status 200
    server app1 app1:8080 check inter 30s
    server app2 app2:8080 check inter 30s
```

## Troubleshooting

### Common Health Check Issues

#### Database Health Check Failures

```bash
# Check database connectivity
docker-compose exec postgres psql -U urlshortener -c "SELECT 1;"

# View database logs
docker-compose logs postgres

# Check connection pool status
curl http://localhost:8080/health/database
```

#### Cache Health Check Failures

```bash
# Check Redis connectivity
docker-compose exec redis redis-cli ping

# View Redis logs
docker-compose logs redis

# Check cache status
curl http://localhost:8080/health/cache
```

#### Memory Issues

```bash
# Check detailed memory metrics
curl http://localhost:8080/health/metrics/system

# View application metrics
curl http://localhost:8080/health/metrics/application

# Monitor memory usage over time
watch -n 5 'curl -s http://localhost:8080/health/metrics/system | jq .memory'
```

### Debug Health Checks

Enable verbose health check logging:

```bash
# Set debug environment variable
export DEBUG_HEALTH_CHECKS=true

# Or use query parameter
curl "http://localhost:8080/health/checks?detailed=true"
```

### Health Check Testing

```bash
# Test all health endpoints
for endpoint in health livez readyz checks database cache external version; do
  echo "Testing /health/$endpoint"
  curl -w "Status: %{http_code}\n" "http://localhost:8080/health/$endpoint"
  echo "---"
done

# Load test health endpoints
ab -n 1000 -c 10 http://localhost:8080/health/livez
```

## Best Practices

### Health Check Design

1. **Keep liveness checks simple** - Only check if the application process is running
2. **Make readiness checks comprehensive** - Verify all dependencies are available
3. **Use appropriate timeouts** - Set reasonable timeouts for external dependencies
4. **Implement circuit breakers** - Prevent cascading failures from dependency checks
5. **Cache health check results** - Avoid overwhelming dependencies with health checks

### Monitoring Strategy

1. **Layer your monitoring** - Use liveness, readiness, and detailed health checks appropriately
2. **Set up alerts** - Monitor critical health check failures
3. **Regular health check validation** - Test health checks during deployments
4. **Document failure scenarios** - Create runbooks for common health check failures

### Performance Considerations

1. **Optimize health check frequency** - Balance monitoring needs with system load
2. **Use efficient health checks** - Prefer lightweight checks over heavy operations
3. **Implement health check caching** - Cache results for high-frequency checks
4. **Monitor health check performance** - Ensure health checks don't impact application performance

---

For more information about the URL Shortener service, see the main documentation in the [docs](.) directory.