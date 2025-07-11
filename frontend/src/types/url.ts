// URL data types
export interface URL {
  id: string
  shortCode: string
  originalUrl: string
  title?: string
  description?: string
  customAlias?: string
  domain?: string
  userId: string
  clickCount: number
  isActive: boolean
  isPublic: boolean
  expiresAt?: string
  createdAt: string
  updatedAt: string
  lastClickedAt?: string
  password?: string
  tags?: string[]
}

// URL creation and update types
export interface CreateURLRequest {
  originalUrl: string
  customAlias?: string
  title?: string
  description?: string
  expiresAt?: string
  password?: string
  isPublic?: boolean
  tags?: string[]
}

export interface UpdateURLRequest {
  title?: string
  description?: string
  isActive?: boolean
  isPublic?: boolean
  expiresAt?: string
  password?: string
  tags?: string[]
}

export interface BulkUpdateRequest {
  urlIds: string[]
  updates: Partial<UpdateURLRequest>
}

// URL query and filter types
export interface URLFilter {
  search?: string
  isActive?: boolean
  isPublic?: boolean
  hasExpiry?: boolean
  hasPassword?: boolean
  tags?: string[]
  sortBy?: 'createdAt' | 'updatedAt' | 'clickCount' | 'title'
  sortOrder?: 'asc' | 'desc'
  page?: number
  limit?: number
}

export interface URLListResponse {
  urls: URL[]
  total: number
  page: number
  limit: number
  hasNext: boolean
  hasPrev: boolean
}

// URL validation types
export interface URLValidation {
  isValid: boolean
  error?: string
  suggestions?: string[]
}

// URL stats types
export interface URLStats {
  url: URL
  totalClicks: number
  uniqueClicks: number
  clicksByDate: { [date: string]: number }
  clicksByCountry: { [country: string]: number }
  clicksByDevice: { [device: string]: number }
  clicksByBrowser: { [browser: string]: number }
  clicksByReferrer: { [referrer: string]: number }
}

// URL analytics time periods
export type AnalyticsPeriod = '1h' | '24h' | '7d' | '30d' | '90d' | '1y' | 'all'

export interface AnalyticsFilter {
  period: AnalyticsPeriod
  startDate?: string
  endDate?: string
}

// Popular URLs response
export interface PopularURL {
  id: string
  shortCode: string
  originalUrl: string
  title?: string
  clickCount: number
  createdAt: string
}

export interface PopularURLsResponse {
  urls: PopularURL[]
  total: number
}

// URL redirection types
export interface RedirectData {
  userAgent?: string
  referer?: string
  ipAddress?: string
  country?: string
  city?: string
  device?: string
  browser?: string
}

// URL import/export types
export interface ImportURLData {
  originalUrl: string
  customAlias?: string
  title?: string
  tags?: string[]
}

export interface ExportURLData {
  shortCode: string
  originalUrl: string
  title?: string
  clickCount: number
  createdAt: string
  isActive: boolean
}

export interface BulkImportRequest {
  urls: ImportURLData[]
  overwriteExisting?: boolean
}

export interface BulkImportResponse {
  success: number
  failed: number
  errors: Array<{
    index: number
    error: string
    url: ImportURLData
  }>
  imported: URL[]
}

// URL sharing types
export interface ShareableURL {
  shortCode: string
  originalUrl: string
  shortUrl: string
  qrCodeUrl?: string
  title?: string
  description?: string
}

// Error types specific to URLs
export interface URLError {
  code: 'URL_NOT_FOUND' | 'URL_EXPIRED' | 'URL_INACTIVE' | 'URL_PASSWORD_REQUIRED' | 'INVALID_URL' | 'CUSTOM_ALIAS_TAKEN'
  message: string
  field?: string
}

// URL preview types
export interface URLPreview {
  url: string
  title?: string
  description?: string
  image?: string
  favicon?: string
  siteName?: string
}

// Domain types
export interface CustomDomain {
  id: string
  domain: string
  isVerified: boolean
  isActive: boolean
  userId: string
  createdAt: string
  updatedAt: string
}

export interface DomainVerification {
  domain: string
  verificationMethod: 'dns' | 'file'
  verificationToken: string
  instructions: string
}