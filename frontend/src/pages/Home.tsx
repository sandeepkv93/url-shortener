import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import URLShortener from '@/components/url/URLShortener'
import { ShortURL } from '@/types/url'
import {
  LinkIcon,
  BarChart3,
  QrCode,
  Shield,
  Zap,
  Globe,
  Users,
  TrendingUp
} from 'lucide-react'

const Home = () => {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [recentUrls, setRecentUrls] = useState<ShortURL[]>([])

  useEffect(() => {
    // If user is already logged in, redirect to dashboard
    if (user) {
      navigate('/dashboard')
    }
  }, [user, navigate])

  const handleURLCreated = (url: ShortURL) => {
    setRecentUrls(prev => [url, ...prev.slice(0, 2)])
    
    // If user is logged in, redirect to dashboard to see the new URL
    if (user) {
      navigate('/dashboard')
    }
  }

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <div className="relative overflow-hidden bg-gradient-to-br from-primary-50 to-primary-100">
        <div className="absolute inset-0 opacity-40 bg-dots-pattern"></div>
        <div className="relative px-4 py-16 sm:py-24">
          <div className="max-w-7xl mx-auto">
            <div className="text-center">
              <h1 className="text-4xl font-bold text-gray-900 sm:text-5xl md:text-6xl lg:text-7xl">
                <span className="block">Shorten Your URLs</span>
                <span className="block text-primary-600">Track Everything</span>
              </h1>
              <p className="mt-6 max-w-3xl mx-auto text-lg text-gray-600 sm:text-xl">
                Create short, memorable links with powerful analytics, QR code generation, 
                and comprehensive link management. Perfect for marketers, businesses, and developers.
              </p>
              
              {!user && (
                <div className="mt-8 flex flex-col sm:flex-row gap-4 justify-center">
                  <Link
                    to="/register"
                    className="inline-flex items-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
                  >
                    Get Started Free
                  </Link>
                  <Link
                    to="/demo"
                    className="inline-flex items-center px-8 py-3 border border-primary-600 text-base font-medium rounded-md text-primary-600 bg-white hover:bg-primary-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
                  >
                    View Demo
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* URL Shortener Section */}
      <div className="py-16 bg-white">
        <div className="max-w-4xl mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Try it now - no registration required
            </h2>
            <p className="text-lg text-gray-600">
              Create your first short link and see the magic happen
            </p>
          </div>
          
          <URLShortener 
            onURLCreated={handleURLCreated}
            className="max-w-2xl mx-auto"
          />
          
          {/* Recent URLs for anonymous users */}
          {!user && recentUrls.length > 0 && (
            <div className="mt-12 max-w-2xl mx-auto">
              <h3 className="text-lg font-medium text-gray-900 mb-4 text-center">
                Your Recent Links
              </h3>
              <div className="space-y-3">
                {recentUrls.map((url, index) => (
                  <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {window.location.origin}/{url.shortCode}
                      </p>
                      <p className="text-sm text-gray-500 truncate">
                        {url.originalUrl}
                      </p>
                    </div>
                    <button
                      onClick={() => navigator.clipboard.writeText(`${window.location.origin}/${url.shortCode}`)}
                      className="ml-4 p-2 text-gray-400 hover:text-gray-600 transition-colors"
                      title="Copy to clipboard"
                    >
                      <LinkIcon className="h-4 w-4" />
                    </button>
                  </div>
                ))}
              </div>
              <div className="mt-6 text-center">
                <p className="text-sm text-gray-600">
                  Want to manage your links? {' '}
                  <Link to="/register" className="font-medium text-primary-600 hover:text-primary-500">
                    Create a free account
                  </Link>
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Features Section */}
      <div className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 sm:text-4xl">
              Everything you need to manage links
            </h2>
            <p className="mt-4 text-lg text-gray-600 max-w-3xl mx-auto">
              From simple URL shortening to advanced analytics and QR codes, 
              we've got all the tools you need to succeed.
            </p>
          </div>

          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <div className="text-center p-6">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-xl bg-primary-100 text-primary-600 mb-6">
                <Zap className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Lightning Fast</h3>
              <p className="text-gray-600">
                Create short links instantly with our optimized infrastructure and global CDN.
              </p>
            </div>

            <div className="text-center p-6">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-xl bg-primary-100 text-primary-600 mb-6">
                <BarChart3 className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Advanced Analytics</h3>
              <p className="text-gray-600">
                Track clicks, geographic data, devices, referrers, and more with real-time insights.
              </p>
            </div>

            <div className="text-center p-6">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-xl bg-primary-100 text-primary-600 mb-6">
                <QrCode className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">QR Code Generation</h3>
              <p className="text-gray-600">
                Generate customizable QR codes for your links with multiple formats and styles.
              </p>
            </div>

            <div className="text-center p-6">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-xl bg-primary-100 text-primary-600 mb-6">
                <Shield className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Secure & Reliable</h3>
              <p className="text-gray-600">
                Enterprise-grade security with password protection, expiration dates, and access controls.
              </p>
            </div>

            <div className="text-center p-6">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-xl bg-primary-100 text-primary-600 mb-6">
                <Globe className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Global Performance</h3>
              <p className="text-gray-600">
                99.9% uptime with worldwide edge locations for the fastest redirect speeds.
              </p>
            </div>

            <div className="text-center p-6">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-xl bg-primary-100 text-primary-600 mb-6">
                <Users className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Team Collaboration</h3>
              <p className="text-gray-600">
                Share links with your team, manage permissions, and collaborate on campaigns.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Stats Section */}
      <div className="py-16 bg-primary-600">
        <div className="max-w-7xl mx-auto px-4">
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-3 text-center">
            <div>
              <div className="text-4xl font-bold text-white mb-2">10M+</div>
              <div className="text-primary-100">Links Created</div>
            </div>
            <div>
              <div className="text-4xl font-bold text-white mb-2">500K+</div>
              <div className="text-primary-100">Active Users</div>
            </div>
            <div>
              <div className="text-4xl font-bold text-white mb-2">99.9%</div>
              <div className="text-primary-100">Uptime</div>
            </div>
          </div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="py-16 bg-white">
        <div className="max-w-4xl mx-auto text-center px-4">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">
            Ready to get started?
          </h2>
          <p className="text-lg text-gray-600 mb-8">
            Join thousands of users who trust us with their links
          </p>
          {!user && (
            <Link
              to="/register"
              className="inline-flex items-center px-8 py-3 border border-transparent text-lg font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
            >
              Start Free Today
              <TrendingUp className="ml-2 h-5 w-5" />
            </Link>
          )}
        </div>
      </div>
    </div>
  )
}

export default Home