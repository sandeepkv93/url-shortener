import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { 
  Home, 
  Search, 
  ArrowLeft, 
  Mail, 
  BarChart3,
  User 
} from 'lucide-react'

const NotFound = () => {
  const { user } = useAuth()
  const navigate = useNavigate()

  const goBack = () => {
    navigate(-1)
  }

  const quickLinks = [
    {
      to: '/',
      label: 'Home',
      icon: <Home className="h-4 w-4" />,
      description: 'Return to the homepage'
    },
    ...(user ? [
      {
        to: '/dashboard',
        label: 'Dashboard',
        icon: <BarChart3 className="h-4 w-4" />,
        description: 'View your URL dashboard'
      },
      {
        to: '/analytics',
        label: 'Analytics',
        icon: <BarChart3 className="h-4 w-4" />,
        description: 'Check your URL analytics'
      },
      {
        to: '/profile',
        label: 'Profile',
        icon: <User className="h-4 w-4" />,
        description: 'Manage your account'
      }
    ] : [])
  ]

  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-gray-50">
      <div className="max-w-lg w-full text-center">
        {/* 404 Display */}
        <div className="mb-8">
          <div className="relative">
            <h1 className="text-9xl font-bold text-primary-600 opacity-20">404</h1>
            <div className="absolute inset-0 flex items-center justify-center">
              <Search className="h-16 w-16 text-primary-500" />
            </div>
          </div>
          <h2 className="text-3xl font-bold text-gray-900 mt-6">Page Not Found</h2>
          <p className="text-gray-600 mt-4 text-lg">
            Sorry, we couldn't find the page you're looking for. The link might be broken or the page may have been moved.
          </p>
        </div>
        
        {/* Action Buttons */}
        <div className="space-y-4 mb-8">
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <button
              onClick={goBack}
              className="inline-flex items-center justify-center px-6 py-3 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Go Back
            </button>
            
            <Link
              to="/"
              className="inline-flex items-center justify-center px-6 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
            >
              <Home className="mr-2 h-4 w-4" />
              Go Home
            </Link>
          </div>
        </div>

        {/* Quick Links */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Navigation</h3>
          <div className="grid grid-cols-1 gap-3">
            {quickLinks.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                className="flex items-center p-3 border border-gray-200 rounded-md hover:bg-gray-50 transition-colors text-left"
              >
                <div className="flex-shrink-0 w-8 h-8 bg-primary-100 rounded-md flex items-center justify-center text-primary-600 mr-3">
                  {link.icon}
                </div>
                <div className="flex-1">
                  <div className="font-medium text-gray-900">{link.label}</div>
                  <div className="text-sm text-gray-500">{link.description}</div>
                </div>
              </Link>
            ))}
          </div>
        </div>

        {/* Help Text */}
        <div className="mt-8 text-center">
          <p className="text-sm text-gray-500">
            Still having trouble? {' '}
            <a 
              href="mailto:support@urlshortener.com" 
              className="font-medium text-primary-600 hover:text-primary-500"
            >
              Contact Support
            </a>
          </p>
        </div>
      </div>
    </div>
  )
}

export default NotFound