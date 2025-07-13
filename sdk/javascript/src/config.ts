export interface ClientConfig {
  baseURL: string;
  apiKey?: string;
  accessToken?: string;
  refreshToken?: string;
  timeout?: number;
  retries?: number;
  onTokenRefresh?: (tokens: { accessToken: string; refreshToken: string }) => void;
  onError?: (error: any) => void;
}

export const DEFAULT_CONFIG: Partial<ClientConfig> = {
  timeout: 30000,
  retries: 3,
};

export const API_ENDPOINTS = {
  // Authentication
  AUTH_REGISTER: '/auth/register',
  AUTH_LOGIN: '/auth/login',
  AUTH_REFRESH: '/auth/refresh',
  AUTH_LOGOUT: '/auth/logout',
  AUTH_PROFILE: '/auth/profile',
  AUTH_UPDATE_PROFILE: '/auth/profile',
  AUTH_CHANGE_PASSWORD: '/auth/change-password',
  AUTH_VALIDATE: '/auth/validate',

  // URLs
  URLS: '/urls',
  URL_BY_ID: (id: number) => `/urls/${id}`,
  URL_POPULAR: '/urls/popular',

  // Analytics
  ANALYTICS_DASHBOARD: '/analytics/dashboard',
  ANALYTICS_GLOBAL: '/analytics/global',
  ANALYTICS_TOP_URLS: '/analytics/top-urls',
  ANALYTICS_EXPORT: '/analytics/export',
  ANALYTICS_URL: (id: number) => `/analytics/urls/${id}`,
  ANALYTICS_URL_TIMELINE: (id: number) => `/analytics/urls/${id}/timeline`,
  ANALYTICS_URL_GEO: (id: number) => `/analytics/urls/${id}/geo`,
  ANALYTICS_URL_DEVICES: (id: number) => `/analytics/urls/${id}/devices`,
  ANALYTICS_URL_REFERRERS: (id: number) => `/analytics/urls/${id}/referrers`,

  // QR Codes
  QR_FORMATS: '/qr/formats',
  QR_SIZES: '/qr/sizes',
  QR_VALIDATE: '/qr/validate',
  QR_PREVIEW: '/qr/preview',
  QR_GENERATE: '/qr/generate',
  QR_FOR_URL: (shortCode: string) => `/qr/${shortCode}`,

  // Webhooks
  WEBHOOKS: '/webhooks',
  WEBHOOK_BY_ID: (id: number) => `/webhooks/${id}`,
  WEBHOOK_ACTIVATE: (id: number) => `/webhooks/${id}/activate`,
  WEBHOOK_DEACTIVATE: (id: number) => `/webhooks/${id}/deactivate`,
  WEBHOOK_TEST: (id: number) => `/webhooks/${id}/test`,
  WEBHOOK_STATS: (id: number) => `/webhooks/${id}/stats`,
  WEBHOOK_DELIVERIES: (id: number) => `/webhooks/${id}/deliveries`,
  WEBHOOK_FAILED_DELIVERIES: (id: number) => `/webhooks/${id}/failed-deliveries`,
  WEBHOOK_DELIVERY: (deliveryId: number) => `/webhooks/deliveries/${deliveryId}`,
  WEBHOOK_DELIVERY_RETRY: (deliveryId: number) => `/webhooks/deliveries/${deliveryId}/retry`,
  WEBHOOK_EVENTS: '/webhooks/events',

  // Health
  HEALTH: '/health',
  HEALTH_LIVEZ: '/health/livez',
  HEALTH_READYZ: '/health/readyz',
} as const;