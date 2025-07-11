import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { 
  Link as LinkIcon, 
  Copy, 
  QrCode, 
  Settings, 
  Calendar, 
  Lock, 
  Eye, 
  EyeOff, 
  Tag, 
  X,
  CheckCircle,
  AlertCircle,
  Zap
} from 'lucide-react'
import { urlService } from '@/services/urls'
import { CreateURLRequest, URL as URLType } from '@/types/url'
import { ButtonLoading } from '@/components/common/Loading'

const urlShortenerSchema = z.object({
  originalUrl: z
    .string()
    .min(1, 'URL is required')
    .url('Please enter a valid URL'),
  customAlias: z
    .string()
    .optional()
    .refine((val) => !val || /^[a-zA-Z0-9-_]+$/.test(val), 'Custom alias can only contain letters, numbers, hyphens, and underscores'),
  title: z
    .string()
    .max(100, 'Title must be less than 100 characters')
    .optional(),
  description: z
    .string()
    .max(500, 'Description must be less than 500 characters')
    .optional(),
  expiresAt: z
    .string()
    .optional()
    .refine((val) => !val || new Date(val) > new Date(), 'Expiration date must be in the future'),
  password: z
    .string()
    .min(4, 'Password must be at least 4 characters')
    .optional()
    .or(z.literal('')),
  isPublic: z.boolean().default(true),
  tags: z.array(z.string()).default([])
})

type URLShortenerFormData = z.infer<typeof urlShortenerSchema>

interface URLShortenerProps {
  onURLCreated?: (url: URLType) => void
  className?: string
}

const URLShortener = ({ onURLCreated, className = '' }: URLShortenerProps) => {
  const [isLoading, setIsLoading] = useState(false)
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [createdURL, setCreatedURL] = useState<URLType | null>(null)
  const [currentTag, setCurrentTag] = useState('')
  const [copied, setCopied] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
    setError,
    clearErrors
  } = useForm<URLShortenerFormData>({
    resolver: zodResolver(urlShortenerSchema),
    defaultValues: {
      isPublic: true,
      tags: []
    }
  })

  const watchedTags = watch('tags') || []
  const watchedPassword = watch('password')

  const handleCreateURL = async (data: URLShortenerFormData) => {
    setIsLoading(true)
    clearErrors()

    try {
      const createRequest: CreateURLRequest = {
        originalUrl: data.originalUrl,
        customAlias: data.customAlias || undefined,
        title: data.title || undefined,
        description: data.description || undefined,
        expiresAt: data.expiresAt || undefined,
        password: data.password || undefined,
        isPublic: data.isPublic,
        tags: data.tags
      }

      const newURL = await urlService.createURL(createRequest)
      setCreatedURL(newURL)
      
      if (onURLCreated) {
        onURLCreated(newURL)
      }

      // Reset form
      reset()
      setShowAdvanced(false)
      setCurrentTag('')
      
    } catch (error: any) {
      if (error.response?.status === 409 && error.response?.data?.message?.includes('alias')) {
        setError('customAlias', {
          type: 'manual',
          message: 'This custom alias is already taken'
        })
      } else if (error.response?.status === 429) {
        setError('root', {
          type: 'manual', 
          message: 'Too many URLs created. Please try again later.'
        })
      } else {
        setError('root', {
          type: 'manual',
          message: error.message || 'Failed to create short URL'
        })
      }
    } finally {
      setIsLoading(false)
    }
  }

  const addTag = () => {
    if (currentTag.trim() && !watchedTags.includes(currentTag.trim())) {
      const newTags = [...watchedTags, currentTag.trim()]
      setValue('tags', newTags)
      setCurrentTag('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    const newTags = watchedTags.filter(tag => tag !== tagToRemove)
    setValue('tags', newTags)
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
    // In a real implementation, this would generate and download a QR code
    const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodeURIComponent(url)}`
    window.open(qrCodeUrl, '_blank')
  }

  const resetForm = () => {
    setCreatedURL(null)
    reset()
    setShowAdvanced(false)
    setCurrentTag('')
  }

  // Success state
  if (createdURL) {
    const shortUrl = `${window.location.origin}/${createdURL.shortCode}`
    
    return (
      <div className={`w-full max-w-2xl mx-auto ${className}`}>
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="text-center mb-6">
            <div className="flex justify-center mb-4">
              <div className="h-12 w-12 bg-green-100 rounded-full flex items-center justify-center">
                <CheckCircle className="h-6 w-6 text-green-600" />
              </div>
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">URL Shortened Successfully!</h2>
            <p className="text-gray-600">Your short URL is ready to share</p>
          </div>

          <div className="space-y-4">
            {/* Original URL */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Original URL
              </label>
              <div className="p-3 bg-gray-50 rounded-md text-sm text-gray-600 break-all">
                {createdURL.originalUrl}
              </div>
            </div>

            {/* Short URL */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Short URL
              </label>
              <div className="flex items-center space-x-2">
                <div className="flex-1 p-3 bg-primary-50 border border-primary-200 rounded-md text-primary-700 font-medium break-all">
                  {shortUrl}
                </div>
                <button
                  onClick={() => copyToClipboard(shortUrl)}
                  className={`p-3 rounded-md transition-colors ${
                    copied 
                      ? 'bg-green-100 text-green-600' 
                      : 'bg-gray-100 hover:bg-gray-200 text-gray-600'
                  }`}
                  title="Copy to clipboard"
                >
                  <Copy className="h-5 w-5" />
                </button>
                <button
                  onClick={() => generateQRCode(shortUrl)}
                  className="p-3 bg-gray-100 hover:bg-gray-200 text-gray-600 rounded-md transition-colors"
                  title="Generate QR Code"
                >
                  <QrCode className="h-5 w-5" />
                </button>
              </div>
              {copied && (
                <p className="text-sm text-green-600 mt-1">✓ Copied to clipboard!</p>
              )}
            </div>

            {/* URL Details */}
            {(createdURL.title || createdURL.description || createdURL.tags?.length) && (
              <div className="border-t pt-4">
                <h3 className="font-medium text-gray-900 mb-2">URL Details</h3>
                {createdURL.title && (
                  <p className="text-sm text-gray-600 mb-1">
                    <span className="font-medium">Title:</span> {createdURL.title}
                  </p>
                )}
                {createdURL.description && (
                  <p className="text-sm text-gray-600 mb-1">
                    <span className="font-medium">Description:</span> {createdURL.description}
                  </p>
                )}
                {createdURL.tags && createdURL.tags.length > 0 && (
                  <div className="flex items-center space-x-2 mt-2">
                    <span className="text-sm font-medium text-gray-700">Tags:</span>
                    <div className="flex flex-wrap gap-1">
                      {createdURL.tags.map(tag => (
                        <span key={tag} className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                          {tag}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}

            {/* Actions */}
            <div className="flex space-x-3 pt-4">
              <button
                onClick={resetForm}
                className="flex-1 bg-primary-600 hover:bg-primary-700 text-white py-2 px-4 rounded-md font-medium transition-colors"
              >
                Create Another URL
              </button>
              <button
                onClick={() => window.open(`/analytics/${createdURL.id}`, '_blank')}
                className="bg-gray-600 hover:bg-gray-700 text-white py-2 px-4 rounded-md font-medium transition-colors"
              >
                View Analytics
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`w-full max-w-2xl mx-auto ${className}`}>
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-center mb-6">
          <div className="flex justify-center mb-4">
            <div className="h-12 w-12 bg-primary-600 rounded-full flex items-center justify-center">
              <Zap className="h-6 w-6 text-white" />
            </div>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Shorten Your URL</h2>
          <p className="text-gray-600">Create a short, memorable link from any URL</p>
        </div>

        {errors.root && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <div className="flex items-center">
              <AlertCircle className="h-5 w-5 text-red-500 mr-2" />
              <p className="text-sm text-red-700">{errors.root.message}</p>
            </div>
          </div>
        )}

        <form onSubmit={handleSubmit(handleCreateURL)} className="space-y-4">
          {/* Original URL Input */}
          <div>
            <label htmlFor="originalUrl" className="block text-sm font-medium text-gray-700 mb-1">
              Original URL *
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <LinkIcon className="h-5 w-5 text-gray-400" />
              </div>
              <input
                id="originalUrl"
                type="url"
                placeholder="https://example.com/very/long/url"
                className={`w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                  errors.originalUrl 
                    ? 'border-red-300 bg-red-50' 
                    : 'border-gray-300 hover:border-gray-400'
                }`}
                {...register('originalUrl')}
              />
            </div>
            {errors.originalUrl && (
              <p className="mt-1 text-sm text-red-600 flex items-center">
                <AlertCircle className="h-4 w-4 mr-1" />
                {errors.originalUrl.message}
              </p>
            )}
          </div>

          {/* Custom Alias */}
          <div>
            <label htmlFor="customAlias" className="block text-sm font-medium text-gray-700 mb-1">
              Custom Alias (Optional)
            </label>
            <div className="flex items-center">
              <span className="text-sm text-gray-500 bg-gray-50 border border-r-0 border-gray-300 rounded-l-md px-3 py-2">
                {window.location.origin}/
              </span>
              <input
                id="customAlias"
                type="text"
                placeholder="my-custom-link"
                className={`flex-1 px-3 py-2 border border-l-0 rounded-r-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                  errors.customAlias 
                    ? 'border-red-300 bg-red-50' 
                    : 'border-gray-300 hover:border-gray-400'
                }`}
                {...register('customAlias')}
              />
            </div>
            {errors.customAlias && (
              <p className="mt-1 text-sm text-red-600 flex items-center">
                <AlertCircle className="h-4 w-4 mr-1" />
                {errors.customAlias.message}
              </p>
            )}
          </div>

          {/* Advanced Options Toggle */}
          <button
            type="button"
            onClick={() => setShowAdvanced(!showAdvanced)}
            className="flex items-center space-x-2 text-primary-600 hover:text-primary-700 font-medium"
          >
            <Settings className="h-4 w-4" />
            <span>{showAdvanced ? 'Hide' : 'Show'} Advanced Options</span>
          </button>

          {/* Advanced Options */}
          {showAdvanced && (
            <div className="space-y-4 p-4 bg-gray-50 rounded-lg">
              {/* Title */}
              <div>
                <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                  Title (Optional)
                </label>
                <input
                  id="title"
                  type="text"
                  placeholder="Descriptive title for your link"
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                    errors.title 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  {...register('title')}
                />
                {errors.title && (
                  <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
                )}
              </div>

              {/* Description */}
              <div>
                <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                  Description (Optional)
                </label>
                <textarea
                  id="description"
                  rows={3}
                  placeholder="Additional description for your link"
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors resize-none ${
                    errors.description 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  {...register('description')}
                />
                {errors.description && (
                  <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
                )}
              </div>

              {/* Expiration Date */}
              <div>
                <label htmlFor="expiresAt" className="block text-sm font-medium text-gray-700 mb-1">
                  Expiration Date (Optional)
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Calendar className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    id="expiresAt"
                    type="datetime-local"
                    className={`w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                      errors.expiresAt 
                        ? 'border-red-300 bg-red-50' 
                        : 'border-gray-300 hover:border-gray-400'
                    }`}
                    {...register('expiresAt')}
                  />
                </div>
                {errors.expiresAt && (
                  <p className="mt-1 text-sm text-red-600">{errors.expiresAt.message}</p>
                )}
              </div>

              {/* Password Protection */}
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                  Password Protection (Optional)
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Lock className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    placeholder="Set a password to protect this link"
                    className={`w-full pl-10 pr-10 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                      errors.password 
                        ? 'border-red-300 bg-red-50' 
                        : 'border-gray-300 hover:border-gray-400'
                    }`}
                    {...register('password')}
                  />
                  {watchedPassword && (
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                    >
                      {showPassword ? (
                        <EyeOff className="h-5 w-5" />
                      ) : (
                        <Eye className="h-5 w-5" />
                      )}
                    </button>
                  )}
                </div>
                {errors.password && (
                  <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
                )}
              </div>

              {/* Tags */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Tags (Optional)
                </label>
                <div className="flex items-center space-x-2 mb-2">
                  <div className="relative flex-1">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Tag className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      type="text"
                      value={currentTag}
                      onChange={(e) => setCurrentTag(e.target.value)}
                      onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
                      placeholder="Add a tag"
                      className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    />
                  </div>
                  <button
                    type="button"
                    onClick={addTag}
                    className="px-3 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md transition-colors"
                  >
                    Add
                  </button>
                </div>
                {watchedTags.length > 0 && (
                  <div className="flex flex-wrap gap-2">
                    {watchedTags.map(tag => (
                      <span key={tag} className="inline-flex items-center px-2 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                        {tag}
                        <button
                          type="button"
                          onClick={() => removeTag(tag)}
                          className="ml-1 text-blue-600 hover:text-blue-800"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </span>
                    ))}
                  </div>
                )}
              </div>

              {/* Public/Private Toggle */}
              <div className="flex items-center justify-between">
                <label className="block text-sm font-medium text-gray-700">
                  Make this URL public
                </label>
                <input
                  type="checkbox"
                  className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                  {...register('isPublic')}
                />
              </div>
            </div>
          )}

          {/* Submit Button */}
          <button
            type="submit"
            disabled={isLoading}
            className={`w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white transition-colors ${
              isLoading
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500'
            }`}
          >
            {isLoading ? (
              <>
                <ButtonLoading className="mr-2" />
                Creating Short URL...
              </>
            ) : (
              <>
                <Zap className="h-4 w-4 mr-2" />
                Shorten URL
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  )
}

export default URLShortener