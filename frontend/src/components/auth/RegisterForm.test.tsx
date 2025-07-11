import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import RegisterForm from './RegisterForm'
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

const mockRegister = vi.fn()
const mockAuthHook = {
  register: mockRegister,
  isLoading: false,
  error: null,
  user: null,
  isAuthenticated: false
}

const renderRegisterForm = (props = {}) => {
  return render(
    <BrowserRouter>
      <RegisterForm {...props} />
    </BrowserRouter>
  )
}

describe('RegisterForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useAuth).mockReturnValue(mockAuthHook)
  })

  it('renders register form with all required fields', () => {
    renderRegisterForm()
    
    expect(screen.getByText('Create Account')).toBeInTheDocument()
    expect(screen.getByText('Join us and start shortening URLs today')).toBeInTheDocument()
    expect(screen.getByLabelText('Full Name')).toBeInTheDocument()
    expect(screen.getByLabelText('Email Address')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByLabelText('Confirm Password')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument()
  })

  it('displays validation errors for empty fields', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Name is required')).toBeInTheDocument()
      expect(screen.getByText('Email is required')).toBeInTheDocument()
      expect(screen.getByText('Password is required')).toBeInTheDocument()
      expect(screen.getByText('Please confirm your password')).toBeInTheDocument()
    })
  })

  it('displays validation error for invalid name', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const nameInput = screen.getByLabelText('Full Name')
    await user.type(nameInput, 'A')
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Name must be at least 2 characters')).toBeInTheDocument()
    })
  })

  it('displays validation error for invalid name format', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const nameInput = screen.getByLabelText('Full Name')
    await user.type(nameInput, 'John123')
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Name can only contain letters and spaces')).toBeInTheDocument()
    })
  })

  it('displays validation error for invalid email format', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const emailInput = screen.getByLabelText('Email Address')
    await user.type(emailInput, 'invalid-email')
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Please enter a valid email address')).toBeInTheDocument()
    })
  })

  it('displays validation error for weak password', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const passwordInput = screen.getByLabelText('Password')
    await user.type(passwordInput, 'weak')
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Password must be at least 8 characters')).toBeInTheDocument()
    })
  })

  it('displays validation error for password without complexity', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const passwordInput = screen.getByLabelText('Password')
    await user.type(passwordInput, 'password')
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Password must contain at least one uppercase letter, one lowercase letter, and one number')).toBeInTheDocument()
    })
  })

  it('displays validation error for password mismatch', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const passwordInput = screen.getByLabelText('Password')
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    
    await user.type(passwordInput, 'Password123')
    await user.type(confirmPasswordInput, 'Password456')
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText("Passwords don't match")).toBeInTheDocument()
    })
  })

  it('shows password strength indicator', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const passwordInput = screen.getByLabelText('Password')
    
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
    renderRegisterForm()
    
    const passwordInput = screen.getByLabelText('Password')
    const toggleButtons = screen.getAllByRole('button', { name: '' })
    const passwordToggle = toggleButtons[0]
    
    // Initially password should be hidden
    expect(passwordInput).toHaveAttribute('type', 'password')
    
    // Click toggle button to show password
    await user.click(passwordToggle)
    expect(passwordInput).toHaveAttribute('type', 'text')
    
    // Click toggle button again to hide password
    await user.click(passwordToggle)
    expect(passwordInput).toHaveAttribute('type', 'password')
  })

  it('toggles confirm password visibility', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    const toggleButtons = screen.getAllByRole('button', { name: '' })
    const confirmPasswordToggle = toggleButtons[1]
    
    // Initially confirm password should be hidden
    expect(confirmPasswordInput).toHaveAttribute('type', 'password')
    
    // Click toggle button to show confirm password
    await user.click(confirmPasswordToggle)
    expect(confirmPasswordInput).toHaveAttribute('type', 'text')
    
    // Click toggle button again to hide confirm password
    await user.click(confirmPasswordToggle)
    expect(confirmPasswordInput).toHaveAttribute('type', 'password')
  })

  it('requires terms agreement', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const submitButton = screen.getByRole('button', { name: /create account/i })
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('You must agree to the terms and conditions')).toBeInTheDocument()
    })
  })

  it('submits form with valid data', async () => {
    const user = userEvent.setup()
    const onSuccess = vi.fn()
    renderRegisterForm({ onSuccess })
    
    const nameInput = screen.getByLabelText('Full Name')
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    const termsCheckbox = screen.getByLabelText(/I agree to the/i)
    const submitButton = screen.getByRole('button', { name: /create account/i })
    
    await user.type(nameInput, 'John Doe')
    await user.type(emailInput, 'john@example.com')
    await user.type(passwordInput, 'Password123')
    await user.type(confirmPasswordInput, 'Password123')
    await user.click(termsCheckbox)
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith({
        name: 'John Doe',
        email: 'john@example.com',
        password: 'Password123',
        confirmPassword: 'Password123'
      })
    })
  })

  it('shows success screen after registration', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const nameInput = screen.getByLabelText('Full Name')
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    const termsCheckbox = screen.getByLabelText(/I agree to the/i)
    const submitButton = screen.getByRole('button', { name: /create account/i })
    
    await user.type(nameInput, 'John Doe')
    await user.type(emailInput, 'john@example.com')
    await user.type(passwordInput, 'Password123')
    await user.type(confirmPasswordInput, 'Password123')
    await user.click(termsCheckbox)
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Welcome!')).toBeInTheDocument()
      expect(screen.getByText('Your account has been created successfully')).toBeInTheDocument()
      expect(screen.getByText('Go to Dashboard')).toBeInTheDocument()
    })
  })

  it('displays loading state during registration', () => {
    vi.mocked(useAuth).mockReturnValue({
      ...mockAuthHook,
      isLoading: true
    })
    
    renderRegisterForm()
    
    expect(screen.getByText('Creating Account...')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /creating account/i })).toBeDisabled()
  })

  it('displays error message from auth hook', () => {
    vi.mocked(useAuth).mockReturnValue({
      ...mockAuthHook,
      error: 'Registration failed'
    })
    
    renderRegisterForm()
    
    expect(screen.getByText('Registration failed')).toBeInTheDocument()
  })

  it('handles 409 error (email already exists)', async () => {
    const user = userEvent.setup()
    const mockError = {
      response: { status: 409 }
    }
    mockRegister.mockRejectedValue(mockError)
    
    renderRegisterForm()
    
    const nameInput = screen.getByLabelText('Full Name')
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    const termsCheckbox = screen.getByLabelText(/I agree to the/i)
    const submitButton = screen.getByRole('button', { name: /create account/i })
    
    await user.type(nameInput, 'John Doe')
    await user.type(emailInput, 'existing@example.com')
    await user.type(passwordInput, 'Password123')
    await user.type(confirmPasswordInput, 'Password123')
    await user.click(termsCheckbox)
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('An account with this email already exists')).toBeInTheDocument()
    })
  })

  it('handles 429 error (rate limiting)', async () => {
    const user = userEvent.setup()
    const mockError = {
      response: { status: 429 }
    }
    mockRegister.mockRejectedValue(mockError)
    
    renderRegisterForm()
    
    const nameInput = screen.getByLabelText('Full Name')
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    const termsCheckbox = screen.getByLabelText(/I agree to the/i)
    const submitButton = screen.getByRole('button', { name: /create account/i })
    
    await user.type(nameInput, 'John Doe')
    await user.type(emailInput, 'john@example.com')
    await user.type(passwordInput, 'Password123')
    await user.type(confirmPasswordInput, 'Password123')
    await user.click(termsCheckbox)
    await user.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('Too many registration attempts. Please try again later.')).toBeInTheDocument()
    })
  })

  it('renders navigation links', () => {
    renderRegisterForm()
    
    expect(screen.getByText('Sign in here')).toBeInTheDocument()
    expect(screen.getByText('Terms of Service')).toBeInTheDocument()
    expect(screen.getByText('Privacy Policy')).toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = renderRegisterForm({ className: 'custom-class' })
    expect(container.firstChild).toHaveClass('custom-class')
  })

  it('validates form on change', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
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

  it('shows password strength colors correctly', async () => {
    const user = userEvent.setup()
    renderRegisterForm()
    
    const passwordInput = screen.getByLabelText('Password')
    
    // Type weak password
    await user.type(passwordInput, 'weak')
    expect(screen.getByText('Weak')).toHaveClass('text-red-600')
    
    // Type fair password
    await user.clear(passwordInput)
    await user.type(passwordInput, 'Password')
    expect(screen.getByText('Fair')).toHaveClass('text-yellow-600')
    
    // Type good password
    await user.clear(passwordInput)
    await user.type(passwordInput, 'Password123')
    expect(screen.getByText('Good')).toHaveClass('text-blue-600')
  })

  it('calls onSuccess callback after successful registration', async () => {
    const user = userEvent.setup()
    const onSuccess = vi.fn()
    renderRegisterForm({ onSuccess })
    
    const nameInput = screen.getByLabelText('Full Name')
    const emailInput = screen.getByLabelText('Email Address')
    const passwordInput = screen.getByLabelText('Password')
    const confirmPasswordInput = screen.getByLabelText('Confirm Password')
    const termsCheckbox = screen.getByLabelText(/I agree to the/i)
    const submitButton = screen.getByRole('button', { name: /create account/i })
    
    await user.type(nameInput, 'John Doe')
    await user.type(emailInput, 'john@example.com')
    await user.type(passwordInput, 'Password123')
    await user.type(confirmPasswordInput, 'Password123')
    await user.click(termsCheckbox)
    await user.click(submitButton)
    
    // Wait for success screen to appear
    await waitFor(() => {
      expect(screen.getByText('Welcome!')).toBeInTheDocument()
    })
    
    // Wait for onSuccess to be called (after timeout)
    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled()
    }, { timeout: 3000 })
  })
})