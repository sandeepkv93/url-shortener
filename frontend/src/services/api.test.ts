import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import axios from 'axios'
import { apiService, handleApiError } from './api'
import { tokenManager } from './auth'

// Mock axios
vi.mock('axios')
const mockedAxios = vi.mocked(axios)

// Mock token manager
vi.mock('./auth', () => ({
  tokenManager: {
    getTokens: vi.fn(),
    clearTokens: vi.fn(),
    clearUser: vi.fn(),
  },
}))

// Mock environment variables
vi.mock('import.meta', () => ({
  env: {
    VITE_API_BASE_URL: 'http://localhost:8080',
    VITE_API_TIMEOUT: '10000',
    VITE_ENABLE_API_LOGGING: 'false',
  },
}))

// Mock axios instance
const mockAxiosInstance = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
  interceptors: {
    request: {
      use: vi.fn(),
    },
    response: {
      use: vi.fn(),
    },
  },
}

describe('API Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockedAxios.create.mockReturnValue(mockAxiosInstance as any)
    vi.mocked(tokenManager.getTokens).mockReturnValue({
      accessToken: 'mock-access-token',
      refreshToken: 'mock-refresh-token',
      expiresAt: new Date(Date.now() + 3600000).toISOString(),
    })
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('apiService.get', () => {
    it('should make GET request and return data', async () => {
      const mockData = { id: 1, name: 'Test' }
      mockAxiosInstance.get.mockResolvedValue({ data: mockData })

      const result = await apiService.get('/test')

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/test', undefined)
      expect(result).toEqual(mockData)
    })

    it('should handle API errors', async () => {
      const mockError = {
        response: {
          status: 404,
          statusText: 'Not Found',
          data: { message: 'Resource not found' },
        },
      }
      mockAxiosInstance.get.mockRejectedValue(mockError)

      await expect(apiService.get('/test')).rejects.toMatchObject({
        message: 'Resource not found',
        status: 404,
      })
    })
  })

  describe('apiService.post', () => {
    it('should make POST request with data', async () => {
      const mockData = { id: 1, name: 'Created' }
      const requestData = { name: 'New Item' }
      mockAxiosInstance.post.mockResolvedValue({ data: mockData })

      const result = await apiService.post('/test', requestData)

      expect(mockAxiosInstance.post).toHaveBeenCalledWith('/test', requestData, undefined)
      expect(result).toEqual(mockData)
    })

    it('should handle network errors', async () => {
      const mockError = {
        request: {},
        message: 'Network Error',
      }
      mockAxiosInstance.post.mockRejectedValue(mockError)

      await expect(apiService.post('/test', {})).rejects.toMatchObject({
        message: 'Network error. Please check your connection.',
        status: 0,
      })
    })
  })

  describe('apiService.put', () => {
    it('should make PUT request with data', async () => {
      const mockData = { id: 1, name: 'Updated' }
      const requestData = { name: 'Updated Item' }
      mockAxiosInstance.put.mockResolvedValue({ data: mockData })

      const result = await apiService.put('/test/1', requestData)

      expect(mockAxiosInstance.put).toHaveBeenCalledWith('/test/1', requestData, undefined)
      expect(result).toEqual(mockData)
    })
  })

  describe('apiService.patch', () => {
    it('should make PATCH request with data', async () => {
      const mockData = { id: 1, name: 'Patched' }
      const requestData = { name: 'Patched Item' }
      mockAxiosInstance.patch.mockResolvedValue({ data: mockData })

      const result = await apiService.patch('/test/1', requestData)

      expect(mockAxiosInstance.patch).toHaveBeenCalledWith('/test/1', requestData, undefined)
      expect(result).toEqual(mockData)
    })
  })

  describe('apiService.delete', () => {
    it('should make DELETE request', async () => {
      const mockData = { success: true }
      mockAxiosInstance.delete.mockResolvedValue({ data: mockData })

      const result = await apiService.delete('/test/1')

      expect(mockAxiosInstance.delete).toHaveBeenCalledWith('/test/1', undefined)
      expect(result).toEqual(mockData)
    })
  })

  describe('apiService.upload', () => {
    it('should upload file with FormData', async () => {
      const mockData = { fileId: 'test-file-id' }
      const mockFile = new File(['test content'], 'test.txt', { type: 'text/plain' })
      mockAxiosInstance.post.mockResolvedValue({ data: mockData })

      const result = await apiService.upload('/upload', mockFile)

      expect(mockAxiosInstance.post).toHaveBeenCalledWith(
        '/upload',
        expect.any(FormData),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'multipart/form-data',
          }),
        })
      )
      expect(result).toEqual(mockData)
    })
  })

  describe('apiService.download', () => {
    it('should download file and trigger download', async () => {
      const mockBlob = new Blob(['test content'], { type: 'text/plain' })
      mockAxiosInstance.get.mockResolvedValue({ data: mockBlob })

      // Mock DOM methods
      const mockCreateElement = vi.spyOn(document, 'createElement')
      const mockAppendChild = vi.spyOn(document.body, 'appendChild').mockImplementation(() => {} as any)
      const mockRemoveChild = vi.spyOn(document.body, 'removeChild').mockImplementation(() => {} as any)
      const mockCreateObjectURL = vi.spyOn(window.URL, 'createObjectURL').mockReturnValue('blob:mock-url')
      const mockRevokeObjectURL = vi.spyOn(window.URL, 'revokeObjectURL').mockImplementation(() => {})

      const mockLink = {
        href: '',
        download: '',
        click: vi.fn(),
      } as any
      mockCreateElement.mockReturnValue(mockLink)

      await apiService.download('/download', 'test.txt')

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/download', {
        responseType: 'blob',
      })
      expect(mockLink.download).toBe('test.txt')
      expect(mockLink.click).toHaveBeenCalled()
      expect(mockCreateObjectURL).toHaveBeenCalledWith(mockBlob)
      expect(mockRevokeObjectURL).toHaveBeenCalledWith('blob:mock-url')

      // Cleanup mocks
      mockCreateElement.mockRestore()
      mockAppendChild.mockRestore()
      mockRemoveChild.mockRestore()
      mockCreateObjectURL.mockRestore()
      mockRevokeObjectURL.mockRestore()
    })
  })

  describe('handleApiError', () => {
    it('should handle axios error with response', () => {
      const axiosError = {
        response: {
          status: 400,
          statusText: 'Bad Request',
          data: {
            message: 'Invalid data',
            code: 'VALIDATION_ERROR',
          },
        },
        isAxiosError: true,
      }

      const result = handleApiError(axiosError)

      expect(result).toEqual({
        message: 'Invalid data',
        status: 400,
        code: 'VALIDATION_ERROR',
      })
    })

    it('should handle axios error without response (network error)', () => {
      const axiosError = {
        request: {},
        isAxiosError: true,
      }

      const result = handleApiError(axiosError)

      expect(result).toEqual({
        message: 'Network error. Please check your connection.',
        status: 0,
      })
    })

    it('should handle generic error', () => {
      const genericError = new Error('Something went wrong')

      const result = handleApiError(genericError)

      expect(result).toEqual({
        message: 'Something went wrong',
        status: 0,
      })
    })

    it('should handle unknown error', () => {
      const unknownError = 'string error'

      const result = handleApiError(unknownError)

      expect(result).toEqual({
        message: 'An unexpected error occurred',
        status: 0,
      })
    })
  })
})