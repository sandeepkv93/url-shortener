# URL Shortener Service

A production-grade URL shortener service built with Go and React, featuring comprehensive analytics, QR code generation, and advanced link management capabilities.

## 🚀 Features

### Core Features
- **URL Shortening**: Create short, memorable links from long URLs
- **Custom Aliases**: Support for custom short codes
- **QR Code Generation**: Generate QR codes for shortened URLs in multiple formats
- **Click Tracking**: Detailed analytics on link clicks and user engagement
- **Password Protection**: Secure URLs with password protection
- **Expiration Dates**: Set automatic expiration for URLs
- **Bulk Operations**: Import/export URLs in bulk

### Analytics & Reporting
- **Real-time Analytics**: Live click tracking and statistics
- **Geographic Analytics**: Location-based click analysis
- **Device Analytics**: Track clicks by device type, browser, and OS
- **Time-based Analytics**: Detailed time-series analysis
- **Export Capabilities**: Export analytics data in multiple formats

### User Management
- **User Registration & Authentication**: Secure JWT-based authentication
- **Profile Management**: User profile and preference management
- **Dashboard**: Comprehensive user dashboard with URL management
- **Email Notifications**: Welcome emails, password reset, and alerts

### Advanced Features
- **Geolocation Services**: IP-based location tracking
- **Health Monitoring**: System health checks and metrics
- **Rate Limiting**: API rate limiting and abuse prevention
- **Caching**: Redis-based caching for performance
- **Security**: Comprehensive security features (CORS, XSS protection, etc.)

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── core/            # Business logic and domain models
│   │   ├── domain/      # Domain entities and errors
│   │   ├── ports/       # Interface definitions
│   │   └── services/    # Business logic implementation
│   ├── infrastructure/  # External dependencies
│   │   ├── database/    # Database connections and repositories
│   │   ├── cache/       # Redis cache implementation
│   │   └── external/    # External API integrations
│   └── config/          # Configuration management
└── tests/               # Test suites
```

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Chi Router
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Testing**: Testify framework

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **State Management**: Context API
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios
- **Charts**: Recharts

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Task Runner**: Taskfile
- **Database**: PostgreSQL 15
- **Cache**: Redis 7

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 15
- Redis 7
- Docker & Docker Compose (optional)

## 🚀 Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd url-shortener
```

2. Start the services:
```bash
docker-compose up -d
```

3. The application will be available at:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/docs

### Manual Setup

1. **Set up the database**:
```bash
# Create PostgreSQL database
createdb url_shortener

# Create Redis instance (or use Docker)
docker run -d -p 6379:6379 redis:7-alpine
```

2. **Configure environment**:
```bash
# Copy environment files
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Edit the .env files with your configuration
```

3. **Start the backend**:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

4. **Start the frontend**:
```bash
cd frontend
npm install
npm run dev
```

## 📁 Project Structure

```
url-shortener/
├── README.md
├── PROGRESS.md                 # Development progress tracking
├── Taskfile.yml               # Task automation
├── docker-compose.yml         # Docker services
├── .env.example              # Environment variables template
├── backend/                  # Go backend application
│   ├── cmd/server/           # Application entry point
│   ├── internal/             # Private application code
│   │   ├── api/             # HTTP layer
│   │   ├── core/            # Business logic
│   │   ├── infrastructure/  # External dependencies
│   │   └── config/          # Configuration
│   ├── tests/               # Test suites
│   └── go.mod              # Go dependencies
├── frontend/                # React frontend application
│   ├── src/                # Source code
│   │   ├── components/     # React components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API services
│   │   ├── hooks/          # Custom React hooks
│   │   ├── types/          # TypeScript types
│   │   └── utils/          # Utility functions
│   └── package.json        # Node.js dependencies
└── scripts/                # Deployment and setup scripts
```

## 🔧 Configuration

### Environment Variables

**Backend Configuration** (`.env`):
```env
# Server
SERVER_HOST=localhost
SERVER_PORT=8080
BASE_URL=http://localhost:8080

# Database
DATABASE_URL=postgres://user:password@localhost:5432/url_shortener

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_TTL=3600
JWT_REFRESH_TOKEN_TTL=604800

# Email (optional)
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-password
```

**Frontend Configuration** (`.env`):
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=URL Shortener
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests
```bash
cd frontend
npm test
npm run test:coverage
```

### Integration Tests
```bash
# Run all tests
task test:all

# Run specific test suites
task backend:test
task frontend:test
```

## 📊 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | User login |
| POST | `/api/auth/refresh` | Refresh JWT token |
| GET | `/api/auth/profile` | Get user profile |

### URL Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/urls` | Create short URL |
| GET | `/api/urls` | List user URLs |
| GET | `/api/urls/:id` | Get URL details |
| PUT | `/api/urls/:id` | Update URL |
| DELETE | `/api/urls/:id` | Delete URL |
| GET | `/:shortCode` | Redirect to original URL |

### Analytics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/analytics/dashboard` | Dashboard statistics |
| GET | `/api/analytics/urls/:id` | URL-specific analytics |
| GET | `/api/analytics/export` | Export analytics data |

### System Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/metrics` | System metrics |

## 🔍 Development

### Code Quality

- **Test Coverage**: Maintain 95%+ test coverage
- **Code Style**: Use `gofmt`, `golint`, and `prettier`
- **Error Handling**: Comprehensive error handling with proper logging
- **Security**: Input validation, SQL injection prevention, XSS protection

### Git Workflow

- Use conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`
- Commit frequently with descriptive messages
- Keep commits atomic and focused

### Development Commands

```bash
# Start development environment
task dev

# Run tests
task test:all

# Build applications
task build

# Run linters
task lint

# Format code
task format

# View logs
task logs
```

## 🚀 Deployment

### Production Deployment

1. **Build the application**:
```bash
task build
```

2. **Deploy with Docker**:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

3. **Set up reverse proxy** (Nginx configuration included)

### Environment-specific Configurations

- **Development**: Full debugging, hot reload
- **Staging**: Production-like environment for testing
- **Production**: Optimized builds, security hardening

## 🔒 Security

### Security Features

- **Authentication**: JWT-based authentication with refresh tokens
- **Authorization**: Role-based access control
- **Input Validation**: Comprehensive input sanitization
- **Rate Limiting**: API rate limiting and abuse prevention
- **HTTPS**: TLS encryption in production
- **CORS**: Configured Cross-Origin Resource Sharing
- **XSS Protection**: Cross-site scripting prevention
- **SQL Injection Protection**: Parameterized queries

### Security Best Practices

- Regular security audits
- Dependency vulnerability scanning
- Secure coding practices
- Environment variable protection
- Regular updates and patches

## 📈 Performance

### Performance Features

- **Caching**: Redis-based caching for frequently accessed data
- **Database Optimization**: Proper indexing and query optimization
- **Connection Pooling**: Efficient database connection management
- **CDN Integration**: Static asset optimization
- **Gzip Compression**: Response compression

### Performance Monitoring

- Health check endpoints
- System metrics collection
- Real-time performance monitoring
- Database query optimization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -am 'Add your feature'`
4. Push to the branch: `git push origin feature/your-feature`
5. Submit a pull request

### Development Guidelines

- Follow clean architecture principles
- Write comprehensive tests
- Use meaningful commit messages
- Update documentation
- Ensure code quality standards

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:

- 📧 Email: support@urlshortener.com
- 📖 Documentation: [docs.urlshortener.com](https://docs.urlshortener.com)
- 🐛 Issues: [GitHub Issues](https://github.com/your-org/url-shortener/issues)

## 🗺️ Roadmap

- [ ] Mobile applications (iOS/Android)
- [ ] Advanced analytics with machine learning
- [ ] Team collaboration features
- [ ] API integrations (Slack, Discord, etc.)
- [ ] Custom domains support
- [ ] Advanced security features
- [ ] Multi-language support
- [ ] Advanced caching strategies

## 📊 Current Status

**Development Progress**: 53% (Step 16 of 30 completed)

See [PROGRESS.md](PROGRESS.md) for detailed development progress and milestones.

---

**Built with ❤️ using Go and React**