import { describe, it, expect, beforeEach, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { render, setupMockApi, mockApiResponses, clearMockAuth } from './testUtils'
import App from '@/App'

describe('End-to-End Component Interaction Workflows', () => {
  let mockAxios: any
  let user: any

  beforeEach(() => {
    mockAxios = setupMockApi()
    user = userEvent.setup()
    clearMockAuth()
  })

  describe('Complete User Journey: Registration to URL Management', () => {
    it('completes full user journey from registration to creating and managing URLs', async () => {
      // Step 1: Mock registration
      mockAxios.post.mockImplementation((url, data) => {
        if (url === '/api/auth/register') {
          return Promise.resolve(mockApiResponses.registrationSuccess)
        }
        if (url === '/api/urls') {
          return Promise.resolve(mockApiResponses.urlCreationSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      // Step 2: Mock dashboard data loading
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/urls')) {
          return Promise.resolve(mockApiResponses.urlListSuccess)
        }
        if (url.includes('/api/analytics')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      // STEP 1: User registers
      const signUpLink = screen.getByRole('link', { name: /sign up/i })
      await user.click(signUpLink)

      await waitFor(() => {
        expect(screen.getByText(/create your account/i)).toBeInTheDocument()
      })

      const nameInput = screen.getByLabelText(/full name/i)
      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/^password$/i)
      const confirmPasswordInput = screen.getByLabelText(/confirm password/i)
      const submitButton = screen.getByRole('button', { name: /create account/i })

      await user.type(nameInput, 'John Doe')
      await user.type(emailInput, 'john@example.com')
      await user.type(passwordInput, 'securePassword123!')
      await user.type(confirmPasswordInput, 'securePassword123!')
      await user.click(submitButton)

      // STEP 2: User lands on dashboard after registration
      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
        expect(screen.getByText(/welcome, john/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // STEP 3: User sees existing URLs and creates a new one
      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      const urlInput = screen.getByLabelText(/original url/i)
      const shortenButton = screen.getByRole('button', { name: /shorten url/i })

      await user.type(urlInput, 'https://my-new-website.com/long-path')
      await user.click(shortenButton)

      // STEP 4: Verify URL creation success
      await waitFor(() => {
        expect(screen.getByText(/url shortened successfully/i)).toBeInTheDocument()
      })

      // STEP 5: Navigate to analytics
      const analyticsLink = screen.getByRole('link', { name: /analytics/i })
      await user.click(analyticsLink)

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
        expect(screen.getByText(/total clicks/i)).toBeInTheDocument()
      })

      // Verify the complete journey worked
      expect(mockAxios.post).toHaveBeenCalledWith('/api/auth/register', expect.any(Object))
      expect(mockAxios.post).toHaveBeenCalledWith('/api/urls', expect.any(Object))
      expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/overview')
    })
  })

  describe('Complete User Journey: Login to Analytics Deep Dive', () => {
    it('completes user journey from login to detailed analytics analysis', async () => {
      // Mock login and subsequent API calls
      mockAxios.post.mockImplementation((url, data) => {
        if (url === '/api/auth/login') {
          return Promise.resolve(mockApiResponses.loginSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/urls')) {
          return Promise.resolve(mockApiResponses.urlListSuccess)
        }
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        if (url.includes('/api/analytics/abc123')) {
          return Promise.resolve({
            data: {
              url: {
                id: 1,
                short_code: 'abc123',
                original_url: 'https://example.com',
                click_count: 25,
              },
              analytics: mockApiResponses.analyticsSuccess.data,
            }
          })
        }
        if (url.includes('/api/analytics/export')) {
          return Promise.resolve({
            data: 'Date,Clicks\n2024-01-01,5\n2024-01-02,10',
            headers: {
              'content-type': 'text/csv'
            }
          })
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      // STEP 1: User logs in
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'john@example.com')
      await user.type(passwordInput, 'securePassword123!')
      await user.click(submitButton)

      // STEP 2: Navigate to analytics
      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      const analyticsLink = screen.getByRole('link', { name: /analytics/i })
      await user.click(analyticsLink)

      // STEP 3: View overall analytics
      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
        expect(screen.getByText(/15/)).toBeInTheDocument() // total clicks
      })

      // STEP 4: Change time range
      const timeRangeSelect = screen.getByRole('combobox', { name: /time range/i })
      await user.selectOptions(timeRangeSelect, '7d')

      // STEP 5: View specific URL analytics
      const viewDetailsButton = screen.getByRole('button', { name: /view details/i })
      await user.click(viewDetailsButton)

      await waitFor(() => {
        expect(screen.getByText(/analytics for abc123/i)).toBeInTheDocument()
      })

      // STEP 6: Export analytics data
      const exportButton = screen.getByRole('button', { name: /export/i })
      await user.click(exportButton)

      const csvOption = screen.getByRole('button', { name: /csv/i })
      await user.click(csvOption)

      // Verify the complete analytics journey
      expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/overview')
      expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/abc123')
      expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/export', expect.any(Object))
    })
  })

  describe('Multi-Component State Synchronization', () => {
    it('synchronizes state changes across multiple components', async () => {
      // Setup authenticated user
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/urls')) {
          return Promise.resolve(mockApiResponses.urlListSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      // Login
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

      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // Verify header shows user info
      expect(screen.getByText(/test user/i)).toBeInTheDocument()

      // Verify dashboard shows URL count
      await waitFor(() => {
        expect(screen.getByText(/2 urls/i)).toBeInTheDocument()
      })

      // Create new URL
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.urlCreationSuccess)
      mockAxios.get.mockResolvedValueOnce({
        data: {
          ...mockApiResponses.urlListSuccess.data,
          total: 3,
          urls: [...mockApiResponses.urlListSuccess.data.urls, mockApiResponses.urlCreationSuccess.data]
        }
      })

      const urlInput = screen.getByLabelText(/original url/i)
      const shortenButton = screen.getByRole('button', { name: /shorten url/i })

      await user.type(urlInput, 'https://new-example.com')
      await user.click(shortenButton)

      // Verify count updates in header and dashboard
      await waitFor(() => {
        expect(screen.getByText(/3 urls/i)).toBeInTheDocument()
      })
    })

    it('maintains state consistency during navigation', async () => {
      // Setup with authentication
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/urls')) {
          return Promise.resolve(mockApiResponses.urlListSuccess)
        }
        if (url.includes('/api/analytics')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      // Login
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

      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // Navigate between pages and verify state persists
      const analyticsLink = screen.getByRole('link', { name: /analytics/i })
      await user.click(analyticsLink)

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Navigate back to dashboard
      const dashboardLink = screen.getByRole('link', { name: /dashboard/i })
      await user.click(dashboardLink)

      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      })

      // Verify user still appears as logged in
      expect(screen.getByText(/test user/i)).toBeInTheDocument()
    })
  })

  describe('Real-time Feature Integration', () => {
    it('integrates real-time updates across components', async () => {
      // Mock WebSocket-like behavior with polling
      let pollCount = 0
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/urls')) {
          return Promise.resolve(mockApiResponses.urlListSuccess)
        }
        if (url.includes('/api/analytics/overview')) {
          pollCount++
          if (pollCount === 1) {
            return Promise.resolve(mockApiResponses.analyticsSuccess)
          } else {
            // Return updated data on subsequent polls
            return Promise.resolve({
              data: {
                ...mockApiResponses.analyticsSuccess.data,
                total_clicks: 20, // Increased from 15
                unique_clicks: 16, // Increased from 12
              }
            })
          }
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)

      render(<App />)

      // Login and navigate to analytics
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

      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      const analyticsLink = screen.getByRole('link', { name: /analytics/i })
      await user.click(analyticsLink)

      // Verify initial data
      await waitFor(() => {
        expect(screen.getByText(/15/)).toBeInTheDocument() // total clicks
        expect(screen.getByText(/12/)).toBeInTheDocument() // unique clicks
      })

      // Simulate real-time update (component should poll for updates)
      await waitFor(() => {
        expect(screen.getByText(/20/)).toBeInTheDocument() // updated total clicks
        expect(screen.getByText(/16/)).toBeInTheDocument() // updated unique clicks
      }, { timeout: 6000 }) // Allow time for polling interval
    })
  })

  describe('Complex Form Interactions', () => {
    it('handles complex multi-step form workflows', async () => {
      mockAxios.post.mockImplementation((url, data) => {
        if (url === '/api/auth/login') {
          return Promise.resolve(mockApiResponses.loginSuccess)
        }
        if (url === '/api/urls') {
          return Promise.resolve({
            data: {
              ...mockApiResponses.urlCreationSuccess.data,
              custom_alias: true,
              short_code: data.custom_alias,
              expires_at: data.expires_at,
              password_protected: !!data.password,
            }
          })
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      // Login
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

      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // Create URL with advanced options
      const urlInput = screen.getByLabelText(/original url/i)
      await user.type(urlInput, 'https://premium-content.com')

      // Enable custom alias
      const customAliasCheckbox = screen.getByLabelText(/custom alias/i)
      await user.click(customAliasCheckbox)

      const aliasInput = screen.getByLabelText(/custom alias/i)
      await user.type(aliasInput, 'premium-link')

      // Enable expiration
      const expirationCheckbox = screen.getByLabelText(/set expiration/i)
      await user.click(expirationCheckbox)

      const expirationInput = screen.getByLabelText(/expiration date/i)
      await user.type(expirationInput, '2024-12-31')

      // Enable password protection
      const passwordCheckbox = screen.getByLabelText(/password protection/i)
      await user.click(passwordCheckbox)

      const linkPasswordInput = screen.getByLabelText(/link password/i)
      await user.type(linkPasswordInput, 'secret123')

      // Submit complex form
      const shortenButton = screen.getByRole('button', { name: /shorten url/i })
      await user.click(shortenButton)

      // Verify complex URL creation
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/urls', {
          original_url: 'https://premium-content.com',
          custom_alias: 'premium-link',
          expires_at: '2024-12-31',
          password: 'secret123',
        })
      })

      await waitFor(() => {
        expect(screen.getByText(/url shortened successfully/i)).toBeInTheDocument()
      })
    })
  })

  describe('Progressive Enhancement Integration', () => {
    it('gracefully degrades when JavaScript features are limited', async () => {
      // Mock limited clipboard support
      const originalClipboard = navigator.clipboard
      // @ts-ignore
      delete navigator.clipboard

      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      // Login
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

      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // Try to copy URL
      const copyButtons = screen.getAllByRole('button', { name: /copy/i })
      await user.click(copyButtons[0])

      // Should show fallback message
      await waitFor(() => {
        expect(screen.getByText(/copy not supported/i)).toBeInTheDocument()
      })

      // Restore clipboard
      navigator.clipboard = originalClipboard
    })
  })
})