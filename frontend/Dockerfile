# Frontend Dockerfile for URL Shortener
# Multi-stage build for optimal production image

# Stage 1: Build stage
FROM node:18-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git python3 make g++

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production --silent

# Copy source code
COPY . .

# Build application
RUN npm run build

# Stage 2: Development stage
FROM node:18-alpine AS development

# Install curl for health checks
RUN apk add --no-cache curl

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies (including dev dependencies)
RUN npm install --silent

# Copy source code
COPY --chown=appuser:appgroup . .

# Create node_modules cache directory
RUN mkdir -p node_modules/.cache && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000 5173

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5173 || exit 1

# Default command for development
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Stage 3: Production build
FROM node:18-alpine AS production-builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production --silent

# Copy source code
COPY . .

# Build for production
ENV NODE_ENV=production
RUN npm run build

# Stage 4: Production stage with Nginx
FROM nginx:1.24-alpine AS production

# Install security updates
RUN apk update && apk upgrade

# Remove default Nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Create nginx user and group if they don't exist
RUN addgroup -g 101 -S nginx 2>/dev/null || true && \
    adduser -u 101 -D -S -G nginx nginx 2>/dev/null || true

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY nginx-default.conf /etc/nginx/conf.d/default.conf

# Copy built application from builder stage
COPY --from=production-builder --chown=nginx:nginx /app/dist /usr/share/nginx/html

# Create necessary directories
RUN mkdir -p /var/cache/nginx/client_temp \
             /var/cache/nginx/proxy_temp \
             /var/cache/nginx/fastcgi_temp \
             /var/cache/nginx/uwsgi_temp \
             /var/cache/nginx/scgi_temp && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/log/nginx

# Switch to non-root user
USER nginx

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]

# Final stage selector
FROM ${BUILD_TARGET:-production} AS final