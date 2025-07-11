import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import { ApiError } from '@/types/auth'
import { tokenManager } from './auth'

// API configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
const API_TIMEOUT = Number(import.meta.env.VITE_API_TIMEOUT) || 10000

// Create base API instance
const apiClient: AxiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  timeout: API_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const tokens = tokenManager.getTokens()
    if (tokens?.accessToken) {
      config.headers.Authorization = `Bearer ${tokens.accessToken}`
    }
    
    // Add request timestamp for debugging
    if (import.meta.env.VITE_ENABLE_API_LOGGING === 'true') {
      console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, {
        data: config.data,
        params: config.params,
        timestamp: new Date().toISOString(),
      })
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for error handling and token refresh
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // Log successful responses in development
    if (import.meta.env.VITE_ENABLE_API_LOGGING === 'true') {
      console.log(`[API Response] ${response.status} ${response.config.url}`, {
        data: response.data,
        timestamp: new Date().toISOString(),
      })
    }
    
    return response
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }
    
    // Log errors in development
    if (import.meta.env.VITE_ENABLE_API_LOGGING === 'true') {
      console.error(`[API Error] ${error.response?.status} ${originalRequest?.url}`, {
        error: error.response?.data,
        timestamp: new Date().toISOString(),
      })
    }
    
    // Handle 401 unauthorized errors with token refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      try {
        // Try to refresh the token
        const { authService } = await import('./auth')
        await authService.refreshToken()
        
        // Retry the original request with new token
        const tokens = tokenManager.getTokens()
        if (tokens?.accessToken && originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${tokens.accessToken}`
        }
        
        return apiClient(originalRequest)
      } catch (refreshError) {
        // Refresh failed, redirect to login
        tokenManager.clearTokens()
        tokenManager.clearUser()
        
        // Dispatch logout event for the auth context
        window.dispatchEvent(new CustomEvent('auth:logout'))
        
        return Promise.reject(refreshError)
      }
    }
    
    return Promise.reject(error)
  }
)

// Error handling utility
export const handleApiError = (error: unknown): ApiError => {
  if (axios.isAxiosError(error)) {
    if (error.response) {
      // Server responded with error status
      return {
        message: error.response.data?.message || error.response.statusText || 'Server error',
        status: error.response.status,
        code: error.response.data?.code,
      }
    } else if (error.request) {
      // Network error
      return {
        message: 'Network error. Please check your connection.',
        status: 0,
      }
    }
  }
  
  // Fallback for other errors
  return {
    message: error instanceof Error ? error.message : 'An unexpected error occurred',
    status: 0,
  }
}

// Generic API methods
export const apiService = {
  // GET request
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await apiClient.get<T>(url, config)
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // POST request
  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await apiClient.post<T>(url, data, config)
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // PUT request
  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await apiClient.put<T>(url, data, config)
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // PATCH request
  async patch<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await apiClient.patch<T>(url, data, config)
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // DELETE request
  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await apiClient.delete<T>(url, config)
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Upload file
  async upload<T>(url: string, file: File, config?: AxiosRequestConfig): Promise<T> {
    try {
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await apiClient.post<T>(url, formData, {
        ...config,
        headers: {
          ...config?.headers,
          'Content-Type': 'multipart/form-data',
        },
      })
      
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Download file
  async download(url: string, filename?: string, config?: AxiosRequestConfig): Promise<void> {
    try {
      const response = await apiClient.get(url, {
        ...config,
        responseType: 'blob',
      })
      
      // Create download link
      const blob = new Blob([response.data])
      const downloadUrl = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = downloadUrl
      link.download = filename || 'download'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(downloadUrl)
    } catch (error) {
      throw handleApiError(error)
    }
  },
}

// Export the configured axios instance for advanced usage
export { apiClient }
export default apiService