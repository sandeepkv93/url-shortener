/**
 * @fileoverview URL Shortener SDK for JavaScript/TypeScript
 * @version 1.0.0
 * @author URL Shortener Team
 * @license MIT
 */

// Main client export
export { URLShortenerClient as default, URLShortenerClient } from './client';

// Client factory and utilities
export { createClient, ClientFactory } from './client';

// Configuration
export { ClientConfig, DEFAULT_CONFIG, API_ENDPOINTS } from './config';

// Services
export { AuthService } from './services/auth';
export { URLService } from './services/urls';
export { AnalyticsService } from './services/analytics';
export { QRService } from './services/qr';
export { WebhookService } from './services/webhooks';

// Types
export * from './types';

// HTTP Client (for advanced usage)
export { HTTPClient } from './http-client';

// Service-specific types
export type { 
  TimelineData, 
  GeoStats, 
  DeviceStats, 
  ReferrerStats, 
  ExportFormat, 
  TimePeriod 
} from './services/analytics';

export type { 
  QRCodeFormats, 
  QRCodeSizes, 
  QRCodeGenerationRequest, 
  QRCodeValidationResult 
} from './services/qr';

export type { 
  WebhookEventCategories, 
  WebhookTestResult 
} from './services/webhooks';

// Constants
export const SDK_VERSION = '1.0.0';
export const API_VERSION = 'v1';

/**
 * SDK Information
 */
export const SDK_INFO = {
  name: '@urlshortener/sdk',
  version: SDK_VERSION,
  apiVersion: API_VERSION,
  userAgent: `URLShortener-SDK-JS/${SDK_VERSION}`,
  repository: 'https://github.com/yourusername/url-shortener',
  documentation: 'https://docs.urlshortener.com/sdk/javascript',
} as const;

/**
 * Common error types for better error handling
 */
export class SDKError extends Error {
  constructor(
    message: string,
    public readonly code?: string,
    public readonly statusCode?: number,
    public readonly field?: string
  ) {
    super(message);
    this.name = 'SDKError';
  }
}

export class AuthenticationError extends SDKError {
  constructor(message = 'Authentication failed') {
    super(message, 'AUTHENTICATION_ERROR', 401);
    this.name = 'AuthenticationError';
  }
}

export class ValidationError extends SDKError {
  constructor(message: string, field?: string) {
    super(message, 'VALIDATION_ERROR', 400, field);
    this.name = 'ValidationError';
  }
}

export class NetworkError extends SDKError {
  constructor(message = 'Network request failed') {
    super(message, 'NETWORK_ERROR');
    this.name = 'NetworkError';
  }
}

export class RateLimitError extends SDKError {
  constructor(message = 'Rate limit exceeded') {
    super(message, 'RATE_LIMIT_ERROR', 429);
    this.name = 'RateLimitError';
  }
}

/**
 * Utility functions for common operations
 */
export const utils = {
  /**
   * Validate URL format
   */
  isValidURL(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  },

  /**
   * Generate random string for custom aliases
   */
  generateRandomAlias(length = 8): string {
    const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  },

  /**
   * Format date for API requests
   */
  formatDate(date: Date): string {
    return date.toISOString().split('T')[0];
  },

  /**
   * Parse API date strings
   */
  parseDate(dateString: string): Date {
    return new Date(dateString);
  },

  /**
   * Format numbers for display
   */
  formatNumber(num: number): string {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
  },

  /**
   * Calculate percentage
   */
  calculatePercentage(part: number, total: number): number {
    if (total === 0) return 0;
    return Math.round((part / total) * 100 * 100) / 100; // Round to 2 decimal places
  },

  /**
   * Debounce function for API calls
   */
  debounce<T extends (...args: any[]) => any>(
    func: T,
    wait: number
  ): (...args: Parameters<T>) => void {
    let timeout: NodeJS.Timeout;
    return (...args: Parameters<T>) => {
      clearTimeout(timeout);
      timeout = setTimeout(() => func.apply(null, args), wait);
    };
  },

  /**
   * Throttle function for API calls
   */
  throttle<T extends (...args: any[]) => any>(
    func: T,
    limit: number
  ): (...args: Parameters<T>) => void {
    let inThrottle: boolean;
    return (...args: Parameters<T>) => {
      if (!inThrottle) {
        func.apply(null, args);
        inThrottle = true;
        setTimeout(() => inThrottle = false, limit);
      }
    };
  },

  /**
   * Retry function with exponential backoff
   */
  async retry<T>(
    fn: () => Promise<T>,
    retries = 3,
    delay = 1000
  ): Promise<T> {
    try {
      return await fn();
    } catch (error) {
      if (retries > 0) {
        await new Promise(resolve => setTimeout(resolve, delay));
        return this.retry(fn, retries - 1, delay * 2);
      }
      throw error;
    }
  },
} as const;

/**
 * Pre-configured client instances for common environments
 */
export const clients = {
  /**
   * Development client (localhost:8080)
   */
  development: (config: Partial<ClientConfig> = {}) => 
    ClientFactory.development(config),

  /**
   * Production client
   */
  production: (config: Partial<ClientConfig> = {}) => 
    ClientFactory.production(config),

  /**
   * Testing client
   */
  testing: (config: Partial<ClientConfig> = {}) => 
    ClientFactory.testing(config),
} as const;