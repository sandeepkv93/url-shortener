import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { ClientConfig, DEFAULT_CONFIG, API_ENDPOINTS } from './config';
import { APIError, RequestConfig } from './types';

export class HTTPClient {
  private client: AxiosInstance;
  private config: ClientConfig;
  private isRefreshing = false;
  private refreshSubscribers: Array<(token: string) => void> = [];

  constructor(config: ClientConfig) {
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.client = this.createAxiosInstance();
    this.setupInterceptors();
  }

  private createAxiosInstance(): AxiosInstance {
    return axios.create({
      baseURL: this.config.baseURL.endsWith('/api/v1') 
        ? this.config.baseURL 
        : `${this.config.baseURL}/api/v1`,
      timeout: this.config.timeout,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  private setupInterceptors(): void {
    // Request interceptor for adding auth headers
    this.client.interceptors.request.use(
      (config) => {
        if (this.config.accessToken) {
          config.headers = config.headers || {};
          config.headers.Authorization = `Bearer ${this.config.accessToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor for handling errors and token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

        if (error.response?.status === 401 && !originalRequest._retry) {
          if (this.config.refreshToken && !this.isRefreshing) {
            return this.handleTokenRefresh(originalRequest);
          }
        }

        return Promise.reject(this.handleError(error));
      }
    );
  }

  private async handleTokenRefresh(originalRequest: AxiosRequestConfig & { _retry?: boolean }): Promise<AxiosResponse> {
    if (this.isRefreshing) {
      return new Promise((resolve) => {
        this.refreshSubscribers.push((token: string) => {
          originalRequest.headers = originalRequest.headers || {};
          originalRequest.headers.Authorization = `Bearer ${token}`;
          resolve(this.client(originalRequest));
        });
      });
    }

    originalRequest._retry = true;
    this.isRefreshing = true;

    try {
      const response = await this.client.post(API_ENDPOINTS.AUTH_REFRESH, {
        refreshToken: this.config.refreshToken,
      });

      const { accessToken, refreshToken } = response.data;
      this.config.accessToken = accessToken;
      this.config.refreshToken = refreshToken;

      // Notify subscribers
      this.refreshSubscribers.forEach((callback) => callback(accessToken));
      this.refreshSubscribers = [];

      // Call the token refresh callback
      if (this.config.onTokenRefresh) {
        this.config.onTokenRefresh({ accessToken, refreshToken });
      }

      // Retry the original request
      originalRequest.headers = originalRequest.headers || {};
      originalRequest.headers.Authorization = `Bearer ${accessToken}`;
      return this.client(originalRequest);
    } catch (refreshError) {
      // Refresh failed, clear tokens
      this.config.accessToken = undefined;
      this.config.refreshToken = undefined;
      this.refreshSubscribers = [];
      throw this.handleError(refreshError as AxiosError);
    } finally {
      this.isRefreshing = false;
    }
  }

  private handleError(error: AxiosError): APIError {
    const apiError: APIError = {
      error: error.message,
      statusCode: error.response?.status || 0,
    };

    if (error.response?.data) {
      const responseData = error.response.data as any;
      apiError.error = responseData.error || responseData.message || error.message;
      apiError.message = responseData.message;
      apiError.field = responseData.field;
      apiError.code = responseData.code;
    }

    // Call error callback if provided
    if (this.config.onError) {
      this.config.onError(apiError);
    }

    return apiError;
  }

  async get<T>(url: string, config?: RequestConfig): Promise<T> {
    const response = await this.client.get(url, this.mergeConfig(config));
    return response.data;
  }

  async post<T, D = any>(url: string, data?: D, config?: RequestConfig): Promise<T> {
    const response = await this.client.post(url, data, this.mergeConfig(config));
    return response.data;
  }

  async put<T, D = any>(url: string, data?: D, config?: RequestConfig): Promise<T> {
    const response = await this.client.put(url, data, this.mergeConfig(config));
    return response.data;
  }

  async patch<T, D = any>(url: string, data?: D, config?: RequestConfig): Promise<T> {
    const response = await this.client.patch(url, data, this.mergeConfig(config));
    return response.data;
  }

  async delete<T>(url: string, config?: RequestConfig): Promise<T> {
    const response = await this.client.delete(url, this.mergeConfig(config));
    return response.data;
  }

  private mergeConfig(config?: RequestConfig): AxiosRequestConfig {
    return {
      ...config,
      headers: {
        ...config?.headers,
      },
      timeout: config?.timeout || this.config.timeout,
      signal: config?.signal,
    };
  }

  // Public methods for updating configuration
  setAccessToken(token: string): void {
    this.config.accessToken = token;
  }

  setRefreshToken(token: string): void {
    this.config.refreshToken = token;
  }

  clearTokens(): void {
    this.config.accessToken = undefined;
    this.config.refreshToken = undefined;
  }

  getConfig(): ClientConfig {
    return { ...this.config };
  }

  updateConfig(newConfig: Partial<ClientConfig>): void {
    this.config = { ...this.config, ...newConfig };
  }
}