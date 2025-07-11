import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  ArrowLeft,
  BarChart3,
  Globe,
  Monitor,
  Users,
  Clock,
  TrendingUp,
  Calendar,
  Download,
  Share2,
  Edit3,
  Eye,
  EyeOff,
  MapPin,
  Smartphone,
  ExternalLink,
  Copy,
  QrCode,
  CheckCircle,
  AlertCircle
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
import { URL as URLType, AnalyticsPeriod } from '@/types/url'
import { urlService, urlAnalyticsService } from '@/services/urls'
import { PageLoading } from '@/components/common/Loading'

interface URLDetailsProps {
  className?: string
}

interface AnalyticsData {
  url: URLType
  timeline: Array<{
    date: string
    clicks: number
    uniqueClicks: number
  }>
  geographic: {
    countries: Array<{
      country: string
      countryCode: string
      clicks: number
      percentage: number
    }>
    cities: Array<{
      city: string
      country: string
      clicks: number
      percentage: number
    }>
  }
  devices: {
    devices: Array<{
      device: string
      clicks: number
      percentage: number
    }>
    browsers: Array<{
      browser: string
      clicks: number
      percentage: number
    }>
    operatingSystems: Array<{
      os: string
      clicks: number
      percentage: number
    }>
  }
  referrers: {
    referrers: Array<{
      referrer: string
      clicks: number
      percentage: number
    }>
    directClicks: number
    totalClicks: number
  }
  recentClicks: Array<{
    id: string
    timestamp: string
    country?: string
    city?: string
    device?: string
    browser?: string
    referer?: string
  }>
}

const URLDetails = ({ className = '' }: URLDetailsProps) => {
  const { id } = useParams<{ id: string }>()
  const [isLoading, setIsLoading] = useState(true)
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData | null>(null)
  const [selectedPeriod, setSelectedPeriod] = useState<AnalyticsPeriod>('30d')
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

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
    if (id) {
      fetchAnalyticsData(id, selectedPeriod)
    }
  }, [id, selectedPeriod])

  const fetchAnalyticsData = async (urlId: string, period: AnalyticsPeriod) => {
    setIsLoading(true)
    setError(null)

    try {
      const [
        url,
        timeline,
        geographic,
        devices,
        referrers,
        clickHistory
      ] = await Promise.all([
        urlService.getURL(urlId),
        urlAnalyticsService.getClickTimeline(urlId, period),
        urlAnalyticsService.getGeographicStats(urlId),
        urlAnalyticsService.getDeviceStats(urlId),
        urlAnalyticsService.getReferrerStats(urlId),
        urlService.getClickHistory(urlId, 1, 10)
      ])

      setAnalyticsData({
        url,
        timeline: timeline.timeline,
        geographic,
        devices,
        referrers,
        recentClicks: clickHistory.clicks
      })
    } catch (err: any) {
      setError(err.message || 'Failed to load analytics data')
    } finally {
      setIsLoading(false)
    }
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  const generateQRCode = (url: string) => {
    const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodeURIComponent(url)}`
    window.open(qrCodeUrl, '_blank')
  }

  const exportData = () => {
    if (!analyticsData) return
    // This would implement data export functionality
    console.log('Exporting analytics data:', analyticsData)
  }

  if (isLoading) {
    return <PageLoading message="Loading analytics..." />
  }

  if (error) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <AlertCircle className="h-5 w-5 text-red-500 mr-2" />
            <p className="text-sm text-red-700">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  if (!analyticsData) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="text-center py-12">
          <h3 className="text-lg font-medium text-gray-900 mb-2">URL not found</h3>
          <p className="text-gray-600">The requested URL could not be found.</p>
        </div>
      </div>
    )
  }

  const { url, timeline, geographic, devices, referrers, recentClicks } = analyticsData
  const shortUrl = `${window.location.origin}/${url.shortCode}`

  const chartColors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#06B6D4', '#84CC16', '#F97316']

  const totalClicks = timeline.reduce((sum, item) => sum + item.clicks, 0)
  const totalUniqueClicks = timeline.reduce((sum, item) => sum + item.uniqueClicks, 0)

  return (
    <div className={`max-w-7xl mx-auto ${className}`}>
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <Link
            to="/dashboard"
            className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
          >
            <ArrowLeft className="h-4 w-4 mr-1" />
            Back to Dashboard
          </Link>
          <div className="flex items-center space-x-3">
            <button
              onClick={exportData}
              className="inline-flex items-center px-3 py-1.5 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
            >
              <Download className="h-4 w-4 mr-1" />
              Export
            </button>
            <Link
              to={`/urls/${url.id}/edit`}
              className="inline-flex items-center px-3 py-1.5 border border-transparent rounded-md text-sm font-medium text-white bg-primary-600 hover:bg-primary-700"
            >
              <Edit3 className="h-4 w-4 mr-1" />
              Edit URL
            </Link>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-start justify-between">
            <div className="flex-1 min-w-0">
              <h1 className="text-2xl font-bold text-gray-900 truncate">
                {url.title || url.originalUrl}
              </h1>
              {url.title && (
                <p className="text-sm text-gray-500 truncate mt-1">
                  {url.originalUrl}
                </p>
              )}
              {url.description && (
                <p className="text-gray-600 mt-2">{url.description}</p>
              )}
            </div>
            <div className="flex items-center space-x-2 ml-4">
              <div className={`flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium ${
                url.isActive ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
              }`}>
                {url.isActive ? <Eye className="h-3 w-3" /> : <EyeOff className="h-3 w-3" />}
                <span>{url.isActive ? 'Active' : 'Inactive'}</span>
              </div>
              {url.isPublic ? (
                <div className="flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  <Globe className="h-3 w-3" />
                  <span>Public</span>
                </div>
              ) : (
                <div className="flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                  <Users className="h-3 w-3" />
                  <span>Private</span>
                </div>
              )}
            </div>
          </div>

          <div className="flex items-center space-x-2 mt-4">
            <div className="flex-1 p-3 bg-gray-50 border border-gray-200 rounded-md text-sm font-mono text-gray-800 truncate">
              {shortUrl}
            </div>
            <button
              onClick={() => copyToClipboard(shortUrl)}
              className={`p-3 rounded-md transition-colors ${
                copied 
                  ? 'bg-green-100 text-green-600' 
                  : 'bg-gray-100 hover:bg-gray-200 text-gray-600'
              }`}
              title="Copy short URL"
            >
              {copied ? <CheckCircle className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
            </button>
            <button
              onClick={() => generateQRCode(shortUrl)}
              className="p-3 bg-gray-100 hover:bg-gray-200 text-gray-600 rounded-md transition-colors"
              title="Generate QR Code"
            >
              <QrCode className="h-4 w-4" />
            </button>
            <a
              href={shortUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="p-3 bg-gray-100 hover:bg-gray-200 text-gray-600 rounded-md transition-colors"
              title="Open short URL"
            >
              <ExternalLink className="h-4 w-4" />
            </a>
          </div>

          {url.tags && url.tags.length > 0 && (
            <div className="flex items-center space-x-2 mt-4">
              <span className="text-sm text-gray-500">Tags:</span>
              <div className="flex flex-wrap gap-1">
                {url.tags.map(tag => (
                  <span key={tag} className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <BarChart3 className="h-8 w-8 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-2xl font-semibold text-gray-900">{url.clickCount}</p>
              <p className="text-sm text-gray-600">Total Clicks</p>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Users className="h-8 w-8 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-2xl font-semibold text-gray-900">{totalUniqueClicks}</p>
              <p className="text-sm text-gray-600">Unique Clicks</p>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <TrendingUp className="h-8 w-8 text-purple-600" />
            </div>
            <div className="ml-4">
              <p className="text-2xl font-semibold text-gray-900">
                {totalClicks > 0 ? ((totalUniqueClicks / totalClicks) * 100).toFixed(1) : 0}%
              </p>
              <p className="text-sm text-gray-600">Click-through Rate</p>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Calendar className="h-8 w-8 text-orange-600" />
            </div>
            <div className="ml-4">
              <p className="text-2xl font-semibold text-gray-900">
                {format(new Date(url.createdAt), 'MMM d')}
              </p>
              <p className="text-sm text-gray-600">Created</p>
            </div>
          </div>
        </div>
      </div>

      {/* Period Selector */}
      <div className="mb-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">Analytics Period</h3>
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

      {/* Click Timeline Chart */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Click Timeline</h3>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={timeline}>
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

      {/* Analytics Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Geographic Stats */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
            <MapPin className="h-5 w-5 mr-2" />
            Geographic Distribution
          </h3>
          <div className="space-y-4">
            <div>
              <h4 className="text-sm font-medium text-gray-700 mb-2">Top Countries</h4>
              <div className="space-y-2">
                {geographic.countries.slice(0, 5).map((country, index) => (
                  <div key={country.country} className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div 
                        className="w-3 h-3 rounded-full mr-2"
                        style={{ backgroundColor: chartColors[index % chartColors.length] }}
                      />
                      <span className="text-sm text-gray-900">{country.country}</span>
                    </div>
                    <div className="text-sm text-gray-600">
                      {country.clicks} ({country.percentage.toFixed(1)}%)
                    </div>
                  </div>
                ))}
              </div>
            </div>
            {geographic.cities.length > 0 && (
              <div>
                <h4 className="text-sm font-medium text-gray-700 mb-2">Top Cities</h4>
                <div className="space-y-2">
                  {geographic.cities.slice(0, 3).map((city, index) => (
                    <div key={`${city.city}-${city.country}`} className="flex items-center justify-between">
                      <span className="text-sm text-gray-900">{city.city}, {city.country}</span>
                      <span className="text-sm text-gray-600">{city.clicks}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Device Stats */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
            <Smartphone className="h-5 w-5 mr-2" />
            Device & Browser Stats
          </h3>
          <div className="grid grid-cols-1 gap-4">
            <div>
              <h4 className="text-sm font-medium text-gray-700 mb-2">Devices</h4>
              <div className="space-y-2">
                {devices.devices.slice(0, 3).map((device, index) => (
                  <div key={device.device} className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div 
                        className="w-3 h-3 rounded-full mr-2"
                        style={{ backgroundColor: chartColors[index % chartColors.length] }}
                      />
                      <span className="text-sm text-gray-900">{device.device}</span>
                    </div>
                    <div className="text-sm text-gray-600">
                      {device.clicks} ({device.percentage.toFixed(1)}%)
                    </div>
                  </div>
                ))}
              </div>
            </div>
            <div>
              <h4 className="text-sm font-medium text-gray-700 mb-2">Browsers</h4>
              <div className="space-y-2">
                {devices.browsers.slice(0, 3).map((browser, index) => (
                  <div key={browser.browser} className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div 
                        className="w-3 h-3 rounded-full mr-2"
                        style={{ backgroundColor: chartColors[(index + 3) % chartColors.length] }}
                      />
                      <span className="text-sm text-gray-900">{browser.browser}</span>
                    </div>
                    <div className="text-sm text-gray-600">
                      {browser.clicks} ({browser.percentage.toFixed(1)}%)
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Referrers and Recent Clicks */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Referrer Stats */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
            <Share2 className="h-5 w-5 mr-2" />
            Traffic Sources
          </h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-900">Direct Traffic</span>
              <span className="text-sm text-gray-600">
                {referrers.directClicks} ({((referrers.directClicks / referrers.totalClicks) * 100).toFixed(1)}%)
              </span>
            </div>
            {referrers.referrers.slice(0, 5).map((referrer, index) => (
              <div key={referrer.referrer} className="flex items-center justify-between">
                <span className="text-sm text-gray-900 truncate">{referrer.referrer}</span>
                <span className="text-sm text-gray-600">
                  {referrer.clicks} ({referrer.percentage.toFixed(1)}%)
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Clicks */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
            <Clock className="h-5 w-5 mr-2" />
            Recent Activity
          </h3>
          <div className="space-y-3">
            {recentClicks.length > 0 ? (
              recentClicks.map((click) => (
                <div key={click.id} className="flex items-center justify-between text-sm">
                  <div>
                    <div className="text-gray-900">
                      {click.city && click.country ? `${click.city}, ${click.country}` : click.country || 'Unknown Location'}
                    </div>
                    <div className="text-gray-500 text-xs">
                      {click.device && click.browser ? `${click.device} • ${click.browser}` : click.device || click.browser || 'Unknown Device'}
                    </div>
                  </div>
                  <div className="text-gray-500 text-xs">
                    {format(new Date(click.timestamp), 'MMM d, h:mm a')}
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-gray-500">No recent activity</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default URLDetails