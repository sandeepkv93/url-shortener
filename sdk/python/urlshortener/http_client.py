"""HTTP client for the URL Shortener API."""

import asyncio
import json
import time
from typing import Any, Dict, Optional, Union
from urllib.parse import urljoin

import httpx
from httpx import AsyncClient, Client, Response

from .types.common import APIError, APIResponse, RequestConfig
from .exceptions import (
    URLShortenerError,
    AuthenticationError,
    RateLimitError,
    ValidationError,
    ServerError,
    NetworkError
)


class HTTPClient:
    """Synchronous HTTP client for URL Shortener API."""

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
        """Initialize HTTP client.
        
        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
            access_token: Access token for authentication
            timeout: Request timeout in seconds
            max_retries: Maximum number of retry attempts
            retry_delay: Initial delay between retries
            debug: Enable debug logging
        """
        self.base_url = base_url.rstrip('/')
        self.api_key = api_key
        self.access_token = access_token
        self.timeout = timeout
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.debug = debug
        
        # Token refresh callback
        self._token_refresh_callback: Optional[callable] = None
        self._refresh_token: Optional[str] = None
        
        # Initialize HTTP client
        self._client = Client(
            timeout=timeout,
            follow_redirects=True
        )
        
    def set_token_refresh_callback(self, callback: callable, refresh_token: str):
        """Set callback for automatic token refresh.
        
        Args:
            callback: Function to call for token refresh
            refresh_token: Refresh token to use
        """
        self._token_refresh_callback = callback
        self._refresh_token = refresh_token
        
    def set_access_token(self, token: str):
        """Set access token for requests."""
        self.access_token = token
        
    def _get_headers(self, additional_headers: Optional[Dict[str, str]] = None) -> Dict[str, str]:
        """Get request headers with authentication."""
        headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'User-Agent': 'URLShortener-Python-SDK/1.0.0'
        }
        
        # Add authentication
        if self.access_token:
            headers['Authorization'] = f'Bearer {self.access_token}'
        elif self.api_key:
            headers['X-API-Key'] = self.api_key
            
        # Add additional headers
        if additional_headers:
            headers.update(additional_headers)
            
        return headers
        
    def _handle_response(self, response: Response) -> Dict[str, Any]:
        """Handle HTTP response and convert to dictionary."""
        if self.debug:
            print(f"Response: {response.status_code} {response.text}")
            
        # Handle successful responses
        if 200 <= response.status_code < 300:
            try:
                return response.json()
            except json.JSONDecodeError:
                return {'data': response.text}
                
        # Handle error responses
        error_data = {}
        try:
            error_data = response.json()
        except json.JSONDecodeError:
            error_data = {'message': response.text}
            
        # Raise appropriate exceptions
        if response.status_code == 401:
            raise AuthenticationError(error_data.get('message', 'Authentication failed'))
        elif response.status_code == 403:
            raise AuthenticationError(error_data.get('message', 'Access forbidden'))
        elif response.status_code == 422:
            raise ValidationError(error_data.get('message', 'Validation failed'), error_data.get('errors', []))
        elif response.status_code == 429:
            raise RateLimitError(error_data.get('message', 'Rate limit exceeded'))
        elif 400 <= response.status_code < 500:
            raise URLShortenerError(error_data.get('message', f'Client error: {response.status_code}'))
        elif 500 <= response.status_code < 600:
            raise ServerError(error_data.get('message', f'Server error: {response.status_code}'))
        else:
            raise URLShortenerError(f'Unknown error: {response.status_code}')
            
    def _should_retry(self, exception: Exception, attempt: int) -> bool:
        """Determine if request should be retried."""
        if attempt >= self.max_retries:
            return False
            
        # Retry on network errors and server errors
        if isinstance(exception, (NetworkError, ServerError)):
            return True
            
        # Retry on rate limits with backoff
        if isinstance(exception, RateLimitError):
            return True
            
        return False
        
    def _get_retry_delay(self, attempt: int) -> float:
        """Calculate retry delay with exponential backoff."""
        return self.retry_delay * (2 ** attempt)
        
    def _attempt_token_refresh(self) -> bool:
        """Attempt to refresh access token."""
        if not self._token_refresh_callback or not self._refresh_token:
            return False
            
        try:
            tokens = self._token_refresh_callback(self._refresh_token)
            self.set_access_token(tokens.access_token)
            self._refresh_token = tokens.refresh_token
            return True
        except Exception:
            return False
            
    def request(
        self,
        method: str,
        endpoint: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        config: Optional[RequestConfig] = None
    ) -> Dict[str, Any]:
        """Make HTTP request with retries and error handling."""
        url = urljoin(self.base_url, endpoint.lstrip('/'))
        request_headers = self._get_headers(headers)
        
        # Apply config overrides
        timeout = self.timeout
        if config and config.timeout:
            timeout = config.timeout
            
        last_exception = None
        
        for attempt in range(self.max_retries + 1):
            try:
                if self.debug:
                    print(f"Request: {method} {url} (attempt {attempt + 1})")
                    
                # Make request
                response = self._client.request(
                    method=method,
                    url=url,
                    json=data,
                    params=params,
                    headers=request_headers,
                    timeout=timeout
                )
                
                return self._handle_response(response)
                
            except AuthenticationError as e:
                # Try token refresh on first auth error
                if attempt == 0 and self._attempt_token_refresh():
                    request_headers = self._get_headers(headers)
                    continue
                raise e
                
            except Exception as e:
                last_exception = e
                
                # Convert network errors
                if isinstance(e, httpx.NetworkError):
                    last_exception = NetworkError(f"Network error: {str(e)}")
                elif isinstance(e, httpx.TimeoutException):
                    last_exception = NetworkError(f"Request timeout: {str(e)}")
                    
                # Check if we should retry
                if not self._should_retry(last_exception, attempt):
                    break
                    
                # Wait before retry
                if attempt < self.max_retries:
                    delay = self._get_retry_delay(attempt)
                    time.sleep(delay)
                    
        # Raise the last exception if all retries failed
        if last_exception:
            raise last_exception
        else:
            raise URLShortenerError("Request failed after all retries")
            
    def get(self, endpoint: str, params: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make GET request."""
        return self.request('GET', endpoint, params=params, **kwargs)
        
    def post(self, endpoint: str, data: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make POST request."""
        return self.request('POST', endpoint, data=data, **kwargs)
        
    def put(self, endpoint: str, data: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make PUT request."""
        return self.request('PUT', endpoint, data=data, **kwargs)
        
    def delete(self, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make DELETE request."""
        return self.request('DELETE', endpoint, **kwargs)
        
    def patch(self, endpoint: str, data: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make PATCH request."""
        return self.request('PATCH', endpoint, data=data, **kwargs)
        
    def close(self):
        """Close the HTTP client."""
        self._client.close()
        
    def __enter__(self):
        """Context manager entry."""
        return self
        
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit."""
        self.close()


class AsyncHTTPClient:
    """Asynchronous HTTP client for URL Shortener API."""

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
        """Initialize async HTTP client.
        
        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
            access_token: Access token for authentication
            timeout: Request timeout in seconds
            max_retries: Maximum number of retry attempts
            retry_delay: Initial delay between retries
            debug: Enable debug logging
        """
        self.base_url = base_url.rstrip('/')
        self.api_key = api_key
        self.access_token = access_token
        self.timeout = timeout
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.debug = debug
        
        # Token refresh callback
        self._token_refresh_callback: Optional[callable] = None
        self._refresh_token: Optional[str] = None
        
        # Initialize async HTTP client
        self._client = AsyncClient(
            timeout=timeout,
            follow_redirects=True
        )
        
    def set_token_refresh_callback(self, callback: callable, refresh_token: str):
        """Set callback for automatic token refresh."""
        self._token_refresh_callback = callback
        self._refresh_token = refresh_token
        
    def set_access_token(self, token: str):
        """Set access token for requests."""
        self.access_token = token
        
    def _get_headers(self, additional_headers: Optional[Dict[str, str]] = None) -> Dict[str, str]:
        """Get request headers with authentication."""
        headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'User-Agent': 'URLShortener-Python-SDK/1.0.0'
        }
        
        # Add authentication
        if self.access_token:
            headers['Authorization'] = f'Bearer {self.access_token}'
        elif self.api_key:
            headers['X-API-Key'] = self.api_key
            
        # Add additional headers
        if additional_headers:
            headers.update(additional_headers)
            
        return headers
        
    def _handle_response(self, response: Response) -> Dict[str, Any]:
        """Handle HTTP response and convert to dictionary."""
        if self.debug:
            print(f"Response: {response.status_code} {response.text}")
            
        # Handle successful responses
        if 200 <= response.status_code < 300:
            try:
                return response.json()
            except json.JSONDecodeError:
                return {'data': response.text}
                
        # Handle error responses
        error_data = {}
        try:
            error_data = response.json()
        except json.JSONDecodeError:
            error_data = {'message': response.text}
            
        # Raise appropriate exceptions
        if response.status_code == 401:
            raise AuthenticationError(error_data.get('message', 'Authentication failed'))
        elif response.status_code == 403:
            raise AuthenticationError(error_data.get('message', 'Access forbidden'))
        elif response.status_code == 422:
            raise ValidationError(error_data.get('message', 'Validation failed'), error_data.get('errors', []))
        elif response.status_code == 429:
            raise RateLimitError(error_data.get('message', 'Rate limit exceeded'))
        elif 400 <= response.status_code < 500:
            raise URLShortenerError(error_data.get('message', f'Client error: {response.status_code}'))
        elif 500 <= response.status_code < 600:
            raise ServerError(error_data.get('message', f'Server error: {response.status_code}'))
        else:
            raise URLShortenerError(f'Unknown error: {response.status_code}')
            
    def _should_retry(self, exception: Exception, attempt: int) -> bool:
        """Determine if request should be retried."""
        if attempt >= self.max_retries:
            return False
            
        # Retry on network errors and server errors
        if isinstance(exception, (NetworkError, ServerError)):
            return True
            
        # Retry on rate limits with backoff
        if isinstance(exception, RateLimitError):
            return True
            
        return False
        
    def _get_retry_delay(self, attempt: int) -> float:
        """Calculate retry delay with exponential backoff."""
        return self.retry_delay * (2 ** attempt)
        
    async def _attempt_token_refresh(self) -> bool:
        """Attempt to refresh access token."""
        if not self._token_refresh_callback or not self._refresh_token:
            return False
            
        try:
            tokens = await self._token_refresh_callback(self._refresh_token)
            self.set_access_token(tokens.access_token)
            self._refresh_token = tokens.refresh_token
            return True
        except Exception:
            return False
            
    async def request(
        self,
        method: str,
        endpoint: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        config: Optional[RequestConfig] = None
    ) -> Dict[str, Any]:
        """Make async HTTP request with retries and error handling."""
        url = urljoin(self.base_url, endpoint.lstrip('/'))
        request_headers = self._get_headers(headers)
        
        # Apply config overrides
        timeout = self.timeout
        if config and config.timeout:
            timeout = config.timeout
            
        last_exception = None
        
        for attempt in range(self.max_retries + 1):
            try:
                if self.debug:
                    print(f"Request: {method} {url} (attempt {attempt + 1})")
                    
                # Make request
                response = await self._client.request(
                    method=method,
                    url=url,
                    json=data,
                    params=params,
                    headers=request_headers,
                    timeout=timeout
                )
                
                return self._handle_response(response)
                
            except AuthenticationError as e:
                # Try token refresh on first auth error
                if attempt == 0 and await self._attempt_token_refresh():
                    request_headers = self._get_headers(headers)
                    continue
                raise e
                
            except Exception as e:
                last_exception = e
                
                # Convert network errors
                if isinstance(e, httpx.NetworkError):
                    last_exception = NetworkError(f"Network error: {str(e)}")
                elif isinstance(e, httpx.TimeoutException):
                    last_exception = NetworkError(f"Request timeout: {str(e)}")
                    
                # Check if we should retry
                if not self._should_retry(last_exception, attempt):
                    break
                    
                # Wait before retry
                if attempt < self.max_retries:
                    delay = self._get_retry_delay(attempt)
                    await asyncio.sleep(delay)
                    
        # Raise the last exception if all retries failed
        if last_exception:
            raise last_exception
        else:
            raise URLShortenerError("Request failed after all retries")
            
    async def get(self, endpoint: str, params: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make async GET request."""
        return await self.request('GET', endpoint, params=params, **kwargs)
        
    async def post(self, endpoint: str, data: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make async POST request."""
        return await self.request('POST', endpoint, data=data, **kwargs)
        
    async def put(self, endpoint: str, data: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make async PUT request."""
        return await self.request('PUT', endpoint, data=data, **kwargs)
        
    async def delete(self, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make async DELETE request."""
        return await self.request('DELETE', endpoint, **kwargs)
        
    async def patch(self, endpoint: str, data: Optional[Dict[str, Any]] = None, **kwargs) -> Dict[str, Any]:
        """Make async PATCH request."""
        return await self.request('PATCH', endpoint, data=data, **kwargs)
        
    async def close(self):
        """Close the async HTTP client."""
        await self._client.aclose()
        
    async def __aenter__(self):
        """Async context manager entry."""
        return self
        
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        await self.close()