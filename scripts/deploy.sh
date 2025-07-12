#!/bin/bash

# Production Deployment Script for URL Shortener
# This script automates the deployment process with safety checks and rollback capabilities

set -euo pipefail  # Exit on error, undefined variables, and pipe failures

# Script Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DEPLOY_USER="${DEPLOY_USER:-urlshortener}"
DEPLOY_HOST="${DEPLOY_HOST:-}"
DEPLOY_PATH="${DEPLOY_PATH:-/opt/urlshortener}"
BACKUP_PATH="${BACKUP_PATH:-/opt/urlshortener/backups}"
LOG_FILE="${LOG_FILE:-/var/log/urlshortener-deploy.log}"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.prod.yml}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

# Help function
show_help() {
    cat << EOF
URL Shortener Deployment Script

Usage: $0 [OPTIONS] COMMAND

Commands:
    deploy          Deploy application to production
    rollback        Rollback to previous version
    status          Check deployment status
    logs            View application logs
    backup          Create backup of data and configuration
    restore         Restore from backup
    health-check    Run health checks on deployed application

Options:
    -h, --help      Show this help message
    -e, --env       Environment file to use (default: .env)
    -f, --force     Force deployment without confirmations
    -d, --dry-run   Show what would be done without executing
    -v, --verbose   Enable verbose output
    --host HOST     Deployment host (required for remote deployment)
    --user USER     Deployment user (default: urlshortener)
    --path PATH     Deployment path (default: /opt/urlshortener)

Environment Variables:
    DEPLOY_HOST     Target deployment host
    DEPLOY_USER     Deployment user
    DEPLOY_PATH     Deployment directory
    BACKUP_PATH     Backup directory
    LOG_FILE        Deployment log file

Examples:
    # Local deployment
    $0 deploy

    # Remote deployment
    $0 --host server.example.com deploy

    # Rollback deployment
    $0 rollback

    # Check status
    $0 status

    # Create backup
    $0 backup

    # Health check
    $0 health-check
EOF
}

# Parse command line arguments
FORCE=false
DRY_RUN=false
VERBOSE=false
ENV_FILE=".env"
COMMAND=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -e|--env)
            ENV_FILE="$2"
            shift 2
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        --host)
            DEPLOY_HOST="$2"
            shift 2
            ;;
        --user)
            DEPLOY_USER="$2"
            shift 2
            ;;
        --path)
            DEPLOY_PATH="$2"
            shift 2
            ;;
        deploy|rollback|status|logs|backup|restore|health-check)
            COMMAND="$1"
            shift
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Validate command
if [[ -z "$COMMAND" ]]; then
    error "No command specified. Use --help for usage information."
fi

# Enable verbose mode if requested
if [[ "$VERBOSE" == true ]]; then
    set -x
fi

# Helper functions
run_command() {
    local cmd="$1"
    local desc="${2:-Running command}"
    
    if [[ "$DRY_RUN" == true ]]; then
        log "[DRY RUN] Would run: $cmd"
        return 0
    fi
    
    log "$desc: $cmd"
    if ! eval "$cmd"; then
        error "Command failed: $cmd"
    fi
}

remote_command() {
    local cmd="$1"
    local desc="${2:-Running remote command}"
    
    if [[ -z "$DEPLOY_HOST" ]]; then
        run_command "$cmd" "$desc"
    else
        run_command "ssh ${DEPLOY_USER}@${DEPLOY_HOST} '$cmd'" "$desc (remote)"
    fi
}

confirm() {
    local message="$1"
    
    if [[ "$FORCE" == true ]]; then
        return 0
    fi
    
    echo -n "$message (y/N): "
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        return 0
    else
        return 1
    fi
}

# Pre-deployment checks
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed"
    fi
    
    # Check environment file
    if [[ ! -f "$ENV_FILE" ]]; then
        error "Environment file not found: $ENV_FILE"
    fi
    
    # Validate environment variables
    source "$ENV_FILE"
    
    local required_vars=(
        "POSTGRES_PASSWORD"
        "REDIS_PASSWORD"
        "JWT_SECRET"
        "BASE_URL"
        "FRONTEND_URL"
    )
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var:-}" ]]; then
            error "Required environment variable not set: $var"
        fi
    done
    
    # Check for default/weak passwords
    if [[ "$POSTGRES_PASSWORD" == *"CHANGE_THIS"* ]] || [[ "$REDIS_PASSWORD" == *"CHANGE_THIS"* ]]; then
        error "Default passwords detected. Please update your environment file."
    fi
    
    success "Prerequisites check passed"
}

# Backup functions
create_backup() {
    log "Creating backup..."
    
    local backup_dir="$BACKUP_PATH/$(date +%Y%m%d_%H%M%S)"
    
    if [[ -z "$DEPLOY_HOST" ]]; then
        # Local backup
        mkdir -p "$backup_dir"
        
        # Backup database
        if docker-compose -f "$COMPOSE_FILE" ps postgres | grep -q Up; then
            run_command "docker-compose -f $COMPOSE_FILE exec -T postgres pg_dump -U \$POSTGRES_USER \$POSTGRES_DB > $backup_dir/database.sql" "Backing up database"
        fi
        
        # Backup Redis
        if docker-compose -f "$COMPOSE_FILE" ps redis | grep -q Up; then
            run_command "docker-compose -f $COMPOSE_FILE exec -T redis redis-cli --rdb - > $backup_dir/redis.rdb" "Backing up Redis"
        fi
        
        # Backup configuration
        cp "$ENV_FILE" "$backup_dir/"
        cp -r ssl "$backup_dir/" 2>/dev/null || true
        
        # Compress backup
        run_command "tar -czf $backup_dir.tar.gz -C $(dirname $backup_dir) $(basename $backup_dir)" "Compressing backup"
        run_command "rm -rf $backup_dir" "Cleaning up temporary files"
        
        success "Backup created: $backup_dir.tar.gz"
        echo "$backup_dir.tar.gz"
    else
        # Remote backup
        remote_command "mkdir -p $backup_dir"
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE exec -T postgres pg_dump -U \\\$POSTGRES_USER \\\$POSTGRES_DB > $backup_dir/database.sql"
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE exec -T redis redis-cli --rdb - > $backup_dir/redis.rdb"
        remote_command "cp $DEPLOY_PATH/$ENV_FILE $backup_dir/"
        remote_command "tar -czf $backup_dir.tar.gz -C $(dirname $backup_dir) $(basename $backup_dir) && rm -rf $backup_dir"
        
        success "Remote backup created: $backup_dir.tar.gz"
    fi
}

# Deployment functions
deploy_application() {
    log "Starting deployment..."
    
    # Pre-deployment checks
    check_prerequisites
    
    # Confirm deployment
    if ! confirm "Deploy to ${DEPLOY_HOST:-localhost}?"; then
        log "Deployment cancelled"
        exit 0
    fi
    
    # Create backup before deployment
    local backup_file
    backup_file=$(create_backup)
    
    if [[ -n "$DEPLOY_HOST" ]]; then
        # Remote deployment
        deploy_remote "$backup_file"
    else
        # Local deployment
        deploy_local "$backup_file"
    fi
    
    # Post-deployment verification
    verify_deployment
    
    success "Deployment completed successfully"
}

deploy_local() {
    local backup_file="$1"
    
    log "Deploying locally..."
    
    # Build images
    run_command "docker-compose -f $COMPOSE_FILE build --no-cache" "Building Docker images"
    
    # Stop existing services
    run_command "docker-compose -f $COMPOSE_FILE down" "Stopping existing services"
    
    # Start new services
    run_command "docker-compose -f $COMPOSE_FILE up -d" "Starting new services"
    
    # Wait for services to be ready
    wait_for_services
}

deploy_remote() {
    local backup_file="$1"
    
    log "Deploying to remote host: $DEPLOY_HOST"
    
    # Sync code to remote
    run_command "rsync -avz --exclude='.git' --exclude='node_modules' --exclude='data' . ${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}/" "Syncing code to remote"
    
    # Deploy on remote
    remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE build --no-cache" "Building Docker images (remote)"
    remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE down" "Stopping existing services (remote)"
    remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE up -d" "Starting new services (remote)"
    
    # Wait for remote services
    wait_for_remote_services
}

wait_for_services() {
    log "Waiting for services to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if docker-compose -f "$COMPOSE_FILE" ps | grep -q "Up (healthy)"; then
            success "Services are ready"
            return 0
        fi
        
        log "Attempt $attempt/$max_attempts: Services not ready yet, waiting..."
        sleep 10
        ((attempt++))
    done
    
    error "Services failed to become ready within expected time"
}

wait_for_remote_services() {
    log "Waiting for remote services to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE ps | grep -q 'Up (healthy)'" "Checking remote service status" 2>/dev/null; then
            success "Remote services are ready"
            return 0
        fi
        
        log "Attempt $attempt/$max_attempts: Remote services not ready yet, waiting..."
        sleep 10
        ((attempt++))
    done
    
    error "Remote services failed to become ready within expected time"
}

# Health check functions
health_check() {
    log "Running health checks..."
    
    local base_url
    if [[ -n "$DEPLOY_HOST" ]]; then
        base_url="https://$DEPLOY_HOST"
    else
        base_url="http://localhost:8080"
    fi
    
    # Check backend health
    if curl -sf "$base_url/health" > /dev/null; then
        success "Backend health check passed"
    else
        error "Backend health check failed"
    fi
    
    # Check frontend (if available)
    if curl -sf "${base_url%:*}:3000/health" > /dev/null 2>&1; then
        success "Frontend health check passed"
    else
        warning "Frontend health check failed or not available"
    fi
    
    # Check database connectivity
    if [[ -z "$DEPLOY_HOST" ]]; then
        if docker-compose -f "$COMPOSE_FILE" exec -T postgres pg_isready > /dev/null; then
            success "Database connectivity check passed"
        else
            error "Database connectivity check failed"
        fi
    else
        if remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE exec -T postgres pg_isready" "Checking remote database" > /dev/null; then
            success "Remote database connectivity check passed"
        else
            error "Remote database connectivity check failed"
        fi
    fi
    
    success "All health checks passed"
}

verify_deployment() {
    log "Verifying deployment..."
    
    # Wait a moment for services to stabilize
    sleep 5
    
    # Run health checks
    health_check
    
    # Check service versions/status
    if [[ -z "$DEPLOY_HOST" ]]; then
        run_command "docker-compose -f $COMPOSE_FILE ps" "Checking service status"
    else
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE ps" "Checking remote service status"
    fi
    
    success "Deployment verification completed"
}

# Rollback functions
rollback_deployment() {
    log "Starting rollback..."
    
    if ! confirm "Rollback deployment on ${DEPLOY_HOST:-localhost}?"; then
        log "Rollback cancelled"
        exit 0
    fi
    
    # Find latest backup
    local latest_backup
    if [[ -z "$DEPLOY_HOST" ]]; then
        latest_backup=$(ls -t "$BACKUP_PATH"/*.tar.gz 2>/dev/null | head -1)
    else
        latest_backup=$(remote_command "ls -t $BACKUP_PATH/*.tar.gz 2>/dev/null | head -1" "Finding latest backup")
    fi
    
    if [[ -z "$latest_backup" ]]; then
        error "No backup found for rollback"
    fi
    
    log "Rolling back to: $latest_backup"
    
    # Stop current services
    if [[ -z "$DEPLOY_HOST" ]]; then
        run_command "docker-compose -f $COMPOSE_FILE down" "Stopping current services"
    else
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE down" "Stopping remote services"
    fi
    
    # Restore from backup
    restore_from_backup "$latest_backup"
    
    # Start services
    if [[ -z "$DEPLOY_HOST" ]]; then
        run_command "docker-compose -f $COMPOSE_FILE up -d" "Starting restored services"
        wait_for_services
    else
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE up -d" "Starting remote restored services"
        wait_for_remote_services
    fi
    
    # Verify rollback
    verify_deployment
    
    success "Rollback completed successfully"
}

restore_from_backup() {
    local backup_file="$1"
    
    log "Restoring from backup: $backup_file"
    
    if [[ -z "$DEPLOY_HOST" ]]; then
        # Local restore
        local restore_dir="/tmp/restore_$(date +%s)"
        run_command "mkdir -p $restore_dir" "Creating restore directory"
        run_command "tar -xzf $backup_file -C $restore_dir" "Extracting backup"
        
        # Restore configuration
        local extracted_dir=$(ls "$restore_dir")
        run_command "cp $restore_dir/$extracted_dir/.env $ENV_FILE" "Restoring configuration"
        
        # Restore database and Redis data would require more complex volume management
        warning "Database and Redis restore requires manual intervention"
        
        run_command "rm -rf $restore_dir" "Cleaning up restore directory"
    else
        # Remote restore
        remote_command "cd /tmp && tar -xzf $backup_file" "Extracting remote backup"
        # Add remote restore logic here
        warning "Remote restore implementation needed"
    fi
}

# Status and monitoring functions
show_status() {
    log "Checking deployment status..."
    
    if [[ -z "$DEPLOY_HOST" ]]; then
        # Local status
        echo "=== Service Status ==="
        docker-compose -f "$COMPOSE_FILE" ps
        
        echo -e "\n=== Resource Usage ==="
        docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
        
    else
        # Remote status
        echo "=== Remote Service Status ==="
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE ps" "Getting remote status"
        
        echo -e "\n=== Remote Resource Usage ==="
        remote_command "docker stats --no-stream --format 'table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}'" "Getting remote resource usage"
    fi
    
    # Run health checks
    health_check
}

show_logs() {
    log "Showing application logs..."
    
    if [[ -z "$DEPLOY_HOST" ]]; then
        docker-compose -f "$COMPOSE_FILE" logs -f --tail=100
    else
        remote_command "cd $DEPLOY_PATH && docker-compose -f $COMPOSE_FILE logs --tail=100" "Getting remote logs"
    fi
}

# Main execution
main() {
    log "Starting deployment script - Command: $COMMAND"
    
    case "$COMMAND" in
        deploy)
            deploy_application
            ;;
        rollback)
            rollback_deployment
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs
            ;;
        backup)
            create_backup
            ;;
        restore)
            if [[ $# -eq 0 ]]; then
                error "Restore command requires backup file path"
            fi
            restore_from_backup "$1"
            ;;
        health-check)
            health_check
            ;;
        *)
            error "Unknown command: $COMMAND"
            ;;
    esac
}

# Create log directory if it doesn't exist
mkdir -p "$(dirname "$LOG_FILE")"

# Run main function
main "$@"