import { HTTPClient } from '../http-client';
import { API_ENDPOINTS } from '../config';
import { 
  ShortURL, 
  CreateURLRequest, 
  UpdateURLRequest, 
  PaginationParams, 
  PaginatedResponse,
  RequestConfig 
} from '../types';

export class URLService {
  constructor(private client: HTTPClient) {}

  /**
   * Create a new short URL
   */
  async create(data: CreateURLRequest, config?: RequestConfig): Promise<ShortURL> {
    return this.client.post<ShortURL>(API_ENDPOINTS.URLS, data, config);
  }

  /**
   * Get all URLs for the authenticated user
   */
  async list(
    pagination?: PaginationParams,
    config?: RequestConfig
  ): Promise<PaginatedResponse<ShortURL>> {
    const params = new URLSearchParams();
    
    if (pagination?.limit) {
      params.set('limit', pagination.limit.toString());
    }
    if (pagination?.offset) {
      params.set('offset', pagination.offset.toString());
    }

    const url = `${API_ENDPOINTS.URLS}${params.toString() ? `?${params.toString()}` : ''}`;
    const response = await this.client.get<{
      urls: ShortURL[];
      total: number;
      limit: number;
      offset: number;
    }>(url, config);

    return {
      data: response.urls,
      total: response.total,
      limit: response.limit,
      offset: response.offset,
    };
  }

  /**
   * Get a specific URL by ID
   */
  async getById(id: number, config?: RequestConfig): Promise<ShortURL> {
    return this.client.get<ShortURL>(API_ENDPOINTS.URL_BY_ID(id), config);
  }

  /**
   * Update an existing URL
   */
  async update(id: number, data: UpdateURLRequest, config?: RequestConfig): Promise<ShortURL> {
    return this.client.put<ShortURL>(API_ENDPOINTS.URL_BY_ID(id), data, config);
  }

  /**
   * Delete a URL
   */
  async delete(id: number, config?: RequestConfig): Promise<void> {
    await this.client.delete(API_ENDPOINTS.URL_BY_ID(id), config);
  }

  /**
   * Get popular/trending URLs (public endpoint)
   */
  async getPopular(limit?: number, config?: RequestConfig): Promise<ShortURL[]> {
    const params = new URLSearchParams();
    if (limit) {
      params.set('limit', limit.toString());
    }

    const url = `${API_ENDPOINTS.URL_POPULAR}${params.toString() ? `?${params.toString()}` : ''}`;
    const response = await this.client.get<{ urls: ShortURL[] }>(url, config);
    return response.urls;
  }

  /**
   * Get URL by short code (for redirect purposes)
   */
  async getByShortCode(shortCode: string, config?: RequestConfig): Promise<ShortURL> {
    return this.client.get<ShortURL>(`/${shortCode}`, config);
  }

  /**
   * Validate a URL before shortening
   */
  async validateURL(url: string): Promise<boolean> {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Check if a custom alias is available
   */
  async checkAliasAvailability(alias: string, config?: RequestConfig): Promise<boolean> {
    try {
      await this.client.get(`/${alias}`, config);
      return false; // If we get a response, alias is taken
    } catch (error: any) {
      // If we get a 404, alias is available
      return error.statusCode === 404;
    }
  }

  /**
   * Bulk create URLs
   */
  async bulkCreate(urls: CreateURLRequest[], config?: RequestConfig): Promise<ShortURL[]> {
    return this.client.post<ShortURL[]>(`${API_ENDPOINTS.URLS}/bulk`, { urls }, config);
  }

  /**
   * Bulk delete URLs
   */
  async bulkDelete(ids: number[], config?: RequestConfig): Promise<void> {
    await this.client.post(`${API_ENDPOINTS.URLS}/bulk-delete`, { ids }, config);
  }

  /**
   * Search URLs by title or original URL
   */
  async search(
    query: string,
    pagination?: PaginationParams,
    config?: RequestConfig
  ): Promise<PaginatedResponse<ShortURL>> {
    const params = new URLSearchParams();
    params.set('q', query);
    
    if (pagination?.limit) {
      params.set('limit', pagination.limit.toString());
    }
    if (pagination?.offset) {
      params.set('offset', pagination.offset.toString());
    }

    const url = `${API_ENDPOINTS.URLS}/search?${params.toString()}`;
    const response = await this.client.get<{
      urls: ShortURL[];
      total: number;
      limit: number;
      offset: number;
    }>(url, config);

    return {
      data: response.urls,
      total: response.total,
      limit: response.limit,
      offset: response.offset,
    };
  }

  /**
   * Get URL statistics summary
   */
  async getStats(id: number, config?: RequestConfig): Promise<{
    totalClicks: number;
    uniqueClicks: number;
    clicksToday: number;
    clicksThisWeek: number;
    clicksThisMonth: number;
  }> {
    return this.client.get<{
      totalClicks: number;
      uniqueClicks: number;
      clicksToday: number;
      clicksThisWeek: number;
      clicksThisMonth: number;
    }>(`${API_ENDPOINTS.URL_BY_ID(id)}/stats`, config);
  }

  /**
   * Toggle URL active status
   */
  async toggleActive(id: number, config?: RequestConfig): Promise<ShortURL> {
    return this.client.post<ShortURL>(`${API_ENDPOINTS.URL_BY_ID(id)}/toggle`, {}, config);
  }

  /**
   * Duplicate an existing URL
   */
  async duplicate(id: number, config?: RequestConfig): Promise<ShortURL> {
    return this.client.post<ShortURL>(`${API_ENDPOINTS.URL_BY_ID(id)}/duplicate`, {}, config);
  }
}