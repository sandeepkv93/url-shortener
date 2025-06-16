# URL Shortener Service - Technical Specification

## Project Overview

A full-stack URL shortener service similar to bit.ly or tinyurl.com, featuring analytics, custom aliases, and comprehensive link management capabilities.

## Core Features

### 1. URL Shortening Engine

- **Basic URL Shortening**: Convert long URLs into short, memorable links
- **Custom Aliases**: Allow users to create custom short codes (e.g., `mysite.com/meeting`)
- **Auto-generated Codes**: Generate random 6-8 character alphanumeric codes
- **URL Validation**: Validate input URLs for proper format and accessibility
- **Duplicate Prevention**: Check for existing shortened URLs to avoid duplicates
- **Expiration Dates**: Optional expiration for temporary links
- **Bulk Shortening**: Process multiple URLs simultaneously

### 2. Analytics & Tracking

- **Click Tracking**: Record every click with timestamp and metadata
- **Geographic Analytics**: Track clicks by country, region, and city
- **Device & Browser Analytics**: Identify user agents, devices, and browsers
- **Referrer Tracking**: Track traffic sources (direct, social media, etc.)
- **Time-based Analytics**: Hourly, daily, weekly, and monthly click patterns
- **Unique vs Total Clicks**: Distinguish between unique visitors and total clicks
- **Real-time Statistics**: Live click counter and recent activity feed

### 3. User Management

- **User Registration/Login**: Email-based authentication with JWT tokens
- **Password Reset**: Secure password recovery via email
- **User Profiles**: Manage personal information and preferences
- **Link Ownership**: Associate shortened URLs with user accounts
- **Guest Mode**: Allow anonymous URL shortening with limited features
- **Account Dashboard**: Overview of user's links and statistics

### 4. Link Management

- **Link History**: View all previously created short links
- **Link Editing**: Modify destination URLs and settings
- **Link Deletion**: Remove unwanted short links
- **Link Status**: Enable/disable links without deletion
- **Search & Filter**: Find links by URL, alias, or creation date
- **Batch Operations**: Bulk enable/disable/delete operations
- **Link Categories**: Organize links with custom tags/categories

### 5. QR Code Generation

- **Automatic QR Codes**: Generate QR codes for all shortened URLs
- **Customizable QR Codes**: Different sizes, colors, and formats
- **SVG/PNG Export**: Download QR codes in various formats
- **Branded QR Codes**: Include logo or custom styling
- **Bulk QR Generation**: Generate QR codes for multiple links

### 6. Security & Privacy

- **Rate Limiting**: Prevent abuse with request throttling
- **Spam Detection**: Identify and block malicious URLs
- **Link Previews**: Safe preview of destination URLs
- **Password Protection**: Optional password protection for sensitive links
- **Access Control**: Restrict access by IP, geographic location, or time
- **Audit Logs**: Track all user actions and system events

### 7. Advanced Features

- **API Access**: RESTful API for third-party integrations
- **Webhook Support**: Real-time notifications for click events
- **Custom Domains**: Allow users to use their own domains
- **Link Retargeting**: Redirect to different URLs based on conditions
- **A/B Testing**: Split traffic between multiple destination URLs
- **UTM Parameter Support**: Automatic UTM parameter handling

## System Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │   Golang API    │    │   PostgreSQL    │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│   (Database)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │
                               ▼
                       ┌─────────────────┐
                       │   Redis Cache   │
                       │   (Sessions)    │
                       └─────────────────┘
```

### Backend Architecture (Golang)

#### Technology Stack

- **Framework**: Gin or Fiber for HTTP routing
- **Database**: PostgreSQL for persistent data
- **Cache**: Redis for sessions and rate limiting
- **Authentication**: JWT tokens
- **Validation**: go-playground/validator
- **Database ORM**: GORM or SQLx
- **Testing**: Go testing package + Testify

#### Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   ├── core/
│   │   ├── domain/
│   │   ├── ports/
│   │   └── services/
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── cache/
│   │   └── external/
│   └── config/
├── migrations/
├── tests/
└── docker/
```

#### Core Components

**1. URL Shortening Service**

```go
type URLService interface {
    ShortenURL(url, customAlias string, userID int) (*ShortURL, error)
    GetOriginalURL(shortCode string) (*ShortURL, error)
    GenerateShortCode() string
    ValidateURL(url string) error
}
```

**2. Analytics Service**

```go
type AnalyticsService interface {
    RecordClick(shortCode, ip, userAgent, referer string) error
    GetClickStats(shortCode string, period string) (*ClickStats, error)
    GetGeoStats(shortCode string) (*GeoStats, error)
}
```

**3. User Service**

```go
type UserService interface {
    CreateUser(email, password string) (*User, error)
    AuthenticateUser(email, password string) (*User, error)
    GetUserByID(id int) (*User, error)
}
```

### Database Schema

#### PostgreSQL Tables

**Users Table**

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Short URLs Table**

```sql
CREATE TABLE short_urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    custom_alias BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Clicks Table**

```sql
CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    short_url_id INTEGER REFERENCES short_urls(id),
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    country VARCHAR(2),
    city VARCHAR(100),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Redis Cache Structure

```
# Rate limiting
rate_limit:{ip}:{endpoint} -> {count, ttl}

# Session storage
session:{token} -> {user_id, expires_at}

# URL cache
url_cache:{short_code} -> {original_url, user_id, is_active}

# Click counters
click_count:{short_code} -> {total_clicks}
unique_clicks:{short_code} -> Set{ip_addresses}
```

### API Endpoints

#### Authentication Endpoints

```
POST   /api/auth/register     - User registration
POST   /api/auth/login        - User login
POST   /api/auth/logout       - User logout
POST   /api/auth/refresh      - Refresh JWT token
POST   /api/auth/forgot       - Password reset request
POST   /api/auth/reset        - Password reset confirmation
```

#### URL Management Endpoints

```
POST   /api/urls              - Create short URL
GET    /api/urls              - List user's URLs
GET    /api/urls/:id          - Get specific URL details
PUT    /api/urls/:id          - Update URL
DELETE /api/urls/:id          - Delete URL
POST   /api/urls/bulk         - Bulk URL operations
```

#### Analytics Endpoints

```
GET    /api/analytics/:code          - Get click statistics
GET    /api/analytics/:code/geo      - Geographic statistics
GET    /api/analytics/:code/devices  - Device/browser stats
GET    /api/analytics/:code/referrers - Referrer statistics
GET    /api/analytics/:code/timeline - Time-based analytics
```

#### Redirect Endpoint

```
GET    /:shortCode            - Redirect to original URL
```

#### QR Code Endpoints

```
GET    /api/qr/:code          - Generate QR code
GET    /api/qr/:code.svg      - SVG format QR code
GET    /api/qr/:code.png      - PNG format QR code
```

### Frontend Architecture (React)

#### Technology Stack

- **Framework**: React 18 with TypeScript
- **Routing**: React Router v6
- **State Management**: Context API + useReducer or Zustand
- **HTTP Client**: Axios
- **UI Framework**: Tailwind CSS + Headless UI
- **Charts**: Chart.js or Recharts
- **Forms**: React Hook Form
- **QR Codes**: qrcode.js
- **Icons**: Lucide React

#### Component Structure

```
src/
├── components/
│   ├── common/
│   │   ├── Header.tsx
│   │   ├── Footer.tsx
│   │   └── Layout.tsx
│   ├── auth/
│   │   ├── LoginForm.tsx
│   │   ├── RegisterForm.tsx
│   │   └── PasswordReset.tsx
│   ├── url/
│   │   ├── URLShortener.tsx
│   │   ├── URLList.tsx
│   │   └── URLCard.tsx
│   ├── analytics/
│   │   ├── Dashboard.tsx
│   │   ├── ClickChart.tsx
│   │   └── GeographicMap.tsx
│   └── qr/
│       ├── QRGenerator.tsx
│       └── QRPreview.tsx
├── pages/
│   ├── Home.tsx
│   ├── Dashboard.tsx
│   ├── Analytics.tsx
│   └── Profile.tsx
├── hooks/
├── services/
├── utils/
└── types/
```

#### Key React Components

**URL Shortener Component**

```typescript
interface URLShortenerProps {
  onURLCreated: (url: ShortURL) => void
}

const URLShortener: React.FC<URLShortenerProps> = ({ onURLCreated }) => {
  // Form handling, validation, API calls
  // Custom alias toggle
  // Expiration date picker
  // Real-time URL validation
}
```

**Analytics Dashboard**

```typescript
interface AnalyticsDashboardProps {
  shortCode: string
}

const AnalyticsDashboard: React.FC<AnalyticsDashboardProps> = ({
  shortCode,
}) => {
  // Click statistics display
  // Interactive charts
  // Geographic visualization
  // Device/browser breakdown
  // Real-time updates
}
```

## Implementation Phases

### Phase 1: Core Functionality (Week 1-2)

- Basic URL shortening and redirection
- User authentication
- Simple frontend with URL creation form
- Basic click tracking

### Phase 2: Analytics & Management (Week 3-4)

- Comprehensive analytics dashboard
- Link management interface
- Click statistics and charts
- User dashboard

### Phase 3: Advanced Features (Week 5-6)

- QR code generation
- Custom aliases
- Rate limiting and security
- Bulk operations

### Phase 4: Polish & Optimization (Week 7-8)

- Performance optimization
- Advanced analytics
- Error handling and validation
- Testing and documentation

## Technical Considerations

### Performance Optimizations

- **Database Indexing**: Optimize queries with proper indexes
- **Caching Strategy**: Cache frequently accessed URLs in Redis
- **CDN Integration**: Serve static assets via CDN
- **Connection Pooling**: Efficient database connection management
- **Horizontal Scaling**: Design for multiple server instances

### Security Measures

- **Input Validation**: Sanitize and validate all user inputs
- **Rate Limiting**: Implement per-IP and per-user rate limits
- **HTTPS Only**: Force secure connections
- **CORS Configuration**: Proper cross-origin resource sharing
- **SQL Injection Prevention**: Use parameterized queries

### Monitoring & Logging

- **Application Metrics**: Track response times, error rates
- **Click Analytics**: Real-time click tracking and aggregation
- **Error Logging**: Comprehensive error tracking and alerting
- **Performance Monitoring**: Database query performance
- **User Activity Logs**: Audit trail for user actions

## Deployment Strategy

### Containerization

```dockerfile
# Golang backend Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
CMD ["./main"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    environment:
      - DATABASE_URL=postgres://...
      - REDIS_URL=redis://...
    depends_on:
      - postgres
      - redis

  frontend:
    build: ./frontend
    ports:
      - '3000:3000'

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: urlshortener
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
```

### Production Deployment

- **Cloud Platform**: AWS, Google Cloud, or DigitalOcean
- **Container Orchestration**: Docker Swarm or Kubernetes
- **Load Balancer**: Nginx or cloud load balancer
- **Database**: Managed PostgreSQL service
- **Cache**: Managed Redis service
- **Monitoring**: Prometheus + Grafana

This comprehensive specification provides a solid foundation for building a production-ready URL shortener service using Claude Code assistance.
