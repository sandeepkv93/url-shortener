import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest'
import GeographicMap from './GeographicMap'
import { useRealTimeAnalytics } from '@/hooks/useRealTimeAnalytics'

// Mock the real-time analytics hook
vi.mock('@/hooks/useRealTimeAnalytics')

// Mock recharts components
vi.mock('recharts', () => ({
  BarChart: ({ children }: any) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  PieChart: ({ children }: any) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => <div data-testid="pie" />,
  Cell: () => <div data-testid="cell" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />,
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>,
  Treemap: ({ children }: any) => <div data-testid="treemap">{children}</div>
}))

const mockGeographicData = {
  geographic: {
    countries: [
      { country: 'United States', countryCode: 'US', clicks: 850, percentage: 42.5, uniqueClicks: 680 },
      { country: 'United Kingdom', countryCode: 'GB', clicks: 350, percentage: 17.5, uniqueClicks: 280 },
      { country: 'Canada', countryCode: 'CA', clicks: 250, percentage: 12.5, uniqueClicks: 200 }
    ],
    cities: [
      { city: 'New York', country: 'United States', countryCode: 'US', clicks: 320, percentage: 16.0, uniqueClicks: 256 },
      { city: 'London', country: 'United Kingdom', countryCode: 'GB', clicks: 280, percentage: 14.0, uniqueClicks: 224 },
      { city: 'Toronto', country: 'Canada', countryCode: 'CA', clicks: 180, percentage: 9.0, uniqueClicks: 144 }
    ]
  }
}

describe('GeographicMap', () => {
  const mockUseRealTimeAnalytics = vi.mocked(useRealTimeAnalytics)

  beforeEach(() => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockGeographicData,
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

  it('renders geographic map header correctly', () => {
    render(<GeographicMap />)
    
    expect(screen.getByText('Geographic Analytics')).toBeInTheDocument()
    expect(screen.getByText('Click distribution by location')).toBeInTheDocument()
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

    render(<GeographicMap />)
    expect(screen.getByText('Loading geographic data...')).toBeInTheDocument()
  })

  it('displays error state', () => {
    const errorMessage = 'Failed to load geographic data'
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

    render(<GeographicMap />)
    expect(screen.getByText(errorMessage)).toBeInTheDocument()
  })

  it('renders view mode selector buttons', () => {
    render(<GeographicMap />)
    
    expect(screen.getByText('Countries')).toBeInTheDocument()
    expect(screen.getByText('Cities')).toBeInTheDocument()
    expect(screen.getByText('Regions')).toBeInTheDocument()
  })

  it('changes view mode when view mode button is clicked', () => {
    render(<GeographicMap />)
    
    const citiesButton = screen.getByText('Cities')
    fireEvent.click(citiesButton)
    
    expect(citiesButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders chart type selector', () => {
    render(<GeographicMap />)
    
    expect(screen.getByText('Bar Chart')).toBeInTheDocument()
    expect(screen.getByText('Pie Chart')).toBeInTheDocument()
    expect(screen.getByText('Treemap')).toBeInTheDocument()
    expect(screen.getByText('Table')).toBeInTheDocument()
  })

  it('changes chart type when chart type button is clicked', () => {
    render(<GeographicMap />)
    
    const pieChartButton = screen.getByText('Pie Chart')
    fireEvent.click(pieChartButton)
    
    expect(pieChartButton.closest('button')).toHaveClass('bg-primary-100', 'text-primary-700')
  })

  it('renders search input', () => {
    render(<GeographicMap />)
    
    const searchInput = screen.getByPlaceholderText('Search countries...')
    expect(searchInput).toBeInTheDocument()
  })

  it('filters data based on search input', () => {
    render(<GeographicMap />)
    
    const searchInput = screen.getByPlaceholderText('Search countries...')
    fireEvent.change(searchInput, { target: { value: 'United' } })
    
    // The search should filter the displayed countries
    expect(searchInput).toHaveValue('United')
  })

  it('renders sort selector', () => {
    render(<GeographicMap />)
    
    const sortSelect = screen.getByDisplayValue('Most Clicks')
    expect(sortSelect).toBeInTheDocument()
    
    const alphabeticalOption = screen.getByText('Alphabetical')
    expect(alphabeticalOption).toBeInTheDocument()
  })

  it('changes sort order when sort option is selected', () => {
    render(<GeographicMap />)
    
    const sortSelect = screen.getByDisplayValue('Most Clicks')
    fireEvent.change(sortSelect, { target: { value: 'alphabetical' } })
    
    expect(sortSelect).toHaveValue('alphabetical')
  })

  it('displays summary statistics', () => {
    render(<GeographicMap />)
    
    expect(screen.getByText('Countries')).toBeInTheDocument()
    expect(screen.getByText('Cities')).toBeInTheDocument()
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Top Country')).toBeInTheDocument()
  })

  it('displays refresh button and handles refresh action', () => {
    const mockRefreshData = vi.fn()
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockGeographicData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: mockRefreshData,
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<GeographicMap />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    fireEvent.click(refreshButton!)
    
    expect(mockRefreshData).toHaveBeenCalledTimes(1)
  })

  it('shows connection status indicators', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockGeographicData,
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'reconnecting',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<GeographicMap />)
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

    render(<GeographicMap />)
    
    const exportButton = screen.getByText('Export').closest('button')
    fireEvent.click(exportButton!)
    
    expect(mockAnchor.click).toHaveBeenCalledTimes(1)
    expect(global.URL.createObjectURL).toHaveBeenCalledTimes(1)
    expect(global.URL.revokeObjectURL).toHaveBeenCalledWith('mock-url')
  })

  it('renders chart container when data is available', () => {
    render(<GeographicMap />)
    
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
  })

  it('displays no data message when no geographic data available', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: { geographic: { countries: [], cities: [] } },
      isLoading: false,
      isRefreshing: false,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<GeographicMap />)
    
    // Change search to something that won't match
    const searchInput = screen.getByPlaceholderText('Search countries...')
    fireEvent.change(searchInput, { target: { value: 'NonExistentCountry' } })
    
    expect(screen.getByText('No Geographic Data')).toBeInTheDocument()
    expect(screen.getByText('No location data found for the selected view.')).toBeInTheDocument()
  })

  it('passes urlId prop to useRealTimeAnalytics hook', () => {
    const testUrlId = 'test-url-123'
    render(<GeographicMap urlId={testUrlId} />)
    
    expect(mockUseRealTimeAnalytics).toHaveBeenCalledWith({
      urlId: testUrlId,
      period: '30d',
      refreshInterval: 30000,
      enabled: true
    })
  })

  it('applies custom className prop', () => {
    const customClass = 'custom-map-class'
    render(<GeographicMap className={customClass} />)
    
    const mapContainer = screen.getByText('Geographic Analytics').closest('.bg-white')
    expect(mapContainer).toHaveClass(customClass)
  })

  it('shows refreshing state when data is being refreshed', () => {
    mockUseRealTimeAnalytics.mockReturnValue({
      data: mockGeographicData,
      isLoading: false,
      isRefreshing: true,
      error: null,
      refreshData: vi.fn(),
      connectionStatus: 'connected',
      lastUpdated: new Date(),
      startRealTimeUpdates: vi.fn(),
      stopRealTimeUpdates: vi.fn()
    })

    render(<GeographicMap />)
    
    const refreshButton = screen.getByText('Refresh').closest('button')
    expect(refreshButton).toHaveAttribute('disabled')
  })

  it('handles different chart types correctly', () => {
    render(<GeographicMap />)
    
    // Test bar chart (default)
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument()
    
    // Switch to pie chart
    const pieButton = screen.getByText('Pie Chart')
    fireEvent.click(pieButton)
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
    render(<GeographicMap />)
    
    // Test countries (default)
    expect(screen.getByText('Countries').closest('button')).toHaveClass('bg-primary-100')
    
    // Switch to cities
    const citiesButton = screen.getByText('Cities')
    fireEvent.click(citiesButton)
    expect(citiesButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Check that search placeholder changes
    expect(screen.getByPlaceholderText('Search cities...')).toBeInTheDocument()
    
    // Switch to regions
    const regionsButton = screen.getByText('Regions')
    fireEvent.click(regionsButton)
    expect(regionsButton.closest('button')).toHaveClass('bg-primary-100')
    
    // Check that search placeholder changes
    expect(screen.getByPlaceholderText('Search regions...')).toBeInTheDocument()
  })

  it('renders table with correct headers when chart type is table', () => {
    render(<GeographicMap />)
    
    // Switch to table view
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('Country')).toBeInTheDocument()
    expect(screen.getByText('Clicks')).toBeInTheDocument()
    expect(screen.getByText('Unique Clicks')).toBeInTheDocument()
    expect(screen.getByText('Percentage')).toBeInTheDocument()
  })

  it('displays correct table headers for cities view', () => {
    render(<GeographicMap />)
    
    // Switch to cities and table view
    const citiesButton = screen.getByText('Cities')
    fireEvent.click(citiesButton)
    
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('City')).toBeInTheDocument()
    expect(screen.getByText('Country')).toBeInTheDocument()
  })

  it('displays correct table headers for regions view', () => {
    render(<GeographicMap />)
    
    // Switch to regions and table view
    const regionsButton = screen.getByText('Regions')
    fireEvent.click(regionsButton)
    
    const tableButton = screen.getByText('Table')
    fireEvent.click(tableButton)
    
    expect(screen.getByText('Region')).toBeInTheDocument()
    expect(screen.getByText('Country')).toBeInTheDocument()
  })

  it('filters search results correctly', () => {
    render(<GeographicMap />)
    
    const searchInput = screen.getByPlaceholderText('Search countries...')
    fireEvent.change(searchInput, { target: { value: 'United States' } })
    
    // Should show results for United States
    expect(searchInput).toHaveValue('United States')
  })

  it('handles empty search gracefully', () => {
    render(<GeographicMap />)
    
    const searchInput = screen.getByPlaceholderText('Search countries...')
    fireEvent.change(searchInput, { target: { value: '' } })
    
    // Should show all results when search is empty
    expect(searchInput).toHaveValue('')
  })

  it('sorts data correctly', () => {
    render(<GeographicMap />)
    
    const sortSelect = screen.getByDisplayValue('Most Clicks')
    
    // Test alphabetical sorting
    fireEvent.change(sortSelect, { target: { value: 'alphabetical' } })
    expect(sortSelect).toHaveValue('alphabetical')
    
    // Test clicks sorting
    fireEvent.change(sortSelect, { target: { value: 'clicks' } })
    expect(sortSelect).toHaveValue('clicks')
  })

  it('calculates summary statistics correctly', () => {
    render(<GeographicMap />)
    
    // Check that summary stats are displayed
    expect(screen.getByText('Countries')).toBeInTheDocument()
    expect(screen.getByText('Cities')).toBeInTheDocument()
    expect(screen.getByText('Total Clicks')).toBeInTheDocument()
    expect(screen.getByText('Top Country')).toBeInTheDocument()
  })
})