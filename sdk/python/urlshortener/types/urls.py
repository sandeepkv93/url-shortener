"""URL management related type definitions."""

from typing import Optional, List
from datetime import datetime
from pydantic import BaseModel, Field, HttpUrl, validator


class ShortURL(BaseModel):
    """Short URL model."""
    id: int = Field(description="URL ID")
    short_code: str = Field(description="Short URL code")
    original_url: HttpUrl = Field(description="Original long URL")
    title: Optional[str] = Field(None, description="URL title")
    description: Optional[str] = Field(None, description="URL description")
    custom_alias: Optional[str] = Field(None, description="Custom alias")
    password: Optional[str] = Field(None, description="Password protection")
    expires_at: Optional[datetime] = Field(None, description="Expiration date")
    is_active: bool = Field(default=True, description="Whether URL is active")
    click_count: int = Field(default=0, description="Total click count")
    user_id: int = Field(description="Owner user ID")
    created_at: datetime = Field(description="Creation timestamp")
    updated_at: datetime = Field(description="Last update timestamp")
    
    # Computed fields
    is_expired: bool = Field(default=False, description="Whether URL has expired")
    short_url: Optional[str] = Field(None, description="Full short URL")


class CreateURLRequest(BaseModel):
    """Request to create a new short URL."""
    original_url: HttpUrl = Field(description="Original long URL")
    title: Optional[str] = Field(None, max_length=255, description="URL title")
    description: Optional[str] = Field(None, max_length=1000, description="URL description")
    custom_alias: Optional[str] = Field(None, max_length=50, description="Custom alias")
    password: Optional[str] = Field(None, max_length=100, description="Password protection")
    expires_at: Optional[datetime] = Field(None, description="Expiration date")
    tags: Optional[List[str]] = Field(None, description="URL tags")
    
    @validator('custom_alias')
    def validate_custom_alias(cls, v):
        if v is not None:
            # Only allow alphanumeric, hyphens, and underscores
            import re
            if not re.match(r'^[a-zA-Z0-9_-]+$', v):
                raise ValueError('Custom alias can only contain letters, numbers, hyphens, and underscores')
        return v


class UpdateURLRequest(BaseModel):
    """Request to update an existing short URL."""
    title: Optional[str] = Field(None, max_length=255, description="URL title")
    description: Optional[str] = Field(None, max_length=1000, description="URL description")
    password: Optional[str] = Field(None, max_length=100, description="Password protection")
    expires_at: Optional[datetime] = Field(None, description="Expiration date")
    is_active: Optional[bool] = Field(None, description="Whether URL is active")
    tags: Optional[List[str]] = Field(None, description="URL tags")


class BulkCreateURLRequest(BaseModel):
    """Request to create multiple URLs."""
    urls: List[CreateURLRequest] = Field(description="List of URLs to create")


class BulkDeleteRequest(BaseModel):
    """Request to delete multiple URLs."""
    ids: List[int] = Field(description="List of URL IDs to delete")


class URLSearchRequest(BaseModel):
    """Request to search URLs."""
    query: str = Field(min_length=1, description="Search query")
    tags: Optional[List[str]] = Field(None, description="Filter by tags")
    is_active: Optional[bool] = Field(None, description="Filter by active status")
    created_after: Optional[datetime] = Field(None, description="Filter by creation date")
    created_before: Optional[datetime] = Field(None, description="Filter by creation date")


class URLStats(BaseModel):
    """URL statistics summary."""
    total_clicks: int = Field(description="Total click count")
    unique_clicks: int = Field(description="Unique visitor count")
    clicks_today: int = Field(description="Clicks today")
    clicks_this_week: int = Field(description="Clicks this week")
    clicks_this_month: int = Field(description="Clicks this month")
    last_clicked_at: Optional[datetime] = Field(None, description="Last click timestamp")
    
    # Geographic stats
    top_countries: List[str] = Field(default_factory=list, description="Top countries")
    top_cities: List[str] = Field(default_factory=list, description="Top cities")
    
    # Device stats
    top_devices: List[str] = Field(default_factory=list, description="Top devices")
    top_browsers: List[str] = Field(default_factory=list, description="Top browsers")
    
    # Referrer stats
    top_referrers: List[str] = Field(default_factory=list, description="Top referrers")


class URLPerformance(BaseModel):
    """URL performance metrics."""
    url: ShortURL = Field(description="Short URL details")
    clicks: int = Field(description="Total clicks")
    unique_clicks: int = Field(description="Unique clicks")
    click_rate: float = Field(description="Click rate percentage")
    last_clicked_at: Optional[datetime] = Field(None, description="Last click timestamp")
    performance_score: float = Field(description="Performance score (0-100)")


class PasswordValidationRequest(BaseModel):
    """Request to validate URL password."""
    password: str = Field(description="Password to validate")


class URLRedirectResponse(BaseModel):
    """URL redirect response."""
    original_url: HttpUrl = Field(description="Original URL to redirect to")
    redirect_type: str = Field(default="302", description="HTTP redirect type")
    requires_password: bool = Field(default=False, description="Whether password is required")


class AliasAvailabilityResponse(BaseModel):
    """Custom alias availability response."""
    alias: str = Field(description="Requested alias")
    available: bool = Field(description="Whether alias is available")
    suggestions: Optional[List[str]] = Field(None, description="Alternative suggestions")