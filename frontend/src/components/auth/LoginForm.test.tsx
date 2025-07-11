import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import LoginForm from './LoginForm'
import { useAuth } from '@/hooks/useAuth'

// Mock the useAuth hook
vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn()
}))

// Mock React Router components
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    Link: ({ children, to, ...props }: any) => (
      <a href={to} {...props}>{children}</a>
    )
  }
})

const mockLogin = vi.fn()
const mockAuthHook = {
  login: mockLogin,
  isLoading: false,
  error: null,
  user: null,
  isAuthenticated: false
}

const renderLoginForm = (props = {}) => {
  return render(
    <BrowserRouter>
      <LoginForm {...props} />
    </BrowserRouter>
  )
}

describe('LoginForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useAuth).mockReturnValue(mockAuthHook)
  })

  it('renders login form with all required fields', () => {
    renderLoginForm()
    
    expect(screen.getByText('Welcome Back')).toBeInTheDocument()
    expect(screen.getByText('Sign in to your account to continue')).toBeInTheDocument()
    expect(screen.getByLabelText('Email Address')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('displays validation errors for empty fields', async () => {
    const user = userEvent.setup()
    renderLoginForm()
    
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Email is required')).toBeInTheDocument()
      expect(screen.getByText('Password is required')).toBeInTheDocument()
    })
  })

  it('displays validation error for invalid email format', async () => {
    const user = userEvent.setup()
    renderLoginForm()
    
    const emailInput = screen.getByLabelText('Email Address')
    await user.type(emailInput, 'invalid-email')
    
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Please enter a valid email address')).toBeInTheDocument()
    })
  })

  it('displays validation error for short password', async () => {
    const user = userEvent.setup()
    renderLoginForm()
    
    const passwordInput = screen.getByLabelText('Password')
    await user.type(passwordInput, '123')
    
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Password must be at least 6 characters')).toBeInTheDocument()
    })
  })

  it('toggles password visibility', async () => {
    const user = userEvent.setup()
    renderLoginForm()
    
    const passwordInput = screen.getByLabelText('Password')
    const toggleButton = screen.getByRole('button', { name: '' })
    
    // Initially password should be hidden
    expect(passwordInput).toHaveAttribute('type', 'password')
    
    // Click toggle button to show password
    await user.click(toggleButton)
    expect(passwordInput).toHaveAttribute('type', 'text')
    
    // Click toggle button again to hide password
    await user.click(toggleButton)
    expect(passwordInput).toHaveAttribute('type', 'password')
  })

  it('submits form with valid credentials', async () => {
    const user = userEvent.setup()
    const onSuccess = vi.fn()
    renderLoginForm({ onSuccess })
    
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123'
      })
    })
  })

  it('calls onSuccess callback after successful login', async () => {
    const user = userEvent.setup()
    const onSuccess = vi.fn()
    renderLoginForm({ onSuccess })
    
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled()
    })
  })

  it('displays loading state during login', () => {
    vi.mocked(useAuth).mockReturnValue({
      ...mockAuthHook,
      isLoading: true
    })
    
    renderLoginForm()
    
    expect(screen.getByText('Signing in...')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /signing in/i })).toBeDisabled()
  })

  it('displays error message from auth hook', () => {
    vi.mocked(useAuth).mockReturnValue({
      ...mockAuthHook,
      error: 'Invalid credentials'
    })
    
    renderLoginForm()
    
    expect(screen.getByText('Invalid credentials')).toBeInTheDocument()
  })

  it('handles 401 error (invalid credentials)', async () => {
    const user = userEvent.setup()
    const mockError = {
      response: { status: 401 }
    }
    mockLogin.mockRejectedValue(mockError)
    
    renderLoginForm()
    
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'wrongpassword')
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Invalid email or password')).toBeInTheDocument()
    })
  })

  it('handles 429 error (rate limiting)', async () => {
    const user = userEvent.setup()
    const mockError = {
      response: { status: 429 }
    }
    mockLogin.mockRejectedValue(mockError)
    
    renderLoginForm()
    
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Too many login attempts. Please try again later.')).toBeInTheDocument()
    })
  })

  it('handles generic errors', async () => {
    const user = userEvent.setup()
    const mockError = {
      message: 'Network error'
    }
    mockLogin.mockRejectedValue(mockError)
    
    renderLoginForm()
    
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeInTheDocument()
    })
  })

  it('renders navigation links', () => {
    renderLoginForm()
    
    expect(screen.getByText('Forgot password?')).toBeInTheDocument()
    expect(screen.getByText('Sign up here')).toBeInTheDocument()
  })

  it('handles remember me checkbox', async () => {
    const user = userEvent.setup()
    renderLoginForm()
    
    const rememberMeCheckbox = screen.getByLabelText('Remember me')
    
    expect(rememberMeCheckbox).not.toBeChecked()
    
    await user.click(rememberMeCheckbox)
    expect(rememberMeCheckbox).toBeChecked()
  })

  it('applies custom className', () => {
    const { container } = renderLoginForm({ className: 'custom-class' })
    expect(container.firstChild).toHaveClass('custom-class')
  })

  it('validates form on change', async () => {
    const user = userEvent.setup()
    renderLoginForm()
    
    const emailInput = screen.getByLabelText('Email Address')
    
    // Type invalid email
    await user.type(emailInput, 'invalid')
    await user.tab() // Trigger blur event
    
    await waitFor(() => {
      expect(screen.getByText('Please enter a valid email address')).toBeInTheDocument()
    })
    
    // Fix the email
    await user.clear(emailInput)
    await user.type(emailInput, 'test@example.com')
    
    await waitFor(() => {
      expect(screen.queryByText('Please enter a valid email address')).not.toBeInTheDocument()
    })
  })
})