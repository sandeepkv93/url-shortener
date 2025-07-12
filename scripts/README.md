# URL Shortener Automation Scripts

This directory contains automation scripts for managing the URL Shortener service deployment, testing, backup, and maintenance.

## Scripts Overview

| Script | Purpose | Environment |
|--------|---------|-------------|
| `setup.sh` | Initial environment setup and configuration | Dev/Prod |
| `deploy.sh` | Production deployment with rollback capabilities | Prod |
| `test.sh` | Comprehensive testing automation | Dev/CI |
| `backup.sh` | Backup, restore, and maintenance operations | Prod |

## Quick Start

```bash
# Make scripts executable (if needed)
chmod +x scripts/*.sh

# Set up development environment
./scripts/setup.sh dev

# Run all tests
./scripts/test.sh all

# Deploy to production
./scripts/deploy.sh --host server.example.com deploy

# Create backup
./scripts/backup.sh backup
```

## Script Details

### setup.sh - Environment Setup

Automates the initial setup process for development and production environments.

#### Usage

```bash
./scripts/setup.sh [OPTIONS] [ENVIRONMENT]
```

#### Examples

```bash
# Development setup
./scripts/setup.sh dev
./scripts/setup.sh --install-tools dev

# Production setup
./scripts/setup.sh prod

# Docker-only setup
./scripts/setup.sh docker

# Check prerequisites only
./scripts/setup.sh --check-only
```

#### Features

- **Dependency Detection**: Automatically detects and installs missing dependencies
- **Cross-Platform**: Supports Linux (Debian/RedHat), macOS, and Windows
- **Environment Files**: Creates and configures environment files
- **Development Tools**: Installs Go and Node.js development tools
- **Production Setup**: Configures systemd services, firewall, and log rotation
- **Docker Integration**: Sets up Docker-based development environment

#### Environment-Specific Actions

**Development (`dev`):**
- Installs Go and Node.js dependencies
- Sets up development environment files
- Configures Docker Compose for development
- Installs development tools (optional)

**Production (`prod`):**
- Creates systemd service files
- Configures firewall rules
- Sets up log rotation
- Configures SSL certificates
- Creates backup directories

**Docker (`docker`):**
- Builds Docker images
- Sets up Docker Compose configuration
- Configures container networking

### deploy.sh - Production Deployment

Handles production deployments with safety checks, backup, and rollback capabilities.

#### Usage

```bash
./scripts/deploy.sh [OPTIONS] COMMAND
```

#### Commands

- `deploy` - Deploy application to production
- `rollback` - Rollback to previous version
- `status` - Check deployment status
- `logs` - View application logs
- `backup` - Create backup before deployment
- `health-check` - Run health checks

#### Examples

```bash
# Local deployment
./scripts/deploy.sh deploy

# Remote deployment
./scripts/deploy.sh --host server.example.com deploy

# Rollback deployment
./scripts/deploy.sh rollback

# Check status
./scripts/deploy.sh status

# Force deployment without confirmation
./scripts/deploy.sh --force deploy
```

#### Features

- **Safety Checks**: Validates environment and dependencies before deployment
- **Automatic Backup**: Creates backup before deployment
- **Health Verification**: Verifies deployment health after completion
- **Rollback Support**: Automatic rollback on failure
- **Remote Deployment**: Supports deployment to remote servers
- **Logging**: Comprehensive logging of all operations
- **Dry Run**: Preview deployment actions without execution

#### Deployment Process

1. **Pre-deployment checks** - Validates prerequisites and environment
2. **Backup creation** - Creates backup of current state
3. **Application deployment** - Builds and deploys new version
4. **Health verification** - Verifies deployment success
5. **Cleanup** - Removes temporary files

### test.sh - Test Automation

Comprehensive testing automation for backend and frontend components.

#### Usage

```bash
./scripts/test.sh [OPTIONS] [COMMAND]
```

#### Commands

- `all` - Run all tests (default)
- `unit` - Run unit tests only
- `integration` - Run integration tests only
- `e2e` - Run end-to-end tests only
- `backend` - Run backend tests only
- `frontend` - Run frontend tests only
- `coverage` - Generate test coverage reports
- `performance` - Run performance tests
- `security` - Run security tests
- `lint` - Run linting checks
- `clean` - Clean test artifacts

#### Examples

```bash
# Run all tests
./scripts/test.sh

# Run unit tests with coverage
./scripts/test.sh unit --coverage

# Run frontend tests in watch mode
./scripts/test.sh frontend --watch

# Run tests in CI mode
./scripts/test.sh --ci

# Run tests in Docker
./scripts/test.sh --docker all
```

#### Features

- **Multi-Language Testing**: Supports Go and Node.js/TypeScript testing
- **Coverage Reports**: Generates comprehensive coverage reports
- **CI/CD Integration**: Optimized for continuous integration
- **Docker Support**: Can run tests in isolated Docker containers
- **Performance Testing**: Includes benchmark and load testing
- **Security Testing**: Runs security scans and vulnerability checks
- **Linting**: Code quality and style checks

#### Test Categories

**Backend Tests:**
- Unit tests with Go's testing framework
- Integration tests with real database
- Benchmark tests for performance
- Security scans with gosec

**Frontend Tests:**
- Unit tests with Vitest
- Component tests with React Testing Library
- E2E tests with Playwright
- Linting with ESLint and TypeScript

### backup.sh - Backup and Maintenance

Handles database backups, file backups, system monitoring, and maintenance tasks.

#### Usage

```bash
./scripts/backup.sh [OPTIONS] COMMAND
```

#### Commands

- `backup` - Create full backup (database + files)
- `restore` - Restore from backup
- `list` - List available backups
- `cleanup` - Clean old backups based on retention policy
- `verify` - Verify backup integrity
- `monitor` - System monitoring and health checks
- `rotate-logs` - Rotate application logs
- `maintenance` - Run full maintenance routine
- `schedule` - Set up automated backup schedule

#### Examples

```bash
# Create full backup
./scripts/backup.sh backup

# Database only backup
./scripts/backup.sh backup --db-only

# Restore from specific date
./scripts/backup.sh restore --backup-date 2024-01-15

# Clean old backups
./scripts/backup.sh cleanup --retention 7

# System monitoring
./scripts/backup.sh monitor

# Schedule automated backups
./scripts/backup.sh schedule
```

#### Features

- **Full System Backup**: Database, configuration files, SSL certificates
- **Incremental Backups**: Space-efficient incremental backup support
- **Compression**: Automatic backup compression with gzip
- **Encryption**: Optional GPG encryption for sensitive data
- **Cloud Storage**: S3 integration for offsite backups
- **Automated Cleanup**: Retention policy-based cleanup
- **System Monitoring**: Comprehensive health and resource monitoring
- **Log Rotation**: Automatic log rotation and cleanup

#### Backup Components

**Database:**
- PostgreSQL full dump
- Schema-only backup
- Data-only backup
- Table listing

**Files:**
- Configuration files (.env, docker-compose.yml)
- SSL certificates
- Nginx configuration
- Application logs
- Custom scripts
- Documentation

## Configuration

### Environment Variables

Scripts can be configured using environment variables:

```bash
# Deployment configuration
export DEPLOY_HOST="server.example.com"
export DEPLOY_USER="urlshortener"
export DEPLOY_PATH="/opt/urlshortener"

# Backup configuration
export BACKUP_DIR="/opt/urlshortener/backups"
export BACKUP_RETENTION_DAYS=30
export S3_BUCKET="my-backup-bucket"

# Logging configuration
export LOG_FILE="/var/log/urlshortener-deploy.log"
```

### Configuration Files

Scripts respect the following configuration files:

- `.env` - Main environment configuration
- `backend/.env` - Backend-specific configuration
- `frontend/.env` - Frontend-specific configuration
- `docker-compose.yml` - Docker Compose configuration
- `docker-compose.prod.yml` - Production Docker Compose

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test and Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run tests
        run: ./scripts/test.sh --ci all

  deploy:
    if: github.ref == 'refs/heads/main'
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to production
        run: ./scripts/deploy.sh --host ${{ secrets.DEPLOY_HOST }} deploy
        env:
          DEPLOY_USER: ${{ secrets.DEPLOY_USER }}
          SSH_PRIVATE_KEY: ${{ secrets.SSH_PRIVATE_KEY }}
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any
    
    stages {
        stage('Test') {
            steps {
                sh './scripts/test.sh --ci all'
            }
        }
        
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh './scripts/deploy.sh --host ${DEPLOY_HOST} deploy'
            }
        }
        
        stage('Backup') {
            steps {
                sh './scripts/backup.sh backup'
            }
        }
    }
    
    post {
        always {
            sh './scripts/test.sh clean'
        }
    }
}
```

## Production Deployment Workflow

### Initial Production Setup

1. **Server Preparation**:
   ```bash
   # On the target server
   ./scripts/setup.sh prod
   ```

2. **Environment Configuration**:
   ```bash
   # Copy and edit production environment
   cp .env.production .env
   # Update with production values
   ```

3. **SSL Certificate Setup**:
   ```bash
   # Configure SSL certificates
   sudo certbot --nginx -d yourdomain.com
   ```

4. **Initial Deployment**:
   ```bash
   ./scripts/deploy.sh --host server.example.com deploy
   ```

### Regular Deployment Process

1. **Run Tests Locally**:
   ```bash
   ./scripts/test.sh all
   ```

2. **Deploy to Production**:
   ```bash
   ./scripts/deploy.sh --host server.example.com deploy
   ```

3. **Verify Deployment**:
   ```bash
   ./scripts/deploy.sh --host server.example.com health-check
   ```

4. **Monitor Status**:
   ```bash
   ./scripts/deploy.sh --host server.example.com status
   ```

### Backup and Maintenance Schedule

Set up automated maintenance with cron jobs:

```bash
# Daily backup at 2 AM
0 2 * * * /opt/urlshortener/scripts/backup.sh backup --quiet

# Weekly cleanup on Sunday at 3 AM
0 3 * * 0 /opt/urlshortener/scripts/backup.sh cleanup --quiet

# Monthly full maintenance on first Sunday at 4 AM
0 4 1-7 * 0 /opt/urlshortener/scripts/backup.sh maintenance --quiet
```

## Troubleshooting

### Common Issues

#### Permission Errors
```bash
# Fix script permissions
chmod +x scripts/*.sh

# Fix data directory permissions
sudo chown -R urlshortener:urlshortener /opt/urlshortener/data
```

#### Docker Issues
```bash
# Check Docker service
sudo systemctl status docker

# Restart Docker
sudo systemctl restart docker

# Clean Docker resources
docker system prune -f
```

#### Database Connection Issues
```bash
# Check database container
docker-compose ps postgres

# View database logs
docker-compose logs postgres

# Test database connection
docker-compose exec postgres psql -U urlshortener -c "SELECT 1;"
```

#### Deployment Failures
```bash
# Check deployment logs
./scripts/deploy.sh logs

# Rollback to previous version
./scripts/deploy.sh rollback

# Check system status
./scripts/backup.sh monitor
```

### Debug Mode

Enable debug mode for detailed output:

```bash
# Enable verbose output
./scripts/deploy.sh --verbose deploy

# Enable dry-run mode
./scripts/deploy.sh --dry-run deploy

# Check prerequisites only
./scripts/setup.sh --check-only
```

### Log Analysis

Scripts generate comprehensive logs:

```bash
# View deployment logs
tail -f /var/log/urlshortener-deploy.log

# View backup logs
tail -f /var/log/urlshortener-backup.log

# View application logs
./scripts/deploy.sh logs
```

## Security Considerations

### Script Security

- Scripts validate input parameters and environment variables
- Sensitive operations require confirmation (unless `--force` is used)
- Backups can be encrypted with GPG
- Remote operations use SSH with key-based authentication

### Production Security

- Scripts create non-root users for application services
- Docker containers run with security restrictions
- Firewall rules are configured automatically
- SSL/TLS certificates are managed and renewed automatically

### Best Practices

1. **Use dedicated deployment user**: Never deploy as root
2. **Secure environment files**: Set proper permissions (600) on .env files
3. **Regular backups**: Automated daily backups with offsite storage
4. **Monitor logs**: Regular log analysis and alerting
5. **Update dependencies**: Keep system and application dependencies updated

## Contributing

### Adding New Scripts

1. Follow the existing script structure and conventions
2. Include comprehensive help text and examples
3. Implement proper error handling and logging
4. Add tests for script functionality
5. Update this README with script documentation

### Script Conventions

- Use `set -euo pipefail` for safety
- Implement colored output with consistent formatting
- Provide `--help` option with comprehensive documentation
- Support `--dry-run` for preview mode
- Include proper error handling and cleanup
- Use confirmation prompts for destructive operations

---

For more information, see the main project documentation in the `docs/` directory.