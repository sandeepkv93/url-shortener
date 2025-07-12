import React, { useCallback, useRef, useEffect, useState } from 'react'

// Performance monitoring utilities
export class PerformanceMonitor {
  private static instance: PerformanceMonitor
  private metrics: Map<string, number[]> = new Map()
  private observers: PerformanceObserver[] = []

  static getInstance(): PerformanceMonitor {
    if (!PerformanceMonitor.instance) {
      PerformanceMonitor.instance = new PerformanceMonitor()
    }
    return PerformanceMonitor.instance
  }

  // Initialize performance monitoring
  init(): void {
    if (typeof window === 'undefined' || !('PerformanceObserver' in window)) {
      return
    }

    // Monitor navigation timing
    this.observeNavigationTiming()
    
    // Monitor resource loading
    this.observeResourceTiming()
    
    // Monitor paint timing
    this.observePaintTiming()
    
    // Monitor layout shifts (CLS)
    this.observeLayoutShifts()
    
    // Monitor long tasks
    this.observeLongTasks()
  }

  private observeNavigationTiming(): void {
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        const navigation = entry as PerformanceNavigationTiming
        this.recordMetric('navigation.domContentLoaded', navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart)
        this.recordMetric('navigation.loadComplete', navigation.loadEventEnd - navigation.loadEventStart)
        this.recordMetric('navigation.firstByte', navigation.responseStart - navigation.requestStart)
      }
    })
    
    observer.observe({ entryTypes: ['navigation'] })
    this.observers.push(observer)
  }

  private observeResourceTiming(): void {
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        const resource = entry as PerformanceResourceTiming
        this.recordMetric('resource.loadTime', resource.responseEnd - resource.startTime)
        
        // Track specific resource types
        if (resource.name.includes('.js')) {
          this.recordMetric('resource.js.loadTime', resource.responseEnd - resource.startTime)
        } else if (resource.name.includes('.css')) {
          this.recordMetric('resource.css.loadTime', resource.responseEnd - resource.startTime)
        }
      }
    })
    
    observer.observe({ entryTypes: ['resource'] })
    this.observers.push(observer)
  }

  private observePaintTiming(): void {
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        this.recordMetric(`paint.${entry.name}`, entry.startTime)
      }
    })
    
    observer.observe({ entryTypes: ['paint'] })
    this.observers.push(observer)
  }

  private observeLayoutShifts(): void {
    let cumulativeLayoutShift = 0
    
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        const shift = entry as any // CLS entries
        if (!shift.hadRecentInput) {
          cumulativeLayoutShift += shift.value
          this.recordMetric('cls.cumulative', cumulativeLayoutShift)
        }
      }
    })
    
    try {
      observer.observe({ entryTypes: ['layout-shift'] })
      this.observers.push(observer)
    } catch (e) {
      // Layout shift not supported in all browsers
      console.debug('Layout shift monitoring not supported')
    }
  }

  private observeLongTasks(): void {
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        this.recordMetric('longTask.duration', entry.duration)
      }
    })
    
    try {
      observer.observe({ entryTypes: ['longtask'] })
      this.observers.push(observer)
    } catch (e) {
      // Long task not supported in all browsers
      console.debug('Long task monitoring not supported')
    }
  }

  // Record a performance metric
  recordMetric(name: string, value: number): void {
    if (!this.metrics.has(name)) {
      this.metrics.set(name, [])
    }
    this.metrics.get(name)!.push(value)
  }

  // Get performance metrics
  getMetrics(): Record<string, { avg: number; min: number; max: number; count: number }> {
    const result: Record<string, { avg: number; min: number; max: number; count: number }> = {}
    
    for (const [name, values] of this.metrics.entries()) {
      if (values.length > 0) {
        result[name] = {
          avg: values.reduce((sum, val) => sum + val, 0) / values.length,
          min: Math.min(...values),
          max: Math.max(...values),
          count: values.length
        }
      }
    }
    
    return result
  }

  // Get Core Web Vitals
  getWebVitals(): Promise<{ fcp?: number; lcp?: number; fid?: number; cls?: number }> {
    return new Promise((resolve) => {
      const vitals: { fcp?: number; lcp?: number; fid?: number; cls?: number } = {}
      
      // First Contentful Paint
      const fcpEntry = performance.getEntriesByName('first-contentful-paint')[0]
      if (fcpEntry) {
        vitals.fcp = fcpEntry.startTime
      }
      
      // Largest Contentful Paint
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries()
        const lastEntry = entries[entries.length - 1] as any
        vitals.lcp = lastEntry.startTime
      })
      
      try {
        observer.observe({ entryTypes: ['largest-contentful-paint'] })
      } catch (e) {
        console.debug('LCP monitoring not supported')
      }
      
      // First Input Delay
      const fidObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          vitals.fid = entry.processingStart - entry.startTime
        }
      })
      
      try {
        fidObserver.observe({ entryTypes: ['first-input'] })
      } catch (e) {
        console.debug('FID monitoring not supported')
      }
      
      // Cumulative Layout Shift
      const clsMetric = this.metrics.get('cls.cumulative')
      if (clsMetric && clsMetric.length > 0) {
        vitals.cls = clsMetric[clsMetric.length - 1]
      }
      
      // Return after a short delay to allow metrics to be collected
      setTimeout(() => resolve(vitals), 1000)
    })
  }

  // Clean up observers
  destroy(): void {
    this.observers.forEach(observer => observer.disconnect())
    this.observers = []
    this.metrics.clear()
  }
}

// Debounce hook for performance optimization
export function useDebounce<T extends (...args: any[]) => any>(
  callback: T,
  delay: number
): T {
  const timeoutRef = useRef<NodeJS.Timeout>()
  
  return useCallback((...args: Parameters<T>) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    
    timeoutRef.current = setTimeout(() => {
      callback(...args)
    }, delay)
  }, [callback, delay]) as T
}

// Throttle hook for performance optimization
export function useThrottle<T extends (...args: any[]) => any>(
  callback: T,
  delay: number
): T {
  const lastRun = useRef<number>(0)
  
  return useCallback((...args: Parameters<T>) => {
    const now = Date.now()
    if (now - lastRun.current >= delay) {
      callback(...args)
      lastRun.current = now
    }
  }, [callback, delay]) as T
}

// Memory usage monitoring
export function useMemoryMonitor(): { memoryUsage: number | null; isMemoryHigh: boolean } {
  const [memoryUsage, setMemoryUsage] = useState<number | null>(null)
  const [isMemoryHigh, setIsMemoryHigh] = useState(false)
  
  useEffect(() => {
    const checkMemory = () => {
      if ('memory' in performance) {
        const memory = (performance as any).memory
        const usage = memory.usedJSHeapSize / memory.totalJSHeapSize
        setMemoryUsage(usage)
        setIsMemoryHigh(usage > 0.8) // Consider high if > 80%
      }
    }
    
    const interval = setInterval(checkMemory, 5000) // Check every 5 seconds
    checkMemory() // Initial check
    
    return () => clearInterval(interval)
  }, [])
  
  return { memoryUsage, isMemoryHigh }
}

// Component render tracking
export function useRenderTracker(componentName: string): void {
  const renderCount = useRef(0)
  
  useEffect(() => {
    renderCount.current += 1
    
    if (process.env.NODE_ENV === 'development') {
      console.debug(`${componentName} rendered ${renderCount.current} times`)
      
      // Warn about excessive re-renders
      if (renderCount.current > 10) {
        console.warn(`${componentName} has rendered ${renderCount.current} times - check for optimization opportunities`)
      }
    }
  })
}

// Virtual scrolling utility for large lists
export function useVirtualScrolling<T>(
  items: T[],
  containerHeight: number,
  itemHeight: number,
  overscan: number = 5
): {
  visibleItems: T[]
  startIndex: number
  endIndex: number
  totalHeight: number
  offsetY: number
  onScroll: (event: Event) => void
} {
  const [scrollTop, setScrollTop] = useState(0)
  
  const startIndex = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan)
  const endIndex = Math.min(
    items.length - 1,
    Math.ceil((scrollTop + containerHeight) / itemHeight) + overscan
  )
  
  const visibleItems = items.slice(startIndex, endIndex + 1)
  const totalHeight = items.length * itemHeight
  const offsetY = startIndex * itemHeight
  
  const handleScroll = useCallback((event: Event) => {
    const target = event.target as HTMLElement
    setScrollTop(target.scrollTop)
  }, [])
  
  return {
    visibleItems,
    startIndex,
    endIndex,
    totalHeight,
    offsetY,
    onScroll: handleScroll
  }
}

// Image lazy loading utility
export function useLazyImage(src: string, options?: IntersectionObserverInit): {
  imageRef: React.RefObject<HTMLImageElement>
  isLoaded: boolean
  isInView: boolean
} {
  const imageRef = useRef<HTMLImageElement>(null)
  const [isLoaded, setIsLoaded] = useState(false)
  const [isInView, setIsInView] = useState(false)
  
  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true)
          observer.disconnect()
        }
      },
      options
    )
    
    if (imageRef.current) {
      observer.observe(imageRef.current)
    }
    
    return () => observer.disconnect()
  }, [options])
  
  useEffect(() => {
    if (isInView && imageRef.current && !isLoaded) {
      imageRef.current.src = src
      imageRef.current.onload = () => setIsLoaded(true)
    }
  }, [isInView, src, isLoaded])
  
  return { imageRef, isLoaded, isInView }
}

// Bundle analyzer helper
export function analyzeBundleSize(): void {
  if (process.env.NODE_ENV === 'development') {
    import('rollup-plugin-visualizer').then(() => {
      console.log('Bundle analysis complete. Check dist/stats.html')
    }).catch(() => {
      console.log('Bundle analyzer not available. Run `npm run build:analyze` to generate bundle analysis.')
    })
  }
}

// Performance warning for slow operations
export function measurePerformance<T extends (...args: any[]) => any>(
  operation: string,
  fn: T,
  threshold: number = 100
): T {
  return ((...args: Parameters<T>) => {
    const start = performance.now()
    const result = fn(...args)
    const end = performance.now()
    
    if (end - start > threshold) {
      console.warn(`Slow operation detected: ${operation} took ${(end - start).toFixed(2)}ms`)
    }
    
    return result
  }) as T
}

// Initialize performance monitoring
export function initializePerformanceMonitoring(): void {
  if (typeof window !== 'undefined') {
    const monitor = PerformanceMonitor.getInstance()
    monitor.init()
    
    // Report metrics periodically in development
    if (process.env.NODE_ENV === 'development') {
      setInterval(() => {
        const metrics = monitor.getMetrics()
        if (Object.keys(metrics).length > 0) {
          console.group('Performance Metrics')
          console.table(metrics)
          console.groupEnd()
        }
      }, 30000) // Report every 30 seconds
    }
    
    // Report web vitals when the page is about to unload
    window.addEventListener('beforeunload', async () => {
      const vitals = await monitor.getWebVitals()
      console.log('Core Web Vitals:', vitals)
    })
  }
}

// React DevTools Profiler wrapper
export function withProfiler<P extends object>(
  Component: React.ComponentType<P>,
  id: string
): React.ComponentType<P> {
  if (process.env.NODE_ENV !== 'development') {
    return Component
  }
  
  const ProfiledComponent = (props: P) => {
    return React.createElement(
      React.Profiler,
      {
        id: id,
        onRender: (id: string, phase: string, actualDuration: number) => {
          if (actualDuration > 16) { // Warn if render takes more than 16ms (60fps)
            console.warn(`Slow render detected in ${id}: ${actualDuration.toFixed(2)}ms (${phase})`)
          }
        }
      },
      React.createElement(Component, props)
    )
  }
  
  ProfiledComponent.displayName = `withProfiler(${Component.displayName || Component.name || 'Component'})`
  
  return ProfiledComponent
}