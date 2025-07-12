# Developer Onboarding Guide

Welcome to the URL Shortener project! This guide will help you get up and running quickly with our development environment and understand the project structure, technologies, and development workflow.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Project Overview](#project-overview)
3. [Quick Start](#quick-start)
4. [Development Environment Setup](#development-environment-setup)
5. [Project Structure](#project-structure)
6. [Technology Stack](#technology-stack)
7. [Development Workflow](#development-workflow)
8. [Testing Strategy](#testing-strategy)
9. [Code Style and Standards](#code-style-and-standards)
10. [Common Tasks](#common-tasks)
11. [Debugging and Troubleshooting](#debugging-and-troubleshooting)
12. [API Documentation](#api-documentation)
13. [Contributing Guidelines](#contributing-guidelines)
14. [Resources and Further Reading](#resources-and-further-reading)

## Prerequisites

### Required Software

Before you begin, ensure you have the following installed on your development machine:

#### Essential Tools
- **Git** (2.30+) - Version control
- **Docker** (20.10+) and **Docker Compose** (2.0+) - Containerization
- **Go** (1.21+) - Backend development
- **Node.js** (18+) and **npm** (9+) - Frontend development
- **PostgreSQL** (15+) - Database (or use Docker)
- **Redis** (7+) - Cache (or use Docker)

#### Recommended Tools
- **VS Code** or **GoLand** - Code editor/IDE
- **Postman** or **Insomnia** - API testing
- **pgAdmin** or **TablePlus** - Database management
- **Git GUI** client (optional) - GitKraken, SourceTree, etc.

#### Optional Tools
- **Task** - Task runner (alternative to make)
- **golangci-lint** - Go linting
- **gosec** - Go security scanner

### System Requirements

- **OS**: Linux, macOS, or Windows with WSL2
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 10GB free space for development environment
- **Network**: Stable internet connection for downloading dependencies

### Account Setup

You'll need accounts on:
- **GitHub** - Code repository access
- **Docker Hub** (optional) - For pushing/pulling container images

## Project Overview

The URL Shortener is a production-grade service built with modern technologies and best practices. It provides:

### Core Features
- **URL Shortening** - Convert long URLs to short, shareable links
- **Custom Aliases** - User-defined short codes
- **Password Protection** - Secure URLs with passwords
- **Expiration Dates** - Time-limited URLs
- **Analytics Dashboard** - Comprehensive click tracking and reporting
- **QR Code Generation** - Dynamic QR codes for shortened URLs
- **User Management** - Authentication and user profiles

### Key Characteristics
- **Microservices Architecture** - Clean separation of concerns
- **RESTful API** - Standard HTTP API design
- **Real-time Analytics** - Live click tracking and metrics
- **Scalable Design** - Horizontal scaling capabilities
- **Security-First** - Built with security best practices
- **Cloud-Ready** - Containerized and deployment-ready

## Quick Start

Get the project running in under 5 minutes:

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/url-shortener.git
cd url-shortener
```

### 2. Automated Setup

```bash
# Make setup script executable
chmod +x scripts/setup.sh

# Run automated development setup
./scripts/setup.sh dev
```

### 3. Start Development Environment

```bash
# Start all services with Docker Compose
docker-compose up -d

# Or use the convenience script
./scripts/deploy.sh status
```

### 4. Verify Installation

```bash
# Check backend health
curl http://localhost:8080/health

# Check frontend (if running)
curl http://localhost:3000

# Run tests
./scripts/test.sh all
```

### 5. Access the Application

- **API Documentation**: http://localhost:8080/docs
- **API Base URL**: http://localhost:8080/api/v1
- **Frontend** (when available): http://localhost:3000
- **Database**: localhost:5432 (postgres/password)
- **Redis**: localhost:6379

## Development Environment Setup

### Method 1: Docker-Based Development (Recommended)

This is the fastest way to get started and ensures environment consistency.

```bash
# 1. Start infrastructure services
docker-compose up -d postgres redis

# 2. Install backend dependencies
cd backend && go mod download

# 3. Install frontend dependencies
cd ../frontend && npm install

# 4. Start backend in development mode
cd ../backend && go run cmd/server/main.go

# 5. Start frontend in development mode (in another terminal)
cd frontend && npm run dev
```

### Method 2: Native Development

For developers who prefer running services natively:

#### Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Set up environment
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
go run cmd/migrate/main.go

# Start the server
go run cmd/server/main.go
```

#### Frontend Setup

```bash
cd frontend

# Install Node.js dependencies
npm install

# Set up environment
cp .env.example .env
# Edit .env with your configuration

# Start development server
npm run dev
```

### Method 3: VS Code Dev Containers

For VS Code users, we provide a complete dev container setup:

```bash
# Install the Dev Containers extension in VS Code
# Open the project in VS Code
# Press F1 and select "Dev Containers: Reopen in Container"
```

## Project Structure

Understanding the project structure is crucial for effective development:

```
url-shortener/
├── README.md                   # Project overview and quick start
├── PROGRESS.md                 # Development progress tracking
├── CLAUDE.md                   # AI assistant instructions
├── Taskfile.yml               # Task automation
├── docker-compose.yml         # Development Docker setup
├── docker-compose.prod.yml    # Production Docker setup
├── .env.example               # Environment template
├── .gitignore                 # Git ignore rules
│
├── backend/                   # Go backend service
│   ├── cmd/                   # Application entry points
│   │   └── server/           
│   │       └── main.go       # Main application
│   ├── internal/             # Private application code
│   │   ├── api/              # HTTP layer
│   │   │   ├── handlers/     # HTTP handlers
│   │   │   ├── middleware/   # HTTP middleware
│   │   │   └── routes/       # Route configuration
│   │   ├── core/             # Business logic
│   │   │   ├── domain/       # Domain models
│   │   │   ├── ports/        # Interface definitions
│   │   │   └── services/     # Business services
│   │   ├── infrastructure/   # External concerns
│   │   │   ├── database/     # Database layer
│   │   │   ├── cache/        # Cache layer
│   │   │   └── external/     # External services
│   │   └── config/           # Configuration management
│   ├── tests/                # Test files
│   ├── Dockerfile            # Backend container definition
│   ├── go.mod               # Go module definition
│   └── .env.example         # Backend environment template
│
├── frontend/                 # React TypeScript frontend
│   ├── public/              # Static assets
│   ├── src/                 # Source code
│   │   ├── components/      # React components
│   │   ├── pages/           # Page components
│   │   ├── hooks/           # Custom React hooks
│   │   ├── services/        # API services
│   │   ├── utils/           # Utility functions
│   │   ├── types/           # TypeScript type definitions
│   │   └── context/         # React context providers
│   ├── tests/               # Frontend tests
│   ├── Dockerfile           # Frontend container definition
│   ├── package.json         # NPM dependencies
│   ├── tsconfig.json        # TypeScript configuration
│   ├── tailwind.config.js   # Tailwind CSS configuration
│   └── vite.config.ts       # Vite build configuration
│
├── docs/                    # Documentation
│   ├── api/                 # API documentation
│   │   ├── openapi.yaml     # OpenAPI specification
│   │   ├── swagger-ui.html  # Interactive API docs
│   │   └── examples/        # API client examples
│   ├── deployment/          # Deployment guides
│   ├── health-checks.md     # Health check documentation
│   └── developer-onboarding.md  # This file
│
└── scripts/                 # Automation scripts
    ├── setup.sh             # Environment setup
    ├── deploy.sh            # Deployment automation
    ├── test.sh              # Testing automation
    ├── backup.sh            # Backup and maintenance
    └── README.md            # Scripts documentation
```

### Key Architecture Principles

#### Backend Architecture (Clean Architecture)

```
┌─────────────────┐
│   HTTP Layer    │  ← api/handlers, api/routes, api/middleware
├─────────────────┤
│ Business Logic  │  ← core/services, core/domain
├─────────────────┤
│ Infrastructure  │  ← infrastructure/database, infrastructure/cache
└─────────────────┘
```

- **Domain Layer**: Pure business logic and entities
- **Service Layer**: Use cases and business rules
- **Infrastructure Layer**: Database, cache, external services
- **API Layer**: HTTP handlers and middleware

#### Frontend Architecture (Component-Based)

```
┌─────────────────┐
│     Pages       │  ← pages/ (route-level components)
├─────────────────┤
│   Components    │  ← components/ (reusable UI components)
├─────────────────┤
│   Services      │  ← services/ (API communication)
├─────────────────┤
│   State/Context │  ← context/ (global state management)
└─────────────────┘
```

## Technology Stack

### Backend Technologies

#### Core Framework
- **Go 1.21+** - Programming language
- **Chi Router** - HTTP router and middleware
- **GORM** - Object-relational mapping
- **JWT** - Authentication tokens

#### Database and Cache
- **PostgreSQL 15+** - Primary database
- **Redis 7+** - Caching and session storage

#### Development Tools
- **Air** - Live reloading for Go development
- **golangci-lint** - Go linting
- **testify** - Testing framework
- **godotenv** - Environment management

#### Production Tools
- **Docker** - Containerization
- **Prometheus** - Metrics collection
- **Grafana** - Monitoring dashboards

### Frontend Technologies

#### Core Framework
- **React 18** - UI framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Build tool and dev server

#### UI and Styling
- **Tailwind CSS** - Utility-first CSS framework
- **Headless UI** - Unstyled, accessible UI components
- **Lucide React** - Icon library

#### State Management
- **React Context** - Global state management
- **React Hook Form** - Form management
- **Zod** - Schema validation

#### Development Tools
- **Vitest** - Unit testing framework
- **React Testing Library** - Component testing
- **ESLint** - JavaScript linting
- **Prettier** - Code formatting

### Infrastructure Technologies

#### Containerization
- **Docker** - Application containerization
- **Docker Compose** - Multi-container orchestration

#### Monitoring and Observability
- **Prometheus** - Metrics collection
- **Grafana** - Visualization and dashboards
- **Health Checks** - Application monitoring

#### Development Tools
- **Task** - Task runner
- **Git** - Version control
- **GitHub Actions** - CI/CD (planned)

## Development Workflow

### Daily Development Process

#### 1. Start Your Development Session

```bash
# Pull latest changes
git pull origin main

# Start development environment
docker-compose up -d postgres redis

# Start backend
cd backend && air # or go run cmd/server/main.go

# Start frontend (in another terminal)
cd frontend && npm run dev
```

#### 2. Make Changes

- Work on features in focused commits
- Write tests as you develop
- Use the health check endpoint to verify your changes
- Test API endpoints using the interactive documentation

#### 3. Test Your Changes

```bash
# Run backend tests
cd backend && go test ./...

# Run frontend tests
cd frontend && npm test

# Run all tests with coverage
./scripts/test.sh all --coverage

# Lint your code
./scripts/test.sh lint
```

#### 4. Commit Your Work

```bash
# Add your changes
git add .

# Commit with descriptive message
git commit -m "feat: add user profile management functionality"

# Push your changes
git push origin feature/user-profiles
```

### Git Workflow

We follow a **GitHub Flow** approach:

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Commits**
   - Use conventional commit messages
   - Keep commits focused and atomic
   - Commit frequently (every 30-60 minutes of work)

3. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   # Create pull request on GitHub
   ```

4. **Code Review and Merge**
   - Address review feedback
   - Ensure CI passes
   - Squash and merge to main

### Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting)
- `refactor:` - Code refactoring
- `test:` - Adding or modifying tests
- `chore:` - Build process or auxiliary tool changes

**Examples:**
```bash
feat: add password protection for shortened URLs
fix: resolve database connection pooling issue
docs: update API documentation for QR code endpoints
test: add integration tests for analytics service
```

## Testing Strategy

We maintain a comprehensive testing strategy with multiple test types:

### Backend Testing

#### Unit Tests
```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/core/services/...

# Run with verbose output
go test -v ./...
```

#### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./tests/integration/...

# Run with database
./scripts/test.sh integration
```

#### Test Structure
```go
// Example unit test
func TestURLService_ShortenURL(t *testing.T) {
    // Arrange
    mockRepo := &mocks.URLRepository{}
    service := NewURLService(mockRepo, nil, nil)
    
    // Act
    result, err := service.ShortenURL(ctx, request)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ShortCode)
}
```

### Frontend Testing

#### Unit Tests
```bash
# Run all frontend tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch
```

#### Component Tests
```typescript
// Example component test
import { render, screen, fireEvent } from '@testing-library/react';
import { URLShortener } from './URLShortener';

test('should shorten URL when form is submitted', async () => {
  render(<URLShortener />);
  
  fireEvent.change(screen.getByLabelText(/original url/i), {
    target: { value: 'https://example.com/very-long-url' }
  });
  
  fireEvent.click(screen.getByText(/shorten/i));
  
  expect(await screen.findByText(/shortened url/i)).toBeInTheDocument();
});
```

### Test Coverage Requirements

- **Backend**: Minimum 95% coverage
- **Frontend**: Minimum 90% coverage
- **Integration**: Cover all critical user workflows
- **E2E**: Cover main user journeys (planned)

## Code Style and Standards

### Go Code Style

#### Formatting
```bash
# Format code
go fmt ./...

# Import organization
goimports -w .

# Linting
golangci-lint run
```

#### Conventions
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful variable and function names
- Keep functions small and focused
- Use interfaces for testability

#### Example Code Style
```go
// Good: Clear, descriptive naming
func (s *URLService) ShortenURL(ctx context.Context, req domain.ShortenURLRequest) (*domain.ShortURL, error) {
    if err := s.validateURL(req.OriginalURL); err != nil {
        return nil, fmt.Errorf("invalid URL: %w", err)
    }
    
    shortCode := s.generateShortCode()
    shortURL := &domain.ShortURL{
        OriginalURL: req.OriginalURL,
        ShortCode:   shortCode,
        UserID:      req.UserID,
        CreatedAt:   time.Now(),
    }
    
    if err := s.urlRepo.Create(ctx, shortURL); err != nil {
        return nil, fmt.Errorf("failed to save URL: %w", err)
    }
    
    return shortURL, nil
}
```

### TypeScript/React Code Style

#### Formatting
```bash
# Format code
npm run format

# Linting
npm run lint

# Type checking
npm run type-check
```

#### Conventions
- Use TypeScript strict mode
- Prefer functional components with hooks
- Use proper TypeScript types
- Follow React best practices
- Use meaningful component and prop names

#### Example Code Style
```typescript
// Good: Proper TypeScript types and React patterns
interface URLShortenerProps {
  onURLShortened: (shortURL: ShortURL) => void;
  isLoading?: boolean;
}

export const URLShortener: React.FC<URLShortenerProps> = ({ 
  onURLShortened, 
  isLoading = false 
}) => {
  const [originalURL, setOriginalURL] = useState('');
  const [error, setError] = useState<string | null>(null);
  
  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError(null);
    
    try {
      const shortURL = await urlService.shortenURL({ originalURL });
      onURLShortened(shortURL);
      setOriginalURL('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    }
  };
  
  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Form implementation */}
    </form>
  );
};
```

### Database Conventions

#### Migration Naming
```
YYYY_MM_DD_HHMMSS_description.sql
```

#### Table Naming
- Use snake_case
- Use plural nouns for table names
- Use descriptive names

#### Column Naming
- Use snake_case
- Be descriptive but concise
- Use consistent naming patterns

## Common Tasks

### Creating a New API Endpoint

#### 1. Define Domain Model
```go
// internal/core/domain/user.go
type UpdateProfileRequest struct {
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
    Email     string `json:"email" validate:"required,email"`
}
```

#### 2. Add Service Method
```go
// internal/core/services/auth.go
func (s *AuthService) UpdateProfile(ctx context.Context, userID uint, req domain.UpdateProfileRequest) (*domain.UserResponse, error) {
    // Implementation
}
```

#### 3. Create Handler
```go
// internal/api/handlers/auth.go
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

#### 4. Add Route
```go
// internal/api/routes/routes.go
protectedRouter.Put("/profile", r.config.AuthHandler.UpdateProfile)
```

#### 5. Write Tests
```go
func TestAuthHandler_UpdateProfile(t *testing.T) {
    // Test implementation
}
```

#### 6. Update API Documentation
```yaml
# docs/api/openapi.yaml
/auth/profile:
  put:
    # OpenAPI specification
```

### Adding a New React Component

#### 1. Create Component File
```typescript
// src/components/user/ProfileEditor.tsx
import React from 'react';

interface ProfileEditorProps {
  user: User;
  onSave: (profile: UpdateProfileRequest) => Promise<void>;
}

export const ProfileEditor: React.FC<ProfileEditorProps> = ({ user, onSave }) => {
  // Component implementation
};
```

#### 2. Add to Index
```typescript
// src/components/user/index.ts
export { ProfileEditor } from './ProfileEditor';
```

#### 3. Write Tests
```typescript
// src/components/user/ProfileEditor.test.tsx
import { render, screen } from '@testing-library/react';
import { ProfileEditor } from './ProfileEditor';

describe('ProfileEditor', () => {
  it('should render user information', () => {
    // Test implementation
  });
});
```

#### 4. Add to Storybook (if available)
```typescript
// src/components/user/ProfileEditor.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { ProfileEditor } from './ProfileEditor';
```

### Database Migration

#### 1. Create Migration File
```bash
cd backend
migrate create -ext sql -dir migrations add_user_preferences
```

#### 2. Write Up Migration
```sql
-- migrations/XXXX_add_user_preferences.up.sql
CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
```

#### 3. Write Down Migration
```sql
-- migrations/XXXX_add_user_preferences.down.sql
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP TABLE IF EXISTS user_preferences;
```

#### 4. Run Migration
```bash
go run cmd/migrate/main.go
```

### Adding Environment Variables

#### 1. Add to .env.example
```bash
# Backend environment
FEATURE_USER_PREFERENCES=true
USER_PREFERENCES_CACHE_TTL=3600
```

#### 2. Update Config Structure
```go
// internal/config/config.go
type Config struct {
    // ...existing fields
    FeatureUserPreferences  bool          `env:"FEATURE_USER_PREFERENCES" envDefault:"false"`
    UserPreferencesCacheTTL time.Duration `env:"USER_PREFERENCES_CACHE_TTL" envDefault:"1h"`
}
```

#### 3. Use in Code
```go
if config.FeatureUserPreferences {
    // Feature implementation
}
```

## Debugging and Troubleshooting

### Backend Debugging

#### Using Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug ./cmd/server/main.go

# Or attach to running process
dlv attach <pid>
```

#### Logging
```go
// Use structured logging
import "log/slog"

slog.Info("Processing request", 
    "user_id", userID, 
    "url", originalURL,
    "remote_addr", r.RemoteAddr)

slog.Error("Database error", 
    "error", err,
    "operation", "create_short_url")
```

#### Common Issues

**Database Connection Issues**
```bash
# Check database status
docker-compose ps postgres
docker-compose logs postgres

# Test connection
psql -h localhost -p 5432 -U urlshortener -d urlshortener
```

**Cache Connection Issues**
```bash
# Check Redis status
docker-compose ps redis
docker-compose logs redis

# Test connection
redis-cli -h localhost -p 6379 ping
```

**Port Conflicts**
```bash
# Check what's using port 8080
lsof -i :8080

# Kill process using port
kill -9 <pid>
```

### Frontend Debugging

#### Browser DevTools
- Use React Developer Tools extension
- Check Network tab for API calls
- Use Console for debugging

#### React Query DevTools (if using)
```typescript
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

function App() {
  return (
    <>
      {/* Your app */}
      <ReactQueryDevtools initialIsOpen={false} />
    </>
  );
}
```

#### Common Issues

**API Connection Issues**
```typescript
// Check API base URL
console.log('API Base URL:', import.meta.env.VITE_API_BASE_URL);

// Test API endpoint
fetch('http://localhost:8080/health')
  .then(r => r.json())
  .then(console.log);
```

**Build Issues**
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vite cache
rm -rf node_modules/.vite
```

### Health Check Debugging

#### Check All Health Endpoints
```bash
# Basic health
curl http://localhost:8080/health

# Detailed checks
curl http://localhost:8080/health/checks | jq

# Component-specific
curl http://localhost:8080/health/database
curl http://localhost:8080/health/cache

# Metrics
curl http://localhost:8080/health/metrics/system | jq
```

#### Monitor Health in Development
```bash
# Watch health status
watch -n 5 'curl -s http://localhost:8080/health | jq .status'

# Monitor database connections
watch -n 2 'curl -s http://localhost:8080/health/database | jq .metadata'
```

## API Documentation

### Interactive Documentation

The project includes comprehensive API documentation accessible at:
- **Development**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/api/openapi.yaml

### Using the API

#### Authentication
```bash
# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'

# Use token in subsequent requests
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### URL Operations
```bash
# Shorten URL
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://example.com/very-long-url"}'

# Get user URLs
curl -X GET http://localhost:8080/api/v1/urls \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Access shortened URL
curl -X GET http://localhost:8080/abc123
```

### API Client Examples

The `docs/api/examples/` directory contains client examples in multiple languages:
- JavaScript/TypeScript
- Python
- Go
- cURL commands

## Contributing Guidelines

### Before Contributing

1. **Read the documentation** - Understand the project structure and conventions
2. **Set up your environment** - Follow the development setup guide
3. **Run tests** - Ensure your environment is working correctly
4. **Check existing issues** - See if your contribution aligns with project goals

### Making Contributions

#### 1. Fork and Clone
```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/url-shortener.git
cd url-shortener
```

#### 2. Create Feature Branch
```bash
git checkout -b feature/your-feature-name
```

#### 3. Make Changes
- Follow coding standards
- Write tests for new functionality
- Update documentation as needed
- Ensure all tests pass

#### 4. Test Your Changes
```bash
# Run all tests
./scripts/test.sh all

# Check code quality
./scripts/test.sh lint

# Verify health checks work
curl http://localhost:8080/health/checks
```

#### 5. Submit Pull Request
- Write clear commit messages
- Provide detailed PR description
- Reference related issues
- Ensure CI checks pass

### Code Review Process

1. **Automated Checks** - CI pipeline runs tests and quality checks
2. **Peer Review** - Team members review code for quality and standards
3. **Testing** - Verify functionality works as expected
4. **Documentation** - Ensure documentation is updated
5. **Approval and Merge** - Maintainer approves and merges PR

### What We Look For

- **Code Quality** - Clean, readable, maintainable code
- **Test Coverage** - Comprehensive tests for new functionality
- **Documentation** - Clear documentation for new features
- **Performance** - No significant performance regressions
- **Security** - Secure coding practices
- **Standards Compliance** - Follows project conventions

## Resources and Further Reading

### Project Documentation
- [API Documentation](./api/) - Complete API reference
- [Health Checks Guide](./health-checks.md) - Health monitoring system
- [Deployment Guide](./deployment/) - Production deployment
- [Scripts Documentation](../scripts/README.md) - Automation scripts

### Technology Documentation

#### Backend (Go)
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Chi Router](https://go-chi.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [testify Testing](https://github.com/stretchr/testify)

#### Frontend (React/TypeScript)
- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Vite Guide](https://vitejs.dev/guide/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [React Hook Form](https://react-hook-form.com/)

#### Infrastructure
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/docs/)

### Best Practices Guides
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [API Design Best Practices](https://restfulapi.net/)
- [React Best Practices](https://react.dev/learn/thinking-in-react)
- [Testing Best Practices](https://kentcdodds.com/blog/common-mistakes-with-react-testing-library)

### Development Tools
- [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [React Developer Tools](https://react.dev/learn/react-developer-tools)
- [Postman](https://www.postman.com/)
- [pgAdmin](https://www.pgadmin.org/)

## Getting Help

### Internal Resources
- **Documentation** - Check this guide and other docs first
- **Code Comments** - Read inline documentation in the codebase
- **Tests** - Look at existing tests for usage examples
- **Health Checks** - Use `/health/checks` to diagnose issues

### Community Resources
- **GitHub Issues** - Search existing issues or create new ones
- **Discussions** - Use GitHub Discussions for questions
- **Code Review** - Ask for help in pull request reviews

### Quick Commands Reference
```bash
# Start development environment
./scripts/setup.sh dev && docker-compose up -d

# Run all tests
./scripts/test.sh all

# Check application health
curl http://localhost:8080/health/checks | jq

# View logs
docker-compose logs -f backend

# Reset environment
docker-compose down -v && docker-compose up -d

# Deploy to production
./scripts/deploy.sh --host server.example.com deploy
```

---

Welcome to the team! We're excited to have you contribute to the URL Shortener project. If you have any questions or need assistance, don't hesitate to reach out through our communication channels.

Happy coding! 🚀