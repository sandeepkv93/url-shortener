#!/bin/bash

# Setup Script for URL Shortener
# This script sets up the development and production environment

set -euo pipefail

# Script Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
URL Shortener Setup Script

Usage: $0 [OPTIONS] [ENVIRONMENT]

Environments:
    dev         Set up development environment (default)
    prod        Set up production environment
    docker      Set up Docker-based environment

Options:
    -h, --help              Show this help message
    -f, --force             Force setup without confirmations
    -s, --skip-deps         Skip dependency installation
    -d, --docker-only       Only set up Docker environment
    -c, --check-only        Only check prerequisites
    --install-tools         Install additional development tools

Examples:
    $0                      # Set up development environment
    $0 dev                  # Set up development environment
    $0 prod                 # Set up production environment
    $0 docker               # Set up Docker environment
    $0 --check-only         # Check prerequisites only
    $0 --install-tools dev  # Install dev tools and set up development
EOF
}

# Configuration
ENVIRONMENT="dev"
FORCE=false
SKIP_DEPS=false
DOCKER_ONLY=false
CHECK_ONLY=false
INSTALL_TOOLS=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -s|--skip-deps)
            SKIP_DEPS=true
            shift
            ;;
        -d|--docker-only)
            DOCKER_ONLY=true
            shift
            ;;
        -c|--check-only)
            CHECK_ONLY=true
            shift
            ;;
        --install-tools)
            INSTALL_TOOLS=true
            shift
            ;;
        dev|prod|docker)
            ENVIRONMENT="$1"
            shift
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Helper functions
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

run_command() {
    local cmd="$1"
    local desc="${2:-Running command}"
    
    log "$desc: $cmd"
    if ! eval "$cmd"; then
        error "Command failed: $cmd"
    fi
}

check_command() {
    local cmd="$1"
    local name="$2"
    
    if command -v "$cmd" &> /dev/null; then
        success "$name is installed"
        return 0
    else
        warning "$name is not installed"
        return 1
    fi
}

# System detection
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [[ -f /etc/debian_version ]]; then
            echo "debian"
        elif [[ -f /etc/redhat-release ]]; then
            echo "redhat"
        else
            echo "linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

# Prerequisite checks
check_prerequisites() {
    log "Checking prerequisites..."
    
    local os=$(detect_os)
    log "Detected OS: $os"
    
    local missing_deps=()
    
    # Check basic tools
    if ! check_command "git" "Git"; then
        missing_deps+=("git")
    fi
    
    if ! check_command "curl" "cURL"; then
        missing_deps+=("curl")
    fi
    
    # Check Go
    if ! check_command "go" "Go"; then
        missing_deps+=("go")
    else
        local go_version=$(go version | awk '{print $3}' | sed 's/go//')
        local required_go="1.21"
        if [[ "$(printf '%s\n' "$required_go" "$go_version" | sort -V | head -n1)" != "$required_go" ]]; then
            warning "Go version $go_version found, but $required_go or higher is recommended"
        else
            success "Go version $go_version is compatible"
        fi
    fi
    
    # Check Node.js
    if ! check_command "node" "Node.js"; then
        missing_deps+=("nodejs")
    else
        local node_version=$(node --version | sed 's/v//')
        local required_node="18.0.0"
        if [[ "$(printf '%s\n' "$required_node" "$node_version" | sort -V | head -n1)" != "$required_node" ]]; then
            warning "Node.js version $node_version found, but $required_node or higher is recommended"
        else
            success "Node.js version $node_version is compatible"
        fi
    fi
    
    # Check npm
    if ! check_command "npm" "npm"; then
        missing_deps+=("npm")
    fi
    
    # Check Docker (optional but recommended)
    if ! check_command "docker" "Docker"; then
        warning "Docker is not installed (optional for development)"
    else
        if ! docker info &> /dev/null; then
            warning "Docker is installed but not running"
        else
            success "Docker is running"
        fi
    fi
    
    # Check Docker Compose
    if ! check_command "docker-compose" "Docker Compose"; then
        warning "Docker Compose is not installed (optional for development)"
    fi
    
    # Environment-specific checks
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        # Production-specific checks
        if ! check_command "systemctl" "systemd"; then
            warning "systemd is not available (required for production service management)"
        fi
        
        if ! check_command "nginx" "Nginx"; then
            missing_deps+=("nginx")
        fi
        
        if ! check_command "postgresql" "PostgreSQL Client" && ! check_command "psql" "PostgreSQL Client"; then
            missing_deps+=("postgresql-client")
        fi
    fi
    
    # Report missing dependencies
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        warning "Missing dependencies: ${missing_deps[*]}"
        
        if [[ "$SKIP_DEPS" == false ]]; then
            if confirm "Install missing dependencies?"; then
                install_dependencies "${missing_deps[@]}"
            else
                warning "Continuing without installing dependencies. Some features may not work."
            fi
        fi
    else
        success "All prerequisites are met"
    fi
}

# Dependency installation
install_dependencies() {
    local deps=("$@")
    local os=$(detect_os)
    
    log "Installing dependencies: ${deps[*]}"
    
    case "$os" in
        "debian")
            run_command "sudo apt update" "Updating package lists"
            for dep in "${deps[@]}"; do
                case "$dep" in
                    "go")
                        install_go_debian
                        ;;
                    "nodejs")
                        install_nodejs_debian
                        ;;
                    "docker")
                        install_docker_debian
                        ;;
                    *)
                        run_command "sudo apt install -y $dep" "Installing $dep"
                        ;;
                esac
            done
            ;;
        "redhat")
            run_command "sudo yum update -y" "Updating package lists"
            for dep in "${deps[@]}"; do
                run_command "sudo yum install -y $dep" "Installing $dep"
            done
            ;;
        "macos")
            if ! command -v brew &> /dev/null; then
                log "Installing Homebrew..."
                run_command '/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"' "Installing Homebrew"
            fi
            
            for dep in "${deps[@]}"; do
                case "$dep" in
                    "nodejs")
                        run_command "brew install node" "Installing Node.js"
                        ;;
                    "postgresql-client")
                        run_command "brew install postgresql" "Installing PostgreSQL"
                        ;;
                    *)
                        run_command "brew install $dep" "Installing $dep"
                        ;;
                esac
            done
            ;;
        *)
            warning "Automatic dependency installation not supported on this OS. Please install manually: ${deps[*]}"
            ;;
    esac
}

install_go_debian() {
    log "Installing Go..."
    local go_version="1.21.0"
    local go_archive="go${go_version}.linux-amd64.tar.gz"
    
    run_command "curl -L https://golang.org/dl/$go_archive -o /tmp/$go_archive" "Downloading Go"
    run_command "sudo rm -rf /usr/local/go" "Removing old Go installation"
    run_command "sudo tar -C /usr/local -xzf /tmp/$go_archive" "Installing Go"
    
    # Add to PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    fi
    
    export PATH=$PATH:/usr/local/go/bin
    success "Go installed successfully"
}

install_nodejs_debian() {
    log "Installing Node.js..."
    run_command "curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -" "Adding Node.js repository"
    run_command "sudo apt-get install -y nodejs" "Installing Node.js"
}

install_docker_debian() {
    log "Installing Docker..."
    run_command "curl -fsSL https://get.docker.com -o get-docker.sh" "Downloading Docker installer"
    run_command "sudo sh get-docker.sh" "Installing Docker"
    run_command "sudo usermod -aG docker $USER" "Adding user to docker group"
    run_command "rm get-docker.sh" "Cleaning up"
    
    # Install Docker Compose
    run_command "sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose" "Installing Docker Compose"
    run_command "sudo chmod +x /usr/local/bin/docker-compose" "Making Docker Compose executable"
    
    warning "Please log out and back in for Docker group membership to take effect"
}

# Development tools installation
install_development_tools() {
    log "Installing development tools..."
    
    # Go tools
    if command -v go &> /dev/null; then
        log "Installing Go development tools..."
        go install golang.org/x/tools/cmd/goimports@latest
        go install golang.org/x/lint/golint@latest
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    fi
    
    # Node.js tools
    if command -v npm &> /dev/null; then
        log "Installing Node.js development tools..."
        npm install -g typescript@latest
        npm install -g @typescript-eslint/parser@latest
        npm install -g prettier@latest
        npm install -g vite@latest
    fi
    
    # Additional tools
    local os=$(detect_os)
    case "$os" in
        "debian")
            run_command "sudo apt install -y jq htop tree" "Installing additional tools"
            ;;
        "macos")
            run_command "brew install jq htop tree" "Installing additional tools"
            ;;
    esac
    
    success "Development tools installed"
}

# Environment setup
setup_environment_files() {
    log "Setting up environment files..."
    
    cd "$PROJECT_ROOT"
    
    # Main environment file
    if [[ ! -f ".env" ]]; then
        if [[ "$ENVIRONMENT" == "prod" ]]; then
            if [[ -f ".env.production" ]]; then
                cp .env.production .env
                log "Copied production environment template"
            else
                cp .env.example .env
                log "Copied example environment file"
            fi
        else
            cp .env.example .env
            log "Copied example environment file"
        fi
        
        warning "Please review and update .env file with your specific configuration"
    else
        log "Environment file already exists"
    fi
    
    # Backend environment
    if [[ ! -f "backend/.env" ]] && [[ -f "backend/.env.example" ]]; then
        cp backend/.env.example backend/.env
        log "Created backend environment file"
    fi
    
    # Frontend environment
    if [[ ! -f "frontend/.env" ]] && [[ -f "frontend/.env.example" ]]; then
        cp frontend/.env.example frontend/.env
        log "Created frontend environment file"
    fi
}

setup_directories() {
    log "Setting up directories..."
    
    cd "$PROJECT_ROOT"
    
    # Data directories
    mkdir -p data/{postgres,redis,prometheus,grafana}
    mkdir -p logs
    mkdir -p backups
    mkdir -p ssl
    
    # Set permissions
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        chmod 700 data
        chmod 700 ssl
        chmod 700 backups
    fi
    
    success "Directories created"
}

setup_development() {
    log "Setting up development environment..."
    
    cd "$PROJECT_ROOT"
    
    # Install backend dependencies
    if [[ -f "backend/go.mod" ]]; then
        cd backend
        run_command "go mod download" "Downloading Go dependencies"
        run_command "go mod verify" "Verifying Go dependencies"
        cd ..
    fi
    
    # Install frontend dependencies
    if [[ -f "frontend/package.json" ]]; then
        cd frontend
        run_command "npm install" "Installing Node.js dependencies"
        cd ..
    fi
    
    success "Development environment setup complete"
}

setup_docker() {
    log "Setting up Docker environment..."
    
    cd "$PROJECT_ROOT"
    
    if [[ ! -f "docker-compose.yml" ]]; then
        if [[ "$ENVIRONMENT" == "prod" ]]; then
            cp docker-compose.prod.yml docker-compose.yml
            log "Using production Docker Compose configuration"
        else
            log "Using development Docker Compose configuration (default)"
        fi
    fi
    
    # Build images
    if confirm "Build Docker images now?"; then
        if [[ "$ENVIRONMENT" == "prod" ]]; then
            run_command "docker-compose -f docker-compose.prod.yml build" "Building production Docker images"
        else
            run_command "docker-compose build" "Building development Docker images"
        fi
    fi
    
    success "Docker environment setup complete"
}

setup_production() {
    log "Setting up production environment..."
    
    # Additional production setup
    if [[ "$EUID" -eq 0 ]]; then
        warning "Running as root. Consider creating a dedicated user for the application."
    fi
    
    # Create systemd service file
    create_systemd_service
    
    # Set up log rotation
    setup_log_rotation
    
    # Set up firewall rules
    setup_firewall
    
    success "Production environment setup complete"
}

create_systemd_service() {
    log "Creating systemd service..."
    
    local service_file="/etc/systemd/system/url-shortener.service"
    
    if [[ ! -f "$service_file" ]]; then
        cat > /tmp/url-shortener.service << EOF
[Unit]
Description=URL Shortener Service
After=network.target docker.service
Wants=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$PROJECT_ROOT
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml down
TimeoutStartSec=0
User=$(whoami)

[Install]
WantedBy=multi-user.target
EOF
        
        if confirm "Install systemd service?"; then
            run_command "sudo mv /tmp/url-shortener.service $service_file" "Installing systemd service"
            run_command "sudo systemctl daemon-reload" "Reloading systemd"
            run_command "sudo systemctl enable url-shortener" "Enabling service"
        fi
    else
        log "Systemd service already exists"
    fi
}

setup_log_rotation() {
    log "Setting up log rotation..."
    
    local logrotate_file="/etc/logrotate.d/url-shortener"
    
    if [[ ! -f "$logrotate_file" ]]; then
        cat > /tmp/url-shortener-logrotate << EOF
$PROJECT_ROOT/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    notifempty
    create 644 $(whoami) $(whoami)
    postrotate
        docker-compose -f $PROJECT_ROOT/docker-compose.prod.yml restart backend frontend
    endscript
}
EOF
        
        if confirm "Install log rotation configuration?"; then
            run_command "sudo mv /tmp/url-shortener-logrotate $logrotate_file" "Installing log rotation"
        fi
    else
        log "Log rotation already configured"
    fi
}

setup_firewall() {
    log "Setting up firewall rules..."
    
    if command -v ufw &> /dev/null; then
        if confirm "Configure UFW firewall?"; then
            run_command "sudo ufw allow ssh" "Allowing SSH"
            run_command "sudo ufw allow 80" "Allowing HTTP"
            run_command "sudo ufw allow 443" "Allowing HTTPS"
            
            if confirm "Enable UFW firewall?"; then
                run_command "sudo ufw --force enable" "Enabling UFW"
            fi
        fi
    else
        warning "UFW firewall not available. Configure firewall manually."
    fi
}

# Validation
validate_setup() {
    log "Validating setup..."
    
    local errors=0
    
    # Check environment files
    if [[ ! -f ".env" ]]; then
        error "Main environment file not found"
        ((errors++))
    fi
    
    # Check Go module
    if [[ -f "backend/go.mod" ]]; then
        cd backend
        if ! go mod verify &> /dev/null; then
            error "Go modules verification failed"
            ((errors++))
        fi
        cd ..
    fi
    
    # Check Node.js dependencies
    if [[ -f "frontend/package.json" ]] && [[ -d "frontend/node_modules" ]]; then
        cd frontend
        if ! npm ls &> /dev/null; then
            warning "Some npm dependencies may be missing"
        fi
        cd ..
    fi
    
    # Check Docker setup
    if [[ "$DOCKER_ONLY" == true ]] || [[ "$ENVIRONMENT" == "docker" ]]; then
        if ! docker info &> /dev/null; then
            error "Docker is not running"
            ((errors++))
        fi
        
        if ! command -v docker-compose &> /dev/null; then
            error "Docker Compose is not available"
            ((errors++))
        fi
    fi
    
    if [[ $errors -eq 0 ]]; then
        success "Setup validation passed"
    else
        error "Setup validation failed with $errors errors"
    fi
}

# Main execution
main() {
    log "Starting setup for environment: $ENVIRONMENT"
    
    # Check prerequisites
    check_prerequisites
    
    if [[ "$CHECK_ONLY" == true ]]; then
        log "Prerequisites check completed"
        exit 0
    fi
    
    # Install development tools if requested
    if [[ "$INSTALL_TOOLS" == true ]]; then
        install_development_tools
    fi
    
    # Setup environment files
    setup_environment_files
    
    # Setup directories
    setup_directories
    
    # Environment-specific setup
    case "$ENVIRONMENT" in
        "dev")
            if [[ "$DOCKER_ONLY" == false ]]; then
                setup_development
            fi
            setup_docker
            ;;
        "prod")
            setup_production
            setup_docker
            ;;
        "docker")
            setup_docker
            ;;
    esac
    
    # Validate setup
    validate_setup
    
    # Final instructions
    log "Setup completed for $ENVIRONMENT environment"
    
    case "$ENVIRONMENT" in
        "dev")
            echo -e "\n${GREEN}Next steps:${NC}"
            echo "1. Review and update .env files"
            echo "2. Start development: docker-compose up -d"
            echo "3. Or run locally: cd backend && go run cmd/server/main.go"
            echo "4. Frontend: cd frontend && npm run dev"
            ;;
        "prod")
            echo -e "\n${GREEN}Next steps:${NC}"
            echo "1. Update .env with production values"
            echo "2. Configure SSL certificates in ssl/ directory"
            echo "3. Review security settings"
            echo "4. Deploy: ./scripts/deploy.sh deploy"
            ;;
        "docker")
            echo -e "\n${GREEN}Next steps:${NC}"
            echo "1. Review docker-compose configuration"
            echo "2. Start services: docker-compose up -d"
            echo "3. Check status: docker-compose ps"
            ;;
    esac
    
    success "URL Shortener setup completed successfully!"
}

# Run main function
main "$@"