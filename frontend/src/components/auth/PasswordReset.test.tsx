import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import PasswordReset from './PasswordReset'
import { authService } from '@/services/auth'

// Mock the auth service
vi.mock('@/services/auth', () => ({
  authService: {
    forgotPassword: vi.fn(),
    resetPassword: vi.fn()
  }
}))

// Mock React Router components
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useSearchParams: vi.fn(() => [new URLSearchParams(), vi.fn()]),
    Link: ({ children, to, ...props }: any) => (
      <a href={to} {...props}>{children}</a>
    )
  }
})

const mockForgotPassword = vi.mocked(authService.forgotPassword)
const mockResetPassword = vi.mocked(authService.resetPassword)

const renderPasswordReset = (searchParams = '', props = {}) => {
  const { useSearchParams } = require('react-router-dom')
  const mockSearchParams = new URLSearchParams(searchParams)
  useSearchParams.mockReturnValue([mockSearchParams, vi.fn()])
  
  return render(
    <BrowserRouter>
      <PasswordReset {...props} />
    </BrowserRouter>
  )
}

describe('PasswordReset', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Forgot Password Step', () => {
    it('renders forgot password form', () => {
      renderPasswordReset()
      
      expect(screen.getByText('Forgot Password?')).toBeInTheDocument()
      expect(screen.getByText("Enter your email address and we'll send you a link to reset your password.")).toBeInTheDocument()
      expect(screen.getByLabelText('Email Address')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /send reset link/i })).toBeInTheDocument()
    })

    it('displays validation error for empty email', async () => {
      const user = userEvent.setup()
      renderPasswordReset()
      
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Email is required')).toBeInTheDocument()
      })
    })

    it('displays validation error for invalid email format', async () => {
      const user = userEvent.setup()
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      await user.type(emailInput, 'invalid-email')
      
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Please enter a valid email address')).toBeInTheDocument()
      })
    })

    it('submits forgot password request', async () => {
      const user = userEvent.setup()
      mockForgotPassword.mockResolvedValue(undefined)
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(mockForgotPassword).toHaveBeenCalledWith({
          email: 'test@example.com'
        })
      })
    })

    it('shows success message after sending reset link', async () => {
      const user = userEvent.setup()
      mockForgotPassword.mockResolvedValue(undefined)
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Check Your Email')).toBeInTheDocument()
        expect(screen.getByText('test@example.com')).toBeInTheDocument()
      })
    })

    it('shows loading state during submission', async () => {
      const user = userEvent.setup()
      mockForgotPassword.mockImplementation(() => new Promise(() => {})) // Never resolves
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)
      
      expect(screen.getByText('Sending Reset Link...')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /sending reset link/i })).toBeDisabled()
    })

    it('handles 404 error (email not found)', async () => {
      const user = userEvent.setup()
      const mockError = { status: 404 }
      mockForgotPassword.mockRejectedValue(mockError)
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'notfound@example.com')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('No account found with this email address')).toBeInTheDocument()
      })
    })

    it('handles 429 error (rate limiting)', async () => {
      const user = userEvent.setup()
      const mockError = { status: 429 }
      mockForgotPassword.mockRejectedValue(mockError)
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Too many password reset requests. Please try again later.')).toBeInTheDocument()
      })
    })

    it('renders navigation links', () => {
      renderPasswordReset()
      expect(screen.getByText('Back to Sign In')).toBeInTheDocument()
    })
  })

  describe('Email Sent Step', () => {
    it('renders email sent confirmation with resend option', async () => {
      const user = userEvent.setup()
      mockForgotPassword.mockResolvedValue(undefined)
      renderPasswordReset()
      
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Check Your Email')).toBeInTheDocument()
        expect(screen.getByText('Send Another Email')).toBeInTheDocument()
        expect(screen.getByText('Back to Sign In')).toBeInTheDocument()
      })
    })

    it('allows resending email', async () => {
      const user = userEvent.setup()
      mockForgotPassword.mockResolvedValue(undefined)
      renderPasswordReset()
      
      // First send
      const emailInput = screen.getByLabelText('Email Address')
      const submitButton = screen.getByRole('button', { name: /send reset link/i })
      
      await user.type(emailInput, 'test@example.com')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Check Your Email')).toBeInTheDocument()
      })
      
      // Resend
      const resendButton = screen.getByText('Send Another Email')
      await user.click(resendButton)
      
      expect(screen.getByText('Forgot Password?')).toBeInTheDocument()
    })
  })

  describe('Reset Password Step', () => {
    it('renders reset password form when token is present', () => {
      renderPasswordReset('token=abc123')
      
      expect(screen.getByText('Reset Password')).toBeInTheDocument()
      expect(screen.getByText('Enter your new password below')).toBeInTheDocument()
      expect(screen.getByLabelText('New Password')).toBeInTheDocument()
      expect(screen.getByLabelText('Confirm New Password')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /reset password/i })).toBeInTheDocument()
    })

    it('displays validation errors for empty fields', async () => {
      const user = userEvent.setup()
      renderPasswordReset('token=abc123')
      
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Password is required')).toBeInTheDocument()
        expect(screen.getByText('Please confirm your password')).toBeInTheDocument()
      })
    })

    it('displays validation error for weak password', async () => {
      const user = userEvent.setup()
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      await user.type(passwordInput, 'weak')
      
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Password must be at least 8 characters')).toBeInTheDocument()
      })
    })

    it('displays validation error for password mismatch', async () => {
      const user = userEvent.setup()
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      
      await user.type(passwordInput, 'Password123')
      await user.type(confirmPasswordInput, 'Password456')
      
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText("Passwords don't match")).toBeInTheDocument()
      })
    })

    it('shows password strength indicator', async () => {
      const user = userEvent.setup()
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      
      // Type a weak password
      await user.type(passwordInput, 'weak')
      expect(screen.getByText('Weak')).toBeInTheDocument()
      
      // Type a stronger password
      await user.clear(passwordInput)
      await user.type(passwordInput, 'Password123')
      expect(screen.getByText('Good')).toBeInTheDocument()
    })

    it('toggles password visibility', async () => {
      const user = userEvent.setup()
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const toggleButtons = screen.getAllByRole('button', { name: '' })
      const passwordToggle = toggleButtons[0]
      
      expect(passwordInput).toHaveAttribute('type', 'password')
      
      await user.click(passwordToggle)
      expect(passwordInput).toHaveAttribute('type', 'text')
      
      await user.click(passwordToggle)
      expect(passwordInput).toHaveAttribute('type', 'password')
    })

    it('submits reset password request', async () => {
      const user = userEvent.setup()
      mockResetPassword.mockResolvedValue(undefined)
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      
      await user.type(passwordInput, 'NewPassword123')
      await user.type(confirmPasswordInput, 'NewPassword123')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(mockResetPassword).toHaveBeenCalledWith({
          token: 'abc123',
          newPassword: 'NewPassword123',
          confirmPassword: 'NewPassword123'
        })
      })
    })

    it('shows success message after password reset', async () => {
      const user = userEvent.setup()
      mockResetPassword.mockResolvedValue(undefined)
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      
      await user.type(passwordInput, 'NewPassword123')
      await user.type(confirmPasswordInput, 'NewPassword123')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Password Reset Successful')).toBeInTheDocument()
        expect(screen.getByText('Your password has been successfully reset')).toBeInTheDocument()
      })
    })

    it('shows loading state during submission', async () => {
      const user = userEvent.setup()
      mockResetPassword.mockImplementation(() => new Promise(() => {})) // Never resolves
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      
      await user.type(passwordInput, 'NewPassword123')
      await user.type(confirmPasswordInput, 'NewPassword123')
      await user.click(submitButton)
      
      expect(screen.getByText('Resetting Password...')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /resetting password/i })).toBeDisabled()
    })

    it('handles 400 error (invalid token)', async () => {
      const user = userEvent.setup()
      const mockError = { status: 400 }
      mockResetPassword.mockRejectedValue(mockError)
      renderPasswordReset('token=invalid')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      
      await user.type(passwordInput, 'NewPassword123')
      await user.type(confirmPasswordInput, 'NewPassword123')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Invalid or expired reset token. Please request a new password reset.')).toBeInTheDocument()
      })
    })

    it('handles 429 error (rate limiting)', async () => {
      const user = userEvent.setup()
      const mockError = { status: 429 }
      mockResetPassword.mockRejectedValue(mockError)
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      
      await user.type(passwordInput, 'NewPassword123')
      await user.type(confirmPasswordInput, 'NewPassword123')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Too many password reset attempts. Please try again later.')).toBeInTheDocument()
      })
    })
  })

  describe('Success Step', () => {
    it('renders success message with sign in link', async () => {
      const user = userEvent.setup()
      mockResetPassword.mockResolvedValue(undefined)
      renderPasswordReset('token=abc123')
      
      const passwordInput = screen.getByLabelText('New Password')
      const confirmPasswordInput = screen.getByLabelText('Confirm New Password')
      const submitButton = screen.getByRole('button', { name: /reset password/i })
      
      await user.type(passwordInput, 'NewPassword123')
      await user.type(confirmPasswordInput, 'NewPassword123')
      await user.click(submitButton)
      
      await waitFor(() => {
        expect(screen.getByText('Password Reset Successful')).toBeInTheDocument()
        expect(screen.getByText('Sign In')).toBeInTheDocument()
      })
    })
  })

  it('applies custom className', () => {
    const { container } = renderPasswordReset('', { className: 'custom-class' })
    expect(container.firstChild).toHaveClass('custom-class')
  })
})