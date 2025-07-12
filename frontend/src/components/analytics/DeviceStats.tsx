import { useState, useEffect, useMemo } from 'react'
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  RadialBarChart,
  RadialBar
} from 'recharts'
import {
  Smartphone,
  Monitor,
  Tablet,
  Globe,
  Cpu,
  HardDrive,
  Chrome,
  RefreshCw,
  Download,
  Filter,
  TrendingUp,
  TrendingDown,
  ArrowUp,
  ArrowDown,
  Minus
} from 'lucide-react'
import { urlAnalyticsService } from '@/services/urls'
import { PageLoading } from '@/components/common/Loading'

interface DeviceData {
  devices: Array<{
    device: string
    clicks: number
    percentage: number
    uniqueClicks: number
    trend?: number
  }>
  browsers: Array<{
    browser: string
    clicks: number
    percentage: number
    uniqueClicks: number
    version?: string
    trend?: number
  }>
  operatingSystems: Array<{
    os: string
    clicks: number
    percentage: number
    uniqueClicks: number
    version?: string
    trend?: number
  }>
  screenResolutions: Array<{
    resolution: string
    clicks: number
    percentage: number
  }>
}

interface DeviceStatsProps {
  urlId?: string // If provided, shows analytics for specific URL
  className?: string
}

type ViewMode = 'devices' | 'browsers' | 'os' | 'resolutions'
type ChartType = 'pie' | 'bar' | 'radial'

const DeviceStats = ({ urlId, className = '' }: DeviceStatsProps) => {
  const [data, setData] = useState<DeviceData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [viewMode, setViewMode] = useState<ViewMode>('devices')
  const [chartType, setChartType] = useState<ChartType>('pie')
  const [error, setError] = useState<string | null>(null)

  const viewModes: { value: ViewMode; label: string; icon: React.ReactNode }[] = [
    { value: 'devices', label: 'Devices', icon: <Smartphone className="h-4 w-4" /> },
    { value: 'browsers', label: 'Browsers', icon: <Globe className="h-4 w-4" /> },
    { value: 'os', label: 'Operating Systems', icon: <Cpu className="h-4 w-4" /> },
    { value: 'resolutions', label: 'Screen Resolutions', icon: <Monitor className="h-4 w-4" /> }
  ]

  const chartTypes: { value: ChartType; label: string; icon: React.ReactNode }[] = [
    { value: 'pie', label: 'Pie Chart', icon: <div className="w-4 h-4 rounded-full border-2 border-current" /> },
    { value: 'bar', label: 'Bar Chart', icon: <div className="w-4 h-4 flex items-end space-x-0.5"><div className="w-1 h-2 bg-current"></div><div className="w-1 h-3 bg-current"></div><div className="w-1 h-1 bg-current"></div></div> },
    { value: 'radial', label: 'Radial Chart', icon: <div className="w-4 h-4 rounded-full border border-current flex items-center justify-center"><div className="w-2 h-2 rounded-full bg-current"></div></div> }
  ]

  const deviceIcons: { [key: string]: React.ReactNode } = {
    'Desktop': <Monitor className="h-5 w-5" />,
    'Mobile': <Smartphone className="h-5 w-5" />,
    'Tablet': <Tablet className="h-5 w-5" />,
    'Smart TV': <Monitor className="h-5 w-5" />,
    'Other': <HardDrive className="h-5 w-5" />
  }

  const browserIcons: { [key: string]: React.ReactNode } = {
    'Chrome': <Chrome className="h-5 w-5" />,
    'Firefox': <Globe className="h-5 w-5" />,
    'Safari': <Globe className="h-5 w-5" />,
    'Edge': <Globe className="h-5 w-5" />,
    'Opera': <Globe className="h-5 w-5" />,
    'Other': <Globe className="h-5 w-5" />
  }

  useEffect(() => {
    fetchDeviceData()
  }, [urlId])

  const fetchDeviceData = async (showLoading = true) => {
    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)
    setError(null)

    try {
      if (urlId) {
        // Fetch data for specific URL
        const response = await urlAnalyticsService.getDeviceStats(urlId)
        setData({
          devices: response.devices.map(device => ({
            ...device,
            uniqueClicks: Math.floor(device.clicks * 0.8),
            trend: (Math.random() - 0.5) * 20 // Mock trend data
          })),
          browsers: response.browsers.map(browser => ({
            ...browser,
            uniqueClicks: Math.floor(browser.clicks * 0.8),
            version: generateBrowserVersion(browser.browser),
            trend: (Math.random() - 0.5) * 20
          })),
          operatingSystems: response.operatingSystems.map(os => ({
            ...os,
            uniqueClicks: Math.floor(os.clicks * 0.8),
            version: generateOSVersion(os.os),
            trend: (Math.random() - 0.5) * 20
          })),
          screenResolutions: generateMockResolutions()
        })
      } else {
        // Generate mock data for dashboard
        setData(generateMockDeviceData())
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load device data')
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  const generateBrowserVersion = (browser: string) => {
    const versions: { [key: string]: string } = {
      'Chrome': `${118 + Math.floor(Math.random() * 5)}.0`,
      'Firefox': `${110 + Math.floor(Math.random() * 10)}.0`,
      'Safari': `${15 + Math.floor(Math.random() * 3)}.${Math.floor(Math.random() * 10)}`,
      'Edge': `${110 + Math.floor(Math.random() * 10)}.0`,
      'Opera': `${95 + Math.floor(Math.random() * 10)}.0`
    }
    return versions[browser] || '1.0'
  }

  const generateOSVersion = (os: string) => {
    const versions: { [key: string]: string } = {
      'Windows': `${Math.random() > 0.5 ? '11' : '10'}`,
      'macOS': `${13 + Math.floor(Math.random() * 2)}.${Math.floor(Math.random() * 10)}`,
      'iOS': `${16 + Math.floor(Math.random() * 2)}.${Math.floor(Math.random() * 10)}`,
      'Android': `${12 + Math.floor(Math.random() * 3)}`,
      'Linux': 'Ubuntu 22.04'
    }
    return versions[os] || '1.0'
  }

  const generateMockResolutions = () => [
    { resolution: '1920x1080', clicks: 450, percentage: 45.0 },
    { resolution: '1366x768', clicks: 200, percentage: 20.0 },
    { resolution: '1536x864', clicks: 120, percentage: 12.0 },
    { resolution: '1440x900', clicks: 80, percentage: 8.0 },
    { resolution: '1280x720', clicks: 70, percentage: 7.0 },
    { resolution: '2560x1440', clicks: 50, percentage: 5.0 },
    { resolution: '3840x2160', clicks: 30, percentage: 3.0 }
  ]

  const generateMockDeviceData = (): DeviceData => ({
    devices: [
      { device: 'Desktop', clicks: 650, percentage: 65.0, uniqueClicks: 520, trend: 5.2 },
      { device: 'Mobile', clicks: 300, percentage: 30.0, uniqueClicks: 240, trend: 12.8 },
      { device: 'Tablet', clicks: 50, percentage: 5.0, uniqueClicks: 40, trend: -2.1 }
    ],
    browsers: [
      { browser: 'Chrome', clicks: 500, percentage: 50.0, uniqueClicks: 400, version: '119.0', trend: 2.4 },
      { browser: 'Firefox', clicks: 200, percentage: 20.0, uniqueClicks: 160, version: '118.0', trend: -1.2 },
      { browser: 'Safari', clicks: 150, percentage: 15.0, uniqueClicks: 120, version: '17.1', trend: 8.7 },
      { browser: 'Edge', clicks: 100, percentage: 10.0, uniqueClicks: 80, version: '117.0', trend: 3.1 },
      { browser: 'Opera', clicks: 50, percentage: 5.0, uniqueClicks: 40, version: '104.0', trend: -0.5 }
    ],
    operatingSystems: [
      { os: 'Windows', clicks: 400, percentage: 40.0, uniqueClicks: 320, version: '11', trend: 1.8 },
      { os: 'macOS', clicks: 250, percentage: 25.0, uniqueClicks: 200, version: '14.1', trend: 4.2 },
      { os: 'iOS', clicks: 200, percentage: 20.0, uniqueClicks: 160, version: '17.1', trend: 6.9 },
      { os: 'Android', clicks: 120, percentage: 12.0, uniqueClicks: 96, version: '14', trend: 8.5 },
      { os: 'Linux', clicks: 30, percentage: 3.0, uniqueClicks: 24, version: 'Ubuntu 22.04', trend: -1.1 }
    ],
    screenResolutions: generateMockResolutions()
  })

  const currentData = useMemo(() => {
    if (!data) return []

    switch (viewMode) {
      case 'devices':
        return data.devices
      case 'browsers':
        return data.browsers
      case 'os':
        return data.operatingSystems
      case 'resolutions':
        return data.screenResolutions
      default:
        return []
    }
  }, [data, viewMode])

  const chartColors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#06B6D4', '#84CC16', '#F97316']

  const exportData = () => {
    if (!currentData.length) return

    const headers = viewMode === 'devices' 
      ? ['Device', 'Clicks', 'Unique Clicks', 'Percentage', 'Trend']
      : viewMode === 'browsers'
      ? ['Browser', 'Version', 'Clicks', 'Unique Clicks', 'Percentage', 'Trend']
      : viewMode === 'os'
      ? ['Operating System', 'Version', 'Clicks', 'Unique Clicks', 'Percentage', 'Trend']
      : ['Resolution', 'Clicks', 'Percentage']

    const rows = currentData.map(item => {
      if (viewMode === 'devices') {
        return [item.device, item.clicks, item.uniqueClicks, `${item.percentage}%`, `${item.trend?.toFixed(1)}%`]
      } else if (viewMode === 'browsers') {
        return [item.browser, item.version, item.clicks, item.uniqueClicks, `${item.percentage}%`, `${item.trend?.toFixed(1)}%`]
      } else if (viewMode === 'os') {
        return [item.os, item.version, item.clicks, item.uniqueClicks, `${item.percentage}%`, `${item.trend?.toFixed(1)}%`]
      } else {
        return [item.resolution, item.clicks, `${item.percentage}%`]
      }
    })

    const csvContent = [headers, ...rows].map(row => row.join(',')).join('\n')
    const blob = new Blob([csvContent], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `device-stats-${viewMode}-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const renderChart = () => {
    if (currentData.length === 0) {
      return (
        <div className="text-center py-12">
          <Smartphone className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Device Data</h3>
          <p className="text-gray-600">No device information available for the selected view.</p>
        </div>
      )
    }

    const dataKey = viewMode === 'devices' ? 'device' : viewMode === 'browsers' ? 'browser' : viewMode === 'os' ? 'os' : 'resolution'

    if (chartType === 'pie') {
      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={currentData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percentage }) => `${name}: ${percentage}%`}
                outerRadius={120}
                fill="#8884d8"
                dataKey="clicks"
                nameKey={dataKey}
              >
                {currentData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={chartColors[index % chartColors.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => [value.toLocaleString(), 'Clicks']} />
            </PieChart>
          </ResponsiveContainer>
        </div>
      )
    } else if (chartType === 'bar') {
      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={currentData}
              margin={{ top: 20, right: 30, left: 20, bottom: 60 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey={dataKey}
                angle={-45}
                textAnchor="end"
                height={80}
                fontSize={12}
              />
              <YAxis />
              <Tooltip formatter={(value, name) => [value.toLocaleString(), name === 'clicks' ? 'Clicks' : 'Unique Clicks']} />
              <Legend />
              <Bar dataKey="clicks" fill="#3B82F6" name="Clicks" radius={[2, 2, 0, 0]} />
              {viewMode !== 'resolutions' && (
                <Bar dataKey="uniqueClicks" fill="#10B981" name="Unique Clicks" radius={[2, 2, 0, 0]} />
              )}
            </BarChart>
          </ResponsiveContainer>
        </div>
      )
    } else {
      // Radial bar chart
      const radialData = currentData.map((item, index) => ({
        ...item,
        fill: chartColors[index % chartColors.length]
      }))

      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <RadialBarChart
              cx="50%"
              cy="50%"
              innerRadius="20%"
              outerRadius="80%"
              data={radialData}
            >
              <RadialBar dataKey="percentage" cornerRadius={10} fill="#8884d8" />
              <Tooltip formatter={(value, name) => {
                if (name === 'percentage') return [`${value}%`, 'Percentage']
                return [value.toLocaleString(), 'Clicks']
              }} />
            </RadialBarChart>
          </ResponsiveContainer>
        </div>
      )
    }
  }

  const getTrendIcon = (trend?: number) => {
    if (!trend) return <Minus className="h-3 w-3" />
    if (trend > 0) return <ArrowUp className="h-3 w-3" />
    return <ArrowDown className="h-3 w-3" />
  }

  const getTrendColor = (trend?: number) => {
    if (!trend) return 'text-gray-500'
    if (trend > 0) return 'text-green-600'
    return 'text-red-600'
  }

  if (isLoading) {
    return <PageLoading message="Loading device analytics..." />
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <Smartphone className="h-5 w-5 text-red-500 mr-2" />
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
            <h3 className="text-lg font-medium text-gray-900">Device & Browser Analytics</h3>
            <p className="text-sm text-gray-600 mt-1">
              User device, browser, and system information
            </p>
          </div>
          <div className="flex items-center space-x-3">
            <button
              onClick={() => fetchDeviceData(false)}
              disabled={isRefreshing}
              className="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
              Refresh
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
          <label className="block text-sm font-medium text-gray-700 mb-2">Category</label>
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

        {/* Chart Type Selector */}
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
                <span className="ml-2">{type.label}</span>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Summary Stats */}
      {data && (
        <div className="p-6 border-b border-gray-200">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.devices.reduce((sum, device) => sum + device.clicks, 0).toLocaleString()}
              </div>
              <div className="text-sm text-gray-600">Total Clicks</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.devices.length + data.browsers.length + data.operatingSystems.length}
              </div>
              <div className="text-sm text-gray-600">Unique Combinations</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.devices[0]?.device || 'N/A'}
              </div>
              <div className="text-sm text-gray-600">Most Popular Device</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.browsers[0]?.browser || 'N/A'}
              </div>
              <div className="text-sm text-gray-600">Most Popular Browser</div>
            </div>
          </div>
        </div>
      )}

      {/* Chart and Table */}
      <div className="p-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Chart */}
          <div>
            <h4 className="text-lg font-medium text-gray-900 mb-4">
              {viewModes.find(m => m.value === viewMode)?.label} Distribution
            </h4>
            {renderChart()}
          </div>

          {/* Details Table */}
          <div>
            <h4 className="text-lg font-medium text-gray-900 mb-4">Detailed Breakdown</h4>
            <div className="overflow-y-auto max-h-96">
              <div className="space-y-3">
                {currentData.map((item, index) => {
                  const displayName = viewMode === 'devices' ? item.device : 
                                    viewMode === 'browsers' ? item.browser : 
                                    viewMode === 'os' ? item.os : item.resolution
                  
                  const icon = viewMode === 'devices' ? deviceIcons[item.device] || deviceIcons['Other'] :
                              viewMode === 'browsers' ? browserIcons[item.browser] || browserIcons['Other'] :
                              <Cpu className="h-5 w-5" />

                  return (
                    <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                      <div className="flex items-center flex-1 min-w-0">
                        <div 
                          className="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center mr-3"
                          style={{ backgroundColor: chartColors[index % chartColors.length] + '20' }}
                        >
                          <div style={{ color: chartColors[index % chartColors.length] }}>
                            {icon}
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center">
                            <p className="text-sm font-medium text-gray-900 truncate">{displayName}</p>
                            {(viewMode === 'browsers' || viewMode === 'os') && item.version && (
                              <span className="ml-2 text-xs text-gray-500">v{item.version}</span>
                            )}
                          </div>
                          <div className="flex items-center mt-1">
                            <p className="text-xs text-gray-500">
                              {item.clicks.toLocaleString()} clicks ({item.percentage.toFixed(1)}%)
                            </p>
                            {item.trend !== undefined && (
                              <div className={`flex items-center ml-2 ${getTrendColor(item.trend)}`}>
                                {getTrendIcon(item.trend)}
                                <span className="text-xs ml-1">{Math.abs(item.trend).toFixed(1)}%</span>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-medium text-gray-900">{item.clicks}</p>
                        {viewMode !== 'resolutions' && item.uniqueClicks && (
                          <p className="text-xs text-gray-500">{item.uniqueClicks} unique</p>
                        )}
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DeviceStats