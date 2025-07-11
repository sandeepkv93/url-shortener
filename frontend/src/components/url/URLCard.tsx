import { useState } from 'react'
import { Link } from 'react-router-dom'
import { 
  ExternalLink, 
  Copy, 
  QrCode, 
  BarChart3, 
  Edit3, 
  Trash2, 
  Eye, 
  EyeOff,
  Lock,
  Calendar,
  Globe,
  Users,
  Tag,
  MoreVertical,
  CheckCircle
} from 'lucide-react'
import { URL as URLType } from '@/types/url'
import { urlService } from '@/services/urls'
import { format } from 'date-fns'

interface URLCardProps {
  url: URLType
  onUpdate?: (url: URLType) => void
  onDelete?: (urlId: string) => void
  className?: string
  showActions?: boolean
}

const URLCard = ({ 
  url, 
  onUpdate, 
  onDelete, 
  className = '',
  showActions = true 
}: URLCardProps) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [copied, setCopied] = useState(false)
  const [isToggling, setIsToggling] = useState(false)

  const shortUrl = `${window.location.origin}/${url.shortCode}`

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  const generateQRCode = () => {
    const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodeURIComponent(shortUrl)}`
    window.open(qrCodeUrl, '_blank')
  }

  const toggleURLStatus = async () => {
    if (!onUpdate) return
    
    setIsToggling(true)
    try {
      const updatedURL = await urlService.toggleURLStatus(url.id)
      onUpdate(updatedURL)
    } catch (error) {
      console.error('Failed to toggle URL status:', error)
    } finally {
      setIsToggling(false)
    }
  }

  const handleDelete = async () => {
    if (!onDelete || !confirm('Are you sure you want to delete this URL? This action cannot be undone.')) {
      return
    }

    try {
      await urlService.deleteURL(url.id)
      onDelete(url.id)
    } catch (error) {
      console.error('Failed to delete URL:', error)
    }
  }

  const getStatusColor = () => {
    if (!url.isActive) return 'text-gray-500'
    if (url.expiresAt && new Date(url.expiresAt) < new Date()) return 'text-orange-500'
    return 'text-green-500'
  }

  const getStatusText = () => {
    if (!url.isActive) return 'Inactive'
    if (url.expiresAt && new Date(url.expiresAt) < new Date()) return 'Expired'
    return 'Active'
  }

  const isExpired = url.expiresAt && new Date(url.expiresAt) < new Date()

  return (
    <div className={`bg-white rounded-lg shadow-md border border-gray-200 hover:shadow-lg transition-shadow ${className}`}>
      <div className="p-6">
        {/* Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1 min-w-0">
            {/* Title or URL */}
            <h3 className="text-lg font-semibold text-gray-900 truncate">
              {url.title || url.originalUrl}
            </h3>
            {url.title && (
              <p className="text-sm text-gray-500 truncate mt-1">
                {url.originalUrl}
              </p>
            )}
          </div>
          
          {showActions && (
            <div className="relative ml-4">
              <button
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                className="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100"
              >
                <MoreVertical className="h-5 w-5" />
              </button>
              
              {isMenuOpen && (
                <>
                  <div 
                    className="fixed inset-0 z-10" 
                    onClick={() => setIsMenuOpen(false)}
                  />
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-20 border border-gray-200">
                    <div className="py-1">
                      <Link
                        to={`/urls/${url.id}/edit`}
                        className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        onClick={() => setIsMenuOpen(false)}
                      >
                        <Edit3 className="h-4 w-4 mr-2" />
                        Edit URL
                      </Link>
                      <button
                        onClick={() => {
                          toggleURLStatus()
                          setIsMenuOpen(false)
                        }}
                        disabled={isToggling}
                        className="w-full flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 disabled:opacity-50"
                      >
                        {url.isActive ? (
                          <>
                            <EyeOff className="h-4 w-4 mr-2" />
                            Deactivate
                          </>
                        ) : (
                          <>
                            <Eye className="h-4 w-4 mr-2" />
                            Activate
                          </>
                        )}
                      </button>
                      <button
                        onClick={() => {
                          handleDelete()
                          setIsMenuOpen(false)
                        }}
                        className="w-full flex items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50"
                      >
                        <Trash2 className="h-4 w-4 mr-2" />
                        Delete
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          )}
        </div>

        {/* Short URL */}
        <div className="mb-4">
          <div className="flex items-center space-x-2">
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
              onClick={generateQRCode}
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
          {copied && (
            <p className="text-sm text-green-600 mt-1">✓ Copied to clipboard!</p>
          )}
        </div>

        {/* Description */}
        {url.description && (
          <p className="text-sm text-gray-600 mb-4 line-clamp-2">
            {url.description}
          </p>
        )}

        {/* Tags */}
        {url.tags && url.tags.length > 0 && (
          <div className="flex items-center space-x-2 mb-4">
            <Tag className="h-4 w-4 text-gray-400" />
            <div className="flex flex-wrap gap-1">
              {url.tags.slice(0, 3).map(tag => (
                <span key={tag} className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {tag}
                </span>
              ))}
              {url.tags.length > 3 && (
                <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-600">
                  +{url.tags.length - 3} more
                </span>
              )}
            </div>
          </div>
        )}

        {/* Metadata */}
        <div className="grid grid-cols-2 gap-4 text-sm text-gray-600 mb-4">
          <div className="flex items-center space-x-2">
            <BarChart3 className="h-4 w-4" />
            <span>{url.clickCount} clicks</span>
          </div>
          <div className="flex items-center space-x-2">
            <Calendar className="h-4 w-4" />
            <span>Created {format(new Date(url.createdAt), 'MMM d, yyyy')}</span>
          </div>
        </div>

        {/* Status Indicators */}
        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center space-x-4">
            {/* Active Status */}
            <div className="flex items-center space-x-1">
              <div className={`h-2 w-2 rounded-full ${
                url.isActive 
                  ? isExpired ? 'bg-orange-500' : 'bg-green-500'
                  : 'bg-gray-400'
              }`} />
              <span className={getStatusColor()}>
                {getStatusText()}
              </span>
            </div>

            {/* Public/Private */}
            <div className="flex items-center space-x-1">
              {url.isPublic ? (
                <>
                  <Globe className="h-4 w-4 text-green-600" />
                  <span className="text-green-600">Public</span>
                </>
              ) : (
                <>
                  <Users className="h-4 w-4 text-gray-600" />
                  <span className="text-gray-600">Private</span>
                </>
              )}
            </div>

            {/* Password Protected */}
            {url.password && (
              <div className="flex items-center space-x-1">
                <Lock className="h-4 w-4 text-amber-600" />
                <span className="text-amber-600">Protected</span>
              </div>
            )}
          </div>

          {/* Analytics Link */}
          <Link
            to={`/analytics/${url.id}`}
            className="flex items-center space-x-1 text-primary-600 hover:text-primary-700 font-medium"
          >
            <BarChart3 className="h-4 w-4" />
            <span>Analytics</span>
          </Link>
        </div>

        {/* Expiration Warning */}
        {url.expiresAt && (
          <div className={`mt-3 p-2 rounded-md text-sm ${
            isExpired 
              ? 'bg-red-50 text-red-700 border border-red-200'
              : new Date(url.expiresAt).getTime() - Date.now() < 7 * 24 * 60 * 60 * 1000
              ? 'bg-yellow-50 text-yellow-700 border border-yellow-200'
              : 'bg-blue-50 text-blue-700 border border-blue-200'
          }`}>
            <div className="flex items-center space-x-2">
              <Calendar className="h-4 w-4" />
              <span>
                {isExpired 
                  ? `Expired on ${format(new Date(url.expiresAt), 'MMM d, yyyy h:mm a')}`
                  : `Expires on ${format(new Date(url.expiresAt), 'MMM d, yyyy h:mm a')}`
                }
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default URLCard