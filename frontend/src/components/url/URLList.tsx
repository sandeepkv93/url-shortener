import { useState, useEffect, useMemo, useCallback, memo } from 'react'
import { 
  Search, 
  Filter, 
  SortAsc, 
  SortDesc, 
  RefreshCw,
  Plus,
  Tag,
  Calendar,
  Globe,
  Lock,
  Eye,
  EyeOff,
  ChevronLeft,
  ChevronRight
} from 'lucide-react'
import { URLFilter, URL as URLType, URLListResponse } from '@/types/url'
import { urlService } from '@/services/urls'
import URLCard from './URLCard'
import { PageLoading } from '@/components/common/Loading'

interface URLListProps {
  onURLCreate?: () => void
  onURLUpdate?: (url: URLType) => void
  onURLDelete?: (urlId: string) => void
  className?: string
}

const URLList = ({ 
  onURLCreate, 
  onURLUpdate, 
  onURLDelete, 
  className = '' 
}: URLListProps) => {
  const [urls, setUrls] = useState<URLType[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 12,
    total: 0,
    hasNext: false,
    hasPrev: false
  })

  const [filters, setFilters] = useState<URLFilter>({
    search: '',
    isActive: undefined,
    isPublic: undefined,
    hasExpiry: undefined,
    hasPassword: undefined,
    tags: [],
    sortBy: 'createdAt',
    sortOrder: 'desc',
    page: 1,
    limit: 12
  })

  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [availableTags, setAvailableTags] = useState<string[]>([])

  // Fetch URLs with current filters
  const fetchURLs = useCallback(async (showLoading = true) => {
    if (showLoading) setIsLoading(true)
    else setIsRefreshing(true)

    try {
      const response = await urlService.getUserURLs(filters)
      setUrls(response.urls)
      setPagination({
        page: response.page,
        limit: response.limit,
        total: response.total,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      })

      // Extract unique tags for filter options
      const allTags = response.urls.flatMap(url => url.tags || [])
      const uniqueTags = Array.from(new Set(allTags))
      setAvailableTags(uniqueTags)

    } catch (error) {
      console.error('Failed to fetch URLs:', error)
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }, [filters])

  // Initial load
  useEffect(() => {
    fetchURLs()
  }, [])

  // Apply filters when they change
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setFilters(prev => ({
        ...prev,
        search: searchTerm,
        tags: selectedTags,
        page: 1 // Reset to first page when filters change
      }))
    }, 300) // Debounce search

    return () => clearTimeout(timeoutId)
  }, [searchTerm, selectedTags])

  // Fetch URLs when filters change
  useEffect(() => {
    fetchURLs(false)
  }, [filters])

  const handleFilterChange = useCallback((key: keyof URLFilter, value: any) => {
    setFilters(prev => ({
      ...prev,
      [key]: value,
      page: 1 // Reset pagination when filters change
    }))
  }, [])

  const handlePageChange = useCallback((newPage: number) => {
    setFilters(prev => ({ ...prev, page: newPage }))
  }, [])

  const handleSortChange = useCallback((sortBy: string) => {
    const newSortOrder = filters.sortBy === sortBy && filters.sortOrder === 'desc' ? 'asc' : 'desc'
    setFilters(prev => ({
      ...prev,
      sortBy: sortBy as any,
      sortOrder: newSortOrder
    }))
  }, [filters.sortBy, filters.sortOrder])

  const handleURLUpdate = useCallback((updatedURL: URLType) => {
    setUrls(prev => prev.map(url => url.id === updatedURL.id ? updatedURL : url))
    if (onURLUpdate) {
      onURLUpdate(updatedURL)
    }
  }, [onURLUpdate])

  const handleURLDelete = useCallback((urlId: string) => {
    setUrls(prev => prev.filter(url => url.id !== urlId))
    if (onURLDelete) {
      onURLDelete(urlId)
    }
  }, [onURLDelete])

  const clearFilters = useCallback(() => {
    setSearchTerm('')
    setSelectedTags([])
    setFilters({
      search: '',
      isActive: undefined,
      isPublic: undefined,
      hasExpiry: undefined,
      hasPassword: undefined,
      tags: [],
      sortBy: 'createdAt',
      sortOrder: 'desc',
      page: 1,
      limit: 12
    })
  }, [])

  const hasActiveFilters = useMemo(() => {
    return searchTerm || 
           selectedTags.length > 0 ||
           filters.isActive !== undefined ||
           filters.isPublic !== undefined ||
           filters.hasExpiry !== undefined ||
           filters.hasPassword !== undefined
  }, [searchTerm, selectedTags, filters])

  const getSortIcon = useCallback((sortKey: string) => {
    if (filters.sortBy !== sortKey) return null
    return filters.sortOrder === 'desc' ? <SortDesc className="h-4 w-4" /> : <SortAsc className="h-4 w-4" />
  }, [filters.sortBy, filters.sortOrder])

  if (isLoading) {
    return <PageLoading message="Loading your URLs..." />
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Your URLs</h2>
          <p className="text-gray-600">
            {pagination.total} URL{pagination.total !== 1 ? 's' : ''} total
            {hasActiveFilters && (
              <span className="ml-2 text-primary-600 font-medium">
                (filtered)
              </span>
            )}
          </p>
        </div>
        <div className="flex items-center space-x-3">
          <button
            onClick={() => fetchURLs(false)}
            disabled={isRefreshing}
            className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md transition-colors"
            title="Refresh"
          >
            <RefreshCw className={`h-5 w-5 ${isRefreshing ? 'animate-spin' : ''}`} />
          </button>
          {onURLCreate && (
            <button
              onClick={onURLCreate}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Create URL
            </button>
          )}
        </div>
      </div>

      {/* Search and Filters */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
        <div className="space-y-4">
          {/* Search Bar */}
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search URLs by title, description, or original URL..."
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>

          {/* Filter Toggle */}
          <div className="flex items-center justify-between">
            <button
              onClick={() => setShowFilters(!showFilters)}
              className="flex items-center space-x-2 text-gray-600 hover:text-gray-900"
            >
              <Filter className="h-4 w-4" />
              <span>Filters</span>
              {hasActiveFilters && (
                <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-800">
                  Active
                </span>
              )}
            </button>
            {hasActiveFilters && (
              <button
                onClick={clearFilters}
                className="text-sm text-primary-600 hover:text-primary-700"
              >
                Clear all filters
              </button>
            )}
          </div>

          {/* Filters Panel */}
          {showFilters && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pt-4 border-t border-gray-200">
              {/* Status Filter */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Status
                </label>
                <select
                  value={filters.isActive === undefined ? '' : filters.isActive.toString()}
                  onChange={(e) => handleFilterChange('isActive', e.target.value === '' ? undefined : e.target.value === 'true')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="">All URLs</option>
                  <option value="true">Active only</option>
                  <option value="false">Inactive only</option>
                </select>
              </div>

              {/* Visibility Filter */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Visibility
                </label>
                <select
                  value={filters.isPublic === undefined ? '' : filters.isPublic.toString()}
                  onChange={(e) => handleFilterChange('isPublic', e.target.value === '' ? undefined : e.target.value === 'true')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="">All URLs</option>
                  <option value="true">Public only</option>
                  <option value="false">Private only</option>
                </select>
              </div>

              {/* Special Properties */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Properties
                </label>
                <div className="space-y-2">
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      checked={filters.hasExpiry === true}
                      onChange={(e) => handleFilterChange('hasExpiry', e.target.checked ? true : undefined)}
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    />
                    <span className="ml-2 text-sm text-gray-700">Has expiration</span>
                  </label>
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      checked={filters.hasPassword === true}
                      onChange={(e) => handleFilterChange('hasPassword', e.target.checked ? true : undefined)}
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    />
                    <span className="ml-2 text-sm text-gray-700">Password protected</span>
                  </label>
                </div>
              </div>

              {/* Tags Filter */}
              {availableTags.length > 0 && (
                <div className="md:col-span-2 lg:col-span-3">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Tags
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {availableTags.map(tag => (
                      <button
                        key={tag}
                        onClick={() => {
                          if (selectedTags.includes(tag)) {
                            setSelectedTags(prev => prev.filter(t => t !== tag))
                          } else {
                            setSelectedTags(prev => [...prev, tag])
                          }
                        }}
                        className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                          selectedTags.includes(tag)
                            ? 'bg-primary-100 text-primary-800 border border-primary-200'
                            : 'bg-gray-100 text-gray-700 border border-gray-200 hover:bg-gray-200'
                        }`}
                      >
                        <Tag className="h-3 w-3 mr-1" />
                        {tag}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Sort Options */}
      <div className="flex items-center space-x-4 text-sm">
        <span className="text-gray-600">Sort by:</span>
        <button
          onClick={() => handleSortChange('createdAt')}
          className={`flex items-center space-x-1 px-3 py-1 rounded-md transition-colors ${
            filters.sortBy === 'createdAt' 
              ? 'bg-primary-100 text-primary-700' 
              : 'text-gray-600 hover:bg-gray-100'
          }`}
        >
          <Calendar className="h-4 w-4" />
          <span>Date Created</span>
          {getSortIcon('createdAt')}
        </button>
        <button
          onClick={() => handleSortChange('clickCount')}
          className={`flex items-center space-x-1 px-3 py-1 rounded-md transition-colors ${
            filters.sortBy === 'clickCount' 
              ? 'bg-primary-100 text-primary-700' 
              : 'text-gray-600 hover:bg-gray-100'
          }`}
        >
          <span>Clicks</span>
          {getSortIcon('clickCount')}
        </button>
        <button
          onClick={() => handleSortChange('title')}
          className={`flex items-center space-x-1 px-3 py-1 rounded-md transition-colors ${
            filters.sortBy === 'title' 
              ? 'bg-primary-100 text-primary-700' 
              : 'text-gray-600 hover:bg-gray-100'
          }`}
        >
          <span>Title</span>
          {getSortIcon('title')}
        </button>
      </div>

      {/* URL Grid */}
      {urls.length === 0 ? (
        <div className="text-center py-12">
          <div className="mx-auto h-24 w-24 text-gray-400 mb-4">
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            {hasActiveFilters ? 'No URLs match your filters' : 'No URLs yet'}
          </h3>
          <p className="text-gray-600 mb-6">
            {hasActiveFilters 
              ? 'Try adjusting your search or filter criteria.'
              : 'Create your first short URL to get started.'
            }
          </p>
          {hasActiveFilters ? (
            <button
              onClick={clearFilters}
              className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
            >
              Clear filters
            </button>
          ) : onURLCreate ? (
            <button
              onClick={onURLCreate}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700"
            >
              <Plus className="h-4 w-4 mr-2" />
              Create your first URL
            </button>
          ) : null}
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
            {urls.map(url => (
              <URLCard
                key={url.id}
                url={url}
                onUpdate={handleURLUpdate}
                onDelete={handleURLDelete}
              />
            ))}
          </div>

          {/* Pagination */}
          {pagination.total > pagination.limit && (
            <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6 rounded-lg">
              <div className="flex flex-1 justify-between sm:hidden">
                <button
                  onClick={() => handlePageChange(pagination.page - 1)}
                  disabled={!pagination.hasPrev}
                  className="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                <button
                  onClick={() => handlePageChange(pagination.page + 1)}
                  disabled={!pagination.hasNext}
                  className="relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </div>
              <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                <div>
                  <p className="text-sm text-gray-700">
                    Showing{' '}
                    <span className="font-medium">
                      {(pagination.page - 1) * pagination.limit + 1}
                    </span>{' '}
                    to{' '}
                    <span className="font-medium">
                      {Math.min(pagination.page * pagination.limit, pagination.total)}
                    </span>{' '}
                    of{' '}
                    <span className="font-medium">{pagination.total}</span> results
                  </p>
                </div>
                <div>
                  <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                    <button
                      onClick={() => handlePageChange(pagination.page - 1)}
                      disabled={!pagination.hasPrev}
                      className="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <ChevronLeft className="h-5 w-5" />
                    </button>
                    <span className="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300">
                      Page {pagination.page} of {Math.ceil(pagination.total / pagination.limit)}
                    </span>
                    <button
                      onClick={() => handlePageChange(pagination.page + 1)}
                      disabled={!pagination.hasNext}
                      className="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <ChevronRight className="h-5 w-5" />
                    </button>
                  </nav>
                </div>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  )
}

export default memo(URLList)