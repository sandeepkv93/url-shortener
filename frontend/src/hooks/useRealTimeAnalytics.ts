import { useState, useEffect, useRef, useCallback } from 'react'
import { urlService, urlAnalyticsService } from '@/services/urls'
import { AnalyticsPeriod } from '@/types/url'

interface UseRealTimeAnalyticsOptions {
  urlId?: string
  period?: AnalyticsPeriod
  refreshInterval?: number // in milliseconds
  enabled?: boolean
}

interface AnalyticsData {
  url?: any
  timeline?: any[]
  geographic?: any
  devices?: any
  referrers?: any
  clickHistory?: any
  lastUpdated: Date
}

export const useRealTimeAnalytics = ({
  urlId,
  period = '30d',
  refreshInterval = 30000, // 30 seconds default
  enabled = true
}: UseRealTimeAnalyticsOptions = {}) => {
  const [data, setData] = useState<AnalyticsData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [connectionStatus, setConnectionStatus] = useState<'connected' | 'disconnected' | 'reconnecting'>('connected')
  
  const intervalRef = useRef<NodeJS.Timeout | null>(null)
  const abortControllerRef = useRef<AbortController | null>(null)
  const retryCountRef = useRef(0)
  const maxRetries = 3

  const fetchAnalyticsData = useCallback(async (showLoading = true): Promise<AnalyticsData | null> => {
    // Cancel any ongoing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    // Create new abort controller for this request
    abortControllerRef.current = new AbortController()
    const signal = abortControllerRef.current.signal

    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)
    setError(null)
    setConnectionStatus('connected')

    try {
      const promises = []

      if (urlId) {
        // Fetch data for specific URL
        promises.push(
          urlService.getURL(urlId),
          urlAnalyticsService.getClickTimeline(urlId, period),
          urlAnalyticsService.getGeographicStats(urlId),
          urlAnalyticsService.getDeviceStats(urlId),
          urlAnalyticsService.getReferrerStats(urlId),
          urlService.getClickHistory(urlId, 1, 10)
        )
      } else {
        // Fetch aggregated data for dashboard
        promises.push(
          urlService.getUserURLs({ page: 1, limit: 100, sortBy: 'clickCount', sortOrder: 'desc' }),
          // Mock analytics calls for dashboard - in real implementation these would be actual API calls
          Promise.resolve({ timeline: [] }),
          Promise.resolve({ countries: [], cities: [] }),
          Promise.resolve({ devices: [], browsers: [], operatingSystems: [] }),
          Promise.resolve({ referrers: [], directClicks: 0, totalClicks: 0 }),
          Promise.resolve({ clicks: [], total: 0, page: 1, limit: 10 })
        )
      }

      // Check if request was aborted
      if (signal.aborted) {
        throw new Error('Request aborted')
      }

      const [
        urlOrUrls,
        timeline,
        geographic,
        devices,
        referrers,
        clickHistory
      ] = await Promise.all(promises)

      // Check again if request was aborted after async operations
      if (signal.aborted) {
        throw new Error('Request aborted')
      }

      const analyticsData: AnalyticsData = {
        url: urlId ? urlOrUrls : undefined,
        timeline: timeline.timeline || timeline,
        geographic,
        devices,
        referrers,
        clickHistory,
        lastUpdated: new Date()
      }

      // Reset retry count on successful fetch
      retryCountRef.current = 0
      setConnectionStatus('connected')
      
      return analyticsData
    } catch (err: any) {
      // Don't treat aborted requests as errors
      if (err.name === 'AbortError' || err.message === 'Request aborted') {
        return null
      }

      console.error('Analytics fetch error:', err)
      
      // Implement exponential backoff for retries
      if (retryCountRef.current < maxRetries) {
        retryCountRef.current++
        setConnectionStatus('reconnecting')
        
        const delay = Math.min(1000 * Math.pow(2, retryCountRef.current), 10000) // Max 10 seconds
        setTimeout(() => {
          if (enabled) {
            fetchAnalyticsData(false)
          }
        }, delay)
        
        return null
      } else {
        setConnectionStatus('disconnected')
        setError(err.message || 'Failed to load analytics data')
        throw err
      }
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }, [urlId, period, enabled])

  const refreshData = useCallback(async () => {
    try {
      const newData = await fetchAnalyticsData(false)
      if (newData) {
        setData(newData)
      }
    } catch (err) {
      // Error is already handled in fetchAnalyticsData
    }
  }, [fetchAnalyticsData])

  const startRealTimeUpdates = useCallback(() => {
    if (!enabled || intervalRef.current) return

    intervalRef.current = setInterval(refreshData, refreshInterval)
  }, [enabled, refreshInterval, refreshData])

  const stopRealTimeUpdates = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current)
      intervalRef.current = null
    }
  }, [])

  // Initial data fetch
  useEffect(() => {
    if (!enabled) return

    const loadInitialData = async () => {
      try {
        const initialData = await fetchAnalyticsData(true)
        if (initialData) {
          setData(initialData)
        }
      } catch (err) {
        // Error is already handled in fetchAnalyticsData
      }
    }

    loadInitialData()

    // Cleanup function
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [fetchAnalyticsData, enabled])

  // Start/stop real-time updates based on enabled flag
  useEffect(() => {
    if (enabled && data) {
      startRealTimeUpdates()
    } else {
      stopRealTimeUpdates()
    }

    return stopRealTimeUpdates
  }, [enabled, data, startRealTimeUpdates, stopRealTimeUpdates])

  // Handle visibility change to pause/resume updates when tab is not visible
  useEffect(() => {
    if (!enabled) return

    const handleVisibilityChange = () => {
      if (document.hidden) {
        stopRealTimeUpdates()
      } else if (data) {
        // Refresh data when tab becomes visible again
        refreshData()
        startRealTimeUpdates()
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)
    
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [enabled, data, refreshData, startRealTimeUpdates, stopRealTimeUpdates])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      stopRealTimeUpdates()
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [stopRealTimeUpdates])

  return {
    data,
    isLoading,
    isRefreshing,
    error,
    connectionStatus,
    lastUpdated: data?.lastUpdated,
    refreshData,
    startRealTimeUpdates,
    stopRealTimeUpdates
  }
}

export default useRealTimeAnalytics