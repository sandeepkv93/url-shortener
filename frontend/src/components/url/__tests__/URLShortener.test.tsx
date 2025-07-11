import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import URLShortener from '../URLShortener'
import { urlService } from '@/services/urls'
import { URL as URLType } from '@/types/url'

// Mock the services
vi.mock('@/services/urls', () => ({
  urlService: {
    createURL: vi.fn()
  }
}))

// Mock clipboard API
Object.assign(navigator, {
  clipboard: {
    writeText: vi.fn().mockImplementation(() => Promise.resolve())
  }
})

// Mock window.open
global.open = vi.fn()

const mockURL: URLType = {
  id: '1',
  shortCode: 'abc123',
  originalUrl: 'https://example.com',
  title: 'Test URL',
  description: 'Test description',
  userId: 'user1',
  clickCount: 0,
  isActive: true,
  isPublic: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  tags: ['test']
}

const renderComponent = (props = {}) => {
  const defaultProps = {
    onURLCreated: vi.fn(),
    ...props
  }
  
  return render(
    <BrowserRouter>
      <URLShortener {...defaultProps} />
    </BrowserRouter>
  )
}

describe('URLShortener', () => {
  const user = userEvent.setup()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('Initial Render', () => {
    it('renders the component with form fields', () => {
      renderComponent()
      
      expect(screen.getByText('Shorten Your URL')).toBeInTheDocument()
      expect(screen.getByLabelText('Original URL *')).toBeInTheDocument()
      expect(screen.getByLabelText('Custom Alias (Optional)')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /shorten url/i })).toBeInTheDocument()
    })

    it('shows advanced options when toggled', async () => {
      renderComponent()
      
      expect(screen.queryByLabelText('Title (Optional)')).not.toBeInTheDocument()
      
      await user.click(screen.getByText('Show Advanced Options'))
      
      expect(screen.getByLabelText('Title (Optional)')).toBeInTheDocument()
      expect(screen.getByLabelText('Description (Optional)')).toBeInTheDocument()
      expect(screen.getByLabelText('Expiration Date (Optional)')).toBeInTheDocument()
      expect(screen.getByLabelText('Password Protection (Optional)')).toBeInTheDocument()
    })

    it('hides advanced options when toggled back', async () => {
      renderComponent()
      
      await user.click(screen.getByText('Show Advanced Options'))
      expect(screen.getByLabelText('Title (Optional)')).toBeInTheDocument()
      
      await user.click(screen.getByText('Hide Advanced Options'))
      expect(screen.queryByLabelText('Title (Optional)')).not.toBeInTheDocument()
    })
  })

  describe('Form Validation', () => {
    it('shows error for empty URL', async () => {
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('URL is required')).toBeInTheDocument()
      })
    })

    it('shows error for invalid URL format', async () => {
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'not-a-valid-url')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('Please enter a valid URL')).toBeInTheDocument()
      })
    })

    it('shows error for invalid custom alias', async () => {
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.type(screen.getByLabelText('Custom Alias (Optional)'), 'invalid alias!')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('Custom alias can only contain letters, numbers, hyphens, and underscores')).toBeInTheDocument()
      })
    })

    it('shows error for short password when provided', async () => {
      renderComponent()
      
      await user.click(screen.getByText('Show Advanced Options'))
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.type(screen.getByLabelText('Password Protection (Optional)'), '123')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('Password must be at least 4 characters')).toBeInTheDocument()
      })
    })

    it('shows error for past expiration date', async () => {
      renderComponent()
      
      await user.click(screen.getByText('Show Advanced Options'))
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      
      // Set a past date
      const pastDate = '2020-01-01T00:00'
      await user.type(screen.getByLabelText('Expiration Date (Optional)'), pastDate)
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('Expiration date must be in the future')).toBeInTheDocument()
      })
    })
  })

  describe('Form Submission', () => {
    it('creates URL with basic information', async () => {
      const mockCreateURL = vi.mocked(urlService.createURL)
      mockCreateURL.mockResolvedValue(mockURL)
      const onURLCreated = vi.fn()
      
      renderComponent({ onURLCreated })
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(mockCreateURL).toHaveBeenCalledWith({
          originalUrl: 'https://example.com',
          customAlias: undefined,
          title: undefined,
          description: undefined,
          expiresAt: undefined,
          password: undefined,
          isPublic: true,
          tags: []
        })
      })
      
      expect(onURLCreated).toHaveBeenCalledWith(mockURL)
    })

    it('creates URL with all advanced options', async () => {
      const mockCreateURL = vi.mocked(urlService.createURL)
      mockCreateURL.mockResolvedValue(mockURL)
      
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.type(screen.getByLabelText('Custom Alias (Optional)'), 'my-alias')
      
      await user.click(screen.getByText('Show Advanced Options'))
      await user.type(screen.getByLabelText('Title (Optional)'), 'My Title')
      await user.type(screen.getByLabelText('Description (Optional)'), 'My Description')
      await user.type(screen.getByLabelText('Password Protection (Optional)'), 'password123')
      
      // Add a tag
      await user.type(screen.getByPlaceholderText('Add a tag'), 'test-tag')
      await user.click(screen.getByText('Add'))
      
      // Set public to false
      await user.click(screen.getByLabelText('Make this URL public'))
      
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(mockCreateURL).toHaveBeenCalledWith({
          originalUrl: 'https://example.com',
          customAlias: 'my-alias',
          title: 'My Title',
          description: 'My Description',
          expiresAt: undefined,
          password: 'password123',
          isPublic: false,
          tags: ['test-tag']
        })
      })
    })

    it('shows loading state during submission', async () => {
      const mockCreateURL = vi.mocked(urlService.createURL)
      mockCreateURL.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)))
      
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      expect(screen.getByText('Creating Short URL...')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /creating short url/i })).toBeDisabled()
    })

    it('handles submission errors', async () => {
      const mockCreateURL = vi.mocked(urlService.createURL)
      mockCreateURL.mockRejectedValue(new Error('Failed to create URL'))
      
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('Failed to create URL')).toBeInTheDocument()
      })
    })

    it('handles custom alias already taken error', async () => {
      const mockCreateURL = vi.mocked(urlService.createURL)
      const error = {
        response: {
          status: 409,
          data: { message: 'Custom alias already taken' }
        }
      }
      mockCreateURL.mockRejectedValue(error)
      
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.type(screen.getByLabelText('Custom Alias (Optional)'), 'taken-alias')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('This custom alias is already taken')).toBeInTheDocument()
      })
    })
  })

  describe('Tags Management', () => {
    beforeEach(async () => {
      renderComponent()
      await user.click(screen.getByText('Show Advanced Options'))
    })

    it('adds tags when Add button is clicked', async () => {
      await user.type(screen.getByPlaceholderText('Add a tag'), 'tag1')
      await user.click(screen.getByText('Add'))
      
      expect(screen.getByText('tag1')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('Add a tag')).toHaveValue('')
    })

    it('adds tags when Enter key is pressed', async () => {
      const tagInput = screen.getByPlaceholderText('Add a tag')
      await user.type(tagInput, 'tag1')
      await user.keyboard('{Enter}')
      
      expect(screen.getByText('tag1')).toBeInTheDocument()
    })

    it('prevents duplicate tags', async () => {
      await user.type(screen.getByPlaceholderText('Add a tag'), 'tag1')
      await user.click(screen.getByText('Add'))
      
      await user.type(screen.getByPlaceholderText('Add a tag'), 'tag1')
      await user.click(screen.getByText('Add'))
      
      expect(screen.getAllByText('tag1')).toHaveLength(1)
    })

    it('removes tags when X button is clicked', async () => {
      await user.type(screen.getByPlaceholderText('Add a tag'), 'tag1')
      await user.click(screen.getByText('Add'))
      
      expect(screen.getByText('tag1')).toBeInTheDocument()
      
      await user.click(screen.getByRole('button', { name: '' })) // X button
      
      expect(screen.queryByText('tag1')).not.toBeInTheDocument()
    })
  })

  describe('Success State', () => {
    beforeEach(async () => {
      const mockCreateURL = vi.mocked(urlService.createURL)
      mockCreateURL.mockResolvedValue(mockURL)
      
      renderComponent()
      
      await user.type(screen.getByLabelText('Original URL *'), 'https://example.com')
      await user.click(screen.getByRole('button', { name: /shorten url/i }))
      
      await waitFor(() => {
        expect(screen.getByText('URL Shortened Successfully!')).toBeInTheDocument()
      })
    })

    it('shows success message and URL details', () => {
      expect(screen.getByText('URL Shortened Successfully!')).toBeInTheDocument()
      expect(screen.getByText('https://example.com')).toBeInTheDocument()
      expect(screen.getByText(`${window.location.origin}/abc123`)).toBeInTheDocument()
    })

    it('copies short URL to clipboard', async () => {
      const copyButton = screen.getByTitle('Copy to clipboard')
      await user.click(copyButton)
      
      expect(navigator.clipboard.writeText).toHaveBeenCalledWith(`${window.location.origin}/abc123`)
      expect(screen.getByText('✓ Copied to clipboard!')).toBeInTheDocument()
    })

    it('generates QR code', async () => {
      const qrButton = screen.getByTitle('Generate QR Code')
      await user.click(qrButton)
      
      expect(global.open).toHaveBeenCalledWith(
        expect.stringContaining('qrserver.com'),
        '_blank'
      )
    })

    it('resets form when Create Another URL is clicked', async () => {
      await user.click(screen.getByText('Create Another URL'))
      
      expect(screen.getByText('Shorten Your URL')).toBeInTheDocument()
      expect(screen.queryByText('URL Shortened Successfully!')).not.toBeInTheDocument()
    })

    it('opens analytics in new window', async () => {
      const analyticsButton = screen.getByText('View Analytics')
      await user.click(analyticsButton)
      
      expect(global.open).toHaveBeenCalledWith('/analytics/1', '_blank')
    })
  })

  describe('Password Visibility Toggle', () => {
    beforeEach(async () => {
      renderComponent()
      await user.click(screen.getByText('Show Advanced Options'))
    })

    it('toggles password visibility', async () => {
      const passwordInput = screen.getByLabelText('Password Protection (Optional)')
      
      await user.type(passwordInput, 'password123')
      
      expect(passwordInput).toHaveAttribute('type', 'password')
      
      await user.click(screen.getByRole('button', { name: '' })) // Eye button
      
      expect(passwordInput).toHaveAttribute('type', 'text')
      
      await user.click(screen.getByRole('button', { name: '' })) // Eye off button
      
      expect(passwordInput).toHaveAttribute('type', 'password')
    })

    it('only shows toggle button when password has value', async () => {
      const passwordInput = screen.getByLabelText('Password Protection (Optional)')
      
      expect(screen.queryByRole('button', { name: '' })).not.toBeInTheDocument()
      
      await user.type(passwordInput, 'password123')
      
      expect(screen.getByRole('button', { name: '' })).toBeInTheDocument()
    })
  })

  describe('Accessibility', () => {
    it('has proper form labels and ARIA attributes', () => {
      renderComponent()
      
      expect(screen.getByLabelText('Original URL *')).toBeInTheDocument()
      expect(screen.getByLabelText('Custom Alias (Optional)')).toBeInTheDocument()
      
      const submitButton = screen.getByRole('button', { name: /shorten url/i })
      expect(submitButton).toBeInTheDocument()
    })

    it('maintains focus management', async () => {
      renderComponent()
      
      const urlInput = screen.getByLabelText('Original URL *')
      const submitButton = screen.getByRole('button', { name: /shorten url/i })
      
      urlInput.focus()
      expect(urlInput).toHaveFocus()
      
      await user.tab()
      expect(screen.getByLabelText('Custom Alias (Optional)')).toHaveFocus()
    })
  })
})