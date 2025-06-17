import axios, { AxiosResponse } from 'axios'
import {
  LoginRequest,
  RegisterRequest,
  RefreshTokenRequest,
  ChangePasswordRequest,
  UpdateProfileRequest,
  AuthResponse,
  TokenResponse,
  User,
  TokenData,
  ApiError,
} from '@/types/auth'

// API base configuration
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const AUTH_API_URL = `${API_BASE_URL}/api/v1/auth`

// Create axios instance for auth requests
const authApi = axios.create({
  baseURL: AUTH_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
})

// Token storage keys
const TOKEN_STORAGE_KEY = 'auth_tokens'
const USER_STORAGE_KEY = 'auth_user'

// Token management functions
export const tokenManager = {
  getTokens(): TokenData | null {
    try {
      const stored = localStorage.getItem(TOKEN_STORAGE_KEY)
      return stored ? JSON.parse(stored) : null
    } catch {
      return null
    }
  },

  setTokens(tokens: TokenData): void {
    localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(tokens))
  },

  clearTokens(): void {
    localStorage.removeItem(TOKEN_STORAGE_KEY)
  },

  getUser(): User | null {
    try {
      const stored = localStorage.getItem(USER_STORAGE_KEY)
      return stored ? JSON.parse(stored) : null
    } catch {
      return null
    }
  },

  setUser(user: User): void {
    localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(user))
  },

  clearUser(): void {
    localStorage.removeItem(USER_STORAGE_KEY)
  },

  isTokenExpired(expiresAt: string): boolean {
    return new Date(expiresAt) <= new Date()
  },

  isTokenExpiringSoon(expiresAt: string, thresholdMinutes = 5): boolean {
    const expirationTime = new Date(expiresAt).getTime()
    const now = new Date().getTime()
    const threshold = thresholdMinutes * 60 * 1000
    return expirationTime - now < threshold
  },
}

// Error handling
const handleApiError = (error: any): ApiError => {
  if (error.response) {
    return {
      message: error.response.data?.message || 'An error occurred',
      status: error.response.status,
      code: error.response.data?.code,
    }
  } else if (error.request) {
    return {
      message: 'Network error. Please check your connection.',
      status: 0,
    }
  } else {
    return {
      message: error.message || 'An unexpected error occurred',
      status: 0,
    }
  }
}

// Auth service API calls
export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    try {
      const response: AxiosResponse<AuthResponse> = await authApi.post('/login', credentials)
      const authData = response.data

      // Store tokens and user data
      tokenManager.setTokens({
        accessToken: authData.accessToken,
        refreshToken: authData.refreshToken,
        expiresAt: authData.expiresAt,
      })
      tokenManager.setUser(authData.user)

      return authData
    } catch (error) {
      throw handleApiError(error)
    }
  },

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    try {
      const response: AxiosResponse<AuthResponse> = await authApi.post('/register', userData)
      const authData = response.data

      // Store tokens and user data
      tokenManager.setTokens({
        accessToken: authData.accessToken,
        refreshToken: authData.refreshToken,
        expiresAt: authData.expiresAt,
      })
      tokenManager.setUser(authData.user)

      return authData
    } catch (error) {
      throw handleApiError(error)
    }
  },

  async logout(): Promise<void> {
    try {
      const tokens = tokenManager.getTokens()
      if (tokens?.refreshToken) {
        // Send logout request to server to invalidate tokens
        await authApi.post('/logout', { refreshToken: tokens.refreshToken })
      }
    } catch (error) {
      // Log the error but don't throw - we want to clear local data regardless
      console.error('Logout API call failed:', error)
    } finally {
      // Always clear local data
      tokenManager.clearTokens()
      tokenManager.clearUser()
    }
  },

  async refreshToken(): Promise<TokenResponse> {
    try {
      const tokens = tokenManager.getTokens()
      if (!tokens?.refreshToken) {
        throw new Error('No refresh token available')
      }

      const request: RefreshTokenRequest = {
        refreshToken: tokens.refreshToken,
      }

      const response: AxiosResponse<TokenResponse> = await authApi.post('/refresh', request)
      const tokenData = response.data

      // Update stored tokens
      tokenManager.setTokens({
        accessToken: tokenData.accessToken,
        refreshToken: tokenData.refreshToken,
        expiresAt: tokenData.expiresAt,
      })

      return tokenData
    } catch (error) {
      // If refresh fails, clear all auth data
      tokenManager.clearTokens()
      tokenManager.clearUser()
      throw handleApiError(error)
    }
  },

  async validateToken(): Promise<boolean> {
    try {
      const tokens = tokenManager.getTokens()
      if (!tokens?.accessToken) {
        return false
      }

      // Check if token is expired
      if (tokenManager.isTokenExpired(tokens.expiresAt)) {
        return false
      }

      // Make a test request to validate token
      const response = await authApi.get('/validate', {
        headers: {
          Authorization: `Bearer ${tokens.accessToken}`,
        },
      })

      return response.status === 200
    } catch {
      return false
    }
  },

  async getProfile(): Promise<User> {
    try {
      const tokens = tokenManager.getTokens()
      if (!tokens?.accessToken) {
        throw new Error('No access token available')
      }

      const response: AxiosResponse<User> = await authApi.get('/profile', {
        headers: {
          Authorization: `Bearer ${tokens.accessToken}`,
        },
      })

      const user = response.data
      tokenManager.setUser(user)
      return user
    } catch (error) {
      throw handleApiError(error)
    }
  },

  async updateProfile(data: UpdateProfileRequest): Promise<User> {
    try {
      const tokens = tokenManager.getTokens()
      if (!tokens?.accessToken) {
        throw new Error('No access token available')
      }

      const response: AxiosResponse<User> = await authApi.put('/profile', data, {
        headers: {
          Authorization: `Bearer ${tokens.accessToken}`,
        },
      })

      const user = response.data
      tokenManager.setUser(user)
      return user
    } catch (error) {
      throw handleApiError(error)
    }
  },

  async changePassword(data: ChangePasswordRequest): Promise<void> {
    try {
      const tokens = tokenManager.getTokens()
      if (!tokens?.accessToken) {
        throw new Error('No access token available')
      }

      await authApi.post('/change-password', data, {
        headers: {
          Authorization: `Bearer ${tokens.accessToken}`,
        },
      })
    } catch (error) {
      throw handleApiError(error)
    }
  },
}

// Automatic token refresh functionality
export const setupTokenRefresh = (onTokenRefreshed?: () => void, onAuthExpired?: () => void) => {
  const checkAndRefreshToken = async () => {
    const tokens = tokenManager.getTokens()
    if (!tokens) return

    // If token is expired, try to refresh
    if (tokenManager.isTokenExpired(tokens.expiresAt)) {
      try {
        await authService.refreshToken()
        onTokenRefreshed?.()
      } catch {
        onAuthExpired?.()
      }
    }
    // If token is expiring soon, proactively refresh
    else if (tokenManager.isTokenExpiringSoon(tokens.expiresAt)) {
      try {
        await authService.refreshToken()
        onTokenRefreshed?.()
      } catch {
        // If proactive refresh fails, we can continue with current token
        console.warn('Proactive token refresh failed')
      }
    }
  }

  // Check token every minute
  const interval = setInterval(checkAndRefreshToken, 60 * 1000)
  
  // Initial check
  checkAndRefreshToken()

  // Return cleanup function
  return () => clearInterval(interval)
}