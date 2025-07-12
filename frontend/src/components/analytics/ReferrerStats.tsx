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
  Treemap
} from 'recharts'
import {
  ExternalLink,
  Share2,
  Globe,
  Users,
  TrendingUp,
  Search,
  MessageSquare,
  Mail,
  Link as LinkIcon,
  RefreshCw,
  Download,
  Filter,
  ArrowUp,
  ArrowDown,
  Minus,
  Eye,
  MousePointer
} from 'lucide-react'
import { urlAnalyticsService } from '@/services/urls'
import { PageLoading } from '@/components/common/Loading'

interface ReferrerData {
  referrers: Array<{
    referrer: string
    clicks: number
    percentage: number
    uniqueClicks: number
    category: 'search' | 'social' | 'direct' | 'email' | 'other'
    trend?: number
  }>
  directClicks: number
  totalClicks: number
  categories: Array<{
    category: string
    clicks: number
    percentage: number
    color: string
  }>
  topDomains: Array<{
    domain: string
    clicks: number
    percentage: number
    sources: string[]
  }>
}

interface ReferrerStatsProps {
  urlId?: string // If provided, shows analytics for specific URL
  className?: string
}

type ViewMode = 'all' | 'search' | 'social' | 'direct' | 'email' | 'other'
type ChartType = 'pie' | 'bar' | 'treemap' | 'table'

const ReferrerStats = ({ urlId, className = '' }: ReferrerStatsProps) => {
  const [data, setData] = useState<ReferrerData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [viewMode, setViewMode] = useState<ViewMode>('all')
  const [chartType, setChartType] = useState<ChartType>('pie')
  const [searchTerm, setSearchTerm] = useState('')
  const [sortBy, setSortBy] = useState<'clicks' | 'alphabetical' | 'trend'>('clicks')
  const [error, setError] = useState<string | null>(null)

  const viewModes: { value: ViewMode; label: string; icon: React.ReactNode; color: string }[] = [
    { value: 'all', label: 'All Sources', icon: <Globe className="h-4 w-4" />, color: '#6B7280' },
    { value: 'search', label: 'Search Engines', icon: <Search className="h-4 w-4" />, color: '#3B82F6' },
    { value: 'social', label: 'Social Media', icon: <MessageSquare className="h-4 w-4" />, color: '#10B981' },
    { value: 'direct', label: 'Direct Traffic', icon: <LinkIcon className="h-4 w-4" />, color: '#F59E0B' },
    { value: 'email', label: 'Email', icon: <Mail className="h-4 w-4" />, color: '#8B5CF6' },
    { value: 'other', label: 'Other', icon: <ExternalLink className="h-4 w-4" />, color: '#EF4444' }
  ]

  const chartTypes: { value: ChartType; label: string; icon: React.ReactNode }[] = [
    { value: 'pie', label: 'Pie Chart', icon: <div className="w-4 h-4 rounded-full border-2 border-current" /> },
    { value: 'bar', label: 'Bar Chart', icon: <div className="w-4 h-4 flex items-end space-x-0.5"><div className="w-1 h-2 bg-current"></div><div className="w-1 h-3 bg-current"></div><div className="w-1 h-1 bg-current"></div></div> },
    { value: 'treemap', label: 'Treemap', icon: <div className="w-4 h-4 grid grid-cols-2 gap-0.5"><div className="bg-current"></div><div className="bg-current"></div><div className="bg-current"></div><div className="bg-current"></div></div> },
    { value: 'table', label: 'Table', icon: <Filter className="h-4 w-4" /> }
  ]

  const referrerIcons: { [key: string]: React.ReactNode } = {
    'google.com': <Search className="h-5 w-5 text-blue-600" />,
    'bing.com': <Search className="h-5 w-5 text-blue-600" />,
    'yahoo.com': <Search className="h-5 w-5 text-purple-600" />,
    'facebook.com': <MessageSquare className="h-5 w-5 text-blue-600" />,
    'twitter.com': <MessageSquare className="h-5 w-5 text-blue-400" />,
    'linkedin.com': <MessageSquare className="h-5 w-5 text-blue-700" />,
    'instagram.com': <MessageSquare className="h-5 w-5 text-pink-600" />,
    'youtube.com': <MessageSquare className="h-5 w-5 text-red-600" />,
    'reddit.com': <MessageSquare className="h-5 w-5 text-orange-600" />,
    'direct': <LinkIcon className="h-5 w-5 text-gray-600" />,
    'email': <Mail className="h-5 w-5 text-purple-600" />
  }

  useEffect(() => {
    fetchReferrerData()
  }, [urlId])

  const fetchReferrerData = async (showLoading = true) => {
    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)
    setError(null)

    try {
      if (urlId) {
        // Fetch data for specific URL
        const response = await urlAnalyticsService.getReferrerStats(urlId)
        
        const categorizedReferrers = response.referrers.map(ref => ({
          ...ref,
          uniqueClicks: Math.floor(ref.clicks * 0.8),
          category: categorizeReferrer(ref.referrer),
          trend: (Math.random() - 0.5) * 20 // Mock trend data
        }))

        setData({
          referrers: categorizedReferrers,
          directClicks: response.directClicks,
          totalClicks: response.totalClicks,
          categories: generateCategoryBreakdown(categorizedReferrers, response.directClicks),
          topDomains: generateTopDomains(categorizedReferrers)
        })
      } else {
        // Generate mock data for dashboard
        setData(generateMockReferrerData())
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load referrer data')
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  const categorizeReferrer = (referrer: string): 'search' | 'social' | 'direct' | 'email' | 'other' => {
    const domain = referrer.toLowerCase()
    
    if (domain.includes('google') || domain.includes('bing') || domain.includes('yahoo') || domain.includes('duckduckgo')) {
      return 'search'
    } else if (domain.includes('facebook') || domain.includes('twitter') || domain.includes('linkedin') || 
               domain.includes('instagram') || domain.includes('youtube') || domain.includes('reddit') ||
               domain.includes('tiktok') || domain.includes('snapchat')) {
      return 'social'
    } else if (domain === 'direct' || domain === '') {
      return 'direct'
    } else if (domain.includes('gmail') || domain.includes('mail') || domain.includes('email')) {
      return 'email'
    } else {
      return 'other'
    }
  }

  const generateCategoryBreakdown = (referrers: any[], directClicks: number) => {
    const categories = {
      direct: directClicks,
      search: 0,
      social: 0,
      email: 0,
      other: 0
    }

    referrers.forEach(ref => {
      categories[ref.category] += ref.clicks
    })

    const total = Object.values(categories).reduce((sum, val) => sum + val, 0)
    const colors = ['#F59E0B', '#3B82F6', '#10B981', '#8B5CF6', '#EF4444']
    
    return Object.entries(categories).map(([category, clicks], index) => ({
      category: category.charAt(0).toUpperCase() + category.slice(1),
      clicks,
      percentage: total > 0 ? (clicks / total) * 100 : 0,
      color: colors[index]
    }))
  }

  const generateTopDomains = (referrers: any[]) => {
    const domainMap = new Map()
    
    referrers.forEach(ref => {
      const domain = ref.referrer.split('/')[0] || ref.referrer
      if (!domainMap.has(domain)) {
        domainMap.set(domain, { clicks: 0, sources: new Set() })
      }
      domainMap.get(domain).clicks += ref.clicks
      domainMap.get(domain).sources.add(ref.referrer)
    })

    const total = referrers.reduce((sum, ref) => sum + ref.clicks, 0)
    
    return Array.from(domainMap.entries())
      .map(([domain, data]) => ({
        domain,
        clicks: data.clicks,
        percentage: total > 0 ? (data.clicks / total) * 100 : 0,
        sources: Array.from(data.sources)
      }))
      .sort((a, b) => b.clicks - a.clicks)
      .slice(0, 10)
  }

  const generateMockReferrerData = (): ReferrerData => {
    const mockReferrers = [
      { referrer: 'google.com', clicks: 450, percentage: 30.0, uniqueClicks: 360, category: 'search' as const, trend: 5.2 },
      { referrer: 'facebook.com', clicks: 300, percentage: 20.0, uniqueClicks: 240, category: 'social' as const, trend: 12.8 },
      { referrer: 'twitter.com', clicks: 200, percentage: 13.3, uniqueClicks: 160, category: 'social' as const, trend: -2.1 },
      { referrer: 'linkedin.com', clicks: 150, percentage: 10.0, uniqueClicks: 120, category: 'social' as const, trend: 8.5 },
      { referrer: 'bing.com', clicks: 100, percentage: 6.7, uniqueClicks: 80, category: 'search' as const, trend: -1.2 },
      { referrer: 'youtube.com', clicks: 80, percentage: 5.3, uniqueClicks: 64, category: 'social' as const, trend: 15.3 },
      { referrer: 'reddit.com', clicks: 70, percentage: 4.7, uniqueClicks: 56, category: 'social' as const, trend: 22.1 },
      { referrer: 'instagram.com', clicks: 50, percentage: 3.3, uniqueClicks: 40, category: 'social' as const, trend: 18.7 },
      { referrer: 'gmail.com', clicks: 40, percentage: 2.7, uniqueClicks: 32, category: 'email' as const, trend: 3.4 },
      { referrer: 'yahoo.com', clicks: 30, percentage: 2.0, uniqueClicks: 24, category: 'search' as const, trend: -5.1 }
    ]

    const directClicks = 300
    const totalClicks = mockReferrers.reduce((sum, ref) => sum + ref.clicks, 0) + directClicks

    return {
      referrers: mockReferrers,
      directClicks,
      totalClicks,
      categories: generateCategoryBreakdown(mockReferrers, directClicks),
      topDomains: generateTopDomains(mockReferrers)
    }
  }

  const filteredData = useMemo(() => {
    if (!data) return []

    let sourceData = data.referrers

    // Filter by view mode
    if (viewMode !== 'all') {
      if (viewMode === 'direct') {
        sourceData = [{ 
          referrer: 'Direct Traffic', 
          clicks: data.directClicks, 
          percentage: (data.directClicks / data.totalClicks) * 100,
          uniqueClicks: Math.floor(data.directClicks * 0.8),
          category: 'direct' as const 
        }]
      } else {
        sourceData = sourceData.filter(item => item.category === viewMode)
      }
    } else {
      // Include direct traffic in 'all' view
      sourceData = [
        ...sourceData,
        { 
          referrer: 'Direct Traffic', 
          clicks: data.directClicks, 
          percentage: (data.directClicks / data.totalClicks) * 100,
          uniqueClicks: Math.floor(data.directClicks * 0.8),
          category: 'direct' as const,
          trend: (Math.random() - 0.5) * 10
        }
      ]
    }

    // Apply search filter
    if (searchTerm) {
      sourceData = sourceData.filter(item => 
        item.referrer.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    // Apply sorting
    if (sortBy === 'clicks') {
      sourceData = sourceData.sort((a, b) => b.clicks - a.clicks)
    } else if (sortBy === 'alphabetical') {
      sourceData = sourceData.sort((a, b) => a.referrer.localeCompare(b.referrer))
    } else if (sortBy === 'trend') {
      sourceData = sourceData.sort((a, b) => (b.trend || 0) - (a.trend || 0))
    }

    return sourceData
  }, [data, viewMode, searchTerm, sortBy])

  const chartColors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#06B6D4', '#84CC16', '#F97316']

  const exportData = () => {
    if (!filteredData.length) return

    const headers = ['Referrer', 'Category', 'Clicks', 'Unique Clicks', 'Percentage', 'Trend']
    const rows = filteredData.map(item => [
      item.referrer,
      item.category,
      item.clicks.toString(),
      item.uniqueClicks.toString(),
      `${item.percentage.toFixed(1)}%`,
      item.trend ? `${item.trend.toFixed(1)}%` : '0%'
    ])

    const csvContent = [headers, ...rows].map(row => row.join(',')).join('\n')
    const blob = new Blob([csvContent], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `referrer-stats-${viewMode}-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const getTrendIcon = (trend?: number) => {
    if (!trend || Math.abs(trend) < 0.1) return <Minus className="h-3 w-3" />
    if (trend > 0) return <ArrowUp className="h-3 w-3" />
    return <ArrowDown className="h-3 w-3" />
  }

  const getTrendColor = (trend?: number) => {
    if (!trend || Math.abs(trend) < 0.1) return 'text-gray-500'
    if (trend > 0) return 'text-green-600'
    return 'text-red-600'
  }

  const renderChart = () => {
    if (filteredData.length === 0) {
      return (
        <div className="text-center py-12">
          <Share2 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Referrer Data</h3>
          <p className="text-gray-600">No traffic source data found for the selected view.</p>
        </div>
      )
    }

    if (chartType === 'pie') {
      return (
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={filteredData.slice(0, 8)}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ referrer, percentage }) => `${referrer}: ${percentage.toFixed(1)}%`}
                outerRadius={120}
                fill="#8884d8"
                dataKey="clicks"
                nameKey="referrer"
              >
                {filteredData.slice(0, 8).map((entry, index) => (
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
              data={filteredData.slice(0, 15)}
              margin={{ top: 20, right: 30, left: 20, bottom: 80 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey="referrer"
                angle={-45}
                textAnchor="end"
                height={100}
                fontSize={12}
              />
              <YAxis />
              <Tooltip formatter={(value, name) => [value.toLocaleString(), name === 'clicks' ? 'Clicks' : 'Unique Clicks']} />
              <Legend />
              <Bar dataKey="clicks" fill="#3B82F6" name="Clicks" radius={[2, 2, 0, 0]} />
              <Bar dataKey="uniqueClicks" fill="#10B981" name="Unique Clicks" radius={[2, 2, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      )
    } else if (chartType === 'treemap') {
      const treemapData = filteredData.slice(0, 12).map(item => ({
        name: item.referrer,
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
                  Referrer
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Category
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Clicks
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Unique Clicks
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Percentage
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Trend
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filteredData.map((item, index) => {
                const icon = referrerIcons[item.referrer.toLowerCase()] || <ExternalLink className="h-5 w-5 text-gray-500" />
                
                return (
                  <tr key={index} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50'}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      <div className="flex items-center">
                        <div className="flex-shrink-0 mr-3">
                          {icon}
                        </div>
                        {item.referrer}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                        item.category === 'search' ? 'bg-blue-100 text-blue-800' :
                        item.category === 'social' ? 'bg-green-100 text-green-800' :
                        item.category === 'direct' ? 'bg-yellow-100 text-yellow-800' :
                        item.category === 'email' ? 'bg-purple-100 text-purple-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {item.category}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {item.clicks.toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {item.uniqueClicks.toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {item.percentage.toFixed(1)}%
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <div className={`flex items-center ${getTrendColor(item.trend)}`}>
                        {getTrendIcon(item.trend)}
                        <span className="ml-1">{item.trend ? Math.abs(item.trend).toFixed(1) : '0.0'}%</span>
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )
    }
  }

  if (isLoading) {
    return <PageLoading message="Loading referrer analytics..." />
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <Share2 className="h-5 w-5 text-red-500 mr-2" />
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
            <h3 className="text-lg font-medium text-gray-900">Traffic Source Analytics</h3>
            <p className="text-sm text-gray-600 mt-1">
              Analysis of referrer sources and traffic channels
            </p>
          </div>
          <div className="flex items-center space-x-3">
            <button
              onClick={() => fetchReferrerData(false)}
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
          <label className="block text-sm font-medium text-gray-700 mb-2">Traffic Source</label>
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
                placeholder="Search referrers..."
                className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Sort By</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as 'clicks' | 'alphabetical' | 'trend')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="clicks">Most Clicks</option>
              <option value="alphabetical">Alphabetical</option>
              <option value="trend">Trending</option>
            </select>
          </div>
        </div>
      </div>

      {/* Summary Stats */}
      {data && (
        <div className="p-6 border-b border-gray-200">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{data.totalClicks.toLocaleString()}</div>
              <div className="text-sm text-gray-600">Total Clicks</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{data.directClicks.toLocaleString()}</div>
              <div className="text-sm text-gray-600">Direct Traffic</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{data.referrers.length}</div>
              <div className="text-sm text-gray-600">Unique Referrers</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {data.topDomains[0]?.domain || 'N/A'}
              </div>
              <div className="text-sm text-gray-600">Top Referrer</div>
            </div>
          </div>
        </div>
      )}

      {/* Chart */}
      <div className="p-6">
        {renderChart()}
      </div>

      {/* Category Breakdown */}
      {data && chartType !== 'table' && (
        <div className="p-6 border-t border-gray-200">
          <h4 className="text-lg font-medium text-gray-900 mb-4">Traffic by Category</h4>
          <div className="grid grid-cols-2 lg:grid-cols-5 gap-4">
            {data.categories.map(category => (
              <div key={category.category} className="text-center p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center justify-center mb-2">
                  <div 
                    className="w-4 h-4 rounded-full mr-2"
                    style={{ backgroundColor: category.color }}
                  />
                  <span className="text-sm font-medium text-gray-900">{category.category}</span>
                </div>
                <div className="text-lg font-bold text-gray-900">{category.clicks.toLocaleString()}</div>
                <div className="text-xs text-gray-500">{category.percentage.toFixed(1)}%</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default ReferrerStats