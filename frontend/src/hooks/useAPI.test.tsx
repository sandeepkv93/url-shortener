import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useAPI, useMutation, usePagination, useOptimisticMutation, usePolling } from './useAPI'
import type { ApiError } from '@/types/auth'

// Mock API calls
const mockSuccessfulApiCall = vi.fn().mockResolvedValue({ data: 'success' })
const mockFailingApiCall = vi.fn().mockRejectedValue({
  message: 'API Error',
  status: 400,
} as ApiError)

describe('useAPI Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('useAPI', () => {
    it('should initialize with correct default state', () => {
      const { result } = renderHook(() => useAPI(mockSuccessfulApiCall))

      expect(result.current.data).toBeNull()
      expect(result.current.loading).toBe(false)
      expect(result.current.error).toBeNull()
      expect(result.current.isLoading).toBe(false)
      expect(result.current.isError).toBe(false)
      expect(result.current.isSuccess).toBe(false)
    })

    it('should execute API call immediately when immediate option is true', async () => {
      const { result } = renderHook(() => 
        useAPI(mockSuccessfulApiCall, { immediate: true })
      )

      expect(result.current.loading).toBe(true)

      await waitFor(() => {
        expect(result.current.loading).toBe(false)
      })

      expect(result.current.data).toEqual({ data: 'success' })
      expect(result.current.isSuccess).toBe(true)
      expect(mockSuccessfulApiCall).toHaveBeenCalledOnce()
    })

    it('should handle successful API call execution', async () => {
      const { result } = renderHook(() => useAPI(mockSuccessfulApiCall))

      await act(async () => {
        await result.current.execute()
      })

      expect(result.current.data).toEqual({ data: 'success' })
      expect(result.current.loading).toBe(false)
      expect(result.current.error).toBeNull()
      expect(result.current.isSuccess).toBe(true)
    })

    it('should handle API call errors', async () => {
      const { result } = renderHook(() => useAPI(mockFailingApiCall))

      await act(async () => {
        await result.current.execute()
      })

      expect(result.current.data).toBeNull()
      expect(result.current.loading).toBe(false)
      expect(result.current.error).toEqual({
        message: 'API Error',
        status: 400,
      })
      expect(result.current.isError).toBe(true)
    })

    it('should call onSuccess callback', async () => {
      const onSuccess = vi.fn()
      const { result } = renderHook(() => 
        useAPI(mockSuccessfulApiCall, { onSuccess })
      )

      await act(async () => {
        await result.current.execute()
      })

      expect(onSuccess).toHaveBeenCalledWith({ data: 'success' })
    })

    it('should call onError callback', async () => {
      const onError = vi.fn()
      const { result } = renderHook(() => 
        useAPI(mockFailingApiCall, { onError })
      )

      await act(async () => {
        await result.current.execute()
      })

      expect(onError).toHaveBeenCalledWith({
        message: 'API Error',
        status: 400,
      })
    })

    it('should reset state', async () => {
      const { result } = renderHook(() => useAPI(mockSuccessfulApiCall))

      await act(async () => {
        await result.current.execute()
      })

      expect(result.current.data).toEqual({ data: 'success' })

      act(() => {
        result.current.reset()
      })

      expect(result.current.data).toBeNull()
      expect(result.current.loading).toBe(false)
      expect(result.current.error).toBeNull()
    })

    it('should cancel ongoing request', async () => {
      const slowApiCall = vi.fn().mockImplementation(
        () => new Promise(resolve => setTimeout(() => resolve({ data: 'slow' }), 1000))
      )
      const { result } = renderHook(() => useAPI(slowApiCall))

      act(() => {
        result.current.execute()
      })

      expect(result.current.loading).toBe(true)

      act(() => {
        result.current.cancel()
      })

      expect(result.current.loading).toBe(false)
    })
  })

  describe('useMutation', () => {
    it('should handle mutation execution', async () => {
      const mutationFn = vi.fn().mockResolvedValue({ id: 1, name: 'Created' })
      const { result } = renderHook(() => useMutation(mutationFn))

      await act(async () => {
        await result.current.mutate({ name: 'Test' })
      })

      expect(result.current.data).toEqual({ id: 1, name: 'Created' })
      expect(result.current.isSuccess).toBe(true)
      expect(mutationFn).toHaveBeenCalledWith({ name: 'Test' })
    })

    it('should handle mutation errors', async () => {
      const mutationFn = vi.fn().mockRejectedValue({
        message: 'Mutation Error',
        status: 500,
      } as ApiError)
      const { result } = renderHook(() => useMutation(mutationFn))

      await act(async () => {
        await result.current.mutate({ name: 'Test' })
      })

      expect(result.current.data).toBeNull()
      expect(result.current.error).toEqual({
        message: 'Mutation Error',
        status: 500,
      })
      expect(result.current.isError).toBe(true)
    })
  })

  describe('usePagination', () => {
    const mockPaginatedApiCall = vi.fn()

    beforeEach(() => {
      mockPaginatedApiCall.mockImplementation((page: number, limit: number) =>
        Promise.resolve({
          data: [`item${page}-1`, `item${page}-2`],
          page,
          limit,
          total: 100,
        })
      )
    })

    it('should initialize with correct pagination state', () => {
      const { result } = renderHook(() => usePagination(mockPaginatedApiCall))

      expect(result.current.page).toBe(1)
      expect(result.current.limit).toBe(20)
    })

    it('should handle page navigation', async () => {
      const { result } = renderHook(() => 
        usePagination(mockPaginatedApiCall, { immediate: true })
      )

      await waitFor(() => {
        expect(result.current.loading).toBe(false)
      })

      expect(result.current.data).toEqual({
        data: ['item1-1', 'item1-2'],
        page: 1,
        limit: 20,
        total: 100,
      })

      act(() => {
        result.current.nextPage()
      })

      await waitFor(() => {
        expect(result.current.page).toBe(2)
      })

      await waitFor(() => {
        expect(result.current.data).toEqual({
          data: ['item2-1', 'item2-2'],
          page: 2,
          limit: 20,
          total: 100,
        })
      })
    })

    it('should go to previous page', async () => {
      const { result } = renderHook(() => usePagination(mockPaginatedApiCall))

      act(() => {
        result.current.goToPage(3)
      })

      await waitFor(() => {
        expect(result.current.page).toBe(3)
      })

      act(() => {
        result.current.prevPage()
      })

      expect(result.current.page).toBe(2)
    })

    it('should not go below page 1', () => {
      const { result } = renderHook(() => usePagination(mockPaginatedApiCall))

      act(() => {
        result.current.prevPage()
      })

      expect(result.current.page).toBe(1)
    })
  })

  describe('useOptimisticMutation', () => {
    it('should apply optimistic update', async () => {
      const mutationFn = vi.fn().mockImplementation(
        (params: any) => 
          new Promise(resolve => 
            setTimeout(() => resolve({ ...params, id: 1, confirmed: true }), 100)
          )
      )
      const optimisticUpdate = (params: any) => ({ ...params, id: 'temp', confirmed: false })

      const { result } = renderHook(() => 
        useOptimisticMutation(mutationFn, optimisticUpdate)
      )

      act(() => {
        result.current.mutate({ name: 'Test Item' })
      })

      // Should show optimistic data immediately
      expect(result.current.data).toEqual({
        name: 'Test Item',
        id: 'temp',
        confirmed: false,
      })
      expect(result.current.isOptimistic).toBe(true)

      // Wait for actual response
      await waitFor(() => {
        expect(result.current.isOptimistic).toBe(false)
      })

      expect(result.current.data).toEqual({
        name: 'Test Item',
        id: 1,
        confirmed: true,
      })
    })

    it('should clear optimistic data on error', async () => {
      const mutationFn = vi.fn().mockRejectedValue(new Error('Mutation failed'))
      const optimisticUpdate = (params: any) => ({ ...params, id: 'temp' })

      const { result } = renderHook(() => 
        useOptimisticMutation(mutationFn, optimisticUpdate)
      )

      await act(async () => {
        await result.current.mutate({ name: 'Test Item' })
      })

      expect(result.current.optimisticData).toBeNull()
      expect(result.current.isOptimistic).toBe(false)
      expect(result.current.isError).toBe(true)
    })
  })

  describe('usePolling', () => {
    beforeEach(() => {
      vi.useFakeTimers()
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('should poll data at specified interval', async () => {
      const pollingApiCall = vi.fn().mockResolvedValue({ timestamp: Date.now() })
      const { result } = renderHook(() => 
        usePolling(pollingApiCall, 1000, { enabled: true })
      )

      // Initial call
      await waitFor(() => {
        expect(pollingApiCall).toHaveBeenCalledTimes(1)
      })

      // Advance timer to trigger polling
      act(() => {
        vi.advanceTimersByTime(1000)
      })

      await waitFor(() => {
        expect(pollingApiCall).toHaveBeenCalledTimes(2)
      })

      // Advance timer again
      act(() => {
        vi.advanceTimersByTime(1000)
      })

      await waitFor(() => {
        expect(pollingApiCall).toHaveBeenCalledTimes(3)
      })
    })

    it('should not poll when disabled', () => {
      const pollingApiCall = vi.fn().mockResolvedValue({ timestamp: Date.now() })
      renderHook(() => 
        usePolling(pollingApiCall, 1000, { enabled: false })
      )

      act(() => {
        vi.advanceTimersByTime(5000)
      })

      expect(pollingApiCall).not.toHaveBeenCalled()
    })

    it('should stop polling when requested', async () => {
      const pollingApiCall = vi.fn().mockResolvedValue({ timestamp: Date.now() })
      const { result } = renderHook(() => 
        usePolling(pollingApiCall, 1000, { enabled: true })
      )

      // Wait for initial call
      await waitFor(() => {
        expect(pollingApiCall).toHaveBeenCalledTimes(1)
      })

      act(() => {
        result.current.stopPolling()
      })

      act(() => {
        vi.advanceTimersByTime(5000)
      })

      // Should not have called again after stopping
      expect(pollingApiCall).toHaveBeenCalledTimes(1)
    })

    it('should restart polling', async () => {
      const pollingApiCall = vi.fn().mockResolvedValue({ timestamp: Date.now() })
      const { result } = renderHook(() => 
        usePolling(pollingApiCall, 1000, { enabled: true })
      )

      await waitFor(() => {
        expect(pollingApiCall).toHaveBeenCalledTimes(1)
      })

      act(() => {
        result.current.stopPolling()
      })

      act(() => {
        result.current.startPolling()
      })

      act(() => {
        vi.advanceTimersByTime(1000)
      })

      await waitFor(() => {
        expect(pollingApiCall).toHaveBeenCalledTimes(2)
      })
    })
  })
})