import React, { createContext, useContext, useReducer, useCallback, useEffect } from 'react'
import {
  ScheduledReport,
  DataExport,
  ReportExecution,
  ReportTemplate,
  ReportInsights,
  CreateScheduledReportRequest,
  UpdateScheduledReportRequest,
  DataExportRequest,
  ReportingState,
  ReportingContextType,
  ReportingError
} from '../types/reporting'
import { reportingService } from '../services/reporting'

// Actions
type ReportingAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_CREATING'; payload: boolean }
  | { type: 'SET_EXECUTING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: ReportingError | null }
  | { type: 'SET_SCHEDULED_REPORTS'; payload: ScheduledReport[] }
  | { type: 'ADD_SCHEDULED_REPORT'; payload: ScheduledReport }
  | { type: 'UPDATE_SCHEDULED_REPORT'; payload: ScheduledReport }
  | { type: 'REMOVE_SCHEDULED_REPORT'; payload: string }
  | { type: 'SET_CURRENT_REPORT'; payload: ScheduledReport | null }
  | { type: 'SET_DATA_EXPORTS'; payload: DataExport[] }
  | { type: 'ADD_DATA_EXPORT'; payload: DataExport }
  | { type: 'UPDATE_DATA_EXPORT'; payload: DataExport }
  | { type: 'SET_CURRENT_EXPORT'; payload: DataExport | null }
  | { type: 'SET_REPORT_HISTORY'; payload: ReportExecution[] }
  | { type: 'SET_TEMPLATES'; payload: ReportTemplate[] }
  | { type: 'SET_INSIGHTS'; payload: ReportInsights | null }

// Initial state
const initialState: ReportingState = {
  scheduledReports: [],
  dataExports: [],
  currentReport: null,
  currentExport: null,
  reportHistory: [],
  templates: [],
  insights: null,
  isLoading: false,
  isCreating: false,
  isExecuting: false,
  error: null
}

// Reducer
function reportingReducer(state: ReportingState, action: ReportingAction): ReportingState {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload }
    
    case 'SET_CREATING':
      return { ...state, isCreating: action.payload }
    
    case 'SET_EXECUTING':
      return { ...state, isExecuting: action.payload }
    
    case 'SET_ERROR':
      return { ...state, error: action.payload, isLoading: false, isCreating: false, isExecuting: false }
    
    case 'SET_SCHEDULED_REPORTS':
      return { ...state, scheduledReports: action.payload }
    
    case 'ADD_SCHEDULED_REPORT':
      return { ...state, scheduledReports: [...state.scheduledReports, action.payload] }
    
    case 'UPDATE_SCHEDULED_REPORT':
      return {
        ...state,
        scheduledReports: state.scheduledReports.map(report =>
          report.id === action.payload.id ? action.payload : report
        ),
        currentReport: state.currentReport?.id === action.payload.id ? action.payload : state.currentReport
      }
    
    case 'REMOVE_SCHEDULED_REPORT':
      return {
        ...state,
        scheduledReports: state.scheduledReports.filter(report => report.id !== action.payload),
        currentReport: state.currentReport?.id === action.payload ? null : state.currentReport
      }
    
    case 'SET_CURRENT_REPORT':
      return { ...state, currentReport: action.payload }
    
    case 'SET_DATA_EXPORTS':
      return { ...state, dataExports: action.payload }
    
    case 'ADD_DATA_EXPORT':
      return { ...state, dataExports: [...state.dataExports, action.payload] }
    
    case 'UPDATE_DATA_EXPORT':
      return {
        ...state,
        dataExports: state.dataExports.map(exp =>
          exp.id === action.payload.id ? action.payload : exp
        ),
        currentExport: state.currentExport?.id === action.payload.id ? action.payload : state.currentExport
      }
    
    case 'SET_CURRENT_EXPORT':
      return { ...state, currentExport: action.payload }
    
    case 'SET_REPORT_HISTORY':
      return { ...state, reportHistory: action.payload }
    
    case 'SET_TEMPLATES':
      return { ...state, templates: action.payload }
    
    case 'SET_INSIGHTS':
      return { ...state, insights: action.payload }
    
    default:
      return state
  }
}

// Context
const ReportingContext = createContext<ReportingContextType | undefined>(undefined)

// Provider component
export function ReportingProvider({ children }: { children: React.ReactNode }) {
  const [state, dispatch] = useReducer(reportingReducer, initialState)

  // Load initial data
  useEffect(() => {
    loadScheduledReports()
    loadDataExports()
    loadTemplates()
    loadInsights()
  }, [])

  // Scheduled Reports operations
  const loadScheduledReports = useCallback(async () => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true })
      const response = await reportingService.getScheduledReports()
      dispatch({ type: 'SET_SCHEDULED_REPORTS', payload: response.reports })
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'report',
          message: 'Failed to load scheduled reports',
          details: error
        }
      })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [])

  const createScheduledReport = useCallback(async (request: CreateScheduledReportRequest): Promise<ScheduledReport> => {
    try {
      dispatch({ type: 'SET_CREATING', payload: true })
      const report = await reportingService.createScheduledReport(request)
      dispatch({ type: 'ADD_SCHEDULED_REPORT', payload: report })
      return report
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'report',
          message: 'Failed to create scheduled report',
          details: error
        }
      })
      throw error
    } finally {
      dispatch({ type: 'SET_CREATING', payload: false })
    }
  }, [])

  const updateScheduledReport = useCallback(async (
    id: string,
    request: UpdateScheduledReportRequest
  ): Promise<ScheduledReport> => {
    try {
      const report = await reportingService.updateScheduledReport(id, request)
      dispatch({ type: 'UPDATE_SCHEDULED_REPORT', payload: report })
      return report
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'report',
          message: 'Failed to update scheduled report',
          details: error
        }
      })
      throw error
    }
  }, [])

  const deleteScheduledReport = useCallback(async (id: string): Promise<void> => {
    try {
      await reportingService.deleteScheduledReport(id)
      dispatch({ type: 'REMOVE_SCHEDULED_REPORT', payload: id })
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'report',
          message: 'Failed to delete scheduled report',
          details: error
        }
      })
      throw error
    }
  }, [])

  const executeReport = useCallback(async (id: string): Promise<void> => {
    try {
      dispatch({ type: 'SET_EXECUTING', payload: true })
      await reportingService.executeReport(id)
      
      // Refresh report history after execution
      setTimeout(() => {
        getReportHistory(id)
      }, 1000)
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'execution',
          message: 'Failed to execute report',
          details: error
        }
      })
    } finally {
      dispatch({ type: 'SET_EXECUTING', payload: false })
    }
  }, [])

  const getReportHistory = useCallback(async (id: string): Promise<void> => {
    try {
      const response = await reportingService.getReportHistory(id)
      dispatch({ type: 'SET_REPORT_HISTORY', payload: response.executions })
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'report',
          message: 'Failed to load report history',
          details: error
        }
      })
    }
  }, [])

  // Data Export operations
  const loadDataExports = useCallback(async () => {
    try {
      const response = await reportingService.getDataExports()
      dispatch({ type: 'SET_DATA_EXPORTS', payload: response.exports })
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'export',
          message: 'Failed to load data exports',
          details: error
        }
      })
    }
  }, [])

  const createDataExport = useCallback(async (request: DataExportRequest): Promise<DataExport> => {
    try {
      dispatch({ type: 'SET_CREATING', payload: true })
      const dataExport = await reportingService.createDataExport(request)
      dispatch({ type: 'ADD_DATA_EXPORT', payload: dataExport })
      return dataExport
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'export',
          message: 'Failed to create data export',
          details: error
        }
      })
      throw error
    } finally {
      dispatch({ type: 'SET_CREATING', payload: false })
    }
  }, [])

  const downloadExport = useCallback(async (id: string): Promise<void> => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true })
      const blob = await reportingService.downloadExport(id)
      
      // Get export details to determine filename
      const exportData = await reportingService.getDataExport(id)
      const filename = `${exportData.export_type}_export_${exportData.id}.${exportData.format}`
      
      await reportingService.downloadFile(blob, filename)
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'export',
          message: 'Failed to download export',
          details: error
        }
      })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [])

  // Template operations
  const loadTemplates = useCallback(async () => {
    try {
      const response = await reportingService.getReportTemplates()
      dispatch({ type: 'SET_TEMPLATES', payload: response.templates })
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'template',
          message: 'Failed to load report templates',
          details: error
        }
      })
    }
  }, [])

  const createReportFromTemplate = useCallback(async (
    templateId: string,
    customization: any
  ): Promise<ScheduledReport> => {
    try {
      dispatch({ type: 'SET_CREATING', payload: true })
      const response = await reportingService.createReportFromTemplate(templateId, customization)
      dispatch({ type: 'ADD_SCHEDULED_REPORT', payload: response.report })
      return response.report
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'template',
          message: 'Failed to create report from template',
          details: error
        }
      })
      throw error
    } finally {
      dispatch({ type: 'SET_CREATING', payload: false })
    }
  }, [])

  // Insights
  const loadInsights = useCallback(async () => {
    try {
      const insights = await reportingService.getReportInsights()
      dispatch({ type: 'SET_INSIGHTS', payload: insights })
    } catch (error) {
      dispatch({
        type: 'SET_ERROR',
        payload: {
          type: 'report',
          message: 'Failed to load reporting insights',
          details: error
        }
      })
    }
  }, [])

  // UI helpers
  const clearError = useCallback(() => {
    dispatch({ type: 'SET_ERROR', payload: null })
  }, [])

  const setCurrentReport = useCallback((report: ScheduledReport | null) => {
    dispatch({ type: 'SET_CURRENT_REPORT', payload: report })
  }, [])

  const setCurrentExport = useCallback((dataExport: DataExport | null) => {
    dispatch({ type: 'SET_CURRENT_EXPORT', payload: dataExport })
  }, [])

  // Context value
  const contextValue: ReportingContextType = {
    state,
    loadScheduledReports,
    createScheduledReport,
    updateScheduledReport,
    deleteScheduledReport,
    executeReport,
    getReportHistory,
    loadDataExports,
    createDataExport,
    downloadExport,
    loadTemplates,
    createReportFromTemplate,
    loadInsights,
    clearError,
    setCurrentReport,
    setCurrentExport
  }

  return (
    <ReportingContext.Provider value={contextValue}>
      {children}
    </ReportingContext.Provider>
  )
}

// Hook to use reporting context
export function useReporting(): ReportingContextType {
  const context = useContext(ReportingContext)
  if (context === undefined) {
    throw new Error('useReporting must be used within a ReportingProvider')
  }
  return context
}

export default ReportingContext