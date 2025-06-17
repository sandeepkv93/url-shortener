const Analytics = () => {
  return (
    <div className="px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Analytics</h1>
        <p className="mt-2 text-gray-600">Detailed insights into your URL performance.</p>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2 mb-8">
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Click Timeline</h3>
          <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
            <p className="text-gray-500">Chart will be implemented in Step 18</p>
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Geographic Distribution</h3>
          <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
            <p className="text-gray-500">Map will be implemented in Step 18</p>
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Device Stats</h3>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Desktop</span>
              <span className="text-sm font-medium text-gray-900">0%</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Mobile</span>
              <span className="text-sm font-medium text-gray-900">0%</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Tablet</span>
              <span className="text-sm font-medium text-gray-900">0%</span>
            </div>
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Top Referrers</h3>
          <div className="space-y-4">
            <div className="text-center py-8">
              <p className="text-gray-500">No data available yet</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Analytics