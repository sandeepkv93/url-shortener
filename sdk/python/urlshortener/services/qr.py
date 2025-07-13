"""QR code service for URL Shortener API."""

from typing import List, Optional, Dict, Any

from ..http_client import HTTPClient, AsyncHTTPClient
from ..types.qr import (
    QRCodeOptions,
    QRCodeResponse,
    QRCodeFormats,
    QRCodeSizes,
    QRCodeGenerationRequest,
    QRCodeValidationResult,
    QRCodeBrandingOptions,
    QRCodeBatchRequest,
    QRCodeBatchResponse,
    QRCodeAnalytics,
    QRCodeMetadata,
    QRCodeTrackingParams,
    QRCodeTemplate,
    QRCodeDataValidation,
    QRCodeUsageStats,
    QRFormat,
    QRErrorCorrectionLevel
)
from ..types.common import RequestConfig


class QRService:
    """Synchronous QR code service."""

    def __init__(self, http_client: HTTPClient):
        """Initialize QR service.
        
        Args:
            http_client: HTTP client instance
        """
        self.http_client = http_client

    def generate(
        self,
        url: str,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate a QR code for a URL.
        
        Args:
            url: URL to encode in QR code
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated QR code
        """
        request = QRCodeGenerationRequest(
            url=url,
            options=options
        )
        
        response = self.http_client.post(
            '/qr/generate',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return QRCodeResponse(**response['data'])

    def generate_for_short_url(
        self,
        short_code: str,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate a QR code for an existing short URL.
        
        Args:
            short_code: Short URL code
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated QR code
        """
        data = {}
        if options:
            data['options'] = options.dict(exclude_none=True)
            
        response = self.http_client.post(
            f'/qr/short/{short_code}',
            data=data,
            config=config
        )
        
        return QRCodeResponse(**response['data'])

    def generate_batch(
        self,
        requests: List[QRCodeGenerationRequest],
        common_options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeBatchResponse:
        """Generate multiple QR codes in batch.
        
        Args:
            requests: List of QR code generation requests
            common_options: Common options for all QR codes
            config: Request configuration
            
        Returns:
            QRCodeBatchResponse: Batch generation results
        """
        request = QRCodeBatchRequest(
            requests=requests,
            common_options=common_options
        )
        
        response = self.http_client.post(
            '/qr/batch',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return QRCodeBatchResponse(**response['data'])

    def get_formats(self, config: Optional[RequestConfig] = None) -> QRCodeFormats:
        """Get available QR code formats.
        
        Args:
            config: Request configuration
            
        Returns:
            QRCodeFormats: Available formats
        """
        response = self.http_client.get('/qr/formats', config=config)
        return QRCodeFormats(**response['data'])

    def get_sizes(self, config: Optional[RequestConfig] = None) -> QRCodeSizes:
        """Get available QR code sizes.
        
        Args:
            config: Request configuration
            
        Returns:
            QRCodeSizes: Available sizes
        """
        response = self.http_client.get('/qr/sizes', config=config)
        return QRCodeSizes(**response['data'])

    def validate_options(
        self,
        options: QRCodeOptions,
        config: Optional[RequestConfig] = None
    ) -> QRCodeValidationResult:
        """Validate QR code options.
        
        Args:
            options: QR code options to validate
            config: Request configuration
            
        Returns:
            QRCodeValidationResult: Validation result
        """
        response = self.http_client.post(
            '/qr/validate',
            data=options.dict(),
            config=config
        )
        
        return QRCodeValidationResult(**response['data'])

    def validate_data(
        self,
        data: str,
        config: Optional[RequestConfig] = None
    ) -> QRCodeDataValidation:
        """Validate QR code data.
        
        Args:
            data: Data to validate
            config: Request configuration
            
        Returns:
            QRCodeDataValidation: Data validation result
        """
        request = {'data': data}
        
        response = self.http_client.post(
            '/qr/validate-data',
            data=request,
            config=config
        )
        
        return QRCodeDataValidation(**response['data'])

    def get_analytics(
        self,
        short_code: str,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeAnalytics:
        """Get QR code scan analytics.
        
        Args:
            short_code: Short URL code
            period: Time period for analytics
            config: Request configuration
            
        Returns:
            QRCodeAnalytics: QR code analytics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            f'/qr/analytics/{short_code}',
            params=params,
            config=config
        )
        
        return QRCodeAnalytics(**response['data'])

    def get_metadata(
        self,
        short_code: str,
        config: Optional[RequestConfig] = None
    ) -> QRCodeMetadata:
        """Get QR code metadata.
        
        Args:
            short_code: Short URL code
            config: Request configuration
            
        Returns:
            QRCodeMetadata: QR code metadata
        """
        response = self.http_client.get(f'/qr/metadata/{short_code}', config=config)
        return QRCodeMetadata(**response['data'])

    def get_usage_stats(
        self,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeUsageStats:
        """Get QR code usage statistics.
        
        Args:
            period: Time period for statistics
            config: Request configuration
            
        Returns:
            QRCodeUsageStats: Usage statistics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = self.http_client.get(
            '/qr/usage-stats',
            params=params,
            config=config
        )
        
        return QRCodeUsageStats(**response['data'])

    def get_templates(self, config: Optional[RequestConfig] = None) -> List[QRCodeTemplate]:
        """Get available QR code templates.
        
        Args:
            config: Request configuration
            
        Returns:
            List[QRCodeTemplate]: Available templates
        """
        response = self.http_client.get('/qr/templates', config=config)
        return [QRCodeTemplate(**item) for item in response['data']]

    def create_template(
        self,
        name: str,
        description: str,
        options: QRCodeOptions,
        branding: Optional[QRCodeBrandingOptions] = None,
        use_case: str = 'general',
        config: Optional[RequestConfig] = None
    ) -> QRCodeTemplate:
        """Create a custom QR code template.
        
        Args:
            name: Template name
            description: Template description
            options: Default QR code options
            branding: Branding options
            use_case: Intended use case
            config: Request configuration
            
        Returns:
            QRCodeTemplate: Created template
        """
        template = QRCodeTemplate(
            name=name,
            description=description,
            options=options,
            branding=branding,
            use_case=use_case
        )
        
        response = self.http_client.post(
            '/qr/templates',
            data=template.dict(exclude_none=True),
            config=config
        )
        
        return QRCodeTemplate(**response['data'])

    def generate_with_tracking(
        self,
        url: str,
        tracking_params: QRCodeTrackingParams,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate QR code with tracking parameters.
        
        Args:
            url: URL to encode
            tracking_params: UTM and custom tracking parameters
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated QR code with tracking
        """
        data = {
            'url': url,
            'tracking_params': tracking_params.dict(exclude_none=True)
        }
        
        if options:
            data['options'] = options.dict(exclude_none=True)
            
        response = self.http_client.post(
            '/qr/generate-with-tracking',
            data=data,
            config=config
        )
        
        return QRCodeResponse(**response['data'])

    def generate_branded(
        self,
        url: str,
        branding: QRCodeBrandingOptions,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate branded QR code.
        
        Args:
            url: URL to encode
            branding: Branding options
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated branded QR code
        """
        data = {
            'url': url,
            'branding': branding.dict(exclude_none=True)
        }
        
        if options:
            data['options'] = options.dict(exclude_none=True)
            
        response = self.http_client.post(
            '/qr/generate-branded',
            data=data,
            config=config
        )
        
        return QRCodeResponse(**response['data'])


class AsyncQRService:
    """Asynchronous QR code service."""

    def __init__(self, http_client: AsyncHTTPClient):
        """Initialize async QR service.
        
        Args:
            http_client: Async HTTP client instance
        """
        self.http_client = http_client

    async def generate(
        self,
        url: str,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate a QR code for a URL.
        
        Args:
            url: URL to encode in QR code
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated QR code
        """
        request = QRCodeGenerationRequest(
            url=url,
            options=options
        )
        
        response = await self.http_client.post(
            '/qr/generate',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return QRCodeResponse(**response['data'])

    async def generate_for_short_url(
        self,
        short_code: str,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate a QR code for an existing short URL.
        
        Args:
            short_code: Short URL code
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated QR code
        """
        data = {}
        if options:
            data['options'] = options.dict(exclude_none=True)
            
        response = await self.http_client.post(
            f'/qr/short/{short_code}',
            data=data,
            config=config
        )
        
        return QRCodeResponse(**response['data'])

    async def generate_batch(
        self,
        requests: List[QRCodeGenerationRequest],
        common_options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeBatchResponse:
        """Generate multiple QR codes in batch.
        
        Args:
            requests: List of QR code generation requests
            common_options: Common options for all QR codes
            config: Request configuration
            
        Returns:
            QRCodeBatchResponse: Batch generation results
        """
        request = QRCodeBatchRequest(
            requests=requests,
            common_options=common_options
        )
        
        response = await self.http_client.post(
            '/qr/batch',
            data=request.dict(exclude_none=True),
            config=config
        )
        
        return QRCodeBatchResponse(**response['data'])

    async def get_formats(self, config: Optional[RequestConfig] = None) -> QRCodeFormats:
        """Get available QR code formats.
        
        Args:
            config: Request configuration
            
        Returns:
            QRCodeFormats: Available formats
        """
        response = await self.http_client.get('/qr/formats', config=config)
        return QRCodeFormats(**response['data'])

    async def get_analytics(
        self,
        short_code: str,
        period: Optional[str] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeAnalytics:
        """Get QR code scan analytics.
        
        Args:
            short_code: Short URL code
            period: Time period for analytics
            config: Request configuration
            
        Returns:
            QRCodeAnalytics: QR code analytics
        """
        params = {}
        if period:
            params['period'] = period
            
        response = await self.http_client.get(
            f'/qr/analytics/{short_code}',
            params=params,
            config=config
        )
        
        return QRCodeAnalytics(**response['data'])

    async def generate_with_tracking(
        self,
        url: str,
        tracking_params: QRCodeTrackingParams,
        options: Optional[QRCodeOptions] = None,
        config: Optional[RequestConfig] = None
    ) -> QRCodeResponse:
        """Generate QR code with tracking parameters.
        
        Args:
            url: URL to encode
            tracking_params: UTM and custom tracking parameters
            options: QR code generation options
            config: Request configuration
            
        Returns:
            QRCodeResponse: Generated QR code with tracking
        """
        data = {
            'url': url,
            'tracking_params': tracking_params.dict(exclude_none=True)
        }
        
        if options:
            data['options'] = options.dict(exclude_none=True)
            
        response = await self.http_client.post(
            '/qr/generate-with-tracking',
            data=data,
            config=config
        )
        
        return QRCodeResponse(**response['data'])