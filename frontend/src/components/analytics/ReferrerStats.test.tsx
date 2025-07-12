import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest'
import ReferrerStats from './ReferrerStats'
import { useRealTimeAnalytics } from '@/hooks/useRealTimeAnalytics'

// Mock the real-time analytics hook
vi.mock('@/hooks/useRealTimeAnalytics')

// Mock recharts components
vi.mock('recharts', () => ({
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
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>,
  Treemap: ({ children }: any) => <div data-testid="treemap">{children}</div>
}))

const mockReferrerData = {
  referrers: {
    referrers: [
      { referrer: 'google.com', clicks: 450, percentage: 30.0, uniqueClicks: 360, category: 'search', trend: 5.2 },
      { referrer: 'facebook.com', clicks: 300, percentage: 20.0, uniqueClicks: 240, category: 'social', trend: 12.8 },
      { referrer: 'twitter.com', clicks: 200, percentage: 13.3, uniqueClicks: 160, category: 'social', trend: -2.1 },
      { referrer: 'linkedin.com', clicks: 150, percentage: 10.0, uniqueClicks: 120, category: 'social', trend: 8.5 },
      { referrer: 'bing.com', clicks: 100, percentage: 6.7, uniqueClicks: 80, category: 'search', trend: -1.2 }
    ],
    directClicks: 300,
    totalClicks: 1500,
    categories: [
      { category: 'Search', clicks: 550, percentage: 36.7, color: '#3B82F6' },
      { category: 'Social', clicks: 650, percentage: 43.3, color: '#10B981' },
      { category: 'Direct', clicks: 300, percentage: 20.0, color: '#F59E0B' }
    ],
    topDomains: [
      { domain: 'google.com', clicks: 450, percentage: 30.0, sources: ['google.com'] },
      { domain: 'facebook.com', clicks: 300, percentage: 20.0, sources: ['facebook.com'] }
    ]
  }
}

describe('ReferrerStats', () => {
  const mockUseRealTimeAnalytics = vi.mocked(useRealTimeAnalytics)

  beforeEach(() => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockReferrerData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('renders referrer stats header correctly', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('Traffic Source Analytics')).toBeInTheDocument()
    expect(screen.getByText('Analysis of referrer sources and traffic channels')).toBeInTheDocument()
  })

  it('displays loading state', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: null,
      isLoading: true,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: undefined,
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ReferrerStats />)
    expect(screen.getByText('Loading referrer analytics...')).toBeInTheDocument()
  })

  it('displays error state', () => {
    const errorMessage = 'Failed to load referrer data'
    mockUseRealTimeAnalytics.mockReturnValue({
      data: null,
      isLoading: false,
      isRefreshing: false,
      error: errorMessage,
      refreshData: vi.fn(),
      connectionStatus: 'disconnected',
      lastUpdated: undefined,
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ReferrerStats />)
    expect(screen.getByText(errorMessage)).toBeInTheDocument()
  })

  it('renders traffic source selector buttons', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('All Sources')).toBeInTheDocument()
    expect(screen.getByText('Search Engines')).toBeInTheDocument()
    expect(screen.getByText('Social Media')).toBeInTheDocument()
    expect(screen.getByText('Direct Traffic')).toBeInTheDocument()
    expect(screen.getByText('Email')).toBeInTheDocument()
    expect(screen.getByText('Other')).toBeInTheDocument()
  })

  it('changes view mode when traffic source button is clicked', () => {
    render(<ReferrerStats />)
    
    const socialButton = screen.getByText('Social Media')
    fireEvent.click(socialButton)
    
    expect(socialButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders chart type selector', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('Pie Chart')).toBeInTheDocument()
    expect(screen.getByText('Bar Chart')).toBeInTheDocument()
    expect(screen.getByText('Treemap')).toBeInTheDocument()
    expect(screen.getByText('Table')).toBeInTheDocument()
  })

  it('changes chart type when chart type button is clicked', () => {
    render(<ReferrerStats />)
    
    const barChartButton = screen.getByText('Bar Chart')
    fireEvent.click(barChartButton)
    
    expect(barChartButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders search input', () => {
    render(<ReferrerStats />)
    
    const searchInput = screen.getByPlaceholderText('Search referrers...')
    expect(searchInput).toBeInTheDocument()
  })

  it('filters data based on search input', () => {
    render(<ReferrerStats />)
    
    const searchInput = screen.getByPlaceholderText('Search referrers...')
    fireEvent.change(searchInput, { target: { value: 'google' } })
    
    expect(searchInput).toHaveValue('google')
  })

  it('renders sort selector', () => {
    render(<ReferrerStats />)
    
    const sortSelect = screen.getByDisplayValue('Most Clicks')
    expect(sortSelect).toBeInTheDocument()
    
    expect(screen.getByText('Alphabetical')).toBeInTheDocument()
    expect(screen.getByText('Trending')).toBeInTheDocument()
  })

  it('changes sort order when sort option is selected', () => {
    render(<ReferrerStats />)
    
    const sortSelect = screen.getByDisplayValue('Most Clicks')
    fireEvent.change(sortSelect, { target: { value: 'alphabetical' } })
    
    expect(sortSelect).toHaveValue('alphabetical')
  })

  it('displays summary statistics', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Direct Traffic')).toBeInTheDocument()
    expect(screen.getByText('Unique Referrers')).toBeInTheDocument()
    expect(screen.getByText('Top Referrer')).toBeInTheDocument()
  })

  it('displays refresh button and handles refresh action', () => {
    const mockRefreshData = vi.fn()
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockReferrerData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: mockRefreshData,
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ReferrerStats />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    fireEvent.click(refreshButton!)
    
    expect(mockRefreshData).toHaveBeenCalledTimes(1)
  })

  it('shows connection status indicators', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockReferrerData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'reconnecting',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ReferrerStats />)
    expect(screen.getByText('Reconnecting...')).toBeInTheDocument()
  })

  it('handles export functionality', () => {
    // Mock URL.createObjectURL and related functions
    global.URL.createObjectURL = vi.fn(() => 'mock-url')
    global.URL.revokeObjectURL = vi.fn()
    
    // Mock document.createElement
    const mockAnchor = {
      href: '',
      download: '',
      click: vi.fn()
    }
    vi.spyOn(document, 'createElement').mockReturnValue(mockAnchor as any)

    render(<ReferrerStats />)
    
    const exportButton = screen.getByText('Export').closest('button')
    fireEvent.click(exportButton!)
    
    expect(mockAnchor.click).toHaveBeenCalledTimes(1)
    expect(global.URL.createObjectURL).toHaveBeenCalledTimes(1)
    expect(global.URL.revokeObjectURL).toHaveBeenCalledWith('mock-url')
  })

  it('renders chart container when data is available', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
  })

  it('displays no data message when no referrer data available', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: { referrers: { referrers: [], directClicks: 0, totalClicks: 0, categories: [], topDomains: [] } },
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ReferrerStats />)
    
    // Search for something that won't match
    const searchInput = screen.getByPlaceholderText('Search referrers...')
    fireEvent.change(searchInput, { target: { value: 'nonexistent' } })
    
    expect(screen.getByText('No Referrer Data')).toBeInTheDocument()
    expect(screen.getByText('No traffic source data found for the selected view.')).toBeInTheDocument()
  })

  it('passes urlId prop to useRealTimeAnalytics hook', () => {
    const testUrlId = 'test-url-123'
    render(<ReferrerStats urlId={testUrlId} />)
    
    expect(mockUseRealTimeAnalytics).toHaveBeenCalledWith({
      urlId: testUrlId,
      period: '30d',
      refreshInterval: 30000,
      enabled: true
    })
  })

  it('applies custom className prop', () => {
    const customClass = 'custom-referrer-class'
    render(<ReferrerStats className={customClass} />)
    
    const referrerContainer = screen.getByText('Traffic Source Analytics').closest('.bg-white')
    expect(referrerContainer).toHaveClass(customClass)
  })

  it('shows refreshing state when data is being refreshed', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockReferrerData,
      isLoading: false,
      isRefreshing: true,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ReferrerStats />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    expect(refreshButton).toHaveAttribute('disabled')
  })

  it('handles different chart types correctly', () => {
    render(<ReferrerStats />)
    
    // Test pie chart (default)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to bar chart
    const barButton = screen.getByText('Bar Chart')
    fireEvent.click(barButton)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to treemap
    const treemapButton = screen.getByText('Treemap')
    fireEvent.click(treemapButton)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to table
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    expect(screen.getByRole('table')).toBeInTheDocument()
  })

  it('handles different view modes correctly', () => {
    render(<ReferrerStats />)
    
    // Test all sources (default)
    expect(screen.getByText('All Sources').closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to search engines
    const searchButton = screen.getByText('Search Engines')
    fireEvent.click(searchButton)
    expect(searchButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to social media
    const socialButton = screen.getByText('Social Media')
    fireEvent.click(socialButton)
    expect(socialButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to direct traffic
    const directButton = screen.getByText('Direct Traffic')
    fireEvent.click(directButton)
    expect(directButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to email
    const emailButton = screen.getByText('Email')
    fireEvent.click(emailButton)
    expect(emailButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to other
    const otherButton = screen.getByText('Other')
    fireEvent.click(otherButton)
    expect(otherButton.closest('button')).toHaveClass('bg-primary-100')
  })

  it('renders table with correct headers when chart type is table', () => {
    render(<ReferrerStats />)
    
    // Switch to table view
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('Referrer')).toBeInTheDocument()
    expect(screen.getByText('Category')).toBeInTheDocument()
    expect(screen.getByText('Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Clicks')).toBeInTheDocument()
    expect(screen.getByText('Percentage')).toBeInTheDocument()
    expect(screen.getByText('Trend')).toBeInTheDocument()
  })

  it('displays category breakdown section', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('Traffic by Category')).toBeInTheDocument()
  })

  it('filters search results correctly', () => {
    render(<ReferrerStats />)
    
    const searchInput = screen.getByPlaceholderText('Search referrers...')
    fireEvent.change(searchInput, { target: { value: 'google' } })
    
    expect(searchInput).toHaveValue('google')
  })

  it('handles empty search gracefully', () => {
    render(<ReferrerStats />)
    
    const searchInput = screen.getByPlaceholderText('Search referrers...')
    fireEvent.change(searchInput, { target: { value: '' } })
    
    expect(searchInput).toHaveValue('')
  })

  it('sorts data correctly', () => {
    render(<ReferrerStats />)
    
    const sortSelect = screen.getByDisplayValue('Most Clicks')
    
    // Test alphabetical sorting
    fireEvent.change(sortSelect, { target: { value: 'alphabetical' } })
    expect(sortSelect).toHaveValue('alphabetical')
    
    // Test trending sorting
    fireEvent.change(sortSelect, { target: { value: 'trend' } })
    expect(sortSelect).toHaveValue('trend')
    
    // Test clicks sorting
    fireEvent.change(sortSelect, { target: { value: 'clicks' } })
    expect(sortSelect).toHaveValue('clicks')
  })

  it('displays trend indicators correctly', () => {
    render(<ReferrerStats />)
    
    // Switch to table to see trend indicators clearly
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('Trend')).toBeInTheDocument()
  })

  it('categorizes referrers correctly', () => {
    render(<ReferrerStats />)
    
    // Switch to table to see categories
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('Category')).toBeInTheDocument()
  })

  it('handles direct traffic view mode', () => {
    render(<ReferrerStats />)
    
    const directButton = screen.getByText('Direct Traffic')
    fireEvent.click(directButton)
    
    expect(directButton.closest('button')).toHaveClass('bg-primary-100')
  })

  it('includes direct traffic in all sources view', () => {
    render(<ReferrerStats />)
    
    // All sources should include direct traffic
    expect(screen.getByText('All Sources').closest('button')).toHaveClass('bg-primary-100')
  })

  it('calculates summary statistics correctly', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Direct Traffic')).toBeInTheDocument()
    expect(screen.getByText('Unique Referrers')).toBeInTheDocument()
    expect(screen.getByText('Top Referrer')).toBeInTheDocument()
  })

  it('displays category breakdown with colors', () => {
    render(<ReferrerStats />)
    
    expect(screen.getByText('Traffic by Category')).toBeInTheDocument()
    
    // Category breakdown should show different traffic sources
    const categorySection = screen.getByText('Traffic by Category').closest('div')
    expect(categorySection).toBeInTheDocument()
  })

  it('shows category badges in table view', () => {
    render(<ReferrerStats />)
    
    // Switch to table view to see category badges
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('Category')).toBeInTheDocument()
  })
})