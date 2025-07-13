"""Main client for the URL Shortener Python SDK."""

from typing import Optional, Dict, Any, Union
from datetime import datetime

from .http_client import HTTPClient, AsyncHTTPClient
from .services.auth import AuthService, AsyncAuthService
from .services.urls import URLService, AsyncURLService
from .services.analytics import AnalyticsService, AsyncAnalyticsService
from .services.qr import QRService, AsyncQRService
from .services.webhooks import WebhookService, AsyncWebhookService
from .types.common import RequestConfig, HealthStatus
from .exceptions import ConfigurationError


class URLShortenerClient:
    """Synchronous URL Shortener API client."""

    def __init__(
        self,
        base_url: str,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        timeout: int = 30,
        max_retries: int = 3,
        retry_delay: float = 1.0,
        debug: bool = False
    ):
        """Initialize URL Shortener client.
        
        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
            access_token: Access token for authentication
            timeout: Request timeout in seconds
            max_retries: Maximum number of retry attempts
            retry_delay: Initial delay between retries
            debug: Enable debug logging
            
        Raises:
            ConfigurationError: If configuration is invalid
        """
        if not base_url:
            raise ConfigurationError("base_url is required")
            
        if not api_key and not access_token:
            raise ConfigurationError("Either api_key or access_token is required")
        
        # Initialize HTTP client
        self.http_client = HTTPClient(
            base_url=base_url,
            api_key=api_key,
            access_token=access_token,
            timeout=timeout,
            max_retries=max_retries,
            retry_delay=retry_delay,
            debug=debug
        )
        
        # Initialize services
        self.auth = AuthService(self.http_client)
        self.urls = URLService(self.http_client)
        self.analytics = AnalyticsService(self.http_client)
        self.qr = QRService(self.http_client)
        self.webhooks = WebhookService(self.http_client)
        
    @classmethod
    def with_api_key(
        cls,
        base_url: str,
        api_key: str,
        **kwargs
    ) -> 'URLShortenerClient':
        """Create client with API key authentication.
        
        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
            **kwargs: Additional client options
            
        Returns:
            URLShortenerClient: Configured client instance
        """
        return cls(base_url=base_url, api_key=api_key, **kwargs)
        
    @classmethod
    def with_access_token(
        cls,
        base_url: str,
        access_token: str,
        **kwargs
    ) -> 'URLShortenerClient':
        """Create client with access token authentication.
        
        Args:
            base_url: Base URL for the API
            access_token: Access token for authentication
            **kwargs: Additional client options
            
        Returns:
            URLShortenerClient: Configured client instance
        """
        return cls(base_url=base_url, access_token=access_token, **kwargs)
        
    @classmethod
    def development(
        cls,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        **kwargs
    ) -> 'URLShortenerClient':
        """Create client for development environment.
        
        Args:
            api_key: API key for authentication
            access_token: Access token for authentication
            **kwargs: Additional client options
            
        Returns:
            URLShortenerClient: Development client instance
        """
        base_url = kwargs.pop('base_url', 'http://localhost:8080/api/v1')
        debug = kwargs.pop('debug', True)
        
        return cls(
            base_url=base_url,
            api_key=api_key,
            access_token=access_token,
            debug=debug,
            **kwargs
        )
        
    @classmethod
    def production(
        cls,
        base_url: str,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        **kwargs
    ) -> 'URLShortenerClient':
        """Create client for production environment.
        
        Args:
            base_url: Production base URL
            api_key: API key for authentication
            access_token: Access token for authentication
            **kwargs: Additional client options
            
        Returns:
            URLShortenerClient: Production client instance
        """
        debug = kwargs.pop('debug', False)
        timeout = kwargs.pop('timeout', 30)
        max_retries = kwargs.pop('max_retries', 5)
        
        return cls(
            base_url=base_url,
            api_key=api_key,
            access_token=access_token,
            debug=debug,
            timeout=timeout,
            max_retries=max_retries,
            **kwargs
        )
        
    def set_access_token(self, token: str):
        """Set access token for authentication.
        
        Args:
            token: Access token
        """
        self.http_client.set_access_token(token)
        
    def set_token_refresh_callback(self, callback: callable, refresh_token: str):
        """Set callback for automatic token refresh.
        
        Args:
            callback: Function to call for token refresh
            refresh_token: Refresh token to use
        """
        self.http_client.set_token_refresh_callback(callback, refresh_token)
        
    def health_check(self, config: Optional[RequestConfig] = None) -> HealthStatus:
        """Check API health status.
        
        Args:
            config: Request configuration
            
        Returns:
            HealthStatus: API health information
        """
        response = self.http_client.get('/health', config=config)
        return HealthStatus(**response)
        
    def close(self):
        """Close the HTTP client and clean up resources."""
        self.http_client.close()
        
    def __enter__(self):
        """Context manager entry."""
        return self
        
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit."""
        self.close()


class AsyncURLShortenerClient:
    """Asynchronous URL Shortener API client."""

    def __init__(
        self,
        base_url: str,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        timeout: int = 30,
        max_retries: int = 3,
        retry_delay: float = 1.0,
        debug: bool = False
    ):
        """Initialize async URL Shortener client.
        
        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
            access_token: Access token for authentication
            timeout: Request timeout in seconds
            max_retries: Maximum number of retry attempts
            retry_delay: Initial delay between retries
            debug: Enable debug logging
            
        Raises:
            ConfigurationError: If configuration is invalid
        """
        if not base_url:
            raise ConfigurationError("base_url is required")
            
        if not api_key and not access_token:
            raise ConfigurationError("Either api_key or access_token is required")
        
        # Initialize async HTTP client
        self.http_client = AsyncHTTPClient(
            base_url=base_url,
            api_key=api_key,
            access_token=access_token,
            timeout=timeout,
            max_retries=max_retries,
            retry_delay=retry_delay,
            debug=debug
        )
        
        # Initialize async services
        self.auth = AsyncAuthService(self.http_client)
        self.urls = AsyncURLService(self.http_client)
        self.analytics = AsyncAnalyticsService(self.http_client)
        self.qr = AsyncQRService(self.http_client)
        self.webhooks = AsyncWebhookService(self.http_client)
        
    @classmethod
    def with_api_key(
        cls,
        base_url: str,
        api_key: str,
        **kwargs
    ) -> 'AsyncURLShortenerClient':
        """Create async client with API key authentication.
        
        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
            **kwargs: Additional client options
            
        Returns:
            AsyncURLShortenerClient: Configured async client instance
        """
        return cls(base_url=base_url, api_key=api_key, **kwargs)
        
    @classmethod
    def with_access_token(
        cls,
        base_url: str,
        access_token: str,
        **kwargs
    ) -> 'AsyncURLShortenerClient':
        """Create async client with access token authentication.
        
        Args:
            base_url: Base URL for the API
            access_token: Access token for authentication
            **kwargs: Additional client options
            
        Returns:
            AsyncURLShortenerClient: Configured async client instance
        """
        return cls(base_url=base_url, access_token=access_token, **kwargs)
        
    @classmethod
    def development(
        cls,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        **kwargs
    ) -> 'AsyncURLShortenerClient':
        """Create async client for development environment.
        
        Args:
            api_key: API key for authentication
            access_token: Access token for authentication
            **kwargs: Additional client options
            
        Returns:
            AsyncURLShortenerClient: Development async client instance
        """
        base_url = kwargs.pop('base_url', 'http://localhost:8080/api/v1')
        debug = kwargs.pop('debug', True)
        
        return cls(
            base_url=base_url,
            api_key=api_key,
            access_token=access_token,
            debug=debug,
            **kwargs
        )
        
    @classmethod
    def production(
        cls,
        base_url: str,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        **kwargs
    ) -> 'AsyncURLShortenerClient':
        """Create async client for production environment.
        
        Args:
            base_url: Production base URL
            api_key: API key for authentication
            access_token: Access token for authentication
            **kwargs: Additional client options
            
        Returns:
            AsyncURLShortenerClient: Production async client instance
        """
        debug = kwargs.pop('debug', False)
        timeout = kwargs.pop('timeout', 30)
        max_retries = kwargs.pop('max_retries', 5)
        
        return cls(
            base_url=base_url,
            api_key=api_key,
            access_token=access_token,
            debug=debug,
            timeout=timeout,
            max_retries=max_retries,
            **kwargs
        )
        
    def set_access_token(self, token: str):
        """Set access token for authentication.
        
        Args:
            token: Access token
        """
        self.http_client.set_access_token(token)
        
    def set_token_refresh_callback(self, callback: callable, refresh_token: str):
        """Set callback for automatic token refresh.
        
        Args:
            callback: Function to call for token refresh
            refresh_token: Refresh token to use
        """
        self.http_client.set_token_refresh_callback(callback, refresh_token)
        
    async def health_check(self, config: Optional[RequestConfig] = None) -> HealthStatus:
        """Check API health status.
        
        Args:
            config: Request configuration
            
        Returns:
            HealthStatus: API health information
        """
        response = await self.http_client.get('/health', config=config)
        return HealthStatus(**response)
        
    async def close(self):
        """Close the async HTTP client and clean up resources."""
        await self.http_client.close()
        
    async def __aenter__(self):
        """Async context manager entry."""
        return self
        
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        await self.close()


# Convenience aliases
Client = URLShortenerClient
AsyncClient = AsyncURLShortenerClient