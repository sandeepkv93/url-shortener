import { useState, useEffect, useMemo } from 'react'
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ReferenceLine
} from 'recharts'
import {
  TrendingUp,
  TrendingDown,
  BarChart3,
  Activity,
  Calendar,
  Clock,
  Users,
  MousePointer,
  Filter,
  Download,
  RefreshCw,
  ArrowUp,
  ArrowDown,
  Minus
} from 'lucide-react'
import { format, subDays, subHours, startOfDay, endOfDay, parseISO, isToday, isYesterday } from 'date-fns'
import { urlAnalyticsService } from '@/services/urls'
import { AnalyticsPeriod } from '@/types/url'
import { PageLoading } from '@/components/common/Loading'

interface ClickChartData {
  date: string
  clicks: number
  uniqueClicks: number
  bounceRate?: number
  avgSessionTime?: number
}

interface ClickChartProps {
  urlId?: string // If provided, shows analytics for specific URL
  period?: AnalyticsPeriod
  className?: string
}

type ChartType = 'area' | 'line' | 'bar'
type MetricType = 'clicks' | 'uniqueClicks' | 'both'

const ClickChart = ({ 
  urlId, 
  period: initialPeriod = '7d',
  className = '' 
}: ClickChartProps) => {
  const [data, setData] = useState<ClickChartData[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [selectedPeriod, setSelectedPeriod] = useState<AnalyticsPeriod>(initialPeriod)
  const [chartType, setChartType] = useState<ChartType>('area')
  const [metricType, setMetricType] = useState<MetricType>('both')
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

  const chartTypes: { value: ChartType; label: string; icon: React.ReactNode }[] = [
    { value: 'area', label: 'Area', icon: <Activity className="h-4 w-4" /> },
    { value: 'line', label: 'Line', icon: <TrendingUp className="h-4 w-4" /> },
    { value: 'bar', label: 'Bar', icon: <BarChart3 className="h-4 w-4" /> }
  ]

  const metricTypes: { value: MetricType; label: string; icon: React.ReactNode }[] = [
    { value: 'both', label: 'Both Metrics', icon: <MousePointer className="h-4 w-4" /> },
    { value: 'clicks', label: 'Total Clicks', icon: <BarChart3 className="h-4 w-4" /> },
    { value: 'uniqueClicks', label: 'Unique Clicks', icon: <Users className="h-4 w-4" /> }
  ]

  useEffect(() => {
    fetchChartData()
  }, [selectedPeriod, urlId])

  const fetchChartData = async (showLoading = true) => {
    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)
    setError(null)

    try {
      let chartData: ClickChartData[]

      if (urlId) {
        // Fetch data for specific URL
        const response = await urlAnalyticsService.getClickTimeline(urlId, selectedPeriod)
        chartData = response.timeline.map(item => ({
          ...item,
          bounceRate: Math.random() * 40 + 30, // Mock bounce rate
          avgSessionTime: Math.random() * 180 + 60 // Mock session time
        }))
      } else {
        // Generate mock aggregated data for dashboard
        chartData = generateMockChartData(selectedPeriod)
      }

      setData(chartData)
    } catch (err: any) {
      setError(err.message || 'Failed to load chart data')
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  const generateMockChartData = (period: AnalyticsPeriod): ClickChartData[] => {
    const now = new Date()
    const data: ClickChartData[] = []

    if (period === '1h') {
      // Generate hourly data for last hour (every 5 minutes)
      for (let i = 11; i >= 0; i--) {
        const date = format(subHours(now, i * 5 / 60), 'yyyy-MM-dd HH:mm')
        const clicks = Math.floor(Math.random() * 20) + 5
        const uniqueClicks = Math.floor(clicks * (0.7 + Math.random() * 0.2))
        
        data.push({
          date,
          clicks,
          uniqueClicks,
          bounceRate: Math.random() * 40 + 30,
          avgSessionTime: Math.random() * 180 + 60
        })
      }
    } else if (period === '24h') {
      // Generate hourly data for last 24 hours
      for (let i = 23; i >= 0; i--) {
        const date = format(subHours(now, i), 'yyyy-MM-dd HH:mm')
        const clicks = Math.floor(Math.random() * 50) + 10
        const uniqueClicks = Math.floor(clicks * (0.7 + Math.random() * 0.2))
        
        data.push({
          date,
          clicks,
          uniqueClicks,
          bounceRate: Math.random() * 40 + 30,
          avgSessionTime: Math.random() * 180 + 60
        })
      }
    } else {
      // Generate daily data
      const days = period === '7d' ? 7 : period === '30d' ? 30 : period === '90d' ? 90 : 365
      
      for (let i = days - 1; i >= 0; i--) {
        const date = format(subDays(now, i), 'yyyy-MM-dd')
        const baseClicks = period === '7d' ? 100 : period === '30d' ? 80 : 60
        const clicks = Math.floor(Math.random() * baseClicks) + 20
        const uniqueClicks = Math.floor(clicks * (0.7 + Math.random() * 0.2))
        
        data.push({
          date,
          clicks,
          uniqueClicks,
          bounceRate: Math.random() * 40 + 30,
          avgSessionTime: Math.random() * 180 + 60
        })
      }
    }

    return data
  }

  const chartMetrics = useMemo(() => {
    if (data.length === 0) return null

    const totalClicks = data.reduce((sum, item) => sum + item.clicks, 0)
    const totalUniqueClicks = data.reduce((sum, item) => sum + item.uniqueClicks, 0)
    const avgBounceRate = data.reduce((sum, item) => sum + (item.bounceRate || 0), 0) / data.length
    const avgSessionTime = data.reduce((sum, item) => sum + (item.avgSessionTime || 0), 0) / data.length

    // Calculate trends (compare last half vs first half)
    const midpoint = Math.floor(data.length / 2)
    const firstHalf = data.slice(0, midpoint)
    const secondHalf = data.slice(midpoint)

    const firstHalfClicks = firstHalf.reduce((sum, item) => sum + item.clicks, 0)
    const secondHalfClicks = secondHalf.reduce((sum, item) => sum + item.clicks, 0)

    const clicksTrend = firstHalfClicks > 0 
      ? ((secondHalfClicks - firstHalfClicks) / firstHalfClicks) * 100 
      : 0

    return {
      totalClicks,
      totalUniqueClicks,
      avgBounceRate,
      avgSessionTime,
      clicksTrend,
      conversionRate: (totalUniqueClicks / totalClicks) * 100
    }
  }, [data])

  const formatXAxisLabel = (dateStr: string) => {
    const date = parseISO(dateStr)
    
    if (selectedPeriod === '1h' || selectedPeriod === '24h') {
      return format(date, 'HH:mm')
    } else if (selectedPeriod === '7d') {
      return format(date, 'EEE')
    } else if (selectedPeriod === '30d' || selectedPeriod === '90d') {
      return format(date, 'MMM d')
    } else {
      return format(date, 'MMM yyyy')
    }
  }

  const formatTooltipLabel = (dateStr: string) => {
    const date = parseISO(dateStr)
    
    if (selectedPeriod === '1h' || selectedPeriod === '24h') {
      if (isToday(date)) {
        return `Today, ${format(date, 'h:mm a')}`
      } else if (isYesterday(date)) {
        return `Yesterday, ${format(date, 'h:mm a')}`
      } else {
        return format(date, 'MMM d, h:mm a')
      }
    } else {
      return format(date, 'MMMM d, yyyy')
    }
  }

  const exportData = () => {
    const csvContent = [
      ['Date', 'Total Clicks', 'Unique Clicks', 'Bounce Rate', 'Avg Session Time'],
      ...data.map(item => [
        item.date,
        item.clicks.toString(),
        item.uniqueClicks.toString(),
        (item.bounceRate || 0).toFixed(1),
        (item.avgSessionTime || 0).toFixed(0)
      ])
    ].map(row => row.join(',')).join('\n')

    const blob = new Blob([csvContent], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `click-chart-${selectedPeriod}-${format(new Date(), 'yyyy-MM-dd')}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const renderChart = () => {
    const chartProps = {
      data,
      margin: { top: 20, right: 30, left: 20, bottom: 5 }
    }

    const commonElements = (
      <>
        <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
        <XAxis 
          dataKey="date" 
          tickFormatter={formatXAxisLabel}
          stroke="#666"
          fontSize={12}
        />
        <YAxis stroke="#666" fontSize={12} />
        <Tooltip 
          labelFormatter={formatTooltipLabel}
          formatter={(value: number, name: string) => {
            if (name === 'clicks') return [value.toLocaleString(), 'Total Clicks']
            if (name === 'uniqueClicks') return [value.toLocaleString(), 'Unique Clicks']
            if (name === 'bounceRate') return [`${value.toFixed(1)}%`, 'Bounce Rate']
            if (name === 'avgSessionTime') return [`${Math.floor(value / 60)}:${(value % 60).toFixed(0).padStart(2, '0')}`, 'Avg Session Time']
            return [value, name]
          }}
          contentStyle={{
            backgroundColor: 'white',
            border: '1px solid #e5e7eb',
            borderRadius: '8px',
            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)'
          }}
        />
        <Legend />
      </>
    )

    if (chartType === 'area') {
      return (
        <AreaChart {...chartProps}>
          {commonElements}
          {(metricType === 'clicks' || metricType === 'both') && (
            <Area
              type="monotone"
              dataKey="clicks"
              stackId="1"
              stroke="#3B82F6"
              fill="#3B82F6"
              fillOpacity={0.6}
              name="Total Clicks"
            />
          )}
          {(metricType === 'uniqueClicks' || metricType === 'both') && (
            <Area
              type="monotone"
              dataKey="uniqueClicks"
              stackId="2"
              stroke="#10B981"
              fill="#10B981"
              fillOpacity={0.6}
              name="Unique Clicks"
            />
          )}
        </AreaChart>
      )
    } else if (chartType === 'line') {
      return (
        <LineChart {...chartProps}>
          {commonElements}
          {(metricType === 'clicks' || metricType === 'both') && (
            <Line
              type="monotone"
              dataKey="clicks"
              stroke="#3B82F6"
              strokeWidth={2}
              dot={{ fill: '#3B82F6', strokeWidth: 2, r: 4 }}
              name="Total Clicks"
            />
          )}
          {(metricType === 'uniqueClicks' || metricType === 'both') && (
            <Line
              type="monotone"
              dataKey="uniqueClicks"
              stroke="#10B981"
              strokeWidth={2}
              dot={{ fill: '#10B981', strokeWidth: 2, r: 4 }}
              name="Unique Clicks"
            />
          )}
        </LineChart>
      )
    } else {
      return (
        <BarChart {...chartProps}>
          {commonElements}
          {(metricType === 'clicks' || metricType === 'both') && (
            <Bar
              dataKey="clicks"
              fill="#3B82F6"
              name="Total Clicks"
              radius={[2, 2, 0, 0]}
            />
          )}
          {(metricType === 'uniqueClicks' || metricType === 'both') && (
            <Bar
              dataKey="uniqueClicks"
              fill="#10B981"
              name="Unique Clicks"
              radius={[2, 2, 0, 0]}
            />
          )}
        </BarChart>
      )
    }
  }

  if (isLoading) {
    return <PageLoading message="Loading chart data..." />
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <BarChart3 className="h-5 w-5 text-red-500 mr-2" />
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
            <h3 className="text-lg font-medium text-gray-900">Click Analytics</h3>
            <p className="text-sm text-gray-600 mt-1">
              Detailed click performance over time
            </p>
          </div>
          <div className="flex items-center space-x-3">
            <button
              onClick={() => fetchChartData(false)}
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
        {/* Period Selector */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Time Period</label>
          <div className="flex flex-wrap gap-2">
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

        {/* Chart Type and Metric Selectors */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Chart Type</label>
            <div className="flex gap-2">
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

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Metrics</label>
            <div className="flex gap-2">
              {metricTypes.map(metric => (
                <button
                  key={metric.value}
                  onClick={() => setMetricType(metric.value)}
                  className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                    metricType === metric.value
                      ? 'bg-primary-100 text-primary-700'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                >
                  {metric.icon}
                  <span className="ml-2 hidden sm:inline">{metric.label}</span>
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Metrics Summary */}
      {chartMetrics && (
        <div className="p-6 border-b border-gray-200">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{chartMetrics.totalClicks.toLocaleString()}</div>
              <div className="text-sm text-gray-600">Total Clicks</div>
              {chartMetrics.clicksTrend !== 0 && (
                <div className={`flex items-center justify-center mt-1 ${
                  chartMetrics.clicksTrend > 0 ? 'text-green-600' : 'text-red-600'
                }`}>
                  {chartMetrics.clicksTrend > 0 ? (
                    <ArrowUp className="h-3 w-3 mr-1" />
                  ) : chartMetrics.clicksTrend < 0 ? (
                    <ArrowDown className="h-3 w-3 mr-1" />
                  ) : (
                    <Minus className="h-3 w-3 mr-1" />
                  )}
                  <span className="text-xs font-medium">
                    {Math.abs(chartMetrics.clicksTrend).toFixed(1)}%
                  </span>
                </div>
              )}
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{chartMetrics.totalUniqueClicks.toLocaleString()}</div>
              <div className="text-sm text-gray-600">Unique Clicks</div>
              <div className="text-xs text-gray-500 mt-1">
                {chartMetrics.conversionRate.toFixed(1)}% conversion
              </div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">{chartMetrics.avgBounceRate.toFixed(1)}%</div>
              <div className="text-sm text-gray-600">Avg Bounce Rate</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {Math.floor(chartMetrics.avgSessionTime / 60)}:{(chartMetrics.avgSessionTime % 60).toFixed(0).padStart(2, '0')}
              </div>
              <div className="text-sm text-gray-600">Avg Session Time</div>
            </div>
          </div>
        </div>
      )}

      {/* Chart */}
      <div className="p-6">
        {data.length > 0 ? (
          <div className="h-96">
            <ResponsiveContainer width="100%" height="100%">
              {renderChart()}
            </ResponsiveContainer>
          </div>
        ) : (
          <div className="text-center py-12">
            <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No Data Available</h3>
            <p className="text-gray-600">No click data found for the selected period.</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default ClickChart