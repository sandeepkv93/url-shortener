import { useState, useMemo } from 'react'
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Treemap
} from 'recharts'
import {
  Globe,
  MapPin,
  Users,
  TrendingUp,
  Filter,
  Download,
  RefreshCw,
  Search,
  Map,
  BarChart3,
  PieChart as PieChartIcon,
  Grid
} from 'lucide-react'
import { PageLoading } from '@/components/common/Loading'
import { useRealTimeAnalytics } from '@/hooks/useRealTimeAnalytics'

interface GeographicData {
  countries: Array<{
    country: string
    countryCode: string
    clicks: number
    percentage: number
    uniqueClicks: number
    coordinates?: { lat: number; lng: number }
  }>
  cities: Array<{
    city: string
    country: string
    countryCode: string
    clicks: number
    percentage: number
    uniqueClicks: number
    coordinates?: { lat: number; lng: number }
  }>
  regions: Array<{
    region: string
    country: string
    clicks: number
    percentage: number
  }>
}

interface GeographicMapProps {
  urlId?: string // If provided, shows analytics for specific URL
  className?: string
}

type ViewMode = 'countries' | 'cities' | 'regions'
type ChartType = 'bar' | 'pie' | 'treemap' | 'table'

const GeographicMap = ({ urlId, className = '' }: GeographicMapProps) => {
  const [viewMode, setViewMode] = useState<ViewMode>('countries')
  const [chartType, setChartType] = useState<ChartType>('bar')
  const [searchTerm, setSearchTerm] = useState('')
  const [sortBy, setSortBy] = useState<'clicks' | 'alphabetical'>('clicks')
  
  // Use real-time analytics hook
  const {
    data: analyticsData,
    isLoading,
    isRefreshing,
    error,
    refreshData,
    connectionStatus
  } = useRealTimeAnalytics({
    urlId,
    period: '30d',
    refreshInterval: 30000, // 30 seconds
    enabled: true
  })
  
  // Transform analytics data to geographic format
  const data = useMemo(() => {
    if (!analyticsData?.geographic) return generateMockGeographicData()
    
    return {
      countries: analyticsData.geographic.countries || [],
      cities: analyticsData.geographic.cities || [],
      regions: generateMockRegions()
    }
  }, [analyticsData])

  const viewModes: { value: ViewMode; label: string; icon: React.ReactNode }[] = [
    { value: 'countries', label: 'Countries', icon: <Globe className="h-4 w-4" /> },
    { value: 'cities', label: 'Cities', icon: <MapPin className="h-4 w-4" /> },
    { value: 'regions', label: 'Regions', icon: <Map className="h-4 w-4" /> }
  ]

  const chartTypes: { value: ChartType; label: string; icon: React.ReactNode }[] = [
    { value: 'bar', label: 'Bar Chart', icon: <BarChart3 className="h-4 w-4" /> },
    { value: 'pie', label: 'Pie Chart', icon: <PieChartIcon className="h-4 w-4" /> },
    { value: 'treemap', label: 'Treemap', icon: <Grid className="h-4 w-4" /> },
    { value: 'table', label: 'Table', icon: <Filter className="h-4 w-4" /> }
  ]


  const generateMockGeographicData = (): GeographicData => ({
    countries: [
      { country: 'United States', countryCode: 'US', clicks: 850, percentage: 42.5, uniqueClicks: 680, coordinates: { lat: 39.8283, lng: -98.5795 } },
      { country: 'United Kingdom', countryCode: 'GB', clicks: 350, percentage: 17.5, uniqueClicks: 280, coordinates: { lat: 55.3781, lng: -3.4360 } },
      { country: 'Canada', countryCode: 'CA', clicks: 250, percentage: 12.5, uniqueClicks: 200, coordinates: { lat: 56.1304, lng: -106.3468 } },
      { country: 'Germany', countryCode: 'DE', clicks: 180, percentage: 9.0, uniqueClicks: 144, coordinates: { lat: 51.1657, lng: 10.4515 } },
      { country: 'France', countryCode: 'FR', clicks: 120, percentage: 6.0, uniqueClicks: 96, coordinates: { lat: 46.2276, lng: 2.2137 } },
      { country: 'Australia', countryCode: 'AU', clicks: 100, percentage: 5.0, uniqueClicks: 80, coordinates: { lat: -25.2744, lng: 133.7751 } },
      { country: 'Japan', countryCode: 'JP', clicks: 80, percentage: 4.0, uniqueClicks: 64, coordinates: { lat: 36.2048, lng: 138.2529 } },
      { country: 'Brazil', countryCode: 'BR', clicks: 70, percentage: 3.5, uniqueClicks: 56, coordinates: { lat: -14.2350, lng: -51.9253 } }
    ],
    cities: [
      { city: 'New York', country: 'United States', countryCode: 'US', clicks: 320, percentage: 16.0, uniqueClicks: 256, coordinates: { lat: 40.7128, lng: -74.0060 } },
      { city: 'London', country: 'United Kingdom', countryCode: 'GB', clicks: 280, percentage: 14.0, uniqueClicks: 224, coordinates: { lat: 51.5074, lng: -0.1278 } },
      { city: 'Los Angeles', country: 'United States', countryCode: 'US', clicks: 220, percentage: 11.0, uniqueClicks: 176, coordinates: { lat: 34.0522, lng: -118.2437 } },
      { city: 'Toronto', country: 'Canada', countryCode: 'CA', clicks: 180, percentage: 9.0, uniqueClicks: 144, coordinates: { lat: 43.6532, lng: -79.3832 } },
      { city: 'Paris', country: 'France', countryCode: 'FR', clicks: 150, percentage: 7.5, uniqueClicks: 120, coordinates: { lat: 48.8566, lng: 2.3522 } },
      { city: 'Berlin', country: 'Germany', countryCode: 'DE', clicks: 130, percentage: 6.5, uniqueClicks: 104, coordinates: { lat: 52.5200, lng: 13.4050 } },
      { city: 'Chicago', country: 'United States', countryCode: 'US', clicks: 120, percentage: 6.0, uniqueClicks: 96, coordinates: { lat: 41.8781, lng: -87.6298 } },
      { city: 'Sydney', country: 'Australia', countryCode: 'AU', clicks: 100, percentage: 5.0, uniqueClicks: 80, coordinates: { lat: -33.8688, lng: 151.2093 } }
    ],
    regions: generateMockRegions()
  })

  const generateMockRegions = () => [
    { region: 'California', country: 'United States', clicks: 450, percentage: 22.5 },
    { region: 'England', country: 'United Kingdom', clicks: 300, percentage: 15.0 },
    { region: 'Ontario', country: 'Canada', clicks: 200, percentage: 10.0 },
    { region: 'Bavaria', country: 'Germany', clicks: 120, percentage: 6.0 },
    { region: 'Île-de-France', country: 'France', clicks: 100, percentage: 5.0 },
    { region: 'New South Wales', country: 'Australia', clicks: 80, percentage: 4.0 }
  ]

  const getCountryCoordinates = (countryCode: string) => {
    const coordinates: { [key: string]: { lat: number; lng: number } } = {
      'US': { lat: 39.8283, lng: -98.5795 },
      'GB': { lat: 55.3781, lng: -3.4360 },
      'CA': { lat: 56.1304, lng: -106.3468 },
      'DE': { lat: 51.1657, lng: 10.4515 },
      'FR': { lat: 46.2276, lng: 2.2137 },
      'AU': { lat: -25.2744, lng: 133.7751 },
      'JP': { lat: 36.2048, lng: 138.2529 },
      'BR': { lat: -14.2350, lng: -51.9253 }
    }
    return coordinates[countryCode]
  }

  const getCountryCode = (country: string) => {
    const codes: { [key: string]: string } = {
      'United States': 'US',
      'United Kingdom': 'GB',
      'Canada': 'CA',
      'Germany': 'DE',
      'France': 'FR',
      'Australia': 'AU',
      'Japan': 'JP',
      'Brazil': 'BR'
    }
    return codes[country] || 'US'
  }

  const getCityCoordinates = (city: string, country: string) => {
    const coordinates: { [key: string]: { lat: number; lng: number } } = {
      'New York': { lat: 40.7128, lng: -74.0060 },
      'London': { lat: 51.5074, lng: -0.1278 },
      'Los Angeles': { lat: 34.0522, lng: -118.2437 },
      'Toronto': { lat: 43.6532, lng: -79.3832 },
      'Paris': { lat: 48.8566, lng: 2.3522 },
      'Berlin': { lat: 52.5200, lng: 13.4050 },
      'Chicago': { lat: 41.8781, lng: -87.6298 },
      'Sydney': { lat: -33.8688, lng: 151.2093 }
    }
    return coordinates[city]
  }

  const currentData = useMemo(() => {
    if (!data) return []

    let sourceData: any[] = []
    
    if (viewMode === 'countries') {
      sourceData = data.countries
    } else if (viewMode === 'cities') {
      sourceData = data.cities
    } else {
      sourceData = data.regions
    }

    // Apply search filter
    if (searchTerm) {
      const searchKey = viewMode === 'countries' ? 'country' : viewMode === 'cities' ? 'city' : 'region'
      sourceData = sourceData.filter(item => 
        item[searchKey].toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    // Apply sorting
    if (sortBy === 'clicks') {
      sourceData = sourceData.sort((a, b) => b.clicks - a.clicks)
    } else {
      const sortKey = viewMode === 'countries' ? 'country' : viewMode === 'cities' ? 'city' : 'region'
      sourceData = sourceData.sort((a, b) => a[sortKey].localeCompare(b[sortKey]))
    }

    return sourceData
  }, [data, viewMode, searchTerm, sortBy])

  const chartColors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#06B6D4', '#84CC16', '#F97316']

  const exportData = () => {
    if (!currentData.length) return

    const headers = viewMode === 'countries' 
      ? ['Country', 'Country Code', 'Clicks', 'Unique Clicks', 'Percentage']
      : viewMode === 'cities'
      ? ['City', 'Country', 'Clicks', 'Unique Clicks', 'Percentage']
      : ['Region', 'Country', 'Clicks', 'Percentage']

    const rows = currentData.map(item => {
      if (viewMode === 'countries') {
        return [item.country, item.countryCode, item.clicks, item.uniqueClicks, `${item.percentage}%`]
      } else if (viewMode === 'cities') {
        return [item.city, item.country, item.clicks, item.uniqueClicks, `${item.percentage}%`]
      } else {
        return [item.region, item.country, item.clicks, `${item.percentage}%`]
      }
    })

    const csvContent = [headers, ...rows].map(row => row.join(',')).join('\n')
    const blob = new Blob([csvContent], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `geographic-${viewMode}-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const renderChart = () => {
    if (currentData.length === 0) {
      return (
        <div className="text-center py-12">
          <Globe className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Geographic Data</h3>
          <p className="text-gray-600">No location data found for the selected view.</p>
        </div>
      )
    }

    const displayKey = viewMode === 'countries' ? 'country' : viewMode === 'cities' ? 'city' : 'region'

    if (chartType === 'bar') {
      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={currentData.slice(0, 15)}
              margin={{ top: 20, right: 30, left: 20, bottom: 60 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey={displayKey}
                angle={-45}
                textAnchor="end"
                height={80}
                fontSize={12}
              />
              <YAxis />
              <Tooltip formatter={(value, name) => [value.toLocaleString(), name === 'clicks' ? 'Clicks' : 'Unique Clicks']} />
              <Legend />
              <Bar dataKey="clicks" fill="#3B82F6" name="Clicks" radius={[2, 2, 0, 0]} />
              {viewMode !== 'regions' && (
                <Bar dataKey="uniqueClicks" fill="#10B981" name="Unique Clicks" radius={[2, 2, 0, 0]} />
              )}
            </BarChart>
          </ResponsiveContainer>
        </div>
      )
    } else if (chartType === 'pie') {
      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={currentData.slice(0, 8)}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percentage }) => `${name}: ${percentage}%`}
                outerRadius={120}
                fill="#8884d8"
                dataKey="clicks"
                nameKey={displayKey}
              >
                {currentData.slice(0, 8).map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={chartColors[index % chartColors.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => [value.toLocaleString(), 'Clicks']} />
            </PieChart>
          </ResponsiveContainer>
        </div>
      )
    } else if (chartType === 'treemap') {
      const treemapData = currentData.slice(0, 12).map(item => ({
        name: item[displayKey],
        size: item.clicks,
        percentage: item.percentage
      }))

      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <Treemap
              data={treemapData}
              dataKey="size"
              ratio={4/3}
              stroke="#fff"
              fill="#3B82F6"
              content={({ root, depth, x, y, width, height, index, name, size }) => {
                if (depth === 1) {
                  return (
                    <g>
                      <rect
                        x={x}
                        y={y}
                        width={width}
                        height={height}
                        style={{
                          fill: chartColors[index % chartColors.length],
                          stroke: '#fff',
                          strokeWidth: 2,
                          strokeOpacity: 1,
                        }}
                      />
                      {width > 60 && height > 30 && (
                        <>
                          <text
                            x={x + width / 2}
                            y={y + height / 2 - 8}
                            textAnchor="middle"
                            fill="#fff"
                            fontSize="12"
                            fontWeight="bold"
                          >
                            {name}
                          </text>
                          <text
                            x={x + width / 2}
                            y={y + height / 2 + 8}
                            textAnchor="middle"
                            fill="#fff"
                            fontSize="10"
                          >
                            {size?.toLocaleString()} clicks
                          </text>
                        </>
                      )}
                    </g>
                  )
                }
              }}
            />
          </ResponsiveContainer>
        </div>
      )
    } else {
      return (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  {viewMode === 'countries' ? 'Country' : viewMode === 'cities' ? 'City' : 'Region'}
                </th>
                {viewMode === 'cities' && (
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Country
                  </th>
                )}
                {viewMode === 'regions' && (
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Country
                  </th>
                )}
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Clicks
                </th>
                {viewMode !== 'regions' && (
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Unique Clicks
                  </th>
                )}
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Percentage
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {currentData.map((item, index) => (
                <tr key={index} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50'}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    <div className="flex items-center">
                      <div 
                        className="w-3 h-3 rounded-full mr-3"
                        style={{ backgroundColor: chartColors[index % chartColors.length] }}
                      />
                      {item[displayKey]}
                    </div>
                  </td>
                  {(viewMode === 'cities' || viewMode === 'regions') && (
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {item.country}
                    </td>
                  )}
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {item.clicks.toLocaleString()}
                  </td>
                  {viewMode !== 'regions' && (
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {item.uniqueClicks.toLocaleString()}
                    </td>
                  )}
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {item.percentage.toFixed(1)}%
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )
    }
  }

  if (isLoading) {
    return <PageLoading message="Loading geographic data..." />
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <Globe className="h-5 w-5 text-red-500 mr-2" />
            <p className="text-sm text-red-700">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`bg-white rounded-lg shadow-sm border border-gray-200 ${className}`}>
      {/* Header */}
      <div className="p-6 border-b border-gray-200">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div>
            <h3 className="text-lg font-medium text-gray-900">Geographic Analytics</h3>
            <p className="text-sm text-gray-600 mt-1">
              Click distribution by location
            </p>
          </div>
          <div className="flex items-center space-x-3">
            <button
              onClick={() => refreshData()}
              disabled={isRefreshing}
              className="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
              Refresh
              {connectionStatus === 'reconnecting' && (
                <span className="ml-2 text-xs text-yellow-600">Reconnecting...</span>
              )}
              {connectionStatus === 'disconnected' && (
                <span className="ml-2 text-xs text-red-600">Disconnected</span>
              )}
            </button>
            <button
              onClick={exportData}
              className="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
            >
              <Download className="h-4 w-4 mr-2" />
              Export
            </button>
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="p-6 border-b border-gray-200 space-y-4">
        {/* View Mode Selector */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">View</label>
          <div className="flex flex-wrap gap-2">
            {viewModes.map(mode => (
              <button
                key={mode.value}
                onClick={() => setViewMode(mode.value)}
                className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                  viewMode === mode.value
                    ? 'bg-primary-100 text-primary-700'
                    : 'text-gray-600 hover:bg-gray-100'
                }`}
              >
                {mode.icon}
                <span className="ml-2">{mode.label}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Chart Type and Controls */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Chart Type</label>
            <div className="flex flex-wrap gap-2">
              {chartTypes.map(type => (
                <button
                  key={type.value}
                  onClick={() => setChartType(type.value)}
                  className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                    chartType === type.value
                      ? 'bg-primary-100 text-primary-700'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                >
                  {type.icon}
                  <span className="ml-2 hidden sm:inline">{type.label}</span>
                </button>
              ))}
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Search</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search className="h-4 w-4 text-gray-400" />
              </div>
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                placeholder={`Search ${viewMode}...`}
                className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Sort By</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as 'clicks' | 'alphabetical')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="clicks">Most Clicks</option>
              <option value="alphabetical">Alphabetical</option>
            </select>
          </div>
        </div>
      </div>

      {/* Summary Stats */}
      {data && (
        <div className="p-6 border-b border-gray-200">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{data.countries.length}</div>
              <div className="text-sm text-gray-600">Countries</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{data.cities.length}</div>
              <div className="text-sm text-gray-600">Cities</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.countries.reduce((sum, country) => sum + country.clicks, 0).toLocaleString()}
              </div>
              <div className="text-sm text-gray-600">Total Clicks</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.countries[0]?.country || 'N/A'}
              </div>
              <div className="text-sm text-gray-600">Top Country</div>
            </div>
          </div>
        </div>
      )}

      {/* Chart */}
      <div className="p-6">
        {renderChart()}
      </div>
    </div>
  )
}

export default GeographicMap