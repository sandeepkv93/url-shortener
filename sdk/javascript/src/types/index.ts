// Authentication Types
export interface AuthCredentials {
  email: string;
  password: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
}

export interface User {
  id: number;
  email: string;
  createdAt: string;
  updatedAt: string;
}

// URL Types
export interface ShortURL {
  id: number;
  shortCode: string;
  originalURL: string;
  title?: string;
  description?: string;
  customAlias?: string;
  password?: string;
  expiresAt?: string;
  isActive: boolean;
  clickCount: number;
  userId: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateURLRequest {
  originalURL: string;
  title?: string;
  description?: string;
  customAlias?: string;
  password?: string;
  expiresAt?: string;
}

export interface UpdateURLRequest {
  title?: string;
  description?: string;
  password?: string;
  expiresAt?: string;
  isActive?: boolean;
}

// Analytics Types
export interface URLAnalytics {
  url: ShortURL;
  totalClicks: number;
  uniqueClicks: number;
  clicksByDate: { [date: string]: number };
  clicksByHour: { [hour: string]: number };
  topCountries: CountryStat[];
  topCities: CityStat[];
  topDevices: DeviceStat[];
  topBrowsers: BrowserStat[];
  topReferrers: ReferrerStat[];
  recentClicks: Click[];
}

export interface CountryStat {
  country: string;
  count: number;
  percentage: number;
}

export interface CityStat {
  city: string;
  country: string;
  count: number;
  percentage: number;
}

export interface DeviceStat {
  device: string;
  count: number;
  percentage: number;
}

export interface BrowserStat {
  browser: string;
  count: number;
  percentage: number;
}

export interface ReferrerStat {
  referrer: string;
  count: number;
  percentage: number;
}

export interface Click {
  id: number;
  shortURLId: number;
  ipAddress: string;
  userAgent: string;
  referer: string;
  country?: string;
  city?: string;
  device?: string;
  browser?: string;
  os?: string;
  clickedAt: string;
}

export interface DashboardStats {
  totalURLs: number;
  totalClicks: number;
  clicksToday: number;
  clicksThisWeek: number;
  clicksThisMonth: number;
  topURLs: URLPerformance[];
  recentURLs: ShortURL[];
}

export interface URLPerformance {
  url: ShortURL;
  clicks: number;
  uniqueClicks: number;
  lastClickedAt?: string;
}

// QR Code Types
export interface QRCodeOptions {
  size?: number;
  format?: 'png' | 'svg';
  level?: 'L' | 'M' | 'Q' | 'H';
  margin?: number;
  darkColor?: string;
  lightColor?: string;
}

export interface QRCodeResponse {
  data: string; // Base64 encoded image or SVG string
  format: string;
  size: number;
}

// Webhook Types
export type WebhookEvent = 
  | 'url.created'
  | 'url.updated'
  | 'url.deleted'
  | 'url.clicked'
  | 'url.expired'
  | 'analytics.threshold'
  | 'analytics.report'
  | 'user.registered'
  | 'user.updated'
  | 'system.error'
  | 'system.alert';

export type WebhookStatus = 'active' | 'inactive' | 'failed' | 'suspended';

export interface Webhook {
  id: number;
  userId: number;
  name: string;
  url: string;
  events: WebhookEvent[];
  secret?: string;
  status: WebhookStatus;
  maxRetries: number;
  timeoutSeconds: number;
  totalDeliveries: number;
  successDeliveries: number;
  failedDeliveries: number;
  lastDeliveryAt?: string;
  lastSuccessAt?: string;
  lastFailureAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateWebhookRequest {
  name: string;
  url: string;
  events: WebhookEvent[];
  secret?: string;
  maxRetries?: number;
  timeoutSeconds?: number;
}

export interface UpdateWebhookRequest {
  name?: string;
  url?: string;
  events?: WebhookEvent[];
  secret?: string;
  status?: WebhookStatus;
  maxRetries?: number;
  timeoutSeconds?: number;
}

export interface WebhookDelivery {
  id: number;
  webhookId: number;
  eventType: WebhookEvent;
  status: 'pending' | 'success' | 'failed' | 'retrying' | 'abandoned';
  requestUrl: string;
  requestHeaders?: any;
  requestBody?: any;
  responseStatus?: number;
  responseHeaders?: any;
  responseBody?: string;
  duration?: number;
  attemptCount: number;
  nextRetryAt?: string;
  errorMessage?: string;
  createdAt: string;
  updatedAt: string;
}

export interface WebhookStats {
  totalDeliveries: number;
  successDeliveries: number;
  failedDeliveries: number;
  successRate: number;
  averageResponseTime?: number;
  lastDeliveryAt?: string;
  lastSuccessAt?: string;
  lastFailureAt?: string;
}

// Pagination Types
export interface PaginationParams {
  limit?: number;
  offset?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}

// Error Types
export interface APIError {
  error: string;
  message?: string;
  field?: string;
  code?: string;
  statusCode: number;
}

// Request Configuration
export interface RequestConfig {
  headers?: { [key: string]: string };
  timeout?: number;
  signal?: AbortSignal;
}