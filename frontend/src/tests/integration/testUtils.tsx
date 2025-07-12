import React from 'react'
import { render, RenderOptions } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from '@/context/AuthContext'
import { vi } from 'vitest'

// Test wrapper component that provides all necessary context
interface TestWrapperProps {
  children: React.ReactNode
}

const TestWrapper: React.FC<TestWrapperProps> = ({ children }) => {
  return (
    <BrowserRouter>
      <AuthProvider>
        {children}
      </AuthProvider>
    </BrowserRouter>
  )
}

// Custom render function that includes providers
const customRender = (
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) => render(ui, { wrapper: TestWrapper, ...options })

// Mock API responses for common scenarios
export const mockApiResponses = {
  // Auth responses
  loginSuccess: {
    data: {
      user: {
        id: 1,
        email: 'test@example.com',
        name: 'Test User',
        created_at: '2024-01-01T00:00:00Z',
      },
      token: 'mock-jwt-token',
    },
  },
  
  registrationSuccess: {
    data: {
      user: {
        id: 1,
        email: 'test@example.com',
        name: 'Test User',
        created_at: '2024-01-01T00:00:00Z',
      },
      token: 'mock-jwt-token',
    },
  },

  // URL responses
  urlCreationSuccess: {
    data: {
      id: 1,
      short_code: 'abc123',
      original_url: 'https://example.com',
      custom_alias: false,
      expires_at: null,
      is_active: true,
      click_count: 0,
      created_at: '2024-01-01T00:00:00Z',
      user_id: 1,
    },
  },

  urlListSuccess: {
    data: {
      urls: [
        {
          id: 1,
          short_code: 'abc123',
          original_url: 'https://example.com',
          custom_alias: false,
          expires_at: null,
          is_active: true,
          click_count: 5,
          created_at: '2024-01-01T00:00:00Z',
          user_id: 1,
        },
        {
          id: 2,
          short_code: 'def456',
          original_url: 'https://google.com',
          custom_alias: true,
          expires_at: null,
          is_active: true,
          click_count: 10,
          created_at: '2024-01-02T00:00:00Z',
          user_id: 1,
        },
      ],
      total: 2,
      page: 1,
      per_page: 10,
    },
  },

  // Analytics responses
  analyticsSuccess: {
    data: {
      total_clicks: 15,
      unique_clicks: 12,
      click_data: [
        { date: '2024-01-01', clicks: 5 },
        { date: '2024-01-02', clicks: 10 },
      ],
      geographic_data: [
        { country: 'US', clicks: 10 },
        { country: 'CA', clicks: 5 },
      ],
      device_data: [
        { device: 'Desktop', clicks: 12 },
        { device: 'Mobile', clicks: 3 },
      ],
      referrer_data: [
        { referrer: 'direct', clicks: 8 },
        { referrer: 'google.com', clicks: 7 },
      ],
    },
  },

  // Error responses
  unauthorizedError: {
    response: {
      status: 401,
      data: { error: 'Unauthorized' },
    },
  },

  validationError: {
    response: {
      status: 400,
      data: { 
        error: 'Validation failed',
        details: ['Email is required', 'Password is too short'],
      },
    },
  },

  serverError: {
    response: {
      status: 500,
      data: { error: 'Internal server error' },
    },
  },
}

// Helper function to setup mock API calls - each test should create its own mock
export const setupMockApi = () => {
  // This will be implemented per test file
  vi.clearAllMocks()
}

// Helper function to simulate user authentication
export const mockAuthenticatedUser = () => {
  localStorage.setItem('auth_token', 'mock-jwt-token')
  localStorage.setItem('auth_user', JSON.stringify({
    id: 1,
    email: 'test@example.com',
    name: 'Test User',
  }))
}

// Helper function to clear authentication
export const clearMockAuth = () => {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('auth_user')
}

// Helper function to wait for async operations
export const waitForAsyncOperations = () => 
  new Promise(resolve => setTimeout(resolve, 0))

// Re-export everything from testing library with custom render
export * from '@testing-library/react'
export { customRender as render }