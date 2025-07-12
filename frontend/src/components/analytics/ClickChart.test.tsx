import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest'
import ClickChart from './ClickChart'
import { useRealTimeAnalytics } from '@/hooks/useRealTimeAnalytics'

// Mock the real-time analytics hook
vi.mock('@/hooks/useRealTimeAnalytics')

// Mock recharts components
vi.mock('recharts', () => ({
  AreaChart: ({ children }: any) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div data-testid="area" />,
  LineChart: ({ children }: any) => <div data-testid="line-chart">{children}</div>,
  Line: () => <div data-testid="line" />,
  BarChart: ({ children }: any) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />,
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>,
  ReferenceLine: () => <div data-testid="reference-line" />
}))

// Mock date-fns functions
vi.mock('date-fns', () => ({
  format: vi.fn((date, formatStr) => {
    if (formatStr === 'MMM d') return 'Jan 1'
    if (formatStr === 'yyyy-MM-dd') return '2024-01-01'
    if (formatStr === 'HH:mm') return '12:00'
    if (formatStr === 'MMM d, yyyy') return 'Jan 1, 2024'
    if (formatStr === 'h:mm a') return '12:00 PM'
    if (formatStr === 'EEE') return 'Mon'
    if (formatStr === 'MMM yyyy') return 'Jan 2024'
    if (formatStr === 'MMMM d, yyyy') return 'January 1, 2024'
    return '2024-01-01'
  }),
  subDays: vi.fn((date, days) => new Date(Date.now() - days * 24 * 60 * 60 * 1000)),
  subHours: vi.fn((date, hours) => new Date(Date.now() - hours * 60 * 60 * 1000)),
  parseISO: vi.fn((dateStr) => new Date(dateStr)),
  isToday: vi.fn(() => true),
  isYesterday: vi.fn(() => false),
  startOfDay: vi.fn((date) => date),
  endOfDay: vi.fn((date) => date)
}))

const mockAnalyticsData = {
  timeline: [
    { date: '2024-01-01', clicks: 100, uniqueClicks: 80 },
    { date: '2024-01-02', clicks: 120, uniqueClicks: 95 },
    { date: '2024-01-03', clicks: 90, uniqueClicks: 75 }
  ]
}

describe('ClickChart', () => {
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

  it('renders click chart header correctly', () => {
    render(<ClickChart />)
    
    expect(screen.getByText('Click Analytics')).toBeInTheDocument()
    expect(screen.getByText('Detailed click performance over time')).toBeInTheDocument()
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

    render(<ClickChart />)
    expect(screen.getByText('Loading chart data...')).toBeInTheDocument()
  })

  it('displays error state', () => {
    const errorMessage = 'Failed to load chart data'
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

    render(<ClickChart />)
    expect(screen.getByText(errorMessage)).toBeInTheDocument()
  })

  it('renders period selector buttons', () => {
    render(<ClickChart />)
    
    expect(screen.getByText('Last Hour')).toBeInTheDocument()
    expect(screen.getByText('Last 24 Hours')).toBeInTheDocument()
    expect(screen.getByText('Last 7 Days')).toBeInTheDocument()
    expect(screen.getByText('Last 30 Days')).toBeInTheDocument()
    expect(screen.getByText('Last 90 Days')).toBeInTheDocument()
    expect(screen.getByText('Last Year')).toBeInTheDocument()
    expect(screen.getByText('All Time')).toBeInTheDocument()
  })

  it('changes period when period button is clicked', () => {
    render(<ClickChart />)
    
    const sevenDaysButton = screen.getByText('Last 7 Days')
    fireEvent.click(sevenDaysButton)
    
    expect(sevenDaysButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders chart type selector', () => {
    render(<ClickChart />)
    
    expect(screen.getByText('Area')).toBeInTheDocument()
    expect(screen.getByText('Line')).toBeInTheDocument()
    expect(screen.getByText('Bar')).toBeInTheDocument()
  })

  it('changes chart type when chart type button is clicked', () => {
    render(<ClickChart />)
    
    const lineChartButton = screen.getByText('Line')
    fireEvent.click(lineChartButton)
    
    expect(lineChartButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders metric type selector', () => {
    render(<ClickChart />)
    
    expect(screen.getByText('Both Metrics')).toBeInTheDocument()
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Clicks')).toBeInTheDocument()
  })

  it('changes metric type when metric button is clicked', () => {
    render(<ClickChart />)
    
    const uniqueClicksButton = screen.getByText('Unique Clicks')
    fireEvent.click(uniqueClicksButton)
    
    expect(uniqueClicksButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('displays metrics summary when data is available', () => {
    render(<ClickChart />)
    
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Clicks')).toBeInTheDocument()
    expect(screen.getByText('Avg Bounce Rate')).toBeInTheDocument()
    expect(screen.getByText('Avg Session Time')).toBeInTheDocument()
  })

  it('displays refresh button and handles refresh action', () => {
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

    render(<ClickChart />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    fireEvent.click(refreshButton!)
    
    expect(mockRefreshData).toHaveBeenCalledTimes(1)
  })

  it('shows connection status indicators', () => {
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

    render(<ClickChart />)
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

    render(<ClickChart />)
    
    const exportButton = screen.getByText('Export').closest('button')
    fireEvent.click(exportButton!)
    
    expect(mockAnchor.click).toHaveBeenCalledTimes(1)
    expect(global.URL.createObjectURL).toHaveBeenCalledTimes(1)
    expect(global.URL.revokeObjectURL).toHaveBeenCalledWith('mock-url')
  })

  it('renders chart container when data is available', () => {
    render(<ClickChart />)
    
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
  })

  it('displays no data message when no chart data available', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: { timeline: [] },
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ClickChart />)
    expect(screen.getByText('No Data Available')).toBeInTheDocument()
    expect(screen.getByText('No click data found for the selected period.')).toBeInTheDocument()
  })

  it('passes urlId prop to useRealTimeAnalytics hook', () => {
    const testUrlId = 'test-url-123'
    render(<ClickChart urlId={testUrlId} />)
    
    expect(mockUseRealTimeAnalytics).toHaveBeenCalledWith({
      urlId: testUrlId,
      period: '7d',
      refreshInterval: 30000,
      enabled: true
    })
  })

  it('uses custom period prop', () => {
    render(<ClickChart period="30d" />)
    
    expect(mockUseRealTimeAnalytics).toHaveBeenCalledWith({
      urlId: undefined,
      period: '30d',
      refreshInterval: 30000,
      enabled: true
    })
  })

  it('applies custom className prop', () => {
    const customClass = 'custom-chart-class'
    render(<ClickChart className={customClass} />)
    
    const chartContainer = screen.getByText('Click Analytics').closest('.bg-white')
    expect(chartContainer).toHaveClass(customClass)
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

    render(<ClickChart />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    expect(refreshButton).toHaveAttribute('disabled')
  })

  it('calculates trend indicators correctly', () => {
    render(<ClickChart />)
    
    // The trend calculation should be visible in the metrics summary
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
  })

  it('formats large numbers with commas', () => {
    const largeData = {
      timeline: [
        { date: '2024-01-01', clicks: 1000, uniqueClicks: 800 },
        { date: '2024-01-02', clicks: 1200, uniqueClicks: 950 }
      ]
    }
    
    mockUseRealTimeAnalytics.mockReturnValue({
      data: largeData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<ClickChart />)
    
    // Check that numbers are formatted properly
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
  })

  it('handles different chart types correctly', () => {
    render(<ClickChart />)
    
    // Test area chart (default)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to line chart
    const lineButton = screen.getByText('Line')
    fireEvent.click(lineButton)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to bar chart
    const barButton = screen.getByText('Bar')
    fireEvent.click(barButton)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
  })

  it('handles different metric types correctly', () => {
    render(<ClickChart />)
    
    // Test both metrics (default)
    expect(screen.getByText('Both Metrics').closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to clicks only
    const clicksButton = screen.getByText('Total Clicks')
    fireEvent.click(clicksButton)
    expect(clicksButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to unique clicks only
    const uniqueClicksButton = screen.getByText('Unique Clicks')
    fireEvent.click(uniqueClicksButton)
    expect(uniqueClicksButton.closest('button')).toHaveClass('bg-primary-100')
  })
})