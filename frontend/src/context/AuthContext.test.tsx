import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { ReactNode } from 'react'
import { AuthProvider } from './AuthContext'
import { useAuth } from '@/hooks/useAuth'
import * as authService from '@/services/auth'

// Mock the entire auth service module
vi.mock('@/services/auth', () => ({
  authService: {
    login: vi.fn(),
    register: vi.fn(),
    logout: vi.fn(),
    refreshToken: vi.fn(),
    validateToken: vi.fn(),
    getProfile: vi.fn(),
    updateProfile: vi.fn(),
    changePassword: vi.fn(),
  },
  tokenManager: {
    getTokens: vi.fn(),
    setTokens: vi.fn(),
    clearTokens: vi.fn(),
    getUser: vi.fn(),
    setUser: vi.fn(),
    clearUser: vi.fn(),
    isTokenExpired: vi.fn(),
    isTokenExpiringSoon: vi.fn(),
  },
  setupTokenRefresh: vi.fn(() => vi.fn()), // Return cleanup function
}))

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('AuthContext', () => {
  const mockUser = {
    id: '1',
    email: 'test@example.com',
    name: 'Test User',
    isActive: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  }

  const mockAuthResponse = {
    user: mockUser,
    accessToken: 'mock-access-token',
    refreshToken: 'mock-refresh-token',
    expiresAt: new Date(Date.now() + 3600000).toISOString(),
  }

  const wrapper = ({ children }: { children: ReactNode }) => (
    <AuthProvider>{children}</AuthProvider>
  )

  beforeEach(() => {
    vi.clearAllMocks()
    // Default mock implementations
    vi.mocked(authService.tokenManager.getTokens).mockReturnValue(null)
    vi.mocked(authService.tokenManager.getUser).mockReturnValue(null)
    vi.mocked(authService.tokenManager.isTokenExpired).mockReturnValue(false)
    vi.mocked(authService.authService.validateToken).mockResolvedValue(true)
  })

  describe('initialization', () => {
    it('should initialize with correct initial state', () => {
      const { result } = renderHook(() => useAuth(), { wrapper })
      expect(typeof result.current.isLoading).toBe('boolean')
      expect(result.current.isAuthenticated).toBe(false)
      expect(result.current.user).toBeNull()
    })

    it('should initialize without auth when no tokens exist', async () => {
      const { result } = renderHook(() => useAuth(), { wrapper })

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      expect(result.current.isAuthenticated).toBe(false)
      expect(result.current.user).toBeNull()
    })

    it('should initialize with auth when valid tokens exist', async () => {
      vi.mocked(authService.tokenManager.getTokens).mockReturnValue({
        accessToken: 'valid-token',
        refreshToken: 'valid-refresh',
        expiresAt: new Date(Date.now() + 3600000).toISOString(),
      })
      vi.mocked(authService.tokenManager.getUser).mockReturnValue(mockUser)
      vi.mocked(authService.authService.validateToken).mockResolvedValue(true)

      const { result } = renderHook(() => useAuth(), { wrapper })

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      expect(result.current.isAuthenticated).toBe(true)
      expect(result.current.user).toEqual(mockUser)
    })

    it('should clear auth when token validation fails', async () => {
      vi.mocked(authService.tokenManager.getTokens).mockReturnValue({
        accessToken: 'invalid-token',
        refreshToken: 'invalid-refresh',
        expiresAt: new Date(Date.now() + 3600000).toISOString(),
      })
      vi.mocked(authService.tokenManager.getUser).mockReturnValue(mockUser)
      vi.mocked(authService.authService.validateToken).mockResolvedValue(false)

      const { result } = renderHook(() => useAuth(), { wrapper })

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      expect(result.current.isAuthenticated).toBe(false)
      expect(result.current.user).toBeNull()
      expect(authService.tokenManager.clearTokens).toHaveBeenCalled()
      expect(authService.tokenManager.clearUser).toHaveBeenCalled()
    })
  })

  describe('login', () => {
    it('should login successfully', async () => {
      vi.mocked(authService.authService.login).mockResolvedValue(mockAuthResponse)

      const { result } = renderHook(() => useAuth(), { wrapper })

      await act(async () => {
        await result.current.login({
          email: 'test@example.com',
          password: 'password123',
        })
      })

      expect(result.current.isAuthenticated).toBe(true)
      expect(result.current.user).toEqual(mockUser)
      expect(result.current.error).toBeNull()
      expect(result.current.isLoading).toBe(false)
    })

    it('should handle login errors', async () => {
      const loginError = { message: 'Invalid credentials', status: 401 }
      vi.mocked(authService.authService.login).mockRejectedValue(loginError)

      const { result } = renderHook(() => useAuth(), { wrapper })

      await act(async () => {
        try {
          await result.current.login({
            email: 'test@example.com',
            password: 'wrong-password',
          })
        } catch (error) {
          // Expected to throw
        }
      })

      expect(result.current.isAuthenticated).toBe(false)
      expect(result.current.user).toBeNull()
      expect(result.current.error).toBe('Invalid credentials')
      expect(result.current.isLoading).toBe(false)
    })
  })

  describe('register', () => {
    it('should register successfully', async () => {
      vi.mocked(authService.authService.register).mockResolvedValue(mockAuthResponse)

      const { result } = renderHook(() => useAuth(), { wrapper })

      await act(async () => {
        await result.current.register({
          name: 'Test User',
          email: 'test@example.com',
          password: 'password123',
          confirmPassword: 'password123',
        })
      })

      expect(result.current.isAuthenticated).toBe(true)
      expect(result.current.user).toEqual(mockUser)
      expect(result.current.error).toBeNull()
    })
  })

  describe('logout', () => {
    it('should logout successfully', async () => {
      // Start with authenticated state
      vi.mocked(authService.authService.login).mockResolvedValue(mockAuthResponse)
      vi.mocked(authService.authService.logout).mockResolvedValue()

      const { result } = renderHook(() => useAuth(), { wrapper })

      // Login first
      await act(async () => {
        await result.current.login({
          email: 'test@example.com',
          password: 'password123',
        })
      })

      expect(result.current.isAuthenticated).toBe(true)

      // Then logout
      await act(async () => {
        await result.current.logout()
      })

      expect(result.current.isAuthenticated).toBe(false)
      expect(result.current.user).toBeNull()
      expect(authService.authService.logout).toHaveBeenCalled()
    })
  })

  describe('updateProfile', () => {
    it('should update profile successfully', async () => {
      const updatedUser = { ...mockUser, name: 'Updated Name' }
      vi.mocked(authService.authService.login).mockResolvedValue(mockAuthResponse)
      vi.mocked(authService.authService.updateProfile).mockResolvedValue(updatedUser)

      const { result } = renderHook(() => useAuth(), { wrapper })

      // Login first
      await act(async () => {
        await result.current.login({
          email: 'test@example.com',
          password: 'password123',
        })
      })

      // Update profile
      await act(async () => {
        await result.current.updateProfile({
          name: 'Updated Name',
          email: 'test@example.com',
        })
      })

      expect(result.current.user).toEqual(updatedUser)
      expect(authService.authService.updateProfile).toHaveBeenCalledWith({
        name: 'Updated Name',
        email: 'test@example.com',
      })
    })
  })

  describe('changePassword', () => {
    it('should change password successfully', async () => {
      vi.mocked(authService.authService.login).mockResolvedValue(mockAuthResponse)
      vi.mocked(authService.authService.changePassword).mockResolvedValue()

      const { result } = renderHook(() => useAuth(), { wrapper })

      // Login first
      await act(async () => {
        await result.current.login({
          email: 'test@example.com',
          password: 'password123',
        })
      })

      // Change password
      await act(async () => {
        await result.current.changePassword({
          currentPassword: 'password123',
          newPassword: 'newpassword123',
          confirmPassword: 'newpassword123',
        })
      })

      expect(authService.authService.changePassword).toHaveBeenCalledWith({
        currentPassword: 'password123',
        newPassword: 'newpassword123',
        confirmPassword: 'newpassword123',
      })
    })
  })

  describe('clearError', () => {
    it('should clear error state', async () => {
      const loginError = { message: 'Test error', status: 500 }
      vi.mocked(authService.authService.login).mockRejectedValue(loginError)

      const { result } = renderHook(() => useAuth(), { wrapper })

      // Cause an error
      await act(async () => {
        try {
          await result.current.login({
            email: 'test@example.com',
            password: 'wrong-password',
          })
        } catch (error) {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Test error')

      // Clear error
      act(() => {
        result.current.clearError()
      })

      expect(result.current.error).toBeNull()
    })
  })

  describe('validateToken', () => {
    it('should validate token successfully', async () => {
      vi.mocked(authService.authService.validateToken).mockResolvedValue(true)

      const { result } = renderHook(() => useAuth(), { wrapper })

      let isValid: boolean
      await act(async () => {
        isValid = await result.current.validateToken()
      })

      expect(isValid!).toBe(true)
      expect(authService.authService.validateToken).toHaveBeenCalled()
    })

    it('should handle token validation failure', async () => {
      vi.mocked(authService.authService.validateToken).mockResolvedValue(false)

      const { result } = renderHook(() => useAuth(), { wrapper })

      let isValid: boolean
      await act(async () => {
        isValid = await result.current.validateToken()
      })

      expect(isValid!).toBe(false)
    })
  })
})