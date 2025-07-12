import { describe, it, expect, beforeEach, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { render, setupMockApi, mockApiResponses, mockAuthenticatedUser, clearMockAuth } from './testUtils'
import App from '@/App'

describe('Error Handling and Edge Cases', () => {
  let mockAxios: any
  let user: any

  beforeEach(() => {
    mockAxios = setupMockApi()
    user = userEvent.setup()
    clearMockAuth()
  })

  describe('Network Error Handling', () => {
    it('handles complete network failure gracefully', async () => {
      // Mock network error
      const networkError = new Error('Network Error')
      networkError.name = 'NetworkError'
      mockAxios.post.mockRejectedValueOnce(networkError)

      render(<App />)

      // Try to login with network failure
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'password123')
      await user.click(submitButton)

      // Verify network error is handled gracefully
      await waitFor(() => {
        expect(screen.getByText(/network error/i)).toBeInTheDocument()
        expect(screen.getByText(/please check your connection/i)).toBeInTheDocument()
      })
    })

    it('handles request timeout errors', async () => {
      // Mock timeout error
      const timeoutError = new Error('Request timeout')
      timeoutError.name = 'TimeoutError'
      mockAxios.get.mockRejectedValueOnce(timeoutError)

      mockAuthenticatedUser()
      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Verify timeout error is handled
      await waitFor(() => {
        expect(screen.getByText(/request timed out/i)).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /retry/i })).toBeInTheDocument()
      })
    })

    it('provides retry mechanism for failed requests', async () => {
      // First call fails, second succeeds
      mockAxios.get
        .mockRejectedValueOnce(new Error('Network Error'))
        .mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      mockAuthenticatedUser()
      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Wait for error to appear
      await waitFor(() => {
        expect(screen.getByText(/network error/i)).toBeInTheDocument()
      })

      // Click retry button
      const retryButton = screen.getByRole('button', { name: /retry/i })
      await user.click(retryButton)

      // Verify successful retry
      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Verify two API calls were made
      expect(mockAxios.get).toHaveBeenCalledTimes(2)
    })
  })

  describe('HTTP Error Status Handling', () => {
    it('handles 401 unauthorized errors with automatic logout', async () => {
      mockAuthenticatedUser()

      // Mock 401 error
      mockAxios.get.mockRejectedValueOnce({
        response: {
          status: 401,
          data: { error: 'Token expired' }
        }
      })

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Verify user is redirected to login
      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
        expect(screen.getByText(/session expired/i)).toBeInTheDocument()
      })

      // Verify auth token is cleared
      expect(localStorage.getItem('auth_token')).toBeNull()
    })

    it('handles 403 forbidden errors', async () => {
      mockAuthenticatedUser()

      mockAxios.get.mockRejectedValueOnce({
        response: {
          status: 403,
          data: { error: 'Access denied' }
        }
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/access denied/i)).toBeInTheDocument()
        expect(screen.getByText(/you don't have permission/i)).toBeInTheDocument()
      })
    })

    it('handles 404 not found errors', async () => {
      mockAuthenticatedUser()

      mockAxios.get.mockRejectedValueOnce({
        response: {
          status: 404,
          data: { error: 'URL not found' }
        }
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics/nonexistent')

      await waitFor(() => {
        expect(screen.getByText(/url not found/i)).toBeInTheDocument()
        expect(screen.getByRole('link', { name: /back to dashboard/i })).toBeInTheDocument()
      })
    })

    it('handles 429 rate limit errors', async () => {
      mockAxios.post.mockRejectedValueOnce({
        response: {
          status: 429,
          data: { 
            error: 'Rate limit exceeded',
            retry_after: 60
          }
        }
      })

      render(<App />)

      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      // Try to login multiple times quickly
      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'password123')
      await user.click(submitButton)

      // Verify rate limit message
      await waitFor(() => {
        expect(screen.getByText(/rate limit exceeded/i)).toBeInTheDocument()
        expect(screen.getByText(/try again in 60 seconds/i)).toBeInTheDocument()
      })
    })

    it('handles 500 server errors', async () => {
      mockAxios.post.mockRejectedValueOnce(mockApiResponses.serverError)

      render(<App />)

      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'password123')
      await user.click(submitButton)

      // Verify server error handling
      await waitFor(() => {
        expect(screen.getByText(/internal server error/i)).toBeInTheDocument()
        expect(screen.getByText(/please try again later/i)).toBeInTheDocument()
      })
    })
  })

  describe('Form Validation Edge Cases', () => {
    it('handles extremely long URLs', async () => {
      mockAuthenticatedUser()
      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      // Try with very long URL (over 2000 characters)
      const longUrl = 'https://example.com/' + 'a'.repeat(2000)
      const urlInput = screen.getByLabelText(/original url/i)
      
      await user.type(urlInput, longUrl)

      // Verify validation error
      await waitFor(() => {
        expect(screen.getByText(/url is too long/i)).toBeInTheDocument()
      })
    })

    it('handles invalid URL formats', async () => {
      mockAuthenticatedUser()
      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      const urlInput = screen.getByLabelText(/original url/i)
      const submitButton = screen.getByRole('button', { name: /shorten url/i })

      // Test various invalid URLs
      const invalidUrls = [
        'not-a-url',
        'ftp://example.com',
        'javascript:alert("xss")',
        'http://',
        'https://',
        'www.example.com', // Missing protocol
      ]

      for (const invalidUrl of invalidUrls) {
        await user.clear(urlInput)
        await user.type(urlInput, invalidUrl)
        await user.click(submitButton)

        await waitFor(() => {
          expect(screen.getByText(/invalid url format/i)).toBeInTheDocument()
        })
      }
    })

    it('handles special characters in custom aliases', async () => {
      mockAuthenticatedUser()
      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      const urlInput = screen.getByLabelText(/original url/i)
      const customAliasCheckbox = screen.getByLabelText(/custom alias/i)

      await user.type(urlInput, 'https://example.com')
      await user.click(customAliasCheckbox)

      const aliasInput = screen.getByLabelText(/custom alias/i)
      const submitButton = screen.getByRole('button', { name: /shorten url/i })

      // Test invalid alias characters
      const invalidAliases = [
        'alias with spaces',
        'alias/with/slashes',
        'alias@with@symbols',
        'alias#with#hash',
        'alias?with?query',
      ]

      for (const invalidAlias of invalidAliases) {
        await user.clear(aliasInput)
        await user.type(aliasInput, invalidAlias)
        await user.click(submitButton)

        await waitFor(() => {
          expect(screen.getByText(/invalid alias format/i)).toBeInTheDocument()
        })
      }
    })
  })

  describe('Data Loading Edge Cases', () => {
    it('handles empty URL list gracefully', async () => {
      mockAuthenticatedUser()

      mockAxios.get.mockResolvedValueOnce({
        data: {
          urls: [],
          total: 0,
          page: 1,
          per_page: 10,
        }
      })

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/no urls found/i)).toBeInTheDocument()
        expect(screen.getByText(/create your first short url/i)).toBeInTheDocument()
      })
    })

    it('handles large URL lists with pagination', async () => {
      mockAuthenticatedUser()

      const largeUrlList = {
        data: {
          urls: Array.from({ length: 50 }, (_, i) => ({
            id: i + 1,
            short_code: `url${i + 1}`,
            original_url: `https://example${i + 1}.com`,
            custom_alias: false,
            expires_at: null,
            is_active: true,
            click_count: i,
            created_at: '2024-01-01T00:00:00Z',
            user_id: 1,
          })),
          total: 250,
          page: 1,
          per_page: 50,
          total_pages: 5,
        }
      }

      mockAxios.get.mockResolvedValueOnce(largeUrlList)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Verify pagination controls appear
      await waitFor(() => {
        expect(screen.getByText(/showing 1-50 of 250/i)).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /next page/i })).toBeInTheDocument()
      })
    })

    it('handles corrupted or malformed data', async () => {
      mockAuthenticatedUser()

      // Mock response with malformed data
      mockAxios.get.mockResolvedValueOnce({
        data: {
          urls: [
            {
              id: 1,
              short_code: null, // Missing required field
              original_url: 'https://example.com',
              // Missing other required fields
            },
            {
              // Completely malformed object
              invalid: 'data'
            }
          ],
          total: 2,
        }
      })

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Verify error handling for malformed data
      await waitFor(() => {
        expect(screen.getByText(/data format error/i)).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /refresh/i })).toBeInTheDocument()
      })
    })
  })

  describe('Browser Compatibility Edge Cases', () => {
    it('handles missing clipboard API gracefully', async () => {
      mockAuthenticatedUser()
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      // Remove clipboard API
      const originalClipboard = navigator.clipboard
      // @ts-ignore
      delete navigator.clipboard

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Try to copy URL
      const copyButtons = screen.getAllByRole('button', { name: /copy/i })
      await user.click(copyButtons[0])

      // Verify fallback mechanism
      await waitFor(() => {
        expect(screen.getByText(/copy not supported/i)).toBeInTheDocument()
      })

      // Restore clipboard API
      navigator.clipboard = originalClipboard
    })

    it('handles missing localStorage gracefully', async () => {
      // Mock localStorage to throw errors
      const originalSetItem = localStorage.setItem
      const originalGetItem = localStorage.getItem
      const originalRemoveItem = localStorage.removeItem

      localStorage.setItem = vi.fn(() => {
        throw new Error('LocalStorage not available')
      })
      localStorage.getItem = vi.fn(() => {
        throw new Error('LocalStorage not available')
      })
      localStorage.removeItem = vi.fn(() => {
        throw new Error('LocalStorage not available')
      })

      render(<App />)

      // Should still render without crashing
      await waitFor(() => {
        expect(screen.getByText(/welcome to url shortener/i)).toBeInTheDocument()
      })

      // Restore localStorage
      localStorage.setItem = originalSetItem
      localStorage.getItem = originalGetItem
      localStorage.removeItem = originalRemoveItem
    })
  })

  describe('Concurrent Operation Handling', () => {
    it('handles multiple simultaneous API calls', async () => {
      mockAuthenticatedUser()

      // Setup multiple API calls that resolve at different times
      const slowResponse = new Promise(resolve => 
        setTimeout(() => resolve(mockApiResponses.urlListSuccess), 1000)
      )
      const fastResponse = Promise.resolve(mockApiResponses.analyticsSuccess)

      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/urls')) {
          return slowResponse
        }
        if (url.includes('/api/analytics')) {
          return fastResponse
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Navigate to analytics quickly
      const analyticsLink = screen.getByRole('link', { name: /analytics/i })
      await user.click(analyticsLink)

      // Verify analytics loads first despite dashboard still loading
      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })
    })

    it('handles race conditions in form submissions', async () => {
      mockAuthenticatedUser()

      // Mock slow API response
      const slowUrlCreation = new Promise(resolve => 
        setTimeout(() => resolve(mockApiResponses.urlCreationSuccess), 2000)
      )

      mockAxios.post.mockReturnValueOnce(slowUrlCreation)
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      const urlInput = screen.getByLabelText(/original url/i)
      const submitButton = screen.getByRole('button', { name: /shorten url/i })

      await user.type(urlInput, 'https://example.com')

      // Click submit multiple times quickly
      await user.click(submitButton)
      await user.click(submitButton)
      await user.click(submitButton)

      // Verify only one API call is made
      expect(mockAxios.post).toHaveBeenCalledTimes(1)

      // Verify button is disabled during submission
      expect(submitButton).toBeDisabled()
    })
  })
})