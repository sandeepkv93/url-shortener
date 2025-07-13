# URL Shortener JavaScript/TypeScript SDK

Official JavaScript/TypeScript SDK for the URL Shortener API. This SDK provides a comprehensive, type-safe interface for interacting with the URL Shortener service.

## Features

- ✅ **TypeScript Support** - Full type definitions and IntelliSense support
- ✅ **Authentication** - Automatic token management and refresh
- ✅ **URL Management** - Create, update, delete, and manage short URLs
- ✅ **Analytics** - Comprehensive analytics and reporting
- ✅ **QR Codes** - Generate and customize QR codes
- ✅ **Webhooks** - Real-time event notifications
- ✅ **Error Handling** - Structured error responses
- ✅ **Retry Logic** - Automatic retry with exponential backoff
- ✅ **Modular Design** - Use only the services you need

## Installation

```bash
npm install @urlshortener/sdk
# or
yarn add @urlshortener/sdk
```

## Quick Start

```typescript
import { URLShortenerClient } from '@urlshortener/sdk';

// Initialize the client
const client = new URLShortenerClient({
  baseURL: 'https://api.urlshortener.com',
  onTokenRefresh: (tokens) => {
    // Save tokens to secure storage
    localStorage.setItem('accessToken', tokens.accessToken);
    localStorage.setItem('refreshToken', tokens.refreshToken);
  }
});

// Authenticate
const { user } = await client.auth.login({
  email: 'user@example.com',
  password: 'password'
});

// Create a short URL
const shortURL = await client.urls.create({
  originalURL: 'https://example.com/very-long-url',
  title: 'My Example URL',
  customAlias: 'my-link'
});

console.log(`Short URL: ${shortURL.shortCode}`);
```

## Client Configuration

### Basic Configuration

```typescript
import { URLShortenerClient } from '@urlshortener/sdk';

const client = new URLShortenerClient({
  baseURL: 'https://api.urlshortener.com',
  timeout: 30000,
  retries: 3,
});
```

### With Token Persistence

```typescript
import { ClientFactory } from '@urlshortener/sdk';

const client = ClientFactory.withTokenPersistence({
  baseURL: 'https://api.urlshortener.com',
  onTokenRefresh: (tokens) => {
    console.log('Tokens refreshed:', tokens);
  },
  onError: (error) => {
    console.error('API Error:', error);
  }
});
```

### Environment-Specific Clients

```typescript
import { ClientFactory } from '@urlshortener/sdk';

// Development
const devClient = ClientFactory.development();

// Production
const prodClient = ClientFactory.production();

// Testing
const testClient = ClientFactory.testing();
```

## Authentication

### Login and Registration

```typescript
// Register a new user
const { user, accessToken, refreshToken } = await client.auth.register({
  email: 'user@example.com',
  password: 'securepassword'
});

// Login existing user
const { user, accessToken, refreshToken } = await client.auth.login({
  email: 'user@example.com',
  password: 'securepassword'
});

// Logout
await client.auth.logout();
```

### Profile Management

```typescript
// Get user profile
const profile = await client.auth.getProfile();

// Update profile
const updatedProfile = await client.auth.updateProfile({
  email: 'newemail@example.com'
});

// Change password
await client.auth.changePassword({
  currentPassword: 'oldpassword',
  newPassword: 'newpassword'
});
```

### Token Management

```typescript
// Check if authenticated
if (client.auth.isAuthenticated()) {
  console.log('User is authenticated');
}

// Set tokens manually
client.auth.setTokens('access-token', 'refresh-token');

// Clear tokens
client.auth.clearTokens();

// Validate current token
const { valid, user } = await client.auth.validateToken();
```

## URL Management

### Creating URLs

```typescript
// Basic URL creation
const url = await client.urls.create({
  originalURL: 'https://example.com'
});

// Advanced URL creation
const url = await client.urls.create({
  originalURL: 'https://example.com/very-long-path',
  title: 'My Link',
  description: 'This is my custom link',
  customAlias: 'my-custom-link',
  password: 'secret123',
  expiresAt: '2024-12-31T23:59:59Z'
});

// Bulk creation
const urls = await client.urls.bulkCreate([
  { originalURL: 'https://example1.com' },
  { originalURL: 'https://example2.com' },
]);
```

### Managing URLs

```typescript
// List user URLs
const { data: urls, total } = await client.urls.list({
  limit: 20,
  offset: 0
});

// Get specific URL
const url = await client.urls.getById(123);

// Update URL
const updatedURL = await client.urls.update(123, {
  title: 'Updated Title',
  isActive: false
});

// Delete URL
await client.urls.delete(123);

// Search URLs
const { data: searchResults } = await client.urls.search('example', {
  limit: 10
});
```

### URL Utilities

```typescript
// Get popular URLs
const popularURLs = await client.urls.getPopular(10);

// Check alias availability
const isAvailable = await client.urls.checkAliasAvailability('my-alias');

// Get URL statistics
const stats = await client.urls.getStats(123);

// Toggle active status
const toggledURL = await client.urls.toggleActive(123);

// Duplicate URL
const duplicatedURL = await client.urls.duplicate(123);
```

## Analytics

### Dashboard Analytics

```typescript
// Get dashboard overview
const dashboard = await client.analytics.getDashboard();

// Get global statistics
const globalStats = await client.analytics.getGlobalStats();

// Get top performing URLs
const topURLs = await client.analytics.getTopURLs(10, '30d');
```

### URL-Specific Analytics

```typescript
// Get detailed URL analytics
const analytics = await client.analytics.getURLAnalytics(123, '7d');

// Get click timeline
const timeline = await client.analytics.getClickTimeline(123, '30d');

// Get geographic stats
const geoStats = await client.analytics.getGeographicStats(123);

// Get device statistics
const deviceStats = await client.analytics.getDeviceStats(123);

// Get referrer statistics
const referrerStats = await client.analytics.getReferrerStats(123);
```

### Advanced Analytics

```typescript
// Real-time statistics
const realTime = await client.analytics.getRealTimeStats();

// Compare periods
const comparison = await client.analytics.getComparison(123, '7d', '7d');

// Bulk analytics
const bulkAnalytics = await client.analytics.getBulkAnalytics([123, 456]);

// Click heatmap
const heatmap = await client.analytics.getClickHeatmap(123, '30d');
```

### Data Export

```typescript
// Export analytics data
const csvBlob = await client.analytics.exportAnalytics('csv', {
  urlIds: [123, 456],
  period: '30d',
  includeDetails: true
});

// Download the file
const url = URL.createObjectURL(csvBlob);
const a = document.createElement('a');
a.href = url;
a.download = 'analytics-export.csv';
a.click();
```

## QR Codes

### Basic QR Code Generation

```typescript
// Generate QR code for any URL
const qrCode = await client.qr.generate({
  url: 'https://example.com',
  options: {
    size: 200,
    format: 'png'
  }
});

// Generate QR code for existing short URL
const qrCode = await client.qr.generateForURL('abc123', {
  size: 300,
  format: 'svg'
});
```

### QR Code Customization

```typescript
// Get available formats and sizes
const formats = await client.qr.getFormats();
const sizes = await client.qr.getSizes();

// Validate options
const validation = await client.qr.validateOptions({
  size: 200,
  format: 'png',
  level: 'H'
});

// Custom QR code with branding
const brandedQR = await client.qr.generateBranded('https://example.com', {
  size: 300,
  backgroundColor: '#ffffff',
  foregroundColor: '#000000',
  logo: 'base64-encoded-logo',
  logoSize: 20
});
```

### QR Code Operations

```typescript
// Get QR code preview
const preview = await client.qr.getPreview('https://example.com', {
  size: 150
});

// Download QR code
const blob = await client.qr.download('abc123', {
  size: 400,
  format: 'png',
  filename: 'my-qr-code.png'
});

// Batch generation
const qrCodes = await client.qr.generateBatch([
  { url: 'https://example1.com' },
  { url: 'https://example2.com' }
]);

// Get QR analytics
const qrAnalytics = await client.qr.getAnalytics('abc123', '30d');
```

## Webhooks

### Webhook Management

```typescript
// Create webhook
const webhook = await client.webhooks.create({
  name: 'My Webhook',
  url: 'https://mysite.com/webhook',
  events: ['url.created', 'url.clicked'],
  secret: 'webhook-secret'
});

// List webhooks
const { data: webhooks } = await client.webhooks.list();

// Update webhook
const updatedWebhook = await client.webhooks.update(123, {
  name: 'Updated Webhook',
  events: ['url.created', 'url.updated', 'url.clicked']
});

// Delete webhook
await client.webhooks.delete(123);
```

### Webhook Operations

```typescript
// Activate/Deactivate webhook
await client.webhooks.activate(123);
await client.webhooks.deactivate(123);

// Test webhook
const testResult = await client.webhooks.test(123);
console.log('Test successful:', testResult.success);

// Get webhook statistics
const stats = await client.webhooks.getStats(123);
console.log('Success rate:', stats.successRate);
```

### Webhook Deliveries

```typescript
// Get webhook deliveries
const { data: deliveries } = await client.webhooks.getDeliveries(123, {
  limit: 50
});

// Get failed deliveries
const failedDeliveries = await client.webhooks.getFailedDeliveries(123);

// Retry failed delivery
await client.webhooks.retryDelivery(456);

// Get delivery details
const delivery = await client.webhooks.getDelivery(456);
```

### Webhook Events and Utilities

```typescript
// Get available events
const events = await client.webhooks.getEvents();

// Validate webhook URL
const { valid, error } = await client.webhooks.validateURL('https://mysite.com/webhook');

// Get webhook health summary
const health = await client.webhooks.getHealthSummary();

// Bulk operations
await client.webhooks.bulkUpdateStatus([123, 456], 'active');
await client.webhooks.bulkDelete([789, 101]);
```

## Error Handling

```typescript
import { 
  SDKError, 
  AuthenticationError, 
  ValidationError, 
  NetworkError,
  RateLimitError 
} from '@urlshortener/sdk';

try {
  const shortURL = await client.urls.create({
    originalURL: 'invalid-url'
  });
} catch (error) {
  if (error instanceof ValidationError) {
    console.error('Validation error:', error.message, error.field);
  } else if (error instanceof AuthenticationError) {
    console.error('Authentication failed:', error.message);
    // Redirect to login
  } else if (error instanceof RateLimitError) {
    console.error('Rate limit exceeded, please try again later');
  } else if (error instanceof NetworkError) {
    console.error('Network error:', error.message);
  } else if (error instanceof SDKError) {
    console.error('SDK error:', error.message, error.code);
  } else {
    console.error('Unexpected error:', error);
  }
}
```

## Utilities

The SDK includes helpful utility functions:

```typescript
import { utils } from '@urlshortener/sdk';

// URL validation
const isValid = utils.isValidURL('https://example.com');

// Generate random alias
const randomAlias = utils.generateRandomAlias(8);

// Format numbers
const formatted = utils.formatNumber(1234567); // "1.2M"

// Calculate percentage
const percentage = utils.calculatePercentage(25, 100); // 25

// Date formatting
const dateStr = utils.formatDate(new Date());

// Debounce API calls
const debouncedSearch = utils.debounce(async (query) => {
  return client.urls.search(query);
}, 300);

// Retry with exponential backoff
const result = await utils.retry(async () => {
  return client.urls.create({ originalURL: 'https://example.com' });
}, 3, 1000);
```

## TypeScript Support

The SDK is written in TypeScript and provides full type definitions:

```typescript
import { 
  ShortURL, 
  URLAnalytics, 
  Webhook, 
  QRCodeOptions,
  CreateURLRequest,
  PaginatedResponse 
} from '@urlshortener/sdk';

// Type-safe API calls
const url: ShortURL = await client.urls.create({
  originalURL: 'https://example.com'
});

const analytics: URLAnalytics = await client.analytics.getURLAnalytics(url.id);

const webhooks: PaginatedResponse<Webhook> = await client.webhooks.list();
```

## Environment Configuration

### Development

```typescript
import { clients } from '@urlshortener/sdk';

const client = clients.development({
  onError: (error) => console.error('Dev Error:', error)
});
```

### Production

```typescript
import { clients } from '@urlshortener/sdk';

const client = clients.production({
  accessToken: process.env.URLSHORTENER_ACCESS_TOKEN,
  onTokenRefresh: (tokens) => {
    // Save to secure storage
  }
});
```

## Browser Usage

```html
<!DOCTYPE html>
<html>
<head>
  <script src="https://unpkg.com/@urlshortener/sdk@1.0.0/dist/index.umd.js"></script>
</head>
<body>
  <script>
    const client = new URLShortenerSDK.URLShortenerClient({
      baseURL: 'https://api.urlshortener.com'
    });
    
    // Use the client
    client.auth.login({ email: 'user@example.com', password: 'password' })
      .then(response => console.log('Logged in:', response));
  </script>
</body>
</html>
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

```bash
# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

## Building

```bash
# Build the SDK
npm run build

# Development build with watch mode
npm run dev
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📚 [Documentation](https://docs.urlshortener.com/sdk/javascript)
- 🐛 [Issue Tracker](https://github.com/yourusername/url-shortener/issues)
- 💬 [Discord Community](https://discord.gg/urlshortener)
- 📧 [Email Support](mailto:support@urlshortener.com)