# CLAUDE.md - URL Shortener Service Development Guide

## Project Overview

You are tasked with building a production-grade URL Shortener Service with comprehensive analytics, QR code generation, and advanced link management features. This is a full-stack application with a Golang backend using Chi router and a React TypeScript frontend.

## Repository Structure

Create the following directory structure:

```
url-shortener/
├── README.md
├── PROGRESS.md
├── Taskfile.yml
├── docker-compose.yml
├── .gitignore
├── .env.example
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   │   ├── auth.go
│   │   │   │   ├── url.go
│   │   │   │   ├── analytics.go
│   │   │   │   └── qr.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── logging.go
│   │   │   │   └── ratelimit.go
│   │   │   └── routes/
│   │   │       └── routes.go
│   │   ├── core/
│   │   │   ├── domain/
│   │   │   │   ├── user.go
│   │   │   │   ├── url.go
│   │   │   │   ├── click.go
│   │   │   │   └── errors.go
│   │   │   ├── ports/
│   │   │   │   ├── repositories.go
│   │   │   │   └── services.go
│   │   │   └── services/
│   │   │       ├── auth.go
│   │   │       ├── url.go
│   │   │       ├── analytics.go
│   │   │       └── qr.go
│   │   ├── infrastructure/
│   │   │   ├── database/
│   │   │   │   ├── postgres.go
│   │   │   │   ├── migrations/
│   │   │   │   └── repositories/
│   │   │   │       ├── user.go
│   │   │   │       ├── url.go
│   │   │   │       └── click.go
│   │   │   ├── cache/
│   │   │   │   └── redis.go
│   │   │   └── external/
│   │   │       ├── geolocation.go
│   │   │       └── qrcode.go
│   │   └── config/
│   │       └── config.go
│   ├── tests/
│   │   ├── integration/
│   │   ├── unit/
│   │   └── testdata/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   │   ├── common/
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Footer.tsx
│   │   │   │   ├── Layout.tsx
│   │   │   │   └── Loading.tsx
│   │   │   ├── auth/
│   │   │   │   ├── LoginForm.tsx
│   │   │   │   ├── RegisterForm.tsx
│   │   │   │   └── PasswordReset.tsx
│   │   │   ├── url/
│   │   │   │   ├── URLShortener.tsx
│   │   │   │   ├── URLList.tsx
│   │   │   │   ├── URLCard.tsx
│   │   │   │   └── URLDetails.tsx
│   │   │   ├── analytics/
│   │   │   │   ├── Dashboard.tsx
│   │   │   │   ├── ClickChart.tsx
│   │   │   │   ├── GeographicMap.tsx
│   │   │   │   └── DeviceStats.tsx
│   │   │   └── qr/
│   │   │       ├── QRGenerator.tsx
│   │   │       └── QRPreview.tsx
│   │   ├── pages/
│   │   │   ├── Home.tsx
│   │   │   ├── Dashboard.tsx
│   │   │   ├── Analytics.tsx
│   │   │   ├── Profile.tsx
│   │   │   └── NotFound.tsx
│   │   ├── hooks/
│   │   │   ├── useAuth.ts
│   │   │   ├── useAPI.ts
│   │   │   └── useLocalStorage.ts
│   │   ├── services/
│   │   │   ├── api.ts
│   │   │   ├── auth.ts
│   │   │   └── urls.ts
│   │   ├── utils/
│   │   │   ├── validation.ts
│   │   │   ├── formatting.ts
│   │   │   └── constants.ts
│   │   ├── types/
│   │   │   ├── auth.ts
│   │   │   ├── url.ts
│   │   │   └── analytics.ts
│   │   ├── context/
│   │   │   └── AuthContext.tsx
│   │   ├── App.tsx
│   │   ├── index.tsx
│   │   └── index.css
│   ├── tests/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/
│   │   └── utils/
│   ├── Dockerfile
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   ├── vite.config.ts
│   └── .env.example
└── scripts/
    ├── setup.sh
    ├── test.sh
    └── deploy.sh
```

## Technology Stack

### Backend (Golang)

- **Router**: Chi (github.com/go-chi/chi/v5)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Validation**: go-playground/validator/v10
- **Testing**: testify/suite, testify/mock
- **Environment**: godotenv
- **Logging**: logrus or slog
- **HTTP Client**: net/http (standard library)
- **QR Code**: github.com/skip2/go-qrcode
- **Rate Limiting**: golang.org/x/time/rate
- **Password Hashing**: golang.org/x/crypto/bcrypt

### Frontend (React TypeScript)

- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v6
- **State Management**: Context API + useReducer
- **HTTP Client**: Axios
- **UI Framework**: Tailwind CSS + Headless UI
- **Charts**: Recharts
- **Forms**: React Hook Form + Zod validation
- **QR Codes**: qrcode.js
- **Icons**: Lucide React
- **Testing**: Vitest + React Testing Library
- **Date Handling**: date-fns

### Infrastructure

- **Containerization**: Docker + Docker Compose
- **Task Runner**: Task (Taskfile.yml)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Reverse Proxy**: Nginx (for production)

## Development Guidelines

### Code Quality Standards

1. **Test Coverage**: Maintain 95%+ test coverage for both frontend and backend
2. **Code Style**: Use gofmt, golint, and Prettier with consistent formatting
3. **Error Handling**: Comprehensive error handling with proper logging
4. **Documentation**: Godoc comments for all public functions and interfaces
5. **Security**: Input validation, SQL injection prevention, XSS protection
6. **Performance**: Efficient database queries, proper indexing, caching strategies

### Git Workflow

- **Commit Frequency**: Commit after each logical unit of work (every 30-60 minutes)
- **Commit Messages**: Use conventional commits (feat:, fix:, docs:, test:, refactor:)
- **Alway follow includeCoAuthoredBy = false. i.e Never include Claude or any AI attribution in git commits. Remain anonymous at all times in regards to git. Do not use co-authored-by tags or any other form of attribution.**
- **Branch Strategy**: Work on main branch for this project
- **Progress Tracking**: Update PROGRESS.md after every major step

### Dependency Injection Pattern

```go
// Example service interface
type URLService interface {
    ShortenURL(ctx context.Context, req ShortenURLRequest) (*ShortURL, error)
    GetOriginalURL(ctx context.Context, shortCode string) (*ShortURL, error)
}

// Implementation with injected dependencies
type urlService struct {
    urlRepo URLRepository
    cache   CacheService
    logger  Logger
}

func NewURLService(urlRepo URLRepository, cache CacheService, logger Logger) URLService {
    return &urlService{
        urlRepo: urlRepo,
        cache:   cache,
        logger:  logger,
    }
}
```

## Step-by-Step Implementation Plan

### Phase 1: Project Foundation (Steps 1-10)

#### Step 1: Initialize Repository Structure

1. Create the repository structure as defined above
2. Initialize git repository
3. Create initial .gitignore file
4. Set up go.mod and package.json files
5. Create basic Taskfile.yml with common tasks
6. Update PROGRESS.md with completion status

#### Step 2: Environment Configuration

1. Create .env.example files for both backend and frontend
2. Set up configuration management in backend/internal/config/config.go
3. Use environment-based configuration loading
4. Create docker-compose.yml for local development
5. Update PROGRESS.md

#### Step 3: Database Setup

1. Set up PostgreSQL connection in infrastructure/database/postgres.go
2. Create migration files for all tables (users, short_urls, clicks)
3. Implement GORM models in core/domain/
4. Create database initialization scripts
5. Write tests for database connections
6. Update PROGRESS.md

#### Step 4: Redis Cache Setup

1. Implement Redis connection in infrastructure/cache/redis.go
2. Create cache service interface and implementation
3. Add caching for URL lookups and rate limiting
4. Write comprehensive tests for cache operations
5. Update PROGRESS.md

#### Step 5: Core Domain Models

1. Define all domain structs in core/domain/ (User, ShortURL, Click)
2. Create custom error types in core/domain/errors.go
3. Define repository interfaces in core/ports/repositories.go
4. Define service interfaces in core/ports/services.go
5. Write unit tests for domain logic
6. Update PROGRESS.md

#### Step 6: Repository Layer Implementation

1. Implement UserRepository in infrastructure/database/repositories/user.go
2. Implement URLRepository in infrastructure/database/repositories/url.go
3. Implement ClickRepository in infrastructure/database/repositories/click.go
4. Write comprehensive unit tests with mocked database
5. Achieve 95%+ test coverage for repository layer
6. Update PROGRESS.md

#### Step 7: Service Layer Implementation

1. Implement AuthService in core/services/auth.go with JWT
2. Implement URLService in core/services/url.go with shortening logic
3. Implement AnalyticsService in core/services/analytics.go
4. Implement QRService in core/services/qr.go
5. Write extensive unit tests with mocked dependencies
6. Update PROGRESS.md

#### Step 8: Middleware Implementation

1. Create authentication middleware in api/middleware/auth.go
2. Implement CORS middleware in api/middleware/cors.go
3. Create logging middleware in api/middleware/logging.go
4. Implement rate limiting middleware in api/middleware/ratelimit.go
5. Write tests for all middleware
6. Update PROGRESS.md

#### Step 9: HTTP Handlers Implementation

1. Implement auth handlers in api/handlers/auth.go (login, register, refresh)
2. Implement URL handlers in api/handlers/url.go (create, list, update, delete)
3. Implement analytics handlers in api/handlers/analytics.go
4. Implement QR code handlers in api/handlers/qr.go
5. Write comprehensive handler tests with mocked services
6. Update PROGRESS.md

#### Step 10: API Routes Setup

1. Set up Chi router in api/routes/routes.go
2. Configure all routes with proper middleware
3. Add request/response validation
4. Implement proper error handling and responses
5. Write integration tests for all endpoints
6. Update PROGRESS.md

### Phase 2: Frontend Foundation (Steps 11-20)

#### Step 11: React Project Setup

1. Initialize React TypeScript project with Vite
2. Configure Tailwind CSS and Headless UI
3. Set up React Router for client-side routing
4. Configure Axios for API communication
5. Set up testing environment with Vitest
6. Update PROGRESS.md

#### Step 12: Authentication Context & Services

1. Create AuthContext in context/AuthContext.tsx
2. Implement auth service in services/auth.ts
3. Create useAuth hook in hooks/useAuth.ts
4. Add token management and automatic refresh
5. Write comprehensive tests for auth functionality
6. Update PROGRESS.md

#### Step 13: API Service Layer

1. Create base API service in services/api.ts with interceptors
2. Implement URL service in services/urls.ts
3. Add error handling and response transformation
4. Create useAPI hook for common API operations
5. Write tests for all service functions
6. Update PROGRESS.md

#### Step 14: Type Definitions

1. Define auth types in types/auth.ts
2. Define URL types in types/url.ts
3. Define analytics types in types/analytics.ts
4. Ensure type safety across the application
5. Update PROGRESS.md

#### Step 15: Common Components

1. Create Layout component in components/common/Layout.tsx
2. Implement Header with navigation in components/common/Header.tsx
3. Create Loading component in components/common/Loading.tsx
4. Add Footer component in components/common/Footer.tsx
5. Write tests for all common components
6. Update PROGRESS.md

#### Step 16: Authentication Components

1. Create LoginForm in components/auth/LoginForm.tsx
2. Implement RegisterForm in components/auth/RegisterForm.tsx
3. Add PasswordReset in components/auth/PasswordReset.tsx
4. Add form validation with React Hook Form + Zod
5. Write comprehensive component tests
6. Update PROGRESS.md

#### Step 17: URL Management Components

1. Create URLShortener in components/url/URLShortener.tsx
2. Implement URLList in components/url/URLList.tsx
3. Add URLCard in components/url/URLCard.tsx
4. Create URLDetails in components/url/URLDetails.tsx
5. Add copy-to-clipboard and QR code features
6. Write tests for all URL components
7. Update PROGRESS.md

#### Step 18: Analytics Components

1. Create Dashboard in components/analytics/Dashboard.tsx
2. Implement ClickChart in components/analytics/ClickChart.tsx with Recharts
3. Add GeographicMap in components/analytics/GeographicMap.tsx
4. Create DeviceStats in components/analytics/DeviceStats.tsx
5. Add real-time data updates
6. Write tests for analytics components
7. Update PROGRESS.md

#### Step 19: Pages Implementation

1. Create Home page in pages/Home.tsx
2. Implement Dashboard page in pages/Dashboard.tsx
3. Add Analytics page in pages/Analytics.tsx
4. Create Profile page in pages/Profile.tsx
5. Add NotFound page in pages/NotFound.tsx
6. Write tests for all pages
7. Update PROGRESS.md

#### Step 20: QR Code Components

1. Create QRGenerator in components/qr/QRGenerator.tsx
2. Implement QRPreview in components/qr/QRPreview.tsx
3. Add download functionality for QR codes
4. Support multiple formats (PNG, SVG)
5. Write tests for QR components
6. Update PROGRESS.md

### Phase 3: Integration & Testing (Steps 21-25)

#### Step 21: Backend Integration Tests

1. Create integration test suite in tests/integration/
2. Test complete API workflows (auth, URL creation, analytics)
3. Test database transactions and rollbacks
4. Test cache invalidation scenarios
5. Achieve 95%+ coverage including integration tests
6. Update PROGRESS.md

#### Step 22: Frontend Integration Tests

1. Create end-to-end component tests
2. Test user workflows (login, create URL, view analytics)
3. Test error handling and edge cases
4. Mock API responses for consistent testing
5. Achieve 95%+ frontend test coverage
6. Update PROGRESS.md

#### Step 23: Performance Optimization

1. Add database indexes for optimal query performance
2. Implement proper caching strategies
3. Optimize frontend bundle size and loading
4. Add connection pooling for database
5. Implement proper error boundaries
6. Update PROGRESS.md

#### Step 24: Security Implementation

1. Add input validation and sanitization
2. Implement rate limiting on all endpoints
3. Add CORS configuration
4. Implement proper password hashing
5. Add XSS and CSRF protection
6. Update PROGRESS.md

#### Step 25: Documentation & Deployment Prep

1. Write comprehensive README.md
2. Add API documentation
3. Create deployment scripts
4. Finalize Docker configurations
5. Add health check endpoints
6. Update PROGRESS.md

### Phase 4: Containerization & Final Polish (Steps 26-30)

#### Step 26: Docker Configuration

1. Create optimized Dockerfile for backend
2. Create optimized Dockerfile for frontend
3. Configure multi-stage builds for production
4. Set up docker-compose for development and production
5. Test container builds and deployments
6. Update PROGRESS.md

#### Step 27: Task Automation

1. Create comprehensive Taskfile.yml with all common tasks
2. Add tasks for testing, building, deploying
3. Include code coverage reporting
4. Add linting and formatting tasks
5. Create setup and teardown tasks
6. Update PROGRESS.md

#### Step 28: Monitoring & Logging

1. Add structured logging throughout the application
2. Implement health check endpoints
3. Add metrics collection
4. Create error tracking and alerting
5. Add request/response logging
6. Update PROGRESS.md

#### Step 29: Production Optimizations

1. Add production environment configurations
2. Implement graceful shutdown
3. Add database migration strategies
4. Configure reverse proxy settings
5. Add SSL/TLS configuration
6. Update PROGRESS.md

#### Step 30: Final Testing & Documentation

1. Run complete test suite and ensure 95%+ coverage
2. Perform security audit
3. Complete final documentation
4. Create deployment guide
5. Mark project as production-ready
6. Final PROGRESS.md update

## Taskfile.yml Structure

Create a comprehensive Taskfile.yml:

```yaml
version: '3'

vars:
  BACKEND_DIR: ./backend
  FRONTEND_DIR: ./frontend
  DOCKER_COMPOSE_FILE: docker-compose.yml

tasks:
  default:
    desc: Show available tasks
    cmds:
      - task --list

  setup:
    desc: Set up the development environment
    cmds:
      - cp .env.example .env
      - cp {{.BACKEND_DIR}}/.env.example {{.BACKEND_DIR}}/.env
      - cp {{.FRONTEND_DIR}}/.env.example {{.FRONTEND_DIR}}/.env
      - docker-compose up -d postgres redis
      - sleep 5
      - task: migrate

  # Backend tasks
  backend:build:
    desc: Build backend application
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - go build -o bin/server cmd/server/main.go

  backend:test:
    desc: Run backend tests
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - go test -v -race -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out -o coverage.html

  backend:test:coverage:
    desc: Check backend test coverage
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - go test -coverprofile=coverage.out ./...
      - go tool cover -func=coverage.out | grep total

  backend:lint:
    desc: Lint backend code
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - golangci-lint run

  backend:run:
    desc: Run backend server
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - go run cmd/server/main.go

  # Frontend tasks
  frontend:install:
    desc: Install frontend dependencies
    dir: '{{.FRONTEND_DIR}}'
    cmds:
      - npm install

  frontend:build:
    desc: Build frontend application
    dir: '{{.FRONTEND_DIR}}'
    cmds:
      - npm run build

  frontend:test:
    desc: Run frontend tests
    dir: '{{.FRONTEND_DIR}}'
    cmds:
      - npm run test:coverage

  frontend:lint:
    desc: Lint frontend code
    dir: '{{.FRONTEND_DIR}}'
    cmds:
      - npm run lint

  frontend:dev:
    desc: Run frontend in development mode
    dir: '{{.FRONTEND_DIR}}'
    cmds:
      - npm run dev

  # Database tasks
  migrate:
    desc: Run database migrations
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - go run cmd/migrate/main.go

  migrate:create:
    desc: Create new migration
    dir: '{{.BACKEND_DIR}}'
    cmds:
      - migrate create -ext sql -dir migrations {{.CLI_ARGS}}

  # Docker tasks
  docker:build:
    desc: Build all Docker images
    cmds:
      - docker-compose build

  docker:up:
    desc: Start all services
    cmds:
      - docker-compose up -d

  docker:down:
    desc: Stop all services
    cmds:
      - docker-compose down

  docker:logs:
    desc: View logs
    cmds:
      - docker-compose logs -f

  # Testing tasks
  test:all:
    desc: Run all tests
    cmds:
      - task: backend:test
      - task: frontend:test

  test:coverage:
    desc: Check test coverage for all components
    cmds:
      - task: backend:test:coverage
      - task: frontend:test

  # Development tasks
  dev:
    desc: Start development environment
    cmds:
      - docker-compose up -d postgres redis
      - task: frontend:dev &
      - task: backend:run

  clean:
    desc: Clean up build artifacts
    cmds:
      - rm -rf {{.BACKEND_DIR}}/bin
      - rm -rf {{.FRONTEND_DIR}}/dist
      - rm -rf {{.FRONTEND_DIR}}/node_modules
      - docker-compose down -v
```

## Testing Requirements

### Backend Testing (Target: 95%+ Coverage)

1. **Unit Tests**: Test all services, repositories, and handlers in isolation
2. **Integration Tests**: Test API endpoints with real database
3. **Test Structure**: Use testify/suite for organized test suites
4. **Mocking**: Use testify/mock for dependency mocking
5. **Test Data**: Use testdata/ directory for test fixtures
6. **Coverage**: Run `go test -coverprofile=coverage.out ./...` to measure coverage

### Frontend Testing (Target: 95%+ Coverage)

1. **Component Tests**: Test all React components with React Testing Library
2. **Hook Tests**: Test custom hooks in isolation
3. **Service Tests**: Test API services with mocked responses
4. **Integration Tests**: Test complete user workflows
5. **Coverage**: Use `npm run test:coverage` to measure coverage

## Progress Tracking

Update PROGRESS.md after each step with the following format:

```markdown
# URL Shortener Development Progress

## Phase 1: Project Foundation

- [x] Step 1: Initialize Repository Structure (Completed: 2024-01-15)
- [x] Step 2: Environment Configuration (Completed: 2024-01-15)
- [ ] Step 3: Database Setup (In Progress)
- [ ] Step 4: Redis Cache Setup
- [ ] Step 5: Core Domain Models
      ...

## Current Status

**Current Step**: Step 3 - Database Setup  
**Completion**: 6.67% (2/30 steps)  
**Next Steps**: Complete PostgreSQL setup and migrations

## Notes

- Database connection established successfully
- All migrations created and tested
- Need to add indexes for performance optimization

## Issues/Blockers

None currently

## Test Coverage Status

- Backend: 0% (tests not implemented yet)
- Frontend: 0% (not started yet)
```

## Quality Checklist

Before marking any step as complete, ensure:

### Code Quality

- [ ] All functions have appropriate error handling
- [ ] All public functions have godoc comments
- [ ] Code follows established patterns and conventions
- [ ] No TODO comments in production code
- [ ] All magic numbers are defined as constants

### Testing

- [ ] Unit tests cover all new functionality
- [ ] Integration tests cover API endpoints
- [ ] Test coverage meets 95% threshold
- [ ] All tests pass consistently
- [ ] Tests include edge cases and error scenarios

### Security

- [ ] All user inputs are validated and sanitized
- [ ] SQL queries use parameterized statements
- [ ] Authentication is properly implemented
- [ ] Rate limiting is in place
- [ ] Sensitive data is not logged

### Performance

- [ ] Database queries are optimized
- [ ] Appropriate indexes are in place
- [ ] Caching is implemented where beneficial
- [ ] No N+1 query problems
- [ ] Frontend bundle size is optimized

### Documentation

- [ ] README is updated with new features
- [ ] API documentation is current
- [ ] Code comments explain complex logic
- [ ] PROGRESS.md is updated
- [ ] Deployment instructions are accurate

## Additional Notes

1. **Error Handling**: Implement comprehensive error handling with proper HTTP status codes and user-friendly error messages.

2. **Logging**: Use structured logging with appropriate log levels. Log all errors and important business events.

3. **Configuration**: Use environment variables for all configuration. Never hardcode secrets or environment-specific values.

4. **Database Migrations**: Create reversible migrations and test them thoroughly. Always backup before running migrations in production.

5. **API Design**: Follow RESTful principles and use consistent response formats. Include proper HTTP status codes and error responses.

6. **Frontend State Management**: Keep state management simple with Context API. Only add complexity if needed.

7. **Performance**: Optimize for performance from the start. Use appropriate indexes, implement caching, and optimize bundle sizes.

8. **Security**: Implement security best practices including input validation, authentication, authorization, and protection against common vulnerabilities.

Start with Step 1 and work through each step methodically. Update PROGRESS.md after each step and ensure all quality requirements are met before proceeding to the next step.
