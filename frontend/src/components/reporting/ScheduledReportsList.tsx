import React, { useState } from 'react'
import { useReporting } from '../../context/ReportingContext'
import { ScheduledReport } from '../../types/reporting'
import { reportingService } from '../../services/reporting'

interface ScheduledReportsListProps {
  onCreateReport: () => void
}

const ScheduledReportsList: React.FC<ScheduledReportsListProps> = ({ onCreateReport }) => {
  const { state, deleteScheduledReport, executeReport, getReportHistory, setCurrentReport } = useReporting()
  const [selectedReport, setSelectedReport] = useState<ScheduledReport | null>(null)
  const [showHistory, setShowHistory] = useState(false)

  const handleDeleteReport = async (reportId: string) => {
    if (window.confirm('Are you sure you want to delete this report? This action cannot be undone.')) {
      try {
        await deleteScheduledReport(reportId)
      } catch (error) {
        console.error('Failed to delete report:', error)
      }
    }
  }

  const handleExecuteReport = async (reportId: string) => {
    try {
      await executeReport(reportId)
      alert('Report execution started. You will receive an email when it\'s complete.')
    } catch (error) {
      console.error('Failed to execute report:', error)
    }
  }

  const handleViewHistory = (report: ScheduledReport) => {
    setSelectedReport(report)
    setCurrentReport(report)
    getReportHistory(report.id)
    setShowHistory(true)
  }

  const getStatusBadge = (isActive: boolean, nextRun?: string) => {
    if (!isActive) {
      return (
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
          Inactive
        </span>
      )
    }

    const now = new Date()
    const next = nextRun ? new Date(nextRun) : null
    
    if (next && next > now) {
      const timeDiff = next.getTime() - now.getTime()
      const hours = Math.floor(timeDiff / (1000 * 60 * 60))
      
      if (hours < 24) {
        return (
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
            Next: {hours}h
          </span>
        )
      }
    }

    return (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
        Active
      </span>
    )
  }

  const formatReportType = (type: string) => {
    const typeMap: Record<string, string> = {
      analytics: 'Analytics',
      dashboard: 'Dashboard',
      funnel: 'Funnel',
      competitive: 'Competitive'
    }
    return typeMap[type] || type
  }

  const formatSchedule = (schedule: string) => {
    return reportingService.formatReportSchedule(schedule)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (state.scheduledReports.length === 0) {
    return (
      <div className="text-center py-12">
        <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">No scheduled reports</h3>
        <p className="mt-1 text-sm text-gray-500">
          Get started by creating your first scheduled report.
        </p>
        <div className="mt-6">
          <button
            onClick={onCreateReport}
            className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Create Report
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="scheduled-reports-list">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-lg font-medium text-gray-900">Scheduled Reports</h2>
          <p className="text-sm text-gray-500">
            {state.scheduledReports.length} report{state.scheduledReports.length !== 1 ? 's' : ''} configured
          </p>
        </div>
        <button
          onClick={onCreateReport}
          disabled={state.isCreating}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {state.isCreating ? (
            <>
              <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Creating...
            </>
          ) : (
            <>
              <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Create Report
            </>
          )}
        </button>
      </div>

      {/* Reports Grid */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {state.scheduledReports.map((report) => (
          <div key={report.id} className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-6">
              {/* Header */}
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div className="flex-shrink-0">
                    <div className="h-10 w-10 rounded-lg bg-blue-100 flex items-center justify-center">
                      <svg className="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                      </svg>
                    </div>
                  </div>
                  <div>
                    <h3 className="text-lg font-medium text-gray-900">{report.name}</h3>
                    <p className="text-sm text-gray-500">{formatReportType(report.report_type)}</p>
                  </div>
                </div>
                {getStatusBadge(report.is_active, report.next_run)}
              </div>

              {/* Description */}
              {report.description && (
                <p className="text-sm text-gray-600 mb-4">{report.description}</p>
              )}

              {/* Details */}
              <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Schedule:</span>
                  <span className="text-gray-900">{formatSchedule(report.schedule)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Format:</span>
                  <span className="text-gray-900 uppercase">{report.format}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Recipients:</span>
                  <span className="text-gray-900">{report.recipients.length} recipient{report.recipients.length !== 1 ? 's' : ''}</span>
                </div>
                {report.last_run && (
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Last run:</span>
                    <span className="text-gray-900">{formatDate(report.last_run)}</span>
                  </div>
                )}
                {report.next_run && (
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Next run:</span>
                    <span className="text-gray-900">{formatDate(report.next_run)}</span>
                  </div>
                )}
              </div>

              {/* Actions */}
              <div className="flex space-x-2">
                <button
                  onClick={() => handleExecuteReport(report.id)}
                  disabled={state.isExecuting}
                  className="flex-1 inline-flex justify-center items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                >
                  {state.isExecuting ? (
                    <svg className="animate-spin h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  ) : (
                    <svg className="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h1m4 0h1m-6-8h1m4 0h1M8 16h1m4 0h1" />
                    </svg>
                  )}
                  <span className="ml-2">Run Now</span>
                </button>
                
                <button
                  onClick={() => handleViewHistory(report)}
                  className="px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  History
                </button>
                
                <button
                  onClick={() => handleDeleteReport(report.id)}
                  className="px-3 py-2 border border-red-300 shadow-sm text-sm leading-4 font-medium rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Report History Modal */}
      {showHistory && selectedReport && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-11/12 max-w-4xl shadow-lg rounded-md bg-white">
            <div className="mt-3">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                  Report History: {selectedReport.name}
                </h3>
                <button
                  onClick={() => setShowHistory(false)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              
              {state.reportHistory.length === 0 ? (
                <p className="text-gray-500 text-center py-8">No execution history found.</p>
              ) : (
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Execution
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Status
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Duration
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          File Size
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Started
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {state.reportHistory.map((execution) => (
                        <tr key={execution.id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            #{execution.id}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                              execution.status === 'completed' ? 'bg-green-100 text-green-800' :
                              execution.status === 'failed' ? 'bg-red-100 text-red-800' :
                              execution.status === 'running' ? 'bg-blue-100 text-blue-800' :
                              'bg-gray-100 text-gray-800'
                            }`}>
                              {execution.status}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {execution.duration ? reportingService.formatDuration(execution.duration) : '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {execution.file_size ? reportingService.formatFileSize(execution.file_size) : '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {formatDate(execution.started_at)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default ScheduledReportsList