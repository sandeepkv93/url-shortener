"""Analytics service for URL Shortener API."""

from typing import List, Optional, Dict, Any
from datetime import datetime

from ..http_client import HTTPClient, AsyncHTTPClient
from ..types.analytics import (
    Click,
    URLAnalytics,
    DashboardStats,
    GeoStats,
    DeviceStats,
    ReferrerStats,
    AnalyticsFilter,
    AnalyticsExportRequest,
    RealTimeStats,
    ComparisonStats,
    ClickHeatmap,
    AnalyticsInsights,
    CountryStat,
    CityStat,
    DeviceStat,
    BrowserStat,
    OStat,
    ReferrerStat,
    TimelineData
)
from ..types.common import RequestConfig


class AnalyticsService:
    """Synchronous analytics service."""

    def __init__(self, http_client: HTTPClient):
        """Initialize analytics service.
        
        Args:
            http_client: HTTP client instance
        """
        self.http_client = http_client

    def get_dashboard_stats(
        self,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> DashboardStats:
        """Get dashboard overview statistics.
        
        Args:
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            DashboardStats: Dashboard statistics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/dashboard',
            params=params,
            config=config
        )
        
        return DashboardStats(**response['data'])

    def get_url_analytics(
        self,
        url_id: int,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
        config: Optional[RequestConfig] = None
    ) -> URLAnalytics:
        """Get comprehensive analytics for a specific URL.
        
        Args:
            url_id: URL ID
            start_date: Start date for analytics
            end_date: End date for analytics
            config: Request configuration
            
        Returns:
            URLAnalytics: Comprehensive URL analytics
        """
        params = {}
        if start_date:
            params['start_date'] = start_date.isoformat()
        if end_date:
            params['end_date'] = end_date.isoformat()
            
        response = self.http_client.get(
            f'/analytics/urls/{url_id}',
            params=params,
            config=config
        )
        
        return URLAnalytics(**response['data'])

    def get_clicks(
        self,
        url_id: Optional[int] = None,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
        limit: int = 100,
        offset: int = 0,
        config: Optional[RequestConfig] = None
    ) -> List[Click]:
        """Get individual click records.
        
        Args:
            url_id: Filter by URL ID
            start_date: Start date filter
            end_date: End date filter
            limit: Maximum number of records
            offset: Number of records to skip
            config: Request configuration
            
        Returns:
            List[Click]: List of click records
        """
        params = {
            'limit': limit,
            'offset': offset
        }
        
        if url_id:
            params['url_id'] = url_id
        if start_date:
            params['start_date'] = start_date.isoformat()
        if end_date:
            params['end_date'] = end_date.isoformat()
            
        response = self.http_client.get(
            '/analytics/clicks',
            params=params,
            config=config
        )
        
        return [Click(**item) for item in response['data']]

    def get_geographic_stats(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> GeoStats:
        """Get geographic statistics.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            GeoStats: Geographic statistics
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/geo',
            params=params,
            config=config
        )
        
        return GeoStats(**response['data'])

    def get_device_stats(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> DeviceStats:
        """Get device and browser statistics.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            DeviceStats: Device statistics
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/devices',
            params=params,
            config=config
        )
        
        return DeviceStats(**response['data'])

    def get_referrer_stats(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> ReferrerStats:
        """Get referrer statistics.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            ReferrerStats: Referrer statistics
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/referrers',
            params=params,
            config=config
        )
        
        return ReferrerStats(**response['data'])

    def get_timeline_data(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        granularity: str = 'day',
        config: Optional[RequestConfig] = None
    ) -> List[TimelineData]:
        """Get timeline data for clicks.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            granularity: Data granularity (hour, day, week, month)
            config: Request configuration
            
        Returns:
            List[TimelineData]: Timeline data points
        """
        params = {
            'granularity': granularity
        }
        
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/timeline',
            params=params,
            config=config
        )
        
        return [TimelineData(**item) for item in response['data']]

    def get_real_time_stats(
        self,
        config: Optional[RequestConfig] = None
    ) -> RealTimeStats:
        """Get real-time analytics statistics.
        
        Args:
            config: Request configuration
            
        Returns:
            RealTimeStats: Real-time statistics
        """
        response = self.http_client.get('/analytics/realtime', config=config)
        return RealTimeStats(**response['data'])

    def get_comparison_stats(
        self,
        url_id: int,
        current_start: datetime,
        current_end: datetime,
        previous_start: datetime,
        previous_end: datetime,
        config: Optional[RequestConfig] = None
    ) -> ComparisonStats:
        """Compare analytics between two time periods.
        
        Args:
            url_id: URL ID
            current_start: Current period start date
            current_end: Current period end date
            previous_start: Previous period start date
            previous_end: Previous period end date
            config: Request configuration
            
        Returns:
            ComparisonStats: Comparison statistics
        """
        params = {
            'current_start': current_start.isoformat(),
            'current_end': current_end.isoformat(),
            'previous_start': previous_start.isoformat(),
            'previous_end': previous_end.isoformat()
        }
        
        response = self.http_client.get(
            f'/analytics/urls/{url_id}/compare',
            params=params,
            config=config
        )
        
        return ComparisonStats(**response['data'])

    def get_click_heatmap(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> ClickHeatmap:
        """Get click heatmap data.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            ClickHeatmap: Heatmap data
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/heatmap',
            params=params,
            config=config
        )
        
        return ClickHeatmap(**response['data'])

    def get_insights(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> AnalyticsInsights:
        """Get analytics insights and recommendations.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            AnalyticsInsights: Analytics insights
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/insights',
            params=params,
            config=config
        )
        
        return AnalyticsInsights(**response['data'])

    def export_analytics(
        self,
        url_ids: Optional[List[int]] = None,
        format: str = 'csv',
        period: Optional[str] = None,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
        include_details: bool = False,
        filters: Optional[AnalyticsFilter] = None,
        config: Optional[RequestConfig] = None
    ) -> Dict[str, Any]:
        """Export analytics data.
        
        Args:
            url_ids: Specific URL IDs to export
            format: Export format (csv, json, xlsx)
            period: Time period (24h, 7d, 30d, etc.)
            start_date: Start date for export
            end_date: End date for export
            include_details: Include detailed click data
            filters: Additional filters
            config: Request configuration
            
        Returns:
            Dict[str, Any]: Export information
        """
        request = AnalyticsExportRequest(
            url_ids=url_ids,
            format=format,
            period=period,
            start_date=start_date,
            end_date=end_date,
            include_details=include_details,
            filters=filters
        )
        
        response = self.http_client.post(
            '/analytics/export',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return response['data']

    def get_top_countries(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        limit: int = 10,
        config: Optional[RequestConfig] = None
    ) -> List[CountryStat]:
        """Get top countries by clicks.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            limit: Number of results
            config: Request configuration
            
        Returns:
            List[CountryStat]: Top countries
        """
        params = {
            'limit': limit
        }
        
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/top/countries',
            params=params,
            config=config
        )
        
        return [CountryStat(**item) for item in response['data']]

    def get_top_cities(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        limit: int = 10,
        config: Optional[RequestConfig] = None
    ) -> List[CityStat]:
        """Get top cities by clicks.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            limit: Number of results
            config: Request configuration
            
        Returns:
            List[CityStat]: Top cities
        """
        params = {
            'limit': limit
        }
        
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/top/cities',
            params=params,
            config=config
        )
        
        return [CityStat(**item) for item in response['data']]

    def get_top_browsers(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        limit: int = 10,
        config: Optional[RequestConfig] = None
    ) -> List[BrowserStat]:
        """Get top browsers by clicks.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            limit: Number of results
            config: Request configuration
            
        Returns:
            List[BrowserStat]: Top browsers
        """
        params = {
            'limit': limit
        }
        
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/analytics/top/browsers',
            params=params,
            config=config
        )
        
        return [BrowserStat(**item) for item in response['data']]


class AsyncAnalyticsService:
    """Asynchronous analytics service."""

    def __init__(self, http_client: AsyncHTTPClient):
        """Initialize async analytics service.
        
        Args:
            http_client: Async HTTP client instance
        """
        self.http_client = http_client

    async def get_dashboard_stats(
        self,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> DashboardStats:
        """Get dashboard overview statistics.
        
        Args:
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            DashboardStats: Dashboard statistics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            '/analytics/dashboard',
            params=params,
            config=config
        )
        
        return DashboardStats(**response['data'])

    async def get_url_analytics(
        self,
        url_id: int,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
        config: Optional[RequestConfig] = None
    ) -> URLAnalytics:
        """Get comprehensive analytics for a specific URL.
        
        Args:
            url_id: URL ID
            start_date: Start date for analytics
            end_date: End date for analytics
            config: Request configuration
            
        Returns:
            URLAnalytics: Comprehensive URL analytics
        """
        params = {}
        if start_date:
            params['start_date'] = start_date.isoformat()
        if end_date:
            params['end_date'] = end_date.isoformat()
            
        response = await self.http_client.get(
            f'/analytics/urls/{url_id}',
            params=params,
            config=config
        )
        
        return URLAnalytics(**response['data'])

    async def get_clicks(
        self,
        url_id: Optional[int] = None,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
        limit: int = 100,
        offset: int = 0,
        config: Optional[RequestConfig] = None
    ) -> List[Click]:
        """Get individual click records.
        
        Args:
            url_id: Filter by URL ID
            start_date: Start date filter
            end_date: End date filter
            limit: Maximum number of records
            offset: Number of records to skip
            config: Request configuration
            
        Returns:
            List[Click]: List of click records
        """
        params = {
            'limit': limit,
            'offset': offset
        }
        
        if url_id:
            params['url_id'] = url_id
        if start_date:
            params['start_date'] = start_date.isoformat()
        if end_date:
            params['end_date'] = end_date.isoformat()
            
        response = await self.http_client.get(
            '/analytics/clicks',
            params=params,
            config=config
        )
        
        return [Click(**item) for item in response['data']]

    async def get_geographic_stats(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> GeoStats:
        """Get geographic statistics.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            GeoStats: Geographic statistics
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            '/analytics/geo',
            params=params,
            config=config
        )
        
        return GeoStats(**response['data'])

    async def get_device_stats(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> DeviceStats:
        """Get device and browser statistics.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            DeviceStats: Device statistics
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            '/analytics/devices',
            params=params,
            config=config
        )
        
        return DeviceStats(**response['data'])

    async def get_referrer_stats(
        self,
        url_id: Optional[int] = None,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> ReferrerStats:
        """Get referrer statistics.
        
        Args:
            url_id: Filter by URL ID
            period: Time period (24h, 7d, 30d, etc.)
            config: Request configuration
            
        Returns:
            ReferrerStats: Referrer statistics
        """
        params = {}
        if url_id:
            params['url_id'] = url_id
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            '/analytics/referrers',
            params=params,
            config=config
        )
        
        return ReferrerStats(**response['data'])

    async def get_real_time_stats(
        self,
        config: Optional[RequestConfig] = None
    ) -> RealTimeStats:
        """Get real-time analytics statistics.
        
        Args:
            config: Request configuration
            
        Returns:
            RealTimeStats: Real-time statistics
        """
        response = await self.http_client.get('/analytics/realtime', config=config)
        return RealTimeStats(**response['data'])

    async def export_analytics(
        self,
        url_ids: Optional[List[int]] = None,
        format: str = 'csv',
        period: Optional[str] = None,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
        include_details: bool = False,
        filters: Optional[AnalyticsFilter] = None,
        config: Optional[RequestConfig] = None
    ) -> Dict[str, Any]:
        """Export analytics data.
        
        Args:
            url_ids: Specific URL IDs to export
            format: Export format (csv, json, xlsx)
            period: Time period (24h, 7d, 30d, etc.)
            start_date: Start date for export
            end_date: End date for export
            include_details: Include detailed click data
            filters: Additional filters
            config: Request configuration
            
        Returns:
            Dict[str, Any]: Export information
        """
        request = AnalyticsExportRequest(
            url_ids=url_ids,
            format=format,
            period=period,
            start_date=start_date,
            end_date=end_date,
            include_details=include_details,
            filters=filters
        )
        
        response = await self.http_client.post(
            '/analytics/export',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return response['data']