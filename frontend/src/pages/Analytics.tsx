import { useState, useEffect } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { useSearchParams, Link } from 'react-router-dom'
import ClickChart from '@/components/analytics/ClickChart'
import AnalyticsDashboard from '@/components/analytics/Dashboard'
import GeographicMap from '@/components/analytics/GeographicMap'
import DeviceStats from '@/components/analytics/DeviceStats'
import ReferrerStats from '@/components/analytics/ReferrerStats'
import { urlService } from '@/services/urls'
import { URL as URLType } from '@/types/url'
import { PageLoading } from '@/components/common/Loading'
import {
  BarChart3,
  Globe,
  Smartphone,
  ExternalLink,
  Calendar,
  TrendingUp,
  Users,
  MousePointer,
  Filter,
  Download,
  RefreshCw,
  Settings,
  ArrowLeft,
  Search
} from 'lucide-react'

type AnalyticsView = 'overview' | 'timeline' | 'geographic' | 'devices' | 'referrers'

const Analytics = () => {
  const { user } = useAuth()
  const [searchParams, setSearchParams] = useSearchParams()
  const [selectedView, setSelectedView] = useState<AnalyticsView>('overview')
  const [selectedURL, setSelectedURL] = useState<URLType | null>(null)
  const [userURLs, setUserURLs] = useState<URLType[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [urlSearchTerm, setUrlSearchTerm] = useState('')
  
  const urlId = searchParams.get('url')
  const view = (searchParams.get('view') as AnalyticsView) || 'overview'

  useEffect(() => {
    setSelectedView(view)
  }, [view])

  useEffect(() => {
    if (user) {
      fetchUserURLs()
    }
  }, [user])

  useEffect(() => {
    if (urlId && userURLs.length > 0) {
      const url = userURLs.find(u => u.id === urlId)
      setSelectedURL(url || null)
    }
  }, [urlId, userURLs])

  const fetchUserURLs = async () => {
    setIsLoading(true)
    try {
      const response = await urlService.getUserURLs({
        page: 1,
        limit: 100, // Get all URLs for dropdown
        sortBy: 'createdAt',
        sortOrder: 'desc'
      })
      setUserURLs(response.urls)
    } catch (error) {
      console.error('Failed to fetch URLs:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleViewChange = (newView: AnalyticsView) => {
    const params = new URLSearchParams(searchParams)
    params.set('view', newView)
    setSearchParams(params)
  }

  const handleURLSelect = (url: URLType | null) => {
    const params = new URLSearchParams(searchParams)
    if (url) {
      params.set('url', url.id)
    } else {
      params.delete('url')
    }
    setSearchParams(params)
  }

  const exportData = () => {
    // Implementation for exporting analytics data
    console.log('Exporting analytics data for:', selectedURL?.id || 'all URLs')
  }

  const filteredURLs = userURLs.filter(url => 
    url.title?.toLowerCase().includes(urlSearchTerm.toLowerCase()) ||
    url.originalUrl.toLowerCase().includes(urlSearchTerm.toLowerCase()) ||
    url.shortCode.toLowerCase().includes(urlSearchTerm.toLowerCase())
  )

  const analyticsViews: { value: AnalyticsView; label: string; icon: React.ReactNode; description: string }[] = [
    { 
      value: 'overview', 
      label: 'Overview', 
      icon: <BarChart3 className="h-5 w-5" />,
      description: 'Complete analytics dashboard'
    },
    { 
      value: 'timeline', 
      label: 'Timeline', 
      icon: <TrendingUp className="h-5 w-5" />,
      description: 'Click performance over time'
    },
    { 
      value: 'geographic', 
      label: 'Geographic', 
      icon: <Globe className="h-5 w-5" />,
      description: 'Geographic distribution'
    },
    { 
      value: 'devices', 
      label: 'Devices', 
      icon: <Smartphone className="h-5 w-5" />,
      description: 'Device and browser statistics'
    },
    { 
      value: 'referrers', 
      label: 'Referrers', 
      icon: <ExternalLink className="h-5 w-5" />,
      description: 'Traffic sources and referrers'
    }
  ]

  if (!user) {
    return (
      <div className="px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Please log in to view analytics</h1>
          <Link to="/login" className="btn-primary">
            Log In
          </Link>
        </div>
      </div>
    )
  }

  if (isLoading) {
    return <PageLoading message="Loading analytics..." />
  }

  const renderAnalyticsContent = () => {
    switch (selectedView) {
      case 'overview':
        return (
          <AnalyticsDashboard 
            className="px-4 py-6"
          />
        )
      case 'timeline':
        return (
          <div className="px-4 py-6">
            <ClickChart 
              urlId={selectedURL?.id}
              period="30d"
              className="mb-6"
            />
          </div>
        )
      case 'geographic':
        return (
          <div className="px-4 py-6">
            <GeographicMap 
              urlId={selectedURL?.id}
              className="mb-6"
            />
          </div>
        )
      case 'devices':
        return (
          <div className="px-4 py-6">
            <DeviceStats 
              urlId={selectedURL?.id}
              className="mb-6"
            />
          </div>
        )
      case 'referrers':
        return (
          <div className="px-4 py-6">
            <ReferrerStats 
              urlId={selectedURL?.id}
              className="mb-6"
            />
          </div>
        )
      default:
        return null
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="px-4 py-6">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
            <div>
              <div className="flex items-center space-x-3">
                <Link 
                  to="/dashboard"
                  className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
                >
                  <ArrowLeft className="h-5 w-5" />
                </Link>
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">Analytics</h1>
                  <p className="mt-1 text-gray-600">
                    {selectedURL 
                      ? `Detailed insights for ${selectedURL.title || selectedURL.shortCode}`
                      : 'Comprehensive analytics for all your URLs'
                    }
                  </p>
                </div>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <button
                onClick={exportData}
                className="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                <Download className="h-4 w-4 mr-2" />
                Export
              </button>
              <Link
                to="/dashboard"
                className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700"
              >
                <Settings className="h-4 w-4 mr-2" />
                Dashboard
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* URL Selector and View Navigation */}
      <div className="bg-white border-b border-gray-200">
        <div className="px-4 py-4">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
            {/* URL Selector */}
            <div className="flex-1 max-w-md">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Analyze Specific URL (Optional)
              </label>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  type="text"
                  value={urlSearchTerm}
                  onChange={(e) => setUrlSearchTerm(e.target.value)}
                  placeholder="Search URLs..."
                  className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
              </div>
              {urlSearchTerm && (
                <div className="absolute z-10 mt-1 w-full bg-white border border-gray-200 rounded-md shadow-lg max-h-60 overflow-auto">
                  <div 
                    className="px-3 py-2 text-sm text-gray-900 hover:bg-gray-50 cursor-pointer border-b border-gray-100"
                    onClick={() => {
                      handleURLSelect(null)
                      setUrlSearchTerm('')
                    }}
                  >
                    <div className="font-medium">All URLs</div>
                    <div className="text-gray-500">View analytics for all your URLs</div>
                  </div>
                  {filteredURLs.map(url => (
                    <div
                      key={url.id}
                      className="px-3 py-2 text-sm text-gray-900 hover:bg-gray-50 cursor-pointer"
                      onClick={() => {
                        handleURLSelect(url)
                        setUrlSearchTerm('')
                      }}
                    >
                      <div className="font-medium truncate">
                        {url.title || url.originalUrl}
                      </div>
                      <div className="text-gray-500 text-xs">
                        /{url.shortCode} • {url.clickCount} clicks
                      </div>
                    </div>
                  ))}
                </div>
              )}
              
              {selectedURL && (
                <div className="mt-2 p-3 bg-primary-50 border border-primary-200 rounded-md">
                  <div className="flex items-center justify-between">
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-medium text-primary-900 truncate">
                        {selectedURL.title || selectedURL.originalUrl}
                      </div>
                      <div className="text-xs text-primary-700">
                        /{selectedURL.shortCode} • {selectedURL.clickCount} total clicks
                      </div>
                    </div>
                    <button
                      onClick={() => handleURLSelect(null)}
                      className="ml-2 text-primary-600 hover:text-primary-800"
                    >
                      ✕
                    </button>
                  </div>
                </div>
              )}
            </div>

            {/* View Navigation */}
            <div className="flex items-center space-x-1 bg-gray-100 p-1 rounded-lg">
              {analyticsViews.map(view => (
                <button
                  key={view.value}
                  onClick={() => handleViewChange(view.value)}
                  className={`flex items-center space-x-2 px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                    selectedView === view.value
                      ? 'bg-white text-primary-700 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-200'
                  }`}
                  title={view.description}
                >
                  {view.icon}
                  <span className="hidden sm:inline">{view.label}</span>
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Analytics Content */}
      <div className="flex-1">
        {userURLs.length === 0 ? (
          <div className="flex items-center justify-center min-h-96">
            <div className="text-center">
              <BarChart3 className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-lg font-medium text-gray-900">No URLs to analyze</h3>
              <p className="mt-1 text-gray-500 mb-6">
                Create some URLs to see detailed analytics data.
              </p>
              <Link
                to="/dashboard"
                className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700"
              >
                Create Your First URL
              </Link>
            </div>
          </div>
        ) : (
          renderAnalyticsContent()
        )}
      </div>

      {/* Analytics Summary Footer */}
      {userURLs.length > 0 && (
        <div className="bg-white border-t border-gray-200 px-4 py-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {userURLs.length}
              </div>
              <div className="text-sm text-gray-600">Total URLs</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {userURLs.reduce((sum, url) => sum + url.clickCount, 0).toLocaleString()}
              </div>
              <div className="text-sm text-gray-600">Total Clicks</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {userURLs.filter(url => url.isActive).length}
              </div>
              <div className="text-sm text-gray-600">Active URLs</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900">
                {userURLs.length > 0 ? 
                  (userURLs.reduce((sum, url) => sum + url.clickCount, 0) / userURLs.length).toFixed(1) : 
                  '0'
                }
              </div>
              <div className="text-sm text-gray-600">Avg Clicks/URL</div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Analytics