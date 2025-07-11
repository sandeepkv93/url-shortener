// Dashboard analytics types
export interface DashboardStats {
  totalUrls: number
  activeUrls: number
  totalClicks: number
  uniqueClicks: number
  clicksToday: number
  clicksThisWeek: number
  clicksThisMonth: number
  urlsCreatedToday: number
  urlsCreatedThisWeek: number
  urlsCreatedThisMonth: number
  topPerformingUrl?: {
    shortCode: string
    originalUrl: string
    clickCount: number
  }
  recentActivity: ActivityItem[]
}

export interface ActivityItem {
  id: string
  type: 'url_created' | 'url_clicked' | 'url_updated' | 'url_deleted'
  description: string
  timestamp: string
  metadata?: Record<string, any>
}

// Time-based analytics
export interface TimelineData {
  date: string
  clicks: number
  uniqueClicks: number
  urls: number
}

export interface TimelineStats {
  period: string
  startDate: string
  endDate: string
  data: TimelineData[]
  totalClicks: number
  averageClicksPerDay: number
  peakDay: {
    date: string
    clicks: number
  }
}

// Geographic analytics
export interface GeoData {
  country: string
  countryCode: string
  clicks: number
  percentage: number
  cities: CityData[]
}

export interface CityData {
  city: string
  clicks: number
  percentage: number
}

export interface GeoStats {
  totalClicks: number
  countries: GeoData[]
  topCountry: {
    country: string
    clicks: number
    percentage: number
  }
  topCity: {
    city: string
    country: string
    clicks: number
    percentage: number
  }
}

// Device and browser analytics
export interface DeviceData {
  device: string
  clicks: number
  percentage: number
  browsers: BrowserData[]
}

export interface BrowserData {
  browser: string
  version?: string
  clicks: number
  percentage: number
}

export interface DeviceStats {
  totalClicks: number
  devices: DeviceData[]
  browsers: BrowserData[]
  operatingSystems: OSData[]
  topDevice: {
    device: string
    clicks: number
    percentage: number
  }
  topBrowser: {
    browser: string
    clicks: number
    percentage: number
  }
}

export interface OSData {
  os: string
  version?: string
  clicks: number
  percentage: number
}

// Referrer analytics
export interface ReferrerData {
  referrer: string
  domain: string
  clicks: number
  percentage: number
}

export interface ReferrerStats {
  totalClicks: number
  directClicks: number
  referrers: ReferrerData[]
  topReferrer: {
    referrer: string
    clicks: number
    percentage: number
  }
  socialMedia: ReferrerData[]
  searchEngines: ReferrerData[]
}

// URL performance analytics
export interface URLPerformance {
  id: string
  shortCode: string
  originalUrl: string
  title?: string
  clickCount: number
  uniqueClicks: number
  clickThroughRate: number
  averageClicksPerDay: number
  createdAt: string
  lastClickedAt?: string
  performance: 'excellent' | 'good' | 'average' | 'poor'
}

export interface TopPerformingURLs {
  period: string
  urls: URLPerformance[]
  total: number
}

// Comparison analytics
export interface ComparisonData {
  current: {
    period: string
    clicks: number
    urls: number
    uniqueClicks: number
  }
  previous: {
    period: string
    clicks: number
    urls: number
    uniqueClicks: number
  }
  growth: {
    clicks: number
    urls: number
    uniqueClicks: number
    clicksPercentage: number
    urlsPercentage: number
    uniqueClicksPercentage: number
  }
}

// Real-time analytics
export interface RealTimeData {
  activeVisitors: number
  clicksLast5Minutes: number
  topUrls: Array<{
    shortCode: string
    originalUrl: string
    recentClicks: number
  }>
  recentClicks: Array<{
    shortCode: string
    originalUrl: string
    timestamp: string
    country?: string
    device?: string
  }>
}

// Analytics export types
export interface AnalyticsExportRequest {
  period: string
  startDate?: string
  endDate?: string
  format: 'csv' | 'json' | 'pdf'
  includeDetails: boolean
  urlIds?: string[]
}

export interface AnalyticsExportResponse {
  downloadUrl: string
  filename: string
  size: number
  expiresAt: string
}

// Analytics filters
export interface AnalyticsFilter {
  period: '1h' | '24h' | '7d' | '30d' | '90d' | '1y' | 'custom'
  startDate?: string
  endDate?: string
  urlIds?: string[]
  countries?: string[]
  devices?: string[]
  browsers?: string[]
  referrers?: string[]
}

// Chart data types
export interface ChartDataPoint {
  label: string
  value: number
  percentage?: number
  color?: string
}

export interface TimeSeriesDataPoint {
  date: string
  value: number
  label?: string
}

export interface HeatmapDataPoint {
  x: string
  y: string
  value: number
  color?: string
}

// Analytics insights
export interface AnalyticsInsight {
  type: 'trend' | 'anomaly' | 'opportunity' | 'warning'
  title: string
  description: string
  value?: number
  change?: number
  recommendation?: string
  priority: 'high' | 'medium' | 'low'
}

export interface AnalyticsInsights {
  insights: AnalyticsInsight[]
  generatedAt: string
  period: string
}

// Global analytics (admin view)
export interface GlobalStats {
  totalUsers: number
  totalUrls: number
  totalClicks: number
  activeUsers: number
  newUsersToday: number
  urlsCreatedToday: number
  clicksToday: number
  topUsers: Array<{
    userId: string
    username: string
    urlCount: number
    clickCount: number
  }>
  systemHealth: {
    uptime: number
    responseTime: number
    errorRate: number
  }
}