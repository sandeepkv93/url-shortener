import { HTTPClient } from './http-client';
import { ClientConfig } from './config';
import { AuthService } from './services/auth';
import { URLService } from './services/urls';
import { AnalyticsService } from './services/analytics';
import { QRService } from './services/qr';
import { WebhookService } from './services/webhooks';

/**
 * Main URL Shortener API Client
 * 
 * @example
 * ```typescript
 * import { URLShortenerClient } from '@urlshortener/sdk';
 * 
 * const client = new URLShortenerClient({
 *   baseURL: 'https://api.urlshortener.com',
 *   onTokenRefresh: (tokens) => {
 *     // Save tokens to localStorage or secure storage
 *     localStorage.setItem('accessToken', tokens.accessToken);
 *     localStorage.setItem('refreshToken', tokens.refreshToken);
 *   }
 * });
 * 
 * // Authenticate
 * await client.auth.login({ email: 'user@example.com', password: 'password' });
 * 
 * // Create a short URL
 * const shortURL = await client.urls.create({
 *   originalURL: 'https://example.com/very-long-url',
 *   title: 'My Example URL'
 * });
 * 
 * // Get analytics
 * const analytics = await client.analytics.getURLAnalytics(shortURL.id);
 * ```
 */
export class URLShortenerClient {
  private httpClient: HTTPClient;
  
  /** Authentication service */
  public readonly auth: AuthService;
  
  /** URL management service */
  public readonly urls: URLService;
  
  /** Analytics service */
  public readonly analytics: AnalyticsService;
  
  /** QR code service */
  public readonly qr: QRService;
  
  /** Webhook service */
  public readonly webhooks: WebhookService;

  constructor(config: ClientConfig) {
    this.httpClient = new HTTPClient(config);
    
    // Initialize services
    this.auth = new AuthService(this.httpClient);
    this.urls = new URLService(this.httpClient);
    this.analytics = new AnalyticsService(this.httpClient);
    this.qr = new QRService(this.httpClient);
    this.webhooks = new WebhookService(this.httpClient);
  }

  /**
   * Update client configuration
   */
  updateConfig(newConfig: Partial<ClientConfig>): void {
    this.httpClient.updateConfig(newConfig);
  }

  /**
   * Get current configuration
   */
  getConfig(): ClientConfig {
    return this.httpClient.getConfig();
  }

  /**
   * Set authentication tokens manually
   */
  setTokens(accessToken: string, refreshToken: string): void {
    this.httpClient.setAccessToken(accessToken);
    this.httpClient.setRefreshToken(refreshToken);
  }

  /**
   * Clear authentication tokens
   */
  clearTokens(): void {
    this.httpClient.clearTokens();
  }

  /**
   * Check if client is authenticated
   */
  isAuthenticated(): boolean {
    return this.auth.isAuthenticated();
  }

  /**
   * Get API health status
   */
  async getHealth(): Promise<{
    status: string;
    version: string;
    uptime: number;
    timestamp: string;
  }> {
    return this.httpClient.get('/health');
  }

  /**
   * Get API version information
   */
  async getVersion(): Promise<{
    version: string;
    buildDate: string;
    gitCommit: string;
    apiVersion: string;
  }> {
    return this.httpClient.get('/version');
  }

  /**
   * Test API connectivity
   */
  async ping(): Promise<{ message: string; timestamp: string }> {
    return this.httpClient.get('/ping');
  }

  /**
   * Create a new client instance with the same configuration
   */
  clone(): URLShortenerClient {
    const config = this.getConfig();
    return new URLShortenerClient(config);
  }

  /**
   * Dispose of the client and cleanup resources
   */
  dispose(): void {
    this.clearTokens();
    // Additional cleanup if needed
  }
}

/**
 * Create a new URL Shortener client instance
 * 
 * @param config Client configuration
 * @returns New client instance
 */
export function createClient(config: ClientConfig): URLShortenerClient {
  return new URLShortenerClient(config);
}

/**
 * Default client factory with common configurations
 */
export class ClientFactory {
  /**
   * Create a client for development environment
   */
  static development(config: Omit<ClientConfig, 'baseURL'> & { baseURL?: string }): URLShortenerClient {
    return new URLShortenerClient({
      baseURL: 'http://localhost:8080',
      timeout: 10000,
      ...config,
    });
  }

  /**
   * Create a client for production environment
   */
  static production(config: Omit<ClientConfig, 'baseURL'> & { baseURL?: string }): URLShortenerClient {
    return new URLShortenerClient({
      baseURL: 'https://api.urlshortener.com',
      timeout: 30000,
      retries: 3,
      ...config,
    });
  }

  /**
   * Create a client with token persistence
   */
  static withTokenPersistence(
    config: ClientConfig,
    storage: {
      getItem: (key: string) => string | null;
      setItem: (key: string, value: string) => void;
      removeItem: (key: string) => void;
    } = localStorage
  ): URLShortenerClient {
    // Load existing tokens
    const accessToken = storage.getItem('urlshortener_access_token');
    const refreshToken = storage.getItem('urlshortener_refresh_token');

    const client = new URLShortenerClient({
      ...config,
      accessToken: accessToken || undefined,
      refreshToken: refreshToken || undefined,
      onTokenRefresh: (tokens) => {
        storage.setItem('urlshortener_access_token', tokens.accessToken);
        storage.setItem('urlshortener_refresh_token', tokens.refreshToken);
        config.onTokenRefresh?.(tokens);
      },
      onError: (error) => {
        // Clear tokens on authentication errors
        if (error.statusCode === 401) {
          storage.removeItem('urlshortener_access_token');
          storage.removeItem('urlshortener_refresh_token');
        }
        config.onError?.(error);
      },
    });

    return client;
  }

  /**
   * Create a client for testing with mocked responses
   */
  static testing(config: Partial<ClientConfig> = {}): URLShortenerClient {
    return new URLShortenerClient({
      baseURL: 'http://localhost:8080',
      timeout: 5000,
      retries: 1,
      ...config,
    });
  }
}

// Export types and client
export * from './types';
export * from './config';
export { 
  AuthService, 
  URLService, 
  AnalyticsService, 
  QRService, 
  WebhookService 
} from './services';

// Default export
export default URLShortenerClient;