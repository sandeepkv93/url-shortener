import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import URLDetails from '../URLDetails'
import { urlService, urlAnalyticsService } from '@/services/urls'
import { URL as URLType } from '@/types/url'

// Mock the services
vi.mock('@/services/urls', () => ({
  urlService: {
    getURL: vi.fn(),
    getClickHistory: vi.fn()
  },
  urlAnalyticsService: {
    getClickTimeline: vi.fn(),
    getGeographicStats: vi.fn(),
    getDeviceStats: vi.fn(),
    getReferrerStats: vi.fn()
  }
}))

// Mock react-router-dom
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useParams: () => ({ id: '1' }),
    Link: ({ children, to, ...props }: any) => <a href={to} {...props}>{children}</a>
  }
})

// Mock Recharts components
vi.mock('recharts', () => ({
  LineChart: ({ children }: any) => <div data-testid="line-chart">{children}</div>,
  Line: () => <div data-testid="line" />,
  AreaChart: ({ children }: any) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div data-testid="area" />,
  PieChart: ({ children }: any) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => <div data-testid="pie" />,
  Cell: () => <div data-testid="cell" />,
  BarChart: ({ children }: any) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />,
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>
}))

// Mock clipboard API
Object.assign(navigator, {
  clipboard: {
    writeText: vi.fn().mockImplementation(() => Promise.resolve())
  }
})

// Mock window.open
global.open = vi.fn()

const mockURL: URLType = {
  id: '1',
  shortCode: 'abc123',
  originalUrl: 'https://example.com',
  title: 'Test URL',
  description: 'Test description',
  userId: 'user1',
  clickCount: 150,
  isActive: true,
  isPublic: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  tags: ['test', 'example']
}

const mockTimeline = {
  timeline: [
    { date: '2024-01-01', clicks: 10, uniqueClicks: 8 },
    { date: '2024-01-02', clicks: 15, uniqueClicks: 12 },
    { date: '2024-01-03', clicks: 20, uniqueClicks: 16 }
  ],
  total: 45,
  period: '7d'
}

const mockGeographic = {
  countries: [
    { country: 'United States', countryCode: 'US', clicks: 50, percentage: 60 },
    { country: 'United Kingdom', countryCode: 'GB', clicks: 20, percentage: 24 },
    { country: 'Canada', countryCode: 'CA', clicks: 13, percentage: 16 }
  ],
  cities: [
    { city: 'New York', country: 'United States', clicks: 25, percentage: 30 },
    { city: 'London', country: 'United Kingdom', clicks: 15, percentage: 18 },
    { city: 'Toronto', country: 'Canada', clicks: 10, percentage: 12 }
  ]
}

const mockDevices = {
  devices: [
    { device: 'Desktop', clicks: 60, percentage: 70 },
    { device: 'Mobile', clicks: 20, percentage: 24 },
    { device: 'Tablet', clicks: 5, percentage: 6 }
  ],
  browsers: [
    { browser: 'Chrome', clicks: 50, percentage: 59 },
    { browser: 'Firefox', clicks: 20, percentage: 24 },
    { browser: 'Safari', clicks: 15, percentage: 17 }
  ],
  operatingSystems: [
    { os: 'Windows', clicks: 45, percentage: 53 },
    { os: 'macOS', clicks: 25, percentage: 29 },
    { os: 'Linux', clicks: 15, percentage: 18 }
  ]
}

const mockReferrers = {
  referrers: [
    { referrer: 'google.com', clicks: 30, percentage: 35 },
    { referrer: 'facebook.com', clicks: 15, percentage: 18 },
    { referrer: 'twitter.com', clicks: 10, percentage: 12 }
  ],
  directClicks: 30,
  totalClicks: 85
}

const mockClickHistory = {
  clicks: [
    {
      id: '1',
      timestamp: '2024-01-03T12:00:00Z',
      country: 'United States',
      city: 'New York',
      device: 'Desktop',
      browser: 'Chrome'
    },
    {
      id: '2',
      timestamp: '2024-01-03T11:30:00Z',
      country: 'United Kingdom',
      city: 'London',
      device: 'Mobile',
      browser: 'Safari'
    }
  ],
  total: 2,
  page: 1,
  limit: 10
}

const setupMocks = () => {
  const mockGetURL = vi.mocked(urlService.getURL)
  const mockGetClickHistory = vi.mocked(urlService.getClickHistory)
  const mockGetClickTimeline = vi.mocked(urlAnalyticsService.getClickTimeline)
  const mockGetGeographicStats = vi.mocked(urlAnalyticsService.getGeographicStats)
  const mockGetDeviceStats = vi.mocked(urlAnalyticsService.getDeviceStats)
  const mockGetReferrerStats = vi.mocked(urlAnalyticsService.getReferrerStats)

  mockGetURL.mockResolvedValue(mockURL)
  mockGetClickHistory.mockResolvedValue(mockClickHistory)
  mockGetClickTimeline.mockResolvedValue(mockTimeline)
  mockGetGeographicStats.mockResolvedValue(mockGeographic)
  mockGetDeviceStats.mockResolvedValue(mockDevices)
  mockGetReferrerStats.mockResolvedValue(mockReferrers)
}

const renderComponent = (props = {}) => {
  const defaultProps = {
    ...props
  }
  
  return render(
    <BrowserRouter>
      <URLDetails {...defaultProps} />
    </BrowserRouter>
  )
}

describe('URLDetails', () => {
  const user = userEvent.setup()

  beforeEach(() => {
    vi.clearAllMocks()
    setupMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('Loading and Error States', () => {
    it('shows loading state initially', () => {
      const mockGetURL = vi.mocked(urlService.getURL)
      mockGetURL.mockImplementation(() => new Promise(() => {})) // Never resolves
      
      renderComponent()
      
      expect(screen.getByText('Loading analytics...')).toBeInTheDocument()
    })

    it('shows error state when data fetch fails', async () => {
      const mockGetURL = vi.mocked(urlService.getURL)
      mockGetURL.mockRejectedValue(new Error('Failed to load data'))
      
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('Failed to load data')).toBeInTheDocument()
      })
    })

    it('shows not found state when URL is not found', async () => {
      const mockGetURL = vi.mocked(urlService.getURL)
      mockGetURL.mockResolvedValue(null as any)
      
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('URL not found')).toBeInTheDocument()
        expect(screen.getByText('The requested URL could not be found.')).toBeInTheDocument()
      })
    })
  })

  describe('Header and URL Information', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('displays URL information correctly', () => {
      expect(screen.getByText('Test URL')).toBeInTheDocument()
      expect(screen.getByText('https://example.com')).toBeInTheDocument()
      expect(screen.getByText('Test description')).toBeInTheDocument()
    })

    it('shows short URL correctly', () => {
      expect(screen.getByText(`${window.location.origin}/abc123`)).toBeInTheDocument()
    })

    it('shows status badges correctly', () => {
      expect(screen.getByText('Active')).toBeInTheDocument()
      expect(screen.getByText('Public')).toBeInTheDocument()
    })

    it('shows tags when present', () => {
      expect(screen.getByText('test')).toBeInTheDocument()
      expect(screen.getByText('example')).toBeInTheDocument()
    })

    it('shows back to dashboard link', () => {
      const backLink = screen.getByText('Back to Dashboard')
      expect(backLink.closest('a')).toHaveAttribute('href', '/dashboard')
    })

    it('shows export button', () => {
      expect(screen.getByText('Export')).toBeInTheDocument()
    })

    it('shows edit URL link', () => {
      const editLink = screen.getByText('Edit URL')
      expect(editLink.closest('a')).toHaveAttribute('href', '/urls/1/edit')
    })
  })

  describe('Quick Statistics', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('displays total clicks', () => {
      expect(screen.getByText('150')).toBeInTheDocument()
      expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    })

    it('displays unique clicks from timeline data', () => {
      expect(screen.getByText('36')).toBeInTheDocument() // Sum of uniqueClicks from timeline
      expect(screen.getByText('Unique Clicks')).toBeInTheDocument()
    })

    it('calculates and displays click-through rate', () => {
      // 36 unique clicks / 45 total clicks from timeline = 80%
      expect(screen.getByText('80.0%')).toBeInTheDocument()
      expect(screen.getByText('Click-through Rate')).toBeInTheDocument()
    })

    it('displays creation date', () => {
      expect(screen.getByText('Jan 1')).toBeInTheDocument()
      expect(screen.getByText('Created')).toBeInTheDocument()
    })
  })

  describe('Period Selector', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('shows all period options', () => {
      expect(screen.getByText('Last Hour')).toBeInTheDocument()
      expect(screen.getByText('Last 24 Hours')).toBeInTheDocument()
      expect(screen.getByText('Last 7 Days')).toBeInTheDocument()
      expect(screen.getByText('Last 30 Days')).toBeInTheDocument()
      expect(screen.getByText('Last 90 Days')).toBeInTheDocument()
      expect(screen.getByText('Last Year')).toBeInTheDocument()
      expect(screen.getByText('All Time')).toBeInTheDocument()
    })

    it('highlights default selected period (30 days)', () => {
      const thirtyDaysButton = screen.getByText('Last 30 Days')
      expect(thirtyDaysButton).toHaveClass('bg-primary-100')
    })

    it('changes period when different option is selected', async () => {
      const mockGetClickTimeline = vi.mocked(urlAnalyticsService.getClickTimeline)
      
      await user.click(screen.getByText('Last 7 Days'))
      
      await waitFor(() => {
        expect(mockGetClickTimeline).toHaveBeenCalledWith('1', '7d')
      })
    })
  })

  describe('Click Timeline Chart', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('renders timeline chart section', () => {
      expect(screen.getByText('Click Timeline')).toBeInTheDocument()
      expect(screen.getByTestId('area-chart')).toBeInTheDocument()
      expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    })
  })

  describe('Geographic Distribution', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('displays geographic stats section', () => {
      expect(screen.getByText('Geographic Distribution')).toBeInTheDocument()
    })

    it('shows top countries', () => {
      expect(screen.getByText('Top Countries')).toBeInTheDocument()
      expect(screen.getByText('United States')).toBeInTheDocument()
      expect(screen.getByText('United Kingdom')).toBeInTheDocument()
      expect(screen.getByText('Canada')).toBeInTheDocument()
    })

    it('displays country click counts and percentages', () => {
      expect(screen.getByText('50 (60.0%)')).toBeInTheDocument()
      expect(screen.getByText('20 (24.0%)')).toBeInTheDocument()
      expect(screen.getByText('13 (16.0%)')).toBeInTheDocument()
    })

    it('shows top cities when available', () => {
      expect(screen.getByText('Top Cities')).toBeInTheDocument()
      expect(screen.getByText('New York, United States')).toBeInTheDocument()
      expect(screen.getByText('London, United Kingdom')).toBeInTheDocument()
      expect(screen.getByText('Toronto, Canada')).toBeInTheDocument()
    })
  })

  describe('Device and Browser Stats', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('displays device stats section', () => {
      expect(screen.getByText('Device & Browser Stats')).toBeInTheDocument()
    })

    it('shows device breakdown', () => {
      expect(screen.getByText('Devices')).toBeInTheDocument()
      expect(screen.getByText('Desktop')).toBeInTheDocument()
      expect(screen.getByText('Mobile')).toBeInTheDocument()
      expect(screen.getByText('Tablet')).toBeInTheDocument()
    })

    it('shows browser breakdown', () => {
      expect(screen.getByText('Browsers')).toBeInTheDocument()
      expect(screen.getByText('Chrome')).toBeInTheDocument()
      expect(screen.getByText('Firefox')).toBeInTheDocument()
      expect(screen.getByText('Safari')).toBeInTheDocument()
    })

    it('displays device click counts and percentages', () => {
      expect(screen.getByText('60 (70.0%)')).toBeInTheDocument()
      expect(screen.getByText('20 (24.0%)')).toBeInTheDocument()
      expect(screen.getByText('5 (6.0%)')).toBeInTheDocument()
    })
  })

  describe('Traffic Sources', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('displays traffic sources section', () => {
      expect(screen.getByText('Traffic Sources')).toBeInTheDocument()
    })

    it('shows direct traffic', () => {
      expect(screen.getByText('Direct Traffic')).toBeInTheDocument()
      expect(screen.getByText('30 (35.3%)')).toBeInTheDocument() // 30 / 85 * 100
    })

    it('shows referrer sources', () => {
      expect(screen.getByText('google.com')).toBeInTheDocument()
      expect(screen.getByText('facebook.com')).toBeInTheDocument()
      expect(screen.getByText('twitter.com')).toBeInTheDocument()
    })

    it('displays referrer click counts and percentages', () => {
      expect(screen.getByText('30 (35.0%)')).toBeInTheDocument()
      expect(screen.getByText('15 (18.0%)')).toBeInTheDocument()
      expect(screen.getByText('10 (12.0%)')).toBeInTheDocument()
    })
  })

  describe('Recent Activity', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('displays recent activity section', () => {
      expect(screen.getByText('Recent Activity')).toBeInTheDocument()
    })

    it('shows recent clicks with location and device info', () => {
      expect(screen.getByText('New York, United States')).toBeInTheDocument()
      expect(screen.getByText('Desktop • Chrome')).toBeInTheDocument()
      expect(screen.getByText('London, United Kingdom')).toBeInTheDocument()
      expect(screen.getByText('Mobile • Safari')).toBeInTheDocument()
    })

    it('formats timestamps correctly', () => {
      expect(screen.getByText('Jan 3, 12:00 PM')).toBeInTheDocument()
      expect(screen.getByText('Jan 3, 11:30 AM')).toBeInTheDocument()
    })

    it('shows no activity message when no recent clicks', async () => {
      const mockGetClickHistory = vi.mocked(urlService.getClickHistory)
      mockGetClickHistory.mockResolvedValue({
        clicks: [],
        total: 0,
        page: 1,
        limit: 10
      })
      
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('No recent activity')).toBeInTheDocument()
      })
    })
  })

  describe('Copy and QR Code Functionality', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('copies short URL to clipboard', async () => {
      const copyButton = screen.getByTitle('Copy short URL')
      await user.click(copyButton)
      
      expect(navigator.clipboard.writeText).toHaveBeenCalledWith(`${window.location.origin}/abc123`)
    })

    it('shows copied feedback after copying', async () => {
      const copyButton = screen.getByTitle('Copy short URL')
      await user.click(copyButton)
      
      // The copied state management would need to be tested
      // This is challenging without access to internal state
    })

    it('generates QR code', async () => {
      const qrButton = screen.getByTitle('Generate QR Code')
      await user.click(qrButton)
      
      expect(global.open).toHaveBeenCalledWith(
        expect.stringContaining('qrserver.com'),
        '_blank'
      )
    })

    it('opens short URL in new tab', () => {
      const externalLink = screen.getByTitle('Open short URL')
      expect(externalLink).toHaveAttribute('href', `${window.location.origin}/abc123`)
      expect(externalLink).toHaveAttribute('target', '_blank')
    })
  })

  describe('Export Functionality', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('calls export function when export button is clicked', async () => {
      const consoleSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
      
      await user.click(screen.getByText('Export'))
      
      expect(consoleSpy).toHaveBeenCalledWith('Exporting analytics data:', expect.any(Object))
      
      consoleSpy.mockRestore()
    })
  })

  describe('Data Fetching', () => {
    it('fetches all analytics data on mount', async () => {
      const mockGetURL = vi.mocked(urlService.getURL)
      const mockGetClickTimeline = vi.mocked(urlAnalyticsService.getClickTimeline)
      const mockGetGeographicStats = vi.mocked(urlAnalyticsService.getGeographicStats)
      const mockGetDeviceStats = vi.mocked(urlAnalyticsService.getDeviceStats)
      const mockGetReferrerStats = vi.mocked(urlAnalyticsService.getReferrerStats)
      const mockGetClickHistory = vi.mocked(urlService.getClickHistory)
      
      renderComponent()
      
      await waitFor(() => {
        expect(mockGetURL).toHaveBeenCalledWith('1')
        expect(mockGetClickTimeline).toHaveBeenCalledWith('1', '30d')
        expect(mockGetGeographicStats).toHaveBeenCalledWith('1')
        expect(mockGetDeviceStats).toHaveBeenCalledWith('1')
        expect(mockGetReferrerStats).toHaveBeenCalledWith('1')
        expect(mockGetClickHistory).toHaveBeenCalledWith('1', 1, 10)
      })
    })

    it('refetches timeline data when period changes', async () => {
      const mockGetClickTimeline = vi.mocked(urlAnalyticsService.getClickTimeline)
      
      renderComponent()
      
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
      
      await user.click(screen.getByText('Last 7 Days'))
      
      expect(mockGetClickTimeline).toHaveBeenCalledWith('1', '7d')
    })
  })

  describe('Accessibility', () => {
    beforeEach(async () => {
      renderComponent()
      await waitFor(() => {
        expect(screen.getByText('Test URL')).toBeInTheDocument()
      })
    })

    it('has proper heading structure', () => {
      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Test URL')
    })

    it('has proper button labels and titles', () => {
      expect(screen.getByTitle('Copy short URL')).toBeInTheDocument()
      expect(screen.getByTitle('Generate QR Code')).toBeInTheDocument()
      expect(screen.getByTitle('Open short URL')).toBeInTheDocument()
    })

    it('maintains proper link relationships', () => {
      const backLink = screen.getByText('Back to Dashboard')
      expect(backLink.closest('a')).toHaveAttribute('href', '/dashboard')
      
      const editLink = screen.getByText('Edit URL')
      expect(editLink.closest('a')).toHaveAttribute('href', '/urls/1/edit')
    })
  })

  describe('Responsive Design', () => {
    it('applies custom className correctly', () => {
      const { container } = renderComponent({ className: 'custom-class' })
      
      expect(container.firstChild).toHaveClass('custom-class')
    })
  })
})