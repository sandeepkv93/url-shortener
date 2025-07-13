"""QR code related type definitions."""

from typing import List, Optional, Dict, Any
from enum import Enum
from pydantic import BaseModel, Field, HttpUrl, validator


class QRFormat(str, Enum):
    """Available QR code formats."""
    PNG = "png"
    SVG = "svg"
    JPEG = "jpeg"
    PDF = "pdf"


class QRErrorCorrectionLevel(str, Enum):
    """QR code error correction levels."""
    LOW = "L"      # ~7% correction
    MEDIUM = "M"   # ~15% correction
    QUARTILE = "Q" # ~25% correction
    HIGH = "H"     # ~30% correction


class QRCodeOptions(BaseModel):
    """QR code generation options."""
    size: Optional[int] = Field(None, ge=50, le=2000, description="QR code size in pixels")
    format: Optional[QRFormat] = Field(QRFormat.PNG, description="Output format")
    level: Optional[QRErrorCorrectionLevel] = Field(QRErrorCorrectionLevel.MEDIUM, description="Error correction level")
    margin: Optional[int] = Field(None, ge=0, le=50, description="Margin size")
    dark_color: Optional[str] = Field(None, description="Dark color (hex code)")
    light_color: Optional[str] = Field(None, description="Light color (hex code)")
    
    @validator('dark_color', 'light_color')
    def validate_color(cls, v):
        """Validate color hex codes."""
        if v is not None:
            import re
            if not re.match(r'^#[0-9A-Fa-f]{6}$', v):
                raise ValueError('Color must be a valid hex code (e.g., #000000)')
        return v


class QRCodeResponse(BaseModel):
    """QR code generation response."""
    data: str = Field(description="Base64 encoded image or SVG string")
    format: QRFormat = Field(description="Output format")
    size: int = Field(description="Generated size in pixels")
    url: str = Field(description="Original URL")
    expires_at: Optional[str] = Field(None, description="Cache expiration")


class QRCodeFormats(BaseModel):
    """Available QR code formats."""
    formats: List[QRFormat] = Field(description="Supported formats")
    default_format: QRFormat = Field(description="Default format")


class QRCodeSizes(BaseModel):
    """Available QR code sizes."""
    sizes: List[int] = Field(description="Supported sizes")
    default_size: int = Field(description="Default size")
    min_size: int = Field(description="Minimum size")
    max_size: int = Field(description="Maximum size")


class QRCodeGenerationRequest(BaseModel):
    """QR code generation request."""
    url: HttpUrl = Field(description="URL to encode")
    options: Optional[QRCodeOptions] = Field(None, description="Generation options")


class QRCodeValidationResult(BaseModel):
    """QR code options validation result."""
    valid: bool = Field(description="Whether options are valid")
    errors: Optional[List[str]] = Field(None, description="Validation errors")
    warnings: Optional[List[str]] = Field(None, description="Validation warnings")


class QRCodeBrandingOptions(BaseModel):
    """QR code branding options."""
    logo: Optional[str] = Field(None, description="Base64 encoded logo")
    logo_size: Optional[int] = Field(None, ge=5, le=30, description="Logo size as percentage")
    background_color: Optional[str] = Field(None, description="Background color")
    foreground_color: Optional[str] = Field(None, description="Foreground color")
    border_radius: Optional[int] = Field(None, ge=0, le=50, description="Border radius")
    gradient: Optional[Dict[str, str]] = Field(None, description="Gradient configuration")
    
    @validator('background_color', 'foreground_color')
    def validate_colors(cls, v):
        """Validate color hex codes."""
        if v is not None:
            import re
            if not re.match(r'^#[0-9A-Fa-f]{6}$', v):
                raise ValueError('Color must be a valid hex code')
        return v


class QRCodeBatchRequest(BaseModel):
    """Batch QR code generation request."""
    requests: List[QRCodeGenerationRequest] = Field(description="QR code requests")
    common_options: Optional[QRCodeOptions] = Field(None, description="Common options for all")


class QRCodeBatchResponse(BaseModel):
    """Batch QR code generation response."""
    qr_codes: List[QRCodeResponse] = Field(description="Generated QR codes")
    total: int = Field(description="Total number generated")
    failed: int = Field(description="Number of failures")
    errors: Optional[List[str]] = Field(None, description="Error messages")


class QRCodeAnalytics(BaseModel):
    """QR code scan analytics."""
    total_scans: int = Field(description="Total number of scans")
    scans_today: int = Field(description="Scans today")
    scans_by_date: List[Dict[str, Any]] = Field(description="Scans by date")
    top_devices: List[Dict[str, Any]] = Field(description="Top devices")
    top_countries: List[Dict[str, Any]] = Field(description="Top countries")
    scan_locations: List[Dict[str, Any]] = Field(description="Scan locations")


class QRCodeMetadata(BaseModel):
    """QR code metadata."""
    short_code: str = Field(description="Short URL code")
    original_url: HttpUrl = Field(description="Original URL")
    title: Optional[str] = Field(None, description="URL title")
    created_at: str = Field(description="Creation timestamp")
    expires_at: Optional[str] = Field(None, description="Expiration timestamp")
    is_active: bool = Field(description="Whether URL is active")
    qr_code_generated: bool = Field(description="Whether QR code was generated")
    total_scans: int = Field(description="Total number of scans")


class QRCodeTrackingParams(BaseModel):
    """QR code tracking parameters."""
    utm_source: Optional[str] = Field(None, description="UTM source")
    utm_medium: Optional[str] = Field(None, description="UTM medium")
    utm_campaign: Optional[str] = Field(None, description="UTM campaign")
    utm_term: Optional[str] = Field(None, description="UTM term")
    utm_content: Optional[str] = Field(None, description="UTM content")
    custom_params: Optional[Dict[str, str]] = Field(None, description="Custom parameters")


class QRCodeTemplate(BaseModel):
    """QR code template."""
    name: str = Field(description="Template name")
    description: str = Field(description="Template description")
    options: QRCodeOptions = Field(description="Default options")
    branding: Optional[QRCodeBrandingOptions] = Field(None, description="Branding options")
    use_case: str = Field(description="Intended use case")


class QRCodeDataValidation(BaseModel):
    """QR code data validation result."""
    is_valid: bool = Field(description="Whether data is valid")
    data_type: str = Field(description="Type of data (url, text, email, etc.)")
    data: str = Field(description="Validated data")
    warnings: Optional[List[str]] = Field(None, description="Validation warnings")


class QRCodeUsageStats(BaseModel):
    """QR code usage statistics."""
    total_generated: int = Field(description="Total QR codes generated")
    generated_today: int = Field(description="Generated today")
    generated_this_week: int = Field(description="Generated this week")
    generated_this_month: int = Field(description="Generated this month")
    most_popular_format: QRFormat = Field(description="Most popular format")
    most_popular_size: int = Field(description="Most popular size")
    average_scans_per_code: float = Field(description="Average scans per QR code")