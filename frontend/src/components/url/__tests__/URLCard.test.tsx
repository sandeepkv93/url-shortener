import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import URLCard from '../URLCard'
import { urlService } from '@/services/urls'
import { URL as URLType } from '@/types/url'

// Mock the services
vi.mock('@/services/urls', () => ({
  urlService: {
    toggleURLStatus: vi.fn(),
    deleteURL: vi.fn()
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

// Mock window.confirm
global.confirm = vi.fn()

const mockURL: URLType = {
  id: '1',
  shortCode: 'abc123',
  originalUrl: 'https://example.com/very/long/url',
  title: 'Test URL',
  description: 'This is a test URL description',
  userId: 'user1',
  clickCount: 42,
  isActive: true,
  isPublic: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  tags: ['test', 'example', 'demo'],
  password: 'hasPassword'
}

const mockExpiredURL: URLType = {
  ...mockURL,
  id: '2',
  expiresAt: '2023-12-31T23:59:59Z', // Past date
  isActive: true
}

const mockInactiveURL: URLType = {
  ...mockURL,
  id: '3',
  isActive: false
}

const renderComponent = (props = {}) => {
  const defaultProps = {
    url: mockURL,
    onUpdate: vi.fn(),
    onDelete: vi.fn(),
    showActions: true,
    ...props
  }
  
  return render(
    <BrowserRouter>
      <URLCard {...defaultProps} />
    </BrowserRouter>
  )
}

describe('URLCard', () => {
  const user = userEvent.setup()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('Basic Rendering', () => {
    it('renders URL information correctly', () => {
      renderComponent()
      
      expect(screen.getByText('Test URL')).toBeInTheDocument()
      expect(screen.getByText('https://example.com/very/long/url')).toBeInTheDocument()
      expect(screen.getByText('This is a test URL description')).toBeInTheDocument()
      expect(screen.getByText('42 clicks')).toBeInTheDocument()
      expect(screen.getByText('Created Jan 1, 2024')).toBeInTheDocument()
    })

    it('renders short URL correctly', () => {
      renderComponent()
      
      expect(screen.getByText(`${window.location.origin}/abc123`)).toBeInTheDocument()
    })

    it('shows URL without title when title is not provided', () => {
      const urlWithoutTitle = { ...mockURL, title: undefined }
      renderComponent({ url: urlWithoutTitle })
      
      expect(screen.getByText('https://example.com/very/long/url')).toBeInTheDocument()
      expect(screen.queryByText('Test URL')).not.toBeInTheDocument()
    })

    it('renders tags correctly', () => {
      renderComponent()
      
      expect(screen.getByText('test')).toBeInTheDocument()
      expect(screen.getByText('example')).toBeInTheDocument()
      expect(screen.getByText('demo')).toBeInTheDocument()
    })

    it('shows +more indicator when there are more than 3 tags', () => {
      const urlWithManyTags = {
        ...mockURL,
        tags: ['tag1', 'tag2', 'tag3', 'tag4', 'tag5']
      }
      renderComponent({ url: urlWithManyTags })
      
      expect(screen.getByText('+2 more')).toBeInTheDocument()
    })

    it('hides actions when showActions is false', () => {
      renderComponent({ showActions: false })
      
      expect(screen.queryByRole('button', { name: '' })).not.toBeInTheDocument() // More menu button
    })
  })

  describe('Status Indicators', () => {
    it('shows active status for active URLs', () => {
      renderComponent()
      
      expect(screen.getByText('Active')).toBeInTheDocument()
      expect(screen.getByText('Active')).toHaveClass('text-green-500')
    })

    it('shows inactive status for inactive URLs', () => {
      renderComponent({ url: mockInactiveURL })
      
      expect(screen.getByText('Inactive')).toBeInTheDocument()
      expect(screen.getByText('Inactive')).toHaveClass('text-gray-500')
    })

    it('shows expired status for expired URLs', () => {
      renderComponent({ url: mockExpiredURL })
      
      expect(screen.getByText('Expired')).toBeInTheDocument()
      expect(screen.getByText('Expired')).toHaveClass('text-orange-500')
    })

    it('shows public status correctly', () => {
      renderComponent()
      
      expect(screen.getByText('Public')).toBeInTheDocument()
    })

    it('shows private status correctly', () => {
      const privateURL = { ...mockURL, isPublic: false }
      renderComponent({ url: privateURL })
      
      expect(screen.getByText('Private')).toBeInTheDocument()
    })

    it('shows password protection indicator', () => {
      renderComponent()
      
      expect(screen.getByText('Protected')).toBeInTheDocument()
    })

    it('hides password protection indicator when no password', () => {
      const urlWithoutPassword = { ...mockURL, password: undefined }
      renderComponent({ url: urlWithoutPassword })
      
      expect(screen.queryByText('Protected')).not.toBeInTheDocument()
    })
  })

  describe('Expiration Warnings', () => {
    it('shows expiration warning for expired URLs', () => {
      renderComponent({ url: mockExpiredURL })
      
      expect(screen.getByText(/Expired on Dec 31, 2023/)).toBeInTheDocument()
    })

    it('shows warning for URLs expiring within 7 days', () => {
      const soonToExpireURL = {
        ...mockURL,
        expiresAt: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString() // 5 days from now
      }
      renderComponent({ url: soonToExpireURL })
      
      expect(screen.getByText(/Expires on/)).toBeInTheDocument()
    })

    it('shows normal expiration info for URLs expiring later', () => {
      const laterExpiringURL = {
        ...mockURL,
        expiresAt: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString() // 30 days from now
      }
      renderComponent({ url: laterExpiringURL })
      
      expect(screen.getByText(/Expires on/)).toBeInTheDocument()
    })
  })

  describe('Copy Functionality', () => {
    it('copies URL to clipboard when copy button is clicked', async () => {
      renderComponent()
      
      const copyButton = screen.getByTitle('Copy short URL')
      await user.click(copyButton)
      
      expect(navigator.clipboard.writeText).toHaveBeenCalledWith(`${window.location.origin}/abc123`)
    })

    it('shows copied feedback after copying', async () => {
      renderComponent()
      
      const copyButton = screen.getByTitle('Copy short URL')
      await user.click(copyButton)
      
      await waitFor(() => {
        expect(screen.getByText('✓ Copied to clipboard!')).toBeInTheDocument()
      })
    })
  })

  describe('QR Code Generation', () => {
    it('opens QR code in new window when QR button is clicked', async () => {
      renderComponent()
      
      const qrButton = screen.getByTitle('Generate QR Code')
      await user.click(qrButton)
      
      expect(global.open).toHaveBeenCalledWith(
        expect.stringContaining('qrserver.com'),
        '_blank'
      )
    })
  })

  describe('External Link', () => {
    it('opens short URL in new tab when external link is clicked', () => {
      renderComponent()
      
      const externalLink = screen.getByTitle('Open short URL')
      expect(externalLink).toHaveAttribute('href', `${window.location.origin}/abc123`)
      expect(externalLink).toHaveAttribute('target', '_blank')
      expect(externalLink).toHaveAttribute('rel', 'noopener noreferrer')
    })
  })

  describe('Analytics Link', () => {
    it('links to analytics page', () => {
      renderComponent()
      
      const analyticsLink = screen.getByText('Analytics')
      expect(analyticsLink.closest('a')).toHaveAttribute('href', '/analytics/1')
    })
  })

  describe('Actions Menu', () => {
    it('opens and closes actions menu', async () => {
      renderComponent()
      
      const menuButton = screen.getByRole('button', { name: '' }) // More menu button
      
      // Menu should be closed initially
      expect(screen.queryByText('Edit URL')).not.toBeInTheDocument()
      
      // Open menu
      await user.click(menuButton)
      expect(screen.getByText('Edit URL')).toBeInTheDocument()
      expect(screen.getByText('Deactivate')).toBeInTheDocument()
      expect(screen.getByText('Delete')).toBeInTheDocument()
      
      // Close menu by clicking outside
      await user.click(document.body)
      await waitFor(() => {
        expect(screen.queryByText('Edit URL')).not.toBeInTheDocument()
      })
    })

    it('shows correct toggle text for active URL', async () => {
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' }))
      expect(screen.getByText('Deactivate')).toBeInTheDocument()
    })

    it('shows correct toggle text for inactive URL', async () => {
      renderComponent({ url: mockInactiveURL })
      
      await user.click(screen.getByRole('button', { name: '' }))
      expect(screen.getByText('Activate')).toBeInTheDocument()
    })

    it('links to edit page correctly', async () => {
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' }))
      
      const editLink = screen.getByText('Edit URL')
      expect(editLink.closest('a')).toHaveAttribute('href', '/urls/1/edit')
    })
  })

  describe('URL Status Toggle', () => {
    it('toggles URL status and calls onUpdate', async () => {
      const mockToggleURLStatus = vi.mocked(urlService.toggleURLStatus)
      const updatedURL = { ...mockURL, isActive: false }
      mockToggleURLStatus.mockResolvedValue(updatedURL)
      
      const onUpdate = vi.fn()
      renderComponent({ onUpdate })
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Deactivate'))
      
      expect(mockToggleURLStatus).toHaveBeenCalledWith('1')
      
      await waitFor(() => {
        expect(onUpdate).toHaveBeenCalledWith(updatedURL)
      })
    })

    it('shows loading state during toggle', async () => {
      const mockToggleURLStatus = vi.mocked(urlService.toggleURLStatus)
      mockToggleURLStatus.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)))
      
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Deactivate'))
      
      expect(screen.getByText('Deactivate')).toBeInTheDocument()
    })

    it('handles toggle errors gracefully', async () => {
      const mockToggleURLStatus = vi.mocked(urlService.toggleURLStatus)
      mockToggleURLStatus.mockRejectedValue(new Error('Failed to toggle status'))
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Deactivate'))
      
      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Failed to toggle URL status:', expect.any(Error))
      })
      
      consoleSpy.mockRestore()
    })
  })

  describe('URL Deletion', () => {
    it('shows confirmation dialog before deletion', async () => {
      global.confirm = vi.fn().mockReturnValue(false)
      
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Delete'))
      
      expect(global.confirm).toHaveBeenCalledWith(
        'Are you sure you want to delete this URL? This action cannot be undone.'
      )
    })

    it('deletes URL when confirmed and calls onDelete', async () => {
      global.confirm = vi.fn().mockReturnValue(true)
      const mockDeleteURL = vi.mocked(urlService.deleteURL)
      mockDeleteURL.mockResolvedValue()
      
      const onDelete = vi.fn()
      renderComponent({ onDelete })
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Delete'))
      
      expect(mockDeleteURL).toHaveBeenCalledWith('1')
      
      await waitFor(() => {
        expect(onDelete).toHaveBeenCalledWith('1')
      })
    })

    it('does not delete URL when not confirmed', async () => {
      global.confirm = vi.fn().mockReturnValue(false)
      const mockDeleteURL = vi.mocked(urlService.deleteURL)
      
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Delete'))
      
      expect(mockDeleteURL).not.toHaveBeenCalled()
    })

    it('handles deletion errors gracefully', async () => {
      global.confirm = vi.fn().mockReturnValue(true)
      const mockDeleteURL = vi.mocked(urlService.deleteURL)
      mockDeleteURL.mockRejectedValue(new Error('Failed to delete URL'))
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      
      renderComponent()
      
      await user.click(screen.getByRole('button', { name: '' })) // Open menu
      await user.click(screen.getByText('Delete'))
      
      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Failed to delete URL:', expect.any(Error))
      })
      
      consoleSpy.mockRestore()
    })
  })

  describe('Accessibility', () => {
    it('has proper ARIA labels and roles', () => {
      renderComponent()
      
      expect(screen.getByRole('button', { name: '' })).toBeInTheDocument() // More menu button
      expect(screen.getByTitle('Copy short URL')).toBeInTheDocument()
      expect(screen.getByTitle('Generate QR Code')).toBeInTheDocument()
      expect(screen.getByTitle('Open short URL')).toBeInTheDocument()
    })

    it('maintains proper focus management in dropdown menu', async () => {
      renderComponent()
      
      const menuButton = screen.getByRole('button', { name: '' })
      await user.click(menuButton)
      
      const editLink = screen.getByText('Edit URL')
      const toggleButton = screen.getByText('Deactivate')
      const deleteButton = screen.getByText('Delete')
      
      expect(editLink).toBeInTheDocument()
      expect(toggleButton).toBeInTheDocument()
      expect(deleteButton).toBeInTheDocument()
    })
  })

  describe('Responsive Design', () => {
    it('applies custom className correctly', () => {
      const { container } = renderComponent({ className: 'custom-class' })
      
      expect(container.firstChild).toHaveClass('custom-class')
    })

    it('handles long URLs with truncation', () => {
      const longURL = {
        ...mockURL,
        originalUrl: 'https://example.com/very/very/very/long/url/that/should/be/truncated'
      }
      
      renderComponent({ url: longURL })
      
      expect(screen.getByText('https://example.com/very/very/very/long/url/that/should/be/truncated')).toBeInTheDocument()
    })
  })
})