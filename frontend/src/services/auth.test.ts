import { describe, it, expect, vi, beforeEach } from 'vitest'
import { tokenManager } from './auth'

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

describe('TokenManager', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('token management', () => {
    it('should store and retrieve tokens', () => {
      const tokenData = {
        accessToken: 'access-token',
        refreshToken: 'refresh-token',
        expiresAt: '2024-12-31T23:59:59Z',
      }

      tokenManager.setTokens(tokenData)
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'auth_tokens',
        JSON.stringify(tokenData)
      )

      localStorageMock.getItem.mockReturnValue(JSON.stringify(tokenData))
      const result = tokenManager.getTokens()
      expect(result).toEqual(tokenData)
    })

    it('should clear tokens', () => {
      tokenManager.clearTokens()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_tokens')
    })

    it('should handle corrupted token data gracefully', () => {
      localStorageMock.getItem.mockReturnValue('invalid-json')
      const result = tokenManager.getTokens()
      expect(result).toBeNull()
    })
  })

  describe('user management', () => {
    it('should store and retrieve user data', () => {
      const userData = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        isActive: true,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }

      tokenManager.setUser(userData)
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'auth_user',
        JSON.stringify(userData)
      )

      localStorageMock.getItem.mockReturnValue(JSON.stringify(userData))
      const result = tokenManager.getUser()
      expect(result).toEqual(userData)
    })

    it('should clear user data', () => {
      tokenManager.clearUser()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_user')
    })
  })

  describe('token expiration', () => {
    it('should correctly identify expired tokens', () => {
      const expiredDate = new Date(Date.now() - 3600000).toISOString() // 1 hour ago
      expect(tokenManager.isTokenExpired(expiredDate)).toBe(true)

      const futureDate = new Date(Date.now() + 3600000).toISOString() // 1 hour from now
      expect(tokenManager.isTokenExpired(futureDate)).toBe(false)
    })

    it('should correctly identify tokens expiring soon', () => {
      const soonDate = new Date(Date.now() + 2 * 60 * 1000).toISOString() // 2 minutes from now
      expect(tokenManager.isTokenExpiringSoon(soonDate, 5)).toBe(true)

      const farDate = new Date(Date.now() + 10 * 60 * 1000).toISOString() // 10 minutes from now
      expect(tokenManager.isTokenExpiringSoon(farDate, 5)).toBe(false)
    })
  })
})