import { useState, useCallback, useRef, useEffect } from 'react'
import { ApiError } from '@/types/auth'

// Generic API state type
export interface APIState<T> {
  data: T | null
  loading: boolean
  error: ApiError | null
}

// Hook options
export interface UseAPIOptions {
  immediate?: boolean
  onSuccess?: (data: any) => void
  onError?: (error: ApiError) => void
}

// Generic API hook for managing async operations
export function useAPI<T = any>(
  apiCall: (...args: any[]) => Promise<T>,
  options: UseAPIOptions = {}
) {
  const { immediate = false, onSuccess, onError } = options
  
  const [state, setState] = useState<APIState<T>>({
    data: null,
    loading: immediate,
    error: null,
  })
  
  const abortControllerRef = useRef<AbortController | null>(null)
  const mountedRef = useRef(true)

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      mountedRef.current = false
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [])

  const execute = useCallback(async (...args: any[]): Promise<T | null> => {
    // Cancel previous request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    // Create new abort controller
    abortControllerRef.current = new AbortController()

    setState(prev => ({
      ...prev,
      loading: true,
      error: null,
    }))

    try {
      const result = await apiCall(...args)
      
      if (mountedRef.current) {
        setState({
          data: result,
          loading: false,
          error: null,
        })
        
        onSuccess?.(result)
      }
      
      return result
    } catch (error) {
      if (mountedRef.current && !abortControllerRef.current?.signal.aborted) {
        const apiError = error as ApiError
        setState({
          data: null,
          loading: false,
          error: apiError,
        })
        
        onError?.(apiError)
      }
      
      return null
    }
  }, [apiCall, onSuccess, onError])

  const reset = useCallback(() => {
    setState({
      data: null,
      loading: false,
      error: null,
    })
  }, [])

  const cancel = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }
    setState(prev => ({
      ...prev,
      loading: false,
    }))
  }, [])

  // Execute immediately if requested
  useEffect(() => {
    if (immediate) {
      execute()
    }
  }, [immediate, execute])

  return {
    ...state,
    execute,
    reset,
    cancel,
    isLoading: state.loading,
    isError: !!state.error,
    isSuccess: !!state.data && !state.loading && !state.error,
  }
}

// Hook for pagination
export interface UsePaginationOptions<T> {
  pageSize?: number
  immediate?: boolean
  onSuccess?: (data: T) => void
  onError?: (error: ApiError) => void
}

export function usePagination<T = any>(
  apiCall: (page: number, limit: number) => Promise<T>,
  options: UsePaginationOptions<T> = {}
) {
  const { pageSize = 20, immediate = false, onSuccess, onError } = options
  
  const [page, setPage] = useState(1)
  const [limit] = useState(pageSize)
  
  const {
    data,
    loading,
    error,
    execute,
    reset,
    cancel,
    isLoading,
    isError,
    isSuccess,
  } = useAPI(
    () => apiCall(page, limit),
    { immediate, onSuccess, onError }
  )

  const nextPage = useCallback(() => {
    setPage(prev => prev + 1)
  }, [])

  const prevPage = useCallback(() => {
    setPage(prev => Math.max(1, prev - 1))
  }, [])

  const goToPage = useCallback((newPage: number) => {
    setPage(Math.max(1, newPage))
  }, [])

  const refresh = useCallback(() => {
    execute()
  }, [execute])

  const resetPagination = useCallback(() => {
    setPage(1)
    reset()
  }, [reset])

  // Re-execute when page changes
  useEffect(() => {
    if (page > 1 || immediate) {
      execute()
    }
  }, [page, execute, immediate])

  return {
    data,
    loading,
    error,
    page,
    limit,
    nextPage,
    prevPage,
    goToPage,
    refresh,
    reset: resetPagination,
    cancel,
    isLoading,
    isError,
    isSuccess,
  }
}

// Hook for mutations (POST, PUT, DELETE operations)
export function useMutation<T = any, P = any>(
  apiCall: (params: P) => Promise<T>,
  options: UseAPIOptions = {}
) {
  const {
    data,
    loading,
    error,
    execute,
    reset,
    cancel,
    isLoading,
    isError,
    isSuccess,
  } = useAPI(apiCall, { ...options, immediate: false })

  const mutate = useCallback(async (params: P): Promise<T | null> => {
    return execute(params)
  }, [execute])

  return {
    data,
    loading,
    error,
    mutate,
    reset,
    cancel,
    isLoading,
    isError,
    isSuccess,
  }
}

// Hook for optimistic updates
export function useOptimisticMutation<T = any, P = any>(
  apiCall: (params: P) => Promise<T>,
  optimisticUpdate: (params: P) => T,
  options: UseAPIOptions = {}
) {
  const [optimisticData, setOptimisticData] = useState<T | null>(null)
  
  const {
    data,
    loading,
    error,
    execute,
    reset: resetMutation,
    cancel,
    isLoading,
    isError,
    isSuccess,
  } = useAPI(apiCall, { ...options, immediate: false })

  const mutate = useCallback(async (params: P): Promise<T | null> => {
    // Apply optimistic update
    const optimistic = optimisticUpdate(params)
    setOptimisticData(optimistic)

    try {
      const result = await execute(params)
      setOptimisticData(null) // Clear optimistic data on success
      return result
    } catch (error) {
      setOptimisticData(null) // Clear optimistic data on error
      throw error
    }
  }, [execute, optimisticUpdate])

  const reset = useCallback(() => {
    setOptimisticData(null)
    resetMutation()
  }, [resetMutation])

  return {
    data: optimisticData || data,
    actualData: data,
    optimisticData,
    loading,
    error,
    mutate,
    reset,
    cancel,
    isLoading,
    isError,
    isSuccess,
    isOptimistic: !!optimisticData,
  }
}

// Hook for polling data
export function usePolling<T = any>(
  apiCall: () => Promise<T>,
  interval: number = 5000,
  options: UseAPIOptions & { enabled?: boolean } = {}
) {
  const { enabled = true, ...apiOptions } = options
  const intervalRef = useRef<NodeJS.Timeout | null>(null)
  
  const {
    data,
    loading,
    error,
    execute,
    reset,
    cancel,
    isLoading,
    isError,
    isSuccess,
  } = useAPI(apiCall, { ...apiOptions, immediate: enabled })

  const startPolling = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current)
    }
    
    intervalRef.current = setInterval(() => {
      execute()
    }, interval)
  }, [execute, interval])

  const stopPolling = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current)
      intervalRef.current = null
    }
  }, [])

  // Start/stop polling based on enabled flag
  useEffect(() => {
    if (enabled) {
      startPolling()
    } else {
      stopPolling()
    }

    return stopPolling
  }, [enabled, startPolling, stopPolling])

  // Cleanup on unmount
  useEffect(() => {
    return stopPolling
  }, [stopPolling])

  return {
    data,
    loading,
    error,
    execute,
    reset,
    cancel,
    startPolling,
    stopPolling,
    isLoading,
    isError,
    isSuccess,
    isPolling: !!intervalRef.current,
  }
}

export default useAPI