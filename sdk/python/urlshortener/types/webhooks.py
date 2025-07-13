"""Webhook related type definitions."""

from typing import Any, Dict, List, Optional
from datetime import datetime
from enum import Enum
from pydantic import BaseModel, Field, HttpUrl, validator


class WebhookEvent(str, Enum):
    """Available webhook events."""
    # URL events
    URL_CREATED = "url.created"
    URL_UPDATED = "url.updated"
    URL_DELETED = "url.deleted"
    URL_CLICKED = "url.clicked"
    URL_EXPIRED = "url.expired"
    
    # Analytics events
    ANALYTICS_THRESHOLD = "analytics.threshold"
    ANALYTICS_REPORT = "analytics.report"
    
    # User events
    USER_REGISTERED = "user.registered"
    USER_UPDATED = "user.updated"
    
    # System events
    SYSTEM_ERROR = "system.error"
    SYSTEM_ALERT = "system.alert"


class WebhookStatus(str, Enum):
    """Webhook status values."""
    ACTIVE = "active"
    INACTIVE = "inactive"
    FAILED = "failed"
    SUSPENDED = "suspended"


class WebhookDeliveryStatus(str, Enum):
    """Webhook delivery status values."""
    PENDING = "pending"
    SUCCESS = "success"
    FAILED = "failed"
    RETRYING = "retrying"
    ABANDONED = "abandoned"


class Webhook(BaseModel):
    """Webhook configuration."""
    id: int = Field(description="Webhook ID")
    user_id: int = Field(description="Owner user ID")
    name: str = Field(description="Webhook name")
    url: HttpUrl = Field(description="Webhook endpoint URL")
    events: List[WebhookEvent] = Field(description="Subscribed events")
    secret: Optional[str] = Field(None, description="Webhook secret for verification")
    status: WebhookStatus = Field(description="Webhook status")
    
    # Configuration
    max_retries: int = Field(default=3, description="Maximum retry attempts")
    timeout_seconds: int = Field(default=30, description="Request timeout")
    retry_backoff_ms: int = Field(default=1000, description="Retry backoff in milliseconds")
    
    # Statistics
    total_deliveries: int = Field(default=0, description="Total delivery attempts")
    success_deliveries: int = Field(default=0, description="Successful deliveries")
    failed_deliveries: int = Field(default=0, description="Failed deliveries")
    last_delivery_at: Optional[datetime] = Field(None, description="Last delivery timestamp")
    last_success_at: Optional[datetime] = Field(None, description="Last successful delivery")
    last_failure_at: Optional[datetime] = Field(None, description="Last failed delivery")
    
    # Metadata
    created_at: datetime = Field(description="Creation timestamp")
    updated_at: datetime = Field(description="Last update timestamp")
    
    @property
    def success_rate(self) -> float:
        """Calculate webhook success rate."""
        if self.total_deliveries == 0:
            return 0.0
        return (self.success_deliveries / self.total_deliveries) * 100.0
    
    @property
    def is_healthy(self) -> bool:
        """Check if webhook is healthy (>95% success rate)."""
        return self.success_rate >= 95.0


class CreateWebhookRequest(BaseModel):
    """Request to create a new webhook."""
    name: str = Field(min_length=1, max_length=255, description="Webhook name")
    url: HttpUrl = Field(description="Webhook endpoint URL")
    events: List[WebhookEvent] = Field(min_items=1, description="Events to subscribe to")
    secret: Optional[str] = Field(None, max_length=255, description="Webhook secret")
    max_retries: Optional[int] = Field(None, ge=0, le=10, description="Maximum retry attempts")
    timeout_seconds: Optional[int] = Field(None, ge=1, le=300, description="Request timeout")
    
    @validator('url')
    def validate_url(cls, v):
        """Validate webhook URL."""
        url_str = str(v)
        if not url_str.startswith(('http://', 'https://')):
            raise ValueError('URL must use HTTP or HTTPS protocol')
        
        # Check for localhost and private IPs
        from urllib.parse import urlparse
        parsed = urlparse(url_str)
        hostname = parsed.hostname.lower()
        
        if hostname in ['localhost', '127.0.0.1', '::1']:
            raise ValueError('Webhook URL cannot point to localhost')
            
        if hostname.startswith(('192.168.', '10.', '172.')):
            raise ValueError('Webhook URL cannot point to private IP addresses')
            
        return v


class UpdateWebhookRequest(BaseModel):
    """Request to update an existing webhook."""
    name: Optional[str] = Field(None, min_length=1, max_length=255, description="Webhook name")
    url: Optional[HttpUrl] = Field(None, description="Webhook endpoint URL")
    events: Optional[List[WebhookEvent]] = Field(None, min_items=1, description="Events to subscribe to")
    secret: Optional[str] = Field(None, max_length=255, description="Webhook secret")
    status: Optional[WebhookStatus] = Field(None, description="Webhook status")
    max_retries: Optional[int] = Field(None, ge=0, le=10, description="Maximum retry attempts")
    timeout_seconds: Optional[int] = Field(None, ge=1, le=300, description="Request timeout")


class WebhookDelivery(BaseModel):
    """Webhook delivery attempt record."""
    id: int = Field(description="Delivery ID")
    webhook_id: int = Field(description="Webhook ID")
    event_type: WebhookEvent = Field(description="Event type")
    status: WebhookDeliveryStatus = Field(description="Delivery status")
    
    # Request details
    request_url: str = Field(description="Request URL")
    request_headers: Optional[Dict[str, Any]] = Field(None, description="Request headers")
    request_body: Optional[Dict[str, Any]] = Field(None, description="Request body")
    
    # Response details
    response_status: Optional[int] = Field(None, description="HTTP response status")
    response_headers: Optional[Dict[str, Any]] = Field(None, description="Response headers")
    response_body: Optional[str] = Field(None, description="Response body")
    
    # Timing and retry info
    duration: Optional[int] = Field(None, description="Request duration in milliseconds")
    attempt_count: int = Field(default=1, description="Attempt number")
    next_retry_at: Optional[datetime] = Field(None, description="Next retry timestamp")
    
    # Error handling
    error_message: Optional[str] = Field(None, description="Error message if failed")
    
    # Metadata
    created_at: datetime = Field(description="Creation timestamp")
    updated_at: datetime = Field(description="Last update timestamp")


class WebhookStats(BaseModel):
    """Webhook delivery statistics."""
    total_deliveries: int = Field(description="Total delivery attempts")
    success_deliveries: int = Field(description="Successful deliveries")
    failed_deliveries: int = Field(description="Failed deliveries")
    success_rate: float = Field(description="Success rate percentage")
    average_response_time: Optional[int] = Field(None, description="Average response time in ms")
    last_delivery_at: Optional[datetime] = Field(None, description="Last delivery timestamp")
    last_success_at: Optional[datetime] = Field(None, description="Last successful delivery")
    last_failure_at: Optional[datetime] = Field(None, description="Last failed delivery")


class WebhookPayload(BaseModel):
    """Webhook payload structure."""
    id: str = Field(description="Event ID")
    event: WebhookEvent = Field(description="Event type")
    data: Dict[str, Any] = Field(description="Event data")
    timestamp: datetime = Field(description="Event timestamp")
    user_id: int = Field(description="User ID")
    version: str = Field(default="1.0", description="Payload version")


class WebhookEventCategories(BaseModel):
    """Available webhook event categories."""
    url_events: List[WebhookEvent] = Field(description="URL-related events")
    analytics_events: List[WebhookEvent] = Field(description="Analytics events")
    user_events: List[WebhookEvent] = Field(description="User events")
    system_events: List[WebhookEvent] = Field(description="System events")


class WebhookTestResult(BaseModel):
    """Webhook test result."""
    delivery: WebhookDelivery = Field(description="Test delivery record")
    success: bool = Field(description="Whether test was successful")
    response_time: int = Field(description="Response time in milliseconds")
    status_code: int = Field(description="HTTP status code")
    error_message: Optional[str] = Field(None, description="Error message if failed")


class WebhookHealthSummary(BaseModel):
    """Webhook health summary."""
    total_webhooks: int = Field(description="Total number of webhooks")
    active_webhooks: int = Field(description="Number of active webhooks")
    healthy_webhooks: int = Field(description="Number of healthy webhooks")
    failing_webhooks: int = Field(description="Number of failing webhooks")
    recent_deliveries: int = Field(description="Recent delivery count")
    average_success_rate: float = Field(description="Average success rate across all webhooks")


class WebhookTemplate(BaseModel):
    """Webhook template for common use cases."""
    name: str = Field(description="Template name")
    description: str = Field(description="Template description")
    events: List[WebhookEvent] = Field(description="Recommended events")
    sample_payload: Dict[str, Any] = Field(description="Sample payload structure")


class WebhookInsights(BaseModel):
    """Webhook delivery insights and recommendations."""
    insights: List[Dict[str, Any]] = Field(description="Generated insights")
    performance: Dict[str, Any] = Field(description="Performance metrics")
    recommendations: List[str] = Field(description="Improvement recommendations")


class WebhookFilter(BaseModel):
    """Webhook filtering options."""
    status: Optional[WebhookStatus] = Field(None, description="Filter by status")
    events: Optional[List[WebhookEvent]] = Field(None, description="Filter by events")
    created_after: Optional[datetime] = Field(None, description="Filter by creation date")
    created_before: Optional[datetime] = Field(None, description="Filter by creation date")


class WebhookDeliveryFilter(BaseModel):
    """Webhook delivery filtering options."""
    webhook_id: Optional[int] = Field(None, description="Filter by webhook ID")
    status: Optional[WebhookDeliveryStatus] = Field(None, description="Filter by status")
    event_type: Optional[WebhookEvent] = Field(None, description="Filter by event type")
    from_date: Optional[datetime] = Field(None, description="Filter from date")
    to_date: Optional[datetime] = Field(None, description="Filter to date")