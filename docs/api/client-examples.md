# API Client Examples

This document provides code examples for integrating with the URL Shortener API in various programming languages.

## Table of Contents

- [Authentication](#authentication)
- [JavaScript/TypeScript](#javascripttypescript)
- [Python](#python)
- [Go](#go)
- [PHP](#php)
- [Java](#java)
- [cURL](#curl)
- [Error Handling](#error-handling)

## Authentication

All authenticated requests require a JWT Bearer token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Getting a Token

First, register or login to get your JWT token:

```bash
# Register
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

## JavaScript/TypeScript

### Setup

```bash
npm install axios
```

### Client Class

```typescript
import axios, { AxiosInstance, AxiosResponse } from 'axios';

interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  is_active: boolean;
  created_at: string;
}

interface ShortURL {
  id: number;
  short_code: string;
  original_url: string;
  short_url: string;
  custom_alias: boolean;
  expires_at?: string;
  is_active: boolean;
  click_count: number;
  created_at: string;
  updated_at: string;
}

interface CreateURLRequest {
  original_url: string;
  custom_alias?: string;
  title?: string;
  description?: string;
  password?: string;
  expires_at?: string;
}

class URLShortenerClient {
  private client: AxiosInstance;
  private accessToken: string | null = null;
  private refreshToken: string | null = null;

  constructor(baseURL: string = 'http://localhost:8080/api/v1') {
    this.client = axios.create({
      baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.client.interceptors.request.use((config) => {
      if (this.accessToken) {
        config.headers.Authorization = `Bearer ${this.accessToken}`;
      }
      return config;
    });

    // Response interceptor for token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401 && this.refreshToken) {
          try {
            await this.refreshAccessToken();
            // Retry the original request
            return this.client.request(error.config);
          } catch (refreshError) {
            // Refresh failed, redirect to login
            this.clearTokens();
            throw refreshError;
          }
        }
        throw error;
      }
    );
  }

  // Authentication methods
  async register(
    email: string,
    password: string,
    firstName: string,
    lastName: string
  ): Promise<AuthResponse> {
    const response = await this.client.post<AuthResponse>('/auth/register', {
      email,
      password,
      first_name: firstName,
      last_name: lastName,
    });
    
    this.setTokens(response.data.access_token, response.data.refresh_token);
    return response.data;
  }

  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await this.client.post<AuthResponse>('/auth/login', {
      email,
      password,
    });
    
    this.setTokens(response.data.access_token, response.data.refresh_token);
    return response.data;
  }

  async logout(): Promise<void> {
    try {
      await this.client.post('/auth/logout');
    } finally {
      this.clearTokens();
    }
  }

  async refreshAccessToken(): Promise<void> {
    if (!this.refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await this.client.post('/auth/refresh', {
      refresh_token: this.refreshToken,
    });

    this.setTokens(response.data.access_token, response.data.refresh_token);
  }

  // URL management methods
  async createShortURL(request: CreateURLRequest): Promise<ShortURL> {
    const response = await this.client.post<ShortURL>('/urls', request);
    return response.data;
  }

  async getURLs(page: number = 1, pageSize: number = 10): Promise<{
    urls: ShortURL[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  }> {
    const response = await this.client.get('/urls', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  }

  async getURL(id: number): Promise<ShortURL> {
    const response = await this.client.get<ShortURL>(`/urls/${id}`);
    return response.data;
  }

  async updateURL(id: number, updates: Partial<CreateURLRequest>): Promise<ShortURL> {
    const response = await this.client.put<ShortURL>(`/urls/${id}`, updates);
    return response.data;
  }

  async deleteURL(id: number): Promise<void> {
    await this.client.delete(`/urls/${id}`);
  }

  // Analytics methods
  async getDashboard(period: string = 'week'): Promise<any> {
    const response = await this.client.get('/analytics/dashboard', {
      params: { period },
    });
    return response.data;
  }

  async getURLAnalytics(id: number, period: string = 'month'): Promise<any> {
    const response = await this.client.get(`/analytics/urls/${id}`, {
      params: { period },
    });
    return response.data;
  }

  // QR Code methods
  async generateQRCode(
    url: string,
    format: string = 'png',
    size: number = 300
  ): Promise<Blob> {
    const response = await this.client.post('/qr/generate', {
      url,
      format,
      size,
    }, {
      responseType: 'blob',
    });
    return response.data;
  }

  // Token management
  private setTokens(accessToken: string, refreshToken: string): void {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    
    // Store in localStorage for persistence
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('access_token', accessToken);
      localStorage.setItem('refresh_token', refreshToken);
    }
  }

  private clearTokens(): void {
    this.accessToken = null;
    this.refreshToken = null;
    
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  }

  // Load tokens from storage
  loadTokensFromStorage(): void {
    if (typeof localStorage !== 'undefined') {
      this.accessToken = localStorage.getItem('access_token');
      this.refreshToken = localStorage.getItem('refresh_token');
    }
  }
}

// Usage example
async function example() {
  const client = new URLShortenerClient();
  
  try {
    // Login
    const authResponse = await client.login('user@example.com', 'password123');
    console.log('Logged in:', authResponse.user);

    // Create a short URL
    const shortURL = await client.createShortURL({
      original_url: 'https://www.example.com/very/long/url',
      custom_alias: 'my-link',
      title: 'My Important Link',
    });
    console.log('Created URL:', shortURL);

    // Get analytics
    const analytics = await client.getURLAnalytics(shortURL.id);
    console.log('Analytics:', analytics);

    // Generate QR code
    const qrBlob = await client.generateQRCode(shortURL.short_url);
    console.log('QR Code generated:', qrBlob);

  } catch (error) {
    console.error('Error:', error);
  }
}

export { URLShortenerClient, example };
```

## Python

### Setup

```bash
pip install requests
```

### Client Class

```python
import requests
from typing import Optional, Dict, Any, List
from datetime import datetime
import json

class URLShortenerClient:
    def __init__(self, base_url: str = "http://localhost:8080/api/v1"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json'
        })
        self.access_token = None
        self.refresh_token = None

    def _make_request(self, method: str, endpoint: str, **kwargs) -> requests.Response:
        """Make HTTP request with automatic token refresh."""
        url = f"{self.base_url}{endpoint}"
        
        if self.access_token:
            self.session.headers['Authorization'] = f"Bearer {self.access_token}"
        
        response = self.session.request(method, url, **kwargs)
        
        # Handle token refresh
        if response.status_code == 401 and self.refresh_token:
            self.refresh_access_token()
            # Retry request
            if self.access_token:
                self.session.headers['Authorization'] = f"Bearer {self.access_token}"
            response = self.session.request(method, url, **kwargs)
        
        return response

    def _handle_response(self, response: requests.Response) -> Dict[str, Any]:
        """Handle API response and raise exceptions for errors."""
        if response.status_code >= 400:
            try:
                error_data = response.json()
                raise Exception(f"API Error: {error_data.get('message', 'Unknown error')}")
            except json.JSONDecodeError:
                raise Exception(f"HTTP {response.status_code}: {response.text}")
        
        return response.json() if response.content else {}

    # Authentication methods
    def register(self, email: str, password: str, first_name: str, last_name: str) -> Dict[str, Any]:
        """Register a new user account."""
        response = self._make_request('POST', '/auth/register', json={
            'email': email,
            'password': password,
            'first_name': first_name,
            'last_name': last_name
        })
        
        data = self._handle_response(response)
        self._set_tokens(data['access_token'], data['refresh_token'])
        return data

    def login(self, email: str, password: str) -> Dict[str, Any]:
        """Login with email and password."""
        response = self._make_request('POST', '/auth/login', json={
            'email': email,
            'password': password
        })
        
        data = self._handle_response(response)
        self._set_tokens(data['access_token'], data['refresh_token'])
        return data

    def logout(self) -> None:
        """Logout and clear tokens."""
        try:
            self._make_request('POST', '/auth/logout')
        finally:
            self._clear_tokens()

    def refresh_access_token(self) -> None:
        """Refresh the access token using refresh token."""
        if not self.refresh_token:
            raise Exception("No refresh token available")
        
        response = self._make_request('POST', '/auth/refresh', json={
            'refresh_token': self.refresh_token
        })
        
        data = self._handle_response(response)
        self._set_tokens(data['access_token'], data['refresh_token'])

    # URL management methods
    def create_short_url(self, original_url: str, custom_alias: Optional[str] = None,
                        title: Optional[str] = None, description: Optional[str] = None,
                        password: Optional[str] = None, expires_at: Optional[str] = None) -> Dict[str, Any]:
        """Create a new short URL."""
        payload = {'original_url': original_url}
        
        if custom_alias:
            payload['custom_alias'] = custom_alias
        if title:
            payload['title'] = title
        if description:
            payload['description'] = description
        if password:
            payload['password'] = password
        if expires_at:
            payload['expires_at'] = expires_at
        
        response = self._make_request('POST', '/urls', json=payload)
        return self._handle_response(response)

    def get_urls(self, page: int = 1, page_size: int = 10, status: Optional[str] = None,
                search: Optional[str] = None) -> Dict[str, Any]:
        """Get user's short URLs with pagination."""
        params = {'page': page, 'page_size': page_size}
        
        if status:
            params['status'] = status
        if search:
            params['search'] = search
        
        response = self._make_request('GET', '/urls', params=params)
        return self._handle_response(response)

    def get_url(self, url_id: int) -> Dict[str, Any]:
        """Get details of a specific URL."""
        response = self._make_request('GET', f'/urls/{url_id}')
        return self._handle_response(response)

    def update_url(self, url_id: int, **kwargs) -> Dict[str, Any]:
        """Update a short URL."""
        response = self._make_request('PUT', f'/urls/{url_id}', json=kwargs)
        return self._handle_response(response)

    def delete_url(self, url_id: int) -> None:
        """Delete a short URL."""
        response = self._make_request('DELETE', f'/urls/{url_id}')
        self._handle_response(response)

    # Analytics methods
    def get_dashboard(self, period: str = 'week') -> Dict[str, Any]:
        """Get user dashboard analytics."""
        response = self._make_request('GET', '/analytics/dashboard', params={'period': period})
        return self._handle_response(response)

    def get_url_analytics(self, url_id: int, period: str = 'month') -> Dict[str, Any]:
        """Get analytics for a specific URL."""
        response = self._make_request('GET', f'/analytics/urls/{url_id}', params={'period': period})
        return self._handle_response(response)

    def get_geographic_stats(self, url_id: int, level: str = 'country') -> Dict[str, Any]:
        """Get geographic statistics for a URL."""
        response = self._make_request('GET', f'/analytics/urls/{url_id}/geo', params={'level': level})
        return self._handle_response(response)

    def export_analytics(self, format: str = 'csv', period: str = 'month') -> bytes:
        """Export analytics data."""
        response = self._make_request('GET', '/analytics/export', 
                                    params={'format': format, 'period': period})
        if response.status_code == 200:
            return response.content
        else:
            self._handle_response(response)

    # QR Code methods
    def generate_qr_code(self, url: str, format: str = 'png', size: int = 300) -> bytes:
        """Generate QR code for a URL."""
        response = self._make_request('POST', '/qr/generate', json={
            'url': url,
            'format': format,
            'size': size
        })
        
        if response.status_code == 200:
            return response.content
        else:
            self._handle_response(response)

    # Token management
    def _set_tokens(self, access_token: str, refresh_token: str) -> None:
        """Set authentication tokens."""
        self.access_token = access_token
        self.refresh_token = refresh_token

    def _clear_tokens(self) -> None:
        """Clear authentication tokens."""
        self.access_token = None
        self.refresh_token = None
        if 'Authorization' in self.session.headers:
            del self.session.headers['Authorization']

# Usage example
def example():
    client = URLShortenerClient()
    
    try:
        # Login
        auth_response = client.login('user@example.com', 'password123')
        print(f"Logged in: {auth_response['user']}")

        # Create a short URL
        short_url = client.create_short_url(
            original_url='https://www.example.com/very/long/url',
            custom_alias='my-link',
            title='My Important Link'
        )
        print(f"Created URL: {short_url}")

        # Get analytics
        analytics = client.get_url_analytics(short_url['id'])
        print(f"Analytics: {analytics}")

        # Generate QR code
        qr_data = client.generate_qr_code(short_url['short_url'])
        with open('qr_code.png', 'wb') as f:
            f.write(qr_data)
        print("QR code saved as qr_code.png")

    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    example()
```

## Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

type Client struct {
    baseURL      string
    httpClient   *http.Client
    accessToken  string
    refreshToken string
}

type AuthResponse struct {
    User         User   `json:"user"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}

type ShortURL struct {
    ID          int        `json:"id"`
    ShortCode   string     `json:"short_code"`
    OriginalURL string     `json:"original_url"`
    ShortURL    string     `json:"short_url"`
    CustomAlias bool       `json:"custom_alias"`
    ExpiresAt   *time.Time `json:"expires_at"`
    IsActive    bool       `json:"is_active"`
    ClickCount  int64      `json:"click_count"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateURLRequest struct {
    OriginalURL string `json:"original_url"`
    CustomAlias string `json:"custom_alias,omitempty"`
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Password    string `json:"password,omitempty"`
    ExpiresAt   string `json:"expires_at,omitempty"`
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader
    
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewReader(jsonData)
    }
    
    req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    if c.accessToken != "" {
        req.Header.Set("Authorization", "Bearer "+c.accessToken)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    
    // Handle token refresh
    if resp.StatusCode == 401 && c.refreshToken != "" {
        resp.Body.Close()
        if err := c.RefreshToken(); err == nil {
            // Retry request
            req.Header.Set("Authorization", "Bearer "+c.accessToken)
            return c.httpClient.Do(req)
        }
    }
    
    return resp, nil
}

func (c *Client) Register(email, password, firstName, lastName string) (*AuthResponse, error) {
    payload := map[string]string{
        "email":      email,
        "password":   password,
        "first_name": firstName,
        "last_name":  lastName,
    }
    
    resp, err := c.makeRequest("POST", "/auth/register", payload)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var authResp AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        return nil, err
    }
    
    c.accessToken = authResp.AccessToken
    c.refreshToken = authResp.RefreshToken
    
    return &authResp, nil
}

func (c *Client) Login(email, password string) (*AuthResponse, error) {
    payload := map[string]string{
        "email":    email,
        "password": password,
    }
    
    resp, err := c.makeRequest("POST", "/auth/login", payload)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var authResp AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        return nil, err
    }
    
    c.accessToken = authResp.AccessToken
    c.refreshToken = authResp.RefreshToken
    
    return &authResp, nil
}

func (c *Client) RefreshToken() error {
    payload := map[string]string{
        "refresh_token": c.refreshToken,
    }
    
    resp, err := c.makeRequest("POST", "/auth/refresh", payload)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var tokenResp struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return err
    }
    
    c.accessToken = tokenResp.AccessToken
    c.refreshToken = tokenResp.RefreshToken
    
    return nil
}

func (c *Client) CreateShortURL(req CreateURLRequest) (*ShortURL, error) {
    resp, err := c.makeRequest("POST", "/urls", req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var shortURL ShortURL
    if err := json.NewDecoder(resp.Body).Decode(&shortURL); err != nil {
        return nil, err
    }
    
    return &shortURL, nil
}

func (c *Client) GetURLs(page, pageSize int) ([]ShortURL, error) {
    params := url.Values{}
    params.Add("page", fmt.Sprintf("%d", page))
    params.Add("page_size", fmt.Sprintf("%d", pageSize))
    
    endpoint := "/urls?" + params.Encode()
    resp, err := c.makeRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var response struct {
        URLs []ShortURL `json:"urls"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }
    
    return response.URLs, nil
}

func (c *Client) GetURLAnalytics(urlID int, period string) (map[string]interface{}, error) {
    params := url.Values{}
    params.Add("period", period)
    
    endpoint := fmt.Sprintf("/analytics/urls/%d?%s", urlID, params.Encode())
    resp, err := c.makeRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var analytics map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
        return nil, err
    }
    
    return analytics, nil
}

// Usage example
func main() {
    client := NewClient("http://localhost:8080/api/v1")
    
    // Login
    authResp, err := client.Login("user@example.com", "password123")
    if err != nil {
        fmt.Printf("Login failed: %v\n", err)
        return
    }
    
    fmt.Printf("Logged in: %+v\n", authResp.User)
    
    // Create short URL
    shortURL, err := client.CreateShortURL(CreateURLRequest{
        OriginalURL: "https://www.example.com/very/long/url",
        CustomAlias: "my-link",
        Title:       "My Important Link",
    })
    if err != nil {
        fmt.Printf("Failed to create URL: %v\n", err)
        return
    }
    
    fmt.Printf("Created URL: %+v\n", shortURL)
    
    // Get analytics
    analytics, err := client.GetURLAnalytics(shortURL.ID, "month")
    if err != nil {
        fmt.Printf("Failed to get analytics: %v\n", err)
        return
    }
    
    fmt.Printf("Analytics: %+v\n", analytics)
}
```

## cURL

### Authentication

```bash
# Register
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Save token for subsequent requests
TOKEN="your-jwt-token-here"
```

### URL Management

```bash
# Create short URL
curl -X POST "http://localhost:8080/api/v1/urls" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://www.example.com/very/long/url",
    "custom_alias": "my-link",
    "title": "My Important Link"
  }'

# Get user URLs
curl -X GET "http://localhost:8080/api/v1/urls?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# Get specific URL
curl -X GET "http://localhost:8080/api/v1/urls/123" \
  -H "Authorization: Bearer $TOKEN"

# Update URL
curl -X PUT "http://localhost:8080/api/v1/urls/123" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "is_active": true
  }'

# Delete URL
curl -X DELETE "http://localhost:8080/api/v1/urls/123" \
  -H "Authorization: Bearer $TOKEN"
```

### Analytics

```bash
# Get dashboard
curl -X GET "http://localhost:8080/api/v1/analytics/dashboard?period=week" \
  -H "Authorization: Bearer $TOKEN"

# Get URL analytics
curl -X GET "http://localhost:8080/api/v1/analytics/urls/123?period=month" \
  -H "Authorization: Bearer $TOKEN"

# Get geographic stats
curl -X GET "http://localhost:8080/api/v1/analytics/urls/123/geo?level=country" \
  -H "Authorization: Bearer $TOKEN"

# Export analytics
curl -X GET "http://localhost:8080/api/v1/analytics/export?format=csv&period=month" \
  -H "Authorization: Bearer $TOKEN" \
  -o analytics.csv
```

### QR Codes

```bash
# Generate QR code
curl -X POST "http://localhost:8080/api/v1/qr/generate" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://short.ly/abc123",
    "format": "png",
    "size": 300
  }' \
  -o qrcode.png

# Get QR code for existing URL
curl -X GET "http://localhost:8080/api/v1/qr/abc123?format=png&size=300" \
  -o qrcode.png
```

## Error Handling

All clients should handle the following HTTP status codes:

- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required or invalid token
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `422 Unprocessable Entity` - Validation errors
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

### Example Error Response

```json
{
  "error": "Validation failed",
  "message": "Request validation failed",
  "code": 422,
  "details": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ],
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Rate Limiting

Check these headers for rate limit information:

- `X-RateLimit-Limit` - Request limit per window
- `X-RateLimit-Remaining` - Remaining requests in current window
- `X-RateLimit-Reset` - Unix timestamp when the window resets
- `Retry-After` - Seconds to wait before retrying (when rate limited)

## Best Practices

1. **Token Management**: Always handle token refresh automatically
2. **Error Handling**: Implement proper error handling for all status codes
3. **Rate Limiting**: Respect rate limits and implement backoff strategies
4. **Timeouts**: Set appropriate request timeouts
5. **Logging**: Log API requests for debugging (but never log tokens)
6. **Security**: Store tokens securely and never expose them in logs or URLs
7. **Validation**: Validate input data before sending requests
8. **Pagination**: Handle paginated responses properly for large datasets

## Support

For additional help:

- [Interactive API Documentation](swagger-ui.html)
- [OpenAPI Specification](openapi.yaml)
- [GitHub Issues](#) (replace with actual URL)