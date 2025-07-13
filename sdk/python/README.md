# URL Shortener Python SDK

A comprehensive Python SDK for the URL Shortener API that provides full access to URL shortening, analytics, QR code generation, and webhook management features.

## Features

- 🔗 **Complete URL Management** - Create, update, delete, and manage short URLs
- 📊 **Advanced Analytics** - Detailed click tracking and reporting
- 🔐 **User Authentication** - Login, registration, and profile management
- 🎯 **QR Code Generation** - Create customizable QR codes for URLs
- 🔔 **Webhook Support** - Real-time event notifications
- ⚡ **Async/Await Support** - Both synchronous and asynchronous clients
- 🛡️ **Type Safety** - Full type hints with Pydantic models
- 🔄 **Auto Token Refresh** - Automatic authentication token management
- 🚨 **Comprehensive Error Handling** - Detailed error types and messages
- 📦 **Easy Integration** - Simple setup and intuitive API

## Installation

```bash
pip install urlshortener-sdk
```

## Quick Start

### Basic Usage

```python
from urlshortener import URLShortenerClient

# Initialize client with API key
client = URLShortenerClient.with_api_key(
    base_url='https://api.urlshortener.com/v1',
    api_key='your-api-key'
)

# Create a short URL
short_url = client.urls.create('https://example.com')
print(f"Short URL: {short_url.short_url}")
print(f"Short Code: {short_url.short_code}")

# Get analytics
analytics = client.analytics.get_url_analytics(short_url.id)
print(f"Total clicks: {analytics.total_clicks}")
```

### User Authentication

```python
# Register new user
user = client.auth.register('user@example.com', 'password123')

# Login existing user
auth_response = client.auth.login('user@example.com', 'password123')
print(f"Access token: {auth_response.access_token}")

# Client automatically uses the access token for subsequent requests
short_url = client.urls.create('https://example.com', title='My Link')
```

### Async Usage

```python
import asyncio
from urlshortener import AsyncURLShortenerClient

async def main():
    async with AsyncURLShortenerClient.with_api_key(
        base_url='https://api.urlshortener.com/v1',
        api_key='your-api-key'
    ) as client:
        # Create short URL
        short_url = await client.urls.create('https://example.com')
        
        # Get analytics
        analytics = await client.analytics.get_url_analytics(short_url.id)
        print(f"Total clicks: {analytics.total_clicks}")

asyncio.run(main())
```

## Core Features

### URL Management

```python
# Create a short URL with options
short_url = client.urls.create(
    original_url='https://example.com',
    title='My Website',
    description='Example website',
    custom_alias='my-site',
    password='secret123',
    expires_at='2024-12-31T23:59:59Z',
    tags=['website', 'example']
)

# List URLs with pagination
urls = client.urls.list(
    pagination={'limit': 10, 'offset': 0},
    tags=['website'],
    is_active=True
)

# Update URL
updated_url = client.urls.update(
    url_id=short_url.id,
    title='Updated Title',
    is_active=False
)

# Search URLs
results = client.urls.search(
    query='example',
    tags=['website'],
    is_active=True
)

# Bulk operations
urls_to_create = [
    {'original_url': 'https://site1.com'},
    {'original_url': 'https://site2.com'},
]
bulk_results = client.urls.bulk_create(urls_to_create)

# Delete URL
client.urls.delete(short_url.id)
```

### Analytics

```python
# Get comprehensive analytics for a URL
analytics = client.analytics.get_url_analytics(
    url_id=short_url.id,
    start_date=datetime(2024, 1, 1),
    end_date=datetime(2024, 12, 31)
)

print(f"Total clicks: {analytics.total_clicks}")
print(f"Unique clicks: {analytics.unique_clicks}")
print(f"Top countries: {analytics.top_countries}")

# Get dashboard stats
dashboard = client.analytics.get_dashboard_stats(period='30d')
print(f"Total URLs: {dashboard.total_urls}")
print(f"Total clicks: {dashboard.total_clicks}")

# Get real-time stats
realtime = client.analytics.get_real_time_stats()
print(f"Active visitors: {realtime.active_visitors}")

# Get geographic stats
geo_stats = client.analytics.get_geographic_stats(
    url_id=short_url.id,
    period='7d'
)

# Export analytics data
export_data = client.analytics.export_analytics(
    url_ids=[short_url.id],
    format='csv',
    period='30d',
    include_details=True
)
```

### QR Code Generation

```python
from urlshortener.types.qr import QRCodeOptions, QRFormat, QRErrorCorrectionLevel

# Generate basic QR code
qr_code = client.qr.generate('https://example.com')
print(f"QR code data: {qr_code.data}")

# Generate QR code with options
options = QRCodeOptions(
    size=300,
    format=QRFormat.PNG,
    level=QRErrorCorrectionLevel.HIGH,
    margin=2,
    dark_color='#000000',
    light_color='#FFFFFF'
)

qr_code = client.qr.generate('https://example.com', options=options)

# Generate QR code for existing short URL
qr_code = client.qr.generate_for_short_url(
    short_code=short_url.short_code,
    options=options
)

# Get QR code analytics
qr_analytics = client.qr.get_analytics(
    short_code=short_url.short_code,
    period='30d'
)
print(f"Total scans: {qr_analytics.total_scans}")
```

### Webhook Management

```python
from urlshortener.types.webhooks import WebhookEvent

# Create webhook
webhook = client.webhooks.create(
    name='My Webhook',
    url='https://mysite.com/webhook',
    events=[
        WebhookEvent.URL_CREATED,
        WebhookEvent.URL_CLICKED,
        WebhookEvent.ANALYTICS_THRESHOLD
    ],
    secret='webhook-secret'
)

# List webhooks
webhooks = client.webhooks.list()

# Test webhook
test_result = client.webhooks.test(webhook.id)
print(f"Test successful: {test_result.success}")

# Get webhook stats
stats = client.webhooks.get_stats(webhook.id, period='7d')
print(f"Success rate: {stats.success_rate}%")

# Get delivery history
deliveries = client.webhooks.get_deliveries(webhook.id)

# Verify webhook signature (for incoming webhooks)
is_valid = client.webhooks.verify_signature(
    payload='{"event": "url.created", "data": {...}}',
    signature='sha256=abc123...',
    secret='webhook-secret'
)
```

## Configuration Options

### Client Configuration

```python
# Development environment
client = URLShortenerClient.development(
    api_key='dev-api-key',
    debug=True
)

# Production environment
client = URLShortenerClient.production(
    base_url='https://api.urlshortener.com/v1',
    api_key='prod-api-key',
    timeout=30,
    max_retries=5
)

# Custom configuration
client = URLShortenerClient(
    base_url='https://api.urlshortener.com/v1',
    api_key='your-api-key',
    timeout=60,
    max_retries=3,
    retry_delay=2.0,
    debug=False
)
```

### Request Configuration

```python
from urlshortener.types.common import RequestConfig

# Custom request configuration
config = RequestConfig(
    headers={'X-Custom-Header': 'value'},
    timeout=120,
    retries=5
)

# Use with any API call
short_url = client.urls.create(
    'https://example.com',
    config=config
)
```

### Authentication

```python
# API key authentication
client = URLShortenerClient.with_api_key(
    base_url='https://api.urlshortener.com/v1',
    api_key='your-api-key'
)

# Access token authentication
client = URLShortenerClient.with_access_token(
    base_url='https://api.urlshortener.com/v1',
    access_token='your-access-token'
)

# Set up automatic token refresh
def refresh_token_callback(refresh_token):
    # Your token refresh logic
    response = requests.post('/auth/refresh', {
        'refresh_token': refresh_token
    })
    return AuthTokens(**response.json())

client.set_token_refresh_callback(refresh_token_callback, 'refresh-token')
```

## Error Handling

```python
from urlshortener.exceptions import (
    URLShortenerError,
    AuthenticationError,
    ValidationError,
    NotFoundError,
    RateLimitError
)

try:
    short_url = client.urls.create('invalid-url')
except ValidationError as e:
    print(f"Validation error: {e.message}")
    print(f"Errors: {e.errors}")
except AuthenticationError as e:
    print(f"Authentication failed: {e.message}")
except RateLimitError as e:
    print(f"Rate limited. Retry after: {e.retry_after} seconds")
except URLShortenerError as e:
    print(f"API error: {e.message}")
```

## Advanced Features

### Pagination

```python
from urlshortener.types.common import PaginationParams

# Paginate through URLs
pagination = PaginationParams(limit=20, offset=0)
while True:
    response = client.urls.list(pagination=pagination)
    
    for url in response.data:
        print(f"URL: {url.short_url}")
    
    if len(response.data) < pagination.limit:
        break
        
    pagination.offset += pagination.limit
```

### Filtering and Searching

```python
from datetime import datetime, timedelta

# Filter URLs by date range
end_date = datetime.now()
start_date = end_date - timedelta(days=30)

urls = client.urls.list(
    created_after=start_date.isoformat(),
    created_before=end_date.isoformat(),
    is_active=True,
    tags=['important']
)

# Search with filters
results = client.urls.search(
    query='example',
    tags=['website'],
    is_active=True,
    created_after=start_date.isoformat()
)
```

### Batch Operations

```python
# Create multiple URLs
urls_data = [
    {'original_url': 'https://site1.com', 'title': 'Site 1'},
    {'original_url': 'https://site2.com', 'title': 'Site 2'},
    {'original_url': 'https://site3.com', 'title': 'Site 3'},
]

created_urls = client.urls.bulk_create(urls_data)

# Delete multiple URLs
url_ids = [url.id for url in created_urls]
result = client.urls.bulk_delete(url_ids)
print(f"Deleted {result['deleted']} URLs")
```

## Type Safety

The SDK provides full type safety with Pydantic models:

```python
from urlshortener.types.urls import ShortURL, CreateURLRequest
from urlshortener.types.analytics import URLAnalytics
from urlshortener.types.qr import QRCodeResponse

# All responses are properly typed
short_url: ShortURL = client.urls.create('https://example.com')
analytics: URLAnalytics = client.analytics.get_url_analytics(short_url.id)
qr_code: QRCodeResponse = client.qr.generate(short_url.short_url)

# IDE will provide proper autocompletion and type checking
print(short_url.short_code)  # ✅ Type-safe
print(analytics.total_clicks)  # ✅ Type-safe
```

## Testing

```python
import pytest
from urlshortener import URLShortenerClient
from urlshortener.exceptions import ValidationError

def test_create_url():
    client = URLShortenerClient.development(api_key='test-key')
    
    short_url = client.urls.create('https://example.com')
    assert short_url.original_url == 'https://example.com'
    assert len(short_url.short_code) > 0

def test_invalid_url():
    client = URLShortenerClient.development(api_key='test-key')
    
    with pytest.raises(ValidationError):
        client.urls.create('not-a-url')
```

## Environment Variables

You can configure the SDK using environment variables:

```bash
export URLSHORTENER_BASE_URL=https://api.urlshortener.com/v1
export URLSHORTENER_API_KEY=your-api-key
export URLSHORTENER_TIMEOUT=30
export URLSHORTENER_MAX_RETRIES=3
export URLSHORTENER_DEBUG=false
```

```python
import os
from urlshortener import URLShortenerClient

# Automatically uses environment variables
client = URLShortenerClient(
    base_url=os.getenv('URLSHORTENER_BASE_URL'),
    api_key=os.getenv('URLSHORTENER_API_KEY'),
    timeout=int(os.getenv('URLSHORTENER_TIMEOUT', 30)),
    max_retries=int(os.getenv('URLSHORTENER_MAX_RETRIES', 3)),
    debug=os.getenv('URLSHORTENER_DEBUG', 'false').lower() == 'true'
)
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run tests: `pytest`
5. Submit a pull request

## License

This SDK is licensed under the MIT License. See the LICENSE file for details.

## Support

- Documentation: [https://docs.urlshortener.com](https://docs.urlshortener.com)
- API Reference: [https://api.urlshortener.com/docs](https://api.urlshortener.com/docs)
- Issues: [GitHub Issues](https://github.com/urlshortener/python-sdk/issues)
- Email: support@urlshortener.com