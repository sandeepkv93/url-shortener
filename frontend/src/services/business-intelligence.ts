import axios from 'axios'
import { api } from './api'
import {
  Dashboard,
  DashboardListResponse,
  CreateDashboardRequest,
  UpdateDashboardRequest,
  DashboardWidget,
  CreateWidgetRequest,
  UpdateWidgetRequest,
  WidgetDataResponse,
  AdvancedAnalytics,
  PerformanceMetrics,
  AudienceInsights,
  ContentAnalytics,
  CompetitiveAnalysis,
  MarketPosition,
  BenchmarkData,
  PredictiveInsights,
  ForecastData,
  Anomaly,
  TrendPrediction,
  RecommendationEngine,
  OptimizationSuggestion,
  ContentRec,
  AudienceRec,
  ExportOptions,
  ExportResult
} from '../types/business-intelligence'

class BusinessIntelligenceService {
  private baseURL = '/api/v1/bi'

  // Dashboard Management

  async getDashboards(): Promise<DashboardListResponse> {
    const response = await api.get<DashboardListResponse>(`${this.baseURL}/dashboards`)
    return response.data
  }

  async getDashboard(id: string): Promise<Dashboard> {
    const response = await api.get<Dashboard>(`${this.baseURL}/dashboards/${id}`)
    return response.data
  }

  async createDashboard(request: CreateDashboardRequest): Promise<Dashboard> {
    const response = await api.post<Dashboard>(`${this.baseURL}/dashboards`, request)
    return response.data
  }

  async updateDashboard(id: string, request: UpdateDashboardRequest): Promise<Dashboard> {
    const response = await api.put<Dashboard>(`${this.baseURL}/dashboards/${id}`, request)
    return response.data
  }

  async deleteDashboard(id: string): Promise<void> {
    await api.delete(`${this.baseURL}/dashboards/${id}`)
  }

  // Widget Management

  async getWidget(id: string): Promise<DashboardWidget> {
    const response = await api.get<DashboardWidget>(`${this.baseURL}/widgets/${id}`)
    return response.data
  }

  async createWidget(request: CreateWidgetRequest): Promise<DashboardWidget> {
    const response = await api.post<DashboardWidget>(`${this.baseURL}/widgets`, request)
    return response.data
  }

  async updateWidget(id: string, request: UpdateWidgetRequest): Promise<DashboardWidget> {
    const response = await api.put<DashboardWidget>(`${this.baseURL}/widgets/${id}`, request)
    return response.data
  }

  async deleteWidget(id: string): Promise<void> {
    await api.delete(`${this.baseURL}/widgets/${id}`)
  }

  async getWidgetData(id: string): Promise<WidgetDataResponse> {
    const response = await api.get<WidgetDataResponse>(`${this.baseURL}/widgets/${id}/data`)
    return response.data
  }

  // Advanced Analytics

  async getAdvancedAnalytics(period = '30d'): Promise<AdvancedAnalytics> {
    const response = await api.get<AdvancedAnalytics>(`${this.baseURL}/analytics/advanced`, {
      params: { period }
    })
    return response.data
  }

  async getPerformanceMetrics(period = '30d'): Promise<PerformanceMetrics> {
    const response = await api.get<PerformanceMetrics>(`${this.baseURL}/analytics/performance`, {
      params: { period }
    })
    return response.data
  }

  async getAudienceInsights(period = '30d'): Promise<AudienceInsights> {
    const response = await api.get<AudienceInsights>(`${this.baseURL}/analytics/audience`, {
      params: { period }
    })
    return response.data
  }

  async getContentAnalytics(period = '30d'): Promise<ContentAnalytics> {
    const response = await api.get<ContentAnalytics>(`${this.baseURL}/analytics/content`, {
      params: { period }
    })
    return response.data
  }

  // Competitive Analysis

  async getCompetitiveAnalysis(): Promise<CompetitiveAnalysis> {
    const response = await api.get<CompetitiveAnalysis>(`${this.baseURL}/competitive/analysis`)
    return response.data
  }

  async getMarketPosition(): Promise<MarketPosition> {
    const response = await api.get<MarketPosition>(`${this.baseURL}/competitive/market-position`)
    return response.data
  }

  async getBenchmarkData(metric = 'all'): Promise<BenchmarkData> {
    const response = await api.get<BenchmarkData>(`${this.baseURL}/competitive/benchmarks`, {
      params: { metric }
    })
    return response.data
  }

  // Predictive Insights

  async getPredictiveInsights(): Promise<PredictiveInsights> {
    const response = await api.get<PredictiveInsights>(`${this.baseURL}/predictive/insights`)
    return response.data
  }

  async getForecastData(metric = 'clicks', period = '30d'): Promise<ForecastData> {
    const response = await api.get<ForecastData>(`${this.baseURL}/predictive/forecast`, {
      params: { metric, period }
    })
    return response.data
  }

  async getAnomalies(metric = 'clicks'): Promise<{ anomalies: Anomaly[]; metric: string; count: number }> {
    const response = await api.get(`${this.baseURL}/predictive/anomalies`, {
      params: { metric }
    })
    return response.data
  }

  async getTrendPrediction(metric = 'clicks'): Promise<TrendPrediction> {
    const response = await api.get<TrendPrediction>(`${this.baseURL}/predictive/trends`, {
      params: { metric }
    })
    return response.data
  }

  // Recommendations

  async getRecommendations(): Promise<RecommendationEngine> {
    const response = await api.get<RecommendationEngine>(`${this.baseURL}/recommendations`)
    return response.data
  }

  async getOptimizationSuggestions(): Promise<{ suggestions: OptimizationSuggestion[]; count: number }> {
    const response = await api.get(`${this.baseURL}/recommendations/optimizations`)
    return response.data
  }

  async getContentRecommendations(): Promise<{ recommendations: ContentRec[]; count: number }> {
    const response = await api.get(`${this.baseURL}/recommendations/content`)
    return response.data
  }

  async getAudienceRecommendations(): Promise<{ recommendations: AudienceRec[]; count: number }> {
    const response = await api.get(`${this.baseURL}/recommendations/audience`)
    return response.data
  }

  // Export

  async exportDashboard(id: string, options: ExportOptions): Promise<ExportResult> {
    const response = await api.post<ExportResult>(`${this.baseURL}/dashboards/${id}/export`, options)
    return response.data
  }

  async exportWidget(id: string, options: ExportOptions): Promise<ExportResult> {
    const response = await api.post<ExportResult>(`${this.baseURL}/widgets/${id}/export`, options)
    return response.data
  }

  // Real-time Data

  async subscribeToWidgetUpdates(
    widgetId: string, 
    callback: (data: WidgetDataResponse) => void
  ): Promise<() => void> {
    // In a real implementation, this would use WebSockets or Server-Sent Events
    // For now, we'll use polling
    const interval = setInterval(async () => {
      try {
        const data = await this.getWidgetData(widgetId)
        callback(data)
      } catch (error) {
        console.error('Error fetching widget data:', error)
      }
    }, 30000) // Update every 30 seconds

    return () => clearInterval(interval)
  }

  // Bulk Operations

  async bulkUpdateWidgets(updates: Array<{ id: string; request: UpdateWidgetRequest }>): Promise<DashboardWidget[]> {
    const promises = updates.map(({ id, request }) => this.updateWidget(id, request))
    return Promise.all(promises)
  }

  async bulkDeleteWidgets(ids: string[]): Promise<void> {
    const promises = ids.map(id => this.deleteWidget(id))
    await Promise.all(promises)
  }

  // Template Operations

  async getDashboardTemplates(): Promise<Dashboard[]> {
    const response = await api.get<Dashboard[]>(`${this.baseURL}/templates/dashboards`)
    return response.data
  }

  async getWidgetTemplates(): Promise<DashboardWidget[]> {
    const response = await api.get<DashboardWidget[]>(`${this.baseURL}/templates/widgets`)
    return response.data
  }

  async createDashboardFromTemplate(templateId: string, name: string): Promise<Dashboard> {
    const response = await api.post<Dashboard>(`${this.baseURL}/templates/dashboards/${templateId}/create`, {
      name
    })
    return response.data
  }

  // Collaboration

  async shareDashboard(id: string, emails: string[], permissions: 'view' | 'edit' = 'view'): Promise<void> {
    await api.post(`${this.baseURL}/dashboards/${id}/share`, {
      emails,
      permissions
    })
  }

  async getDashboardCollaborators(id: string): Promise<Array<{ email: string; permissions: string; addedAt: string }>> {
    const response = await api.get(`${this.baseURL}/dashboards/${id}/collaborators`)
    return response.data
  }

  async removeDashboardCollaborator(dashboardId: string, email: string): Promise<void> {
    await api.delete(`${this.baseURL}/dashboards/${dashboardId}/collaborators/${email}`)
  }

  // Comments and Annotations

  async addDashboardComment(dashboardId: string, comment: string, position?: { x: number; y: number }): Promise<void> {
    await api.post(`${this.baseURL}/dashboards/${dashboardId}/comments`, {
      comment,
      position
    })
  }

  async getDashboardComments(dashboardId: string): Promise<Array<{
    id: string
    author: string
    comment: string
    position?: { x: number; y: number }
    createdAt: string
  }>> {
    const response = await api.get(`${this.baseURL}/dashboards/${dashboardId}/comments`)
    return response.data
  }

  // Version Control

  async saveDashboardVersion(dashboardId: string, description?: string): Promise<{
    versionId: string
    version: number
    description: string
    createdAt: string
  }> {
    const response = await api.post(`${this.baseURL}/dashboards/${dashboardId}/versions`, {
      description
    })
    return response.data
  }

  async getDashboardVersions(dashboardId: string): Promise<Array<{
    versionId: string
    version: number
    description: string
    createdAt: string
  }>> {
    const response = await api.get(`${this.baseURL}/dashboards/${dashboardId}/versions`)
    return response.data
  }

  async restoreDashboardVersion(dashboardId: string, versionId: string): Promise<Dashboard> {
    const response = await api.post<Dashboard>(`${this.baseURL}/dashboards/${dashboardId}/versions/${versionId}/restore`)
    return response.data
  }

  // Performance and Optimization

  async getPerformanceReport(dashboardId: string): Promise<{
    loadTime: number
    renderTime: number
    memoryUsage: number
    widgetPerformance: Array<{
      widgetId: string
      loadTime: number
      renderTime: number
      errorRate: number
    }>
    recommendations: string[]
  }> {
    const response = await api.get(`${this.baseURL}/dashboards/${dashboardId}/performance`)
    return response.data
  }

  async optimizeDashboard(dashboardId: string): Promise<{
    optimizations: string[]
    estimatedImprovement: {
      loadTime: number
      renderTime: number
      memoryUsage: number
    }
  }> {
    const response = await api.post(`${this.baseURL}/dashboards/${dashboardId}/optimize`)
    return response.data
  }

  // Usage Analytics

  async getDashboardUsageStats(dashboardId: string, period = '30d'): Promise<{
    views: number
    uniqueViewers: number
    averageViewTime: number
    interactionRate: number
    mostUsedWidgets: Array<{
      widgetId: string
      title: string
      interactions: number
    }>
    usage_trend: Array<{
      date: string
      views: number
      interactions: number
    }>
  }> {
    const response = await api.get(`${this.baseURL}/dashboards/${dashboardId}/usage`, {
      params: { period }
    })
    return response.data
  }

  // Error Handling and Retry Logic

  private async retryRequest<T>(
    operation: () => Promise<T>,
    maxRetries = 3,
    delay = 1000
  ): Promise<T> {
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      try {
        return await operation()
      } catch (error) {
        if (attempt === maxRetries) {
          throw error
        }
        
        // Exponential backoff
        await new Promise(resolve => setTimeout(resolve, delay * Math.pow(2, attempt - 1)))
      }
    }
    
    throw new Error('Max retries exceeded')
  }

  // Batch Operations with Error Handling

  async batchOperation<T, R>(
    items: T[],
    operation: (item: T) => Promise<R>,
    batchSize = 5
  ): Promise<Array<{ success: boolean; result?: R; error?: string; item: T }>> {
    const results: Array<{ success: boolean; result?: R; error?: string; item: T }> = []
    
    for (let i = 0; i < items.length; i += batchSize) {
      const batch = items.slice(i, i + batchSize)
      
      const batchResults = await Promise.allSettled(
        batch.map(async item => {
          try {
            const result = await operation(item)
            return { success: true, result, item }
          } catch (error) {
            return { 
              success: false, 
              error: error instanceof Error ? error.message : 'Unknown error',
              item 
            }
          }
        })
      )
      
      results.push(...batchResults.map(result => 
        result.status === 'fulfilled' 
          ? result.value 
          : { success: false, error: 'Promise rejected', item: batch[0] }
      ))
    }
    
    return results
  }
}

export const biService = new BusinessIntelligenceService()
export default biService