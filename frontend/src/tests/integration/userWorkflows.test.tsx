import { describe, it, expect, beforeEach, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { render, mockApiResponses, clearMockAuth } from './testUtils'
import App from '@/App'
import axios from 'axios'

// Setup mock axios module
vi.mock('axios', () => {
  const mockAxios = {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    patch: vi.fn(),
    create: vi.fn(),
    interceptors: {
      request: { use: vi.fn(), eject: vi.fn() },
      response: { use: vi.fn(), eject: vi.fn() },
    },
    defaults: {
      headers: {
        common: {},
        delete: {},
        get: {},
        head: {},
        post: {},
        put: {},
        patch: {},
      },
    },
  }
  
  // Re-assign create method to return the same instance
  mockAxios.create = vi.fn(() => mockAxios)
  
  return {
    default: mockAxios,
    create: vi.fn(() => mockAxios),
  }
})

describe('User Authentication Workflows', () => {
  let user: any
  let mockAxios: any

  beforeEach(() => {
    vi.clearAllMocks()
    mockAxios = vi.mocked(axios)
    user = userEvent.setup()
    clearMockAuth()
  })

  describe('User Registration Workflow', () => {
    it('completes full user registration flow', async () => {
      // Mock successful registration API call
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.registrationSuccess)

      // Render the app
      render(<App />)

      // Navigate to registration page
      const signUpLink = screen.getByRole('link', { name: /sign up/i })
      await user.click(signUpLink)

      // Fill out registration form
      await waitFor(() => {
        expect(screen.getByText(/create your account/i)).toBeInTheDocument()
      })

      const nameInput = screen.getByLabelText(/full name/i)
      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/^password$/i)
      const confirmPasswordInput = screen.getByLabelText(/confirm password/i)
      const submitButton = screen.getByRole('button', { name: /create account/i })

      await user.type(nameInput, 'Test User')
      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'securePassword123!')
      await user.type(confirmPasswordInput, 'securePassword123!')

      // Submit the form
      await user.click(submitButton)

      // Verify API call was made with correct data
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/auth/register', {
          name: 'Test User',
          email: 'test@example.com',
          password: 'securePassword123!',
          confirmPassword: 'securePassword123!',
        })
      })

      // Verify user is redirected to dashboard after successful registration
      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })
    })

    it('displays validation errors for invalid registration data', async () => {
      // Mock validation error response
      mockAxios.post.mockRejectedValueOnce(mockApiResponses.validationError)

      render(<App />)

      // Navigate to registration and submit with invalid data
      const signUpLink = screen.getByRole('link', { name: /sign up/i })
      await user.click(signUpLink)

      await waitFor(() => {
        expect(screen.getByText(/create your account/i)).toBeInTheDocument()
      })

      const submitButton = screen.getByRole('button', { name: /create account/i })
      await user.click(submitButton)

      // Verify validation errors are displayed
      await waitFor(() => {
        expect(screen.getByText(/email is required/i)).toBeInTheDocument()
        expect(screen.getByText(/password is too short/i)).toBeInTheDocument()
      })
    })
  })

  describe('User Login Workflow', () => {
    it('completes full user login flow', async () => {
      // Mock successful login API call
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)

      render(<App />)

      // Navigate to login page  
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      // Fill out login form
      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'securePassword123!')

      // Submit the form
      await user.click(submitButton)

      // Verify API call was made with correct data
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/auth/login', {
          email: 'test@example.com',
          password: 'securePassword123!',
        })
      })

      // Verify user is redirected to dashboard after successful login
      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })
    })

    it('displays error for invalid login credentials', async () => {
      // Mock unauthorized error response
      mockAxios.post.mockRejectedValueOnce(mockApiResponses.unauthorizedError)

      render(<App />)

      // Navigate to login and submit with invalid credentials
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'wrong@example.com')
      await user.type(passwordInput, 'wrongpassword')
      await user.click(submitButton)

      // Verify error message is displayed
      await waitFor(() => {
        expect(screen.getByText(/unauthorized/i)).toBeInTheDocument()
      })
    })
  })

  describe('User Logout Workflow', () => {
    it('completes user logout flow', async () => {
      // Mock successful logout API call
      mockAxios.post.mockResolvedValueOnce({ data: { message: 'Logged out successfully' } })

      // Start with an authenticated user
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)

      render(<App />)

      // Login first
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'securePassword123!')
      await user.click(submitButton)

      // Wait for login to complete
      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // Now logout
      const userMenuButton = screen.getByRole('button', { name: /test user/i })
      await user.click(userMenuButton)

      const logoutButton = screen.getByRole('button', { name: /sign out/i })
      await user.click(logoutButton)

      // Verify logout API call was made
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/auth/logout')
      })

      // Verify user is redirected to home page
      await waitFor(() => {
        expect(screen.getByText(/welcome to url shortener/i)).toBeInTheDocument()
      }, { timeout: 3000 })
    })
  })

  describe('Password Reset Workflow', () => {
    it('completes password reset request flow', async () => {
      // Mock successful password reset request
      mockAxios.post.mockResolvedValueOnce({ 
        data: { message: 'Password reset email sent' } 
      })

      render(<App />)

      // Navigate to login page then forgot password
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const forgotPasswordLink = screen.getByRole('link', { name: /forgot your password/i })
      await user.click(forgotPasswordLink)

      // Fill out forgot password form
      await waitFor(() => {
        expect(screen.getByText(/reset your password/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const submitButton = screen.getByRole('button', { name: /send reset link/i })

      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)

      // Verify API call was made
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/auth/forgot-password', {
          email: 'test@example.com',
        })
      })

      // Verify success message is displayed
      await waitFor(() => {
        expect(screen.getByText(/password reset email sent/i)).toBeInTheDocument()
      })
    })
  })

  describe('Protected Route Access', () => {
    it('redirects unauthenticated users to login', async () => {
      render(<App />)

      // Try to access protected dashboard route directly
      window.history.pushState({}, '', '/dashboard')

      // Should be redirected to login
      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })
    })

    it('allows authenticated users to access protected routes', async () => {
      // Mock successful login
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.loginSuccess)
      
      render(<App />)

      // Login first
      const signInLink = screen.getByRole('link', { name: /sign in/i })
      await user.click(signInLink)

      await waitFor(() => {
        expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument()
      })

      const emailInput = screen.getByLabelText(/email address/i)
      const passwordInput = screen.getByLabelText(/password/i)
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await user.type(emailInput, 'test@example.com')
      await user.type(passwordInput, 'securePassword123!')
      await user.click(submitButton)

      // Wait for login and dashboard access
      await waitFor(() => {
        expect(screen.getByText(/dashboard/i)).toBeInTheDocument()
      }, { timeout: 3000 })

      // Now navigate to other protected routes
      const analyticsLink = screen.getByRole('link', { name: /analytics/i })
      await user.click(analyticsLink)

      await waitFor(() => {
        expect(screen.getByText(/analytics/i)).toBeInTheDocument()
      })
    })
  })
})