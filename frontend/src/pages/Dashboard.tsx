import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import URLShortener from '@/components/url/URLShortener'
import URLList from '@/components/url/URLList'
import ClickChart from '@/components/analytics/ClickChart'
import { urlService } from '@/services/urls'
import { URL as URLType, URLListResponse } from '@/types/url'
import { PageLoading } from '@/components/common/Loading'
import {
  Link as LinkIcon,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Eye,
  Calendar,
  Plus,
  Settings,
  Users,
  Activity,
  ArrowRight,
  RefreshCw
} from 'lucide-react'

interface DashboardStats {
  totalURLs: number
  totalClicks: number
  clicksThisMonth: number
  activeURLs: number
  trend: {
    urls: number
    clicks: number
  }
}

const Dashboard = () => {
  const { user } = useAuth()
  const [isLoading, setIsLoading] = useState(true)
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [recentUrls, setRecentUrls] = useState<URLType[]>([])
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [showURLShortener, setShowURLShortener] = useState(false)

  // Fetch dashboard data
  const fetchDashboardData = async (showLoading = true) => {
    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)

    try {
      // Fetch recent URLs and overall stats
      const [urlsResponse] = await Promise.all([
        urlService.getUserURLs({ 
          page: 1, 
          limit: 6, 
          sortBy: 'createdAt', 
          sortOrder: 'desc' 
        })
      ])

      setRecentUrls(urlsResponse.urls)

      // Calculate stats from the fetched data
      const totalClicks = urlsResponse.urls.reduce((sum, url) => sum + url.clickCount, 0)
      const activeURLs = urlsResponse.urls.filter(url => url.isActive).length

      // For demo purposes, we'll calculate some mock trends and stats
      const mockStats: DashboardStats = {
        totalURLs: urlsResponse.total,
        totalClicks: totalClicks,
        clicksThisMonth: Math.floor(totalClicks * 0.3), // Assuming 30% of clicks are from this month
        activeURLs: activeURLs,
        trend: {
          urls: Math.floor(Math.random() * 20) - 10, // Random trend between -10 and +10
          clicks: Math.floor(Math.random() * 50) - 25 // Random trend between -25 and +25
        }
      }

      setStats(mockStats)

    } catch (error) {
      console.error('Failed to fetch dashboard data:', error)
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  useEffect(() => {
    if (user) {
      fetchDashboardData()
    }
  }, [user])

  const handleURLCreated = (newURL: URLType) => {
    setRecentUrls(prev => [newURL, ...prev.slice(0, 5)]) // Keep only 6 most recent
    if (stats) {
      setStats(prev => prev ? {
        ...prev,
        totalURLs: prev.totalURLs + 1,
        activeURLs: newURL.isActive ? prev.activeURLs + 1 : prev.activeURLs
      } : null)
    }
    setShowURLShortener(false)
  }

  const handleURLUpdate = (updatedURL: URLType) => {
    setRecentUrls(prev => prev.map(url => url.id === updatedURL.id ? updatedURL : url))
  }

  const handleURLDelete = (urlId: string) => {
    setRecentUrls(prev => prev.filter(url => url.id !== urlId))
    if (stats) {
      setStats(prev => prev ? {
        ...prev,
        totalURLs: Math.max(0, prev.totalURLs - 1)
      } : null)
    }
  }

  const formatNumber = (num: number) => {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + 'M'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K'
    }
    return num.toString()
  }

  const getTrendIcon = (trend: number) => {
    if (trend > 0) return <TrendingUp className="h-4 w-4 text-green-500" />
    if (trend < 0) return <TrendingDown className="h-4 w-4 text-red-500" />
    return <Activity className="h-4 w-4 text-gray-400" />
  }

  const getTrendColor = (trend: number) => {
    if (trend > 0) return 'text-green-600'
    if (trend < 0) return 'text-red-600'
    return 'text-gray-500'
  }

  if (!user) {
    return (
      <div className="px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Please log in to access your dashboard</h1>
          <Link to="/login" className="btn-primary">
            Log In
          </Link>
        </div>
      </div>
    )
  }

  if (isLoading) {
    return <PageLoading message="Loading your dashboard..." />
  }

  return (
    <div className="px-4 py-8 space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            Welcome back, {user.name || user.email}!
          </h1>
          <p className="mt-2 text-gray-600">
            Here's an overview of your URL shortening activity
          </p>
        </div>
        <div className="mt-4 sm:mt-0 flex items-center space-x-3">
          <button
            onClick={() => fetchDashboardData(false)}
            disabled={isRefreshing}
            className="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
            Refresh
          </button>
          <button
            onClick={() => setShowURLShortener(!showURLShortener)}
            className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700"
          >
            <Plus className="h-4 w-4 mr-2" />
            Create URL
          </button>
        </div>
      </div>

      {/* Stats Overview */}
      {stats && (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-primary-100 rounded-lg flex items-center justify-center">
                  <LinkIcon className="h-5 w-5 text-primary-600" />
                </div>
              </div>
              <div className="ml-4 flex-1">
                <div className="text-2xl font-bold text-gray-900">
                  {formatNumber(stats.totalURLs)}
                </div>
                <div className="text-sm text-gray-600">Total URLs</div>
                <div className="flex items-center mt-1">
                  {getTrendIcon(stats.trend.urls)}
                  <span className={`text-xs font-medium ml-1 ${getTrendColor(stats.trend.urls)}`}>
                    {Math.abs(stats.trend.urls)}% this month
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-green-100 rounded-lg flex items-center justify-center">
                  <Eye className="h-5 w-5 text-green-600" />
                </div>
              </div>
              <div className="ml-4 flex-1">
                <div className="text-2xl font-bold text-gray-900">
                  {formatNumber(stats.totalClicks)}
                </div>
                <div className="text-sm text-gray-600">Total Clicks</div>
                <div className="flex items-center mt-1">
                  {getTrendIcon(stats.trend.clicks)}
                  <span className={`text-xs font-medium ml-1 ${getTrendColor(stats.trend.clicks)}`}>
                    {Math.abs(stats.trend.clicks)}% this month
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-yellow-100 rounded-lg flex items-center justify-center">
                  <Calendar className="h-5 w-5 text-yellow-600" />
                </div>
              </div>
              <div className="ml-4 flex-1">
                <div className="text-2xl font-bold text-gray-900">
                  {formatNumber(stats.clicksThisMonth)}
                </div>
                <div className="text-sm text-gray-600">This Month</div>
                <div className="text-xs text-gray-500 mt-1">
                  {stats.totalClicks > 0 ? ((stats.clicksThisMonth / stats.totalClicks) * 100).toFixed(1) : 0}% of total
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
                  <Activity className="h-5 w-5 text-blue-600" />
                </div>
              </div>
              <div className="ml-4 flex-1">
                <div className="text-2xl font-bold text-gray-900">
                  {stats.activeURLs}
                </div>
                <div className="text-sm text-gray-600">Active URLs</div>
                <div className="text-xs text-gray-500 mt-1">
                  {stats.totalURLs > 0 ? ((stats.activeURLs / stats.totalURLs) * 100).toFixed(1) : 0}% of total
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* URL Shortener */}
      {showURLShortener && (
        <div className="bg-gray-50 rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-medium text-gray-900">Create New Short URL</h2>
            <button
              onClick={() => setShowURLShortener(false)}
              className="text-gray-400 hover:text-gray-600"
            >
              ✕
            </button>
          </div>
          <URLShortener 
            onURLCreated={handleURLCreated}
            className="max-w-none"
          />
        </div>
      )}

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
        {/* Recent URLs */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-medium text-gray-900">Recent URLs</h3>
                <p className="text-sm text-gray-600 mt-1">Your latest shortened links</p>
              </div>
              <Link
                to="/urls"
                className="inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-700"
              >
                View all
                <ArrowRight className="h-4 w-4 ml-1" />
              </Link>
            </div>
          </div>
          
          <div className="p-6">
            {recentUrls.length > 0 ? (
              <div className="space-y-4">
                {recentUrls.slice(0, 4).map(url => (
                  <div key={url.id} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center space-x-2">
                        <h4 className="text-sm font-medium text-gray-900 truncate">
                          {url.title || url.originalUrl}
                        </h4>
                        {!url.isActive && (
                          <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                            Inactive
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-gray-500 truncate">
                        {window.location.origin}/{url.shortCode}
                      </p>
                      <div className="flex items-center space-x-4 mt-1">
                        <span className="text-xs text-gray-500">
                          {url.clickCount} clicks
                        </span>
                        <span className="text-xs text-gray-500">
                          {new Date(url.createdAt).toLocaleDateString()}
                        </span>
                      </div>
                    </div>
                    <Link
                      to={`/analytics/${url.id}`}
                      className="ml-4 p-2 text-gray-400 hover:text-gray-600 transition-colors"
                      title="View analytics"
                    >
                      <BarChart3 className="h-4 w-4" />
                    </Link>
                  </div>
                ))}
                <div className="text-center pt-4">
                  <Link
                    to="/urls"
                    className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                  >
                    <Users className="h-4 w-4 mr-2" />
                    Manage All URLs
                  </Link>
                </div>
              </div>
            ) : (
              <div className="text-center py-8">
                <LinkIcon className="mx-auto h-12 w-12 text-gray-400" />
                <h3 className="mt-2 text-sm font-medium text-gray-900">No URLs yet</h3>
                <p className="mt-1 text-sm text-gray-500">Get started by creating your first short URL.</p>
                <div className="mt-6">
                  <button
                    onClick={() => setShowURLShortener(true)}
                    className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Create Short URL
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Analytics Overview */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-medium text-gray-900">Analytics Overview</h3>
                <p className="text-sm text-gray-600 mt-1">Click performance over the last 7 days</p>
              </div>
              <Link
                to="/analytics"
                className="inline-flex items-center text-sm font-medium text-primary-600 hover:text-primary-700"
              >
                View details
                <ArrowRight className="h-4 w-4 ml-1" />
              </Link>
            </div>
          </div>
          
          <div className="p-6">
            {stats && stats.totalClicks > 0 ? (
              <ClickChart 
                period="7d"
                className="h-64"
              />
            ) : (
              <div className="text-center py-12">
                <BarChart3 className="mx-auto h-12 w-12 text-gray-400" />
                <h3 className="mt-2 text-sm font-medium text-gray-900">No analytics data</h3>
                <p className="mt-1 text-sm text-gray-500">
                  Analytics will appear here once your URLs start getting clicks.
                </p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <Link
            to="/analytics"
            className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <BarChart3 className="h-8 w-8 text-primary-600 mr-3" />
            <div>
              <div className="font-medium text-gray-900">View Analytics</div>
              <div className="text-sm text-gray-500">Detailed performance insights</div>
            </div>
          </Link>
          
          <Link
            to="/profile"
            className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <Settings className="h-8 w-8 text-gray-600 mr-3" />
            <div>
              <div className="font-medium text-gray-900">Account Settings</div>
              <div className="text-sm text-gray-500">Manage your profile</div>
            </div>
          </Link>
          
          <button
            onClick={() => setShowURLShortener(true)}
            className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors text-left"
          >
            <Plus className="h-8 w-8 text-green-600 mr-3" />
            <div>
              <div className="font-medium text-gray-900">Create URL</div>
              <div className="text-sm text-gray-500">Shorten a new link</div>
            </div>
          </button>
        </div>
      </div>
    </div>
  )
}

export default Dashboard