import { HTTPClient } from '../http-client';
import { API_ENDPOINTS } from '../config';
import { QRCodeOptions, QRCodeResponse, RequestConfig } from '../types';

export interface QRCodeFormats {
  formats: string[];
  defaultFormat: string;
}

export interface QRCodeSizes {
  sizes: number[];
  defaultSize: number;
  minSize: number;
  maxSize: number;
}

export interface QRCodeGenerationRequest {
  url: string;
  options?: QRCodeOptions;
}

export interface QRCodeValidationResult {
  valid: boolean;
  errors?: string[];
  warnings?: string[];
}

export class QRService {
  constructor(private client: HTTPClient) {}

  /**
   * Get available QR code formats
   */
  async getFormats(config?: RequestConfig): Promise<QRCodeFormats> {
    return this.client.get<QRCodeFormats>(API_ENDPOINTS.QR_FORMATS, config);
  }

  /**
   * Get available QR code sizes
   */
  async getSizes(config?: RequestConfig): Promise<QRCodeSizes> {
    return this.client.get<QRCodeSizes>(API_ENDPOINTS.QR_SIZES, config);
  }

  /**
   * Validate QR code options
   */
  async validateOptions(
    options: QRCodeOptions,
    config?: RequestConfig
  ): Promise<QRCodeValidationResult> {
    return this.client.post<QRCodeValidationResult>(
      API_ENDPOINTS.QR_VALIDATE,
      options,
      config
    );
  }

  /**
   * Get QR code preview without saving
   */
  async getPreview(
    url: string,
    options?: QRCodeOptions,
    config?: RequestConfig
  ): Promise<QRCodeResponse> {
    return this.client.post<QRCodeResponse>(
      API_ENDPOINTS.QR_PREVIEW,
      { url, options },
      config
    );
  }

  /**
   * Generate QR code for any URL
   */
  async generate(
    request: QRCodeGenerationRequest,
    config?: RequestConfig
  ): Promise<QRCodeResponse> {
    return this.client.post<QRCodeResponse>(
      API_ENDPOINTS.QR_GENERATE,
      request,
      config
    );
  }

  /**
   * Generate QR code for a short URL by short code
   */
  async generateForURL(
    shortCode: string,
    options?: QRCodeOptions,
    config?: RequestConfig
  ): Promise<QRCodeResponse> {
    const params = new URLSearchParams();
    if (options?.size) params.set('size', options.size.toString());
    if (options?.format) params.set('format', options.format);
    if (options?.level) params.set('level', options.level);
    if (options?.margin) params.set('margin', options.margin.toString());
    if (options?.darkColor) params.set('darkColor', options.darkColor);
    if (options?.lightColor) params.set('lightColor', options.lightColor);

    const url = `${API_ENDPOINTS.QR_FOR_URL(shortCode)}${
      params.toString() ? `?${params.toString()}` : ''
    }`;
    
    return this.client.get<QRCodeResponse>(url, config);
  }

  /**
   * Download QR code as file
   */
  async download(
    shortCode: string,
    options?: QRCodeOptions & { filename?: string },
    config?: RequestConfig
  ): Promise<Blob> {
    const qrResponse = await this.generateForURL(shortCode, options, config);
    
    // Convert base64 to blob for download
    if (qrResponse.format === 'svg') {
      return new Blob([qrResponse.data], { type: 'image/svg+xml' });
    } else {
      // Handle base64 PNG data
      const binaryString = atob(qrResponse.data.replace(/^data:image\/png;base64,/, ''));
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      return new Blob([bytes], { type: 'image/png' });
    }
  }

  /**
   * Generate QR code with custom branding
   */
  async generateBranded(
    url: string,
    options: QRCodeOptions & {
      logo?: string; // Base64 encoded logo
      logoSize?: number; // Logo size as percentage of QR code
      backgroundColor?: string;
      foregroundColor?: string;
      borderRadius?: number;
    },
    config?: RequestConfig
  ): Promise<QRCodeResponse> {
    return this.client.post<QRCodeResponse>(
      `${API_ENDPOINTS.QR_GENERATE}/branded`,
      { url, options },
      config
    );
  }

  /**
   * Generate batch QR codes for multiple URLs
   */
  async generateBatch(
    requests: QRCodeGenerationRequest[],
    config?: RequestConfig
  ): Promise<QRCodeResponse[]> {
    const response = await this.client.post<{ qrCodes: QRCodeResponse[] }>(
      `${API_ENDPOINTS.QR_GENERATE}/batch`,
      { requests },
      config
    );
    return response.qrCodes;
  }

  /**
   * Get QR code analytics (if URL belongs to user)
   */
  async getAnalytics(
    shortCode: string,
    period?: '24h' | '7d' | '30d' | '90d',
    config?: RequestConfig
  ): Promise<{
    totalScans: number;
    scansToday: number;
    scansByDate: { date: string; scans: number }[];
    topDevices: { device: string; scans: number }[];
    topCountries: { country: string; scans: number }[];
  }> {
    const params = new URLSearchParams();
    if (period) params.set('period', period);

    const url = `${API_ENDPOINTS.QR_FOR_URL(shortCode)}/analytics${
      params.toString() ? `?${params.toString()}` : ''
    }`;
    
    return this.client.get(url, config);
  }

  /**
   * Create QR code with tracking parameters
   */
  async generateWithTracking(
    url: string,
    trackingParams: { [key: string]: string },
    options?: QRCodeOptions,
    config?: RequestConfig
  ): Promise<QRCodeResponse> {
    // Add tracking parameters to URL
    const urlObj = new URL(url);
    Object.entries(trackingParams).forEach(([key, value]) => {
      urlObj.searchParams.set(key, value);
    });

    return this.generate(
      {
        url: urlObj.toString(),
        options,
      },
      config
    );
  }

  /**
   * Get QR code metadata
   */
  async getMetadata(
    shortCode: string,
    config?: RequestConfig
  ): Promise<{
    shortCode: string;
    originalURL: string;
    title?: string;
    createdAt: string;
    expiresAt?: string;
    isActive: boolean;
    qrCodeGenerated: boolean;
    totalScans: number;
  }> {
    return this.client.get(`${API_ENDPOINTS.QR_FOR_URL(shortCode)}/metadata`, config);
  }

  /**
   * Validate QR code data
   */
  validateQRData(data: string): {
    isValid: boolean;
    type: 'url' | 'text' | 'email' | 'phone' | 'sms' | 'wifi' | 'other';
    data: string;
  } {
    // Basic URL validation
    try {
      new URL(data);
      return { isValid: true, type: 'url', data };
    } catch {
      // Check for other formats
      if (data.startsWith('mailto:')) {
        return { isValid: true, type: 'email', data };
      }
      if (data.startsWith('tel:')) {
        return { isValid: true, type: 'phone', data };
      }
      if (data.startsWith('sms:')) {
        return { isValid: true, type: 'sms', data };
      }
      if (data.startsWith('WIFI:')) {
        return { isValid: true, type: 'wifi', data };
      }
      
      return { isValid: true, type: 'text', data };
    }
  }

  /**
   * Get default QR code options
   */
  getDefaultOptions(): QRCodeOptions {
    return {
      size: 200,
      format: 'png',
      level: 'M',
      margin: 4,
      darkColor: '#000000',
      lightColor: '#FFFFFF',
    };
  }

  /**
   * Merge options with defaults
   */
  mergeOptions(options?: QRCodeOptions): QRCodeOptions {
    return {
      ...this.getDefaultOptions(),
      ...options,
    };
  }
}