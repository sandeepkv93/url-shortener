"""Webhook service for URL Shortener API."""

from typing import List, Optional, Dict, Any

from ..http_client import HTTPClient, AsyncHTTPClient
from ..types.webhooks import (
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
    WebhookDeliveryFilter,
    WebhookEvent,
    WebhookStatus
)
from ..types.common import RequestConfig, PaginationParams, PaginatedResponse


class WebhookService:
    """Synchronous webhook service."""

    def __init__(self, http_client: HTTPClient):
        """Initialize webhook service.
        
        Args:
            http_client: HTTP client instance
        """
        self.http_client = http_client

    def create(
        self,
        name: str,
        url: str,
        events: List[WebhookEvent],
        secret: Optional[str] = None,
        max_retries: Optional[int] = None,
        timeout_seconds: Optional[int] = None,
        config: Optional[RequestConfig] = None
    ) -> Webhook:
        """Create a new webhook.
        
        Args:
            name: Webhook name
            url: Webhook endpoint URL
            events: Events to subscribe to
            secret: Webhook secret for verification
            max_retries: Maximum retry attempts
            timeout_seconds: Request timeout
            config: Request configuration
            
        Returns:
            Webhook: Created webhook
        """
        request = CreateWebhookRequest(
            name=name,
            url=url,
            events=events,
            secret=secret,
            max_retries=max_retries,
            timeout_seconds=timeout_seconds
        )
        
        response = self.http_client.post(
            '/webhooks',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return Webhook(**response['data'])

    def get(self, webhook_id: int, config: Optional[RequestConfig] = None) -> Webhook:
        """Get a webhook by ID.
        
        Args:
            webhook_id: Webhook ID
            config: Request configuration
            
        Returns:
            Webhook: Webhook details
        """
        response = self.http_client.get(f'/webhooks/{webhook_id}', config=config)
        return Webhook(**response['data'])

    def list(
        self,
        pagination: Optional[PaginationParams] = None,
        filter_params: Optional[WebhookFilter] = None,
        config: Optional[RequestConfig] = None
    ) -> PaginatedResponse[Webhook]:
        """List user's webhooks.
        
        Args:
            pagination: Pagination parameters
            filter_params: Filter parameters
            config: Request configuration
            
        Returns:
            PaginatedResponse[Webhook]: Paginated list of webhooks
        """
        params = {}
        
        if pagination:
            params.update(pagination.dict(exclude_none=True))
            
        if filter_params:
            params.update(filter_params.dict(exclude_none=True))
            
        response = self.http_client.get('/webhooks', params=params, config=config)
        
        webhooks = [Webhook(**item) for item in response['data']]
        return PaginatedResponse[Webhook](
            data=webhooks,
            total=response['total'],
            limit=response['limit'],
            offset=response['offset']
        )

    def update(
        self,
        webhook_id: int,
        name: Optional[str] = None,
        url: Optional[str] = None,
        events: Optional[List[WebhookEvent]] = None,
        secret: Optional[str] = None,
        status: Optional[WebhookStatus] = None,
        max_retries: Optional[int] = None,
        timeout_seconds: Optional[int] = None,
        config: Optional[RequestConfig] = None
    ) -> Webhook:
        """Update a webhook.
        
        Args:
            webhook_id: Webhook ID
            name: Webhook name
            url: Webhook endpoint URL
            events: Events to subscribe to
            secret: Webhook secret
            status: Webhook status
            max_retries: Maximum retry attempts
            timeout_seconds: Request timeout
            config: Request configuration
            
        Returns:
            Webhook: Updated webhook
        """
        request = UpdateWebhookRequest(
            name=name,
            url=url,
            events=events,
            secret=secret,
            status=status,
            max_retries=max_retries,
            timeout_seconds=timeout_seconds
        )
        
        response = self.http_client.put(
            f'/webhooks/{webhook_id}',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return Webhook(**response['data'])

    def delete(self, webhook_id: int, config: Optional[RequestConfig] = None) -> bool:
        """Delete a webhook.
        
        Args:
            webhook_id: Webhook ID
            config: Request configuration
            
        Returns:
            bool: True if deleted successfully
        """
        response = self.http_client.delete(f'/webhooks/{webhook_id}', config=config)
        return response.get('success', False)

    def test(
        self,
        webhook_id: int,
        event_type: Optional[WebhookEvent] = None,
        config: Optional[RequestConfig] = None
    ) -> WebhookTestResult:
        """Test a webhook endpoint.
        
        Args:
            webhook_id: Webhook ID
            event_type: Event type to test (uses default if not specified)
            config: Request configuration
            
        Returns:
            WebhookTestResult: Test result
        """
        data = {}
        if event_type:
            data['event_type'] = event_type
            
        response = self.http_client.post(
            f'/webhooks/{webhook_id}/test',
            data=data,
            config=config
        )
        
        return WebhookTestResult(**response['data'])

    def get_deliveries(
        self,
        webhook_id: int,
        pagination: Optional[PaginationParams] = None,
        filter_params: Optional[WebhookDeliveryFilter] = None,
        config: Optional[RequestConfig] = None
    ) -> PaginatedResponse[WebhookDelivery]:
        """Get webhook delivery history.
        
        Args:
            webhook_id: Webhook ID
            pagination: Pagination parameters
            filter_params: Filter parameters
            config: Request configuration
            
        Returns:
            PaginatedResponse[WebhookDelivery]: Paginated delivery history
        """
        params = {}
        
        if pagination:
            params.update(pagination.dict(exclude_none=True))
            
        if filter_params:
            params.update(filter_params.dict(exclude_none=True))
            
        response = self.http_client.get(
            f'/webhooks/{webhook_id}/deliveries',
            params=params,
            config=config
        )
        
        deliveries = [WebhookDelivery(**item) for item in response['data']]
        return PaginatedResponse[WebhookDelivery](
            data=deliveries,
            total=response['total'],
            limit=response['limit'],
            offset=response['offset']
        )

    def get_delivery(
        self,
        webhook_id: int,
        delivery_id: int,
        config: Optional[RequestConfig] = None
    ) -> WebhookDelivery:
        """Get a specific webhook delivery.
        
        Args:
            webhook_id: Webhook ID
            delivery_id: Delivery ID
            config: Request configuration
            
        Returns:
            WebhookDelivery: Delivery details
        """
        response = self.http_client.get(
            f'/webhooks/{webhook_id}/deliveries/{delivery_id}',
            config=config
        )
        
        return WebhookDelivery(**response['data'])

    def retry_delivery(
        self,
        webhook_id: int,
        delivery_id: int,
        config: Optional[RequestConfig] = None
    ) -> WebhookDelivery:
        """Retry a failed webhook delivery.
        
        Args:
            webhook_id: Webhook ID
            delivery_id: Delivery ID
            config: Request configuration
            
        Returns:
            WebhookDelivery: Retry result
        """
        response = self.http_client.post(
            f'/webhooks/{webhook_id}/deliveries/{delivery_id}/retry',
            config=config
        )
        
        return WebhookDelivery(**response['data'])

    def get_stats(
        self,
        webhook_id: int,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> WebhookStats:
        """Get webhook delivery statistics.
        
        Args:
            webhook_id: Webhook ID
            period: Time period for statistics
            config: Request configuration
            
        Returns:
            WebhookStats: Delivery statistics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            f'/webhooks/{webhook_id}/stats',
            params=params,
            config=config
        )
        
        return WebhookStats(**response['data'])

    def get_event_categories(
        self,
        config: Optional[RequestConfig] = None
    ) -> WebhookEventCategories:
        """Get available webhook event categories.
        
        Args:
            config: Request configuration
            
        Returns:
            WebhookEventCategories: Available event categories
        """
        response = self.http_client.get('/webhooks/events', config=config)
        return WebhookEventCategories(**response['data'])

    def get_health_summary(
        self,
        config: Optional[RequestConfig] = None
    ) -> WebhookHealthSummary:
        """Get webhook health summary.
        
        Args:
            config: Request configuration
            
        Returns:
            WebhookHealthSummary: Health summary
        """
        response = self.http_client.get('/webhooks/health', config=config)
        return WebhookHealthSummary(**response['data'])

    def get_templates(
        self,
        config: Optional[RequestConfig] = None
    ) -> List[WebhookTemplate]:
        """Get webhook templates for common use cases.
        
        Args:
            config: Request configuration
            
        Returns:
            List[WebhookTemplate]: Available templates
        """
        response = self.http_client.get('/webhooks/templates', config=config)
        return [WebhookTemplate(**item) for item in response['data']]

    def get_insights(
        self,
        webhook_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> WebhookInsights:
        """Get webhook insights and recommendations.
        
        Args:
            webhook_id: Specific webhook ID (or all webhooks if not specified)
            period: Time period for analysis
            config: Request configuration
            
        Returns:
            WebhookInsights: Insights and recommendations
        """
        params = {}
        if webhook_id:
            params['webhook_id'] = webhook_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/webhooks/insights',
            params=params,
            config=config
        )
        
        return WebhookInsights(**response['data'])

    def pause(self, webhook_id: int, config: Optional[RequestConfig] = None) -> Webhook:
        """Pause a webhook (set status to inactive).
        
        Args:
            webhook_id: Webhook ID
            config: Request configuration
            
        Returns:
            Webhook: Updated webhook
        """
        response = self.http_client.post(f'/webhooks/{webhook_id}/pause', config=config)
        return Webhook(**response['data'])

    def resume(self, webhook_id: int, config: Optional[RequestConfig] = None) -> Webhook:
        """Resume a webhook (set status to active).
        
        Args:
            webhook_id: Webhook ID
            config: Request configuration
            
        Returns:
            Webhook: Updated webhook
        """
        response = self.http_client.post(f'/webhooks/{webhook_id}/resume', config=config)
        return Webhook(**response['data'])

    def verify_signature(
        self,
        payload: str,
        signature: str,
        secret: str
    ) -> bool:
        """Verify webhook signature locally.
        
        Args:
            payload: Raw webhook payload
            signature: Webhook signature header
            secret: Webhook secret
            
        Returns:
            bool: True if signature is valid
        """
        import hmac
        import hashlib
        
        # Compute expected signature
        expected_signature = hmac.new(
            secret.encode('utf-8'),
            payload.encode('utf-8'),
            hashlib.sha256
        ).hexdigest()
        
        # Compare signatures
        expected_signature = f"sha256={expected_signature}"
        return hmac.compare_digest(expected_signature, signature)


class AsyncWebhookService:
    """Asynchronous webhook service."""

    def __init__(self, http_client: AsyncHTTPClient):
        """Initialize async webhook service.
        
        Args:
            http_client: Async HTTP client instance
        """
        self.http_client = http_client

    async def create(
        self,
        name: str,
        url: str,
        events: List[WebhookEvent],
        secret: Optional[str] = None,
        max_retries: Optional[int] = None,
        timeout_seconds: Optional[int] = None,
        config: Optional[RequestConfig] = None
    ) -> Webhook:
        """Create a new webhook.
        
        Args:
            name: Webhook name
            url: Webhook endpoint URL
            events: Events to subscribe to
            secret: Webhook secret for verification
            max_retries: Maximum retry attempts
            timeout_seconds: Request timeout
            config: Request configuration
            
        Returns:
            Webhook: Created webhook
        """
        request = CreateWebhookRequest(
            name=name,
            url=url,
            events=events,
            secret=secret,
            max_retries=max_retries,
            timeout_seconds=timeout_seconds
        )
        
        response = await self.http_client.post(
            '/webhooks',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return Webhook(**response['data'])

    async def get(self, webhook_id: int, config: Optional[RequestConfig] = None) -> Webhook:
        """Get a webhook by ID.
        
        Args:
            webhook_id: Webhook ID
            config: Request configuration
            
        Returns:
            Webhook: Webhook details
        """
        response = await self.http_client.get(f'/webhooks/{webhook_id}', config=config)
        return Webhook(**response['data'])

    async def list(
        self,
        pagination: Optional[PaginationParams] = None,
        filter_params: Optional[WebhookFilter] = None,
        config: Optional[RequestConfig] = None
    ) -> PaginatedResponse[Webhook]:
        """List user's webhooks.
        
        Args:
            pagination: Pagination parameters
            filter_params: Filter parameters
            config: Request configuration
            
        Returns:
            PaginatedResponse[Webhook]: Paginated list of webhooks
        """
        params = {}
        
        if pagination:
            params.update(pagination.dict(exclude_none=True))
            
        if filter_params:
            params.update(filter_params.dict(exclude_none=True))
            
        response = await self.http_client.get('/webhooks', params=params, config=config)
        
        webhooks = [Webhook(**item) for item in response['data']]
        return PaginatedResponse[Webhook](
            data=webhooks,
            total=response['total'],
            limit=response['limit'],
            offset=response['offset']
        )

    async def update(
        self,
        webhook_id: int,
        name: Optional[str] = None,
        url: Optional[str] = None,
        events: Optional[List[WebhookEvent]] = None,
        secret: Optional[str] = None,
        status: Optional[WebhookStatus] = None,
        max_retries: Optional[int] = None,
        timeout_seconds: Optional[int] = None,
        config: Optional[RequestConfig] = None
    ) -> Webhook:
        """Update a webhook.
        
        Args:
            webhook_id: Webhook ID
            name: Webhook name
            url: Webhook endpoint URL
            events: Events to subscribe to
            secret: Webhook secret
            status: Webhook status
            max_retries: Maximum retry attempts
            timeout_seconds: Request timeout
            config: Request configuration
            
        Returns:
            Webhook: Updated webhook
        """
        request = UpdateWebhookRequest(
            name=name,
            url=url,
            events=events,
            secret=secret,
            status=status,
            max_retries=max_retries,
            timeout_seconds=timeout_seconds
        )
        
        response = await self.http_client.put(
            f'/webhooks/{webhook_id}',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return Webhook(**response['data'])

    async def delete(self, webhook_id: int, config: Optional[RequestConfig] = None) -> bool:
        """Delete a webhook.
        
        Args:
            webhook_id: Webhook ID
            config: Request configuration
            
        Returns:
            bool: True if deleted successfully
        """
        response = await self.http_client.delete(f'/webhooks/{webhook_id}', config=config)
        return response.get('success', False)

    async def test(
        self,
        webhook_id: int,
        event_type: Optional[WebhookEvent] = None,
        config: Optional[RequestConfig] = None
    ) -> WebhookTestResult:
        """Test a webhook endpoint.
        
        Args:
            webhook_id: Webhook ID
            event_type: Event type to test (uses default if not specified)
            config: Request configuration
            
        Returns:
            WebhookTestResult: Test result
        """
        data = {}
        if event_type:
            data['event_type'] = event_type
            
        response = await self.http_client.post(
            f'/webhooks/{webhook_id}/test',
            data=data,
            config=config
        )
        
        return WebhookTestResult(**response['data'])

    async def get_stats(
        self,
        webhook_id: int,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> WebhookStats:
        """Get webhook delivery statistics.
        
        Args:
            webhook_id: Webhook ID
            period: Time period for statistics
            config: Request configuration
            
        Returns:
            WebhookStats: Delivery statistics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            f'/webhooks/{webhook_id}/stats',
            params=params,
            config=config
        )
        
        return WebhookStats(**response['data'])

    def verify_signature(
        self,
        payload: str,
        signature: str,
        secret: str
    ) -> bool:
        """Verify webhook signature locally.
        
        Args:
            payload: Raw webhook payload
            signature: Webhook signature header
            secret: Webhook secret
            
        Returns:
            bool: True if signature is valid
        """
        import hmac
        import hashlib
        
        # Compute expected signature
        expected_signature = hmac.new(
            secret.encode('utf-8'),
            payload.encode('utf-8'),
            hashlib.sha256
        ).hexdigest()
        
        # Compare signatures
        expected_signature = f"sha256={expected_signature}"
        return hmac.compare_digest(expected_signature, signature)