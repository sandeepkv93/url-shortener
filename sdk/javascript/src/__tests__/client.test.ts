import { URLShortenerClient, ClientFactory } from '../client';

describe('URLShortenerClient', () => {
  const mockConfig = {
    baseURL: 'http://localhost:8080',
    accessToken: 'test-access-token',
    refreshToken: 'test-refresh-token',
  };

  let client: URLShortenerClient;

  beforeEach(() => {
    client = new URLShortenerClient(mockConfig);
  });

  describe('constructor', () => {
    it('should create client with config', () => {
      expect(client).toBeInstanceOf(URLShortenerClient);
      expect(client.getConfig().baseURL).toBe(mockConfig.baseURL);
    });

    it('should initialize all services', () => {
      expect(client.auth).toBeDefined();
      expect(client.urls).toBeDefined();
      expect(client.analytics).toBeDefined();
      expect(client.qr).toBeDefined();
      expect(client.webhooks).toBeDefined();
    });
  });

  describe('authentication', () => {
    it('should check authentication status', () => {
      expect(client.isAuthenticated()).toBe(true);
    });

    it('should set tokens', () => {
      const newAccessToken = 'new-access-token';
      const newRefreshToken = 'new-refresh-token';
      
      client.setTokens(newAccessToken, newRefreshToken);
      
      expect(client.getConfig().accessToken).toBe(newAccessToken);
      expect(client.getConfig().refreshToken).toBe(newRefreshToken);
    });

    it('should clear tokens', () => {
      client.clearTokens();
      
      expect(client.getConfig().accessToken).toBeUndefined();
      expect(client.getConfig().refreshToken).toBeUndefined();
      expect(client.isAuthenticated()).toBe(false);
    });
  });

  describe('configuration', () => {
    it('should update config', () => {
      const newTimeout = 5000;
      client.updateConfig({ timeout: newTimeout });
      
      expect(client.getConfig().timeout).toBe(newTimeout);
    });

    it('should clone client', () => {
      const clonedClient = client.clone();
      
      expect(clonedClient).toBeInstanceOf(URLShortenerClient);
      expect(clonedClient.getConfig()).toEqual(client.getConfig());
      expect(clonedClient).not.toBe(client);
    });
  });

  describe('disposal', () => {
    it('should dispose client and clear tokens', () => {
      client.dispose();
      
      expect(client.getConfig().accessToken).toBeUndefined();
      expect(client.getConfig().refreshToken).toBeUndefined();
    });
  });
});

describe('ClientFactory', () => {
  describe('development', () => {
    it('should create development client', () => {
      const client = ClientFactory.development({});
      
      expect(client.getConfig().baseURL).toBe('http://localhost:8080');
      expect(client.getConfig().timeout).toBe(10000);
    });
  });

  describe('production', () => {
    it('should create production client', () => {
      const client = ClientFactory.production({});
      
      expect(client.getConfig().baseURL).toBe('https://api.urlshortener.com');
      expect(client.getConfig().timeout).toBe(30000);
      expect(client.getConfig().retries).toBe(3);
    });
  });

  describe('testing', () => {
    it('should create testing client', () => {
      const client = ClientFactory.testing();
      
      expect(client.getConfig().baseURL).toBe('http://localhost:8080');
      expect(client.getConfig().timeout).toBe(5000);
      expect(client.getConfig().retries).toBe(1);
    });
  });

  describe('withTokenPersistence', () => {
    const mockStorage = {
      data: {} as Record<string, string>,
      getItem: jest.fn((key: string) => mockStorage.data[key] || null),
      setItem: jest.fn((key: string, value: string) => {
        mockStorage.data[key] = value;
      }),
      removeItem: jest.fn((key: string) => {
        delete mockStorage.data[key];
      }),
    };

    beforeEach(() => {
      jest.clearAllMocks();
      mockStorage.data = {};
    });

    it('should load existing tokens from storage', () => {
      mockStorage.data['urlshortener_access_token'] = 'stored-access-token';
      mockStorage.data['urlshortener_refresh_token'] = 'stored-refresh-token';

      const client = ClientFactory.withTokenPersistence(
        { baseURL: 'http://localhost:8080' },
        mockStorage
      );

      expect(client.getConfig().accessToken).toBe('stored-access-token');
      expect(client.getConfig().refreshToken).toBe('stored-refresh-token');
    });

    it('should handle token refresh callback', () => {
      const mockOnTokenRefresh = jest.fn();
      const client = ClientFactory.withTokenPersistence(
        { 
          baseURL: 'http://localhost:8080',
          onTokenRefresh: mockOnTokenRefresh,
        },
        mockStorage
      );

      const newTokens = {
        accessToken: 'new-access-token',
        refreshToken: 'new-refresh-token',
      };

      // Simulate token refresh
      client.getConfig().onTokenRefresh?.(newTokens);

      expect(mockStorage.setItem).toHaveBeenCalledWith('urlshortener_access_token', 'new-access-token');
      expect(mockStorage.setItem).toHaveBeenCalledWith('urlshortener_refresh_token', 'new-refresh-token');
      expect(mockOnTokenRefresh).toHaveBeenCalledWith(newTokens);
    });

    it('should handle authentication errors', () => {
      const mockOnError = jest.fn();
      const client = ClientFactory.withTokenPersistence(
        { 
          baseURL: 'http://localhost:8080',
          onError: mockOnError,
        },
        mockStorage
      );

      const authError = { statusCode: 401, error: 'Unauthorized' };

      // Simulate auth error
      client.getConfig().onError?.(authError);

      expect(mockStorage.removeItem).toHaveBeenCalledWith('urlshortener_access_token');
      expect(mockStorage.removeItem).toHaveBeenCalledWith('urlshortener_refresh_token');
      expect(mockOnError).toHaveBeenCalledWith(authError);
    });
  });
});