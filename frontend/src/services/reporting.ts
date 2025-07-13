import axios, { AxiosResponse } from 'axios'
import {
  ScheduledReport,
  CreateScheduledReportRequest,
  UpdateScheduledReportRequest,
  ScheduledReportsResponse,
  ReportExecution,
  ReportHistoryResponse,
  DataExport,
  DataExportRequest,
  DataExportsResponse,
  ReportTemplate,
  ReportTemplatesResponse,
  ReportInsights,
  JobStatus
} from '../types/reporting'

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080'

class ReportingService {
  private baseURL: string

  constructor() {
    this.baseURL = `${API_BASE_URL}/api/v1/reports`
  }

  private getAuthHeaders() {
    const token = localStorage.getItem('access_token')
    return {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  }

  // Scheduled Reports
  async getScheduledReports(): Promise<ScheduledReportsResponse> {
    const response: AxiosResponse<ScheduledReportsResponse> = await axios.get(
      `${this.baseURL}/scheduled`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async createScheduledReport(request: CreateScheduledReportRequest): Promise<ScheduledReport> {
    const response: AxiosResponse<ScheduledReport> = await axios.post(
      `${this.baseURL}/scheduled`,
      request,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async getScheduledReport(id: string): Promise<ScheduledReport> {
    const response: AxiosResponse<ScheduledReport> = await axios.get(
      `${this.baseURL}/scheduled/${id}`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async updateScheduledReport(id: string, request: UpdateScheduledReportRequest): Promise<ScheduledReport> {
    const response: AxiosResponse<ScheduledReport> = await axios.put(
      `${this.baseURL}/scheduled/${id}`,
      request,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async deleteScheduledReport(id: string): Promise<void> {
    await axios.delete(
      `${this.baseURL}/scheduled/${id}`,
      { headers: this.getAuthHeaders() }
    )
  }

  async executeReport(id: string): Promise<{ message: string; report_id: string; status: string }> {
    const response = await axios.post(
      `${this.baseURL}/scheduled/${id}/execute`,
      {},
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async getReportHistory(id: string): Promise<ReportHistoryResponse> {
    const response: AxiosResponse<ReportHistoryResponse> = await axios.get(
      `${this.baseURL}/scheduled/${id}/history`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  // Data Exports
  async createDataExport(request: DataExportRequest): Promise<DataExport> {
    const response: AxiosResponse<DataExport> = await axios.post(
      `${this.baseURL}/exports`,
      request,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async getDataExports(): Promise<DataExportsResponse> {
    const response: AxiosResponse<DataExportsResponse> = await axios.get(
      `${this.baseURL}/exports`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async getDataExport(id: string): Promise<DataExport> {
    const response: AxiosResponse<DataExport> = await axios.get(
      `${this.baseURL}/exports/${id}`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async downloadExport(id: string): Promise<Blob> {
    const response = await axios.get(
      `${this.baseURL}/exports/${id}/download`,
      {
        headers: this.getAuthHeaders(),
        responseType: 'blob'
      }
    )
    return response.data
  }

  // Report Generation
  async generateAnalyticsReport(config: any): Promise<Blob> {
    const response = await axios.post(
      `${this.baseURL}/analytics`,
      config,
      {
        headers: this.getAuthHeaders(),
        responseType: 'blob'
      }
    )
    return response.data
  }

  async generateDashboardReport(dashboardId: string, format = 'json'): Promise<Blob> {
    const response = await axios.get(
      `${this.baseURL}/dashboards/${dashboardId}?format=${format}`,
      {
        headers: this.getAuthHeaders(),
        responseType: 'blob'
      }
    )
    return response.data
  }

  async generateFunnelReport(funnelId: string, format = 'json'): Promise<Blob> {
    const response = await axios.get(
      `${this.baseURL}/funnels/${funnelId}?format=${format}`,
      {
        headers: this.getAuthHeaders(),
        responseType: 'blob'
      }
    )
    return response.data
  }

  // Templates
  async getReportTemplates(): Promise<ReportTemplatesResponse> {
    const response: AxiosResponse<ReportTemplatesResponse> = await axios.get(
      `${this.baseURL}/templates`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  async createReportFromTemplate(
    templateId: string,
    customization: {
      name: string
      description?: string
      recipients: string[]
      schedule?: string
    }
  ): Promise<{ report: ScheduledReport; template_id: string; message: string }> {
    const response = await axios.post(
      `${this.baseURL}/templates/${templateId}`,
      customization,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  // Insights
  async getReportInsights(): Promise<ReportInsights> {
    const response: AxiosResponse<ReportInsights> = await axios.get(
      `${this.baseURL}/insights`,
      { headers: this.getAuthHeaders() }
    )
    return response.data
  }

  // Utility methods
  async downloadFile(blob: Blob, filename: string): Promise<void> {
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }

  formatReportSchedule(cronExpression: string): string {
    const scheduleMap: Record<string, string> = {
      '0 8 * * *': 'Daily at 8:00 AM',
      '0 9 * * MON': 'Weekly on Monday at 9:00 AM',
      '0 10 1 * *': 'Monthly on 1st at 10:00 AM',
      '0 0 * * 0': 'Weekly on Sunday at midnight',
      '0 6 * * 1-5': 'Weekdays at 6:00 AM'
    }
    
    return scheduleMap[cronExpression] || cronExpression
  }

  validateCronExpression(cronExpression: string): boolean {
    // Basic cron validation - in a real app you'd use a proper cron parser
    const cronRegex = /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]|1[0-9]|2[0-3])) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\*|([1-9]|1[0-2])|\*\/([1-9]|1[0-2])) (\*|([0-6])|\*\/([0-6]))$/
    return cronRegex.test(cronExpression)
  }

  getNextRunTime(cronExpression: string): Date | null {
    // Simple implementation - in production use a proper cron parser
    const now = new Date()
    
    switch (cronExpression) {
      case '0 8 * * *': // Daily at 8 AM
        const tomorrow8am = new Date(now)
        tomorrow8am.setDate(now.getDate() + 1)
        tomorrow8am.setHours(8, 0, 0, 0)
        return tomorrow8am
        
      case '0 9 * * MON': // Weekly on Monday at 9 AM
        const nextMonday = new Date(now)
        const daysUntilMonday = (1 + 7 - now.getDay()) % 7 || 7
        nextMonday.setDate(now.getDate() + daysUntilMonday)
        nextMonday.setHours(9, 0, 0, 0)
        return nextMonday
        
      case '0 10 1 * *': // Monthly on 1st at 10 AM
        const nextMonth = new Date(now)
        nextMonth.setMonth(now.getMonth() + 1, 1)
        nextMonth.setHours(10, 0, 0, 0)
        return nextMonth
        
      default:
        return null
    }
  }

  formatFileSize(bytes: number): string {
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    if (bytes === 0) return '0 Bytes'
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
  }

  formatDuration(milliseconds: number): string {
    const seconds = Math.floor(milliseconds / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    
    if (hours > 0) {
      return `${hours}h ${minutes % 60}m ${seconds % 60}s`
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`
    } else {
      return `${seconds}s`
    }
  }
}

export const reportingService = new ReportingService()
export default reportingService