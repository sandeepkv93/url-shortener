import { useState, useEffect, useMemo } from 'react'
import { 
  BarChart3, 
  TrendingUp, 
  TrendingDown, 
  Users, 
  Globe, 
  Clock,
  MousePointer,
  Calendar,
  Filter,
  Download,
  RefreshCw,
  ArrowUp,
  ArrowDown,
  Link as LinkIcon,
  Eye,
  Smartphone,
  MapPin
} from 'lucide-react'
import { 
  LineChart, 
  Line, 
  AreaChart, 
  Area,
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
  ResponsiveContainer 
} from 'recharts'
import { format, subDays, parseISO } from 'date-fns'
import { urlService, urlAnalyticsService } from '@/services/urls'
import { AnalyticsPeriod } from '@/types/url'
import { PageLoading } from '@/components/common/Loading'

interface DashboardData {
  totalURLs: number
  totalClicks: number
  uniqueClicks: number
  avgClicksPerURL: number
  topURLs: Array<{
    id: string
    title: string
    shortCode: string
    originalUrl: string
    clickCount: number
    createdAt: string
  }>
  clickTimeline: Array<{
    date: string
    clicks: number
    uniqueClicks: number
  }>
  topCountries: Array<{
    country: string
    clicks: number
    percentage: number
  }>
  topDevices: Array<{
    device: string
    clicks: number
    percentage: number
  }>
  topReferrers: Array<{
    referrer: string
    clicks: number
    percentage: number
  }>
  recentActivity: Array<{
    id: string
    urlTitle: string
    shortCode: string
    timestamp: string
    country: string
    device: string
    browser: string
  }>
}

interface DashboardProps {
  className?: string
}

const Dashboard = ({ className = '' }: DashboardProps) => {
  const [data, setData] = useState<DashboardData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [selectedPeriod, setSelectedPeriod] = useState<AnalyticsPeriod>('30d')
  const [error, setError] = useState<string | null>(null)

  const periods: { value: AnalyticsPeriod; label: string }[] = [
    { value: '1h', label: 'Last Hour' },
    { value: '24h', label: 'Last 24 Hours' },
    { value: '7d', label: 'Last 7 Days' },
    { value: '30d', label: 'Last 30 Days' },
    { value: '90d', label: 'Last 90 Days' },
    { value: '1y', label: 'Last Year' },
    { value: 'all', label: 'All Time' }
  ]

  useEffect(() => {
    fetchDashboardData()
  }, [selectedPeriod])

  const fetchDashboardData = async (showLoading = true) => {
    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)
    setError(null)

    try {
      // Fetch user's URLs and aggregate analytics
      const urlsResponse = await urlService.getUserURLs({ 
        page: 1, 
        limit: 100,
        sortBy: 'clickCount',
        sortOrder: 'desc'
      })

      // Calculate basic metrics
      const totalURLs = urlsResponse.total
      const totalClicks = urlsResponse.urls.reduce((sum, url) => sum + url.clickCount, 0)
      const uniqueClicks = Math.floor(totalClicks * 0.8) // Estimate unique clicks
      const avgClicksPerURL = totalURLs > 0 ? totalClicks / totalURLs : 0

      // Get top performing URLs
      const topURLs = urlsResponse.urls.slice(0, 5).map(url => ({
        id: url.id,
        title: url.title || url.originalUrl,
        shortCode: url.shortCode,
        originalUrl: url.originalUrl,
        clickCount: url.clickCount,
        createdAt: url.createdAt
      }))

      // Mock analytics data - in real implementation, this would come from backend
      const clickTimeline = generateMockTimeline(selectedPeriod)
      const topCountries = generateMockCountryData()
      const topDevices = generateMockDeviceData()
      const topReferrers = generateMockReferrerData()
      const recentActivity = generateMockActivity()

      setData({
        totalURLs,
        totalClicks,
        uniqueClicks,
        avgClicksPerURL,
        topURLs,
        clickTimeline,
        topCountries,
        topDevices,
        topReferrers,
        recentActivity
      })
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard data')
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  const generateMockTimeline = (period: AnalyticsPeriod) => {
    const days = period === '7d' ? 7 : period === '30d' ? 30 : period === '90d' ? 90 : 7
    const timeline = []
    
    for (let i = days - 1; i >= 0; i--) {
      const date = format(subDays(new Date(), i), 'yyyy-MM-dd')
      const clicks = Math.floor(Math.random() * 100) + 10
      const uniqueClicks = Math.floor(clicks * 0.8)
      
      timeline.push({ date, clicks, uniqueClicks })
    }
    
    return timeline
  }

  const generateMockCountryData = () => [
    { country: 'United States', clicks: 450, percentage: 45 },
    { country: 'United Kingdom', clicks: 200, percentage: 20 },
    { country: 'Canada', clicks: 150, percentage: 15 },
    { country: 'Germany', clicks: 100, percentage: 10 },
    { country: 'Australia', clicks: 100, percentage: 10 }
  ]

  const generateMockDeviceData = () => [
    { device: 'Desktop', clicks: 500, percentage: 50 },
    { device: 'Mobile', clicks: 350, percentage: 35 },
    { device: 'Tablet', clicks: 150, percentage: 15 }
  ]

  const generateMockReferrerData = () => [
    { referrer: 'Direct', clicks: 400, percentage: 40 },
    { referrer: 'Google', clicks: 300, percentage: 30 },
    { referrer: 'Twitter', clicks: 150, percentage: 15 },
    { referrer: 'Facebook', clicks: 100, percentage: 10 },
    { referrer: 'LinkedIn', clicks: 50, percentage: 5 }
  ]

  const generateMockActivity = () => [
    {
      id: '1',
      urlTitle: 'Product Launch Page',
      shortCode: 'prod123',
      timestamp: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
      country: 'United States',
      device: 'Desktop',
      browser: 'Chrome'
    },
    {
      id: '2',
      urlTitle: 'Marketing Campaign',
      shortCode: 'camp456',
      timestamp: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
      country: 'United Kingdom',
      device: 'Mobile',
      browser: 'Safari'
    },
    {
      id: '3',
      urlTitle: 'Blog Post',
      shortCode: 'blog789',
      timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
      country: 'Canada',
      device: 'Desktop',
      browser: 'Firefox'
    }
  ]

  const chartColors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#06B6D4']

  const kpiChange = useMemo(() => {
    if (!data?.clickTimeline || data.clickTimeline.length < 2) return null
    
    const recent = data.clickTimeline.slice(-7).reduce((sum, day) => sum + day.clicks, 0)
    const previous = data.clickTimeline.slice(-14, -7).reduce((sum, day) => sum + day.clicks, 0)
    
    if (previous === 0) return null
    
    const change = ((recent - previous) / previous) * 100
    return {
      clicks: change,
      isPositive: change > 0
    }
  }, [data])

  const exportData = () => {
    // Implementation for exporting dashboard data
    console.log('Exporting dashboard data:', data)
  }

  if (isLoading) {
    return <PageLoading message="Loading dashboard..." />
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <BarChart3 className="h-5 w-5 text-red-500 mr-2" />
            <p className="text-sm text-red-700">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="max-w-7xl mx-auto p-6">
        <div className="text-center py-12">
          <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Data Available</h3>
          <p className="text-gray-600">Create some URLs to see analytics data.</p>
        </div>
      </div>
    )
  }

  return (
    <div className={`max-w-7xl mx-auto ${className}`}>
      {/* Header */}
      <div className="mb-8">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Analytics Dashboard</h1>
            <p className="text-gray-600 mt-1">
              Overview of your URL performance and user engagement
            </p>
          </div>
          <div className="flex items-center space-x-3">
            <button
              onClick={() => fetchDashboardData(false)}
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

      {/* Period Selector */}
      <div className="mb-8">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">Time Period</h3>
            <div className="flex items-center space-x-2">
              {periods.map(period => (
                <button
                  key={period.value}
                  onClick={() => setSelectedPeriod(period.value)}
                  className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                    selectedPeriod === period.value
                      ? 'bg-primary-100 text-primary-700'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                >
                  {period.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total URLs</p>
              <p className="text-3xl font-bold text-gray-900">{data.totalURLs}</p>
            </div>
            <div className="h-12 w-12 bg-blue-100 rounded-lg flex items-center justify-center">
              <LinkIcon className="h-6 w-6 text-blue-600" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Clicks</p>
              <p className="text-3xl font-bold text-gray-900">{data.totalClicks.toLocaleString()}</p>
              {kpiChange && (
                <div className={`flex items-center mt-1 ${kpiChange.isPositive ? 'text-green-600' : 'text-red-600'}`}>
                  {kpiChange.isPositive ? <ArrowUp className="h-3 w-3 mr-1" /> : <ArrowDown className="h-3 w-3 mr-1" />}
                  <span className="text-xs font-medium">
                    {Math.abs(kpiChange.clicks).toFixed(1)}% vs last week
                  </span>
                </div>
              )}
            </div>
            <div className="h-12 w-12 bg-green-100 rounded-lg flex items-center justify-center">
              <MousePointer className="h-6 w-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Unique Visitors</p>
              <p className="text-3xl font-bold text-gray-900">{data.uniqueClicks.toLocaleString()}</p>
              <p className="text-xs text-gray-500 mt-1">
                {((data.uniqueClicks / data.totalClicks) * 100).toFixed(1)}% of total clicks
              </p>
            </div>
            <div className="h-12 w-12 bg-purple-100 rounded-lg flex items-center justify-center">
              <Users className="h-6 w-6 text-purple-600" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Avg. Clicks/URL</p>
              <p className="text-3xl font-bold text-gray-900">{data.avgClicksPerURL.toFixed(1)}</p>
            </div>
            <div className="h-12 w-12 bg-orange-100 rounded-lg flex items-center justify-center">
              <BarChart3 className="h-6 w-6 text-orange-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Click Timeline */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Click Timeline</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={data.clickTimeline}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="date" 
                  tickFormatter={(value) => format(parseISO(value), 'MMM d')}
                />
                <YAxis />
                <Tooltip 
                  labelFormatter={(value) => format(parseISO(value), 'MMM d, yyyy')}
                  formatter={(value, name) => [value, name === 'clicks' ? 'Total Clicks' : 'Unique Clicks']}
                />
                <Legend />
                <Area
                  type="monotone"
                  dataKey="clicks"
                  stackId="1"
                  stroke="#3B82F6"
                  fill="#3B82F6"
                  fillOpacity={0.6}
                  name="Total Clicks"
                />
                <Area
                  type="monotone"
                  dataKey="uniqueClicks"
                  stackId="2"
                  stroke="#10B981"
                  fill="#10B981"
                  fillOpacity={0.6}
                  name="Unique Clicks"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Top Countries */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Top Countries</h3>
          <div className="space-y-4">
            {data.topCountries.map((country, index) => (
              <div key={country.country} className="flex items-center justify-between">
                <div className="flex items-center">
                  <div 
                    className="w-3 h-3 rounded-full mr-3"
                    style={{ backgroundColor: chartColors[index % chartColors.length] }}
                  />
                  <span className="text-sm font-medium text-gray-900">{country.country}</span>
                </div>
                <div className="text-right">
                  <div className="text-sm font-medium text-gray-900">{country.clicks}</div>
                  <div className="text-xs text-gray-500">{country.percentage}%</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Bottom Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Top URLs */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Top Performing URLs</h3>
          <div className="space-y-4">
            {data.topURLs.map((url, index) => (
              <div key={url.id} className="flex items-center justify-between">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center">
                    <span className="text-xs font-medium text-gray-500 mr-2">#{index + 1}</span>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">{url.title}</p>
                      <p className="text-xs text-gray-500 truncate">/{url.shortCode}</p>
                    </div>
                  </div>
                </div>
                <div className="text-right ml-4">
                  <div className="text-sm font-medium text-gray-900">{url.clickCount}</div>
                  <div className="text-xs text-gray-500">clicks</div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Device Breakdown */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Device Types</h3>
          <div className="space-y-4">
            {data.topDevices.map((device, index) => (
              <div key={device.device} className="flex items-center justify-between">
                <div className="flex items-center">
                  <div 
                    className="w-3 h-3 rounded-full mr-3"
                    style={{ backgroundColor: chartColors[index % chartColors.length] }}
                  />
                  <span className="text-sm font-medium text-gray-900">{device.device}</span>
                </div>
                <div className="text-right">
                  <div className="text-sm font-medium text-gray-900">{device.clicks}</div>
                  <div className="text-xs text-gray-500">{device.percentage}%</div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
          <div className="space-y-4">
            {data.recentActivity.map(activity => (
              <div key={activity.id} className="flex items-start space-x-3">
                <div className="flex-shrink-0">
                  <div className="h-8 w-8 bg-blue-100 rounded-full flex items-center justify-center">
                    <Eye className="h-4 w-4 text-blue-600" />
                  </div>
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900 truncate">
                    {activity.urlTitle}
                  </p>
                  <div className="flex items-center text-xs text-gray-500 mt-1">
                    <MapPin className="h-3 w-3 mr-1" />
                    <span>{activity.country}</span>
                    <span className="mx-1">•</span>
                    <Smartphone className="h-3 w-3 mr-1" />
                    <span>{activity.device}</span>
                  </div>
                </div>
                <div className="flex-shrink-0 text-xs text-gray-500">
                  {format(new Date(activity.timestamp), 'h:mm a')}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard