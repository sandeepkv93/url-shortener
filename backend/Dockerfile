# Backend Dockerfile for URL Shortener
# Multi-stage build for optimal production image

# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bin/server \
    cmd/server/main.go

# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bin/migrate \
    cmd/migrate/main.go

# Stage 2: Development stage
FROM alpine:3.18 AS development

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl bash

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy source code for development
COPY --chown=appuser:appgroup . .

# Copy built binaries
COPY --from=builder --chown=appuser:appgroup /app/bin/ ./bin/

# Create necessary directories
RUN mkdir -p logs tmp && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command for development
CMD ["./bin/server"]

# Stage 3: Production stage
FROM scratch AS production

# Copy timezone data and CA certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy user info from development stage
COPY --from=development /etc/passwd /etc/passwd
COPY --from=development /etc/group /etc/group

# Set working directory
WORKDIR /app

# Copy built binaries
COPY --from=builder --chown=1001:1001 /app/bin/server ./server
COPY --from=builder --chown=1001:1001 /app/bin/migrate ./migrate

# Copy static files if any
COPY --from=builder --chown=1001:1001 /app/docs ./docs

# Switch to non-root user
USER 1001

# Expose port
EXPOSE 8080

# Production command
CMD ["./server"]

# Final stage selector
FROM ${BUILD_TARGET:-production} AS final