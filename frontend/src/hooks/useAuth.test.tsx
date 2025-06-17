import { describe, it, expect, vi } from 'vitest'
import { renderHook } from '@testing-library/react'
import { ReactNode } from 'react'
import { useAuth } from './useAuth'
import { AuthProvider } from '@/context/AuthContext'

// Mock the auth service
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
  setupTokenRefresh: vi.fn(),
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

describe('useAuth Hook', () => {
  const wrapper = ({ children }: { children: ReactNode }) => (
    <AuthProvider>{children}</AuthProvider>
  )

  it('should throw error when used outside AuthProvider', () => {
    // Suppress console.error for this test since we expect an error
    const originalError = console.error
    console.error = vi.fn()

    expect(() => {
      renderHook(() => useAuth())
    }).toThrow('useAuth must be used within an AuthProvider')

    console.error = originalError
  })

  it('should return auth context when used within AuthProvider', () => {
    const { result } = renderHook(() => useAuth(), { wrapper })

    expect(result.current).toMatchObject({
      user: null,
      isAuthenticated: false,
      isLoading: expect.any(Boolean),
      error: null,
      login: expect.any(Function),
      register: expect.any(Function),
      logout: expect.any(Function),
      refreshAuth: expect.any(Function),
      updateProfile: expect.any(Function),
      changePassword: expect.any(Function),
      clearError: expect.any(Function),
      validateToken: expect.any(Function),
    })
  })

  it('should provide all required auth methods', () => {
    const { result } = renderHook(() => useAuth(), { wrapper })

    // Check that all methods are functions
    expect(typeof result.current.login).toBe('function')
    expect(typeof result.current.register).toBe('function')
    expect(typeof result.current.logout).toBe('function')
    expect(typeof result.current.refreshAuth).toBe('function')
    expect(typeof result.current.updateProfile).toBe('function')
    expect(typeof result.current.changePassword).toBe('function')
    expect(typeof result.current.clearError).toBe('function')
    expect(typeof result.current.validateToken).toBe('function')
  })

  it('should provide correct initial state', () => {
    localStorageMock.getItem.mockReturnValue(null) // No stored auth data

    const { result } = renderHook(() => useAuth(), { wrapper })

    expect(result.current.user).toBeNull()
    expect(result.current.isAuthenticated).toBe(false)
    expect(result.current.error).toBeNull()
    expect(typeof result.current.isLoading).toBe('boolean')
  })
})