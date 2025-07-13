"""URL management service for URL Shortener API."""

from typing import List, Optional, Dict, Any

from ..http_client import HTTPClient, AsyncHTTPClient
from ..types.urls import (
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
from ..types.common import RequestConfig, PaginationParams, PaginatedResponse


class URLService:
    """Synchronous URL management service."""

    def __init__(self, http_client: HTTPClient):
        """Initialize URL service.
        
        Args:
            http_client: HTTP client instance
        """
        self.http_client = http_client

    def create(
        self,
        original_url: str,
        title: Optional[str] = None,
        description: Optional[str] = None,
        custom_alias: Optional[str] = None,
        password: Optional[str] = None,
        expires_at: Optional[str] = None,
        tags: Optional[List[str]] = None,
        config: Optional[RequestConfig] = None
    ) -> ShortURL:
        """Create a new short URL.
        
        Args:
            original_url: Original long URL
            title: URL title
            description: URL description
            custom_alias: Custom alias
            password: Password protection
            expires_at: Expiration date (ISO format)
            tags: URL tags
            config: Request configuration
            
        Returns:
            ShortURL: Created short URL
        """
        request = CreateURLRequest(
            original_url=original_url,
            title=title,
            description=description,
            custom_alias=custom_alias,
            password=password,
            expires_at=expires_at,
            tags=tags
        )
        
        response = self.http_client.post(
            '/urls',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return ShortURL(**response['data'])

    def get(self, url_id: int, config: Optional[RequestConfig] = None) -> ShortURL:
        """Get a short URL by ID.
        
        Args:
            url_id: URL ID
            config: Request configuration
            
        Returns:
            ShortURL: Short URL details
        """
        response = self.http_client.get(f'/urls/{url_id}', config=config)
        return ShortURL(**response['data'])

    def get_by_code(self, short_code: str, config: Optional[RequestConfig] = None) -> ShortURL:
        """Get a short URL by code.
        
        Args:
            short_code: Short URL code
            config: Request configuration
            
        Returns:
            ShortURL: Short URL details
        """
        response = self.http_client.get(f'/urls/code/{short_code}', config=config)
        return ShortURL(**response['data'])

    def list(
        self,
        pagination: Optional[PaginationParams] = None,
        tags: Optional[List[str]] = None,
        is_active: Optional[bool] = None,
        created_after: Optional[str] = None,
        created_before: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> PaginatedResponse[ShortURL]:
        """List user's short URLs.
        
        Args:
            pagination: Pagination parameters
            tags: Filter by tags
            is_active: Filter by active status
            created_after: Filter by creation date
            created_before: Filter by creation date
            config: Request configuration
            
        Returns:
            PaginatedResponse[ShortURL]: Paginated list of URLs
        """
        params = {}
        
        if pagination:
            params.update(pagination.dict(exclude_none=True))
            
        if tags:
            params['tags'] = ','.join(tags)
        if is_active is not None:
            params['is_active'] = is_active
        if created_after:
            params['created_after'] = created_after
        if created_before:
            params['created_before'] = created_before
            
        response = self.http_client.get('/urls', params=params, config=config)
        
        urls = [ShortURL(**item) for item in response['data']]
        return PaginatedResponse[ShortURL](
            data=urls,
            total=response['total'],
            limit=response['limit'],
            offset=response['offset']
        )

    def update(
        self,
        url_id: int,
        title: Optional[str] = None,
        description: Optional[str] = None,
        password: Optional[str] = None,
        expires_at: Optional[str] = None,
        is_active: Optional[bool] = None,
        tags: Optional[List[str]] = None,
        config: Optional[RequestConfig] = None
    ) -> ShortURL:
        """Update a short URL.
        
        Args:
            url_id: URL ID
            title: URL title
            description: URL description
            password: Password protection
            expires_at: Expiration date (ISO format)
            is_active: Whether URL is active
            tags: URL tags
            config: Request configuration
            
        Returns:
            ShortURL: Updated short URL
        """
        request = UpdateURLRequest(
            title=title,
            description=description,
            password=password,
            expires_at=expires_at,
            is_active=is_active,
            tags=tags
        )
        
        response = self.http_client.put(
            f'/urls/{url_id}',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return ShortURL(**response['data'])

    def delete(self, url_id: int, config: Optional[RequestConfig] = None) -> bool:
        """Delete a short URL.
        
        Args:
            url_id: URL ID
            config: Request configuration
            
        Returns:
            bool: True if deleted successfully
        """
        response = self.http_client.delete(f'/urls/{url_id}', config=config)
        return response.get('success', False)

    def bulk_create(
        self,
        urls: List[CreateURLRequest],
        config: Optional[RequestConfig] = None
    ) -> List[ShortURL]:
        """Create multiple URLs in bulk.
        
        Args:
            urls: List of URL creation requests
            config: Request configuration
            
        Returns:
            List[ShortURL]: List of created URLs
        """
        request = BulkCreateURLRequest(urls=urls)
        
        response = self.http_client.post(
            '/urls/bulk',
            data=request.dict(),
            config=config
        )
        
        return [ShortURL(**item) for item in response['data']]

    def bulk_delete(
        self,
        url_ids: List[int],
        config: Optional[RequestConfig] = None
    ) -> Dict[str, int]:
        """Delete multiple URLs in bulk.
        
        Args:
            url_ids: List of URL IDs to delete
            config: Request configuration
            
        Returns:
            Dict[str, int]: Deletion results with counts
        """
        request = BulkDeleteRequest(ids=url_ids)
        
        response = self.http_client.delete(
            '/urls/bulk',
            data=request.dict(),
            config=config
        )
        
        return response['data']

    def search(
        self,
        query: str,
        tags: Optional[List[str]] = None,
        is_active: Optional[bool] = None,
        created_after: Optional[str] = None,
        created_before: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> List[ShortURL]:
        """Search URLs.
        
        Args:
            query: Search query
            tags: Filter by tags
            is_active: Filter by active status
            created_after: Filter by creation date
            created_before: Filter by creation date
            config: Request configuration
            
        Returns:
            List[ShortURL]: Search results
        """
        request = URLSearchRequest(
            query=query,
            tags=tags,
            is_active=is_active,
            created_after=created_after,
            created_before=created_before
        )
        
        response = self.http_client.post(
            '/urls/search',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return [ShortURL(**item) for item in response['data']]

    def get_stats(self, url_id: int, config: Optional[RequestConfig] = None) -> URLStats:
        """Get URL statistics.
        
        Args:
            url_id: URL ID
            config: Request configuration
            
        Returns:
            URLStats: URL statistics
        """
        response = self.http_client.get(f'/urls/{url_id}/stats', config=config)
        return URLStats(**response['data'])

    def get_performance(
        self,
        url_id: int,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> URLPerformance:
        """Get URL performance metrics.
        
        Args:
            url_id: URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            URLPerformance: Performance metrics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            f'/urls/{url_id}/performance',
            params=params,
            config=config
        )
        
        return URLPerformance(**response['data'])

    def validate_password(
        self,
        url_id: int,
        password: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Validate URL password.
        
        Args:
            url_id: URL ID
            password: Password to validate
            config: Request configuration
            
        Returns:
            bool: True if password is valid
        """
        request = PasswordValidationRequest(password=password)
        
        response = self.http_client.post(
            f'/urls/{url_id}/validate-password',
            data=request.dict(),
            config=config
        )
        
        return response.get('valid', False)

    def get_redirect_info(
        self,
        short_code: str,
        password: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> URLRedirectResponse:
        """Get redirect information for a short URL.
        
        Args:
            short_code: Short URL code
            password: Password if URL is protected
            config: Request configuration
            
        Returns:
            URLRedirectResponse: Redirect information
        """
        params = {}
        if password:
            params['password'] = password
            
        response = self.http_client.get(
            f'/r/{short_code}/info',
            params=params,
            config=config
        )
        
        return URLRedirectResponse(**response['data'])

    def check_alias_availability(
        self,
        alias: str,
        config: Optional[RequestConfig] = None
    ) -> AliasAvailabilityResponse:
        """Check if custom alias is available.
        
        Args:
            alias: Custom alias to check
            config: Request configuration
            
        Returns:
            AliasAvailabilityResponse: Availability information
        """
        response = self.http_client.get(
            f'/urls/alias/{alias}/availability',
            config=config
        )
        
        return AliasAvailabilityResponse(**response['data'])


class AsyncURLService:
    """Asynchronous URL management service."""

    def __init__(self, http_client: AsyncHTTPClient):
        """Initialize async URL service.
        
        Args:
            http_client: Async HTTP client instance
        """
        self.http_client = http_client

    async def create(
        self,
        original_url: str,
        title: Optional[str] = None,
        description: Optional[str] = None,
        custom_alias: Optional[str] = None,
        password: Optional[str] = None,
        expires_at: Optional[str] = None,
        tags: Optional[List[str]] = None,
        config: Optional[RequestConfig] = None
    ) -> ShortURL:
        """Create a new short URL.
        
        Args:
            original_url: Original long URL
            title: URL title
            description: URL description
            custom_alias: Custom alias
            password: Password protection
            expires_at: Expiration date (ISO format)
            tags: URL tags
            config: Request configuration
            
        Returns:
            ShortURL: Created short URL
        """
        request = CreateURLRequest(
            original_url=original_url,
            title=title,
            description=description,
            custom_alias=custom_alias,
            password=password,
            expires_at=expires_at,
            tags=tags
        )
        
        response = await self.http_client.post(
            '/urls',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return ShortURL(**response['data'])

    async def get(self, url_id: int, config: Optional[RequestConfig] = None) -> ShortURL:
        """Get a short URL by ID.
        
        Args:
            url_id: URL ID
            config: Request configuration
            
        Returns:
            ShortURL: Short URL details
        """
        response = await self.http_client.get(f'/urls/{url_id}', config=config)
        return ShortURL(**response['data'])

    async def get_by_code(self, short_code: str, config: Optional[RequestConfig] = None) -> ShortURL:
        """Get a short URL by code.
        
        Args:
            short_code: Short URL code
            config: Request configuration
            
        Returns:
            ShortURL: Short URL details
        """
        response = await self.http_client.get(f'/urls/code/{short_code}', config=config)
        return ShortURL(**response['data'])

    async def list(
        self,
        pagination: Optional[PaginationParams] = None,
        tags: Optional[List[str]] = None,
        is_active: Optional[bool] = None,
        created_after: Optional[str] = None,
        created_before: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> PaginatedResponse[ShortURL]:
        """List user's short URLs.
        
        Args:
            pagination: Pagination parameters
            tags: Filter by tags
            is_active: Filter by active status
            created_after: Filter by creation date
            created_before: Filter by creation date
            config: Request configuration
            
        Returns:
            PaginatedResponse[ShortURL]: Paginated list of URLs
        """
        params = {}
        
        if pagination:
            params.update(pagination.dict(exclude_none=True))
            
        if tags:
            params['tags'] = ','.join(tags)
        if is_active is not None:
            params['is_active'] = is_active
        if created_after:
            params['created_after'] = created_after
        if created_before:
            params['created_before'] = created_before
            
        response = await self.http_client.get('/urls', params=params, config=config)
        
        urls = [ShortURL(**item) for item in response['data']]
        return PaginatedResponse[ShortURL](
            data=urls,
            total=response['total'],
            limit=response['limit'],
            offset=response['offset']
        )

    async def update(
        self,
        url_id: int,
        title: Optional[str] = None,
        description: Optional[str] = None,
        password: Optional[str] = None,
        expires_at: Optional[str] = None,
        is_active: Optional[bool] = None,
        tags: Optional[List[str]] = None,
        config: Optional[RequestConfig] = None
    ) -> ShortURL:
        """Update a short URL.
        
        Args:
            url_id: URL ID
            title: URL title
            description: URL description
            password: Password protection
            expires_at: Expiration date (ISO format)
            is_active: Whether URL is active
            tags: URL tags
            config: Request configuration
            
        Returns:
            ShortURL: Updated short URL
        """
        request = UpdateURLRequest(
            title=title,
            description=description,
            password=password,
            expires_at=expires_at,
            is_active=is_active,
            tags=tags
        )
        
        response = await self.http_client.put(
            f'/urls/{url_id}',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return ShortURL(**response['data'])

    async def delete(self, url_id: int, config: Optional[RequestConfig] = None) -> bool:
        """Delete a short URL.
        
        Args:
            url_id: URL ID
            config: Request configuration
            
        Returns:
            bool: True if deleted successfully
        """
        response = await self.http_client.delete(f'/urls/{url_id}', config=config)
        return response.get('success', False)

    async def bulk_create(
        self,
        urls: List[CreateURLRequest],
        config: Optional[RequestConfig] = None
    ) -> List[ShortURL]:
        """Create multiple URLs in bulk.
        
        Args:
            urls: List of URL creation requests
            config: Request configuration
            
        Returns:
            List[ShortURL]: List of created URLs
        """
        request = BulkCreateURLRequest(urls=urls)
        
        response = await self.http_client.post(
            '/urls/bulk',
            data=request.dict(),
            config=config
        )
        
        return [ShortURL(**item) for item in response['data']]

    async def bulk_delete(
        self,
        url_ids: List[int],
        config: Optional[RequestConfig] = None
    ) -> Dict[str, int]:
        """Delete multiple URLs in bulk.
        
        Args:
            url_ids: List of URL IDs to delete
            config: Request configuration
            
        Returns:
            Dict[str, int]: Deletion results with counts
        """
        request = BulkDeleteRequest(ids=url_ids)
        
        response = await self.http_client.delete(
            '/urls/bulk',
            data=request.dict(),
            config=config
        )
        
        return response['data']

    async def search(
        self,
        query: str,
        tags: Optional[List[str]] = None,
        is_active: Optional[bool] = None,
        created_after: Optional[str] = None,
        created_before: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> List[ShortURL]:
        """Search URLs.
        
        Args:
            query: Search query
            tags: Filter by tags
            is_active: Filter by active status
            created_after: Filter by creation date
            created_before: Filter by creation date
            config: Request configuration
            
        Returns:
            List[ShortURL]: Search results
        """
        request = URLSearchRequest(
            query=query,
            tags=tags,
            is_active=is_active,
            created_after=created_after,
            created_before=created_before
        )
        
        response = await self.http_client.post(
            '/urls/search',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return [ShortURL(**item) for item in response['data']]

    async def get_stats(self, url_id: int, config: Optional[RequestConfig] = None) -> URLStats:
        """Get URL statistics.
        
        Args:
            url_id: URL ID
            config: Request configuration
            
        Returns:
            URLStats: URL statistics
        """
        response = await self.http_client.get(f'/urls/{url_id}/stats', config=config)
        return URLStats(**response['data'])

    async def get_performance(
        self,
        url_id: int,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> URLPerformance:
        """Get URL performance metrics.
        
        Args:
            url_id: URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            URLPerformance: Performance metrics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            f'/urls/{url_id}/performance',
            params=params,
            config=config
        )
        
        return URLPerformance(**response['data'])

    async def validate_password(
        self,
        url_id: int,
        password: str,
        config: Optional[RequestConfig] = None
    ) -> bool:
        """Validate URL password.
        
        Args:
            url_id: URL ID
            password: Password to validate
            config: Request configuration
            
        Returns:
            bool: True if password is valid
        """
        request = PasswordValidationRequest(password=password)
        
        response = await self.http_client.post(
            f'/urls/{url_id}/validate-password',
            data=request.dict(),
            config=config
        )
        
        return response.get('valid', False)

    async def get_redirect_info(
        self,
        short_code: str,
        password: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> URLRedirectResponse:
        """Get redirect information for a short URL.
        
        Args:
            short_code: Short URL code
            password: Password if URL is protected
            config: Request configuration
            
        Returns:
            URLRedirectResponse: Redirect information
        """
        params = {}
        if password:
            params['password'] = password
            
        response = await self.http_client.get(
            f'/r/{short_code}/info',
            params=params,
            config=config
        )
        
        return URLRedirectResponse(**response['data'])

    async def check_alias_availability(
        self,
        alias: str,
        config: Optional[RequestConfig] = None
    ) -> AliasAvailabilityResponse:
        """Check if custom alias is available.
        
        Args:
            alias: Custom alias to check
            config: Request configuration
            
        Returns:
            AliasAvailabilityResponse: Availability information
        """
        response = await self.http_client.get(
            f'/urls/alias/{alias}/availability',
            config=config
        )
        
        return AliasAvailabilityResponse(**response['data'])