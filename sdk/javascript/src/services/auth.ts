import { HTTPClient } from '../http-client';
import { API_ENDPOINTS } from '../config';
import { AuthCredentials, AuthTokens, User, RequestConfig } from '../types';

export class AuthService {
  constructor(private client: HTTPClient) {}

  /**
   * Register a new user account
   */
  async register(credentials: AuthCredentials, config?: RequestConfig): Promise<AuthTokens & { user: User }> {
    const response = await this.client.post<AuthTokens & { user: User }>(
      API_ENDPOINTS.AUTH_REGISTER,
      credentials,
      config
    );
    
    // Update client tokens
    this.client.setAccessToken(response.accessToken);
    this.client.setRefreshToken(response.refreshToken);
    
    return response;
  }

  /**
   * Login with email and password
   */
  async login(credentials: AuthCredentials, config?: RequestConfig): Promise<AuthTokens & { user: User }> {
    const response = await this.client.post<AuthTokens & { user: User }>(
      API_ENDPOINTS.AUTH_LOGIN,
      credentials,
      config
    );
    
    // Update client tokens
    this.client.setAccessToken(response.accessToken);
    this.client.setRefreshToken(response.refreshToken);
    
    return response;
  }

  /**
   * Refresh access token using refresh token
   */
  async refreshToken(refreshToken?: string, config?: RequestConfig): Promise<AuthTokens> {
    const token = refreshToken || this.client.getConfig().refreshToken;
    if (!token) {
      throw new Error('No refresh token available');
    }

    const response = await this.client.post<AuthTokens>(
      API_ENDPOINTS.AUTH_REFRESH,
      { refreshToken: token },
      config
    );
    
    // Update client tokens
    this.client.setAccessToken(response.accessToken);
    this.client.setRefreshToken(response.refreshToken);
    
    return response;
  }

  /**
   * Logout and invalidate tokens
   */
  async logout(config?: RequestConfig): Promise<void> {
    try {
      await this.client.post(API_ENDPOINTS.AUTH_LOGOUT, {}, config);
    } finally {
      // Clear tokens regardless of API response
      this.client.clearTokens();
    }
  }

  /**
   * Get current user profile
   */
  async getProfile(config?: RequestConfig): Promise<User> {
    return this.client.get<User>(API_ENDPOINTS.AUTH_PROFILE, config);
  }

  /**
   * Update user profile
   */
  async updateProfile(updates: Partial<Pick<User, 'email'>>, config?: RequestConfig): Promise<User> {
    return this.client.put<User>(API_ENDPOINTS.AUTH_UPDATE_PROFILE, updates, config);
  }

  /**
   * Change user password
   */
  async changePassword(
    data: { currentPassword: string; newPassword: string },
    config?: RequestConfig
  ): Promise<void> {
    await this.client.post(API_ENDPOINTS.AUTH_CHANGE_PASSWORD, data, config);
  }

  /**
   * Validate current access token
   */
  async validateToken(config?: RequestConfig): Promise<{ valid: boolean; user?: User }> {
    try {
      const user = await this.client.get<User>(API_ENDPOINTS.AUTH_VALIDATE, config);
      return { valid: true, user };
    } catch (error) {
      return { valid: false };
    }
  }

  /**
   * Check if user is currently authenticated
   */
  isAuthenticated(): boolean {
    return !!this.client.getConfig().accessToken;
  }

  /**
   * Get current access token
   */
  getAccessToken(): string | undefined {
    return this.client.getConfig().accessToken;
  }

  /**
   * Get current refresh token
   */
  getRefreshToken(): string | undefined {
    return this.client.getConfig().refreshToken;
  }

  /**
   * Manually set authentication tokens
   */
  setTokens(accessToken: string, refreshToken: string): void {
    this.client.setAccessToken(accessToken);
    this.client.setRefreshToken(refreshToken);
  }

  /**
   * Clear authentication tokens
   */
  clearTokens(): void {
    this.client.clearTokens();
  }
}