"""URL Shortener Python SDK.

A comprehensive Python SDK for the URL Shortener API that provides:
- Synchronous and asynchronous clients
- Full type safety with Pydantic models
- Automatic token refresh and retry logic
- Comprehensive error handling
- Support for all API endpoints

Examples:
    Basic usage with API key:
    
    ```python
    from urlshortener import URLShortenerClient
    
    client = URLShortenerClient.with_api_key(
        base_url='https://api.urlshortener.com/v1',
        api_key='your-api-key'
    )
    
    # Create a short URL
    short_url = client.urls.create('https://example.com')
    print(f"Short URL: {short_url.short_url}")
    
    # Get analytics
    analytics = client.analytics.get_url_analytics(short_url.id)
    print(f"Total clicks: {analytics.total_clicks}")
    ```
    
    User authentication:
    
    ```python
    # Login with email/password
    auth_response = client.auth.login('user@example.com', 'password')
    
    # Client automatically uses the access token
    # Create URLs as authenticated user
    short_url = client.urls.create('https://example.com', title='My Link')
    ```
    
    Async usage:
    
    ```python
    from urlshortener import AsyncURLShortenerClient
    
    async def main():
        async with AsyncURLShortenerClient.with_api_key(
            base_url='https://api.urlshortener.com/v1',
            api_key='your-api-key'
        ) as client:
            short_url = await client.urls.create('https://example.com')
            analytics = await client.analytics.get_url_analytics(short_url.id)
    ```
"""

from .client import URLShortenerClient, AsyncURLShortenerClient, Client, AsyncClient

# Services
from .services.auth import AuthService, AsyncAuthService
from .services.urls import URLService, AsyncURLService
from .services.analytics import AnalyticsService, AsyncAnalyticsService
from .services.qr import QRService, AsyncQRService
from .services.webhooks import WebhookService, AsyncWebhookService

# HTTP Client
from .http_client import HTTPClient, AsyncHTTPClient

# Common types
from .types.common import (
    PaginationParams,
    PaginatedResponse,
    RequestConfig,
    APIResponse,
    APIError,
    AuthenticationError,
    ValidationError,
    NotFoundError,
    RateLimitError,
    NetworkError,
    HealthStatus,
    VersionInfo,
    TimePeriod,
    ExportFormat,
    SortOrder
)

# Auth types
from .types.auth import (
    AuthCredentials,
    AuthTokens,
    User,
    RegisterRequest,
    LoginRequest,
    AuthResponse,
    UpdateProfileRequest,
    ChangePasswordRequest,
    TokenValidationResponse,
    RefreshTokenRequest,
    PasswordResetRequest,
    PasswordResetConfirmRequest,
    EmailVerificationRequest,
    ResendVerificationRequest
)

# URL types
from .types.urls import (
    ShortURL,
    CreateURLRequest,
    UpdateURLRequest,
    BulkCreateURLRequest,
    BulkDeleteRequest,
    URLSearchRequest,
    URLStats,
    URLPerformance,
    PasswordValidationRequest,
    URLRedirectResponse,
    AliasAvailabilityResponse
)

# Analytics types
from .types.analytics import (
    Click,
    CountryStat,
    CityStat,
    DeviceStat,
    BrowserStat,
    OStat,
    ReferrerStat,
    TimelineData,
    URLAnalytics,
    DashboardStats,
    GeoStats,
    DeviceStats,
    ReferrerStats,
    AnalyticsFilter,
    AnalyticsExportRequest,
    RealTimeStats,
    ComparisonStats,
    ClickHeatmap,
    AnalyticsInsights
)

# QR code types
from .types.qr import (
    QRFormat,
    QRErrorCorrectionLevel,
    QRCodeOptions,
    QRCodeResponse,
    QRCodeFormats,
    QRCodeSizes,
    QRCodeGenerationRequest,
    QRCodeValidationResult,
    QRCodeBrandingOptions,
    QRCodeBatchRequest,
    QRCodeBatchResponse,
    QRCodeAnalytics,
    QRCodeMetadata,
    QRCodeTrackingParams,
    QRCodeTemplate,
    QRCodeDataValidation,
    QRCodeUsageStats
)

# Webhook types
from .types.webhooks import (
    WebhookEvent,
    WebhookStatus,
    WebhookDeliveryStatus,
    Webhook,
    CreateWebhookRequest,
    UpdateWebhookRequest,
    WebhookDelivery,
    WebhookStats,
    WebhookPayload,
    WebhookEventCategories,
    WebhookTestResult,
    WebhookHealthSummary,
    WebhookTemplate,
    WebhookInsights,
    WebhookFilter,
    WebhookDeliveryFilter
)

# Exceptions
from .exceptions import (
    URLShortenerError,
    AuthenticationError as AuthError,
    AuthorizationError,
    ValidationError as ValidError,
    NotFoundError as NotFound,
    ConflictError,
    RateLimitError as RateLimit,
    ServerError,
    NetworkError as NetError,
    TimeoutError,
    ConfigurationError,
    TokenExpiredError,
    RefreshTokenExpiredError,
    InvalidTokenError,
    QuotaExceededError,
    WebhookError,
    QRCodeError,
    AnalyticsError,
    URLError,
    InvalidURLError,
    URLExpiredError,
    URLPasswordProtectedError,
    CustomAliasUnavailableError,
    BulkOperationError,
    SDKError
)

__version__ = '1.0.0'
__author__ = 'URL Shortener Team'
__email__ = 'support@urlshortener.com'
__description__ = 'Python SDK for URL Shortener API'

__all__ = [
    # Main clients
    'URLShortenerClient',
    'AsyncURLShortenerClient', 
    'Client',
    'AsyncClient',
    
    # Services
    'AuthService',
    'AsyncAuthService',
    'URLService',
    'AsyncURLService',
    'AnalyticsService',
    'AsyncAnalyticsService',
    'QRService',
    'AsyncQRService',
    'WebhookService',
    'AsyncWebhookService',
    
    # HTTP clients
    'HTTPClient',
    'AsyncHTTPClient',
    
    # Common types
    'PaginationParams',
    'PaginatedResponse',
    'RequestConfig',
    'APIResponse',
    'HealthStatus',
    'VersionInfo',
    'TimePeriod',
    'ExportFormat',
    'SortOrder',
    
    # Auth types
    'AuthCredentials',
    'AuthTokens',
    'User',
    'RegisterRequest',
    'LoginRequest',
    'AuthResponse',
    'UpdateProfileRequest',
    'ChangePasswordRequest',
    'TokenValidationResponse',
    'RefreshTokenRequest',
    'PasswordResetRequest',
    'PasswordResetConfirmRequest',
    'EmailVerificationRequest',
    'ResendVerificationRequest',
    
    # URL types
    'ShortURL',
    'CreateURLRequest',
    'UpdateURLRequest',
    'BulkCreateURLRequest',
    'BulkDeleteRequest',
    'URLSearchRequest',
    'URLStats',
    'URLPerformance',
    'PasswordValidationRequest',
    'URLRedirectResponse',
    'AliasAvailabilityResponse',
    
    # Analytics types
    'Click',
    'CountryStat',
    'CityStat', 
    'DeviceStat',
    'BrowserStat',
    'OStat',
    'ReferrerStat',
    'TimelineData',
    'URLAnalytics',
    'DashboardStats',
    'GeoStats',
    'DeviceStats',
    'ReferrerStats',
    'AnalyticsFilter',
    'AnalyticsExportRequest',
    'RealTimeStats',
    'ComparisonStats',
    'ClickHeatmap',
    'AnalyticsInsights',
    
    # QR types
    'QRFormat',
    'QRErrorCorrectionLevel',
    'QRCodeOptions',
    'QRCodeResponse',
    'QRCodeFormats',
    'QRCodeSizes',
    'QRCodeGenerationRequest',
    'QRCodeValidationResult',
    'QRCodeBrandingOptions',
    'QRCodeBatchRequest',
    'QRCodeBatchResponse',
    'QRCodeAnalytics',
    'QRCodeMetadata',
    'QRCodeTrackingParams',
    'QRCodeTemplate',
    'QRCodeDataValidation',
    'QRCodeUsageStats',
    
    # Webhook types
    'WebhookEvent',
    'WebhookStatus',
    'WebhookDeliveryStatus',
    'Webhook',
    'CreateWebhookRequest',
    'UpdateWebhookRequest',
    'WebhookDelivery',
    'WebhookStats',
    'WebhookPayload',
    'WebhookEventCategories',
    'WebhookTestResult',
    'WebhookHealthSummary',
    'WebhookTemplate',
    'WebhookInsights',
    'WebhookFilter',
    'WebhookDeliveryFilter',
    
    # Exceptions
    'URLShortenerError',
    'AuthError',
    'AuthorizationError',
    'ValidError',
    'NotFound',
    'ConflictError',
    'RateLimit',
    'ServerError',
    'NetError',
    'TimeoutError',
    'ConfigurationError',
    'TokenExpiredError',
    'RefreshTokenExpiredError',
    'InvalidTokenError',
    'QuotaExceededError',
    'WebhookError',
    'QRCodeError',
    'AnalyticsError',
    'URLError',
    'InvalidURLError',
    'URLExpiredError',
    'URLPasswordProtectedError',
    'CustomAliasUnavailableError',
    'BulkOperationError',
    'SDKError',
]