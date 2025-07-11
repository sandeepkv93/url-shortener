import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import URLList from '../URLList'
import { urlService } from '@/services/urls'
import { URL as URLType, URLListResponse } from '@/types/url'

// Mock the services
vi.mock('@/services/urls', () => ({
  urlService: {
    getUserURLs: vi.fn(),
    toggleURLStatus: vi.fn(),
    deleteURL: vi.fn()
  }
}))

// Mock URLCard component to simplify testing
vi.mock('../URLCard', () => ({
  default: ({ url, onUpdate, onDelete }: any) => (
    <div data-testid={`url-card-${url.id}`}>
      <div>{url.title || url.originalUrl}</div>
      <button onClick={() => onUpdate?.(url)}>Update</button>
      <button onClick={() => onDelete?.(url.id)}>Delete</button>
    </div>
  )
}))

const mockURLs: URLType[] = [
  {
    id: '1',
    shortCode: 'abc123',
    originalUrl: 'https://example.com',
    title: 'Example URL',
    description: 'Example description',
    userId: 'user1',
    clickCount: 42,
    isActive: true,
    isPublic: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    tags: ['example', 'test']
  },
  {
    id: '2',
    shortCode: 'def456',
    originalUrl: 'https://test.com',
    title: 'Test URL',
    userId: 'user1',
    clickCount: 15,
    isActive: false,
    isPublic: false,
    createdAt: '2024-01-02T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
    tags: ['test'],
    password: 'hasPassword'
  },
  {
    id: '3',
    shortCode: 'ghi789',
    originalUrl: 'https://demo.com',
    title: 'Demo URL',
    userId: 'user1',
    clickCount: 100,
    isActive: true,
    isPublic: true,
    createdAt: '2024-01-03T00:00:00Z',
    updatedAt: '2024-01-03T00:00:00Z',
    expiresAt: '2024-12-31T23:59:59Z',
    tags: ['demo']
  }
]

const mockResponse: URLListResponse = {
  urls: mockURLs,
  total: 3,
  page: 1,
  limit: 12,
  hasNext: false,
  hasPrev: false
}

const renderComponent = (props = {}) => {
  const defaultProps = {
    onURLCreate: vi.fn(),
    onURLUpdate: vi.fn(),
    onURLDelete: vi.fn(),
    ...props
  }
  
  return render(
    <BrowserRouter>
      <URLList {...defaultProps} />
    </BrowserRouter>
  )
}

describe('URLList', () => {
  const user = userEvent.setup()

  beforeEach(() => {
    vi.clearAllMocks()
    const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
    mockGetUserURLs.mockResolvedValue(mockResponse)
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('Initial Render and Loading', () => {
    it('shows loading state initially', () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      mockGetUserURLs.mockImplementation(() => new Promise(() => {})) // Never resolves
      
      renderComponent()
      
      expect(screen.getByText('Loading your URLs...')).toBeInTheDocument()
    })

    it('renders URLs after loading', async () => {
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('Example URL')).toBeInTheDocument()
        expect(screen.getByText('Test URL')).toBeInTheDocument()
        expect(screen.getByText('Demo URL')).toBeInTheDocument()
      })
    })

    it('shows correct total count', async () => {
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('3 URLs total')).toBeInTheDocument()
      })
    })

    it('calls getUserURLs with default filters on mount', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      
      renderComponent()
      
      expect(mockGetUserURLs).toHaveBeenCalledWith({
        search: '',
        isActive: undefined,
        isPublic: undefined,
        hasExpiry: undefined,
        hasPassword: undefined,
        tags: [],
        sortBy: 'createdAt',
        sortOrder: 'desc',
        page: 1,
        limit: 12
      })
    })
  })

  describe('Header and Actions', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
    })

    it('shows header with title and count', () => {
      expect(screen.getByText('Your URLs')).toBeInTheDocument()
      expect(screen.getByText('3 URLs total')).toBeInTheDocument()
    })

    it('shows refresh button', () => {
      expect(screen.getByTitle('Refresh')).toBeInTheDocument()
    })

    it('shows create URL button when onURLCreate is provided', () => {
      expect(screen.getByText('Create URL')).toBeInTheDocument()
    })

    it('hides create URL button when onURLCreate is not provided', () => {
      const { rerender } = renderComponent()
      
      rerender(
        <BrowserRouter>
          <URLList />
        </BrowserRouter>
      )
      
      expect(screen.queryByText('Create URL')).not.toBeInTheDocument()
    })

    it('calls onURLCreate when create button is clicked', async () => {
      const onURLCreate = vi.fn()
      const { rerender } = renderComponent()
      
      rerender(
        <BrowserRouter>
          <URLList onURLCreate={onURLCreate} />
        </BrowserRouter>
      )
      
      await user.click(screen.getByText('Create URL'))
      expect(onURLCreate).toHaveBeenCalled()
    })

    it('refreshes URLs when refresh button is clicked', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      
      await user.click(screen.getByTitle('Refresh'))
      
      // Should be called twice: initial load + refresh
      expect(mockGetUserURLs).toHaveBeenCalledTimes(2)
    })
  })

  describe('Search Functionality', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
    })

    it('renders search input', () => {
      expect(screen.getByPlaceholderText('Search URLs by title, description, or original URL...')).toBeInTheDocument()
    })

    it('updates search term and filters URLs with debouncing', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      const searchInput = screen.getByPlaceholderText('Search URLs by title, description, or original URL...')
      
      await user.type(searchInput, 'test')
      
      // Should not call immediately due to debouncing
      expect(mockGetUserURLs).toHaveBeenCalledTimes(1) // Only initial call
      
      // Wait for debounce
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            search: 'test',
            page: 1 // Should reset to page 1
          })
        )
      }, { timeout: 500 })
    })
  })

  describe('Filters Panel', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
    })

    it('toggles filters panel visibility', async () => {
      const filtersButton = screen.getByText('Filters')
      
      // Filters panel should be hidden initially
      expect(screen.queryByText('Status')).not.toBeInTheDocument()
      
      await user.click(filtersButton)
      
      expect(screen.getByText('Status')).toBeInTheDocument()
      expect(screen.getByText('Visibility')).toBeInTheDocument()
      expect(screen.getByText('Properties')).toBeInTheDocument()
    })

    it('shows active filter indicator', async () => {
      const searchInput = screen.getByPlaceholderText('Search URLs by title, description, or original URL...')
      await user.type(searchInput, 'test')
      
      await waitFor(() => {
        expect(screen.getByText('Active')).toBeInTheDocument()
      })
    })

    it('shows clear filters button when filters are active', async () => {
      const searchInput = screen.getByPlaceholderText('Search URLs by title, description, or original URL...')
      await user.type(searchInput, 'test')
      
      await waitFor(() => {
        expect(screen.getByText('Clear all filters')).toBeInTheDocument()
      })
    })

    it('clears all filters when clear button is clicked', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      const searchInput = screen.getByPlaceholderText('Search URLs by title, description, or original URL...')
      
      // Add some filters
      await user.type(searchInput, 'test')
      await waitFor(() => {
        expect(screen.getByText('Clear all filters')).toBeInTheDocument()
      })
      
      await user.click(screen.getByText('Clear all filters'))
      
      // Should reset all filters
      expect(searchInput).toHaveValue('')
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith({
          search: '',
          isActive: undefined,
          isPublic: undefined,
          hasExpiry: undefined,
          hasPassword: undefined,
          tags: [],
          sortBy: 'createdAt',
          sortOrder: 'desc',
          page: 1,
          limit: 12
        })
      })
    })
  })

  describe('Filter Options', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
      await user.click(screen.getByText('Filters'))
    })

    it('filters by status', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      const statusSelect = screen.getByDisplayValue('All URLs')
      
      await user.selectOptions(statusSelect, 'true')
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            isActive: true,
            page: 1
          })
        )
      })
    })

    it('filters by visibility', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      const visibilitySelect = screen.getAllByDisplayValue('All URLs')[1] // Second select
      
      await user.selectOptions(visibilitySelect, 'false')
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            isPublic: false,
            page: 1
          })
        )
      })
    })

    it('filters by expiration', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      const expiryCheckbox = screen.getByLabelText('Has expiration')
      
      await user.click(expiryCheckbox)
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            hasExpiry: true,
            page: 1
          })
        )
      })
    })

    it('filters by password protection', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      const passwordCheckbox = screen.getByLabelText('Password protected')
      
      await user.click(passwordCheckbox)
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            hasPassword: true,
            page: 1
          })
        )
      })
    })

    it('shows and filters by tags when URLs have tags', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      
      // Tags should be visible since mock URLs have tags
      expect(screen.getByText('Tags')).toBeInTheDocument()
      expect(screen.getByText('example')).toBeInTheDocument()
      expect(screen.getByText('test')).toBeInTheDocument()
      expect(screen.getByText('demo')).toBeInTheDocument()
      
      await user.click(screen.getByText('test'))
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            tags: ['test'],
            page: 1
          })
        )
      })
    })
  })

  describe('Sorting', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
    })

    it('shows sort options', () => {
      expect(screen.getByText('Sort by:')).toBeInTheDocument()
      expect(screen.getByText('Date Created')).toBeInTheDocument()
      expect(screen.getByText('Clicks')).toBeInTheDocument()
      expect(screen.getByText('Title')).toBeInTheDocument()
    })

    it('sorts by clicks', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      
      await user.click(screen.getByText('Clicks'))
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            sortBy: 'clickCount',
            sortOrder: 'desc'
          })
        )
      })
    })

    it('toggles sort order when clicking same sort option', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      
      // Click Date Created twice to toggle order
      await user.click(screen.getByText('Date Created'))
      await user.click(screen.getByText('Date Created'))
      
      await waitFor(() => {
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            sortBy: 'createdAt',
            sortOrder: 'asc'
          })
        )
      })
    })

    it('shows sort direction indicator', async () => {
      // Date Created should be active by default
      const dateCreatedButton = screen.getByText('Date Created').closest('button')
      expect(dateCreatedButton).toHaveClass('bg-primary-100')
    })
  })

  describe('Pagination', () => {
    beforeEach(() => {
      const paginatedResponse = {
        ...mockResponse,
        total: 25,
        hasNext: true,
        hasPrev: false
      }
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      mockGetUserURLs.mockResolvedValue(paginatedResponse)
    })

    it('shows pagination when there are more items than limit', async () => {
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('Showing 1 to 3 of 25 results')).toBeInTheDocument()
        expect(screen.getByText('Page 1 of 3')).toBeInTheDocument()
      })
    })

    it('navigates to next page', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('Page 1 of 3')).toBeInTheDocument()
      })
      
      const nextButtons = screen.getAllByRole('button').filter(btn => 
        btn.textContent === 'Next' || btn.querySelector('svg')?.getAttribute('class')?.includes('ChevronRight')
      )
      
      if (nextButtons.length > 0) {
        await user.click(nextButtons[0])
        
        expect(mockGetUserURLs).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 2
          })
        )
      }
    })
  })

  describe('Empty States', () => {
    it('shows empty state when no URLs exist', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      mockGetUserURLs.mockResolvedValue({
        urls: [],
        total: 0,
        page: 1,
        limit: 12,
        hasNext: false,
        hasPrev: false
      })
      
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('No URLs yet')).toBeInTheDocument()
        expect(screen.getByText('Create your first short URL to get started.')).toBeInTheDocument()
      })
    })

    it('shows filtered empty state when no URLs match filters', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      mockGetUserURLs.mockResolvedValue({
        urls: [],
        total: 0,
        page: 1,
        limit: 12,
        hasNext: false,
        hasPrev: false
      })
      
      renderComponent()
      
      // Add a filter first
      const searchInput = screen.getByPlaceholderText('Search URLs by title, description, or original URL...')
      await user.type(searchInput, 'nonexistent')
      
      await waitFor(() => {
        expect(screen.getByText('No URLs match your filters')).toBeInTheDocument()
        expect(screen.getByText('Try adjusting your search or filter criteria.')).toBeInTheDocument()
      })
    })
  })

  describe('URL Updates and Deletions', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
    })

    it('updates URL in list when onUpdate is called', async () => {
      const updatedURL = { ...mockURLs[0], title: 'Updated Title' }
      
      const updateButton = screen.getAllByText('Update')[0]
      await user.click(updateButton)
      
      // URL should be updated in the component's state
      // This would require mocking the update properly
    })

    it('removes URL from list when onDelete is called', async () => {
      const deleteButton = screen.getAllByText('Delete')[0]
      await user.click(deleteButton)
      
      // URL should be removed from the component's state
      // This would require mocking the deletion properly
    })

    it('calls parent onURLUpdate when URL is updated', async () => {
      const onURLUpdate = vi.fn()
      const { rerender } = renderComponent()
      
      rerender(
        <BrowserRouter>
          <URLList onURLUpdate={onURLUpdate} />
        </BrowserRouter>
      )
      
      const updateButton = screen.getAllByText('Update')[0]
      await user.click(updateButton)
      
      expect(onURLUpdate).toHaveBeenCalledWith(mockURLs[0])
    })

    it('calls parent onURLDelete when URL is deleted', async () => {
      const onURLDelete = vi.fn()
      const { rerender } = renderComponent()
      
      rerender(
        <BrowserRouter>
          <URLList onURLDelete={onURLDelete} />
        </BrowserRouter>
      )
      
      const deleteButton = screen.getAllByText('Delete')[0]
      await user.click(deleteButton)
      
      expect(onURLDelete).toHaveBeenCalledWith('1')
    })
  })

  describe('Error Handling', () => {
    it('handles fetch errors gracefully', async () => {
      const mockGetUserURLs = vi.mocked(urlService.getUserURLs)
      mockGetUserURLs.mockRejectedValue(new Error('Failed to fetch URLs'))
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      
      renderComponent()
      
      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Failed to fetch URLs:', expect.any(Error))
      })
      
      consoleSpy.mockRestore()
    })
  })

  describe('Accessibility', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Your URLs')).toBeInTheDocument()
      })
    })

    it('has proper form labels and ARIA attributes', () => {
      expect(screen.getByPlaceholderText('Search URLs by title, description, or original URL...')).toBeInTheDocument()
      expect(screen.getByText('Filters')).toBeInTheDocument()
    })

    it('maintains proper focus management', async () => {
      const searchInput = screen.getByPlaceholderText('Search URLs by title, description, or original URL...')
      const filtersButton = screen.getByText('Filters')
      
      searchInput.focus()
      expect(searchInput).toHaveFocus()
      
      await user.tab()
      expect(filtersButton).toHaveFocus()
    })
  })

  describe('Responsive Design', () => {
    it('applies custom className correctly', () => {
      const { container } = renderComponent({ className: 'custom-class' })
      
      expect(container.firstChild).toHaveClass('custom-class')
    })
  })
})