import { describe, it, expect, beforeEach, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { render, setupMockApi, mockApiResponses, mockAuthenticatedUser, clearMockAuth } from './testUtils'
import App from '@/App'

describe('Analytics Workflows', () => {
  let mockAxios: any
  let user: any

  beforeEach(() => {
    mockAxios = setupMockApi()
    user = userEvent.setup()
    clearMockAuth()
    
    // Setup authenticated user for analytics tests
    mockAuthenticatedUser()
  })

  describe('Analytics Dashboard Workflow', () => {
    it('displays comprehensive analytics dashboard', async () => {
      // Mock analytics API calls
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        if (url.includes('/api/urls')) {
          return Promise.resolve(mockApiResponses.urlListSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      // Navigate to analytics page
      window.history.pushState({}, '', '/analytics')

      // Wait for analytics dashboard to load
      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Verify key metrics are displayed
      await waitFor(() => {
        expect(screen.getByText(/total clicks/i)).toBeInTheDocument()
        expect(screen.getByText(/15/)).toBeInTheDocument() // total clicks
        expect(screen.getByText(/unique clicks/i)).toBeInTheDocument()
        expect(screen.getByText(/12/)).toBeInTheDocument() // unique clicks
      })

      // Verify API calls were made
      expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/overview')
    })

    it('displays click timeline chart', async () => {
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/click timeline/i)).toBeInTheDocument()
      })

      // Verify chart data points are rendered
      await waitFor(() => {
        expect(screen.getByText(/2024-01-01/)).toBeInTheDocument()
        expect(screen.getByText(/2024-01-02/)).toBeInTheDocument()
      })
    })

    it('displays geographic distribution', async () => {
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/geographic distribution/i)).toBeInTheDocument()
      })

      // Verify country data is displayed
      await waitFor(() => {
        expect(screen.getByText(/US/)).toBeInTheDocument()
        expect(screen.getByText(/CA/)).toBeInTheDocument()
      })
    })

    it('displays device and browser statistics', async () => {
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/device statistics/i)).toBeInTheDocument()
      })

      // Verify device data is displayed
      await waitFor(() => {
        expect(screen.getByText(/Desktop/)).toBeInTheDocument()
        expect(screen.getByText(/Mobile/)).toBeInTheDocument()
      })
    })

    it('displays referrer statistics', async () => {
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/referrer sources/i)).toBeInTheDocument()
      })

      // Verify referrer data is displayed
      await waitFor(() => {
        expect(screen.getByText(/direct/)).toBeInTheDocument()
        expect(screen.getByText(/google\.com/)).toBeInTheDocument()
      })
    })
  })

  describe('Individual URL Analytics Workflow', () => {
    it('displays analytics for specific URL', async () => {
      // Mock URL-specific analytics
      const urlAnalyticsResponse = {
        data: {
          url: {
            id: 1,
            short_code: 'abc123',
            original_url: 'https://example.com',
            click_count: 25,
          },
          analytics: {
            total_clicks: 25,
            unique_clicks: 20,
            click_data: [
              { date: '2024-01-01', clicks: 10 },
              { date: '2024-01-02', clicks: 15 },
            ],
            geographic_data: [
              { country: 'US', clicks: 15 },
              { country: 'CA', clicks: 10 },
            ],
            device_data: [
              { device: 'Desktop', clicks: 20 },
              { device: 'Mobile', clicks: 5 },
            ],
          },
        },
      }

      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/abc123')) {
          return Promise.resolve(urlAnalyticsResponse)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      // Navigate to specific URL analytics
      window.history.pushState({}, '', '/analytics/abc123')

      // Wait for URL analytics to load
      await waitFor(() => {
        expect(screen.getByText(/analytics for abc123/i)).toBeInTheDocument()
      })

      // Verify URL details are displayed
      await waitFor(() => {
        expect(screen.getByText(/https:\/\/example\.com/)).toBeInTheDocument()
        expect(screen.getByText(/25 total clicks/i)).toBeInTheDocument()
        expect(screen.getByText(/20 unique clicks/i)).toBeInTheDocument()
      })

      // Verify API call was made
      expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/abc123')
    })

    it('handles analytics loading errors gracefully', async () => {
      // Mock analytics error
      mockAxios.get.mockRejectedValueOnce(mockApiResponses.serverError)

      render(<App />)

      window.history.pushState({}, '', '/analytics/abc123')

      // Verify error message is displayed
      await waitFor(() => {
        expect(screen.getByText(/failed to load analytics/i)).toBeInTheDocument()
      })
    })
  })

  describe('Analytics Time Range Filtering', () => {
    it('filters analytics by time range', async () => {
      const timeRangeAnalyticsResponse = {
        data: {
          ...mockApiResponses.analyticsSuccess.data,
          time_range: '7d',
          click_data: [
            { date: '2024-01-08', clicks: 8 },
            { date: '2024-01-09', clicks: 12 },
          ],
        },
      }

      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          if (url.includes('period=7d')) {
            return Promise.resolve(timeRangeAnalyticsResponse)
          }
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Change time range
      const timeRangeSelect = screen.getByRole('combobox', { name: /time range/i })
      await user.selectOptions(timeRangeSelect, '7d')

      // Verify new API call with time range parameter
      await waitFor(() => {
        expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/overview', {
          params: { period: '7d' }
        })
      })

      // Verify updated data is displayed
      await waitFor(() => {
        expect(screen.getByText(/2024-01-08/)).toBeInTheDocument()
        expect(screen.getByText(/2024-01-09/)).toBeInTheDocument()
      })
    })

    it('supports custom date range selection', async () => {
      const customRangeResponse = {
        data: {
          ...mockApiResponses.analyticsSuccess.data,
          start_date: '2024-01-01',
          end_date: '2024-01-31',
        },
      }

      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          if (url.includes('start_date=2024-01-01') && url.includes('end_date=2024-01-31')) {
            return Promise.resolve(customRangeResponse)
          }
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Select custom date range
      const timeRangeSelect = screen.getByRole('combobox', { name: /time range/i })
      await user.selectOptions(timeRangeSelect, 'custom')

      // Fill in custom date inputs
      const startDateInput = screen.getByLabelText(/start date/i)
      const endDateInput = screen.getByLabelText(/end date/i)

      await user.type(startDateInput, '2024-01-01')
      await user.type(endDateInput, '2024-01-31')

      const applyButton = screen.getByRole('button', { name: /apply/i })
      await user.click(applyButton)

      // Verify API call with custom date range
      await waitFor(() => {
        expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/overview', {
          params: {
            start_date: '2024-01-01',
            end_date: '2024-01-31'
          }
        })
      })
    })
  })

  describe('Analytics Export Workflow', () => {
    it('exports analytics data as CSV', async () => {
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        if (url.includes('/api/analytics/export')) {
          return Promise.resolve({
            data: 'Date,Clicks\n2024-01-01,5\n2024-01-02,10',
            headers: {
              'content-type': 'text/csv',
              'content-disposition': 'attachment; filename="analytics.csv"'
            }
          })
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      // Mock download functionality
      const mockCreateObjectURL = vi.fn(() => 'blob:mock-url')
      const mockRevokeObjectURL = vi.fn()
      global.URL.createObjectURL = mockCreateObjectURL
      global.URL.revokeObjectURL = mockRevokeObjectURL

      // Mock link click for download
      const mockClick = vi.fn()
      const mockLink = {
        href: '',
        download: '',
        click: mockClick,
      }
      vi.spyOn(document, 'createElement').mockReturnValue(mockLink as any)

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Click export button
      const exportButton = screen.getByRole('button', { name: /export/i })
      await user.click(exportButton)

      // Select CSV format
      const csvOption = screen.getByRole('button', { name: /csv/i })
      await user.click(csvOption)

      // Verify export API call
      await waitFor(() => {
        expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/export', {
          params: { format: 'csv' }
        })
      })

      // Verify download was triggered
      expect(mockCreateObjectURL).toHaveBeenCalled()
      expect(mockClick).toHaveBeenCalled()
    })

    it('exports analytics data as PDF report', async () => {
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          return Promise.resolve(mockApiResponses.analyticsSuccess)
        }
        if (url.includes('/api/analytics/export')) {
          return Promise.resolve({
            data: new ArrayBuffer(1024), // Mock PDF data
            headers: {
              'content-type': 'application/pdf',
              'content-disposition': 'attachment; filename="analytics-report.pdf"'
            }
          })
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      const mockCreateObjectURL = vi.fn(() => 'blob:mock-url')
      global.URL.createObjectURL = mockCreateObjectURL

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Export as PDF
      const exportButton = screen.getByRole('button', { name: /export/i })
      await user.click(exportButton)

      const pdfOption = screen.getByRole('button', { name: /pdf/i })
      await user.click(pdfOption)

      // Verify export API call
      await waitFor(() => {
        expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/export', {
          params: { format: 'pdf' }
        })
      })
    })
  })

  describe('Real-time Analytics Updates', () => {
    it('updates analytics data in real-time', async () => {
      let resolveInitialCall: any
      let resolveUpdateCall: any

      // Setup promises for controlling API call timing
      const initialCallPromise = new Promise(resolve => {
        resolveInitialCall = resolve
      })
      const updateCallPromise = new Promise(resolve => {
        resolveUpdateCall = resolve
      })

      let callCount = 0
      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/overview')) {
          callCount++
          if (callCount === 1) {
            return initialCallPromise
          } else {
            return updateCallPromise
          }
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      // Resolve initial call
      resolveInitialCall(mockApiResponses.analyticsSuccess)

      await waitFor(() => {
        expect(screen.getByText(/total clicks/i)).toBeInTheDocument()
        expect(screen.getByText(/15/)).toBeInTheDocument()
      })

      // Simulate real-time update with new data
      const updatedAnalyticsResponse = {
        data: {
          ...mockApiResponses.analyticsSuccess.data,
          total_clicks: 18,
          unique_clicks: 15,
        },
      }

      // Resolve update call
      resolveUpdateCall(updatedAnalyticsResponse)

      // Verify updated data appears
      await waitFor(() => {
        expect(screen.getByText(/18/)).toBeInTheDocument() // updated total clicks
        expect(screen.getByText(/15/)).toBeInTheDocument() // updated unique clicks
      })
    })
  })

  describe('Analytics Comparison Workflow', () => {
    it('compares analytics across different time periods', async () => {
      const comparisonResponse = {
        data: {
          current_period: {
            total_clicks: 25,
            unique_clicks: 20,
            start_date: '2024-01-08',
            end_date: '2024-01-14',
          },
          previous_period: {
            total_clicks: 15,
            unique_clicks: 12,
            start_date: '2024-01-01',
            end_date: '2024-01-07',
          },
          comparison: {
            clicks_change: 10,
            clicks_change_percent: 66.67,
            unique_clicks_change: 8,
            unique_clicks_change_percent: 66.67,
          },
        },
      }

      mockAxios.get.mockImplementation((url) => {
        if (url.includes('/api/analytics/comparison')) {
          return Promise.resolve(comparisonResponse)
        }
        return Promise.reject(new Error('Unknown endpoint'))
      })

      render(<App />)

      window.history.pushState({}, '', '/analytics')

      await waitFor(() => {
        expect(screen.getByText(/analytics dashboard/i)).toBeInTheDocument()
      })

      // Enable comparison mode
      const comparisonToggle = screen.getByRole('checkbox', { name: /compare periods/i })
      await user.click(comparisonToggle)

      // Verify comparison API call
      await waitFor(() => {
        expect(mockAxios.get).toHaveBeenCalledWith('/api/analytics/comparison', {
          params: expect.objectContaining({
            current_period: expect.any(String),
            previous_period: expect.any(String),
          })
        })
      })

      // Verify comparison data is displayed
      await waitFor(() => {
        expect(screen.getByText(/66\.67% increase/i)).toBeInTheDocument()
        expect(screen.getByText(/\+10 clicks/i)).toBeInTheDocument()
      })
    })
  })
})