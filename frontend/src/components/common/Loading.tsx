import { Loader2 } from 'lucide-react'

interface LoadingProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  variant?: 'spinner' | 'dots' | 'pulse' | 'skeleton'
  message?: string
  fullScreen?: boolean
  overlay?: boolean
  className?: string
}

const Loading = ({
  size = 'md',
  variant = 'spinner',
  message,
  fullScreen = false,
  overlay = false,
  className = ''
}: LoadingProps) => {
  const sizeClasses = {
    sm: 'h-4 w-4',
    md: 'h-6 w-6',
    lg: 'h-8 w-8',
    xl: 'h-12 w-12'
  }

  const renderSpinner = () => (
    <Loader2 className={`${sizeClasses[size]} animate-spin text-primary-600`} />
  )

  const renderDots = () => (
    <div className="flex space-x-1">
      {[0, 1, 2].map((i) => (
        <div
          key={i}
          className={`${size === 'sm' ? 'h-2 w-2' : size === 'md' ? 'h-3 w-3' : size === 'lg' ? 'h-4 w-4' : 'h-6 w-6'} 
                     bg-primary-600 rounded-full animate-bounce`}
          style={{
            animationDelay: `${i * 0.1}s`,
            animationDuration: '0.6s'
          }}
        />
      ))}
    </div>
  )

  const renderPulse = () => (
    <div className={`${sizeClasses[size]} bg-primary-200 rounded-full animate-pulse`} />
  )

  const renderSkeleton = () => (
    <div className="animate-pulse space-y-3">
      <div className="h-4 bg-gray-200 rounded w-3/4"></div>
      <div className="h-4 bg-gray-200 rounded w-1/2"></div>
      <div className="h-4 bg-gray-200 rounded w-5/6"></div>
    </div>
  )

  const renderLoadingContent = () => {
    switch (variant) {
      case 'dots':
        return renderDots()
      case 'pulse':
        return renderPulse()
      case 'skeleton':
        return renderSkeleton()
      default:
        return renderSpinner()
    }
  }

  const loadingContent = (
    <div className={`flex flex-col items-center justify-center space-y-3 ${className}`}>
      {variant !== 'skeleton' && renderLoadingContent()}
      {variant === 'skeleton' && renderLoadingContent()}
      {message && (
        <p className={`text-gray-600 ${size === 'sm' ? 'text-sm' : size === 'lg' || size === 'xl' ? 'text-lg' : 'text-base'}`}>
          {message}
        </p>
      )}
    </div>
  )

  if (fullScreen) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-white z-50">
        <div className="text-center">
          {loadingContent}
        </div>
      </div>
    )
  }

  if (overlay) {
    return (
      <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-40">
        <div className="text-center">
          {loadingContent}
        </div>
      </div>
    )
  }

  return loadingContent
}

// Specialized loading components for common use cases
export const ButtonLoading = ({ className = '' }: { className?: string }) => (
  <Loading size="sm" className={`${className}`} />
)

export const PageLoading = ({ message = 'Loading...' }: { message?: string }) => (
  <Loading 
    size="lg" 
    message={message} 
    className="min-h-[200px]" 
  />
)

export const FullScreenLoading = ({ message = 'Loading...' }: { message?: string }) => (
  <Loading 
    size="xl" 
    message={message} 
    fullScreen 
  />
)

export const OverlayLoading = ({ message }: { message?: string }) => (
  <Loading 
    size="lg" 
    message={message} 
    overlay 
  />
)

export const SkeletonLoading = ({ className = '' }: { className?: string }) => (
  <Loading 
    variant="skeleton" 
    className={className} 
  />
)

export const InlineLoading = ({ size = 'sm', message }: { size?: 'sm' | 'md'; message?: string }) => (
  <Loading 
    size={size} 
    message={message} 
    className="inline-flex items-center" 
  />
)

export default Loading