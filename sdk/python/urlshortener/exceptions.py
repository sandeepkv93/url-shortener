"""Custom exceptions for the URL Shortener SDK."""

from typing import List, Optional, Dict, Any


class URLShortenerError(Exception):
    """Base exception for all URL Shortener SDK errors."""
    
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        """Initialize exception with message and optional details.
        
        Args:
            message: Error message
            details: Additional error details
        """
        super().__init__(message)
        self.message = message
        self.details = details or {}


class AuthenticationError(URLShortenerError):
    """Raised when authentication fails."""
    pass


class AuthorizationError(URLShortenerError):
    """Raised when authorization fails."""
    pass


class ValidationError(URLShortenerError):
    """Raised when request validation fails."""
    
    def __init__(self, message: str, errors: Optional[List[str]] = None, details: Optional[Dict[str, Any]] = None):
        """Initialize validation error.
        
        Args:
            message: Error message
            errors: List of validation error messages
            details: Additional error details
        """
        super().__init__(message, details)
        self.errors = errors or []


class NotFoundError(URLShortenerError):
    """Raised when a resource is not found."""
    pass


class ConflictError(URLShortenerError):
    """Raised when a resource conflict occurs."""
    pass


class RateLimitError(URLShortenerError):
    """Raised when rate limit is exceeded."""
    
    def __init__(self, message: str, retry_after: Optional[int] = None, details: Optional[Dict[str, Any]] = None):
        """Initialize rate limit error.
        
        Args:
            message: Error message
            retry_after: Seconds to wait before retrying
            details: Additional error details
        """
        super().__init__(message, details)
        self.retry_after = retry_after


class ServerError(URLShortenerError):
    """Raised when a server error occurs."""
    pass


class NetworkError(URLShortenerError):
    """Raised when a network error occurs."""
    pass


class TimeoutError(URLShortenerError):
    """Raised when a request times out."""
    pass


class ConfigurationError(URLShortenerError):
    """Raised when SDK configuration is invalid."""
    pass


class TokenExpiredError(AuthenticationError):
    """Raised when access token has expired."""
    pass


class RefreshTokenExpiredError(AuthenticationError):
    """Raised when refresh token has expired."""
    pass


class InvalidTokenError(AuthenticationError):
    """Raised when token is invalid or malformed."""
    pass


class QuotaExceededError(URLShortenerError):
    """Raised when API quota is exceeded."""
    
    def __init__(self, message: str, quota_limit: Optional[int] = None, quota_remaining: Optional[int] = None, details: Optional[Dict[str, Any]] = None):
        """Initialize quota exceeded error.
        
        Args:
            message: Error message
            quota_limit: Total quota limit
            quota_remaining: Remaining quota
            details: Additional error details
        """
        super().__init__(message, details)
        self.quota_limit = quota_limit
        self.quota_remaining = quota_remaining


class WebhookError(URLShortenerError):
    """Raised when webhook operations fail."""
    pass


class QRCodeError(URLShortenerError):
    """Raised when QR code operations fail."""
    pass


class AnalyticsError(URLShortenerError):
    """Raised when analytics operations fail."""
    pass


class URLError(URLShortenerError):
    """Raised when URL operations fail."""
    pass


class InvalidURLError(URLError):
    """Raised when provided URL is invalid."""
    pass


class URLExpiredError(URLError):
    """Raised when trying to access an expired URL."""
    pass


class URLPasswordProtectedError(URLError):
    """Raised when URL requires password authentication."""
    pass


class CustomAliasUnavailableError(URLError):
    """Raised when custom alias is not available."""
    
    def __init__(self, message: str, suggestions: Optional[List[str]] = None, details: Optional[Dict[str, Any]] = None):
        """Initialize custom alias error.
        
        Args:
            message: Error message
            suggestions: Alternative alias suggestions
            details: Additional error details
        """
        super().__init__(message, details)
        self.suggestions = suggestions or []


class BulkOperationError(URLShortenerError):
    """Raised when bulk operations partially fail."""
    
    def __init__(self, message: str, successful_count: int = 0, failed_count: int = 0, failures: Optional[List[Dict[str, Any]]] = None, details: Optional[Dict[str, Any]] = None):
        """Initialize bulk operation error.
        
        Args:
            message: Error message
            successful_count: Number of successful operations
            failed_count: Number of failed operations
            failures: List of failure details
            details: Additional error details
        """
        super().__init__(message, details)
        self.successful_count = successful_count
        self.failed_count = failed_count
        self.failures = failures or []


class SDKError(URLShortenerError):
    """Raised when SDK internal errors occur."""
    pass


def handle_api_error(status_code: int, response_data: Dict[str, Any]) -> URLShortenerError:
    """Create appropriate exception based on API response.
    
    Args:
        status_code: HTTP status code
        response_data: Response data from API
        
    Returns:
        Appropriate exception instance
    """
    message = response_data.get('message', f'API error: {status_code}')
    details = response_data.get('details', {})
    
    if status_code == 400:
        if 'validation' in message.lower():
            return ValidationError(message, response_data.get('errors', []), details)
        return URLShortenerError(message, details)
    elif status_code == 401:
        if 'token' in message.lower() and 'expired' in message.lower():
            return TokenExpiredError(message, details)
        elif 'token' in message.lower() and 'invalid' in message.lower():
            return InvalidTokenError(message, details)
        return AuthenticationError(message, details)
    elif status_code == 403:
        return AuthorizationError(message, details)
    elif status_code == 404:
        return NotFoundError(message, details)
    elif status_code == 409:
        if 'alias' in message.lower():
            return CustomAliasUnavailableError(
                message, 
                response_data.get('suggestions', []), 
                details
            )
        return ConflictError(message, details)
    elif status_code == 422:
        return ValidationError(message, response_data.get('errors', []), details)
    elif status_code == 429:
        return RateLimitError(
            message, 
            response_data.get('retry_after'), 
            details
        )
    elif 500 <= status_code < 600:
        return ServerError(message, details)
    else:
        return URLShortenerError(message, details)