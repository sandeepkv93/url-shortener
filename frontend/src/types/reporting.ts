export interface ScheduledReport {
  id: string
  user_id: number
  name: string
  description: string
  report_type: 'analytics' | 'dashboard' | 'funnel' | 'competitive'
  schedule: string
  recipients: string[]
  format: 'pdf' | 'excel' | 'csv' | 'json'
  config: ReportConfig
  is_active: boolean
  last_run?: string
  next_run?: string
  created_at: string
  updated_at: string
}

export interface ReportConfig {
  date_range: string
  metrics: string[]
  grouping?: string
  comparison?: boolean
  charts?: boolean
  raw_data?: boolean
  filters?: Record<string, any>
}

export interface CreateScheduledReportRequest {
  name: string
  description?: string
  report_type: string
  schedule: string
  recipients: string[]
  format: string
  config: ReportConfig
}

export interface UpdateScheduledReportRequest extends CreateScheduledReportRequest {}

export interface ReportExecution {
  id: number
  report_id: number
  status: 'pending' | 'running' | 'completed' | 'failed'
  started_at: string
  completed_at?: string
  duration?: number
  file_path?: string
  file_size?: number
  error_message?: string
  created_at: string
  updated_at: string
}

export interface DataExport {
  id: string
  user_id: number
  export_type: 'analytics' | 'urls' | 'clicks' | 'users'
  format: 'csv' | 'excel' | 'json'
  status: 'pending' | 'processing' | 'completed' | 'failed'
  file_path?: string
  file_size?: number
  record_count?: number
  config: ExportConfig
  expires_at: string
  created_at: string
  updated_at: string
}

export interface ExportConfig {
  date_range: {
    start_date: string
    end_date: string
  }
  columns: string[]
  compression?: boolean
  filters?: Record<string, any>
}

export interface DataExportRequest {
  export_type: string
  format: string
  config: ExportConfig
}

export interface ReportTemplate {
  id: string
  name: string
  description: string
  type: string
  format: string
  schedule: string
  config: Record<string, any>
}

export interface ReportInsights {
  total_reports: number
  active_reports: number
  reports_sent: number
  last_30_days: number
  success_rate: number
  avg_generation_time: number
  most_popular_format: string
  most_popular_schedule: string
  recent_activity: Array<{
    date: string
    action: string
    report: string
    status: string
  }>
  recommendations: string[]
}

export interface JobStatus {
  id: string
  type: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  progress: number
  message?: string
  result?: Record<string, any>
  error?: string
  started_at: string
  completed_at?: string
  created_at: string
  updated_at: string
}

// API Response types
export interface ScheduledReportsResponse {
  reports: ScheduledReport[]
  total: number
}

export interface DataExportsResponse {
  exports: DataExport[]
  total: number
}

export interface ReportHistoryResponse {
  report_id: number
  executions: ReportExecution[]
  total: number
}

export interface ReportTemplatesResponse {
  templates: ReportTemplate[]
  total: number
}

// Error types
export interface ReportingError {
  type: 'report' | 'export' | 'template' | 'execution'
  message: string
  details?: any
}

// UI State types
export interface ReportingState {
  scheduledReports: ScheduledReport[]
  dataExports: DataExport[]
  currentReport: ScheduledReport | null
  currentExport: DataExport | null
  reportHistory: ReportExecution[]
  templates: ReportTemplate[]
  insights: ReportInsights | null
  isLoading: boolean
  isCreating: boolean
  isExecuting: boolean
  error: ReportingError | null
}

export interface ReportingContextType {
  state: ReportingState
  
  // Report operations
  loadScheduledReports: () => Promise<void>
  createScheduledReport: (request: CreateScheduledReportRequest) => Promise<ScheduledReport>
  updateScheduledReport: (id: string, request: UpdateScheduledReportRequest) => Promise<ScheduledReport>
  deleteScheduledReport: (id: string) => Promise<void>
  executeReport: (id: string) => Promise<void>
  getReportHistory: (id: string) => Promise<void>
  
  // Export operations
  loadDataExports: () => Promise<void>
  createDataExport: (request: DataExportRequest) => Promise<DataExport>
  downloadExport: (id: string) => Promise<void>
  
  // Template operations
  loadTemplates: () => Promise<void>
  createReportFromTemplate: (templateId: string, customization: any) => Promise<ScheduledReport>
  
  // Insights
  loadInsights: () => Promise<void>
  
  // UI state
  clearError: () => void
  setCurrentReport: (report: ScheduledReport | null) => void
  setCurrentExport: (exportData: DataExport | null) => void
}