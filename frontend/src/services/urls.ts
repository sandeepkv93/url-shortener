import { apiService } from './api'
import {
  URL,
  CreateURLRequest,
  UpdateURLRequest,
  URLFilter,
  URLListResponse,
  URLStats,
  URLValidation,
  PopularURLsResponse,
  BulkUpdateRequest,
  BulkImportRequest,
  BulkImportResponse,
  ShareableURL,
  URLPreview,
  ExportURLData,
} from '@/types/url'

// URL management service
export const urlService = {
  // Create a new short URL
  async createURL(data: CreateURLRequest): Promise<URL> {
    return apiService.post<URL>('/urls', data)
  },

  // Get user's URLs with filtering and pagination
  async getUserURLs(filter?: URLFilter): Promise<URLListResponse> {
    const params = new URLSearchParams()
    
    if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            value.forEach(item => params.append(key, item.toString()))
          } else {
            params.append(key, value.toString())
          }
        }
      })
    }

    const queryString = params.toString()
    const url = queryString ? `/urls?${queryString}` : '/urls'
    
    return apiService.get<URLListResponse>(url)
  },

  // Get a specific URL by ID
  async getURL(id: string): Promise<URL> {
    return apiService.get<URL>(`/urls/${id}`)
  },

  // Get a URL by short code
  async getURLByShortCode(shortCode: string): Promise<URL> {
    return apiService.get<URL>(`/urls/code/${shortCode}`)
  },

  // Update a URL
  async updateURL(id: string, data: UpdateURLRequest): Promise<URL> {
    return apiService.put<URL>(`/urls/${id}`, data)
  },

  // Delete a URL
  async deleteURL(id: string): Promise<void> {
    return apiService.delete<void>(`/urls/${id}`)
  },

  // Bulk update URLs
  async bulkUpdateURLs(data: BulkUpdateRequest): Promise<{ updated: number; failed: number }> {
    return apiService.patch<{ updated: number; failed: number }>('/urls/bulk', data)
  },

  // Bulk delete URLs
  async bulkDeleteURLs(urlIds: string[]): Promise<{ deleted: number; failed: number }> {
    return apiService.delete<{ deleted: number; failed: number }>('/urls/bulk', {
      data: { urlIds }
    })
  },

  // Get URL statistics
  async getURLStats(id: string, period?: string): Promise<URLStats> {
    const url = period ? `/urls/${id}/stats?period=${period}` : `/urls/${id}/stats`
    return apiService.get<URLStats>(url)
  },

  // Validate URL before creation
  async validateURL(url: string): Promise<URLValidation> {
    return apiService.post<URLValidation>('/urls/validate', { url })
  },

  // Check custom alias availability
  async checkAliasAvailability(alias: string): Promise<{ available: boolean; suggestions?: string[] }> {
    return apiService.get<{ available: boolean; suggestions?: string[] }>(`/urls/alias/check?alias=${alias}`)
  },

  // Get popular URLs (public)
  async getPopularURLs(limit: number = 10): Promise<PopularURLsResponse> {
    return apiService.get<PopularURLsResponse>(`/urls/popular?limit=${limit}`)
  },

  // Get URL preview/metadata
  async getURLPreview(url: string): Promise<URLPreview> {
    return apiService.post<URLPreview>('/urls/preview', { url })
  },

  // Import URLs from CSV/JSON
  async importURLs(data: BulkImportRequest): Promise<BulkImportResponse> {
    return apiService.post<BulkImportResponse>('/urls/import', data)
  },

  // Export URLs to CSV/JSON
  async exportURLs(format: 'csv' | 'json', filter?: URLFilter): Promise<void> {
    const params = new URLSearchParams()
    params.append('format', format)
    
    if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            value.forEach(item => params.append(key, item.toString()))
          } else {
            params.append(key, value.toString())
          }
        }
      })
    }

    const filename = `urls-${format}-${new Date().toISOString().split('T')[0]}.${format}`
    return apiService.download(`/urls/export?${params.toString()}`, filename)
  },

  // Get shareable URL data
  async getShareableURL(shortCode: string): Promise<ShareableURL> {
    return apiService.get<ShareableURL>(`/urls/${shortCode}/share`)
  },

  // Toggle URL active status
  async toggleURLStatus(id: string): Promise<URL> {
    return apiService.patch<URL>(`/urls/${id}/toggle`)
  },

  // Duplicate URL
  async duplicateURL(id: string, customizations?: Partial<CreateURLRequest>): Promise<URL> {
    return apiService.post<URL>(`/urls/${id}/duplicate`, customizations)
  },

  // Get URL click history
  async getClickHistory(id: string, page: number = 1, limit: number = 50): Promise<{
    clicks: Array<{
      id: string
      timestamp: string
      ipAddress?: string
      userAgent?: string
      referer?: string
      country?: string
      city?: string
      device?: string
      browser?: string
    }>
    total: number
    page: number
    limit: number
  }> {
    return apiService.get(`/urls/${id}/clicks?page=${page}&limit=${limit}`)
  },

  // Reset URL statistics
  async resetURLStats(id: string): Promise<void> {
    return apiService.delete<void>(`/urls/${id}/stats`)
  },

  // Set URL password
  async setURLPassword(id: string, password: string): Promise<URL> {
    return apiService.post<URL>(`/urls/${id}/password`, { password })
  },

  // Remove URL password
  async removeURLPassword(id: string): Promise<URL> {
    return apiService.delete<URL>(`/urls/${id}/password`)
  },

  // Set URL expiration
  async setURLExpiration(id: string, expiresAt: string): Promise<URL> {
    return apiService.post<URL>(`/urls/${id}/expiration`, { expiresAt })
  },

  // Remove URL expiration
  async removeURLExpiration(id: string): Promise<URL> {
    return apiService.delete<URL>(`/urls/${id}/expiration`)
  },

  // Add tags to URL
  async addURLTags(id: string, tags: string[]): Promise<URL> {
    return apiService.post<URL>(`/urls/${id}/tags`, { tags })
  },

  // Remove tags from URL
  async removeURLTags(id: string, tags: string[]): Promise<URL> {
    return apiService.delete<URL>(`/urls/${id}/tags`, { data: { tags } })
  },

  // Get all user tags
  async getUserTags(): Promise<{ tags: string[] }> {
    return apiService.get<{ tags: string[] }>('/urls/tags')
  },

  // Archive URL (soft delete)
  async archiveURL(id: string): Promise<URL> {
    return apiService.post<URL>(`/urls/${id}/archive`)
  },

  // Restore archived URL
  async restoreURL(id: string): Promise<URL> {
    return apiService.post<URL>(`/urls/${id}/restore`)
  },

  // Get archived URLs
  async getArchivedURLs(page: number = 1, limit: number = 20): Promise<URLListResponse> {
    return apiService.get<URLListResponse>(`/urls/archived?page=${page}&limit=${limit}`)
  },

  // Permanently delete URL
  async permanentlyDeleteURL(id: string): Promise<void> {
    return apiService.delete<void>(`/urls/${id}/permanent`)
  },
}

// URL analytics service
export const urlAnalyticsService = {
  // Get detailed analytics for a URL
  async getDetailedAnalytics(urlId: string, period: string = '30d'): Promise<URLStats> {
    return apiService.get<URLStats>(`/analytics/urls/${urlId}?period=${period}`)
  },

  // Get click timeline for a URL
  async getClickTimeline(urlId: string, period: string = '7d'): Promise<{
    timeline: Array<{
      date: string
      clicks: number
      uniqueClicks: number
    }>
    total: number
    period: string
  }> {
    return apiService.get(`/analytics/urls/${urlId}/timeline?period=${period}`)
  },

  // Get geographic stats for a URL
  async getGeographicStats(urlId: string): Promise<{
    countries: Array<{
      country: string
      countryCode: string
      clicks: number
      percentage: number
    }>
    cities: Array<{
      city: string
      country: string
      clicks: number
      percentage: number
    }>
  }> {
    return apiService.get(`/analytics/urls/${urlId}/geo`)
  },

  // Get device stats for a URL
  async getDeviceStats(urlId: string): Promise<{
    devices: Array<{
      device: string
      clicks: number
      percentage: number
    }>
    browsers: Array<{
      browser: string
      clicks: number
      percentage: number
    }>
    operatingSystems: Array<{
      os: string
      clicks: number
      percentage: number
    }>
  }> {
    return apiService.get(`/analytics/urls/${urlId}/devices`)
  },

  // Get referrer stats for a URL
  async getReferrerStats(urlId: string): Promise<{
    referrers: Array<{
      referrer: string
      clicks: number
      percentage: number
    }>
    directClicks: number
    totalClicks: number
  }> {
    return apiService.get(`/analytics/urls/${urlId}/referrers`)
  },
}

export default urlService