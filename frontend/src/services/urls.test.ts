import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { urlService, urlAnalyticsService } from './urls'
import { apiService } from './api'
import type { CreateURLRequest, UpdateURLRequest, URLFilter } from '@/types/url'

// Mock the API service
vi.mock('./api', () => ({
  apiService: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
    download: vi.fn(),
  },
}))

describe('URL Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('urlService.createURL', () => {
    it('should create a new URL', async () => {
      const requestData: CreateURLRequest = {
        originalUrl: 'https://example.com',
        customAlias: 'test',
        title: 'Test URL',
      }
      const mockResponse = {
        id: '1',
        shortCode: 'test',
        originalUrl: 'https://example.com',
        title: 'Test URL',
        userId: 'user1',
        clickCount: 0,
        isActive: true,
        isPublic: true,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
      }

      vi.mocked(apiService.post).mockResolvedValue(mockResponse)

      const result = await urlService.createURL(requestData)

      expect(apiService.post).toHaveBeenCalledWith('/urls', requestData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.getUserURLs', () => {
    it('should get user URLs without filter', async () => {
      const mockResponse = {
        urls: [],
        total: 0,
        page: 1,
        limit: 20,
        hasNext: false,
        hasPrev: false,
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlService.getUserURLs()

      expect(apiService.get).toHaveBeenCalledWith('/urls')
      expect(result).toEqual(mockResponse)
    })

    it('should get user URLs with filter', async () => {
      const filter: URLFilter = {
        search: 'test',
        isActive: true,
        sortBy: 'createdAt',
        sortOrder: 'desc',
        page: 2,
        limit: 10,
      }
      const mockResponse = {
        urls: [],
        total: 0,
        page: 2,
        limit: 10,
        hasNext: false,
        hasPrev: true,
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlService.getUserURLs(filter)

      expect(apiService.get).toHaveBeenCalledWith(
        '/urls?search=test&isActive=true&sortBy=createdAt&sortOrder=desc&page=2&limit=10'
      )
      expect(result).toEqual(mockResponse)
    })

    it('should handle array filters', async () => {
      const filter: URLFilter = {
        tags: ['work', 'personal'],
      }

      vi.mocked(apiService.get).mockResolvedValue({
        urls: [],
        total: 0,
        page: 1,
        limit: 20,
        hasNext: false,
        hasPrev: false,
      })

      await urlService.getUserURLs(filter)

      expect(apiService.get).toHaveBeenCalledWith('/urls?tags=work&tags=personal')
    })
  })

  describe('urlService.getURL', () => {
    it('should get URL by ID', async () => {
      const mockURL = {
        id: '1',
        shortCode: 'test',
        originalUrl: 'https://example.com',
        userId: 'user1',
        clickCount: 5,
        isActive: true,
        isPublic: true,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
      }

      vi.mocked(apiService.get).mockResolvedValue(mockURL)

      const result = await urlService.getURL('1')

      expect(apiService.get).toHaveBeenCalledWith('/urls/1')
      expect(result).toEqual(mockURL)
    })
  })

  describe('urlService.updateURL', () => {
    it('should update URL', async () => {
      const updateData: UpdateURLRequest = {
        title: 'Updated Title',
        isActive: false,
      }
      const mockResponse = {
        id: '1',
        shortCode: 'test',
        originalUrl: 'https://example.com',
        title: 'Updated Title',
        userId: 'user1',
        clickCount: 5,
        isActive: false,
        isPublic: true,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-02T00:00:00Z',
      }

      vi.mocked(apiService.put).mockResolvedValue(mockResponse)

      const result = await urlService.updateURL('1', updateData)

      expect(apiService.put).toHaveBeenCalledWith('/urls/1', updateData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.deleteURL', () => {
    it('should delete URL', async () => {
      vi.mocked(apiService.delete).mockResolvedValue(undefined)

      await urlService.deleteURL('1')

      expect(apiService.delete).toHaveBeenCalledWith('/urls/1')
    })
  })

  describe('urlService.bulkUpdateURLs', () => {
    it('should bulk update URLs', async () => {
      const bulkData = {
        urlIds: ['1', '2'],
        updates: { isActive: false },
      }
      const mockResponse = { updated: 2, failed: 0 }

      vi.mocked(apiService.patch).mockResolvedValue(mockResponse)

      const result = await urlService.bulkUpdateURLs(bulkData)

      expect(apiService.patch).toHaveBeenCalledWith('/urls/bulk', bulkData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.validateURL', () => {
    it('should validate URL', async () => {
      const mockResponse = {
        isValid: true,
        error: undefined,
        suggestions: undefined,
      }

      vi.mocked(apiService.post).mockResolvedValue(mockResponse)

      const result = await urlService.validateURL('https://example.com')

      expect(apiService.post).toHaveBeenCalledWith('/urls/validate', {
        url: 'https://example.com',
      })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.checkAliasAvailability', () => {
    it('should check alias availability', async () => {
      const mockResponse = {
        available: false,
        suggestions: ['test1', 'test2'],
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlService.checkAliasAvailability('test')

      expect(apiService.get).toHaveBeenCalledWith('/urls/alias/check?alias=test')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.exportURLs', () => {
    it('should export URLs as CSV', async () => {
      const filter: URLFilter = {
        isActive: true,
        sortBy: 'createdAt',
      }

      vi.mocked(apiService.download).mockResolvedValue(undefined)

      await urlService.exportURLs('csv', filter)

      const expectedDate = new Date().toISOString().split('T')[0]
      expect(apiService.download).toHaveBeenCalledWith(
        '/urls/export?format=csv&isActive=true&sortBy=createdAt',
        `urls-csv-${expectedDate}.csv`
      )
    })

    it('should export URLs as JSON without filter', async () => {
      vi.mocked(apiService.download).mockResolvedValue(undefined)

      await urlService.exportURLs('json')

      const expectedDate = new Date().toISOString().split('T')[0]
      expect(apiService.download).toHaveBeenCalledWith(
        '/urls/export?format=json',
        `urls-json-${expectedDate}.json`
      )
    })
  })

  describe('urlService.toggleURLStatus', () => {
    it('should toggle URL status', async () => {
      const mockResponse = {
        id: '1',
        shortCode: 'test',
        originalUrl: 'https://example.com',
        userId: 'user1',
        clickCount: 5,
        isActive: false,
        isPublic: true,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-02T00:00:00Z',
      }

      vi.mocked(apiService.patch).mockResolvedValue(mockResponse)

      const result = await urlService.toggleURLStatus('1')

      expect(apiService.patch).toHaveBeenCalledWith('/urls/1/toggle')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.setURLPassword', () => {
    it('should set URL password', async () => {
      const mockResponse = {
        id: '1',
        shortCode: 'test',
        originalUrl: 'https://example.com',
        userId: 'user1',
        clickCount: 5,
        isActive: true,
        isPublic: false,
        password: 'hashed_password',
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-02T00:00:00Z',
      }

      vi.mocked(apiService.post).mockResolvedValue(mockResponse)

      const result = await urlService.setURLPassword('1', 'secret123')

      expect(apiService.post).toHaveBeenCalledWith('/urls/1/password', {
        password: 'secret123',
      })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlService.addURLTags', () => {
    it('should add tags to URL', async () => {
      const tags = ['work', 'important']
      const mockResponse = {
        id: '1',
        shortCode: 'test',
        originalUrl: 'https://example.com',
        userId: 'user1',
        clickCount: 5,
        isActive: true,
        isPublic: true,
        tags: ['work', 'important'],
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-02T00:00:00Z',
      }

      vi.mocked(apiService.post).mockResolvedValue(mockResponse)

      const result = await urlService.addURLTags('1', tags)

      expect(apiService.post).toHaveBeenCalledWith('/urls/1/tags', { tags })
      expect(result).toEqual(mockResponse)
    })
  })
})

describe('URL Analytics Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('urlAnalyticsService.getDetailedAnalytics', () => {
    it('should get detailed analytics with default period', async () => {
      const mockResponse = {
        url: {
          id: '1',
          shortCode: 'test',
          originalUrl: 'https://example.com',
          userId: 'user1',
          clickCount: 5,
          isActive: true,
          isPublic: true,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
        },
        totalClicks: 5,
        uniqueClicks: 3,
        clicksByDate: { '2023-01-01': 5 },
        clicksByCountry: { US: 3, CA: 2 },
        clicksByDevice: { desktop: 4, mobile: 1 },
        clicksByBrowser: { chrome: 3, firefox: 2 },
        clicksByReferrer: { direct: 5 },
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlAnalyticsService.getDetailedAnalytics('1')

      expect(apiService.get).toHaveBeenCalledWith('/analytics/urls/1?period=30d')
      expect(result).toEqual(mockResponse)
    })

    it('should get detailed analytics with custom period', async () => {
      const mockResponse = {
        url: {
          id: '1',
          shortCode: 'test',
          originalUrl: 'https://example.com',
          userId: 'user1',
          clickCount: 5,
          isActive: true,
          isPublic: true,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
        },
        totalClicks: 5,
        uniqueClicks: 3,
        clicksByDate: { '2023-01-01': 5 },
        clicksByCountry: { US: 3, CA: 2 },
        clicksByDevice: { desktop: 4, mobile: 1 },
        clicksByBrowser: { chrome: 3, firefox: 2 },
        clicksByReferrer: { direct: 5 },
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlAnalyticsService.getDetailedAnalytics('1', '7d')

      expect(apiService.get).toHaveBeenCalledWith('/analytics/urls/1?period=7d')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlAnalyticsService.getClickTimeline', () => {
    it('should get click timeline', async () => {
      const mockResponse = {
        timeline: [
          { date: '2023-01-01', clicks: 5, uniqueClicks: 3 },
          { date: '2023-01-02', clicks: 3, uniqueClicks: 2 },
        ],
        total: 8,
        period: '7d',
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlAnalyticsService.getClickTimeline('1')

      expect(apiService.get).toHaveBeenCalledWith('/analytics/urls/1/timeline?period=7d')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlAnalyticsService.getGeographicStats', () => {
    it('should get geographic stats', async () => {
      const mockResponse = {
        countries: [
          { country: 'United States', countryCode: 'US', clicks: 10, percentage: 50 },
          { country: 'Canada', countryCode: 'CA', clicks: 10, percentage: 50 },
        ],
        cities: [
          { city: 'New York', country: 'United States', clicks: 8, percentage: 40 },
          { city: 'Toronto', country: 'Canada', clicks: 6, percentage: 30 },
        ],
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlAnalyticsService.getGeographicStats('1')

      expect(apiService.get).toHaveBeenCalledWith('/analytics/urls/1/geo')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlAnalyticsService.getDeviceStats', () => {
    it('should get device stats', async () => {
      const mockResponse = {
        devices: [
          { device: 'desktop', clicks: 15, percentage: 75 },
          { device: 'mobile', clicks: 5, percentage: 25 },
        ],
        browsers: [
          { browser: 'Chrome', clicks: 12, percentage: 60 },
          { browser: 'Firefox', clicks: 8, percentage: 40 },
        ],
        operatingSystems: [
          { os: 'Windows', clicks: 10, percentage: 50 },
          { os: 'macOS', clicks: 10, percentage: 50 },
        ],
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlAnalyticsService.getDeviceStats('1')

      expect(apiService.get).toHaveBeenCalledWith('/analytics/urls/1/devices')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('urlAnalyticsService.getReferrerStats', () => {
    it('should get referrer stats', async () => {
      const mockResponse = {
        referrers: [
          { referrer: 'google.com', clicks: 10, percentage: 50 },
          { referrer: 'facebook.com', clicks: 5, percentage: 25 },
        ],
        directClicks: 5,
        totalClicks: 20,
      }

      vi.mocked(apiService.get).mockResolvedValue(mockResponse)

      const result = await urlAnalyticsService.getReferrerStats('1')

      expect(apiService.get).toHaveBeenCalledWith('/analytics/urls/1/referrers')
      expect(result).toEqual(mockResponse)
    })
  })
})