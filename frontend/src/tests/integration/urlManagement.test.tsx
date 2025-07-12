import { describe, it, expect, beforeEach, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { render, setupMockApi, mockApiResponses, mockAuthenticatedUser, clearMockAuth } from './testUtils'
import App from '@/App'

describe('URL Management Workflows', () => {
  let mockAxios: any
  let user: any

  beforeEach(() => {
    mockAxios = setupMockApi()
    user = userEvent.setup()
    clearMockAuth()
    
    // Setup authenticated user for URL management tests
    mockAuthenticatedUser()
  })

  describe('URL Creation Workflow', () => {
    it('creates a new short URL successfully', async () => {
      // Mock successful URL creation
      mockAxios.post.mockResolvedValueOnce(mockApiResponses.urlCreationSuccess)
      // Mock URL list for dashboard
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      // Navigate to dashboard
      window.history.pushState({}, '', '/dashboard')

      // Wait for dashboard to load
      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      // Fill out URL creation form
      const urlInput = screen.getByLabelText(/original url/i)
      const submitButton = screen.getByRole('button', { name: /shorten url/i })

      await user.type(urlInput, 'https://example.com/very-long-url-that-needs-shortening')
      await user.click(submitButton)

      // Verify API call was made
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/urls', {
          original_url: 'https://example.com/very-long-url-that-needs-shortening',
          custom_alias: '',
          expires_at: null,
        })
      })

      // Verify success feedback
      await waitFor(() => {
        expect(screen.getByText(/url shortened successfully/i)).toBeInTheDocument()
      })

      // Verify the new URL appears in the list
      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
        expect(screen.getByText(/https:\/\/example\.com/)).toBeInTheDocument()
      })
    })

    it('creates a short URL with custom alias', async () => {
      // Mock successful URL creation with custom alias
      const customUrlResponse = {
        data: {
          ...mockApiResponses.urlCreationSuccess.data,
          short_code: 'my-custom-link',
          custom_alias: true,
        }
      }
      mockAxios.post.mockResolvedValueOnce(customUrlResponse)
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      // Fill out form with custom alias
      const urlInput = screen.getByLabelText(/original url/i)
      const customAliasCheckbox = screen.getByLabelText(/custom alias/i)
      const aliasInput = screen.getByLabelText(/custom alias/i)
      const submitButton = screen.getByRole('button', { name: /shorten url/i })

      await user.type(urlInput, 'https://example.com/marketing-page')
      await user.click(customAliasCheckbox)
      await user.type(aliasInput, 'my-custom-link')
      await user.click(submitButton)

      // Verify API call with custom alias
      await waitFor(() => {
        expect(mockAxios.post).toHaveBeenCalledWith('/api/urls', {
          original_url: 'https://example.com/marketing-page',
          custom_alias: 'my-custom-link',
          expires_at: null,
        })
      })

      // Verify custom alias is displayed
      await waitFor(() => {
        expect(screen.getByText(/my-custom-link/i)).toBeInTheDocument()
      })
    })

    it('handles URL creation validation errors', async () => {
      // Mock validation error
      mockAxios.post.mockRejectedValueOnce(mockApiResponses.validationError)
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/create short url/i)).toBeInTheDocument()
      })

      // Submit form with invalid URL
      const urlInput = screen.getByLabelText(/original url/i)
      const submitButton = screen.getByRole('button', { name: /shorten url/i })

      await user.type(urlInput, 'invalid-url')
      await user.click(submitButton)

      // Verify validation errors are displayed
      await waitFor(() => {
        expect(screen.getByText(/validation failed/i)).toBeInTheDocument()
      })
    })
  })

  describe('URL List and Display Workflow', () => {
    it('displays user URLs in dashboard', async () => {
      // Mock URL list API call
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Verify URLs are displayed
      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
        expect(screen.getByText(/def456/i)).toBeInTheDocument()
        expect(screen.getByText(/https:\/\/example\.com/)).toBeInTheDocument()
        expect(screen.getByText(/https:\/\/google\.com/)).toBeInTheDocument()
      })

      // Verify click counts are displayed
      expect(screen.getByText(/5 clicks/i)).toBeInTheDocument()
      expect(screen.getByText(/10 clicks/i)).toBeInTheDocument()
    })

    it('allows copying short URLs to clipboard', async () => {
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      // Mock clipboard API
      const mockWriteText = vi.fn()
      Object.assign(navigator, {
        clipboard: {
          writeText: mockWriteText,
        },
      })

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Click copy button for first URL
      const copyButtons = screen.getAllByRole('button', { name: /copy/i })
      await user.click(copyButtons[0])

      // Verify clipboard write was called
      expect(mockWriteText).toHaveBeenCalledWith(expect.stringContaining('abc123'))

      // Verify success feedback
      await waitFor(() => {
        expect(screen.getByText(/copied to clipboard/i)).toBeInTheDocument()
      })
    })

    it('handles URL list loading errors', async () => {
      // Mock server error
      mockAxios.get.mockRejectedValueOnce(mockApiResponses.serverError)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      // Verify error message is displayed
      await waitFor(() => {
        expect(screen.getByText(/failed to load urls/i)).toBeInTheDocument()
      })
    })
  })

  describe('URL Edit and Update Workflow', () => {
    it('updates URL details successfully', async () => {
      // Mock initial URL list
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)
      
      // Mock successful update
      const updatedUrlResponse = {
        data: {
          ...mockApiResponses.urlCreationSuccess.data,
          original_url: 'https://updated-example.com',
        }
      }
      mockAxios.put.mockResolvedValueOnce(updatedUrlResponse)
      
      // Mock updated URL list
      const updatedListResponse = {
        data: {
          ...mockApiResponses.urlListSuccess.data,
          urls: [{
            ...mockApiResponses.urlListSuccess.data.urls[0],
            original_url: 'https://updated-example.com',
          }, ...mockApiResponses.urlListSuccess.data.urls.slice(1)]
        }
      }
      mockAxios.get.mockResolvedValueOnce(updatedListResponse)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Click edit button for first URL
      const editButtons = screen.getAllByRole('button', { name: /edit/i })
      await user.click(editButtons[0])

      // Wait for edit modal/form to appear
      await waitFor(() => {
        expect(screen.getByText(/edit url/i)).toBeInTheDocument()
      })

      // Update the URL
      const urlInput = screen.getByDisplayValue(/https:\/\/example\.com/)
      await user.clear(urlInput)
      await user.type(urlInput, 'https://updated-example.com')

      const saveButton = screen.getByRole('button', { name: /save changes/i })
      await user.click(saveButton)

      // Verify update API call
      await waitFor(() => {
        expect(mockAxios.put).toHaveBeenCalledWith('/api/urls/1', {
          original_url: 'https://updated-example.com',
        })
      })

      // Verify updated URL appears in list
      await waitFor(() => {
        expect(screen.getByText(/https:\/\/updated-example\.com/)).toBeInTheDocument()
      })
    })

    it('toggles URL active status', async () => {
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)
      
      // Mock successful status toggle
      const toggledUrlResponse = {
        data: {
          ...mockApiResponses.urlCreationSuccess.data,
          is_active: false,
        }
      }
      mockAxios.put.mockResolvedValueOnce(toggledUrlResponse)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Click toggle button
      const toggleButtons = screen.getAllByRole('button', { name: /disable/i })
      await user.click(toggleButtons[0])

      // Verify API call
      await waitFor(() => {
        expect(mockAxios.put).toHaveBeenCalledWith('/api/urls/1', {
          is_active: false,
        })
      })

      // Verify status change feedback
      await waitFor(() => {
        expect(screen.getByText(/url disabled/i)).toBeInTheDocument()
      })
    })
  })

  describe('URL Deletion Workflow', () => {
    it('deletes URL successfully', async () => {
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)
      
      // Mock successful deletion
      mockAxios.delete.mockResolvedValueOnce({ data: { message: 'URL deleted successfully' } })
      
      // Mock updated URL list without deleted URL
      const deletedListResponse = {
        data: {
          ...mockApiResponses.urlListSuccess.data,
          urls: mockApiResponses.urlListSuccess.data.urls.slice(1),
          total: 1,
        }
      }
      mockAxios.get.mockResolvedValueOnce(deletedListResponse)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Click delete button
      const deleteButtons = screen.getAllByRole('button', { name: /delete/i })
      await user.click(deleteButtons[0])

      // Confirm deletion in modal
      await waitFor(() => {
        expect(screen.getByText(/are you sure you want to delete/i)).toBeInTheDocument()
      })

      const confirmButton = screen.getByRole('button', { name: /confirm delete/i })
      await user.click(confirmButton)

      // Verify deletion API call
      await waitFor(() => {
        expect(mockAxios.delete).toHaveBeenCalledWith('/api/urls/1')
      })

      // Verify URL is removed from list
      await waitFor(() => {
        expect(screen.queryByText(/abc123/i)).not.toBeInTheDocument()
      })

      // Verify success message
      await waitFor(() => {
        expect(screen.getByText(/url deleted successfully/i)).toBeInTheDocument()
      })
    })
  })

  describe('QR Code Generation Workflow', () => {
    it('generates and displays QR code for URL', async () => {
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Click QR code button
      const qrButtons = screen.getAllByRole('button', { name: /qr code/i })
      await user.click(qrButtons[0])

      // Verify QR code modal appears
      await waitFor(() => {
        expect(screen.getByText(/qr code for/i)).toBeInTheDocument()
      })

      // Verify QR code image is displayed
      expect(screen.getByRole('img', { name: /qr code/i })).toBeInTheDocument()

      // Test download functionality
      const downloadButton = screen.getByRole('button', { name: /download/i })
      expect(downloadButton).toBeInTheDocument()
    })
  })

  describe('URL Search and Filter Workflow', () => {
    it('filters URLs by search term', async () => {
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
        expect(screen.getByText(/def456/i)).toBeInTheDocument()
      })

      // Use search functionality
      const searchInput = screen.getByPlaceholderText(/search urls/i)
      await user.type(searchInput, 'google')

      // Verify filtering works (this would typically trigger a new API call)
      await waitFor(() => {
        expect(screen.getByText(/def456/i)).toBeInTheDocument()
        expect(screen.queryByText(/abc123/i)).not.toBeInTheDocument()
      })
    })

    it('sorts URLs by different criteria', async () => {
      mockAxios.get.mockResolvedValueOnce(mockApiResponses.urlListSuccess)

      render(<App />)

      window.history.pushState({}, '', '/dashboard')

      await waitFor(() => {
        expect(screen.getByText(/abc123/i)).toBeInTheDocument()
      })

      // Change sort order
      const sortSelect = screen.getByRole('combobox', { name: /sort by/i })
      await user.selectOptions(sortSelect, 'clicks')

      // This would typically trigger a new API call with sort parameters
      await waitFor(() => {
        expect(mockAxios.get).toHaveBeenCalledWith('/api/urls', {
          params: expect.objectContaining({
            sort: 'clicks',
            order: 'desc',
          })
        })
      })
    })
  })
})