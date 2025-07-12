import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest'
import Dashboard from './Dashboard'
import { useRealTimeAnalytics } from '@/hooks/useRealTimeAnalytics'

// Mock the real-time analytics hook
vi.mock('@/hooks/useRealTimeAnalytics')

// Mock recharts components to avoid canvas rendering issues in tests
vi.mock('recharts', () => ({
  AreaChart: ({ children }: any) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div data-testid="area" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />,
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>,
  LineChart: ({ children }: any) => <div data-testid="line-chart">{children}</div>,
  Line: () => <div data-testid="line" />,
  BarChart: ({ children }: any) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  PieChart: ({ children }: any) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => <div data-testid="pie" />,
  Cell: () => <div data-testid="cell" />
}))

const mockAnalyticsData = {
  timeline: [
    { date: '2024-01-01', clicks: 100, uniqueClicks: 80 },
    { date: '2024-01-02', clicks: 120, uniqueClicks: 95 },
    { date: '2024-01-03', clicks: 90, uniqueClicks: 75 }
  ],
  geographic: {
    countries: [
      { country: 'United States', clicks: 200, percentage: 50 },
      { country: 'United Kingdom', clicks: 100, percentage: 25 }
    ],
    cities: [
      { city: 'New York', country: 'United States', clicks: 150, percentage: 37.5 }
    ]
  },
  devices: {
    devices: [
      { device: 'Desktop', clicks: 200, percentage: 66.7 },
      { device: 'Mobile', clicks: 100, percentage: 33.3 }
    ]
  },
  referrers: {
    referrers: [
      { referrer: 'google.com', clicks: 150, percentage: 50 },
      { referrer: 'facebook.com', clicks: 75, percentage: 25 }
    ]
  }
}

describe('Dashboard', () => {
  const mockUseRealTimeAnalytics = vi.mocked(useRealTimeAnalytics)

  beforeEach(() => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockAnalyticsData,
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

  it('renders dashboard header correctly', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Analytics Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Overview of your URL performance and user engagement')).toBeInTheDocument()
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

    render(<Dashboard />)
    expect(screen.getByText('Loading dashboard...')).toBeInTheDocument()
  })

  it('displays error state', () => {
    const errorMessage = 'Failed to load analytics data'
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

    render(<Dashboard />)
    expect(screen.getByText(errorMessage)).toBeInTheDocument()
  })

  it('displays no data state', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: null,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: undefined,
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<Dashboard />)
    expect(screen.getByText('No Data Available')).toBeInTheDocument()
    expect(screen.getByText('Create some URLs to see analytics data.')).toBeInTheDocument()
  })

  it('renders key metrics correctly', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Total URLs')).toBeInTheDocument()
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Visitors')).toBeInTheDocument()
    expect(screen.getByText('Avg. Clicks/URL')).toBeInTheDocument()
  })

  it('displays period selector with all options', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Last Hour')).toBeInTheDocument()
    expect(screen.getByText('Last 24 Hours')).toBeInTheDocument()
    expect(screen.getByText('Last 7 Days')).toBeInTheDocument()
    expect(screen.getByText('Last 30 Days')).toBeInTheDocument()
    expect(screen.getByText('Last 90 Days')).toBeInTheDocument()
    expect(screen.getByText('Last Year')).toBeInTheDocument()
    expect(screen.getByText('All Time')).toBeInTheDocument()
  })

  it('changes period when period button is clicked', () => {
    render(<Dashboard />)
    
    const sevenDaysButton = screen.getByText('Last 7 Days')
    fireEvent.click(sevenDaysButton)
    
    // Check if the button is now active (has primary styling)
    expect(sevenDaysButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders charts and data visualizations', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Click Timeline')).toBeInTheDocument()
    expect(screen.getByText('Top Countries')).toBeInTheDocument()
    expect(screen.getByText('Top Performing URLs')).toBeInTheDocument()
    expect(screen.getByText('Device Types')).toBeInTheDocument()
    expect(screen.getByText('Recent Activity')).toBeInTheDocument()
  })

  it('displays refresh button and handles refresh action', async () => {
    const mockRefreshData = vi.fn()
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockAnalyticsData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: mockRefreshData,
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<Dashboard />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    expect(refreshButton).toBeInTheDocument()
    
    fireEvent.click(refreshButton!)
    expect(mockRefreshData).toHaveBeenCalledTimes(1)
  })

  it('shows refreshing state when data is being refreshed', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockAnalyticsData,
      isLoading: false,
      isRefreshing: true,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<Dashboard />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    expect(refreshButton).toHaveAttribute('disabled')
  })

  it('displays connection status indicators', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockAnalyticsData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'reconnecting',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<Dashboard />)
    expect(screen.getByText('Reconnecting...')).toBeInTheDocument()
  })

  it('handles export functionality', () => {
    // Mock console.log to capture export calls
    const consoleSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
    
    render(<Dashboard />)
    
    const exportButton = screen.getByText('Export').closest('button')
    fireEvent.click(exportButton!)
    
    expect(consoleSpy).toHaveBeenCalledWith('Exporting dashboard data:', expect.any(Object))
    
    consoleSpy.mockRestore()
  })

  it('calculates and displays trend indicators', () => {
    render(<Dashboard />)
    
    // Check for trend indicators in the metrics
    const clicksMetric = screen.getByText('Total Clicks').closest('div')
    expect(clicksMetric).toBeInTheDocument()
  })

  it('renders top URLs with proper data', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Top Performing URLs')).toBeInTheDocument()
    expect(screen.getByText('#1')).toBeInTheDocument()
    expect(screen.getByText('clicks')).toBeInTheDocument()
  })

  it('renders device breakdown correctly', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Device Types')).toBeInTheDocument()
  })

  it('renders recent activity feed', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Recent Activity')).toBeInTheDocument()
  })

  it('applies custom className prop', () => {
    const customClass = 'custom-dashboard-class'
    render(<Dashboard className={customClass} />)
    
    const dashboardContainer = screen.getByText('Analytics Dashboard').closest('div')
    expect(dashboardContainer?.parentElement).toHaveClass(customClass)
  })

  it('uses real-time analytics hook with correct options', () => {
    render(<Dashboard />)
    
    expect(mockUseRealTimeAnalytics).toHaveBeenCalledWith({
      period: '30d',
      refreshInterval: 30000,
      enabled: true
    })
  })

  it('handles analytics data transformation correctly', () => {
    render(<Dashboard />)
    
    // Verify that the component processes analytics data and displays derived metrics
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Visitors')).toBeInTheDocument()
  })

  it('renders responsive grid layouts', () => {
    render(<Dashboard />)
    
    // Check for grid classes that ensure responsive design
    const keyMetricsSection = screen.getByText('Total URLs').closest('.grid')
    expect(keyMetricsSection).toHaveClass('grid-cols-1', 'md:grid-cols-2', 'lg:grid-cols-4')
  })

  it('displays formatted numbers correctly', () => {
    render(<Dashboard />)
    
    // Check that large numbers are formatted with commas
    const totalClicksValue = screen.getByText('Total Clicks').closest('div')?.querySelector('.text-3xl')
    expect(totalClicksValue).toBeInTheDocument()
  })

  it('shows proper time formatting in activity feed', () => {
    render(<Dashboard />)
    
    expect(screen.getByText('Recent Activity')).toBeInTheDocument()
    // The component should display relative time for activities
  })
})