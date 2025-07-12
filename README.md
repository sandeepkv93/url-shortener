# URL Shortener Service

A production-grade URL shortener service built with Go and React, featuring comprehensive analytics, QR code generation, and advanced link management capabilities.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-18+-blue.svg)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-blue.svg)](https://typescriptlang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](#testing)

## 🚀 Features

### Core Features
- **URL Shortening**: Create short, memorable links from long URLs
- **Custom Aliases**: Support for custom short codes and vanity URLs
- **QR Code Generation**: Generate QR codes for shortened URLs in multiple formats (PNG, SVG, PDF)
- **Password Protection**: Secure URLs with password protection
- **Expiration Dates**: Set automatic expiration for URLs
- **Bulk Operations**: Import/export URLs in bulk

### Analytics & Reporting
- **Real-time Analytics**: Live click tracking and statistics
- **Geographic Analytics**: Location-based click analysis with interactive maps
- **Device Analytics**: Track clicks by device type, browser, and operating system
- **Time-based Analytics**: Detailed time-series analysis with customizable date ranges
- **Referrer Analytics**: Track traffic sources and campaigns
- **Export Capabilities**: Export analytics data in JSON, CSV, and Excel formats

### User Management
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **User Profiles**: Comprehensive profile management with preferences
- **Dashboard**: Interactive dashboard with URL management and analytics
- **Role-based Access**: Support for different user roles and permissions

### Advanced Features
- **Health Monitoring**: Comprehensive system health checks and metrics
- **Rate Limiting**: Intelligent API rate limiting and abuse prevention
- **Caching**: Multi-layer Redis-based caching for optimal performance
- **Security**: Enterprise-grade security features (CORS, XSS protection, CSRF protection)
- **Geolocation Services**: IP-based location tracking with privacy controls
- **API Documentation**: Interactive OpenAPI/Swagger documentation

## 📚 Documentation

### Quick Start
- **[Developer Onboarding Guide](docs/developer-onboarding.md)** - Complete guide for new developers
- **[Production Deployment Guide](docs/production-deployment.md)** - Comprehensive production deployment and operations

### API Documentation
- **[Interactive API Documentation](http://localhost:8080/docs)** - Swagger UI (when running)
- **[OpenAPI Specification](docs/api/openapi.yaml)** - Complete API specification
- **[API Client Examples](docs/api/examples/)** - Client examples in multiple languages

### Operations and Monitoring
- **[Health Check Documentation](docs/health-checks.md)** - Health monitoring and diagnostics
- **[Automation Scripts Guide](scripts/README.md)** - Deployment and maintenance automation

### Development
- **[Project Progress](PROGRESS.md)** - Development progress and milestones
- **[Environment Setup](docs/developer-onboarding.md#development-environment-setup)** - Detailed setup instructions

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    URL Shortener Service                    │
├─────────────────┬─────────────────┬─────────────────────────┤
│   Frontend      │     Backend     │    Infrastructure       │
│   (React/TS)    │     (Go)        │                         │
├─────────────────┼─────────────────┼─────────────────────────┤
│ • Components    │ • HTTP Handlers │ • PostgreSQL Database   │
│ • Pages         │ • Business      │ • Redis Cache           │
│ • Services      │   Logic         │ • Docker Containers     │
│ • State Mgmt    │ • Domain Models │ • Nginx Load Balancer   │
│ • UI/UX         │ • Repositories  │ • Monitoring Stack      │
└─────────────────┴─────────────────┴─────────────────────────┘

              API Communication (REST/JSON)
                     Health Checks
                  Performance Monitoring
```

### Backend Architecture (Clean Architecture)

```
backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP layer
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # HTTP middleware (auth, logging, security)
│   │   └── routes/      # Route configuration and setup
│   ├── core/            # Business logic layer
│   │   ├── domain/      # Domain entities and business rules
│   │   ├── ports/       # Interface definitions (ports)
│   │   └── services/    # Business logic implementation
│   ├── infrastructure/  # External dependencies layer
│   │   ├── database/    # Database connections and repositories
│   │   ├── cache/       # Redis cache implementation
│   │   └── external/    # External API integrations
│   └── config/          # Configuration management
└── tests/               # Comprehensive test suites
```

## 🛠️ Technology Stack

### Backend Technologies
- **Language**: Go 1.21+ with modern Go features
- **Framework**: Chi Router v5 - Lightweight, fast HTTP router
- **Database**: PostgreSQL 15+ with GORM ORM
- **Cache**: Redis 7+ for high-performance caching
- **Authentication**: JWT tokens with secure refresh mechanism
- **Testing**: Testify framework with 95%+ coverage
- **Monitoring**: Built-in health checks and metrics
- **Security**: Comprehensive security middleware

### Frontend Technologies
- **Framework**: React 18 with TypeScript 5.0+
- **Build Tool**: Vite for fast development and optimized builds
- **State Management**: React Context API with useReducer
- **Styling**: Tailwind CSS v3 with custom design system
- **HTTP Client**: Axios with interceptors and error handling
- **Charts**: Recharts for beautiful, responsive analytics
- **Forms**: React Hook Form with Zod validation
- **Testing**: Vitest + React Testing Library

### Infrastructure & DevOps
- **Containerization**: Docker & Docker Compose with multi-stage builds
- **Task Automation**: Task runner (Taskfile) for development workflows
- **Monitoring**: Prometheus metrics + Grafana dashboards
- **Logging**: Structured logging with JSON format
- **Security**: SSL/TLS, firewall configuration, security headers
- **Deployment**: Automated deployment scripts with rollback support

## 📋 Prerequisites

### Required Software
- **Go** 1.21 or higher
- **Node.js** 18 or higher with npm 9+
- **PostgreSQL** 15+ (or use Docker)
- **Redis** 7+ (or use Docker)
- **Docker** 20.10+ & **Docker Compose** 2.0+ (recommended)

### System Requirements
- **OS**: Linux, macOS, or Windows with WSL2
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 10GB free space for development
- **Network**: Stable internet connection

## 🚀 Quick Start

### Option 1: Automated Setup (Recommended)

```bash
# 1. Clone the repository
git clone https://github.com/your-org/url-shortener.git
cd url-shortener

# 2. Run automated setup
chmod +x scripts/setup.sh
./scripts/setup.sh dev

# 3. Start all services
docker-compose up -d

# 4. Verify installation
curl http://localhost:8080/health/checks
```

### Option 2: Manual Setup

```bash
# 1. Clone and configure
git clone https://github.com/your-org/url-shortener.git
cd url-shortener
cp .env.example .env

# 2. Start infrastructure
docker-compose up -d postgres redis

# 3. Start backend
cd backend
go mod download
go run cmd/server/main.go

# 4. Start frontend (in another terminal)
cd frontend
npm install
npm run dev
```

### 🌐 Access Points

Once running, access the application at:

- **🖥️ Frontend**: http://localhost:3000 (development)
- **🔧 Backend API**: http://localhost:8080/api/v1
- **📖 API Documentation**: http://localhost:8080/docs
- **🏥 Health Checks**: http://localhost:8080/health
- **📊 System Metrics**: http://localhost:8080/health/metrics/system

## 📁 Project Structure

```
url-shortener/
├── 📄 README.md                      # Project overview (this file)
├── 📊 PROGRESS.md                    # Development progress tracking
├── ⚙️ Taskfile.yml                   # Task automation definitions
├── 🐳 docker-compose.yml             # Development Docker setup
├── 🐳 docker-compose.prod.yml        # Production Docker configuration
├── 🔧 .env.example                   # Environment variables template
│
├── 📁 backend/                       # 🟢 Go Backend Service
│   ├── 🚀 cmd/server/               # Application entry point
│   ├── 🏗️ internal/                 # Private application code
│   │   ├── 🌐 api/                  # HTTP layer (handlers, middleware, routes)
│   │   ├── 💼 core/                 # Business logic (domain, services, ports)
│   │   ├── 🏭 infrastructure/       # External dependencies (DB, cache, APIs)
│   │   └── ⚙️ config/               # Configuration management
│   ├── 🧪 tests/                    # Comprehensive test suites
│   ├── 🐳 Dockerfile                # Backend container definition
│   ├── 📦 go.mod                    # Go module dependencies
│   └── 🔧 .env.example              # Backend environment template
│
├── 📁 frontend/                      # ⚛️ React Frontend Application
│   ├── 🌍 public/                   # Static assets and index.html
│   ├── 💻 src/                      # TypeScript source code
│   │   ├── 🧩 components/           # Reusable React components
│   │   ├── 📄 pages/                # Page-level components
│   │   ├── 🪝 hooks/                # Custom React hooks
│   │   ├── 🔌 services/             # API communication services
│   │   ├── 🛠️ utils/                # Utility functions and helpers
│   │   ├── 📋 types/                # TypeScript type definitions
│   │   └── 🗃️ context/              # React context providers
│   ├── 🧪 tests/                    # Frontend test suites
│   ├── 🐳 Dockerfile                # Frontend container definition
│   ├── 📦 package.json              # NPM dependencies and scripts
│   ├── ⚙️ tsconfig.json             # TypeScript configuration
│   ├── 🎨 tailwind.config.js        # Tailwind CSS configuration
│   └── ⚡ vite.config.ts            # Vite build tool configuration
│
├── 📁 docs/                          # 📚 Comprehensive Documentation
│   ├── 🔗 api/                      # API documentation and examples
│   │   ├── 📋 openapi.yaml          # OpenAPI 3.0 specification
│   │   ├── 🌐 swagger-ui.html       # Interactive API documentation
│   │   └── 💡 examples/             # Client examples (JS, Python, Go, cURL)
│   ├── 👨‍💻 developer-onboarding.md    # Complete developer guide
│   ├── 🚀 production-deployment.md   # Production deployment guide
│   └── 🏥 health-checks.md          # Health monitoring documentation
│
└── 📁 scripts/                       # 🤖 Automation Scripts
    ├── 🔧 setup.sh                  # Environment setup automation
    ├── 🚀 deploy.sh                 # Deployment automation with rollback
    ├── 🧪 test.sh                   # Comprehensive testing automation
    ├── 💾 backup.sh                 # Backup and maintenance automation
    └── 📖 README.md                 # Scripts documentation
```

## 🔧 Configuration

### Environment Variables

**Backend Configuration** (`.env`):
```env
# Application Settings
APP_ENV=development
APP_DEBUG=true
PORT=8080
BASE_URL=http://localhost:8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=urlshortener
DB_PASSWORD=your_secure_password
DB_NAME=urlshortener
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Authentication
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_TOKEN_TTL=24h
JWT_REFRESH_TOKEN_TTL=7d

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_GLOBAL=100
RATE_LIMIT_AUTH=5
RATE_LIMIT_URL_CREATION=10

# External Services
GEOLOCATION_API_KEY=your_api_key
```

**Frontend Configuration** (`.env`):
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_NAME=URL Shortener
VITE_ENABLE_ANALYTICS=true
```

For complete configuration details, see the [Developer Onboarding Guide](docs/developer-onboarding.md#configuration).

## 🧪 Testing

### Comprehensive Testing Strategy

We maintain **95%+ test coverage** across all components:

#### Backend Tests
```bash
# Run all backend tests
cd backend && go test ./... -v

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test suites
go test ./internal/core/services/... -v
go test ./internal/api/handlers/... -v
```

#### Frontend Tests
```bash
# Run all frontend tests
cd frontend && npm test

# Run with coverage
npm run test:coverage

# Run in watch mode during development
npm run test:watch
```

#### Automated Testing
```bash
# Run all tests with automation script
./scripts/test.sh all

# Run specific test types
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh e2e

# Run tests with coverage reporting
./scripts/test.sh all --coverage

# Run tests in CI mode
./scripts/test.sh --ci all
```

### Test Categories

- **Unit Tests**: Business logic, utility functions, components
- **Integration Tests**: API endpoints with real database
- **Component Tests**: React components with user interactions
- **E2E Tests**: Complete user workflows (planned)
- **Performance Tests**: Load testing and benchmarks
- **Security Tests**: Vulnerability scanning and penetration testing

For detailed testing procedures, see the [Developer Onboarding Guide](docs/developer-onboarding.md#testing-strategy).

## 📊 API Documentation

### Interactive Documentation

**Live API Documentation**: http://localhost:8080/docs (when running)

The service provides comprehensive API documentation with:
- **Interactive Testing**: Try API endpoints directly from the browser
- **Authentication Guide**: Complete authentication flow examples
- **Request/Response Examples**: Real examples for all endpoints
- **Error Handling**: Detailed error codes and troubleshooting

### API Overview

#### Authentication Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/v1/auth/register` | Register new user | ❌ |
| `POST` | `/api/v1/auth/login` | User login | ❌ |
| `POST` | `/api/v1/auth/refresh` | Refresh JWT token | ❌ |
| `GET` | `/api/v1/auth/profile` | Get user profile | ✅ |
| `PUT` | `/api/v1/auth/profile` | Update user profile | ✅ |
| `POST` | `/api/v1/auth/logout` | User logout | ✅ |

#### URL Management Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/v1/urls` | Create short URL | ✅ |
| `GET` | `/api/v1/urls` | List user URLs | ✅ |
| `GET` | `/api/v1/urls/{id}` | Get URL details | ✅ |
| `PUT` | `/api/v1/urls/{id}` | Update URL | ✅ |
| `DELETE` | `/api/v1/urls/{id}` | Delete URL | ✅ |
| `GET` | `/{shortCode}` | Redirect to original URL | ❌ |

#### Analytics Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/v1/analytics/dashboard` | Dashboard statistics | ✅ |
| `GET` | `/api/v1/analytics/urls/{id}` | URL-specific analytics | ✅ |
| `GET` | `/api/v1/analytics/urls/{id}/timeline` | Click timeline | ✅ |
| `GET` | `/api/v1/analytics/urls/{id}/geo` | Geographic statistics | ✅ |
| `GET` | `/api/v1/analytics/export` | Export analytics data | ✅ |

#### QR Code Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/v1/qr/generate` | Generate QR code | Optional |
| `GET` | `/api/v1/qr/{shortCode}` | Get QR code for URL | ❌ |
| `GET` | `/api/v1/qr/formats` | Available QR formats | ❌ |

#### Health & Monitoring Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Overall health status |
| `GET` | `/health/livez` | Kubernetes liveness probe |
| `GET` | `/health/readyz` | Kubernetes readiness probe |
| `GET` | `/health/checks` | Detailed health checks |
| `GET` | `/health/metrics/system` | System metrics |
| `GET` | `/health/metrics/application` | Application metrics |

For complete API documentation, see:
- **[OpenAPI Specification](docs/api/openapi.yaml)**
- **[API Client Examples](docs/api/examples/)**
- **[Health Check Documentation](docs/health-checks.md)**

## 🔍 Development

### Development Workflow

#### Daily Development Commands
```bash
# Start development environment
./scripts/setup.sh dev
docker-compose up -d

# Run all tests
./scripts/test.sh all

# Check application health
curl http://localhost:8080/health/checks | jq

# View real-time logs
docker-compose logs -f backend frontend

# Stop environment
docker-compose down
```

#### Code Quality Standards

- **Test Coverage**: Maintain 95%+ test coverage for backend, 90%+ for frontend
- **Code Style**: Automated formatting with `gofmt`, `goimports`, and `prettier`
- **Error Handling**: Comprehensive error handling with structured logging
- **Security**: Input validation, SQL injection prevention, XSS protection
- **Performance**: Efficient database queries, proper indexing, caching strategies

#### Git Workflow

We use **Conventional Commits** for clear commit history:

```bash
# Feature development
git checkout -b feature/user-profile-management
git commit -m "feat: add user profile update functionality"

# Bug fixes
git commit -m "fix: resolve database connection pooling issue"

# Documentation
git commit -m "docs: update API documentation for QR endpoints"

# Tests
git commit -m "test: add integration tests for analytics service"
```

For complete development guidelines, see the [Developer Onboarding Guide](docs/developer-onboarding.md).

## 🚀 Deployment

### Production Deployment

#### Automated Deployment (Recommended)
```bash
# Deploy to production server
./scripts/deploy.sh --host your-server.com deploy

# Deploy with backup
./scripts/deploy.sh --host your-server.com --backup deploy

# Check deployment status
./scripts/deploy.sh --host your-server.com status

# Rollback if needed
./scripts/deploy.sh --host your-server.com rollback
```

#### Manual Deployment
```bash
# 1. Build production images
docker-compose -f docker-compose.prod.yml build

# 2. Deploy with production configuration
docker-compose -f docker-compose.prod.yml up -d

# 3. Run database migrations
docker-compose -f docker-compose.prod.yml exec backend go run cmd/migrate/main.go

# 4. Verify deployment
curl https://your-domain.com/health/checks
```

### Environment-Specific Configurations

- **Development**: Full debugging, hot reload, development tools
- **Staging**: Production-like environment for integration testing
- **Production**: Optimized builds, security hardening, monitoring

For comprehensive deployment instructions, see the [Production Deployment Guide](docs/production-deployment.md).

## 🔒 Security

### Security Features

#### Application Security
- **Authentication**: JWT-based authentication with secure refresh tokens
- **Authorization**: Role-based access control with granular permissions
- **Input Validation**: Comprehensive input sanitization and validation
- **Rate Limiting**: Intelligent API rate limiting and abuse prevention
- **Session Management**: Secure session handling with automatic expiration

#### Infrastructure Security
- **HTTPS/TLS**: End-to-end encryption with modern TLS protocols
- **CORS**: Configured Cross-Origin Resource Sharing
- **Security Headers**: HSTS, XSS protection, content type sniffing prevention
- **SQL Injection Protection**: Parameterized queries and ORM protection
- **XSS Protection**: Content Security Policy and output encoding

#### Operational Security
- **Environment Variables**: Secure secret management
- **Container Security**: Non-root users, read-only filesystems
- **Network Security**: Firewall configuration, VPC isolation
- **Backup Encryption**: Encrypted backups with secure key management

### Security Best Practices

- Regular security audits and penetration testing
- Dependency vulnerability scanning with automated updates
- Secure coding practices and security-focused code reviews
- Environment variable protection and secret rotation
- Regular security updates and patch management

## 📈 Performance & Monitoring

### Performance Features

#### Application Performance
- **Caching**: Multi-layer Redis caching for optimal response times
- **Database Optimization**: Proper indexing and query optimization
- **Connection Pooling**: Efficient database connection management
- **Compression**: Gzip compression for API responses
- **CDN Integration**: Static asset optimization and delivery

#### Monitoring & Observability
- **Health Checks**: Comprehensive health monitoring system
- **Metrics Collection**: Prometheus-compatible metrics
- **Real-time Monitoring**: Live performance dashboards
- **Error Tracking**: Structured error logging and alerting
- **Performance Profiling**: Built-in Go profiling endpoints

### Health Monitoring

The service includes a comprehensive health monitoring system:

```bash
# Check overall health
curl http://localhost:8080/health

# Detailed health checks
curl http://localhost:8080/health/checks | jq

# System metrics
curl http://localhost:8080/health/metrics/system | jq

# Application metrics
curl http://localhost:8080/health/metrics/application | jq
```

For complete monitoring setup, see the [Health Check Documentation](docs/health-checks.md).

## 🤝 Contributing

We welcome contributions! Please follow our contribution guidelines:

### Getting Started
1. **Read the Documentation**: Start with the [Developer Onboarding Guide](docs/developer-onboarding.md)
2. **Set up Development Environment**: Follow the setup instructions
3. **Check Existing Issues**: Look for open issues or create new ones
4. **Understand the Architecture**: Review the codebase structure

### Contribution Process
1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/your-feature`
3. **Develop** your changes following our coding standards
4. **Test** your changes thoroughly: `./scripts/test.sh all`
5. **Commit** your changes with conventional commits
6. **Push** to your branch: `git push origin feature/your-feature`
7. **Submit** a pull request with detailed description

### Development Guidelines

- **Code Quality**: Follow clean architecture principles and maintain test coverage
- **Testing**: Write comprehensive tests for new functionality
- **Documentation**: Update documentation for new features
- **Security**: Follow security best practices and review guidelines
- **Performance**: Ensure no significant performance regressions

For detailed contributing guidelines, see the [Developer Onboarding Guide](docs/developer-onboarding.md#contributing-guidelines).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support & Resources

### Getting Help

- **📖 Documentation**: Start with our comprehensive [Developer Onboarding Guide](docs/developer-onboarding.md)
- **🔍 API Reference**: Use the [Interactive API Documentation](http://localhost:8080/docs)
- **🏥 Health Monitoring**: Check the [Health Check Documentation](docs/health-checks.md)
- **🚀 Deployment**: Follow the [Production Deployment Guide](docs/production-deployment.md)
- **🐛 Issues**: Report bugs on [GitHub Issues](https://github.com/your-org/url-shortener/issues)

### Community & Support

- **📧 Email**: support@urlshortener.com
- **💬 Discussions**: [GitHub Discussions](https://github.com/your-org/url-shortener/discussions)
- **📱 Community**: Join our developer community
- **📰 Blog**: Follow our [development blog](https://blog.urlshortener.com) for updates

### Resources

- **🎯 Examples**: [API Client Examples](docs/api/examples/) in multiple languages
- **🔧 Scripts**: [Automation Scripts](scripts/README.md) for deployment and maintenance
- **📊 Monitoring**: [Performance monitoring](docs/health-checks.md) and alerting setup
- **🔐 Security**: Security best practices and vulnerability reporting

## 🗺️ Roadmap

### Current Development (Phase 1) ✅
- [x] Core URL shortening functionality
- [x] User authentication and management
- [x] Basic analytics and reporting
- [x] QR code generation
- [x] Comprehensive API documentation
- [x] Health monitoring system
- [x] Production deployment automation

### Upcoming Features (Phase 2) 🚧
- [ ] Advanced analytics with machine learning insights
- [ ] Team collaboration and workspace features
- [ ] Custom domains support
- [ ] Advanced security features (2FA, SSO)
- [ ] API integrations (Slack, Discord, webhooks)
- [ ] Mobile applications (iOS/Android)

### Future Enhancements (Phase 3) 📅
- [ ] Multi-language support (i18n)
- [ ] Advanced caching strategies and CDN integration
- [ ] Enterprise features (LDAP, audit logs)
- [ ] AI-powered link optimization
- [ ] Advanced reporting and business intelligence
- [ ] White-label solutions

## 📊 Current Status

**Development Progress**: **83% Complete** (Step 25 of 30 completed)

### Completed Milestones ✅
- ✅ **Project Foundation** (Steps 1-10): Complete architecture and core services
- ✅ **Frontend Development** (Steps 11-20): React application with full UI
- ✅ **Integration & Testing** (Steps 21-25): Comprehensive testing and documentation

### In Progress 🚧
- 🚧 **Documentation & Deployment** (Steps 26-30): Production optimization and final polish

### Recent Achievements 🎉
- ✅ **Comprehensive Health Monitoring**: Multi-level health checks with Kubernetes support
- ✅ **Interactive API Documentation**: Complete OpenAPI specification with Swagger UI
- ✅ **Production Deployment Automation**: Automated deployment scripts with rollback
- ✅ **Developer Onboarding Guide**: Complete documentation for new developers
- ✅ **Security Hardening**: Enterprise-grade security implementation

For detailed progress tracking, see [PROGRESS.md](PROGRESS.md).

---

<div align="center">

**🌟 Built with ❤️ using Go and React 🌟**

[⭐ Star us on GitHub](https://github.com/your-org/url-shortener) | [🐛 Report Issue](https://github.com/your-org/url-shortener/issues) | [💡 Request Feature](https://github.com/your-org/url-shortener/discussions)

</div>