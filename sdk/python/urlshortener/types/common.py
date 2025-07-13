"""Common type definitions."""

from typing import Any, Dict, Generic, List, Optional, TypeVar, Union
from datetime import datetime
from pydantic import BaseModel, Field

T = TypeVar('T')


class PaginationParams(BaseModel):
    """Pagination parameters for list requests."""
    limit: Optional[int] = Field(None, ge=1, le=1000, description="Maximum number of items to return")
    offset: Optional[int] = Field(None, ge=0, description="Number of items to skip")


class PaginatedResponse(BaseModel, Generic[T]):
    """Generic paginated response wrapper."""
    data: List[T] = Field(description="List of items")
    total: int = Field(description="Total number of items")
    limit: int = Field(description="Number of items per page")
    offset: int = Field(description="Number of items skipped")


class RequestConfig(BaseModel):
    """Configuration for individual API requests."""
    headers: Optional[Dict[str, str]] = Field(None, description="Additional request headers")
    timeout: Optional[float] = Field(None, gt=0, description="Request timeout in seconds")
    retries: Optional[int] = Field(None, ge=0, le=10, description="Number of retry attempts")


class APIResponse(BaseModel, Generic[T]):
    """Generic API response wrapper."""
    success: bool = Field(description="Whether request was successful")
    data: Optional[T] = Field(None, description="Response data")
    message: Optional[str] = Field(None, description="Response message")
    timestamp: datetime = Field(description="Response timestamp")


class APIError(Exception):
    """Base exception for API errors."""
    
    def __init__(
        self,
        message: str,
        status_code: Optional[int] = None,
        error_code: Optional[str] = None,
        field: Optional[str] = None,
        details: Optional[Dict[str, Any]] = None
    ):
        super().__init__(message)
        self.message = message
        self.status_code = status_code
        self.error_code = error_code
        self.field = field
        self.details = details or {}

    def __str__(self) -> str:
        parts = [self.message]
        if self.status_code:
            parts.append(f"(HTTP {self.status_code})")
        if self.error_code:
            parts.append(f"[{self.error_code}]")
        if self.field:
            parts.append(f"Field: {self.field}")
        return " ".join(parts)


class AuthenticationError(APIError):
    """Authentication related errors."""
    pass


class ValidationError(APIError):
    """Validation errors."""
    pass


class NotFoundError(APIError):
    """Resource not found errors."""
    pass


class RateLimitError(APIError):
    """Rate limit exceeded errors."""
    pass


class NetworkError(APIError):
    """Network related errors."""
    pass


class HealthStatus(BaseModel):
    """API health status response."""
    status: str = Field(description="Overall health status")
    version: str = Field(description="API version")
    uptime: int = Field(description="Uptime in seconds")
    timestamp: datetime = Field(description="Current server timestamp")
    components: Optional[Dict[str, Any]] = Field(None, description="Component health details")


class VersionInfo(BaseModel):
    """API version information."""
    version: str = Field(description="API version")
    build_date: str = Field(description="Build date")
    git_commit: str = Field(description="Git commit hash")
    api_version: str = Field(description="API version")


class TimePeriod(str):
    """Valid time periods for analytics."""
    HOUR_24 = "24h"
    DAYS_7 = "7d"
    DAYS_30 = "30d"
    DAYS_90 = "90d"
    YEAR_1 = "1y"
    ALL = "all"


class ExportFormat(str):
    """Valid export formats."""
    CSV = "csv"
    JSON = "json"
    XLSX = "xlsx"


class SortOrder(str):
    """Sort order options."""
    ASC = "asc"
    DESC = "desc"