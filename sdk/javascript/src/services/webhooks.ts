import { HTTPClient } from '../http-client';
import { API_ENDPOINTS } from '../config';
import { 
  Webhook,
  WebhookEvent,
  CreateWebhookRequest,
  UpdateWebhookRequest,
  WebhookDelivery,
  WebhookStats,
  PaginationParams,
  PaginatedResponse,
  RequestConfig 
} from '../types';

export interface WebhookEventCategories {
  url_events: WebhookEvent[];
  analytics_events: WebhookEvent[];
  user_events: WebhookEvent[];
  system_events: WebhookEvent[];
}

export interface WebhookTestResult {
  delivery: WebhookDelivery;
  success: boolean;
  responseTime: number;
  statusCode: number;
  errorMessage?: string;
}

export class WebhookService {
  constructor(private client: HTTPClient) {}

  /**
   * Create a new webhook
   */
  async create(request: CreateWebhookRequest, config?: RequestConfig): Promise<Webhook> {
    return this.client.post<Webhook>(API_ENDPOINTS.WEBHOOKS, request, config);
  }

  /**
   * Get all webhooks for the authenticated user
   */
  async list(
    pagination?: PaginationParams,
    config?: RequestConfig
  ): Promise<PaginatedResponse<Webhook>> {
    const params = new URLSearchParams();
    
    if (pagination?.limit) {
      params.set('limit', pagination.limit.toString());
    }
    if (pagination?.offset) {
      params.set('offset', pagination.offset.toString());
    }

    const url = `${API_ENDPOINTS.WEBHOOKS}${params.toString() ? `?${params.toString()}` : ''}`;
    const response = await this.client.get<{
      webhooks: Webhook[];
      total: number;
      limit: number;
      offset: number;
    }>(url, config);

    return {
      data: response.webhooks,
      total: response.total,
      limit: response.limit,
      offset: response.offset,
    };
  }

  /**
   * Get a specific webhook by ID
   */
  async getById(id: number, config?: RequestConfig): Promise<Webhook> {
    return this.client.get<Webhook>(API_ENDPOINTS.WEBHOOK_BY_ID(id), config);
  }

  /**
   * Update an existing webhook
   */
  async update(id: number, request: UpdateWebhookRequest, config?: RequestConfig): Promise<Webhook> {
    return this.client.put<Webhook>(API_ENDPOINTS.WEBHOOK_BY_ID(id), request, config);
  }

  /**
   * Delete a webhook
   */
  async delete(id: number, config?: RequestConfig): Promise<void> {
    await this.client.delete(API_ENDPOINTS.WEBHOOK_BY_ID(id), config);
  }

  /**
   * Activate a webhook
   */
  async activate(id: number, config?: RequestConfig): Promise<void> {
    await this.client.post(API_ENDPOINTS.WEBHOOK_ACTIVATE(id), {}, config);
  }

  /**
   * Deactivate a webhook
   */
  async deactivate(id: number, config?: RequestConfig): Promise<void> {
    await this.client.post(API_ENDPOINTS.WEBHOOK_DEACTIVATE(id), {}, config);
  }

  /**
   * Test a webhook with a sample payload
   */
  async test(id: number, config?: RequestConfig): Promise<WebhookTestResult> {
    const delivery = await this.client.post<WebhookDelivery>(
      API_ENDPOINTS.WEBHOOK_TEST(id),
      {},
      config
    );

    return {
      delivery,
      success: delivery.status === 'success',
      responseTime: delivery.duration || 0,
      statusCode: delivery.responseStatus || 0,
      errorMessage: delivery.errorMessage,
    };
  }

  /**
   * Get webhook statistics
   */
  async getStats(id: number, config?: RequestConfig): Promise<WebhookStats> {
    return this.client.get<WebhookStats>(API_ENDPOINTS.WEBHOOK_STATS(id), config);
  }

  /**
   * Get webhook deliveries
   */
  async getDeliveries(
    id: number,
    pagination?: PaginationParams,
    config?: RequestConfig
  ): Promise<PaginatedResponse<WebhookDelivery>> {
    const params = new URLSearchParams();
    
    if (pagination?.limit) {
      params.set('limit', pagination.limit.toString());
    }
    if (pagination?.offset) {
      params.set('offset', pagination.offset.toString());
    }

    const url = `${API_ENDPOINTS.WEBHOOK_DELIVERIES(id)}${
      params.toString() ? `?${params.toString()}` : ''
    }`;
    
    const response = await this.client.get<{
      deliveries: WebhookDelivery[];
      total: number;
      limit: number;
      offset: number;
    }>(url, config);

    return {
      data: response.deliveries,
      total: response.total,
      limit: response.limit,
      offset: response.offset,
    };
  }

  /**
   * Get failed webhook deliveries
   */
  async getFailedDeliveries(
    id: number,
    limit?: number,
    config?: RequestConfig
  ): Promise<WebhookDelivery[]> {
    const params = new URLSearchParams();
    if (limit) params.set('limit', limit.toString());

    const url = `${API_ENDPOINTS.WEBHOOK_FAILED_DELIVERIES(id)}${
      params.toString() ? `?${params.toString()}` : ''
    }`;
    
    const response = await this.client.get<{ failed_deliveries: WebhookDelivery[] }>(url, config);
    return response.failed_deliveries;
  }

  /**
   * Get a specific webhook delivery
   */
  async getDelivery(deliveryId: number, config?: RequestConfig): Promise<WebhookDelivery> {
    return this.client.get<WebhookDelivery>(API_ENDPOINTS.WEBHOOK_DELIVERY(deliveryId), config);
  }

  /**
   * Retry a failed webhook delivery
   */
  async retryDelivery(deliveryId: number, config?: RequestConfig): Promise<void> {
    await this.client.post(API_ENDPOINTS.WEBHOOK_DELIVERY_RETRY(deliveryId), {}, config);
  }

  /**
   * Get available webhook events
   */
  async getEvents(config?: RequestConfig): Promise<WebhookEventCategories> {
    return this.client.get<WebhookEventCategories>(API_ENDPOINTS.WEBHOOK_EVENTS, config);
  }

  /**
   * Validate webhook URL
   */
  async validateURL(url: string): Promise<{ valid: boolean; error?: string }> {
    try {
      const urlObj = new URL(url);
      
      // Check protocol
      if (!['http:', 'https:'].includes(urlObj.protocol)) {
        return { valid: false, error: 'URL must use HTTP or HTTPS protocol' };
      }
      
      // Check for localhost or private IPs (basic check)
      const hostname = urlObj.hostname.toLowerCase();
      if (hostname === 'localhost' || 
          hostname === '127.0.0.1' || 
          hostname.startsWith('192.168.') ||
          hostname.startsWith('10.') ||
          hostname.startsWith('172.')) {
        return { valid: false, error: 'Webhook URL cannot point to localhost or private IP addresses' };
      }
      
      return { valid: true };
    } catch (error) {
      return { valid: false, error: 'Invalid URL format' };
    }
  }

  /**
   * Generate webhook signature for verification
   */
  generateSignature(payload: string, secret: string): string {
    // This would typically be done server-side, but including for completeness
    // In a real implementation, you'd use a crypto library
    const crypto = require('crypto');
    const hmac = crypto.createHmac('sha256', secret);
    hmac.update(payload);
    return `sha256=${hmac.digest('hex')}`;
  }

  /**
   * Verify webhook signature
   */
  verifySignature(payload: string, signature: string, secret: string): boolean {
    const expectedSignature = this.generateSignature(payload, secret);
    return signature === expectedSignature;
  }

  /**
   * Get webhook health summary
   */
  async getHealthSummary(config?: RequestConfig): Promise<{
    totalWebhooks: number;
    activeWebhooks: number;
    healthyWebhooks: number;
    failingWebhooks: number;
    recentDeliveries: number;
    averageSuccessRate: number;
  }> {
    const webhooks = await this.list({ limit: 1000 }, config);
    
    let activeCount = 0;
    let healthyCount = 0;
    let failingCount = 0;
    let totalSuccessRate = 0;
    let recentDeliveries = 0;

    for (const webhook of webhooks.data) {
      if (webhook.status === 'active') {
        activeCount++;
        
        const successRate = webhook.totalDeliveries > 0 
          ? (webhook.successDeliveries / webhook.totalDeliveries) * 100 
          : 100;
        
        totalSuccessRate += successRate;
        recentDeliveries += webhook.totalDeliveries;
        
        if (successRate >= 95) {
          healthyCount++;
        } else if (successRate < 80) {
          failingCount++;
        }
      }
    }

    return {
      totalWebhooks: webhooks.total,
      activeWebhooks: activeCount,
      healthyWebhooks: healthyCount,
      failingWebhooks: failingCount,
      recentDeliveries,
      averageSuccessRate: activeCount > 0 ? totalSuccessRate / activeCount : 0,
    };
  }

  /**
   * Bulk update webhook status
   */
  async bulkUpdateStatus(
    ids: number[],
    status: 'active' | 'inactive',
    config?: RequestConfig
  ): Promise<void> {
    await this.client.post(
      `${API_ENDPOINTS.WEBHOOKS}/bulk-status`,
      { ids, status },
      config
    );
  }

  /**
   * Bulk delete webhooks
   */
  async bulkDelete(ids: number[], config?: RequestConfig): Promise<void> {
    await this.client.post(
      `${API_ENDPOINTS.WEBHOOKS}/bulk-delete`,
      { ids },
      config
    );
  }

  /**
   * Get webhook templates for common use cases
   */
  async getTemplates(config?: RequestConfig): Promise<{
    templates: Array<{
      name: string;
      description: string;
      events: WebhookEvent[];
      samplePayload: any;
    }>;
  }> {
    return this.client.get(`${API_ENDPOINTS.WEBHOOKS}/templates`, config);
  }

  /**
   * Create webhook from template
   */
  async createFromTemplate(
    templateName: string,
    url: string,
    name: string,
    config?: RequestConfig
  ): Promise<Webhook> {
    return this.client.post<Webhook>(
      `${API_ENDPOINTS.WEBHOOKS}/from-template`,
      { templateName, url, name },
      config
    );
  }

  /**
   * Get delivery insights and recommendations
   */
  async getDeliveryInsights(id: number, config?: RequestConfig): Promise<{
    insights: Array<{
      type: 'error' | 'warning' | 'info';
      message: string;
      recommendation?: string;
    }>;
    performance: {
      averageResponseTime: number;
      successRate: number;
      errorRate: number;
      timeoutRate: number;
    };
    recommendations: string[];
  }> {
    return this.client.get(`${API_ENDPOINTS.WEBHOOK_BY_ID(id)}/insights`, config);
  }
}