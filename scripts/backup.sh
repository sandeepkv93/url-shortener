#!/bin/bash

# Backup and Maintenance Script for URL Shortener
# Handles database backups, log rotation, and system maintenance

set -euo pipefail

# Script Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="${BACKUP_DIR:-$PROJECT_ROOT/backups}"
LOG_FILE="${LOG_FILE:-/var/log/urlshortener-backup.log}"

# Default configuration
DEFAULT_RETENTION_DAYS=30
DEFAULT_S3_BUCKET=""
DEFAULT_BACKUP_SCHEDULE="daily"

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
URL Shortener Backup and Maintenance Script

Usage: $0 [OPTIONS] COMMAND

Commands:
    backup          Create full backup (database + files)
    restore         Restore from backup
    list            List available backups
    cleanup         Clean old backups based on retention policy
    verify          Verify backup integrity
    monitor         System monitoring and health checks
    rotate-logs     Rotate application logs
    maintenance     Run full maintenance routine
    schedule        Set up automated backup schedule

Options:
    -h, --help              Show this help message
    -d, --backup-dir DIR    Backup directory (default: ./backups)
    -r, --retention DAYS    Backup retention in days (default: 30)
    -c, --compress          Compress backups (default: enabled)
    -e, --encrypt           Encrypt backups (requires GPG)
    -s, --s3-bucket BUCKET  Upload to S3 bucket
    -f, --force             Force operation without confirmation
    -v, --verbose           Enable verbose output
    -q, --quiet             Suppress non-essential output
    --dry-run               Show what would be done without executing

Backup Commands:
    --db-only               Backup database only
    --files-only            Backup files only
    --incremental           Create incremental backup

Restore Commands:
    --backup-file FILE      Specific backup file to restore
    --backup-date DATE      Restore from specific date (YYYY-MM-DD)
    --db-only               Restore database only
    --files-only            Restore files only

Examples:
    $0 backup                           # Full backup
    $0 backup --db-only                 # Database only
    $0 restore --backup-date 2024-01-15 # Restore from date
    $0 cleanup --retention 7            # Clean backups older than 7 days
    $0 schedule                         # Set up automated backups
    $0 maintenance                      # Run full maintenance
EOF
}

# Configuration
COMMAND=""
BACKUP_DIR_ARG=""
RETENTION_DAYS=$DEFAULT_RETENTION_DAYS
COMPRESS=true
ENCRYPT=false
S3_BUCKET="$DEFAULT_S3_BUCKET"
FORCE=false
VERBOSE=false
QUIET=false
DRY_RUN=false
DB_ONLY=false
FILES_ONLY=false
INCREMENTAL=false
BACKUP_FILE=""
BACKUP_DATE=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -d|--backup-dir)
            BACKUP_DIR_ARG="$2"
            shift 2
            ;;
        -r|--retention)
            RETENTION_DAYS="$2"
            shift 2
            ;;
        -c|--compress)
            COMPRESS=true
            shift
            ;;
        -e|--encrypt)
            ENCRYPT=true
            shift
            ;;
        -s|--s3-bucket)
            S3_BUCKET="$2"
            shift 2
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --db-only)
            DB_ONLY=true
            shift
            ;;
        --files-only)
            FILES_ONLY=true
            shift
            ;;
        --incremental)
            INCREMENTAL=true
            shift
            ;;
        --backup-file)
            BACKUP_FILE="$2"
            shift 2
            ;;
        --backup-date)
            BACKUP_DATE="$2"
            shift 2
            ;;
        backup|restore|list|cleanup|verify|monitor|rotate-logs|maintenance|schedule)
            COMMAND="$1"
            shift
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Override backup directory if specified
if [[ -n "$BACKUP_DIR_ARG" ]]; then
    BACKUP_DIR="$BACKUP_DIR_ARG"
fi

# Validate command
if [[ -z "$COMMAND" ]]; then
    error "No command specified. Use --help for usage information."
fi

# Helper functions
run_command() {
    local cmd="$1"
    local desc="${2:-Running command}"
    
    if [[ "$DRY_RUN" == true ]]; then
        log "[DRY RUN] Would run: $cmd"
        return 0
    fi
    
    if [[ "$QUIET" == false ]]; then
        log "$desc"
    fi
    
    if [[ "$VERBOSE" == true ]]; then
        log "Command: $cmd"
    fi
    
    if ! eval "$cmd"; then
        error "Command failed: $cmd"
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

check_dependencies() {
    log "Checking dependencies..."
    
    # Check Docker and Docker Compose
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed"
    fi
    
    # Check compression tools
    if [[ "$COMPRESS" == true ]]; then
        if ! command -v gzip &> /dev/null; then
            warning "gzip not found, disabling compression"
            COMPRESS=false
        fi
    fi
    
    # Check encryption tools
    if [[ "$ENCRYPT" == true ]]; then
        if ! command -v gpg &> /dev/null; then
            warning "GPG not found, disabling encryption"
            ENCRYPT=false
        fi
    fi
    
    # Check AWS CLI for S3
    if [[ -n "$S3_BUCKET" ]]; then
        if ! command -v aws &> /dev/null; then
            error "AWS CLI is not installed but S3 bucket specified"
        fi
    fi
    
    success "Dependencies check passed"
}

# Backup functions
create_full_backup() {
    log "Creating full backup..."
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_name="urlshortener_backup_$timestamp"
    local backup_path="$BACKUP_DIR/$backup_name"
    
    # Create backup directory
    mkdir -p "$backup_path"
    
    # Create metadata file
    create_backup_metadata "$backup_path"
    
    # Backup database
    if [[ "$FILES_ONLY" == false ]]; then
        backup_database "$backup_path"
    fi
    
    # Backup files
    if [[ "$DB_ONLY" == false ]]; then
        backup_files "$backup_path"
    fi
    
    # Compress backup
    if [[ "$COMPRESS" == true ]]; then
        compress_backup "$backup_path"
        backup_path="$backup_path.tar.gz"
    fi
    
    # Encrypt backup
    if [[ "$ENCRYPT" == true ]]; then
        encrypt_backup "$backup_path"
        backup_path="$backup_path.gpg"
    fi
    
    # Upload to S3
    if [[ -n "$S3_BUCKET" ]]; then
        upload_to_s3 "$backup_path"
    fi
    
    # Verify backup
    verify_backup "$backup_path"
    
    success "Full backup created: $backup_path"
    echo "$backup_path"
}

create_backup_metadata() {
    local backup_path="$1"
    
    log "Creating backup metadata..."
    
    cat > "$backup_path/metadata.json" << EOF
{
    "backup_type": "$(if [[ "$DB_ONLY" == true ]]; then echo "database"; elif [[ "$FILES_ONLY" == true ]]; then echo "files"; else echo "full"; fi)",
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "version": "$(git rev-parse HEAD 2>/dev/null || echo 'unknown')",
    "hostname": "$(hostname)",
    "user": "$(whoami)",
    "compressed": $COMPRESS,
    "encrypted": $ENCRYPT,
    "incremental": $INCREMENTAL
}
EOF
}

backup_database() {
    local backup_path="$1"
    
    log "Backing up database..."
    
    # Check if database is running
    if ! docker-compose ps postgres | grep -q Up; then
        error "PostgreSQL container is not running"
    fi
    
    # Get database credentials from environment
    source "$PROJECT_ROOT/.env" 2>/dev/null || true
    
    local db_name="${POSTGRES_DB:-urlshortener}"
    local db_user="${POSTGRES_USER:-urlshortener}"
    
    # Create database backup
    run_command "docker-compose exec -T postgres pg_dump -U $db_user -d $db_name --verbose > $backup_path/database.sql" "Creating database backup"
    
    # Create schema backup
    run_command "docker-compose exec -T postgres pg_dump -U $db_user -d $db_name --schema-only > $backup_path/schema.sql" "Creating schema backup"
    
    # Create data-only backup
    run_command "docker-compose exec -T postgres pg_dump -U $db_user -d $db_name --data-only > $backup_path/data.sql" "Creating data backup"
    
    # Backup database metadata
    run_command "docker-compose exec -T postgres psql -U $db_user -d $db_name -c '\\dt' > $backup_path/tables.txt" "Backing up table list"
    
    success "Database backup completed"
}

backup_files() {
    local backup_path="$1"
    
    log "Backing up files..."
    
    # Configuration files
    mkdir -p "$backup_path/config"
    cp "$PROJECT_ROOT/.env" "$backup_path/config/" 2>/dev/null || true
    cp "$PROJECT_ROOT/docker-compose.yml" "$backup_path/config/" 2>/dev/null || true
    cp "$PROJECT_ROOT/docker-compose.prod.yml" "$backup_path/config/" 2>/dev/null || true
    
    # SSL certificates
    if [[ -d "$PROJECT_ROOT/ssl" ]]; then
        mkdir -p "$backup_path/ssl"
        cp -r "$PROJECT_ROOT/ssl"/* "$backup_path/ssl/" 2>/dev/null || true
    fi
    
    # Nginx configuration
    if [[ -d "$PROJECT_ROOT/nginx-sites" ]]; then
        mkdir -p "$backup_path/nginx"
        cp -r "$PROJECT_ROOT/nginx-sites"/* "$backup_path/nginx/" 2>/dev/null || true
    fi
    
    # Application logs
    if [[ -d "$PROJECT_ROOT/logs" ]]; then
        mkdir -p "$backup_path/logs"
        cp -r "$PROJECT_ROOT/logs"/* "$backup_path/logs/" 2>/dev/null || true
    fi
    
    # Custom scripts
    if [[ -d "$PROJECT_ROOT/scripts" ]]; then
        mkdir -p "$backup_path/scripts"
        cp -r "$PROJECT_ROOT/scripts"/* "$backup_path/scripts/" 2>/dev/null || true
    fi
    
    # Documentation
    if [[ -d "$PROJECT_ROOT/docs" ]]; then
        mkdir -p "$backup_path/docs"
        cp -r "$PROJECT_ROOT/docs"/* "$backup_path/docs/" 2>/dev/null || true
    fi
    
    success "Files backup completed"
}

compress_backup() {
    local backup_path="$1"
    
    log "Compressing backup..."
    
    run_command "tar -czf $backup_path.tar.gz -C $(dirname $backup_path) $(basename $backup_path)" "Compressing backup"
    run_command "rm -rf $backup_path" "Removing uncompressed backup"
    
    success "Backup compressed"
}

encrypt_backup() {
    local backup_path="$1"
    
    log "Encrypting backup..."
    
    # Use symmetric encryption for simplicity
    run_command "gpg --symmetric --cipher-algo AES256 --compress-algo 1 --output $backup_path.gpg $backup_path" "Encrypting backup"
    run_command "rm $backup_path" "Removing unencrypted backup"
    
    success "Backup encrypted"
}

upload_to_s3() {
    local backup_path="$1"
    local backup_name=$(basename "$backup_path")
    
    log "Uploading backup to S3..."
    
    run_command "aws s3 cp $backup_path s3://$S3_BUCKET/urlshortener-backups/$backup_name" "Uploading to S3"
    
    success "Backup uploaded to S3"
}

# Restore functions
restore_backup() {
    log "Restoring from backup..."
    
    local backup_file="$BACKUP_FILE"
    
    # Find backup file if date specified
    if [[ -n "$BACKUP_DATE" ]]; then
        backup_file=$(find "$BACKUP_DIR" -name "*$BACKUP_DATE*" -type f | head -1)
        if [[ -z "$backup_file" ]]; then
            error "No backup found for date: $BACKUP_DATE"
        fi
    fi
    
    if [[ -z "$backup_file" ]]; then
        # Use latest backup
        backup_file=$(ls -t "$BACKUP_DIR"/*.tar.gz 2>/dev/null | head -1)
        if [[ -z "$backup_file" ]]; then
            error "No backup file found"
        fi
    fi
    
    log "Restoring from: $backup_file"
    
    if ! confirm "This will replace current data. Continue?"; then
        log "Restore cancelled"
        exit 0
    fi
    
    # Stop services
    run_command "docker-compose down" "Stopping services"
    
    # Extract backup
    local restore_dir="/tmp/restore_$(date +%s)"
    mkdir -p "$restore_dir"
    
    if [[ "$backup_file" == *.gpg ]]; then
        run_command "gpg --decrypt $backup_file | tar -xzf - -C $restore_dir" "Decrypting and extracting backup"
    else
        run_command "tar -xzf $backup_file -C $restore_dir" "Extracting backup"
    fi
    
    # Restore database
    if [[ "$FILES_ONLY" == false ]]; then
        restore_database "$restore_dir"
    fi
    
    # Restore files
    if [[ "$DB_ONLY" == false ]]; then
        restore_files "$restore_dir"
    fi
    
    # Start services
    run_command "docker-compose up -d" "Starting services"
    
    # Clean up
    rm -rf "$restore_dir"
    
    success "Restore completed"
}

restore_database() {
    local restore_dir="$1"
    
    log "Restoring database..."
    
    # Find extracted backup directory
    local backup_extracted=$(find "$restore_dir" -type d -name "urlshortener_backup_*" | head -1)
    if [[ -z "$backup_extracted" ]]; then
        error "Could not find extracted backup directory"
    fi
    
    # Start database
    run_command "docker-compose up -d postgres" "Starting PostgreSQL"
    sleep 10
    
    # Get database credentials
    source "$PROJECT_ROOT/.env" 2>/dev/null || true
    local db_name="${POSTGRES_DB:-urlshortener}"
    local db_user="${POSTGRES_USER:-urlshortener}"
    
    # Drop and recreate database
    run_command "docker-compose exec -T postgres psql -U $db_user -c 'DROP DATABASE IF EXISTS $db_name;'" "Dropping existing database"
    run_command "docker-compose exec -T postgres psql -U $db_user -c 'CREATE DATABASE $db_name;'" "Creating database"
    
    # Restore database
    run_command "docker-compose exec -T postgres psql -U $db_user -d $db_name < $backup_extracted/database.sql" "Restoring database"
    
    success "Database restored"
}

restore_files() {
    local restore_dir="$1"
    
    log "Restoring files..."
    
    local backup_extracted=$(find "$restore_dir" -type d -name "urlshortener_backup_*" | head -1)
    
    # Restore configuration
    if [[ -d "$backup_extracted/config" ]]; then
        cp "$backup_extracted/config/.env" "$PROJECT_ROOT/" 2>/dev/null || true
        cp "$backup_extracted/config/docker-compose.yml" "$PROJECT_ROOT/" 2>/dev/null || true
    fi
    
    # Restore SSL certificates
    if [[ -d "$backup_extracted/ssl" ]]; then
        mkdir -p "$PROJECT_ROOT/ssl"
        cp -r "$backup_extracted/ssl"/* "$PROJECT_ROOT/ssl/" 2>/dev/null || true
    fi
    
    # Restore Nginx configuration
    if [[ -d "$backup_extracted/nginx" ]]; then
        mkdir -p "$PROJECT_ROOT/nginx-sites"
        cp -r "$backup_extracted/nginx"/* "$PROJECT_ROOT/nginx-sites/" 2>/dev/null || true
    fi
    
    success "Files restored"
}

# List backups
list_backups() {
    log "Listing available backups..."
    
    if [[ ! -d "$BACKUP_DIR" ]]; then
        warning "Backup directory does not exist: $BACKUP_DIR"
        return
    fi
    
    echo "Available backups in $BACKUP_DIR:"
    echo "=================================="
    
    local backups=($(ls -t "$BACKUP_DIR"/*.tar.gz 2>/dev/null || true))
    
    if [[ ${#backups[@]} -eq 0 ]]; then
        echo "No backups found"
        return
    fi
    
    for backup in "${backups[@]}"; do
        local size=$(du -h "$backup" | cut -f1)
        local date=$(stat -c %y "$backup" 2>/dev/null || stat -f %Sm "$backup" 2>/dev/null || echo "unknown")
        echo "$(basename "$backup") - $size - $date"
    done
}

# Cleanup old backups
cleanup_backups() {
    log "Cleaning up old backups (retention: $RETENTION_DAYS days)..."
    
    if [[ ! -d "$BACKUP_DIR" ]]; then
        warning "Backup directory does not exist: $BACKUP_DIR"
        return
    fi
    
    local count=0
    while IFS= read -r -d '' backup; do
        run_command "rm '$backup'" "Removing old backup: $(basename "$backup")"
        ((count++))
    done < <(find "$BACKUP_DIR" -name "*.tar.gz" -mtime +$RETENTION_DAYS -print0 2>/dev/null)
    
    if [[ $count -eq 0 ]]; then
        log "No old backups to clean up"
    else
        success "Cleaned up $count old backups"
    fi
    
    # Also clean S3 if configured
    if [[ -n "$S3_BUCKET" ]]; then
        log "Cleaning up old S3 backups..."
        run_command "aws s3 ls s3://$S3_BUCKET/urlshortener-backups/ | while read -r line; do
            backup_date=\$(echo \$line | awk '{print \$1}')
            backup_file=\$(echo \$line | awk '{print \$4}')
            if [[ -n \"\$backup_date\" ]] && [[ \$(date -d \"\$backup_date\" +%s) -lt \$(date -d \"$RETENTION_DAYS days ago\" +%s) ]]; then
                aws s3 rm s3://$S3_BUCKET/urlshortener-backups/\$backup_file
            fi
        done" "Cleaning up old S3 backups"
    fi
}

# Verify backup integrity
verify_backup() {
    local backup_file="$1"
    
    log "Verifying backup integrity..."
    
    if [[ ! -f "$backup_file" ]]; then
        error "Backup file not found: $backup_file"
    fi
    
    # Test archive integrity
    if [[ "$backup_file" == *.tar.gz ]]; then
        if ! tar -tzf "$backup_file" &> /dev/null; then
            error "Backup archive is corrupted"
        fi
    fi
    
    # Test decryption if encrypted
    if [[ "$backup_file" == *.gpg ]]; then
        if ! gpg --list-packets "$backup_file" &> /dev/null; then
            error "Backup encryption is corrupted"
        fi
    fi
    
    success "Backup integrity verified"
}

# System monitoring
system_monitor() {
    log "Running system monitoring..."
    
    echo "System Health Report - $(date)"
    echo "=============================="
    
    # Disk usage
    echo -e "\nDisk Usage:"
    df -h
    
    # Memory usage
    echo -e "\nMemory Usage:"
    free -h
    
    # Docker status
    echo -e "\nDocker Services:"
    docker-compose ps
    
    # Container resource usage
    echo -e "\nContainer Resource Usage:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
    
    # Application health
    echo -e "\nApplication Health:"
    if curl -sf http://localhost:8080/health &> /dev/null; then
        echo "✓ Backend is healthy"
    else
        echo "✗ Backend is unhealthy"
    fi
    
    # Database status
    echo -e "\nDatabase Status:"
    if docker-compose exec -T postgres pg_isready &> /dev/null; then
        echo "✓ Database is healthy"
    else
        echo "✗ Database is unhealthy"
    fi
    
    # Redis status
    echo -e "\nRedis Status:"
    if docker-compose exec -T redis redis-cli ping &> /dev/null; then
        echo "✓ Redis is healthy"
    else
        echo "✗ Redis is unhealthy"
    fi
    
    # Backup status
    echo -e "\nBackup Status:"
    local latest_backup=$(ls -t "$BACKUP_DIR"/*.tar.gz 2>/dev/null | head -1)
    if [[ -n "$latest_backup" ]]; then
        local backup_age=$((($(date +%s) - $(stat -c %Y "$latest_backup")) / 86400))
        echo "✓ Latest backup: $(basename "$latest_backup") ($backup_age days old)"
    else
        echo "✗ No backups found"
    fi
}

# Log rotation
rotate_logs() {
    log "Rotating application logs..."
    
    local log_dir="$PROJECT_ROOT/logs"
    
    if [[ ! -d "$log_dir" ]]; then
        warning "Log directory does not exist: $log_dir"
        return
    fi
    
    # Rotate logs older than 1 day
    find "$log_dir" -name "*.log" -mtime +1 -exec gzip {} \;
    
    # Remove compressed logs older than 30 days
    find "$log_dir" -name "*.log.gz" -mtime +30 -delete
    
    # Restart services to reopen log files
    run_command "docker-compose restart" "Restarting services for log rotation"
    
    success "Log rotation completed"
}

# Full maintenance routine
run_maintenance() {
    log "Running full maintenance routine..."
    
    # System monitoring
    system_monitor
    
    # Create backup
    create_full_backup
    
    # Clean old backups
    cleanup_backups
    
    # Rotate logs
    rotate_logs
    
    # Docker cleanup
    log "Cleaning up Docker resources..."
    run_command "docker system prune -f" "Removing unused Docker resources"
    
    success "Maintenance routine completed"
}

# Schedule automated backups
schedule_backups() {
    log "Setting up automated backup schedule..."
    
    local cron_job="0 2 * * * $SCRIPT_DIR/backup.sh backup --quiet >> $LOG_FILE 2>&1"
    local cleanup_job="0 3 * * 0 $SCRIPT_DIR/backup.sh cleanup --quiet >> $LOG_FILE 2>&1"
    
    if confirm "Add daily backup at 2 AM?"; then
        (crontab -l 2>/dev/null; echo "$cron_job") | crontab -
        success "Daily backup scheduled"
    fi
    
    if confirm "Add weekly cleanup on Sunday at 3 AM?"; then
        (crontab -l 2>/dev/null; echo "$cleanup_job") | crontab -
        success "Weekly cleanup scheduled"
    fi
    
    echo "Current crontab:"
    crontab -l
}

# Main execution
main() {
    log "Starting backup script - Command: $COMMAND"
    
    # Create backup directory
    mkdir -p "$BACKUP_DIR"
    mkdir -p "$(dirname "$LOG_FILE")"
    
    # Check dependencies
    check_dependencies
    
    # Execute command
    case "$COMMAND" in
        "backup")
            create_full_backup
            ;;
        "restore")
            restore_backup
            ;;
        "list")
            list_backups
            ;;
        "cleanup")
            cleanup_backups
            ;;
        "verify")
            if [[ -n "$BACKUP_FILE" ]]; then
                verify_backup "$BACKUP_FILE"
            else
                error "Backup file must be specified for verification"
            fi
            ;;
        "monitor")
            system_monitor
            ;;
        "rotate-logs")
            rotate_logs
            ;;
        "maintenance")
            run_maintenance
            ;;
        "schedule")
            schedule_backups
            ;;
        *)
            error "Unknown command: $COMMAND"
            ;;
    esac
}

# Run main function
main "$@"