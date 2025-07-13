import { AuthService } from '../auth';
import { HTTPClient } from '../../http-client';
import { AuthCredentials, AuthTokens, User } from '../../types';

// Mock the HTTPClient
jest.mock('../../http-client');

describe('AuthService', () => {
  let authService: AuthService;
  let mockHttpClient: jest.Mocked<HTTPClient>;

  const mockUser: User = {
    id: 1,
    email: 'test@example.com',
    createdAt: '2023-01-01T00:00:00Z',
    updatedAt: '2023-01-01T00:00:00Z',
  };

  const mockTokens: AuthTokens = {
    accessToken: 'access-token-123',
    refreshToken: 'refresh-token-123',
    expiresAt: '2023-01-01T01:00:00Z',
  };

  beforeEach(() => {
    mockHttpClient = new HTTPClient({
      baseURL: 'http://localhost:8080',
    }) as jest.Mocked<HTTPClient>;

    authService = new AuthService(mockHttpClient);

    // Reset all mocks
    jest.clearAllMocks();
  });

  describe('register', () => {
    it('should register user and set tokens', async () => {
      const credentials: AuthCredentials = {
        email: 'test@example.com',
        password: 'password123',
      };

      const expectedResponse = { ...mockTokens, user: mockUser };
      mockHttpClient.post.mockResolvedValue(expectedResponse);

      const result = await authService.register(credentials);

      expect(mockHttpClient.post).toHaveBeenCalledWith(
        '/auth/register',
        credentials,
        undefined
      );
      expect(mockHttpClient.setAccessToken).toHaveBeenCalledWith(mockTokens.accessToken);
      expect(mockHttpClient.setRefreshToken).toHaveBeenCalledWith(mockTokens.refreshToken);
      expect(result).toEqual(expectedResponse);
    });

    it('should handle registration errors', async () => {
      const credentials: AuthCredentials = {
        email: 'invalid@example.com',
        password: 'weak',
      };

      const error = new Error('Registration failed');
      mockHttpClient.post.mockRejectedValue(error);

      await expect(authService.register(credentials)).rejects.toThrow('Registration failed');
      expect(mockHttpClient.setAccessToken).not.toHaveBeenCalled();
      expect(mockHttpClient.setRefreshToken).not.toHaveBeenCalled();
    });
  });

  describe('login', () => {
    it('should login user and set tokens', async () => {
      const credentials: AuthCredentials = {
        email: 'test@example.com',
        password: 'password123',
      };

      const expectedResponse = { ...mockTokens, user: mockUser };
      mockHttpClient.post.mockResolvedValue(expectedResponse);

      const result = await authService.login(credentials);

      expect(mockHttpClient.post).toHaveBeenCalledWith(
        '/auth/login',
        credentials,
        undefined
      );
      expect(mockHttpClient.setAccessToken).toHaveBeenCalledWith(mockTokens.accessToken);
      expect(mockHttpClient.setRefreshToken).toHaveBeenCalledWith(mockTokens.refreshToken);
      expect(result).toEqual(expectedResponse);
    });
  });

  describe('refreshToken', () => {
    it('should refresh token using stored refresh token', async () => {
      mockHttpClient.getConfig.mockReturnValue({
        baseURL: 'http://localhost:8080',
        refreshToken: 'stored-refresh-token',
      });
      mockHttpClient.post.mockResolvedValue(mockTokens);

      const result = await authService.refreshToken();

      expect(mockHttpClient.post).toHaveBeenCalledWith(
        '/auth/refresh',
        { refreshToken: 'stored-refresh-token' },
        undefined
      );
      expect(mockHttpClient.setAccessToken).toHaveBeenCalledWith(mockTokens.accessToken);
      expect(mockHttpClient.setRefreshToken).toHaveBeenCalledWith(mockTokens.refreshToken);
      expect(result).toEqual(mockTokens);
    });

    it('should refresh token using provided refresh token', async () => {
      const customRefreshToken = 'custom-refresh-token';
      mockHttpClient.post.mockResolvedValue(mockTokens);

      const result = await authService.refreshToken(customRefreshToken);

      expect(mockHttpClient.post).toHaveBeenCalledWith(
        '/auth/refresh',
        { refreshToken: customRefreshToken },
        undefined
      );
      expect(result).toEqual(mockTokens);
    });

    it('should throw error when no refresh token available', async () => {
      mockHttpClient.getConfig.mockReturnValue({
        baseURL: 'http://localhost:8080',
      });

      await expect(authService.refreshToken()).rejects.toThrow('No refresh token available');
    });
  });

  describe('logout', () => {
    it('should logout and clear tokens', async () => {
      mockHttpClient.post.mockResolvedValue(undefined);

      await authService.logout();

      expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/logout', {}, undefined);
      expect(mockHttpClient.clearTokens).toHaveBeenCalled();
    });

    it('should clear tokens even if API call fails', async () => {
      mockHttpClient.post.mockRejectedValue(new Error('Network error'));

      await authService.logout();

      expect(mockHttpClient.clearTokens).toHaveBeenCalled();
    });
  });

  describe('getProfile', () => {
    it('should get user profile', async () => {
      mockHttpClient.get.mockResolvedValue(mockUser);

      const result = await authService.getProfile();

      expect(mockHttpClient.get).toHaveBeenCalledWith('/auth/profile', undefined);
      expect(result).toEqual(mockUser);
    });
  });

  describe('updateProfile', () => {
    it('should update user profile', async () => {
      const updates = { email: 'newemail@example.com' };
      const updatedUser = { ...mockUser, ...updates };
      mockHttpClient.put.mockResolvedValue(updatedUser);

      const result = await authService.updateProfile(updates);

      expect(mockHttpClient.put).toHaveBeenCalledWith('/auth/profile', updates, undefined);
      expect(result).toEqual(updatedUser);
    });
  });

  describe('changePassword', () => {
    it('should change user password', async () => {
      const passwordData = {
        currentPassword: 'oldpassword',
        newPassword: 'newpassword',
      };
      mockHttpClient.post.mockResolvedValue(undefined);

      await authService.changePassword(passwordData);

      expect(mockHttpClient.post).toHaveBeenCalledWith(
        '/auth/change-password',
        passwordData,
        undefined
      );
    });
  });

  describe('validateToken', () => {
    it('should validate token and return user', async () => {
      mockHttpClient.get.mockResolvedValue(mockUser);

      const result = await authService.validateToken();

      expect(mockHttpClient.get).toHaveBeenCalledWith('/auth/validate', undefined);
      expect(result).toEqual({ valid: true, user: mockUser });
    });

    it('should handle invalid token', async () => {
      mockHttpClient.get.mockRejectedValue(new Error('Invalid token'));

      const result = await authService.validateToken();

      expect(result).toEqual({ valid: false });
    });
  });

  describe('authentication state', () => {
    it('should check if authenticated', () => {
      mockHttpClient.getConfig.mockReturnValue({
        baseURL: 'http://localhost:8080',
        accessToken: 'some-token',
      });

      expect(authService.isAuthenticated()).toBe(true);

      mockHttpClient.getConfig.mockReturnValue({
        baseURL: 'http://localhost:8080',
      });

      expect(authService.isAuthenticated()).toBe(false);
    });

    it('should get access token', () => {
      const accessToken = 'test-access-token';
      mockHttpClient.getConfig.mockReturnValue({
        baseURL: 'http://localhost:8080',
        accessToken,
      });

      expect(authService.getAccessToken()).toBe(accessToken);
    });

    it('should get refresh token', () => {
      const refreshToken = 'test-refresh-token';
      mockHttpClient.getConfig.mockReturnValue({
        baseURL: 'http://localhost:8080',
        refreshToken,
      });

      expect(authService.getRefreshToken()).toBe(refreshToken);
    });

    it('should set tokens manually', () => {
      const accessToken = 'manual-access-token';
      const refreshToken = 'manual-refresh-token';

      authService.setTokens(accessToken, refreshToken);

      expect(mockHttpClient.setAccessToken).toHaveBeenCalledWith(accessToken);
      expect(mockHttpClient.setRefreshToken).toHaveBeenCalledWith(refreshToken);
    });

    it('should clear tokens manually', () => {
      authService.clearTokens();

      expect(mockHttpClient.clearTokens).toHaveBeenCalled();
    });
  });
});