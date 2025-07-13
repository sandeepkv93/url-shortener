import { HTTPClient } from '../http-client';
import { API_ENDPOINTS } from '../config';
import { 
  URLAnalytics,
  DashboardStats,
  URLPerformance,
  CountryStat,
  CityStat,
  DeviceStat,
  BrowserStat,
  ReferrerStat,
  Click,
  RequestConfig 
} from '../types';

export interface TimelineData {
  date: string;
  clicks: number;
  uniqueClicks: number;
}

export interface GeoStats {
  countries: CountryStat[];
  cities: CityStat[];
  totalClicks: number;
  uniqueCountries: number;
  uniqueCities: number;
}

export interface DeviceStats {
  devices: DeviceStat[];
  browsers: BrowserStat[];
  operatingSystems: DeviceStat[];
  totalClicks: number;
}

export interface ReferrerStats {
  referrers: ReferrerStat[];
  domains: ReferrerStat[];
  totalClicks: number;
  directClicks: number;
}

export type ExportFormat = 'csv' | 'json' | 'xlsx';
export type TimePeriod = '24h' | '7d' | '30d' | '90d' | '1y' | 'all';

export class AnalyticsService {
  constructor(private client: HTTPClient) {}

  /**
   * Get dashboard analytics overview
   */
  async getDashboard(config?: RequestConfig): Promise<DashboardStats> {
    return this.client.get<DashboardStats>(API_ENDPOINTS.ANALYTICS_DASHBOARD, config);
  }

  /**
   * Get global analytics statistics
   */
  async getGlobalStats(config?: RequestConfig): Promise<{
    totalURLs: number;
    totalClicks: number;
    totalUsers: number;
    clicksToday: number;
    urlsCreatedToday: number;
    averageClicksPerURL: number;
  }> {
    return this.client.get(API_ENDPOINTS.ANALYTICS_GLOBAL, config);
  }

  /**
   * Get top performing URLs
   */
  async getTopURLs(
    limit?: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<URLPerformance[]> {
    const params = new URLSearchParams();
    if (limit) params.set('limit', limit.toString());
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_TOP_URLS}${params.toString() ? `?${params.toString()}` : ''}`;
    const response = await this.client.get<{ urls: URLPerformance[] }>(url, config);
    return response.urls;
  }

  /**
   * Get detailed analytics for a specific URL
   */
  async getURLAnalytics(
    urlId: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<URLAnalytics> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_URL(urlId)}${params.toString() ? `?${params.toString()}` : ''}`;
    return this.client.get<URLAnalytics>(url, config);
  }

  /**
   * Get click timeline for a URL
   */
  async getClickTimeline(
    urlId: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<TimelineData[]> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_URL_TIMELINE(urlId)}${params.toString() ? `?${params.toString()}` : ''}`;
    const response = await this.client.get<{ timeline: TimelineData[] }>(url, config);
    return response.timeline;
  }

  /**
   * Get geographic statistics for a URL
   */
  async getGeographicStats(
    urlId: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<GeoStats> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_URL_GEO(urlId)}${params.toString() ? `?${params.toString()}` : ''}`;
    return this.client.get<GeoStats>(url, config);
  }

  /**
   * Get device and browser statistics for a URL
   */
  async getDeviceStats(
    urlId: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<DeviceStats> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_URL_DEVICES(urlId)}${params.toString() ? `?${params.toString()}` : ''}`;
    return this.client.get<DeviceStats>(url, config);
  }

  /**
   * Get referrer statistics for a URL
   */
  async getReferrerStats(
    urlId: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<ReferrerStats> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_URL_REFERRERS(urlId)}${params.toString() ? `?${params.toString()}` : ''}`;
    return this.client.get<ReferrerStats>(url, config);
  }

  /**
   * Export analytics data
   */
  async exportAnalytics(
    format: ExportFormat,
    options: {
      urlIds?: number[];
      period?: TimePeriod;
      startDate?: string;
      endDate?: string;
      includeDetails?: boolean;
    } = {},
    config?: RequestConfig
  ): Promise<Blob> {
    const params = new URLSearchParams();
    params.set('format', format);
    
    if (options.urlIds?.length) {
      params.set('urlIds', options.urlIds.join(','));
    }
    if (options.period) params.set('period', options.period);
    if (options.startDate) params.set('startDate', options.startDate);
    if (options.endDate) params.set('endDate', options.endDate);
    if (options.includeDetails) params.set('includeDetails', 'true');

    const url = `${API_ENDPOINTS.ANALYTICS_EXPORT}?${params.toString()}`;
    
    // Override request config to handle blob response
    const requestConfig = {
      ...config,
      headers: {
        ...config?.headers,
        'Accept': this.getAcceptHeader(format),
      },
    };

    const response = await this.client.get<ArrayBuffer>(url, requestConfig);
    return new Blob([response], { type: this.getMimeType(format) });
  }

  /**
   * Get real-time analytics (last 30 minutes)
   */
  async getRealTimeStats(config?: RequestConfig): Promise<{
    activeClicks: number;
    recentClicks: Click[];
    topActiveURLs: URLPerformance[];
    clicksPerMinute: { time: string; clicks: number }[];
  }> {
    return this.client.get(`${API_ENDPOINTS.ANALYTICS_DASHBOARD}/realtime`, config);
  }

  /**
   * Get analytics comparison between two periods
   */
  async getComparison(
    urlId: number,
    currentPeriod: TimePeriod,
    previousPeriod: TimePeriod,
    config?: RequestConfig
  ): Promise<{
    current: URLAnalytics;
    previous: URLAnalytics;
    growth: {
      clicks: number;
      uniqueClicks: number;
      clicksPercentage: number;
      uniqueClicksPercentage: number;
    };
  }> {
    const params = new URLSearchParams();
    params.set('current', currentPeriod);
    params.set('previous', previousPeriod);

    const url = `${API_ENDPOINTS.ANALYTICS_URL(urlId)}/compare?${params.toString()}`;
    return this.client.get(url, config);
  }

  /**
   * Get analytics for multiple URLs
   */
  async getBulkAnalytics(
    urlIds: number[],
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<URLAnalytics[]> {
    const params = new URLSearchParams();
    params.set('urlIds', urlIds.join(','));
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_DASHBOARD}/bulk?${params.toString()}`;
    const response = await this.client.get<{ analytics: URLAnalytics[] }>(url, config);
    return response.analytics;
  }

  /**
   * Get click heatmap data for time-based visualization
   */
  async getClickHeatmap(
    urlId: number,
    period?: TimePeriod,
    config?: RequestConfig
  ): Promise<{
    hourly: { hour: number; clicks: number }[];
    daily: { day: number; clicks: number }[];
    weekly: { week: number; clicks: number }[];
  }> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.ANALYTICS_URL(urlId)}/heatmap${params.toString() ? `?${params.toString()}` : ''}`;
    return this.client.get(url, config);
  }

  private getAcceptHeader(format: ExportFormat): string {
    switch (format) {
      case 'csv':
        return 'text/csv';
      case 'json':
        return 'application/json';
      case 'xlsx':
        return 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
      default:
        return 'application/octet-stream';
    }
  }

  private getMimeType(format: ExportFormat): string {
    switch (format) {
      case 'csv':
        return 'text/csv';
      case 'json':
        return 'application/json';
      case 'xlsx':
        return 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
      default:
        return 'application/octet-stream';
    }
  }
}