"""Type definitions for URL Shortener SDK."""

from .auth import *
from .urls import *
from .analytics import *
from .webhooks import *
from .qr import *
from .common import *

__all__ = [
    # Auth types
    "AuthCredentials",
    "AuthTokens", 
    "User",
    
    # URL types
    "ShortURL",
    "CreateURLRequest",
    "UpdateURLRequest",
    
    # Analytics types
    "URLAnalytics",
    "DashboardStats",
    "CountryStat",
    "DeviceStat",
    "BrowserStat",
    "Click",
    
    # Webhook types
    "Webhook",
    "WebhookEvent",
    "CreateWebhookRequest",
    "UpdateWebhookRequest",
    "WebhookDelivery",
    "WebhookStats",
    
    # QR types
    "QRCodeOptions",
    "QRCodeResponse",
    
    # Common types
    "PaginationParams",
    "PaginatedResponse",
    "APIError",
    "RequestConfig",
]