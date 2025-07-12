import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest'
import DeviceStats from './DeviceStats'
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
  RadialBarChart: ({ children }: any) => <div data-testid="radial-bar-chart">{children}</div>,
  RadialBar: () => <div data-testid="radial-bar" />
}))

const mockDeviceData = {
  devices: {
    devices: [
      { device: 'Desktop', clicks: 650, percentage: 65.0, uniqueClicks: 520, trend: 5.2 },
      { device: 'Mobile', clicks: 300, percentage: 30.0, uniqueClicks: 240, trend: 12.8 },
      { device: 'Tablet', clicks: 50, percentage: 5.0, uniqueClicks: 40, trend: -2.1 }
    ],
    browsers: [
      { browser: 'Chrome', clicks: 500, percentage: 50.0, uniqueClicks: 400, version: '119.0', trend: 2.4 },
      { browser: 'Firefox', clicks: 200, percentage: 20.0, uniqueClicks: 160, version: '118.0', trend: -1.2 },
      { browser: 'Safari', clicks: 150, percentage: 15.0, uniqueClicks: 120, version: '17.1', trend: 8.7 }
    ],
    operatingSystems: [
      { os: 'Windows', clicks: 400, percentage: 40.0, uniqueClicks: 320, version: '11', trend: 1.8 },
      { os: 'macOS', clicks: 250, percentage: 25.0, uniqueClicks: 200, version: '14.1', trend: 4.2 },
      { os: 'iOS', clicks: 200, percentage: 20.0, uniqueClicks: 160, version: '17.1', trend: 6.9 }
    ]
  }
}

describe('DeviceStats', () => {
  const mockUseRealTimeAnalytics = vi.mocked(useRealTimeAnalytics)

  beforeEach(() => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockDeviceData,
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

  it('renders device stats header correctly', () => {
    render(<DeviceStats />)
    
    expect(screen.getByText('Device & Browser Analytics')).toBeInTheDocument()
    expect(screen.getByText('User device, browser, and system information')).toBeInTheDocument()
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

    render(<DeviceStats />)
    expect(screen.getByText('Loading device analytics...')).toBeInTheDocument()
  })

  it('displays error state', () => {
    const errorMessage = 'Failed to load device data'
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

    render(<DeviceStats />)
    expect(screen.getByText(errorMessage)).toBeInTheDocument()
  })

  it('renders category selector buttons', () => {
    render(<DeviceStats />)
    
    expect(screen.getByText('Devices')).toBeInTheDocument()
    expect(screen.getByText('Browsers')).toBeInTheDocument()
    expect(screen.getByText('Operating Systems')).toBeInTheDocument()
    expect(screen.getByText('Screen Resolutions')).toBeInTheDocument()
  })

  it('changes view mode when category button is clicked', () => {
    render(<DeviceStats />)
    
    const browsersButton = screen.getByText('Browsers')
    fireEvent.click(browsersButton)
    
    expect(browsersButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders chart type selector', () => {
    render(<DeviceStats />)
    
    expect(screen.getByText('Pie Chart')).toBeInTheDocument()
    expect(screen.getByText('Bar Chart')).toBeInTheDocument()
    expect(screen.getByText('Radial Chart')).toBeInTheDocument()
  })

  it('changes chart type when chart type button is clicked', () => {
    render(<DeviceStats />)
    
    const barChartButton = screen.getByText('Bar Chart')
    fireEvent.click(barChartButton)
    
    expect(barChartButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('displays summary statistics', () => {
    render(<DeviceStats />)
    
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Combinations')).toBeInTheDocument()
    expect(screen.getByText('Most Popular Device')).toBeInTheDocument()
    expect(screen.getByText('Most Popular Browser')).toBeInTheDocument()
  })

  it('displays refresh button and handles refresh action', () => {
    const mockRefreshData = vi.fn()
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockDeviceData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: mockRefreshData,
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<DeviceStats />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    fireEvent.click(refreshButton!)
    
    expect(mockRefreshData).toHaveBeenCalledTimes(1)
  })

  it('shows connection status indicators', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockDeviceData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'reconnecting',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<DeviceStats />)
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

    render(<DeviceStats />)
    
    const exportButton = screen.getByText('Export').closest('button')
    fireEvent.click(exportButton!)
    
    expect(mockAnchor.click).toHaveBeenCalledTimes(1)
    expect(global.URL.createObjectURL).toHaveBeenCalledTimes(1)
    expect(global.URL.revokeObjectURL).toHaveBeenCalledWith('mock-url')
  })

  it('renders chart container when data is available', () => {
    render(<DeviceStats />)
    
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
  })

  it('displays no data message when no device data available', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: { devices: { devices: [], browsers: [], operatingSystems: [] } },
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<DeviceStats />)
    expect(screen.getByText('No Device Data')).toBeInTheDocument()
    expect(screen.getByText('No device information available for the selected view.')).toBeInTheDocument()
  })

  it('passes urlId prop to useRealTimeAnalytics hook', () => {
    const testUrlId = 'test-url-123'
    render(<DeviceStats urlId={testUrlId} />)
    
    expect(mockUseRealTimeAnalytics).toHaveBeenCalledWith({
      urlId: testUrlId,
      period: '30d',
      refreshInterval: 30000,
      enabled: true
    })
  })

  it('applies custom className prop', () => {
    const customClass = 'custom-device-class'
    render(<DeviceStats className={customClass} />)
    
    const deviceContainer = screen.getByText('Device & Browser Analytics').closest('.bg-white')
    expect(deviceContainer).toHaveClass(customClass)
  })

  it('shows refreshing state when data is being refreshed', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockDeviceData,
      isLoading: false,
      isRefreshing: true,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<DeviceStats />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    expect(refreshButton).toHaveAttribute('disabled')
  })

  it('handles different chart types correctly', () => {
    render(<DeviceStats />)
    
    // Test pie chart (default)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to bar chart
    const barButton = screen.getByText('Bar Chart')
    fireEvent.click(barButton)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to radial chart
    const radialButton = screen.getByText('Radial Chart')
    fireEvent.click(radialButton)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
  })

  it('handles different view modes correctly', () => {
    render(<DeviceStats />)
    
    // Test devices (default)
    expect(screen.getByText('Devices').closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to browsers
    const browsersButton = screen.getByText('Browsers')
    fireEvent.click(browsersButton)
    expect(browsersButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to operating systems
    const osButton = screen.getByText('Operating Systems')
    fireEvent.click(osButton)
    expect(osButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to screen resolutions
    const resolutionsButton = screen.getByText('Screen Resolutions')
    fireEvent.click(resolutionsButton)
    expect(resolutionsButton.closest('button')).toHaveClass('bg-primary-100')
  })

  it('renders detailed breakdown with device icons', () => {
    render(<DeviceStats />)
    
    expect(screen.getByText('Detailed Breakdown')).toBeInTheDocument()
    expect(screen.getByText('Desktop')).toBeInTheDocument()
    expect(screen.getByText('Mobile')).toBeInTheDocument()
    expect(screen.getByText('Tablet')).toBeInTheDocument()
  })

  it('displays trend indicators correctly', () => {
    render(<DeviceStats />)
    
    // Look for trend information in the detailed breakdown
    expect(screen.getByText('Detailed Breakdown')).toBeInTheDocument()
    
    // The trends should be visible in the device breakdown cards
    const detailsSection = screen.getByText('Detailed Breakdown').closest('div')
    expect(detailsSection).toBeInTheDocument()
  })

  it('shows version information for browsers and OS', () => {
    render(<DeviceStats />)
    
    // Switch to browsers to see version info
    const browsersButton = screen.getByText('Browsers')
    fireEvent.click(browsersButton)
    
    expect(screen.getByText('Chrome')).toBeInTheDocument()
    expect(screen.getByText('Firefox')).toBeInTheDocument()
    expect(screen.getByText('Safari')).toBeInTheDocument()
  })

  it('displays correct summary statistics for each view mode', () => {
    render(<DeviceStats />)
    
    // Check devices view stats
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Most Popular Device')).toBeInTheDocument()
    
    // Switch to browsers and check stats update
    const browsersButton = screen.getByText('Browsers')
    fireEvent.click(browsersButton)
    
    expect(screen.getByText('Most Popular Browser')).toBeInTheDocument()
  })

  it('handles empty data gracefully for different categories', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: { devices: { devices: [], browsers: [], operatingSystems: [] } },
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<DeviceStats />)
    
    // Should show no data message for all categories
    expect(screen.getByText('No Device Data')).toBeInTheDocument()
    
    // Switch to browsers - should still show no data
    const browsersButton = screen.getByText('Browsers')
    fireEvent.click(browsersButton)
    expect(screen.getByText('No Device Data')).toBeInTheDocument()
    
    // Switch to OS - should still show no data
    const osButton = screen.getByText('Operating Systems')
    fireEvent.click(osButton)
    expect(screen.getByText('No Device Data')).toBeInTheDocument()
  })

  it('calculates correct summary statistics', () => {
    render(<DeviceStats />)
    
    // Check that total clicks and other stats are calculated correctly
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Combinations')).toBeInTheDocument()
    expect(screen.getByText('Most Popular Device')).toBeInTheDocument()
    expect(screen.getByText('Most Popular Browser')).toBeInTheDocument()
  })

  it('displays formatted click counts', () => {
    render(<DeviceStats />)
    
    // Check that click counts are displayed in the detailed breakdown
    expect(screen.getByText('Detailed Breakdown')).toBeInTheDocument()
    
    // Device click counts should be visible
    const detailsSection = screen.getByText('Detailed Breakdown').closest('div')
    expect(detailsSection).toBeInTheDocument()
  })

  it('shows unique clicks information where available', () => {
    render(<DeviceStats />)
    
    // Unique clicks should be shown for devices, browsers, and OS but not resolutions
    expect(screen.getByText('Detailed Breakdown')).toBeInTheDocument()
    
    // Switch to resolutions to verify unique clicks are not shown
    const resolutionsButton = screen.getByText('Screen Resolutions')
    fireEvent.click(resolutionsButton)
    
    // Resolution view should not show unique clicks
    expect(screen.getByText('Detailed Breakdown')).toBeInTheDocument()
  })
})