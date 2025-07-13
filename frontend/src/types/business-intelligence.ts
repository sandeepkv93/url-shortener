// Business Intelligence Dashboard Types

export interface Dashboard {
  id: string
  userId: string
  name: string
  description: string
  isDefault: boolean
  isPublic: boolean
  layout: DashboardLayout
  widgets?: DashboardWidget[]
  createdAt: string
  updatedAt: string
}

export interface DashboardLayout {
  gridSize: GridSize
  breakpoints: Breakpoint
  margin: [number, number]
  padding: [number, number, number, number]
  rowHeight: number
}

export interface GridSize {
  columns: number
  rows: number
}

export interface Breakpoint {
  xs: number
  sm: number
  md: number
  lg: number
  xl: number
}

export interface DashboardWidget {
  id: string
  userId: string
  type: WidgetType
  title: string
  description: string
  position: WidgetPosition
  size: WidgetSize
  config: WidgetConfig
  isVisible: boolean
  createdAt: string
  updatedAt: string
}

export interface WidgetPosition {
  x: number
  y: number
  z: number
}

export interface WidgetSize {
  width: number
  height: number
}

export type WidgetType = 'chart' | 'metric' | 'table' | 'map' | 'funnel'

export interface WidgetConfig {
  chartType?: ChartType
  dataSource?: DataSource
  filters?: Record<string, any>
  aggregationType?: AggregationType
  groupBy?: string
  timeRange?: TimeRange
  refreshRate?: number
  colors?: string[]
  showLegend?: boolean
  showGrid?: boolean
  limit?: number
}

export type ChartType = 'line' | 'bar' | 'pie' | 'donut' | 'area' | 'scatter' | 'heatmap'
export type DataSource = 'clicks' | 'urls' | 'devices' | 'geo' | 'referrers' | 'funnel' | 'performance'
export type AggregationType = 'sum' | 'avg' | 'count' | 'unique' | 'min' | 'max'
export type TimeRange = '1h' | '24h' | '7d' | '30d' | '90d' | '1y' | 'custom'

// Widget Data and Responses

export interface WidgetDataResponse {
  widgetId: string
  data: any
  metadata: Record<string, any>
  updatedAt: string
}

export interface WidgetData {
  [key: string]: any
}

export interface TimeSeriesData {
  date: string
  value: number
  label?: string
}

export interface CategoryData {
  label: string
  value: number
  percentage?: number
  color?: string
}

export interface TableData {
  headers: string[]
  rows: any[][]
}

export interface MapData {
  country: string
  value: number
  percentage: number
}

export interface FunnelData {
  step: string
  value: number
  percentage: number
}

// Advanced Analytics Types

export interface AdvancedAnalytics {
  userId: string
  period: string
  generatedAt: string
  performanceMetrics: PerformanceMetrics
  audienceInsights: AudienceInsights
  contentAnalytics: ContentAnalytics
  competitiveAnalysis: CompetitiveAnalysis
  predictiveInsights: PredictiveInsights
  recommendationEngine: RecommendationEngine
}

export interface PerformanceMetrics {
  ctr: number
  conversionRate: number
  bounceRate: number
  engagementScore: number
  qualityScore: number
  performanceTrend: Record<string, number>
  benchmarkComparison: Record<string, any>
}

export interface AudienceInsights {
  demographics: Demographics
  behaviorPatterns: BehaviorPatterns
  segmentAnalysis: AudienceSegment[]
  retentionMetrics: RetentionMetrics
  userJourney: UserJourneyStep[]
}

export interface Demographics {
  ageGroups: Record<string, number>
  gender: Record<string, number>
  locations: Record<string, number>
  devices: Record<string, number>
  languages: Record<string, number>
  interests: Record<string, number>
}

export interface BehaviorPatterns {
  clickPatterns: Record<string, any>
  timeOfDay: Record<number, number>
  dayOfWeek: Record<string, number>
  seasonalTrends: Record<string, number>
  returnVisitors: number
  averageSession: number
}

export interface AudienceSegment {
  segmentId: string
  name: string
  size: number
  characteristics: Record<string, any>
  performance: Record<string, number>
  growthRate: number
}

export interface RetentionMetrics {
  retentionRate: number
  churnRate: number
  lifetimeValue: number
  retentionCohort: Record<string, number>
  reactivationRate: number
}

export interface UserJourneyStep {
  stepNumber: number
  action: string
  url: string
  timestamp: string
  duration: number
  interactionData: Record<string, any>
}

export interface ContentAnalytics {
  topPerformingContent: ContentPerformance[]
  contentCategories: Record<string, number>
  engagementMetrics: ContentEngagement
  contentLifecycle: ContentLifecycle
  viralityMetrics: ViralityMetrics
}

export interface ContentPerformance {
  contentId: string
  title: string
  url: string
  clicks: number
  shares: number
  engagementScore: number
  category: string
  tags: string[]
}

export interface ContentEngagement {
  averageTimeOnPage: number
  shareRate: number
  commentRate: number
  likeRate: number
  engagementTrend: Record<string, number>
}

export interface ContentLifecycle {
  creationDate: string
  peakDate: string
  declineDate: string
  lifespanDays: number
  peakTraffic: number
  currentStatus: 'growing' | 'stable' | 'declining'
}

export interface ViralityMetrics {
  viralCoefficient: number
  shareVelocity: number
  reachAmplification: number
  viralPaths: ViralPath[]
}

export interface ViralPath {
  source: string
  destination: string
  shares: number
  strength: number
}

export interface CompetitiveAnalysis {
  marketPosition: MarketPosition
  competitorMetrics: CompetitorData[]
  benchmarkData: BenchmarkData
  marketTrends: MarketTrends
  opportunityGaps: OpportunityGap[]
}

export interface MarketPosition {
  rank: number
  marketShare: number
  growthRate: number
  competitiveGap: number
  strengths: string[]
  weaknesses: string[]
}

export interface CompetitorData {
  competitorId: string
  name: string
  marketShare: number
  growthRate: number
  performanceGap: number
  keyMetrics: Record<string, any>
}

export interface BenchmarkData {
  industryAverage: Record<string, number>
  topPerformers: Record<string, number>
  yourPerformance: Record<string, number>
  percentile: Record<string, number>
}

export interface MarketTrends {
  emergingTrends: Trend[]
  seasonalPatterns: Record<string, number>
  consumerBehavior: Record<string, any>
  technologyTrends: string[]
}

export interface Trend {
  name: string
  direction: 'up' | 'down' | 'stable'
  impact: 'high' | 'medium' | 'low'
  confidence: number
  description: string
}

export interface OpportunityGap {
  area: string
  impact: string
  effort: string
  priority: string
  description: string
  potential: number
}

export interface PredictiveInsights {
  forecastData: ForecastData
  anomalyDetection: Anomaly[]
  trendPrediction: TrendPrediction
  riskAssessment: RiskAssessment
  optimizationSuggestions: OptimizationSuggestion[]
}

export interface ForecastData {
  period: string
  predictions: Record<string, number>
  confidence: Record<string, number>
  scenarios: Record<string, Scenario>
}

export interface Scenario {
  name: string
  probability: number
  predictions: Record<string, number>
  impact: string
}

export interface Anomaly {
  type: 'spike' | 'drop' | 'outlier'
  metric: string
  severity: 'high' | 'medium' | 'low'
  detected: string
  value: number
  expected: number
  deviation: number
  description: string
  impact: string
}

export interface TrendPrediction {
  shortTerm: TrendForecast
  mediumTerm: TrendForecast
  longTerm: TrendForecast
}

export interface TrendForecast {
  direction: string
  magnitude: number
  confidence: number
  factorsInfluencing: string[]
}

export interface RiskAssessment {
  overallRisk: RiskLevel
  riskFactors: RiskFactor[]
  mitigation: Mitigation[]
}

export interface RiskLevel {
  level: 'low' | 'medium' | 'high' | 'critical'
  score: number
  description: string
  trend: 'increasing' | 'stable' | 'decreasing'
}

export interface RiskFactor {
  factor: string
  impact: string
  probability: number
  description: string
}

export interface Mitigation {
  strategy: string
  priority: string
  effort: string
  effectiveness: number
  description: string
}

export interface OptimizationSuggestion {
  area: string
  suggestion: string
  impact: string
  effort: string
  priority: string
  potential: number
  confidence: number
  implementation: string
}

export interface RecommendationEngine {
  personalizedRecommendations: Recommendation[]
  contentRecommendations: ContentRec[]
  optimizationRecommendations: OptimizationRec[]
  audienceRecommendations: AudienceRec[]
}

export interface Recommendation {
  type: string
  title: string
  description: string
  priority: string
  impact: string
  confidence: number
  actionItems: string[]
}

export interface ContentRec {
  contentType: string
  topics: string[]
  timing: string
  channels: string[]
  potential: number
}

export interface OptimizationRec {
  metric: string
  currentValue: number
  targetValue: number
  strategy: string
  timeline: string
}

export interface AudienceRec {
  segmentName: string
  characteristics: string[]
  size: number
  potential: number
  approach: string
}

// Request/Response Types

export interface CreateDashboardRequest {
  name: string
  description: string
  isDefault: boolean
  isPublic: boolean
  layout: DashboardLayout
}

export interface UpdateDashboardRequest {
  name?: string
  description?: string
  isDefault?: boolean
  isPublic?: boolean
  layout?: DashboardLayout
}

export interface CreateWidgetRequest {
  type: WidgetType
  title: string
  description: string
  position: WidgetPosition
  size: WidgetSize
  config: WidgetConfig
}

export interface UpdateWidgetRequest {
  title?: string
  description?: string
  position?: WidgetPosition
  size?: WidgetSize
  config?: WidgetConfig
  isVisible?: boolean
}

export interface DashboardResponse {
  dashboard: Dashboard
  widgetCount: number
}

export interface DashboardListResponse {
  dashboards: DashboardResponse[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// UI State Types

export interface DashboardState {
  currentDashboard: Dashboard | null
  dashboards: Dashboard[]
  widgets: DashboardWidget[]
  isLoading: boolean
  isEditing: boolean
  selectedWidget: DashboardWidget | null
  draggedWidget: DashboardWidget | null
  showWidgetPalette: boolean
  showSettings: boolean
  error: BIError | null
}

export interface WidgetPaletteItem {
  type: WidgetType
  title: string
  description: string
  icon: React.ReactNode
  defaultSize: WidgetSize
  defaultConfig: Partial<WidgetConfig>
}

export interface DragItem {
  type: 'widget' | 'palette-item'
  widget?: DashboardWidget
  paletteItem?: WidgetPaletteItem
  position?: WidgetPosition
}

// Event Types

export interface WidgetMoveEvent {
  widgetId: string
  newPosition: WidgetPosition
}

export interface WidgetResizeEvent {
  widgetId: string
  newSize: WidgetSize
}

export interface WidgetConfigChangeEvent {
  widgetId: string
  newConfig: WidgetConfig
}

// Error Types

export interface BIError {
  type: 'dashboard' | 'widget' | 'data' | 'network'
  message: string
  details?: any
}

// Filter and Search Types

export interface DashboardFilter {
  timeRange: TimeRange
  startDate?: string
  endDate?: string
  userId?: string
  tags?: string[]
}

export interface WidgetFilter {
  type?: WidgetType[]
  dataSource?: DataSource[]
  isVisible?: boolean
  userId?: string
}

export interface AnalyticsFilter {
  period: TimeRange
  startDate?: string
  endDate?: string
  urlIds?: string[]
  countries?: string[]
  devices?: string[]
  browsers?: string[]
  referrers?: string[]
}

// Export Types

export interface ExportOptions {
  format: 'pdf' | 'png' | 'svg' | 'json'
  includeData: boolean
  includeCharts: boolean
  includeMetadata: boolean
  quality?: 'low' | 'medium' | 'high'
}

export interface ExportResult {
  url: string
  filename: string
  size: number
  expiresAt: string
}

// Utility Types

export type DashboardAction = 
  | { type: 'SET_CURRENT_DASHBOARD'; payload: Dashboard }
  | { type: 'SET_DASHBOARDS'; payload: Dashboard[] }
  | { type: 'ADD_DASHBOARD'; payload: Dashboard }
  | { type: 'UPDATE_DASHBOARD'; payload: Dashboard }
  | { type: 'REMOVE_DASHBOARD'; payload: string }
  | { type: 'ADD_WIDGET'; payload: DashboardWidget }
  | { type: 'UPDATE_WIDGET'; payload: DashboardWidget }
  | { type: 'REMOVE_WIDGET'; payload: string }
  | { type: 'SET_EDITING'; payload: boolean }
  | { type: 'SET_SELECTED_WIDGET'; payload: DashboardWidget | null }
  | { type: 'SET_DRAGGED_WIDGET'; payload: DashboardWidget | null }
  | { type: 'SET_WIDGET_PALETTE_VISIBLE'; payload: boolean }
  | { type: 'SET_SETTINGS_VISIBLE'; payload: boolean }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: BIError | null }

export interface DashboardContextType {
  state: DashboardState
  dispatch: React.Dispatch<DashboardAction>
  
  // Dashboard operations
  createDashboard: (request: CreateDashboardRequest) => Promise<Dashboard>
  updateDashboard: (id: string, request: UpdateDashboardRequest) => Promise<Dashboard>
  deleteDashboard: (id: string) => Promise<void>
  loadDashboard: (id: string) => Promise<void>
  loadDashboards: () => Promise<void>
  
  // Widget operations
  createWidget: (request: CreateWidgetRequest) => Promise<DashboardWidget>
  updateWidget: (id: string, request: UpdateWidgetRequest) => Promise<DashboardWidget>
  deleteWidget: (id: string) => Promise<void>
  getWidgetData: (id: string) => Promise<WidgetDataResponse>
  moveWidget: (widgetId: string, newPosition: { x: number; y: number }) => void
  resizeWidget: (widgetId: string, newSize: { width: number; height: number }) => void
  
  // Export operations
  exportDashboard: (id: string, options: ExportOptions) => Promise<ExportResult>
  
  // UI state
  toggleEditing: () => void
  selectWidget: (widget: DashboardWidget | null) => void
  toggleWidgetPalette: () => void
  toggleSettings: () => void
  clearError: () => void
  
  // Bulk operations
  bulkUpdateWidgets: (updates: Array<{ id: string; request: UpdateWidgetRequest }>) => Promise<DashboardWidget[]>
  bulkDeleteWidgets: (ids: string[]) => Promise<void>
}