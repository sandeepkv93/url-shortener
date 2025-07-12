import { renderHook, act, waitFor } from '@testing-library/react'
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest'
import { useRealTimeAnalytics } from './useRealTimeAnalytics'
import { urlService, urlAnalyticsService } from '@/services/urls'

// Mock the URL services
vi.mock('@/services/urls', () => ({
  urlService: {
    getURL: vi.fn(),
    getUserURLs: vi.fn(),
    getClickHistory: vi.fn()
  },
  urlAnalyticsService: {
    getClickTimeline: vi.fn(),
    getGeographicStats: vi.fn(),
    getDeviceStats: vi.fn(),
    getReferrerStats: vi.fn()
  }
}))

// Mock timers
vi.useFakeTimers()

const mockUrlService = vi.mocked(urlService)
const mockUrlAnalyticsService = vi.mocked(urlAnalyticsService)

const mockAnalyticsData = {
  url: { id: 'test-url', shortCode: 'abc123', originalUrl: 'https://example.com' },
  timeline: [
    { date: '2024-01-01', clicks: 100, uniqueClicks: 80 },
    { date: '2024-01-02', clicks: 120, uniqueClicks: 95 }
  ],
  geographic: {
    countries: [
      { country: 'United States', clicks: 150, percentage: 60 }
    ]
  },
  devices: {
    devices: [
      { device: 'Desktop', clicks: 100, percentage: 50 }
    ]
  },
  referrers: {
    referrers: [
      { referrer: 'google.com', clicks: 80, percentage: 40 }
    ]
  },
  clickHistory: {
    clicks: [],
    total: 0,
    page: 1,
    limit: 10
  }
}

describe('useRealTimeAnalytics', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    
    // Setup default mock implementations
    mockUrlService.getURL.mockResolvedValue(mockAnalyticsData.url)
    mockUrlAnalyticsService.getClickTimeline.mockResolvedValue(mockAnalyticsData.timeline)
    mockUrlAnalyticsService.getGeographicStats.mockResolvedValue(mockAnalyticsData.geographic)
    mockUrlAnalyticsService.getDeviceStats.mockResolvedValue(mockAnalyticsData.devices)
    mockUrlAnalyticsService.getReferrerStats.mockResolvedValue(mockAnalyticsData.referrers)
    mockUrlService.getClickHistory.mockResolvedValue(mockAnalyticsData.clickHistory)
    mockUrlService.getUserURLs.mockResolvedValue({ urls: [], total: 0, page: 1, limit: 100 })
  })

  afterEach(() => {
    vi.clearAllTimers()
    vi.clearAllMocks()
  })

  it('initializes with default values', () => {
    const { result } = renderHook(() => useRealTimeAnalytics())
    
    expect(result.current.data).toBeNull()
    expect(result.current.isLoading).toBe(true)
    expect(result.current.isRefreshing).toBe(false)
    expect(result.current.error).toBeNull()
    expect(result.current.connectionStatus).toBe('connected')
    expect(result.current.lastUpdated).toBeUndefined()
  })

  it('fetches analytics data on mount when enabled', async () => {
    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        period: '7d',
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(mockUrlService.getURL).toHaveBeenCalledWith('test-url')
    expect(mockUrlAnalyticsService.getClickTimeline).toHaveBeenCalledWith('test-url', '7d')
    expect(mockUrlAnalyticsService.getGeographicStats).toHaveBeenCalledWith('test-url')
    expect(mockUrlAnalyticsService.getDeviceStats).toHaveBeenCalledWith('test-url')
    expect(mockUrlAnalyticsService.getReferrerStats).toHaveBeenCalledWith('test-url')
    expect(mockUrlService.getClickHistory).toHaveBeenCalledWith('test-url', 1, 10)
    
    expect(result.current.data).toBeTruthy()
    expect(result.current.error).toBeNull()
  })

  it('does not fetch data when disabled', () => {
    renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: false 
      })
    )

    expect(mockUrlService.getURL).not.toHaveBeenCalled()
    expect(mockUrlAnalyticsService.getClickTimeline).not.toHaveBeenCalled()
  })

  it('fetches dashboard data when no urlId provided', async () => {
    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(mockUrlService.getUserURLs).toHaveBeenCalledWith({
      page: 1,
      limit: 100,
      sortBy: 'clickCount',
      sortOrder: 'desc'
    })
    
    expect(result.current.data).toBeTruthy()
  })

  it('handles fetch errors gracefully', async () => {
    const errorMessage = 'Network error'
    mockUrlService.getURL.mockRejectedValue(new Error(errorMessage))

    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.error).toBe(errorMessage)
    expect(result.current.data).toBeNull()
  })

  it('implements retry mechanism with exponential backoff', async () => {
    const errorMessage = 'Network error'
    mockUrlService.getURL.mockRejectedValue(new Error(errorMessage))

    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    // Wait for initial error
    await waitFor(() => {
      expect(result.current.connectionStatus).toBe('reconnecting')
    })

    // Fast-forward time to trigger retry
    act(() => {
      vi.advanceTimersByTime(2000) // First retry after 2 seconds
    })

    // Should have attempted retry
    expect(mockUrlService.getURL).toHaveBeenCalledTimes(2)
  })

  it('stops retrying after max retries reached', async () => {
    const errorMessage = 'Network error'
    mockUrlService.getURL.mockRejectedValue(new Error(errorMessage))

    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    // Wait for initial error and retries
    await waitFor(() => {
      expect(result.current.connectionStatus).toBe('reconnecting')
    })

    // Fast-forward through all retries
    act(() => {
      vi.advanceTimersByTime(30000) // Advance enough for all retries
    })

    await waitFor(() => {
      expect(result.current.connectionStatus).toBe('disconnected')
    })

    // Should have attempted max retries (1 initial + 3 retries = 4 total)
    expect(mockUrlService.getURL).toHaveBeenCalledTimes(4)
  })

  it('refreshes data manually', async () => {
    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Clear previous calls
    vi.clearAllMocks()

    // Manual refresh
    act(() => {
      result.current.refreshData()
    })

    expect(result.current.isRefreshing).toBe(true)

    await waitFor(() => {
      expect(result.current.isRefreshing).toBe(false)
    })

    expect(mockUrlService.getURL).toHaveBeenCalledTimes(1)
  })

  it('sets up automatic refresh interval', async () => {
    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        refreshInterval: 5000, // 5 seconds
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Clear initial fetch calls
    vi.clearAllMocks()

    // Fast-forward past refresh interval
    act(() => {
      vi.advanceTimersByTime(5000)
    })

    await waitFor(() => {
      expect(mockUrlService.getURL).toHaveBeenCalledTimes(1)
    })
  })

  it('stops automatic refresh when disabled', async () => {
    const { result, rerender } = renderHook(
      (props) => useRealTimeAnalytics(props),
      {
        initialProps: {
          urlId: 'test-url',
          refreshInterval: 5000,
          enabled: true
        }
      }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Disable the hook
    rerender({
      urlId: 'test-url',
      refreshInterval: 5000,
      enabled: false
    })

    // Clear calls
    vi.clearAllMocks()

    // Fast-forward past refresh interval
    act(() => {
      vi.advanceTimersByTime(10000)
    })

    // Should not have made any new calls
    expect(mockUrlService.getURL).not.toHaveBeenCalled()
  })

  it('handles visibility change events', async () => {
    Object.defineProperty(document, 'hidden', {
      writable: true,
      value: false
    })

    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        refreshInterval: 5000,
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Clear initial calls
    vi.clearAllMocks()

    // Simulate tab becoming hidden
    Object.defineProperty(document, 'hidden', { value: true })
    act(() => {
      document.dispatchEvent(new Event('visibilitychange'))
    })

    // Fast-forward past refresh interval
    act(() => {
      vi.advanceTimersByTime(10000)
    })

    // Should not refresh while hidden
    expect(mockUrlService.getURL).not.toHaveBeenCalled()

    // Simulate tab becoming visible again
    Object.defineProperty(document, 'hidden', { value: false })
    act(() => {
      document.dispatchEvent(new Event('visibilitychange'))
    })

    // Should refresh when visible again
    await waitFor(() => {
      expect(mockUrlService.getURL).toHaveBeenCalledTimes(1)
    })
  })

  it('cancels ongoing requests when component unmounts', async () => {
    const abortSpy = vi.fn()
    const mockAbortController = {
      abort: abortSpy,
      signal: { aborted: false }
    }
    
    global.AbortController = vi.fn(() => mockAbortController) as any

    const { unmount } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    unmount()

    expect(abortSpy).toHaveBeenCalled()
  })

  it('handles aborted requests gracefully', async () => {
    const abortError = new Error('Request aborted')
    abortError.name = 'AbortError'
    mockUrlService.getURL.mockRejectedValue(abortError)

    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Should not treat abort as error
    expect(result.current.error).toBeNull()
    expect(result.current.connectionStatus).toBe('connected')
  })

  it('transforms analytics data correctly', async () => {
    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data).toEqual({
      url: mockAnalyticsData.url,
      timeline: mockAnalyticsData.timeline,
      geographic: mockAnalyticsData.geographic,
      devices: mockAnalyticsData.devices,
      referrers: mockAnalyticsData.referrers,
      clickHistory: mockAnalyticsData.clickHistory,
      lastUpdated: expect.any(Date)
    })
  })

  it('uses default options when none provided', () => {
    const { result } = renderHook(() => useRealTimeAnalytics())
    
    expect(result.current.data).toBeNull()
    expect(result.current.isLoading).toBe(true)
    expect(result.current.enabled).toBe(true)
  })

  it('updates period correctly', async () => {
    const { result, rerender } = renderHook(
      (props) => useRealTimeAnalytics(props),
      {
        initialProps: {
          urlId: 'test-url',
          period: '7d' as const,
          enabled: true
        }
      }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Clear initial calls
    vi.clearAllMocks()

    // Change period
    rerender({
      urlId: 'test-url',
      period: '30d' as const,
      enabled: true
    })

    await waitFor(() => {
      expect(mockUrlAnalyticsService.getClickTimeline).toHaveBeenCalledWith('test-url', '30d')
    })
  })

  it('provides start and stop functions for real-time updates', () => {
    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        enabled: true 
      })
    )

    expect(typeof result.current.startRealTimeUpdates).toBe('function')
    expect(typeof result.current.stopRealTimeUpdates).toBe('function')
  })

  it('handles timeline data with proper fallback', async () => {
    mockUrlAnalyticsService.getClickTimeline.mockResolvedValue({
      timeline: mockAnalyticsData.timeline
    })

    const { result } = renderHook(() => 
      useRealTimeAnalytics({ 
        urlId: 'test-url',
        enabled: true 
      })
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data?.timeline).toEqual(mockAnalyticsData.timeline)
  })
})