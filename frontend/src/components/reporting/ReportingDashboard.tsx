import React, { useState } from 'react'
import { useReporting } from '../../context/ReportingContext'
import ScheduledReportsList from './ScheduledReportsList'
import CreateReportModal from './CreateReportModal'
import DataExportsList from './DataExportsList'
import ReportTemplates from './ReportTemplates'
import ReportingInsights from './ReportingInsights'

interface ReportingDashboardProps {
  className?: string
}

type TabType = 'scheduled' | 'exports' | 'templates' | 'insights'

const ReportingDashboard: React.FC<ReportingDashboardProps> = ({ className = '' }) => {
  const { state, clearError } = useReporting()
  const [activeTab, setActiveTab] = useState<TabType>('scheduled')
  const [showCreateModal, setShowCreateModal] = useState(false)

  const tabs = [
    { id: 'scheduled' as TabType, label: 'Scheduled Reports', icon: '📅' },
    { id: 'exports' as TabType, label: 'Data Exports', icon: '📊' },
    { id: 'templates' as TabType, label: 'Templates', icon: '📋' },
    { id: 'insights' as TabType, label: 'Insights', icon: '💡' }
  ]

  const handleCreateReport = () => {
    setShowCreateModal(true)
  }

  const handleCloseCreateModal = () => {
    setShowCreateModal(false)
  }

  const renderTabContent = () => {
    switch (activeTab) {
      case 'scheduled':
        return <ScheduledReportsList onCreateReport={handleCreateReport} />
      case 'exports':
        return <DataExportsList />
      case 'templates':
        return <ReportTemplates />
      case 'insights':
        return <ReportingInsights />
      default:
        return <ScheduledReportsList onCreateReport={handleCreateReport} />
    }
  }

  return (
    <div className={`reporting-dashboard ${className}`}>
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Reporting & Analytics
        </h1>
        <p className="text-gray-600">
          Create, schedule, and manage your reports and data exports
        </p>
      </div>

      {/* Error Alert */}
      {state.error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-start justify-between">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">
                  {state.error.message}
                </h3>
                {state.error.details && (
                  <div className="mt-2 text-sm text-red-700">
                    <p>{JSON.stringify(state.error.details, null, 2)}</p>
                  </div>
                )}
              </div>
            </div>
            <div className="ml-auto pl-3">
              <div className="-mx-1.5 -my-1.5">
                <button
                  onClick={clearError}
                  className="inline-flex bg-red-50 rounded-md p-1.5 text-red-500 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-red-50 focus:ring-red-600"
                >
                  <span className="sr-only">Dismiss</span>
                  <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Loading Overlay */}
      {state.isLoading && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <div className="mt-3 text-center">
              <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-blue-100">
                <svg className="animate-spin h-6 w-6 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <h3 className="text-lg leading-6 font-medium text-gray-900 mt-4">Loading...</h3>
              <div className="mt-2 px-7 py-3">
                <p className="text-sm text-gray-500">
                  Please wait while we load your reporting data.
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Navigation Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <nav className="-mb-px flex space-x-8">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`${
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm flex items-center space-x-2`}
            >
              <span>{tab.icon}</span>
              <span>{tab.label}</span>
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="tab-content">
        {renderTabContent()}
      </div>

      {/* Create Report Modal */}
      {showCreateModal && (
        <CreateReportModal
          isOpen={showCreateModal}
          onClose={handleCloseCreateModal}
        />
      )}

      {/* Quick Actions Floating Button (Mobile) */}
      <div className="fixed bottom-6 right-6 md:hidden">
        <button
          onClick={handleCreateReport}
          className="bg-blue-600 hover:bg-blue-700 text-white rounded-full p-4 shadow-lg focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
    </div>
  )
}

export default ReportingDashboard