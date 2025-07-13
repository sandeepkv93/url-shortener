"""Analytics related type definitions."""

from typing import Dict, List, Optional, Any
from datetime import datetime
from pydantic import BaseModel, Field

from .urls import ShortURL, URLPerformance


class Click(BaseModel):
    """Individual click record."""
    id: int = Field(description="Click ID")
    short_url_id: int = Field(description="Short URL ID")
    ip_address: str = Field(description="Visitor IP address")
    user_agent: str = Field(description="User agent string")
    referer: Optional[str] = Field(None, description="Referrer URL")
    country: Optional[str] = Field(None, description="Country code")
    city: Optional[str] = Field(None, description="City name")
    device: Optional[str] = Field(None, description="Device type")
    browser: Optional[str] = Field(None, description="Browser name")
    os: Optional[str] = Field(None, description="Operating system")
    clicked_at: datetime = Field(description="Click timestamp")
    
    # Additional metadata
    is_unique: bool = Field(default=True, description="Whether this is a unique visitor")
    session_id: Optional[str] = Field(None, description="Session identifier")


class CountryStat(BaseModel):
    """Country-based statistics."""
    country: str = Field(description="Country name")
    country_code: str = Field(description="ISO country code")
    count: int = Field(description="Number of clicks")
    percentage: float = Field(description="Percentage of total clicks")


class CityStat(BaseModel):
    """City-based statistics."""
    city: str = Field(description="City name")
    country: str = Field(description="Country name")
    count: int = Field(description="Number of clicks")
    percentage: float = Field(description="Percentage of total clicks")


class DeviceStat(BaseModel):
    """Device-based statistics."""
    device: str = Field(description="Device type")
    count: int = Field(description="Number of clicks")
    percentage: float = Field(description="Percentage of total clicks")


class BrowserStat(BaseModel):
    """Browser-based statistics."""
    browser: str = Field(description="Browser name")
    version: Optional[str] = Field(None, description="Browser version")
    count: int = Field(description="Number of clicks")
    percentage: float = Field(description="Percentage of total clicks")


class OStat(BaseModel):
    """Operating system statistics."""
    os: str = Field(description="Operating system")
    version: Optional[str] = Field(None, description="OS version")
    count: int = Field(description="Number of clicks")
    percentage: float = Field(description="Percentage of total clicks")


class ReferrerStat(BaseModel):
    """Referrer-based statistics."""
    referrer: str = Field(description="Referrer domain")
    count: int = Field(description="Number of clicks")
    percentage: float = Field(description="Percentage of total clicks")
    is_social: bool = Field(default=False, description="Whether referrer is social media")
    is_search: bool = Field(default=False, description="Whether referrer is search engine")


class TimelineData(BaseModel):
    """Timeline data point."""
    date: datetime = Field(description="Date/time")
    clicks: int = Field(description="Number of clicks")
    unique_clicks: int = Field(description="Number of unique clicks")


class URLAnalytics(BaseModel):
    """Comprehensive URL analytics."""
    url: ShortURL = Field(description="Short URL details")
    
    # Basic metrics
    total_clicks: int = Field(description="Total click count")
    unique_clicks: int = Field(description="Unique visitor count")
    
    # Time-based metrics
    clicks_by_date: List[TimelineData] = Field(description="Daily click timeline")
    clicks_by_hour: List[Dict[str, Any]] = Field(description="Hourly distribution")
    
    # Geographic metrics
    top_countries: List[CountryStat] = Field(description="Top countries")
    top_cities: List[CityStat] = Field(description="Top cities")
    
    # Device metrics
    top_devices: List[DeviceStat] = Field(description="Top devices")
    top_browsers: List[BrowserStat] = Field(description="Top browsers")
    top_os: List[OStat] = Field(description="Top operating systems")
    
    # Traffic sources
    top_referrers: List[ReferrerStat] = Field(description="Top referrers")
    direct_clicks: int = Field(description="Direct traffic count")
    
    # Recent activity
    recent_clicks: List[Click] = Field(description="Recent click records")
    
    # Performance metrics
    bounce_rate: float = Field(description="Bounce rate percentage")
    avg_session_duration: float = Field(description="Average session duration in seconds")


class DashboardStats(BaseModel):
    """Dashboard overview statistics."""
    total_urls: int = Field(description="Total number of URLs")
    total_clicks: int = Field(description="Total number of clicks")
    clicks_today: int = Field(description="Clicks today")
    clicks_this_week: int = Field(description="Clicks this week")
    clicks_this_month: int = Field(description="Clicks this month")
    
    # Performance metrics
    top_urls: List[URLPerformance] = Field(description="Top performing URLs")
    recent_urls: List[ShortURL] = Field(description="Recently created URLs")
    
    # Growth metrics
    click_growth_rate: float = Field(description="Click growth rate percentage")
    url_creation_rate: float = Field(description="URL creation rate")
    
    # Geographic distribution
    top_countries: List[CountryStat] = Field(description="Top countries")
    
    # Time-based metrics
    clicks_timeline: List[TimelineData] = Field(description="Click timeline")


class GeoStats(BaseModel):
    """Geographic statistics."""
    countries: List[CountryStat] = Field(description="Country statistics")
    cities: List[CityStat] = Field(description="City statistics")
    total_clicks: int = Field(description="Total clicks")
    unique_countries: int = Field(description="Number of unique countries")
    unique_cities: int = Field(description="Number of unique cities")


class DeviceStats(BaseModel):
    """Device and browser statistics."""
    devices: List[DeviceStat] = Field(description="Device statistics")
    browsers: List[BrowserStat] = Field(description="Browser statistics")
    operating_systems: List[OStat] = Field(description="OS statistics")
    total_clicks: int = Field(description="Total clicks")


class ReferrerStats(BaseModel):
    """Referrer statistics."""
    referrers: List[ReferrerStat] = Field(description="Referrer statistics")
    domains: List[ReferrerStat] = Field(description="Domain statistics")
    social_media_clicks: int = Field(description="Social media clicks")
    search_engine_clicks: int = Field(description="Search engine clicks")
    direct_clicks: int = Field(description="Direct clicks")
    total_clicks: int = Field(description="Total clicks")


class AnalyticsFilter(BaseModel):
    """Analytics filtering options."""
    start_date: Optional[datetime] = Field(None, description="Start date filter")
    end_date: Optional[datetime] = Field(None, description="End date filter")
    country: Optional[str] = Field(None, description="Country filter")
    device: Optional[str] = Field(None, description="Device filter")
    browser: Optional[str] = Field(None, description="Browser filter")
    referrer: Optional[str] = Field(None, description="Referrer filter")


class AnalyticsExportRequest(BaseModel):
    """Analytics export request."""
    url_ids: Optional[List[int]] = Field(None, description="Specific URL IDs")
    format: str = Field(default="csv", description="Export format")
    period: Optional[str] = Field(None, description="Time period")
    start_date: Optional[datetime] = Field(None, description="Start date")
    end_date: Optional[datetime] = Field(None, description="End date")
    include_details: bool = Field(default=False, description="Include detailed data")
    filters: Optional[AnalyticsFilter] = Field(None, description="Additional filters")


class RealTimeStats(BaseModel):
    """Real-time analytics statistics."""
    active_visitors: int = Field(description="Current active visitors")
    clicks_last_minute: int = Field(description="Clicks in last minute")
    clicks_last_hour: int = Field(description="Clicks in last hour")
    top_active_urls: List[URLPerformance] = Field(description="Currently active URLs")
    recent_clicks: List[Click] = Field(description="Most recent clicks")
    clicks_per_minute: List[Dict[str, Any]] = Field(description="Clicks per minute timeline")


class ComparisonStats(BaseModel):
    """Comparison between two time periods."""
    current_period: URLAnalytics = Field(description="Current period analytics")
    previous_period: URLAnalytics = Field(description="Previous period analytics")
    growth_metrics: Dict[str, float] = Field(description="Growth percentages")
    
    # Specific growth metrics
    clicks_growth: float = Field(description="Clicks growth percentage")
    unique_clicks_growth: float = Field(description="Unique clicks growth percentage")
    countries_growth: float = Field(description="Countries growth percentage")
    devices_growth: float = Field(description="Devices growth percentage")


class ClickHeatmap(BaseModel):
    """Click heatmap data."""
    hourly: List[Dict[str, Any]] = Field(description="Hourly distribution")
    daily: List[Dict[str, Any]] = Field(description="Daily distribution")
    weekly: List[Dict[str, Any]] = Field(description="Weekly distribution")
    peak_hours: List[int] = Field(description="Peak hour indicators")
    peak_days: List[str] = Field(description="Peak day indicators")


class AnalyticsInsights(BaseModel):
    """Analytics insights and recommendations."""
    insights: List[Dict[str, Any]] = Field(description="Generated insights")
    recommendations: List[str] = Field(description="Improvement recommendations")
    performance_score: float = Field(description="Overall performance score")
    trend_analysis: Dict[str, Any] = Field(description="Trend analysis data")